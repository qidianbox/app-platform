#!/usr/bin/env python3
"""
自动视觉UI测试脚本
覆盖所有页面的截图和视觉回归测试
修复版本：支持登录状态持久化和更多页面测试
"""

import os
import json
import time
import hashlib
import shutil
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from PIL import Image, ImageChops, ImageDraw, ImageFont
import math

# 配置
BASE_URL = "http://localhost:5173"
RESULTS_DIR = "/home/ubuntu/app-platform/test_results/visual_test"
BASELINE_DIR = os.path.join(RESULTS_DIR, "baseline")
CURRENT_DIR = os.path.join(RESULTS_DIR, "current")
DIFF_DIR = os.path.join(RESULTS_DIR, "diff")

# 测试页面配置 - 扩展版本
TEST_PAGES = [
    # 公开页面
    {
        "name": "登录页",
        "path": "/login",
        "requires_auth": False,
        "wait_for": ".login-container, .login-form, form",
        "wait_time": 2
    },
    # 需要认证的页面
    {
        "name": "APP列表页",
        "path": "/apps",
        "requires_auth": True,
        "wait_for": ".app-list, .apps-container, table, .el-table",
        "wait_time": 3
    },
    {
        "name": "APP详情_概览",
        "path": "/apps/2/config",
        "requires_auth": True,
        "wait_for": ".page-content, .stats-cards, .overview",
        "wait_time": 3
    },
    {
        "name": "APP详情_基础配置",
        "path": "/apps/2/config",
        "requires_auth": True,
        "wait_for": ".config-form, .el-form",
        "wait_time": 2,
        "action": "click_basic_config"
    },
    {
        "name": "APP详情_工作台",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".workspace-content, .workspace-container",
        "wait_time": 3
    },
    # 工作台子页面 - 新增
    {
        "name": "工作台_监控告警",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".monitor-content, .alert-list",
        "wait_time": 3,
        "action": "click_monitor"
    },
    {
        "name": "工作台_审计日志",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".audit-content, .audit-log",
        "wait_time": 3,
        "action": "click_audit"
    },
    {
        "name": "工作台_消息推送",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".push-content, .message-list",
        "wait_time": 3,
        "action": "click_push"
    },
    {
        "name": "工作台_版本管理",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".version-content, .version-list",
        "wait_time": 3,
        "action": "click_version"
    },
    {
        "name": "工作台_用户管理",
        "path": "/apps/2/workspace",
        "requires_auth": True,
        "wait_for": ".user-content, .user-list",
        "wait_time": 3,
        "action": "click_user"
    }
]

# 设备配置
DEVICES = [
    {"name": "Desktop_1920x1080", "width": 1920, "height": 1080},
    {"name": "Laptop_1366x768", "width": 1366, "height": 768},
    {"name": "Tablet_768x1024", "width": 768, "height": 1024},
    {"name": "Mobile_375x667", "width": 375, "height": 667}
]

class VisualUITester:
    def __init__(self):
        self.driver = None
        self.results = []
        self.timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        self.is_logged_in = False
        self.login_cookies = None
        
        # 创建目录
        for dir_path in [RESULTS_DIR, BASELINE_DIR, CURRENT_DIR, DIFF_DIR]:
            os.makedirs(dir_path, exist_ok=True)
    
    def setup_driver(self, width, height):
        """设置浏览器驱动"""
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument(f"--window-size={width},{height}")
        chrome_options.add_argument("--force-device-scale-factor=1")
        # 添加用户数据目录以保持登录状态
        chrome_options.add_argument("--user-data-dir=/tmp/chrome_visual_test")
        
        service = Service("/usr/bin/chromedriver")
        self.driver = webdriver.Chrome(service=service, options=chrome_options)
        self.driver.set_window_size(width, height)
        
        # 恢复登录Cookie
        if self.login_cookies:
            self.driver.get(BASE_URL)
            time.sleep(1)
            for cookie in self.login_cookies:
                try:
                    self.driver.add_cookie(cookie)
                except Exception as e:
                    print(f"恢复Cookie失败: {e}")
    
    def login(self):
        """执行登录并保存Cookie"""
        try:
            print("正在执行登录...")
            self.driver.get(f"{BASE_URL}/login")
            time.sleep(3)
            
            # 检查是否已经登录（通过URL或页面元素判断）
            current_url = self.driver.current_url
            if "/apps" in current_url or "/dashboard" in current_url:
                print("已经处于登录状态")
                self.is_logged_in = True
                self.login_cookies = self.driver.get_cookies()
                return True
            
            # 查找并填写登录表单
            try:
                # 尝试多种选择器
                username_selectors = [
                    "input[type='text']",
                    "input[placeholder*='用户']",
                    "input[name='username']",
                    ".el-input__inner"
                ]
                
                username_input = None
                for selector in username_selectors:
                    try:
                        inputs = self.driver.find_elements(By.CSS_SELECTOR, selector)
                        for inp in inputs:
                            if inp.is_displayed() and inp.get_attribute("type") != "password":
                                username_input = inp
                                break
                        if username_input:
                            break
                    except:
                        continue
                
                if not username_input:
                    print("找不到用户名输入框")
                    return False
                
                password_input = self.driver.find_element(By.CSS_SELECTOR, "input[type='password']")
                
                # 清空并输入
                username_input.clear()
                time.sleep(0.5)
                username_input.send_keys("admin")
                time.sleep(0.5)
                
                password_input.clear()
                time.sleep(0.5)
                password_input.send_keys("admin123")
                time.sleep(0.5)
                
                # 点击登录按钮
                login_btn_selectors = [
                    "button[type='submit']",
                    ".login-btn",
                    "button.el-button--primary",
                    ".el-button--primary"
                ]
                
                login_btn = None
                for selector in login_btn_selectors:
                    try:
                        btn = self.driver.find_element(By.CSS_SELECTOR, selector)
                        if btn.is_displayed():
                            login_btn = btn
                            break
                    except:
                        continue
                
                if login_btn:
                    login_btn.click()
                else:
                    # 尝试通过文本查找
                    login_btn = self.driver.find_element(By.XPATH, "//button[contains(text(), '登') or contains(text(), 'Login')]")
                    login_btn.click()
                
                # 等待登录完成
                time.sleep(4)
                
                # 验证登录成功
                current_url = self.driver.current_url
                print(f"登录后URL: {current_url}")
                
                if "/login" not in current_url or "/apps" in current_url:
                    print("登录成功!")
                    self.is_logged_in = True
                    self.login_cookies = self.driver.get_cookies()
                    print(f"保存了 {len(self.login_cookies)} 个Cookie")
                    return True
                else:
                    print("登录可能失败，仍在登录页")
                    # 尝试手动跳转
                    self.driver.get(f"{BASE_URL}/apps")
                    time.sleep(2)
                    if "/login" not in self.driver.current_url:
                        self.is_logged_in = True
                        self.login_cookies = self.driver.get_cookies()
                        return True
                    return False
                    
            except Exception as e:
                print(f"登录过程出错: {e}")
                return False
                
        except Exception as e:
            print(f"登录失败: {e}")
            return False
    
    def wait_for_element(self, selector, timeout=10):
        """等待元素出现"""
        try:
            selectors = selector.split(", ")
            for sel in selectors:
                try:
                    WebDriverWait(self.driver, timeout).until(
                        EC.presence_of_element_located((By.CSS_SELECTOR, sel.strip()))
                    )
                    return True
                except:
                    continue
            return False
        except:
            return False
    
    def execute_action(self, action):
        """执行页面操作"""
        try:
            time.sleep(1)
            
            if action == "click_basic_config":
                # 点击基础配置菜单
                selectors = [
                    "//*[contains(text(), '基础配置')]",
                    "//span[contains(text(), '基础配置')]",
                    "//*[@class='menu-item' and contains(., '基础配置')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_workspace":
                # 点击工作台Tab
                selectors = [
                    "//*[contains(text(), '工作台')]",
                    "//span[contains(text(), '工作台')]",
                    "//*[@class='el-tabs__item' and contains(., '工作台')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_monitor":
                # 点击监控告警菜单
                selectors = [
                    "//*[contains(text(), '监控告警')]",
                    "//span[contains(text(), '监控告警')]",
                    "//*[@class='menu-item' and contains(., '监控')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_audit":
                # 点击审计日志菜单
                selectors = [
                    "//*[contains(text(), '审计日志')]",
                    "//span[contains(text(), '审计日志')]",
                    "//*[@class='menu-item' and contains(., '审计')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_push":
                # 点击消息推送菜单
                selectors = [
                    "//*[contains(text(), '消息推送')]",
                    "//span[contains(text(), '消息推送')]",
                    "//*[@class='menu-item' and contains(., '推送')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_version":
                # 点击版本管理菜单
                selectors = [
                    "//*[contains(text(), '版本管理')]",
                    "//span[contains(text(), '版本管理')]",
                    "//*[@class='menu-item' and contains(., '版本')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
                        
            elif action == "click_user":
                # 点击用户管理菜单
                selectors = [
                    "//*[contains(text(), '用户管理')]",
                    "//span[contains(text(), '用户管理')]",
                    "//*[@class='menu-item' and contains(., '用户')]"
                ]
                for xpath in selectors:
                    try:
                        elem = self.driver.find_element(By.XPATH, xpath)
                        if elem.is_displayed():
                            elem.click()
                            time.sleep(2)
                            return True
                    except:
                        continue
            
            print(f"执行操作: {action}")
            return True
            
        except Exception as e:
            print(f"执行操作失败: {action}, {e}")
            return False
    
    def take_screenshot(self, name, device_name):
        """截图"""
        # 清理文件名中的特殊字符
        safe_name = name.replace(" ", "_").replace("-", "_").replace("/", "_")
        filename = f"{safe_name}_{device_name}_{self.timestamp}.png"
        filepath = os.path.join(CURRENT_DIR, filename)
        self.driver.save_screenshot(filepath)
        print(f"截图保存: {filename}")
        return filepath
    
    def calculate_image_diff(self, img1_path, img2_path):
        """计算两张图片的差异"""
        try:
            img1 = Image.open(img1_path).convert('RGB')
            img2 = Image.open(img2_path).convert('RGB')
            
            # 调整尺寸一致
            if img1.size != img2.size:
                img2 = img2.resize(img1.size, Image.Resampling.LANCZOS)
            
            # 计算差异
            diff = ImageChops.difference(img1, img2)
            
            # 计算差异百分比
            diff_pixels = 0
            total_pixels = img1.size[0] * img1.size[1]
            
            for pixel in diff.getdata():
                if pixel != (0, 0, 0):
                    diff_pixels += 1
            
            diff_percentage = (diff_pixels / total_pixels) * 100
            
            return diff_percentage, diff
        except Exception as e:
            print(f"计算图片差异失败: {e}")
            return -1, None
    
    def create_diff_image(self, img1_path, img2_path, diff_path):
        """创建差异对比图"""
        try:
            img1 = Image.open(img1_path).convert('RGB')
            img2 = Image.open(img2_path).convert('RGB')
            
            if img1.size != img2.size:
                img2 = img2.resize(img1.size, Image.Resampling.LANCZOS)
            
            # 创建差异图
            diff = ImageChops.difference(img1, img2)
            
            # 增强差异可见度
            diff = diff.point(lambda x: min(255, x * 10))
            
            # 创建并排对比图
            width = img1.size[0] * 3
            height = img1.size[1]
            comparison = Image.new('RGB', (width, height))
            
            comparison.paste(img1, (0, 0))
            comparison.paste(img2, (img1.size[0], 0))
            comparison.paste(diff, (img1.size[0] * 2, 0))
            
            # 添加标签
            draw = ImageDraw.Draw(comparison)
            try:
                font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 20)
            except:
                font = ImageFont.load_default()
            
            draw.text((10, 10), "Baseline", fill="white", font=font)
            draw.text((img1.size[0] + 10, 10), "Current", fill="white", font=font)
            draw.text((img1.size[0] * 2 + 10, 10), "Diff", fill="white", font=font)
            
            comparison.save(diff_path)
            return True
        except Exception as e:
            print(f"创建差异图失败: {e}")
            return False
    
    def run_test_batch(self, pages, device):
        """批量运行测试（保持登录状态）"""
        results = []
        
        try:
            # 设置浏览器
            self.setup_driver(device["width"], device["height"])
            
            # 先执行登录
            login_success = False
            for page in pages:
                if page["requires_auth"]:
                    login_success = self.login()
                    break
            
            # 遍历所有页面
            for page in pages:
                result = {
                    "page": page["name"],
                    "device": device["name"],
                    "status": "unknown",
                    "diff_percentage": 0,
                    "screenshot": "",
                    "baseline": "",
                    "diff_image": "",
                    "error": ""
                }
                
                try:
                    print(f"\n测试页面: {page['name']} @ {device['name']}")
                    
                    # 检查认证状态
                    if page["requires_auth"] and not login_success:
                        result["status"] = "skipped"
                        result["error"] = "登录失败，跳过需要认证的页面"
                        results.append(result)
                        continue
                    
                    # 访问页面
                    self.driver.get(f"{BASE_URL}{page['path']}")
                    time.sleep(page.get("wait_time", 2))
                    
                    # 等待元素
                    self.wait_for_element(page["wait_for"], timeout=5)
                    
                    # 执行操作（如果有）
                    if "action" in page:
                        self.execute_action(page["action"])
                        time.sleep(2)
                    
                    # 额外等待确保页面完全加载
                    time.sleep(1)
                    
                    # 截图
                    screenshot_path = self.take_screenshot(page["name"], device["name"])
                    result["screenshot"] = screenshot_path
                    
                    # 查找基准图
                    safe_name = page['name'].replace(' ', '_').replace('-', '_').replace('/', '_')
                    baseline_pattern = f"{safe_name}_{device['name']}_"
                    baseline_files = [f for f in os.listdir(BASELINE_DIR) if f.startswith(baseline_pattern)]
                    
                    if baseline_files:
                        # 有基准图，进行对比
                        baseline_path = os.path.join(BASELINE_DIR, sorted(baseline_files)[-1])
                        result["baseline"] = baseline_path
                        
                        diff_percentage, _ = self.calculate_image_diff(baseline_path, screenshot_path)
                        result["diff_percentage"] = round(diff_percentage, 2)
                        
                        if diff_percentage < 0:
                            result["status"] = "error"
                            result["error"] = "图片对比失败"
                        elif diff_percentage < 1:
                            result["status"] = "passed"
                        elif diff_percentage < 5:
                            result["status"] = "warning"
                        else:
                            result["status"] = "failed"
                        
                        # 创建差异图
                        if diff_percentage > 0:
                            diff_filename = f"diff_{safe_name}_{device['name']}_{self.timestamp}.png"
                            diff_path = os.path.join(DIFF_DIR, diff_filename)
                            self.create_diff_image(baseline_path, screenshot_path, diff_path)
                            result["diff_image"] = diff_path
                    else:
                        # 没有基准图，保存为新基准
                        baseline_filename = f"{safe_name}_{device['name']}_{self.timestamp}.png"
                        baseline_path = os.path.join(BASELINE_DIR, baseline_filename)
                        shutil.copy(screenshot_path, baseline_path)
                        result["baseline"] = baseline_path
                        result["status"] = "new_baseline"
                    
                except Exception as e:
                    result["status"] = "error"
                    result["error"] = str(e)
                    print(f"测试出错: {e}")
                
                results.append(result)
                
        except Exception as e:
            print(f"批量测试出错: {e}")
        finally:
            if self.driver:
                self.driver.quit()
                self.driver = None
        
        return results
    
    def run_all_tests(self):
        """运行所有测试"""
        print("=" * 60)
        print("开始视觉UI测试")
        print(f"测试时间: {self.timestamp}")
        print(f"测试页面: {len(TEST_PAGES)}")
        print(f"测试设备: {len(DEVICES)}")
        print("=" * 60)
        
        all_results = []
        
        for device in DEVICES:
            print(f"\n{'='*40}")
            print(f"设备: {device['name']} ({device['width']}x{device['height']})")
            print("=" * 40)
            
            # 清理Chrome用户数据目录
            chrome_data_dir = "/tmp/chrome_visual_test"
            if os.path.exists(chrome_data_dir):
                shutil.rmtree(chrome_data_dir, ignore_errors=True)
            
            # 重置登录状态
            self.is_logged_in = False
            self.login_cookies = None
            
            # 批量运行测试
            results = self.run_test_batch(TEST_PAGES, device)
            all_results.extend(results)
        
        self.results = all_results
        return all_results
    
    def generate_report(self):
        """生成测试报告"""
        report_path = os.path.join(RESULTS_DIR, f"visual_test_report_{self.timestamp}.md")
        
        # 统计
        total = len(self.results)
        passed = len([r for r in self.results if r["status"] == "passed"])
        failed = len([r for r in self.results if r["status"] == "failed"])
        warning = len([r for r in self.results if r["status"] == "warning"])
        new_baseline = len([r for r in self.results if r["status"] == "new_baseline"])
        errors = len([r for r in self.results if r["status"] in ["error", "skipped"]])
        
        report = f"""# 视觉UI测试报告

**生成时间**: {datetime.now().strftime("%Y-%m-%d %H:%M:%S")}

## 测试概览

| 指标 | 数值 |
|------|------|
| 总测试数 | {total} |
| 通过 | {passed} ✅ |
| 失败 | {failed} ❌ |
| 警告 | {warning} ⚠️ |
| 新基准 | {new_baseline} 🆕 |
| 错误/跳过 | {errors} |
| 通过率 | {round(passed/total*100, 1) if total > 0 else 0}% |

## 测试页面

| 页面 | 路径 | 需要认证 |
|------|------|----------|
"""
        
        for page in TEST_PAGES:
            report += f"| {page['name']} | {page['path']} | {'是' if page['requires_auth'] else '否'} |\n"
        
        report += """
## 测试设备

| 设备 | 分辨率 | 类型 |
|------|--------|------|
"""
        
        for device in DEVICES:
            device_type = "桌面端" if device["width"] >= 1366 else ("平板" if device["width"] >= 768 else "手机")
            report += f"| {device['name']} | {device['width']}x{device['height']} | {device_type} |\n"
        
        report += """
## 详细结果

"""
        
        # 按页面分组
        pages_tested = {}
        for result in self.results:
            page_name = result["page"]
            if page_name not in pages_tested:
                pages_tested[page_name] = []
            pages_tested[page_name].append(result)
        
        for page_name, page_results in pages_tested.items():
            report += f"### {page_name}\n\n"
            report += "| 设备 | 状态 | 差异率 | 备注 |\n"
            report += "|------|------|--------|------|\n"
            
            for result in page_results:
                status_emoji = {
                    "passed": "✅ 通过",
                    "failed": "❌ 失败",
                    "warning": "⚠️ 警告",
                    "new_baseline": "🆕 新基准",
                    "error": "💥 错误",
                    "skipped": "⏭️ 跳过"
                }.get(result["status"], result["status"])
                
                diff = f"{result['diff_percentage']}%" if result["diff_percentage"] > 0 else "-"
                error = result.get("error", "")
                
                report += f"| {result['device']} | {status_emoji} | {diff} | {error} |\n"
            
            report += "\n"
        
        report += """
## 使用说明

1. **通过 (✅)**: 当前截图与基准图差异小于1%
2. **警告 (⚠️)**: 差异在1%-5%之间，可能是细微变化
3. **失败 (❌)**: 差异超过5%，需要检查UI变化
4. **新基准 (🆕)**: 首次运行，已保存为基准图
5. **错误 (💥)**: 测试过程出现异常
6. **跳过 (⏭️)**: 因依赖条件不满足而跳过

## 文件位置

- **基准图目录**: `test_results/visual_test/baseline/`
- **当前截图目录**: `test_results/visual_test/current/`
- **差异图目录**: `test_results/visual_test/diff/`

---
*报告由自动视觉UI测试工具生成*
"""
        
        with open(report_path, "w", encoding="utf-8") as f:
            f.write(report)
        
        print(f"\n报告已生成: {report_path}")
        return report_path
    
    def save_results_json(self):
        """保存JSON格式结果"""
        json_path = os.path.join(RESULTS_DIR, f"results_{self.timestamp}.json")
        with open(json_path, "w", encoding="utf-8") as f:
            json.dump({
                "timestamp": self.timestamp,
                "total_tests": len(self.results),
                "results": self.results
            }, f, ensure_ascii=False, indent=2)
        print(f"JSON结果已保存: {json_path}")
        return json_path


def main():
    """主函数"""
    tester = VisualUITester()
    
    # 运行所有测试
    results = tester.run_all_tests()
    
    # 生成报告
    tester.generate_report()
    tester.save_results_json()
    
    # 打印摘要
    print("\n" + "=" * 60)
    print("测试完成!")
    print("=" * 60)
    
    total = len(results)
    passed = len([r for r in results if r["status"] == "passed"])
    failed = len([r for r in results if r["status"] == "failed"])
    new_baseline = len([r for r in results if r["status"] == "new_baseline"])
    
    print(f"总测试: {total}")
    print(f"通过: {passed}")
    print(f"失败: {failed}")
    print(f"新基准: {new_baseline}")
    print(f"通过率: {round(passed/total*100, 1) if total > 0 else 0}%")


if __name__ == "__main__":
    main()

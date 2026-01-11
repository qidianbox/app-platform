#!/usr/bin/env python3
"""
自动视觉UI测试脚本 V2
修复版本：使用API登录获取token，通过localStorage注入认证状态
"""

import os
import json
import time
import hashlib
import shutil
import requests
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from PIL import Image, ImageChops, ImageDraw, ImageFont

# 配置
BASE_URL = "http://localhost:5173"
API_URL = "http://localhost:8080"
RESULTS_DIR = "/home/ubuntu/app-platform/test_results/visual_test"
BASELINE_DIR = os.path.join(RESULTS_DIR, "baseline")
CURRENT_DIR = os.path.join(RESULTS_DIR, "current")
DIFF_DIR = os.path.join(RESULTS_DIR, "diff")

# 测试页面配置
# 注意：工作台是APP详情页面内的Tab，不是独立路由
TEST_PAGES = [
    {"name": "登录页", "path": "/login", "requires_auth": False, "wait_time": 2},
    {"name": "APP列表页", "path": "/apps", "requires_auth": True, "wait_time": 2},
    {"name": "APP详情_概览", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3},
    {"name": "APP详情_工作台", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_tab"},
    {"name": "工作台_监控告警", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_then_monitor"},
    {"name": "工作台_审计日志", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_then_audit"},
    {"name": "工作台_消息推送", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_then_push"},
    {"name": "工作台_版本管理", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_then_version"},
    {"name": "工作台_用户管理", "path": "/apps/2/config", "requires_auth": True, "wait_time": 3, "action": "click_workspace_then_user"},
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
        self.auth_token = None
        
        for dir_path in [RESULTS_DIR, BASELINE_DIR, CURRENT_DIR, DIFF_DIR]:
            os.makedirs(dir_path, exist_ok=True)
    
    def get_auth_token(self):
        """通过API获取认证token"""
        try:
            response = requests.post(
                f"{API_URL}/api/v1/admin/login",
                json={"username": "admin", "password": "admin123"},
                headers={"Content-Type": "application/json"},
                timeout=10
            )
            data = response.json()
            if data.get("code") == 0 and data.get("data", {}).get("token"):
                self.auth_token = data["data"]["token"]
                print(f"获取token成功: {self.auth_token[:50]}...")
                return True
            else:
                print(f"登录失败: {data}")
                return False
        except Exception as e:
            print(f"获取token失败: {e}")
            return False
    
    def setup_driver(self, width, height):
        """设置浏览器驱动"""
        chrome_options = Options()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument(f"--window-size={width},{height}")
        chrome_options.add_argument("--force-device-scale-factor=1")
        
        service = Service("/usr/bin/chromedriver")
        self.driver = webdriver.Chrome(service=service, options=chrome_options)
        self.driver.set_window_size(width, height)
    
    def inject_auth(self):
        """注入认证token到localStorage"""
        if not self.auth_token:
            return False
        
        try:
            # 先访问首页以初始化localStorage
            self.driver.get(BASE_URL)
            time.sleep(1)
            
            # 注入token到localStorage
            self.driver.execute_script(f"""
                localStorage.setItem('token', '{self.auth_token}');
                localStorage.setItem('admin_token', '{self.auth_token}');
            """)
            print("Token已注入localStorage")
            return True
        except Exception as e:
            print(f"注入token失败: {e}")
            return False
    
    def click_workspace_tab(self):
        """点击工作台Tab"""
        try:
            # 尝试多种方式点击工作台Tab
            selectors = [
                ".nav-item:contains('工作台')",
                "//div[contains(@class, 'nav-item') and contains(text(), '工作台')]",
                "//div[contains(text(), '工作台')]"
            ]
            
            # 使用JavaScript点击
            js_click = """
                var items = document.querySelectorAll('.nav-item');
                for (var i = 0; i < items.length; i++) {
                    if (items[i].textContent.includes('工作台')) {
                        items[i].click();
                        return true;
                    }
                }
                return false;
            """
            result = self.driver.execute_script(js_click)
            if result:
                print("点击工作台Tab成功")
                time.sleep(2)
                return True
            
            # 备用：XPath方式
            try:
                elem = self.driver.find_element(By.XPATH, "//div[contains(@class, 'nav-item') and contains(text(), '工作台')]")
                self.driver.execute_script("arguments[0].click();", elem)
                print("通过XPath点击工作台Tab成功")
                time.sleep(2)
                return True
            except:
                pass
            
            print("未能找到工作台Tab")
            return False
        except Exception as e:
            print(f"点击工作台Tab失败: {e}")
            return False
    
    def click_workspace_menu(self, menu_key):
        """在工作台中点击菜单项"""
        try:
            # 先点击工作台Tab
            self.click_workspace_tab()
            time.sleep(1)
            
            # 然后点击具体菜单项
            text_map = {
                "monitor": "监控告警",
                "audit": "审计日志",
                "messages": "消息推送",
                "versions": "版本管理",
                "users": "用户管理"
            }
            
            menu_text = text_map.get(menu_key, "")
            if not menu_text:
                return False
            
            # 使用JavaScript点击菜单项
            js_click = f"""
                var items = document.querySelectorAll('.menu-item');
                for (var i = 0; i < items.length; i++) {{
                    if (items[i].textContent.includes('{menu_text}')) {{
                        items[i].click();
                        return true;
                    }}
                }}
                return false;
            """
            result = self.driver.execute_script(js_click)
            if result:
                print(f"点击菜单成功: {menu_text}")
                time.sleep(2)
                return True
            
            # 备用：使用data-testid
            try:
                elem = self.driver.find_element(By.CSS_SELECTOR, f"[data-testid='menu-{menu_key}']")
                self.driver.execute_script("arguments[0].click();", elem)
                print(f"通过data-testid点击菜单成功: {menu_key}")
                time.sleep(2)
                return True
            except:
                pass
            
            print(f"未能找到菜单项: {menu_text}")
            return False
        except Exception as e:
            print(f"点击菜单失败: {e}")
            return False
    
    def execute_action(self, action):
        """执行页面操作"""
        try:
            time.sleep(1)
            
            # 工作台Tab点击
            if action == "click_workspace_tab":
                return self.click_workspace_tab()
            
            # 工作台内菜单点击
            workspace_menu_actions = {
                "click_workspace_then_monitor": "monitor",
                "click_workspace_then_audit": "audit",
                "click_workspace_then_push": "messages",
                "click_workspace_then_version": "versions",
                "click_workspace_then_user": "users",
            }
            
            if action in workspace_menu_actions:
                return self.click_workspace_menu(workspace_menu_actions[action])
            
            # 旧的菜单点击逻辑（备用）
            action_map = {
                "click_monitor": "monitor",
                "click_audit": "audit",
                "click_push": "messages",
                "click_version": "versions",
                "click_user": "users",
            }
            
            if action in action_map:
                menu_key = action_map[action]
                # 尝试多种定位方式
                selectors = [
                    f"[data-testid='menu-{menu_key}']",
                    f"[data-menu-key='{menu_key}']",
                    f".menu-item[data-menu-key='{menu_key}']"
                ]
                
                for selector in selectors:
                    try:
                        elem = self.driver.find_element(By.CSS_SELECTOR, selector)
                        if elem.is_displayed():
                            self.driver.execute_script("arguments[0].click();", elem)
                            print(f"点击菜单成功: {menu_key}")
                            time.sleep(2)
                            return True
                    except:
                        continue
                
                # 备用方案：使用文本匹配
                text_map = {
                    "monitor": "监控告警",
                    "audit": "审计日志",
                    "messages": "消息推送",
                    "versions": "版本管理",
                    "users": "用户管理"
                }
                if menu_key in text_map:
                    try:
                        elem = self.driver.find_element(By.XPATH, f"//span[contains(text(), '{text_map[menu_key]}')]")
                        if elem.is_displayed():
                            self.driver.execute_script("arguments[0].click();", elem)
                            print(f"通过文本点击菜单成功: {menu_key}")
                            time.sleep(2)
                            return True
                    except:
                        pass
                        
                print(f"未能找到菜单项: {menu_key}")
            return True
        except Exception as e:
            print(f"执行操作失败: {action}, {e}")
            return False
    
    def take_screenshot(self, name, device_name):
        """截图"""
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
            
            if img1.size != img2.size:
                img2 = img2.resize(img1.size, Image.Resampling.LANCZOS)
            
            diff = ImageChops.difference(img1, img2)
            diff_pixels = sum(1 for pixel in diff.getdata() if pixel != (0, 0, 0))
            total_pixels = img1.size[0] * img1.size[1]
            
            return (diff_pixels / total_pixels) * 100, diff
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
            
            diff = ImageChops.difference(img1, img2)
            diff = diff.point(lambda x: min(255, x * 10))
            
            width = img1.size[0] * 3
            height = img1.size[1]
            comparison = Image.new('RGB', (width, height))
            
            comparison.paste(img1, (0, 0))
            comparison.paste(img2, (img1.size[0], 0))
            comparison.paste(diff, (img1.size[0] * 2, 0))
            
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
    
    def run_test_for_device(self, device):
        """为单个设备运行所有测试"""
        results = []
        
        try:
            self.setup_driver(device["width"], device["height"])
            
            # 注入认证
            if self.auth_token:
                self.inject_auth()
            
            for page in TEST_PAGES:
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
                    
                    # 检查认证
                    if page["requires_auth"] and not self.auth_token:
                        result["status"] = "skipped"
                        result["error"] = "无认证token"
                        results.append(result)
                        continue
                    
                    # 对于需要认证的页面，先访问首页并注入token
                    if page["requires_auth"] and self.auth_token:
                        self.driver.get(BASE_URL)
                        time.sleep(0.5)
                        self.driver.execute_script(f"""
                            localStorage.setItem('token', '{self.auth_token}');
                        """)
                        time.sleep(0.3)
                    
                    # 访问页面
                    self.driver.get(f"{BASE_URL}{page['path']}")
                    time.sleep(page.get("wait_time", 2))
                    
                    # 执行操作
                    if "action" in page:
                        self.execute_action(page["action"])
                        time.sleep(2)
                    
                    # 额外等待
                    time.sleep(1)
                    
                    # 截图
                    screenshot_path = self.take_screenshot(page["name"], device["name"])
                    result["screenshot"] = screenshot_path
                    
                    # 查找基准图
                    safe_name = page['name'].replace(' ', '_').replace('-', '_').replace('/', '_')
                    baseline_pattern = f"{safe_name}_{device['name']}_"
                    baseline_files = [f for f in os.listdir(BASELINE_DIR) if f.startswith(baseline_pattern)]
                    
                    if baseline_files:
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
                        
                        if diff_percentage > 0:
                            diff_filename = f"diff_{safe_name}_{device['name']}_{self.timestamp}.png"
                            diff_path = os.path.join(DIFF_DIR, diff_filename)
                            self.create_diff_image(baseline_path, screenshot_path, diff_path)
                            result["diff_image"] = diff_path
                    else:
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
            print(f"设备测试出错: {e}")
        finally:
            if self.driver:
                self.driver.quit()
                self.driver = None
        
        return results
    
    def run_all_tests(self):
        """运行所有测试"""
        print("=" * 60)
        print("开始视觉UI测试 V2")
        print(f"测试时间: {self.timestamp}")
        print(f"测试页面: {len(TEST_PAGES)}")
        print(f"测试设备: {len(DEVICES)}")
        print("=" * 60)
        
        # 获取认证token
        print("\n获取认证token...")
        self.get_auth_token()
        
        all_results = []
        
        for device in DEVICES:
            print(f"\n{'='*40}")
            print(f"设备: {device['name']} ({device['width']}x{device['height']})")
            print("=" * 40)
            
            results = self.run_test_for_device(device)
            all_results.extend(results)
        
        self.results = all_results
        return all_results
    
    def generate_report(self):
        """生成测试报告"""
        report_path = os.path.join(RESULTS_DIR, f"visual_test_report_{self.timestamp}.md")
        
        total = len(self.results)
        passed = len([r for r in self.results if r["status"] == "passed"])
        failed = len([r for r in self.results if r["status"] == "failed"])
        warning = len([r for r in self.results if r["status"] == "warning"])
        new_baseline = len([r for r in self.results if r["status"] == "new_baseline"])
        errors = len([r for r in self.results if r["status"] in ["error", "skipped"]])
        
        report = f"""# 视觉UI测试报告 V2

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
        
        report += "\n## 详细结果\n\n"
        
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

---
*报告由自动视觉UI测试工具 V2 生成*
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
    
    results = tester.run_all_tests()
    tester.generate_report()
    tester.save_results_json()
    
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

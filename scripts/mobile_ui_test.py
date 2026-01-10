#!/usr/bin/env python3
"""
移动端UI兼容性自动化测试脚本
测试APP中台管理系统在不同移动设备上的响应式布局
"""

import os
import json
import time
from datetime import datetime
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException

# 测试设备配置
DEVICES = [
    {"name": "iPhone_SE", "width": 375, "height": 667, "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1"},
    {"name": "iPhone_12", "width": 390, "height": 844, "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1"},
    {"name": "iPhone_14_Pro_Max", "width": 430, "height": 932, "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1"},
    {"name": "iPad_Mini", "width": 768, "height": 1024, "user_agent": "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Mobile/15E148 Safari/604.1"},
    {"name": "Android_Phone", "width": 360, "height": 800, "user_agent": "Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36"},
    {"name": "Android_Tablet", "width": 800, "height": 1280, "user_agent": "Mozilla/5.0 (Linux; Android 12; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"},
]

# 测试页面配置
TEST_PAGES = [
    {"name": "登录页", "path": "/", "requires_auth": False},
    {"name": "APP列表", "path": "/apps", "requires_auth": True},
    {"name": "APP详情", "path": "/apps/1/workspace", "requires_auth": True},
]

# 测试检查项
UI_CHECKS = [
    {"name": "页面可访问", "type": "page_load"},
    {"name": "无水平滚动", "type": "no_horizontal_scroll"},
    {"name": "文字可读", "type": "text_readable"},
    {"name": "按钮可点击", "type": "buttons_clickable"},
    {"name": "导航可用", "type": "navigation_usable"},
    {"name": "表单可用", "type": "forms_usable"},
    {"name": "图片适配", "type": "images_responsive"},
]

class MobileUITester:
    def __init__(self, base_url):
        self.base_url = base_url
        self.results = []
        self.screenshots_dir = "/home/ubuntu/app-platform/test_results/screenshots"
        self.driver = None
        
        # 创建截图目录
        os.makedirs(self.screenshots_dir, exist_ok=True)
    
    def create_driver(self, device):
        """创建模拟移动设备的浏览器驱动"""
        options = Options()
        options.add_argument('--headless')
        options.add_argument('--no-sandbox')
        options.add_argument('--disable-dev-shm-usage')
        options.add_argument('--disable-gpu')
        options.add_argument(f'--window-size={device["width"]},{device["height"]}')
        options.add_argument(f'--user-agent={device["user_agent"]}')
        
        # 模拟移动设备
        mobile_emulation = {
            "deviceMetrics": {"width": device["width"], "height": device["height"], "pixelRatio": 2.0},
            "userAgent": device["user_agent"]
        }
        options.add_experimental_option("mobileEmulation", mobile_emulation)
        
        service = Service('/usr/bin/chromedriver')
        driver = webdriver.Chrome(service=service, options=options)
        driver.set_page_load_timeout(30)
        return driver
    
    def take_screenshot(self, device_name, page_name):
        """截取当前页面截图"""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"{device_name}_{page_name}_{timestamp}.png"
        filepath = os.path.join(self.screenshots_dir, filename)
        self.driver.save_screenshot(filepath)
        return filepath
    
    def check_page_load(self):
        """检查页面是否成功加载"""
        try:
            WebDriverWait(self.driver, 10).until(
                lambda d: d.execute_script("return document.readyState") == "complete"
            )
            return {"passed": True, "message": "页面加载成功"}
        except TimeoutException:
            return {"passed": False, "message": "页面加载超时"}
    
    def check_no_horizontal_scroll(self):
        """检查是否存在水平滚动条"""
        try:
            body_width = self.driver.execute_script("return document.body.scrollWidth")
            viewport_width = self.driver.execute_script("return window.innerWidth")
            
            if body_width > viewport_width + 10:  # 允许10px误差
                return {"passed": False, "message": f"存在水平滚动: body宽度{body_width}px > 视口宽度{viewport_width}px"}
            return {"passed": True, "message": "无水平滚动"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def check_text_readable(self):
        """检查文字是否可读（字体大小>=12px）"""
        try:
            script = """
            var elements = document.querySelectorAll('p, span, div, a, button, label, h1, h2, h3, h4, h5, h6');
            var smallTexts = [];
            for (var i = 0; i < elements.length; i++) {
                var style = window.getComputedStyle(elements[i]);
                var fontSize = parseFloat(style.fontSize);
                if (fontSize < 12 && elements[i].innerText.trim().length > 0) {
                    smallTexts.push({
                        tag: elements[i].tagName,
                        text: elements[i].innerText.substring(0, 50),
                        fontSize: fontSize
                    });
                }
            }
            return smallTexts.slice(0, 5);
            """
            small_texts = self.driver.execute_script(script)
            
            if len(small_texts) > 0:
                return {"passed": False, "message": f"发现{len(small_texts)}个小于12px的文字元素", "details": small_texts}
            return {"passed": True, "message": "所有文字大小>=12px"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def check_buttons_clickable(self):
        """检查按钮是否可点击（尺寸>=44x44px）"""
        try:
            script = """
            var buttons = document.querySelectorAll('button, a, [role="button"], input[type="submit"], input[type="button"]');
            var smallButtons = [];
            for (var i = 0; i < buttons.length; i++) {
                var rect = buttons[i].getBoundingClientRect();
                if (rect.width > 0 && rect.height > 0 && (rect.width < 44 || rect.height < 44)) {
                    smallButtons.push({
                        tag: buttons[i].tagName,
                        text: buttons[i].innerText.substring(0, 30) || buttons[i].getAttribute('aria-label') || '',
                        width: Math.round(rect.width),
                        height: Math.round(rect.height)
                    });
                }
            }
            return smallButtons.slice(0, 10);
            """
            small_buttons = self.driver.execute_script(script)
            
            if len(small_buttons) > 0:
                return {"passed": False, "message": f"发现{len(small_buttons)}个小于44x44px的可点击元素", "details": small_buttons}
            return {"passed": True, "message": "所有可点击元素尺寸>=44x44px"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def check_navigation_usable(self):
        """检查导航是否可用"""
        try:
            # 检查是否有可见的导航元素
            nav_selectors = ['nav', '[role="navigation"]', '.nav', '.navbar', '.sidebar', '.menu']
            nav_found = False
            
            for selector in nav_selectors:
                try:
                    elements = self.driver.find_elements(By.CSS_SELECTOR, selector)
                    for el in elements:
                        if el.is_displayed():
                            nav_found = True
                            break
                except:
                    continue
            
            # 检查是否有汉堡菜单（移动端常见）
            hamburger_selectors = ['.hamburger', '.menu-toggle', '[aria-label*="menu"]', '.mobile-menu', 'button[class*="menu"]']
            hamburger_found = False
            
            for selector in hamburger_selectors:
                try:
                    elements = self.driver.find_elements(By.CSS_SELECTOR, selector)
                    for el in elements:
                        if el.is_displayed():
                            hamburger_found = True
                            break
                except:
                    continue
            
            if nav_found or hamburger_found:
                return {"passed": True, "message": "导航可用"}
            return {"passed": False, "message": "未找到可用的导航元素"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def check_forms_usable(self):
        """检查表单是否可用"""
        try:
            script = """
            var inputs = document.querySelectorAll('input, textarea, select');
            var issues = [];
            for (var i = 0; i < inputs.length; i++) {
                var rect = inputs[i].getBoundingClientRect();
                if (rect.width > 0 && rect.height > 0) {
                    // 检查输入框高度是否足够
                    if (rect.height < 40) {
                        issues.push({
                            type: inputs[i].type || inputs[i].tagName.toLowerCase(),
                            name: inputs[i].name || inputs[i].placeholder || '',
                            height: Math.round(rect.height),
                            issue: '高度不足'
                        });
                    }
                    // 检查输入框是否太窄
                    if (rect.width < 100 && inputs[i].type !== 'checkbox' && inputs[i].type !== 'radio') {
                        issues.push({
                            type: inputs[i].type || inputs[i].tagName.toLowerCase(),
                            name: inputs[i].name || inputs[i].placeholder || '',
                            width: Math.round(rect.width),
                            issue: '宽度不足'
                        });
                    }
                }
            }
            return issues.slice(0, 10);
            """
            issues = self.driver.execute_script(script)
            
            if len(issues) > 0:
                return {"passed": False, "message": f"发现{len(issues)}个表单元素存在问题", "details": issues}
            return {"passed": True, "message": "表单元素尺寸合适"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def check_images_responsive(self):
        """检查图片是否响应式"""
        try:
            script = """
            var images = document.querySelectorAll('img');
            var issues = [];
            var viewportWidth = window.innerWidth;
            for (var i = 0; i < images.length; i++) {
                var rect = images[i].getBoundingClientRect();
                if (rect.width > viewportWidth) {
                    issues.push({
                        src: images[i].src.substring(0, 50),
                        width: Math.round(rect.width),
                        viewportWidth: viewportWidth,
                        issue: '图片超出视口'
                    });
                }
            }
            return issues.slice(0, 5);
            """
            issues = self.driver.execute_script(script)
            
            if len(issues) > 0:
                return {"passed": False, "message": f"发现{len(issues)}个图片超出视口", "details": issues}
            return {"passed": True, "message": "所有图片适配视口"}
        except Exception as e:
            return {"passed": False, "message": f"检查失败: {str(e)}"}
    
    def run_check(self, check_type):
        """运行指定类型的检查"""
        check_methods = {
            "page_load": self.check_page_load,
            "no_horizontal_scroll": self.check_no_horizontal_scroll,
            "text_readable": self.check_text_readable,
            "buttons_clickable": self.check_buttons_clickable,
            "navigation_usable": self.check_navigation_usable,
            "forms_usable": self.check_forms_usable,
            "images_responsive": self.check_images_responsive,
        }
        
        if check_type in check_methods:
            return check_methods[check_type]()
        return {"passed": False, "message": f"未知检查类型: {check_type}"}
    
    def test_page(self, device, page):
        """测试单个页面"""
        page_result = {
            "device": device["name"],
            "page": page["name"],
            "path": page["path"],
            "viewport": f"{device['width']}x{device['height']}",
            "checks": [],
            "screenshot": None,
            "timestamp": datetime.now().isoformat()
        }
        
        try:
            # 访问页面
            url = f"{self.base_url}{page['path']}"
            self.driver.get(url)
            time.sleep(2)  # 等待页面渲染
            
            # 截图
            screenshot_path = self.take_screenshot(device["name"], page["name"].replace(" ", "_"))
            page_result["screenshot"] = screenshot_path
            
            # 运行所有检查
            for check in UI_CHECKS:
                check_result = self.run_check(check["type"])
                check_result["name"] = check["name"]
                check_result["type"] = check["type"]
                page_result["checks"].append(check_result)
            
        except Exception as e:
            page_result["error"] = str(e)
        
        return page_result
    
    def run_tests(self):
        """运行所有测试"""
        print("=" * 60)
        print("移动端UI兼容性测试")
        print("=" * 60)
        print(f"测试URL: {self.base_url}")
        print(f"测试设备: {len(DEVICES)}个")
        print(f"测试页面: {len(TEST_PAGES)}个")
        print("=" * 60)
        
        for device in DEVICES:
            print(f"\n测试设备: {device['name']} ({device['width']}x{device['height']})")
            print("-" * 40)
            
            try:
                self.driver = self.create_driver(device)
                
                for page in TEST_PAGES:
                    print(f"  测试页面: {page['name']} ({page['path']})")
                    result = self.test_page(device, page)
                    self.results.append(result)
                    
                    # 打印检查结果摘要
                    passed = sum(1 for c in result["checks"] if c.get("passed", False))
                    total = len(result["checks"])
                    status = "✓" if passed == total else "✗"
                    print(f"    {status} 通过: {passed}/{total}")
                
            except Exception as e:
                print(f"  错误: {str(e)}")
            finally:
                if self.driver:
                    self.driver.quit()
                    self.driver = None
        
        return self.results
    
    def generate_report(self):
        """生成测试报告"""
        report = {
            "title": "移动端UI兼容性测试报告",
            "generated_at": datetime.now().isoformat(),
            "base_url": self.base_url,
            "summary": {
                "total_tests": len(self.results),
                "devices_tested": len(DEVICES),
                "pages_tested": len(TEST_PAGES),
                "total_checks": 0,
                "passed_checks": 0,
                "failed_checks": 0,
            },
            "device_summary": {},
            "page_summary": {},
            "issues": [],
            "results": self.results
        }
        
        # 统计结果
        for result in self.results:
            device = result["device"]
            page = result["page"]
            
            if device not in report["device_summary"]:
                report["device_summary"][device] = {"passed": 0, "failed": 0, "total": 0}
            if page not in report["page_summary"]:
                report["page_summary"][page] = {"passed": 0, "failed": 0, "total": 0}
            
            for check in result.get("checks", []):
                report["summary"]["total_checks"] += 1
                report["device_summary"][device]["total"] += 1
                report["page_summary"][page]["total"] += 1
                
                if check.get("passed", False):
                    report["summary"]["passed_checks"] += 1
                    report["device_summary"][device]["passed"] += 1
                    report["page_summary"][page]["passed"] += 1
                else:
                    report["summary"]["failed_checks"] += 1
                    report["device_summary"][device]["failed"] += 1
                    report["page_summary"][page]["failed"] += 1
                    
                    # 记录问题
                    report["issues"].append({
                        "device": device,
                        "page": page,
                        "check": check["name"],
                        "message": check.get("message", ""),
                        "details": check.get("details", None)
                    })
        
        # 计算通过率
        if report["summary"]["total_checks"] > 0:
            report["summary"]["pass_rate"] = round(
                report["summary"]["passed_checks"] / report["summary"]["total_checks"] * 100, 2
            )
        else:
            report["summary"]["pass_rate"] = 0
        
        return report
    
    def save_report(self, report, output_dir="/home/ubuntu/app-platform/test_results"):
        """保存测试报告"""
        os.makedirs(output_dir, exist_ok=True)
        
        # 保存JSON报告
        json_path = os.path.join(output_dir, "mobile_ui_test_report.json")
        with open(json_path, "w", encoding="utf-8") as f:
            json.dump(report, f, ensure_ascii=False, indent=2)
        
        # 生成Markdown报告
        md_path = os.path.join(output_dir, "mobile_ui_test_report.md")
        with open(md_path, "w", encoding="utf-8") as f:
            f.write(f"# {report['title']}\n\n")
            f.write(f"**生成时间**: {report['generated_at']}\n\n")
            f.write(f"**测试URL**: {report['base_url']}\n\n")
            
            f.write("## 测试摘要\n\n")
            f.write(f"| 指标 | 数值 |\n")
            f.write(f"|------|------|\n")
            f.write(f"| 测试设备数 | {report['summary']['devices_tested']} |\n")
            f.write(f"| 测试页面数 | {report['summary']['pages_tested']} |\n")
            f.write(f"| 总检查项 | {report['summary']['total_checks']} |\n")
            f.write(f"| 通过项 | {report['summary']['passed_checks']} |\n")
            f.write(f"| 失败项 | {report['summary']['failed_checks']} |\n")
            f.write(f"| **通过率** | **{report['summary']['pass_rate']}%** |\n\n")
            
            f.write("## 设备测试结果\n\n")
            f.write("| 设备 | 通过 | 失败 | 通过率 |\n")
            f.write("|------|------|------|--------|\n")
            for device, stats in report["device_summary"].items():
                rate = round(stats["passed"] / stats["total"] * 100, 1) if stats["total"] > 0 else 0
                status = "✓" if rate >= 80 else "⚠" if rate >= 60 else "✗"
                f.write(f"| {device} | {stats['passed']} | {stats['failed']} | {status} {rate}% |\n")
            
            f.write("\n## 页面测试结果\n\n")
            f.write("| 页面 | 通过 | 失败 | 通过率 |\n")
            f.write("|------|------|------|--------|\n")
            for page, stats in report["page_summary"].items():
                rate = round(stats["passed"] / stats["total"] * 100, 1) if stats["total"] > 0 else 0
                status = "✓" if rate >= 80 else "⚠" if rate >= 60 else "✗"
                f.write(f"| {page} | {stats['passed']} | {stats['failed']} | {status} {rate}% |\n")
            
            if report["issues"]:
                f.write("\n## 发现的问题\n\n")
                for i, issue in enumerate(report["issues"], 1):
                    f.write(f"### {i}. {issue['check']} - {issue['device']} - {issue['page']}\n\n")
                    f.write(f"**问题描述**: {issue['message']}\n\n")
                    if issue.get("details"):
                        f.write(f"**详细信息**:\n```json\n{json.dumps(issue['details'], ensure_ascii=False, indent=2)}\n```\n\n")
            
            f.write("\n## 截图\n\n")
            f.write("截图保存在 `test_results/screenshots/` 目录下\n")
        
        print(f"\n报告已保存:")
        print(f"  JSON: {json_path}")
        print(f"  Markdown: {md_path}")
        
        return json_path, md_path


def main():
    # 测试URL（使用本地前端服务）
    base_url = "http://localhost:5173"
    
    tester = MobileUITester(base_url)
    
    # 运行测试
    tester.run_tests()
    
    # 生成报告
    report = tester.generate_report()
    
    # 保存报告
    tester.save_report(report)
    
    # 打印摘要
    print("\n" + "=" * 60)
    print("测试完成!")
    print("=" * 60)
    print(f"通过率: {report['summary']['pass_rate']}%")
    print(f"通过: {report['summary']['passed_checks']}/{report['summary']['total_checks']}")
    if report["issues"]:
        print(f"发现 {len(report['issues'])} 个问题")


if __name__ == "__main__":
    main()

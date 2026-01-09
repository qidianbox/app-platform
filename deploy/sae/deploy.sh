#!/bin/bash
# APP中台管理系统 - 阿里云SAE自动化部署脚本
# 作者：Manus AI
# 版本：1.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step() { echo -e "${BLUE}[STEP]${NC} $1"; }

# 配置变量
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
REGION="${REGION:-cn-hangzhou}"
PROJECT_NAME="${PROJECT_NAME:-app-platform}"

# 检查必要的工具
check_prerequisites() {
    log_step "检查必要的工具..."
    
    local missing_tools=()
    
    command -v terraform &> /dev/null || missing_tools+=("terraform")
    command -v docker &> /dev/null || missing_tools+=("docker")
    command -v aliyun &> /dev/null || missing_tools+=("aliyun-cli")
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "缺少以下工具: ${missing_tools[*]}"
        log_info "安装指南："
        log_info "  Terraform: https://www.terraform.io/downloads"
        log_info "  Docker: https://docs.docker.com/get-docker/"
        log_info "  阿里云CLI: https://help.aliyun.com/document_detail/139508.html"
        exit 1
    fi
    
    log_info "所有必要工具已安装 ✓"
}

# 检查阿里云凭证
check_credentials() {
    log_step "检查阿里云凭证..."
    
    if [ -z "$ALICLOUD_ACCESS_KEY" ] || [ -z "$ALICLOUD_SECRET_KEY" ]; then
        log_error "请设置阿里云凭证环境变量："
        log_info "  export ALICLOUD_ACCESS_KEY=your-access-key"
        log_info "  export ALICLOUD_SECRET_KEY=your-secret-key"
        exit 1
    fi
    
    # 验证凭证有效性
    if ! aliyun sts GetCallerIdentity &> /dev/null; then
        log_error "阿里云凭证无效，请检查AccessKey"
        exit 1
    fi
    
    log_info "阿里云凭证验证通过 ✓"
}

# 阶段一：部署基础设施
deploy_infrastructure() {
    log_step "========== 阶段一：部署基础设施 =========="
    
    cd "$SCRIPT_DIR"
    
    # 检查terraform.tfvars
    if [ ! -f "terraform.tfvars" ]; then
        log_warn "未找到terraform.tfvars，从示例文件创建..."
        cp terraform.tfvars.example terraform.tfvars
        log_error "请编辑 terraform.tfvars 填写实际配置后重新运行"
        exit 1
    fi
    
    log_info "初始化 Terraform..."
    terraform init
    
    log_info "验证 Terraform 配置..."
    terraform validate
    
    log_info "规划基础设施变更..."
    terraform plan -out=tfplan
    
    log_info "应用基础设施变更（预计需要10-15分钟）..."
    terraform apply tfplan
    
    # 保存输出
    terraform output -json > outputs.json
    
    log_info "基础设施部署完成 ✓"
}

# 阶段二：构建并推送Docker镜像
build_and_push_image() {
    log_step "========== 阶段二：构建并推送Docker镜像 =========="
    
    # 获取ACR信息
    local acr_repo=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.acr_repo_domain.value')
    local image_tag="${acr_repo}:$(date +%Y%m%d%H%M%S)"
    local image_latest="${acr_repo}:latest"
    
    log_info "构建Docker镜像..."
    cd "$PROJECT_ROOT/backend"
    docker build -t "$image_tag" -t "$image_latest" .
    
    log_info "登录ACR..."
    docker login --username="$ALICLOUD_ACCESS_KEY" --password="$ALICLOUD_SECRET_KEY" "registry.${REGION}.aliyuncs.com"
    
    log_info "推送镜像到ACR..."
    docker push "$image_tag"
    docker push "$image_latest"
    
    log_info "镜像推送完成: $image_tag ✓"
    
    # 保存镜像标签
    echo "$image_tag" > "$SCRIPT_DIR/image_tag.txt"
}

# 阶段三：更新SAE应用
update_sae_application() {
    log_step "========== 阶段三：更新SAE应用 =========="
    
    local image_tag=$(cat "$SCRIPT_DIR/image_tag.txt")
    local app_id=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.sae_app_id.value')
    
    log_info "更新SAE应用镜像..."
    aliyun sae DeployApplication \
        --AppId "$app_id" \
        --ImageUrl "$image_tag" \
        --region "$REGION"
    
    log_info "等待部署完成..."
    sleep 30
    
    # 检查部署状态
    local status=$(aliyun sae DescribeApplicationStatus --AppId "$app_id" --region "$REGION" | jq -r '.Data.CurrentStatus')
    
    if [ "$status" == "RUNNING" ]; then
        log_info "SAE应用部署成功 ✓"
    else
        log_warn "SAE应用状态: $status，请在控制台查看详情"
    fi
}

# 阶段四：部署前端
deploy_frontend() {
    log_step "========== 阶段四：部署前端 =========="
    
    local oss_bucket=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.oss_bucket.value')
    local backend_url=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.backend_slb_address.value')
    
    log_info "构建前端..."
    cd "$PROJECT_ROOT/frontend"
    
    # 更新API地址
    echo "VITE_API_BASE_URL=http://${backend_url}" > .env.production
    
    npm install
    npm run build
    
    log_info "上传到OSS..."
    aliyun oss cp -r dist/ "oss://${oss_bucket}/" --force --region "$REGION"
    
    log_info "前端部署完成 ✓"
}

# 阶段五：数据库初始化
init_database() {
    log_step "========== 阶段五：数据库初始化 =========="
    
    local rds_host=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.rds_connection_string.value')
    local rds_port=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.rds_port.value')
    
    log_info "数据库连接信息："
    log_info "  主机: $rds_host"
    log_info "  端口: $rds_port"
    log_info "  数据库: app_platform"
    
    log_warn "请手动执行数据库迁移脚本，或通过应用自动迁移"
    log_info "数据库初始化步骤完成 ✓"
}

# 显示部署结果
show_results() {
    log_step "========== 部署完成 =========="
    
    local frontend_url=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.frontend_url.value')
    local backend_url=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.backend_slb_address.value')
    local sae_app_name=$(cat "$SCRIPT_DIR/outputs.json" | jq -r '.sae_app_name.value')
    
    echo ""
    echo "=================================================="
    echo "  🎉 APP中台管理系统部署成功！"
    echo "=================================================="
    echo ""
    echo "  前端地址: $frontend_url"
    echo "  后端API: http://${backend_url}"
    echo "  SAE应用: $sae_app_name"
    echo ""
    echo "  默认管理员账号: admin"
    echo "  默认管理员密码: admin123"
    echo ""
    echo "  SAE控制台: https://sae.console.aliyun.com"
    echo ""
    echo "=================================================="
}

# 清理资源
cleanup() {
    log_step "清理所有资源..."
    
    cd "$SCRIPT_DIR"
    terraform destroy -auto-approve
    
    log_info "资源清理完成 ✓"
}

# 显示帮助
show_help() {
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  deploy      完整部署（默认）"
    echo "  infra       仅部署基础设施"
    echo "  image       仅构建和推送镜像"
    echo "  update      更新SAE应用"
    echo "  frontend    仅部署前端"
    echo "  cleanup     清理所有资源"
    echo "  help        显示帮助"
    echo ""
    echo "环境变量:"
    echo "  ALICLOUD_ACCESS_KEY   阿里云AccessKey ID（必需）"
    echo "  ALICLOUD_SECRET_KEY   阿里云AccessKey Secret（必需）"
    echo "  REGION                部署区域（默认: cn-hangzhou）"
    echo "  PROJECT_NAME          项目名称（默认: app-platform）"
}

# 主函数
main() {
    local command="${1:-deploy}"
    
    case "$command" in
        deploy)
            check_prerequisites
            check_credentials
            deploy_infrastructure
            build_and_push_image
            update_sae_application
            deploy_frontend
            init_database
            show_results
            ;;
        infra)
            check_prerequisites
            check_credentials
            deploy_infrastructure
            ;;
        image)
            check_prerequisites
            check_credentials
            build_and_push_image
            ;;
        update)
            check_prerequisites
            check_credentials
            build_and_push_image
            update_sae_application
            ;;
        frontend)
            check_prerequisites
            check_credentials
            deploy_frontend
            ;;
        cleanup)
            check_prerequisites
            check_credentials
            cleanup
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"

#!/bin/bash

BASE_DIR=/opt/goi_example
SERVICE_NAME="goi 示例项目"
SYSTEM_BIN=goi_example
SYSTEM_SERVICE=goi_example.service

function log() {
    message="[${SERVICE_NAME}]: $1 "
    echo -e "${message}" 2>&1 | tee -a ${BASE_DIR}/install.log
}

function check_root() {
    if [[ $EUID -ne 0 ]]; then
        log "错误: 此脚本必须以 root 权限运行"
        exit 1
    fi
}

function Install() {
    log "开始安装 ${SERVICE_NAME}..."

    chmod +x ./${SYSTEM_BIN}

    if [[ ! -e /usr/bin/${SYSTEM_BIN} ]]; then
        ln -sf ${BASE_DIR}/${SYSTEM_BIN} /usr/bin/${SYSTEM_BIN} || {
            log "错误: 无法创建软链接"
            exit 1
        }
    fi

    cp ${BASE_DIR}/${SYSTEM_SERVICE} /etc/systemd/system/

    # 重载系统服务
    systemctl daemon-reload 2>&1 | tee -a ${BASE_DIR}/install.log

    # 启用服务
    systemctl enable ${SYSTEM_SERVICE}

    log "启动 ${SERVICE_NAME} 服务"
    systemctl start ${SYSTEM_SERVICE} 2>&1 | tee -a ${BASE_DIR}/install.log || {
        log "错误: 无法启动服务"
        exit 1
    }

    # 检查服务状态
    for b in {1..10}
    do
        sleep 3
        service_status=$(systemctl status ${SYSTEM_SERVICE} 2>&1 | grep Active)
        if systemctl is-active --quiet ${SYSTEM_SERVICE}; then
            log "${SERVICE_NAME} 服务启动成功!"
            return 0
        else
            echo "$service_status" >> ${BASE_DIR}/install.log
            log "等待服务启动... (${b}/10)"
        fi
    done

    log "错误: 服务启动失败，请检查 ${BASE_DIR}/install.log"
    systemctl status ${SYSTEM_SERVICE} 2>&1 | tee -a ${BASE_DIR}/install.log
    exit 1
}

function Uninstall() {
    echo "开始卸载 ${SERVICE_NAME}..."
    
    systemctl stop ${SYSTEM_SERVICE} 2>/dev/null
    systemctl disable ${SYSTEM_SERVICE} 2>/dev/null
    
    sudo rm -rf /etc/systemd/system/${SYSTEM_SERVICE}
    sudo rm -rf /usr/bin/${SYSTEM_BIN}
    
    if [[ -d ${BASE_DIR} ]]; then
        read -p "是否删除 ${BASE_DIR} 目录? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf ${BASE_DIR}
            echo "已删除 ${BASE_DIR} 目录"
        fi
    fi
    
    systemctl daemon-reload
    echo "${SERVICE_NAME} 卸载完成"
    cd ../
}

function main() {
    check_root

    timedatectl set-timezone Asia/Shanghai

    if [[ "$1" == "uninstall" ]]; then
        Uninstall
    else
        Install
    fi
}

# 接收命令行参数
main "$@"

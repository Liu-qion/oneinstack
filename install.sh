#!/bin/bash
export PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:$PATH

LOGO="+----------------------------------------------------\n| one面板安装脚本\n| \n+----------------------------------------------------\n| Copyright © 2022-"$(date +%Y)" oneinstack All rights reserved.\n+----------------------------------------------------"
current_path=$(pwd)
ssh_port=$(cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}')
in_china=$(curl --retry 2 -m 10 -L https://www.qualcomm.cn/cdn-cgi/trace 2>/dev/null | grep -qx 'loc=CN' && echo "true" || echo "false")

# 检查系统
Prepare_System() {
    if [ $(whoami) != "root" ]; then
        error "请使用root用户运行安装命令(Please run the installation command using the root user)"
    fi

    if [ ${OS} == "unknown" ]; then
        error "系统不支持安装面板(The system does not support installing the one)"
    fi
    if [ ${ARCH} != "x86_64" ] && [ ${ARCH} != "aarch64" ]; then
        error "系统架构不支持安装面板(The system architecture does not support installing the one)"
    fi

    if [ ${ARCH} == "x86_64" ]; then
        if [ "$(cat /proc/cpuinfo | grep -c ssse3)" -lt "1" ]; then
            error "CPU至少需支持x86-64-v2指令集(CPU must support at least the x86-64-v2 instruction set)"
        fi
    fi

    kernel_version=$(uname -r | awk -F '.' '{print $1}')
    if [ "${kernel_version}" -lt "4" ]; then
        error "系统内核版本太低，请升级到4.x以上版本(The system kernel version is too low, please upgrade to version 4.x or above)"
    fi

    is_64bit=$(getconf LONG_BIT)
    if [ "${is_64bit}" != '64' ]; then
        error "请更换64位系统安装面板(Please switch to a 64-bit system to install the one)"
    fi

    if [ -f "${setup_path}/one/web" ]; then
        error "面板已安装，无需重复安装(one is already installed, no need to install again)"
    fi

    if ! id -u "www" >/dev/null 2>&1; then
        groupadd www
        useradd -s /sbin/nologin -g www www
    fi

    if [ ! -d ${setup_path} ]; then
        mkdir ${setup_path}
    fi

    timedatectl set-timezone Asia/Shanghai

    [ -s /etc/selinux/config ] && sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    setenforce 0 >/dev/null 2>&1

    ulimit -n 1048576
    echo 2147483584 >/proc/sys/fs/file-max
    soft_nofile_check=$(cat /etc/security/limits.conf | grep '^* soft nofile .*$')
    hard_nofile_check=$(cat /etc/security/limits.conf | grep '^* hard nofile .*$')
    soft_nproc_check=$(cat /etc/security/limits.conf | grep '^* soft nproc .*$')
    hard_nproc_check=$(cat /etc/security/limits.conf | grep '^* hard nproc .*$')
    fs_file_max_check=$(cat /etc/sysctl.conf | grep '^fs.file-max.*$')
    if [ "${soft_nofile_check}" == "" ]; then
        echo "* soft nofile 1048576" >>/etc/security/limits.conf
    fi
    if [ "${hard_nofile_check}" == "" ]; then
        echo "* hard nofile 1048576" >>/etc/security/limits.conf
    fi
    if [ "${soft_nproc_check}" == "" ]; then
        echo "* soft nproc 1048576" >>/etc/security/limits.conf
    fi
    if [ "${hard_nproc_check}" == "" ]; then
        echo "* hard nproc 1048576" >>/etc/security/limits.conf
    fi
    if [ "${fs_file_max_check}" == "" ]; then
        echo fs.file-max = 2147483584 >>/etc/sysctl.conf
    fi

    # 自动开启 BBR
    # 清理旧配置，防止干扰
    sed -i '/net.core.default_qdisc/d' /etc/sysctl.conf
    sed -i '/net.ipv4.tcp_congestion_control/d' /etc/sysctl.conf
    sysctl -p
    bbr_support_check=$(ls -l /lib/modules/*/kernel/net/ipv4 | grep -c tcp_bbr)
    bbr_open_check=$(sysctl net.ipv4.tcp_congestion_control | grep -c bbr)
    if [ "${bbr_support_check}" != "0" ] && [ "${bbr_open_check}" == "0" ]; then
        qdisc=$(sysctl net.core.default_qdisc | awk '{print $3}')
        # cake 是 fq_codel 的升级版，优先使用
        if cat /boot/config-$(uname -r) | grep CONFIG_NET_SCH_CAKE | grep -q "="; then
            qdisc="cake"
        elif cat /boot/config-$(uname -r) | grep CONFIG_NET_SCH_FQ_CODEL | grep -q "="; then
            qdisc="fq_codel"
        elif cat /boot/config-$(uname -r) | grep CONFIG_NET_SCH_FQ_PIE | grep -q "="; then
            qdisc="fq_pie"
        elif cat /boot/config-$(uname -r) | grep CONFIG_NET_SCH_FQ | grep -q "="; then
            qdisc="fq"
        fi
        echo "net.core.default_qdisc=${qdisc}" >>/etc/sysctl.conf
        echo "net.ipv4.tcp_congestion_control=bbr" >>/etc/sysctl.conf
        sysctl -p
    fi

    if [ ${OS} == "rhel" ]; then
        if ${in_china}; then
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^#baseurl=http://dl.rockylinux.org/$contentdir|baseurl=https://mirrors.tencent.com/rocky|g' \
                -e 's|^# baseurl=http://dl.rockylinux.org/$contentdir|baseurl=https://mirrors.tencent.com/rocky|g' \
                -i.bak \
                /etc/yum.repos.d/[Rr]ocky*.repo >/dev/null 2>&1
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^#baseurl=https://repo.almalinux.org|baseurl=https://mirrors.tencent.com|g' \
                -e 's|^# baseurl=https://repo.almalinux.org|baseurl=https://mirrors.tencent.com|g' \
                -i.bak \
                /etc/yum.repos.d/[Aa]lmalinux*.repo >/dev/null 2>&1
            sed -e 's|^mirrorlist=|#mirrorlist=|g' \
                -e 's|^#baseurl=http://mirror.centos.org/$contentdir|baseurl=https://mirrors.tencent.com/centos-stream|g' \
                -e 's|^# baseurl=http://mirror.centos.org/$contentdir|baseurl=https://mirrors.tencent.com/centos-stream|g' \
                -i.bak \
                /etc/yum.repos.d/[Cc]ent*.repo >/dev/null 2>&1
        fi
        dnf makecache -y
        dnf install dnf-plugins-core -y
        dnf config-manager --set-enabled epel
        if ${in_china}; then
            sed -i 's|^#baseurl=https://download.example/pub|baseurl=https://mirrors.tencent.com|' /etc/yum.repos.d/epel* >/dev/null 2>&1
            sed -i 's|^# baseurl=https://download.example/pub|baseurl=https://mirrors.tencent.com|' /etc/yum.repos.d/epel* >/dev/null 2>&1
            sed -i 's|^metalink|#metalink|' /etc/yum.repos.d/epel* >/dev/null 2>&1
            dnf makecache -y
        fi
        # EL 9
        dnf config-manager --set-enabled crb
        dnf install epel-release epel-next-release -y
        # 部分系统可能没有这两个包，需要手动安装
        # 对于 openEuler 这种大改的系统，下载会 404，这没有影响
        if [ "$?" != "0" ]; then
            if ${in_china}; then
                dnf install -y https://mirrors.tencent.com/epel/epel-release-latest-$(rpm -E %{rhel}).noarch.rpm
                dnf install -y https://mirrors.tencent.com/epel/epel-next-release-latest-$(rpm -E %{rhel}).noarch.rpm
            else
                dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-$(rpm -E %{rhel}).noarch.rpm
                dnf install -y https://dl.fedoraproject.org/pub/epel/epel-next-release-latest-$(rpm -E %{rhel}).noarch.rpm
            fi
        fi
        # Rocky Linux
        /usr/bin/crb enable >/dev/null 2>&1
        # openEuler
        if [ -f /etc/openEuler-release ]; then
            # 清理旧配置
            grep -rl '^baseurl=https://repo.oepkgs.net' /etc/yum.repos.d/ | xargs -I {} rm -f {}
            dnf config-manager --add-repo https://repo.oepkgs.net/openeuler/rpm/$(awk '{print $1"-"$3"-"$4}' /etc/openEuler-release | sed 's/[()]//g')/extras/$(uname -m)/
            oe_version=$(awk '{print $3}' /etc/openEuler-release | cut -d '.' -f 1)
            case ${oe_version} in
            22)
                dnf config-manager --add-repo https://repo.oepkgs.net/openeuler/rpm/$(awk '{print $1"-"$3"-"$4}' /etc/openEuler-release | sed 's/[()]//g')/extras/$(uname -m)/
                dnf config-manager --add-repo https://repo.oepkgs.net/openeuler/rpm/$(awk '{print $1"-"$3"-"$4}' /etc/openEuler-release | sed 's/[()]//g')/compatible/f33/$(uname -m)/
                ;;
            24)
                # openEuler 24 目前没有 p7zip-plugins，等他们自己修复
                error "openEuler 24 暂不支持安装面板(one installation is not supported on openEuler 24)"
                ;;
            *)
                error "不支持的 openEuler 版本(Unsupported openEuler version)"
                ;;
            esac
            # 禁用gpcheck，这仓库缺签名
            grep -rl '^baseurl=https://repo.oepkgs.net' /etc/yum.repos.d/ | xargs -I {} sh -c 'echo "gpgcheck=0" >> "{}"'
        fi
        dnf makecache -y
        dnf install -y bash curl wget zip unzip tar p7zip p7zip-plugins git jq git-core dos2unix make sudo
    elif [ ${OS} == "debian" ] || [ ${OS} == "ubuntu" ]; then
        if ${in_china}; then
            # Debian
            sed -i 's/deb.debian.org/mirrors.tencent.com/g' /etc/apt/sources.list >/dev/null 2>&1
            sed -i 's/deb.debian.org/mirrors.tencent.com/g' /etc/apt/sources.list.d/debian.sources >/dev/null 2>&1
            sed -i -e 's|security.debian.org/\? |security.debian.org/debian-security |g' \
                -e 's|security.debian.org|mirrors.tencent.com|g' \
                -e 's|deb.debian.org/debian-security|mirrors.tencent.com/debian-security|g' \
                /etc/apt/sources.list >/dev/null 2>&1
            # Ubuntu
            sed -i 's@//.*archive.ubuntu.com@//mirrors.tencent.com@g' /etc/apt/sources.list >/dev/null 2>&1
            sed -i 's@//.*archive.ubuntu.com@//mirrors.tencent.com@g' /etc/apt/sources.list.d/ubuntu.sources >/dev/null 2>&1
            sed -i 's/security.ubuntu.com/mirrors.tencent.com/g' /etc/apt/sources.list >/dev/null 2>&1
            sed -i 's/security.ubuntu.com/mirrors.tencent.com/g' /etc/apt/sources.list.d/ubuntu.sources >/dev/null 2>&1
        fi
        apt-get update -y
        apt-get install -y bash curl wget zip unzip tar p7zip p7zip-full git jq git dos2unix make sudo
    fi
    if [ "$?" != "0" ]; then
        error "安装面板依赖软件失败(Installation of one dependency software failed)"
    fi
}



Install_One() {
    local url="https://github.com/jimbirthday/oneinstack/releases/download/test/one"
    local dest="/usr/local/one/one"
    local timeout=30  # 设置下载超时时间为30秒

    # 创建目录
    if ! mkdir -p /usr/local/one; then
        error "创建目录失败(Failed to create directory)"
        return 1
    fi

    # 下载 one 二进制文件，设置超时和等待时间
    if ! curl --max-time $timeout -L -o "$dest" "$url"; then
        error "下载 one 二进制文件失败(Failed to download the one binary)"
        return 1
    fi

    # 赋予执行权限
    if ! chmod +x "$dest"; then
        error "设置执行权限失败(Failed to set execute permissions)"
        return 1
    fi

    # 创建 systemd 服务文件
    local service_file="/etc/systemd/system/one.service"
    cat <<EOF > "$service_file"
[Unit]
Description=One Service

[Service]
ExecStart=$dest server start
ExecStop=$dest server stop

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载 systemd 管理器配置
    if ! systemctl daemon-reload; then
        error "重新加载 systemd 配置失败(Failed to reload systemd configuration)"
        return 1
    fi

    # 启用并启动服务
    if ! systemctl enable one || ! systemctl start one; then
        error "启用或启动 one 服务失败(Failed to enable or start the one service)"
        return 1
    fi

    # 添加到 PATH
    if ! grep -q "/usr/local/one" /etc/profile; then
        echo "export PATH=\$PATH:/usr/local/one" >> /etc/profile
        echo "已将 /usr/local/one 添加到 PATH(Added /usr/local/one to PATH)"
    fi

    echo "one 二进制文件安装成功并已添加到 systemctl(One binary installed successfully and added to systemctl)"
}

clear
echo -e $LOGO


# 安装确认
read -p "面板将安装至 ${setup_path} 目录，请输入 y 并回车以开始安装 (Enter 'y' to start installation): " install
if [ "$install" != 'y' ]; then
    echo "输入不正确，已退出安装。"
    echo "Incorrect input, installation has been exited."
    exit
fi

clear
echo -e $LOGO
echo "安装面板依赖软件（如报错请检查软件源是否正常）"
echo "Installing one dependency software (if error, please check the software source)"
echo -e $HR
sleep 1s
Prepare_System

echo -e $LOGO
echo "安装面板运行环境（视网络情况可能需要较长时间）"
echo "Installing the one running environment (may take a long time depending on the network)"
echo -e $HR
sleep 1s
Install_One

clear
echo -e $LOGO
echo '面板安装成功！'
echo -e $HR
# 查看 one 启动日志
echo "启动日志(Startup logs):"
journalctl -u one --no-pager

cd ${current_path}
rm -f install.sh

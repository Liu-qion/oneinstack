package software

import (
	"fmt"
	"oneinstack/utils"
	"oneinstack/web/input"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var mysql55 = `
#!/bin/bash

# 默认值
ROOT_PASSWORD=""

# 函数：显示帮助信息
usage() {
  echo "用法: $0 -p <root_password>"
  exit 1
}

# 解析命令行参数
while getopts "p:" opt; do
  case "$opt" in
    p) ROOT_PASSWORD="$OPTARG" ;;  # 设置 root 密码
    *) usage ;;  # 不支持的选项
  esac
done

# 检查是否传入了 root 密码
if [ -z "$ROOT_PASSWORD" ]; then
  echo "请通过 -p 参数传入 MySQL root 密码，例如：$0 -p <root_password>"
  exit 1
fi

# 确保脚本以 root 用户执行
if [ "$(id -u)" -ne 0 ]; then
  echo "请使用 root 用户执行该脚本"
  exit 1
fi

# 定义函数来检测和选择包管理器
setup_package_manager() {
  if command -v apt-get > /dev/null 2>&1; then
    PACKAGE_MANAGER="apt-get"
  elif command -v yum > /dev/null 2>&1; then
    PACKAGE_MANAGER="yum"
  elif command -v dnf > /dev/null 2>&1; then
    PACKAGE_MANAGER="dnf"
  else
    echo "不支持的包管理器"
    exit 1
  fi
}

# 更新系统包
update_packages() {
  echo "更新系统包..."
  $PACKAGE_MANAGER update -y
}

# 安装依赖项
install_dependencies() {
  echo "安装 MySQL 所需的依赖..."
  case $PACKAGE_MANAGER in
    apt-get)
      $PACKAGE_MANAGER install -y build-essential cmake libncurses5-dev libssl-dev libboost-all-dev bison
      ;;
    yum|dnf)
      $PACKAGE_MANAGER install -y gcc gcc-c++ make cmake ncurses-devel openssl-devel boost-devel bison
      ;;
  esac
}

# 下载并解压 MySQL 源码包
download_and_extract_mysql() {
  local MYSQL_VERSION="mysql-5.5.62"
  local MYSQL_TAR="$MYSQL_VERSION.tar.gz"
  
  echo "下载 MySQL 5.5 源码包..."
  wget https://dev.mysql.com/get/Downloads/MySQL-5.5/$MYSQL_TAR
  
  echo "解压 MySQL 源码包..."
  tar -xvzf $MYSQL_TAR
  cd $MYSQL_VERSION
}

# 编译并安装 MySQL
compile_and_install_mysql() {
  echo "创建 MySQL 安装目录..."
  sudo mkdir -p /usr/local/mysql
  
  echo "编译 MySQL..."
  cmake . -DCMAKE_INSTALL_PREFIX=/usr/local/mysql \
          -DMYSQL_DATADIR=/usr/local/mysql/data \
          -DDEFAULT_CHARSET=utf8 \
          -DDEFAULT_COLLATION=utf8_general_ci \
          -DWITH_INNOBASE_STORAGE_ENGINE=1 \
          -DWITH_PARTITION_STORAGE_ENGINE=1 \
          -DWITH_FEDERATED_STORAGE_ENGINE=1 \
          -DWITH_BLACKHOLE_STORAGE_ENGINE=1 \
          -DWITH_MYISAM_STORAGE_ENGINE=1 \
          -DWITH_ARCHIVE_STORAGE_ENGINE=1 \
          -DEXTRA_CHARSETS=all \
          -DENABLED_LOCAL_INFILE=1
  
  echo "安装 MySQL..."
  make -j$(nproc)
  sudo make install
}

# 创建 MySQL 用户并设置权限
create_mysql_user_and_set_permissions() {
  echo "创建 MySQL 用户..."
  sudo useradd -r -s /bin/false mysql
  
  echo "设置 MySQL 目录权限..."
  sudo chown -R mysql:mysql /usr/local/mysql
  sudo chown -R mysql:mysql /usr/local/mysql/data
}

# 初始化数据库
initialize_database() {
  echo "初始化 MySQL 数据库..."
  sudo /usr/local/mysql/scripts/mysql_install_db --user=mysql --basedir=/usr/local/mysql --datadir=/usr/local/mysql/data
}

# 配置 MySQL 服务
configure_mysql_service() {
  echo "配置 MySQL 服务..."
  if command -v systemctl > /dev/null 2>&1; then
    sudo cp /usr/local/mysql/support-files/mysql.server /etc/init.d/mysql
    sudo chmod +x /etc/init.d/mysql
    sudo systemctl daemon-reload
    sudo systemctl enable mysql
  else
    sudo chkconfig --add mysql
    sudo chkconfig mysql on
  fi
}

# 设置环境变量
set_env_vars() {
  echo 'export PATH=$PATH:/usr/local/mysql/bin' >> ~/.bashrc
  source ~/.bashrc
}

# 启动 MySQL 服务并设置 root 密码
start_mysql_and_set_root_password() {
  echo "启动 MySQL 服务..."
  sudo service mysql start || sudo /etc/init.d/mysql start
  
  echo "设置 MySQL root 密码..."
  sudo /usr/local/mysql/bin/mysqladmin -u root password "$ROOT_PASSWORD"
}

# 主函数
main() {
  setup_package_manager
  update_packages
  install_dependencies
  download_and_extract_mysql
  compile_and_install_mysql
  create_mysql_user_and_set_permissions
  initialize_database
  configure_mysql_service
  set_env_vars
  start_mysql_and_set_root_password

  echo "MySQL 5.5 安装完成，root 密码已设置为：$ROOT_PASSWORD"
}

main

`

var mysql57 = `
#!/bin/bash

# MySQL版本
MYSQL_VERSION="5.7.40"
MYSQL_DOWNLOAD_URL="https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-${MYSQL_VERSION}.tar.gz"

# Boost版本
BOOST_VERSION="1.59.0"
BOOST_DOWNLOAD_URL="https://zenlayer.dl.sourceforge.net/project/boost/boost/1.59.0/boost_1_59_0.tar.gz?viasf=1"
BOOST_INSTALL_DIR="/usr/local/boost_${BOOST_VERSION}"

# OpenSSL版本
OPENSSL_VERSION="1.1.1u"
OPENSSL_DOWNLOAD_URL="https://www.openssl.org/source/openssl-${OPENSSL_VERSION}.tar.gz"
OPENSSL_INSTALL_DIR="/usr/local/openssl"

# Zlib版本
ZLIB_VERSION="1.2.13"
ZLIB_DOWNLOAD_URL="https://www.zlib.net/fossils/zlib-${ZLIB_VERSION}.tar.gz"
ZLIB_INSTALL_DIR="/usr/local/zlib"

# 安装目录
MYSQL_INSTALL_DIR="/usr/local/mysql"
MYSQL_DATA_DIR="/data/mysql"

# 检查是否为 root 用户
if [ "$(id -u)" != "0" ]; then
  echo "请以 root 用户运行该脚本"
  exit 1
fi

# 检测系统类型
Detect_OS() {
  if [ -f /etc/redhat-release ]; then
    OS="CentOS"
    PM="yum"
  elif [ -f /etc/debian_version ]; then
    OS="Debian"
    PM="apt"
  else
    echo "不支持的操作系统"
    exit 1
  fi
}

# 安装依赖
Install_Dependencies() {
  echo "安装必要依赖包..."
  if [ "$PM" == "yum" ]; then
    yum -y install gcc gcc-c++ cmake ncurses-devel bison wget perl make libaio-devel
  elif [ "$PM" == "apt" ]; then
    apt update
    apt -y install build-essential cmake libncurses5-dev libaio-dev bison wget
  fi
}

Install_Zlib() {
  if [ ! -d "${ZLIB_INSTALL_DIR}" ]; then
    echo "安装 Zlib ${ZLIB_VERSION}..."
    wget -c "${ZLIB_DOWNLOAD_URL}" || { echo "下载 Zlib 失败"; exit 1; }
    tar xf "zlib-${ZLIB_VERSION}.tar.gz"
    cd "zlib-${ZLIB_VERSION}" || exit
    ./configure --prefix=${ZLIB_INSTALL_DIR}
    make -j"$(nproc)"
    make install
    cd ..
  else
    echo "Zlib ${ZLIB_VERSION} 已安装，路径：${ZLIB_INSTALL_DIR}"
  fi
}

# 安装 Boost
Install_Boost() {
  if [ ! -d "${BOOST_INSTALL_DIR}" ]; then
    echo "安装 Boost ${BOOST_VERSION}..."
    wget -c "${BOOST_DOWNLOAD_URL}" || { echo "下载 Boost 失败"; exit 1; }
    tar xf "boost_1_59_0.tar.gz"
    cd "boost_1_59_0" || exit
    ./bootstrap.sh --prefix="${BOOST_INSTALL_DIR}"
    ./b2 install
    cd ..
  else
    echo "Boost ${BOOST_VERSION} 已安装，路径：${BOOST_INSTALL_DIR}"
  fi
}

# 安装 OpenSSL
Install_OpenSSL() {
  if [ ! -d "${OPENSSL_INSTALL_DIR}" ]; then
    echo "安装 OpenSSL ${OPENSSL_VERSION}..."
    wget -c "${OPENSSL_DOWNLOAD_URL}" || { echo "下载 OpenSSL 失败"; exit 1; }
    tar xf "openssl-${OPENSSL_VERSION}.tar.gz"
    cd "openssl-${OPENSSL_VERSION}" || exit
    ./config --prefix=${OPENSSL_INSTALL_DIR} no-ssl2 no-ssl3
    make -j"$(nproc)"
    make install
    cd ..
  else
    echo "OpenSSL ${OPENSSL_VERSION} 已安装，路径：${OPENSSL_INSTALL_DIR}"
  fi
}

# 下载 MySQL 源码
Download_MySQL() {
  if [ ! -f "mysql-${MYSQL_VERSION}.tar.gz" ]; then
    echo "下载 MySQL ${MYSQL_VERSION}..."
    wget -c "${MYSQL_DOWNLOAD_URL}" || { echo "下载 MySQL 失败"; exit 1; }
  fi
  tar xf "mysql-${MYSQL_VERSION}.tar.gz"
}

# 编译并安装 MySQL
Install_MySQL() {
  cd "mysql-${MYSQL_VERSION}" || exit
  cmake . \
  -DCMAKE_INSTALL_PREFIX=${MYSQL_INSTALL_DIR} \
  -DMYSQL_DATADIR=${MYSQL_DATA_DIR} \
  -DWITH_INNOBASE_STORAGE_ENGINE=1 \
  -DWITH_ARCHIVE_STORAGE_ENGINE=1 \
  -DWITH_BLACKHOLE_STORAGE_ENGINE=1 \
  -DWITH_FEDERATED_STORAGE_ENGINE=1 \
  -DWITH_PARTITION_STORAGE_ENGINE=1 \
  -DENABLED_LOCAL_INFILE=1 \
  -DWITH_SSL=${OPENSSL_INSTALL_DIR} \
  -DWITH_ZLIB=${ZLIB_INSTALL_DIR} \
  -DWITH_BOOST=${BOOST_INSTALL_DIR} \
  -DCMAKE_C_FLAGS="-fPIC" \
  -DDEFAULT_CHARSET=utf8 \
  -DDEFAULT_COLLATION=utf8_general_ci \
  -DMYSQL_TCP_PORT=3306 \
  -DMYSQL_UNIX_ADDR=/tmp/mysql.sock

  make -j"$(nproc)"
  make install
  cd ..
}

# 初始化 MySQL
Initialize_MySQL() {
  echo "初始化 MySQL 数据目录..."
  ${MYSQL_INSTALL_DIR}/bin/mysqld --initialize-insecure --user=mysql --basedir=${MYSQL_INSTALL_DIR} --datadir=${MYSQL_DATA_DIR}
  echo "MySQL 初始化完成"
}

# 设置环境变量
Configure_Environment() {
  echo "配置环境变量..."
  if ! grep -q "${MYSQL_INSTALL_DIR}/bin" /etc/profile; then
    echo "export PATH=\$PATH:${MYSQL_INSTALL_DIR}/bin" >> /etc/profile
    source /etc/profile
  fi
  echo "环境变量配置完成"
}

# 启动 MySQL
Start_MySQL() {
  echo "启动 MySQL 服务..."
  ${MYSQL_INSTALL_DIR}/bin/mysqld_safe --user=mysql &
  echo "MySQL 启动完成"
}

# 主函数
Main() {
  Detect_OS
  Install_Dependencies
  Install_OpenSSL
  Install_Boost
  Install_Zlib
  Download_MySQL
  Install_MySQL
  Initialize_MySQL
  Configure_Environment
  Start_MySQL
  echo "MySQL ${MYSQL_VERSION} 安装完成"
}

# 执行主函数
Main



`

var mysql80 = ``

var redis = `

#!/bin/bash

# 脚本名称：install_redis.sh
# 用途：从源码安装 Redis 6 或 Redis 7

# 检查是否有 root 权限
if [[ $EUID -ne 0 ]]; then
   echo "请使用 root 权限运行此脚本" 
   exit 1
fi

# 检查参数是否传递正确
if [ -z "$1" ]; then
    echo "使用方法：$0 {6|7}"
    echo "6: 安装 Redis 6"
    echo "7: 安装 Redis 7"
    exit 1
fi

# 安装依赖
echo "正在安装依赖..."
apt-get update && apt-get install -y build-essential tcl wget

# 设置版本
VERSION="$1"

# 下载 Redis 源码
if [ "$VERSION" == "6" ]; then
    echo "正在下载 Redis 6.x 源码..."
    wget https://mirrors.huaweicloud.com/redis/redis-6.2.0.tar.gz -O /tmp/redis-6.tar.gz
    tar -zxvf /tmp/redis-6.tar.gz -C /tmp
    cd /tmp/redis-6.2.0
elif [ "$VERSION" == "7" ]; then
    echo "正在下载 Redis 7.x 源码..."
    wget https://mirrors.huaweicloud.com/redis/redis-7.0.5.tar.gz -O /tmp/redis-7.tar.gz
    tar -zxvf /tmp/redis-7.tar.gz -C /tmp
    cd /tmp/redis-7.0.5
else
    echo "无效的版本号：$VERSION。请指定 6 或 7"
    exit 1
fi

# 编译和安装 Redis
echo "正在编译 Redis..."
make
make install

# 配置 Redis
echo "正在配置 Redis..."
cp /tmp/redis-*/redis.conf /etc/redis.conf

# 创建 Redis 用户
useradd -r -s /bin/false redis

# 创建 Redis 数据目录
mkdir /var/lib/redis
chown redis:redis /var/lib/redis

# 创建 Redis 启动脚本
echo "[Unit]
Description=Redis In-Memory Data Store
After=network.target

[Service]
ExecStart=/usr/local/bin/redis-server /etc/redis.conf
ExecStop=/usr/local/bin/redis-cli shutdown
User=redis
Group=redis
WorkingDirectory=/var/lib/redis
Restart=always

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/redis.service

# 设置 Redis 服务为开机自启
systemctl enable redis
systemctl start redis

# 清理安装文件
rm -rf /tmp/redis-*

echo "Redis $VERSION 安装完成！"

`

var nginx = `
#!/bin/bash

# 脚本名称：install_nginx.sh
# 用途：从源码安装 Nginx（使用国内源下载）

# 检查是否有 root 权限
if [[ $EUID -ne 0 ]]; then
   echo "请使用 root 权限运行此脚本" 
   exit 1
fi

# 安装依赖
echo "正在安装依赖..."
apt-get update && apt-get install -y build-essential libpcre3 libpcre3-dev libssl-dev zlib1g-dev wget

# 下载 Nginx 源码（使用阿里云镜像）
NGINX_VERSION="1.24.0"  # 指定 Nginx 版本
echo "正在从国内源下载 Nginx $NGINX_VERSION 源码..."
wget https://mirrors.huaweicloud.com/nginx/nginx-$NGINX_VERSION.tar.gz -O /tmp/nginx.tar.gz
tar -zxvf /tmp/nginx.tar.gz -C /tmp
cd /tmp/nginx-$NGINX_VERSION

# 编译和安装 Nginx
echo "正在编译 Nginx..."
./configure --prefix=/usr/local/nginx --with-http_ssl_module --with-http_v2_module --with-pcre
make
make install

# 配置 Nginx 为系统服务
echo "正在配置 Nginx 服务..."

# 创建 Nginx 启动脚本
echo "[Unit]
Description=NGINX
After=network.target

[Service]
ExecStart=/usr/local/nginx/sbin/nginx
ExecReload=/usr/local/nginx/sbin/nginx -s reload
ExecStop=/usr/local/nginx/sbin/nginx -s stop
PIDFile=/usr/local/nginx/logs/nginx.pid
User=www-data
Group=www-data
WorkingDirectory=/usr/local/nginx

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/nginx.service

# 创建必要的文件夹
echo "创建日志和临时目录..."
mkdir -p /usr/local/nginx/logs
mkdir -p /usr/local/nginx/client_body_temp
mkdir -p /usr/local/nginx/proxy_temp
mkdir -p /usr/local/nginx/fastcgi_temp
mkdir -p /usr/local/nginx/scgi_temp
mkdir -p /usr/local/nginx/uwsgi_temp

# 设置文件夹权限
chown -R www-data:www-data /usr/local/nginx
chmod -R 755 /usr/local/nginx

# 设置 Nginx 服务为开机自启
systemctl enable nginx

# 启动 Nginx 服务
systemctl start nginx

# 配置 nginx 环境变量
echo "正在将 nginx 添加到环境变量中..."
ln -sf /usr/local/nginx/sbin/nginx /usr/bin/nginx

# 清理安装文件
rm -rf /tmp/nginx*

# 输出安装信息
echo "Nginx $NGINX_VERSION 安装完成！"
echo "访问 http://<your_server_ip> 来查看 Nginx 默认页面"

`

type InstallOPI interface {
	Install() (string, error)
}

type InstallOP struct {
	Params     *input.InstallationParams
	BashParams *input.InstallParams
}

func NewInstallOP(p *input.InstallParams) (InstallOP, error) {
	//params := buildIParams(p)
	return InstallOP{Params: nil, BashParams: p}, nil
}

func buildIParams(p *input.InstallParams) *input.InstallationParams {
	ps := &input.InstallationParams{}
	switch p.Key {
	case "webserver":
		ps.NginxOption = p.Version
	case "db":
		ps.DBOption = p.Version
		ps.DBRootPWD = p.Pwd
	case "redis":
		ps.Redis = true
	case "php":
		ps.PHPOption = p.Version
	case "java":
		ps.JDKOption = p.Version
	}
	return ps
}

func (ps InstallOP) Install() (string, error) {
	bash := ""
	switch ps.BashParams.Key {
	case "webserver":
		bash = nginx
	case "db":
		if ps.BashParams.Version == "5.5" {
			bash = mysql55
		}
		if ps.BashParams.Version == "5.7" {
			bash = mysql57
		}
		if ps.BashParams.Version == "8.0" {
			bash = mysql80
		}
	case "redis":
		bash = redis
	case "php":
	case "java":
	default:
		return "", fmt.Errorf("未知的软件类型")
	}
	fn, err := createShScript(bash, ps.BashParams.Key+ps.BashParams.Version+".sh")
	if err != nil {
		return "", err

	}

	switch ps.BashParams.Key {
	case "webserver":
		return executeShScript(fn)
	case "db":
		if ps.BashParams.Version == "5.5" {
			return executeShScript(fn, "-p", ps.BashParams.Pwd)
		}
		if ps.BashParams.Version == "5.7" {
			return executeShScript(fn, "-p", ps.BashParams.Pwd)
		}
		if ps.BashParams.Version == "8.0" {
			return executeShScript(fn, "-p", ps.BashParams.Pwd)
		}
		return "", fmt.Errorf("未知的db类型")
	case "redis":
		if ps.BashParams.Version == "6.2.0" {
			return executeShScript(fn, "6")
		}
		if ps.BashParams.Version == "7.0.5" {
			return executeShScript(fn, "7")
		}
		return "", fmt.Errorf("未知的redis类型")
	case "php":
		return "", nil
	case "java":
		return "", nil
	default:
		return "", fmt.Errorf("未知的软件类型")
	}
}

// createShScript 将字符串内容保存为.sh脚本文件，如果文件已存在则覆盖
func createShScript(scriptContent, filename string) (string, error) {
	// 打开文件，如果文件不存在则创建，权限设置为可读可写可执行
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 写入脚本内容
	_, err = file.WriteString(scriptContent)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	// 打印成功信息
	fmt.Printf("脚本已保存为 %s\n", filename)
	return filename, nil
}

// executeShScript 执行指定的脚本文件，并支持传递命令行参数
func executeShScript(scriptName string, args ...string) (string, error) {
	// 拼接完整的命令：bash scriptName args...
	cmdArgs := append([]string{scriptName}, args...)
	cmd := exec.Command("bash", cmdArgs...)

	logFileName := "install_" + time.Now().Format("2006-01-02_15-04-05") + ".log"
	logFile, err := os.Create(logFileName)
	if err != nil {
		return "", fmt.Errorf("无法创建日志文件: %v", err)
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	go func() {
		err = cmd.Wait()
		if err != nil {
			fmt.Println("cmd wait err:" + fmt.Sprintf("%v", err))
		}
	}()
	return logFileName, nil
}

func runInstall(params *input.InstallationParams) (string, error) {
	err := downloadshell()
	if err != nil {
		return "", err
	}

	// 构建命令行参数列表
	cmdArgs := params.BuildCmdArgs()
	argsWithSudo := append([]string{"./oneinstack/oneinstack/install.sh"}, cmdArgs...)

	// 添加执行权限
	dirPath := "./oneinstack/oneinstack/include"
	err = utils.SetExecPermissions(dirPath)
	if err != nil {
		return "", fmt.Errorf("设置 include 目录下文件的执行权限失败: %v", err)
	}

	scriptPath := "./oneinstack/oneinstack/install.sh"
	err = os.Chmod(scriptPath, 0755)
	if err != nil {
		return "", fmt.Errorf("无法设置脚本执行权限: %v", err)
	}

	cmdInstall := exec.Command("sudo", argsWithSudo...)

	logFileName := "install_" + time.Now().Format("2006-01-02_15-04-05") + ".log"
	logFile, err := os.Create(logFileName)
	if err != nil {
		return "", fmt.Errorf("无法创建日志文件: %v", err)
	}
	defer logFile.Close()

	cmdInstall.Stdout = logFile
	cmdInstall.Stderr = logFile
	err = cmdInstall.Start()
	if err != nil {
		return "", err
	}
	go func() {
		err = cmdInstall.Wait()
		if err != nil {
			fmt.Println("cmd wait err:" + fmt.Sprintf("%v", err))
		}
	}()

	return logFileName, nil
}

// checkIfFileExists 检查文件是否存在。
func checkIfFileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// downloadOneInstack 如果 oneinstack.tar.gz 不存在则下载它。
func downloadshell() error {
	tarFilePath := filepath.Join(".", "oneinstack.tar.gz")
	if !checkIfFileExists(tarFilePath) {
		fmt.Println("oneinstack.tar.gz does not exist. Downloading...")
		err := utils.DownloadFile("https://mirrors.oneinstack.com/oneinstack.tar.gz", tarFilePath)
		if err != nil {
			return err
		}
		fmt.Printf("Download completed.\n")
	} else {
		fmt.Println("oneinstack.tar.gz already exists, skipping download.")
	}
	return utils.DecompressTarGz(tarFilePath, filepath.Join(".", "oneinstack"))
}

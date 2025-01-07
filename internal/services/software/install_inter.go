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

# 检查是否传入密码
if [ -z "$1" ]; then
  echo "请通过参数传入 MySQL root 密码，如：$0 <root_password>"
  exit 1
fi

ROOT_PASSWORD=$1

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

var mysql56 = ``
var mysql57 = ``
var mysql80 = ``

type InstallOPI interface {
	Install() (string, error)
}

type InstallOP struct {
	Params *input.InstallationParams
}

func NewInstallOP(p *input.InstallParams) (InstallOP, error) {
	params := buildIParams(p)
	return InstallOP{Params: params}, nil
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

func (s InstallOP) Install() (string, error) {
	return runInstall(s.Params)
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

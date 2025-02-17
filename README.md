<h1 align=center>Oneinstack</h1>
<h2 align=center><a href="">Demo Live</a></h1>

[![GitHub forks](https://img.shields.io/github/forks/guangzhengli/k8s-tutorials)](https://github.com/guangzhengli/k8s-tutorials/network)[![GitHub stars](https://img.shields.io/github/stars/guangzhengli/k8s-tutorials)](https://github.com/guangzhengli/k8s-tutorials/stargazers)[![GitHub issues](https://img.shields.io/github/issues/guangzhengli/k8s-tutorials)](https://github.com/guangzhengli/k8s-tutorials/issues)[![GitHub license](https://img.shields.io/github/license/guangzhengli/k8s-tutorials)](https://github.com/guangzhengli/k8s-tutorials/blob/main/LICENSE)![Docker Pulls](https://img.shields.io/docker/pulls/guangzhengli/hellok8s)

<h4 align=center>🌈 Oneinstack</h4>

xxxx

这里是文档的索引：
* [准备工作](docs/pre.md)
* [container](docs/container.md)
* [pod](docs/pod.md)


# Oneinstack

## 准备工作

在开始本教程之前，需要配置好本地环境，以下是需要安装的依赖和包。

### 安装 docker

首先我们需要安装 `docker` 来打包镜像，如果你本地已经安装了 `docker`，那么你可以选择跳过这一小节。

#### 推荐安装方法

目前使用 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 来安装 docker 还是最简单的方案，打开官网下载对应你电脑操作系统的包即可 (https://www.docker.com/products/docker-desktop/)，

当安装完成后，可以通过 `docker run hello-world` 来快速校验是否安装成功！

#### 其它安装方法

目前  Docker 公司宣布  [Docker Desktop](https://www.docker.com/products/docker-desktop/) 只对个人开发者或者小型团体免费 (2021年起对大型公司不再免费)，所以如果你不能通过  [Docker Desktop](https://www.docker.com/products/docker-desktop/) 的方式下载安装 `docker`，可以参考 [这篇文章](https://dhwaneetbhatt.com/blog/run-docker-without-docker-desktop-on-macos) 只安装  [Docker CLI](https://github.com/docker/cli)。

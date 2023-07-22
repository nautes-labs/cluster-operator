# Cluster Operator
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![golang](https://img.shields.io/badge/golang-v1.20.5-brightgreen)](https://go.dev/doc/install)
[![version](https://img.shields.io/badge/version-v0.3.6-green)]()

Cluster Operator 项目提供了一个用于调谐 Cluster 资源事件的 Controller，调谐内容主要是管理 Cluster 资源所声明的 Kubernetes 集群的密钥信息，使参与集群管理的其他组件可以从租户的密钥管理系统中正确获取到集群的密钥。

如果从集群形态进行区分，Cluster 资源可以分为物理集群和虚拟集群；如果从集群的用途进行区分，Cluster 资源可以分为宿主集群和运行时集群。针对不同的集群属性以及操作，Controller 的具体工作内容有所差异。

## 功能简介

### 注册集群

|            | 物理集群                                                     | 虚拟集群                                                     |
| ---------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 宿主集群   | 不处理                                                       | 不支持                                                       |
| 运行时集群 | 在密钥管理系统中创建物理集群的认证、将物理集群的密钥授权给 Runtime Operator、将租户配置库的密钥授权给物理集群的 Argo Operator、收集集群中可以作为环境入口的k8s服务的端口号 | 将从虚拟集群所属的宿主集群中获取的密钥写入密钥管理系统、在密钥管理系统中创建虚拟集群的认证、将虚拟集群的密钥授权给 Runtime Operator、将租户配置库的密钥授权给虚拟集群的 Argo Operator、收集宿主集群中可以作为环境入口的k8s服务的端口号 |

### 删除集群

|            | 物理集群                                                     | 虚拟集群                                                     |
| ---------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 宿主集群   | 从密钥管理系统中删除物理集群的密钥                           | 不支持                                                       |
| 运行时集群 | 从密钥管理系统中删除物理集群的密钥和认证、从密钥管理系统中删除 Runtime Operator 对物理集群密钥的权限 | 从密钥管理系统中删除物理集群的密钥和认证、从密钥管理系统中删除 Runtime Operator 对物理集群密钥的权限 |

## 快速开始

### 准备

安装以下工具，并配置 GOBIN 环境变量：

- [go](https://golang.org/dl/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

准备一个 kubernetes 实例，复制 kubeconfig 文件到 {$HOME}/.kube/config

### 构建

```shell
go mod tidy
go build -o manager main.go
```

### 运行
```shell
./manager
```

### 单元测试

安装 Vault

```shell
wget https://releases.hashicorp.com/vault/1.10.4/vault_1.10.4_linux_amd64.zip
unzip vault_1.10.4_linux_amd64.zip
sudo mv vault /usr/local/bin/
```

安装 Ginkgo

```shell
go install github.com/onsi/ginkgo/v2/ginkgo@v2.10.0
```

执行单元测试

```shell
make test
```

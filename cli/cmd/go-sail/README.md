# go-sail 命令行工具  

简体中文 | [English](./README_EN.md)

#### 安装  
```shell
go install github.com/keepchen/go-sail/cli/cmd/go-sail/v2@latest
```  

#### 使用方法  
- 初始化工程目录
```shell
go-sail init -v ${goVersion} -d ${workDir} -n ${appName} -s ${serviceName}  

# example
go-sail init -v 1.17 -d /var/www/kupo-ha -n pets -s corgi
```  

#### 完成  
工程目录被成功创建，目录结构大致如下:  
```text
pets
├── cmd
│   ├── corgi.go
│   └── root.go
├── go.mod
├── go.sum
├── logs
│   └── running.log
├── main.go
└── pkg
    ├── app
    │   └── corgi
    ├── common
    ├── constants
    ├── lib
    └── utils
```  

- 在已有的工程中添加新的服务 (即子应用)
```shell
go-sail add -v ${goVersion} -d ${workDir} -n ${appName} -s ${serviceName}  

# example
go-sail add -v 1.17 -d /var/www/kupo-ha -n pets -s bourne
```  


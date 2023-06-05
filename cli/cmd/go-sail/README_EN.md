# go-sail cli  

English | [简体中文](./README.md)

#### Installation  
```shell
go install github.com/keepchen/go-sail/cli/cmd/go-sail@latest
```  

#### Usage  
- Init project
```shell
go-sail init -v ${goVersion} -d ${workDir} -n ${appName} -s ${serviceName}  

# example
go-sail init -v 1.17 -d /var/www/kupo-ha -n pets -s corgi
```  

#### Finish  
A golang project will be created,structures like this:  
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

- Add new service (sub application)
```shell
go-sail add -v ${goVersion} -d ${workDir} -n ${appName} -s ${serviceName}  

# example
go-sail add -v 1.17 -d /var/www/kupo-ha -n pets -s bourne
```  


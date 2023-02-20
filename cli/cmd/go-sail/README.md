# go-sail cli  

#### Installation  
```shell
go install github.com/keepchen/go-sail/cli/cmd/go-sail@latest
```  

#### Usage  
```shell
go-sail gen -v ${goVersion} -d ${workDir} -n ${appName} -s ${serviceName}  

# example
gen -v 1.17 -d /var/www/kupo-ha -n pets -s corgi
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


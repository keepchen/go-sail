package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/keepchen/go-sail/v3/lib/nats"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/config"
	"github.com/keepchen/go-sail/v3/examples/pkg/common/db/mock"
	"github.com/keepchen/go-sail/v3/examples/pkg/common/db/models/users"
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/spf13/cobra"
)

func userCMD() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "user",
		Short: "启动user服务",
		Run: func(cmd *cobra.Command, args []string) {
			//启动时要执行的操作写在这里
			wg := &sync.WaitGroup{}

			//启动http接口服务
			wg.Add(1)
			go user.StartServer(wg)

			//更多服务...
			utils.Signal().ListeningExit(wg)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			//启动前要执行的方法写在这里，例如加载配置文件、初始化数据库连接等
			if len(os.Getenv("nacosAddrs")) != 0 && len(os.Getenv("nacosNamespaceID")) != 0 {
				nacos.InitClient(cmd.Use, os.Getenv("nacosAddrs"), os.Getenv("nacosNamespaceID"))
			}
			//::解析配置
			config.ParseConfig(cfgPath)
			//::初始化redis集群连接
			redis.InitRedisCluster(config.GetGlobalConfig().RedisCluster)
			//::初始化nats
			nats.Init(config.GetGlobalConfig().Nats)
			//::初始化日志组件
			logger.InitLoggerZap(config.GetGlobalConfig().Logger, config.GetGlobalConfig().AppName)
			//::初始化数据库
			db.InitDB(config.GetGlobalConfig().Datasource)
			//当数据库配置了主从自动同步的情况下，只对写库进行结构同步
			if config.GetGlobalConfig().Datasource.AutoMigrate {
				users.AutoMigrate(db.GetInstance().W)
				//当开启了自动同步结构且处于调试模式时，进行数据mock操作
				if config.GetGlobalConfig().Debug {
					log.Println("[!] Tips: mock user and wallet data")
					mock.CreateUserAndWalletData()
				}
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "配置文件路径")
	return cmd
}

func init() {
	RootCMD.AddCommand(userCMD())
}

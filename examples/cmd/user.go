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
			//Place startup operations here
			wg := &sync.WaitGroup{}

			//Start the HTTP API service
			wg.Add(1)
			go user.StartServer(wg)

			//More services...
			utils.Signal().ListeningExit(wg)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			//Place pre-startup operations here, such as loading configuration and initializing database connections
			if len(os.Getenv("nacosAddrs")) != 0 && len(os.Getenv("nacosNamespaceID")) != 0 {
				nacos.InitClient(cmd.Use, os.Getenv("nacosAddrs"), os.Getenv("nacosNamespaceID"))
			}
			//::Parse configuration
			config.ParseConfig(cfgPath)
			//::Initialize the Redis cluster connection
			redis.InitRedisCluster(config.GetGlobalConfig().RedisCluster)
			//::Initialize NATS
			nats.Init(config.GetGlobalConfig().Nats)
			//::Initialize the logging component
			logger.InitLoggerZap(config.GetGlobalConfig().Logger, config.GetGlobalConfig().AppName)
			//::Initialize the database
			db.InitDB(config.GetGlobalConfig().Datasource)
			//Migrate only the write database when primary-replica synchronization is configured
			if config.GetGlobalConfig().Datasource.AutoMigrate {
				users.AutoMigrate(db.GetInstance().W)
				//Mock data when automatic schema migration and debug mode are enabled
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

// ----------- api doc definition -----------------

// @title          user - <go-sail>
// @version        1.0
// @description    This is an api document of go-sail.
// @termsOfService https://blog.keepchen.com

// @contact.name  keepchen
// @contact.url   https://blog.keepchen.com
// @contact.email keepchen2016@gmail.com

// @license.name MIT
// @license.url  https://github.com/keepchen/go-sail/LICENSE

// @Host     localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Access Token protects our entity endpoints

// ----------- api doc definition -----------------

package user

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/jwt"

	"go.uber.org/zap"

	"github.com/keepchen/go-sail/v3/schedule"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/routes"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	userConfig "github.com/keepchen/go-sail/v3/examples/pkg/app/user/config"
)

// StartServer 启动服务
//
//go:generate swag init --dir ./ --output ./http/docs --parseDependency --parseInternal --generalInfo user.go
//go:generate redoc-cli bundle ./http/docs/*.yaml -o ./http/docs/docs.html
//go:generate node ../../../../plugins/redocly/redocly-copy.js ./http/docs/docs.html
func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		conf = &config.Config{
			LoggerConf: logger.Conf{
				ConsoleOutput: true,
				Filename:      "examples/logs/running.log",
			},
			HttpServer: config.HttpServerConf{
				Debug: true,
				Addr:  ":8000",
				Swagger: config.SwaggerConf{
					Enable:      true,
					RedocUIPath: "examples/pkg/app/user/http/docs/docs.html",
					JsonPath:    "examples/pkg/app/user/http/docs/swagger.json",
				},
				Prometheus: config.PrometheusConf{
					Enable:         true,
					Addr:           ":19100",
					AccessPath:     "/metrics",
					DiskPath:       "/var",
					SampleInterval: "10s",
				},
				WebSocketRoutePath: "go-sail-ws",
			},
			JwtConf: &jwt.Conf{
				Enable:      true,
				Algorithm:   "RS256",
				TokenIssuer: "go-sail-examples",
				PrivateKey:  "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDUvUDx+LPQ0S+L+5UmtD2EJw1L953mVCMWBJktBbqPTIhDmrd33+3cNq0t7rXuALhoqZS/53nDchU1wsCveieNDR7SsdO4HMS4bnxgyuYCkC1ugAdyvJ2FCv7xUppc7PvyIQ1gQS/nOP0wKcZiFiqxpVBoVKzSv/Tw4ct8p/WL2u75xakZj5oM6ztTdwYwxnRcs5EylWZ1QD7m/y9pwLO79arvbZggQff1GkvbJ3FM/arlsE2st4NZ3HIHFmU/3Nn9PsBb5uiogYN08coGZaspMlD4YbNSo4bOu5hDmwzdOdTwC9vg4Xfq7IAMHADhs/3ji0pI1IRPsELc3RR6tGhDAgMBAAECggEAQal4VjcxKQ6n4kjwrFWNdzCmhgATmHf3rGAW9zKBdqFknZkvb6yKOiIWKcs4FBHc2VEePG0xxAV+Tm2iE4dclciq7tU8R+N5RIO1mBqIC9p8a1LQ+bUF2X6fWdTpGC19Riq1ejQkmPWaEDeUp8m3u8UOoGUiQppE++R1bjBZNaT5S16qbfDOV9plF550wnwbq6fNZlWT4PdiI6ox/4KPZdhIKGKnKkh4xX4mHk3E9fl+udHbXiT3qjSDOEchUpHglzNZG1LMD2BWb+zcxUbJzm2r5BviZmKHPd4+w5mt+kfPbDHnFCgjnlZoFswNFO2s/ZHk99NveoDa1i0OVGbuwQKBgQD045x067oKAvBcRawtP5H6DaifBFnwp3xTg8GIPwitD8+bQMIi6s0jbN7HI7A7S5BTFnukCweURhBiXmbrs29ImiaxcIVGdBKDhkHABq3oino4oHVs8bLpD9moaQccon79aPQ38j8KeU1UJ7a4R6Jbd2eLoZPmj5bPrQHkbqDWWwKBgQDeZDjS0tLtOjbY9CHdB9+BWSyL2DylqSyec/Ew/c8sr2SBK59Db0W2Vgc62iTOYjzWYBTBUrYWRRoAnSoQLkePjQ+mpzGMtpR9BKq3ADrIremgJGRFIo+NL9qjpJu238na+FGp1DgfTSXMMxzLC23lTgh6PXIIcF03kL/yN9lqOQKBgQDPCmCMuX9gV3u/h2g6GTThpAqb5qHjxLZoJUzKVACR0HxFVkrMGpe1C6aN1q54czpiBPAjkO+nfFT91bJONDYxu6JbAjarihbc+/U61GrT37/VgFPG99G7GZt7ttA8dWXH+aQAaN7DjCrEq47f3jB2BE2Wz9SraVqn2i1vY9i3YQKBgQCu0d4RbIU+0upWteMA27WI+s6XyA40s75NeRr6xipcGCxLlj0GR6xnX00jqGQSkQr+Al2OczSMYRnFrcZpHdhHMj5BZWEAGm6zsD16ygVrx7rFlpXz+u0ZsaqPxVBa+6S0K0wW0qqjgIPb97oEqyFihmsHnNHNbHb6vSEGiXyxkQKBgAwb/3lWqp1Zpj6hMw9NdB0c6huQYLqX2INkKj9PcIlFq0nOeHMZfMisuQKhvcGsPQsHMP2NbPjZiLnbpRHPvplU0p7ayaXuNF2t73k/L5f92+8VBuYECEUOXw2xST5gvkPdKGK1xM1cLT6y8TrFRIXvUK2duHjDxiaPKtANi2P4",
				PublicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1L1A8fiz0NEvi/uVJrQ9hCcNS/ed5lQjFgSZLQW6j0yIQ5q3d9/t3DatLe617gC4aKmUv+d5w3IVNcLAr3onjQ0e0rHTuBzEuG58YMrmApAtboAHcrydhQr+8VKaXOz78iENYEEv5zj9MCnGYhYqsaVQaFSs0r/08OHLfKf1i9ru+cWpGY+aDOs7U3cGMMZ0XLORMpVmdUA+5v8vacCzu/Wq722YIEH39RpL2ydxTP2q5bBNrLeDWdxyBxZlP9zZ/T7AW+boqIGDdPHKBmWrKTJQ+GGzUqOGzruYQ5sM3TnU8Avb4OF36uyADBwA4bP944tKSNSET7BC3N0UerRoQwIDAQAB",
			},
		}
		//apiOption = &api.Option{
		//	EmptyDataStruct:  api.DefaultEmptyDataStructObject,
		//	ErrNoneCode:      constants.CodeType(200),
		//	ErrNoneCodeMsg:   "SUCCEED",
		//	ForceHttpCode200: true,
		//}
		beforeFunc = func() {
			fmt.Println("call user function [before] to do something...")
		}
		afterFunc = func() {
			sail.GetLogger().Info("print log info to console")

			token, err := sail.JWT().MakeToken("10000", time.Now().Add(time.Hour).Unix())
			fmt.Println("issue token:", token, err)
			valid, claims, err := sail.JWT().ValidToken(token)
			fmt.Println("valid token:", valid, claims, err)

			fmt.Println("call user function [after] to do something...")
			job0 := "print now datetime"
			cancel0 := schedule.NewJob(job0, func() {
				fmt.Println("now: ", utils.Datetime().FormatDate(time.Now(), utils.YYYY_MM_DD_HH_MM_SS_EN))
			}).RunAt(schedule.EveryMinute)
			time.AfterFunc(time.Minute*3, cancel0)

			job1 := "print hello"
			cancel1 := schedule.NewJob(job1, func() {
				time.Sleep(time.Second * 10)
				fmt.Println(utils.Datetime().FormatDate(time.Now(), utils.YYYY_MM_DD_HH_MM_SS_EN), "hello")
				sail.GetLogger().Info("print log info to console",
					zap.String("value", "go-sail"),
					zap.Errors("errors", []error{nil}))
			}).EverySecond()
			//duplicate job name
			//schedule.NewJob(job1, func() {
			//	fmt.Println("This task will be panic")
			//}).EverySecond()
			time.AfterFunc(time.Second*33, cancel1)

			ticker := time.NewTicker(time.Second)
			times := 0
		LOOP:
			for range ticker.C {
				times++
				fmt.Printf("job: {%s} is running: %t | job: {%s} is running: %t\n",
					job0, schedule.JobIsRunning(job0), job1, schedule.JobIsRunning(job1))
				if times > 50 {
					break LOOP
				}
			}
		}
	)

	//挂载处理方法后启动
	sail.WakeupHttp("go-sail", conf).
		//SetupApiOption(apiOption).
		//EnableWebsocket(nil, nil).
		Hook(routes.RegisterRoutes, beforeFunc, afterFunc).
		Launch()
}

// RegisterServicesToNacos 将服务注册到注册中心
func RegisterServicesToNacos(wg *sync.WaitGroup) {
	defer wg.Done()

	nc := nacos.GetNamingClient()
	var param vo.RegisterInstanceParam
	localIp, err := utils.IP().GetLocal()
	if err == nil {
		param.Ip = localIp
	}
	portSplit := strings.Split(userConfig.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "go-sail-user"
	param.Weight = 100
	param.Enable = true
	param.Healthy = true
	param.Ephemeral = true
	param.Metadata = map[string]string{"description": "go-sail-user服务"}
	_, err = nc.RegisterInstance(param)
	if err != nil {
		panic(err)
	}
}

// UnregisterServiceFromNacos 将服务从注册中心下线
func UnregisterServiceFromNacos() {
	nc := nacos.GetNamingClient()
	var param vo.DeregisterInstanceParam
	localIp, err := utils.IP().GetLocal()
	if err == nil {
		param.Ip = localIp
	}
	portSplit := strings.Split(userConfig.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "go-sail-user"
	param.Ephemeral = true

	_, _ = nc.DeregisterInstance(param)
}

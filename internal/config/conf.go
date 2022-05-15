package config

import (
	"log"
	"time"

	"github.com/daheige/tigago/logger"
	"github.com/daheige/tigago/setting"
	"go.uber.org/zap"
)

var (
	// AppServerConf app config
	AppServerConf = &AppServerSettingS{}
)

// AppServerSettingS server config.
type AppServerSettingS struct {
	AppEnv       string
	AppDebug     bool
	GRPCPort     int
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogDir       string
	LogFileName  string
	PProfPort    int // pprof port
}

// InitConfig 读取rpc配置文件
func InitConfig(configDir string) error {
	s, err := setting.NewSetting(configDir, "app")
	if err != nil {
		return err
	}

	err = s.ReadSection("AppServer", &AppServerConf)
	if err != nil {
		return err
	}

	AppServerConf.ReadTimeout *= time.Second
	AppServerConf.WriteTimeout *= time.Second

	if AppServerConf.AppDebug {
		log.Println("app server config: ", AppServerConf)
	}

	// init logger
	InitLogger()

	// init db
	err = InitDatabase(s, "DbDefault", "default")
	if err != nil {
		return err
	}

	// 初始化redis
	err = InitRedis(s, "RedisCommon", "default")
	if err != nil {
		return err
	}

	return nil
}

// InitLogger 初始化日志句柄
func InitLogger() {
	if AppServerConf.LogDir == "" {
		AppServerConf.LogDir = "./logs"
	}

	if AppServerConf.LogFileName == "" {
		AppServerConf.LogFileName = "go-rpc.log"
	}

	opts := []logger.Option{
		logger.WithLogDir(AppServerConf.LogDir),           // 日志目录
		logger.WithLogFilename(AppServerConf.LogFileName), // 日志文件名，默认zap.log
		logger.WithJsonFormat(true),                       // json格式化
		logger.WithAddCaller(true),                        // 打印行号
		logger.WithLogLevel(zap.DebugLevel),               // 设置日志打印最低级别,如果不设置默认为info级别
		logger.WithMaxAge(7),                              // 最大保存3天
		logger.WithMaxSize(200),                           // 每个日志文件最大20MB
	}

	// 生成默认的日志句柄对象
	logger.Default(opts...)
}

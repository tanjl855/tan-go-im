package conf

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/servers/push_server/internal/db"
)

type (
	Config struct {
		App        `yaml:"app"`
		Data       ` yaml:"data"`
		Server     ` yaml:"server"`
		Log        `yaml:"log"`
		ThirdParty `yaml:"third_party"`
	}
	App struct {
		Auth    `yaml:"auth"`
		Name    string `yaml:"name" env:"APP_NAME"`
		Host    string `yaml:"host" env:"APP_HOST"`
		RunMode string `yaml:"run_mode" env:"RUN_MODE"`
	}
	Data struct {
		//Mongo `yaml:"mongo"`
		Redis `yaml:"redis"`
		Kafka `yaml:"kafka"`
	}
	Server struct {
		Http `yaml:"http"`
		Grpc `yaml:"grpc"`
	}
	Http struct {
		Addr         int `yaml:"port" env:"SERVER_HTTP_PORT"`
		ReadTimeout  int `yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`
		WriteTimeout int `yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
		Ws           `yaml:"ws"`
	}
	Ws struct {
		PrivateIp           string `yaml:"private_ip"`
		WebsocketMaxConnNum int    `yaml:"websocket_max_conn_num"`
		WebsocketTimeOut    int    `yaml:"websocket_time_out"`
		WebsocketMaxMsgLen  int    `yaml:"websocket_max_msg_len"`
	}
	Grpc struct {
		Addr         string `yaml:"addr"`
		RegisterName string `yaml:"register_name"`
	}
	Kafka struct {
		Addr          string   `yaml:"addr"`
		Topics        []string `yaml:"topics"`
		ConsumerGroup string   `yaml:"consumer_group"`
	}
	ThirdParty struct {
		Email `yaml:"email"`
	}
	Mongo struct {
		Url      string `yaml:"url" env:"DATA_MONGO_URL"`
		Database string `yaml:"database" env:"DATA_MONGO_DATABASE"`
	}
	Redis struct {
		ADDR     string `yaml:"addr" env:"DATA_REDIS_ADDR"`
		Password string `yaml:"password" env:"DATA_REDIS_PASSWORD"`
	}
	Log struct {
		LogPath      string `yaml:"log_path"`
		LogLevel     string `yaml:"log_level"`
		LogEncodeMod string `yaml:"log_encode_mod"`
		IsConsole    bool   `yaml:"is_console"`
	}
	Email struct {
		Sender string `yaml:"sender"` //发送人邮箱（邮箱以自己的为准）
		Name   string `yaml:"name"`
		Pass   string `yaml:"pass"` //发送人邮箱的密码，现在可能会需要邮箱 开启授权密码后在pass填写授权码
		Host   string `yaml:"host"` //邮箱服务器（此时用的是qq邮箱）
		Port   int    `yaml:"port"` //邮箱服务器端口
	}
	Auth struct {
		Secret string `yaml:"secret"`
	}
)

var Bootstrap = &Config{}

func init() {
	err := cleanenv.ReadConfig("D:\\goland_project\\tan_go_im\\configs\\push_server_conf.yml", Bootstrap)
	if err != nil {
		panic(err)
	}
	if Bootstrap.App.RunMode != "debug" {
		//初始化日志
		log.InitLog(Bootstrap.Log.LogPath, Bootstrap.Log.LogLevel, Bootstrap.Log.LogEncodeMod, Bootstrap.Log.IsConsole)
	} else {
		//初始化日志
		log.InitLog("", "", "", true)
	}
	db.InitDB(Bootstrap.Data.Redis.ADDR, Bootstrap.Data.Redis.Password, Bootstrap.Kafka.Addr, Bootstrap.Kafka.ConsumerGroup)
}

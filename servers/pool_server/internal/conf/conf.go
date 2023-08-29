package conf

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tanjl855/tan_go_im/pkg/im_auth"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/db"
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
		OutAddr      string `yaml:"out_addr""`
		Addr         int    `yaml:"port" env:"SERVER_HTTP_PORT"`
		ReadTimeout  int    `yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`
		WriteTimeout int    `yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
		Ws           `yaml:"ws"`
	}
	Ws struct {
		WebsocketMaxConnNum int `yaml:"websocket_max_conn_num"`
		WebsocketTimeOut    int `yaml:"websocket_time_out"`
		WebsocketMaxMsgLen  int `yaml:"websocket_max_msg_len"`
	}
	Grpc struct {
		Port         string `yaml:"port"`
		IP           string `yaml:"ip"`
		RegisterName string `yaml:"register_name"`
	}
	Kafka struct {
		Addr   string   `yaml:"addr"`
		Topics []string `yaml:"topics"`
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
	err := cleanenv.ReadConfig("D:\\Goland_project\\tan-go-im\\configs\\pool_server_conf.yml", Bootstrap)
	if err != nil {
		panic(err)
	}
	if Bootstrap.App.RunMode != "debug" {
		log.InitLog(Bootstrap.Log.LogPath, Bootstrap.Log.LogLevel, Bootstrap.Log.LogEncodeMod, Bootstrap.Log.IsConsole)
	} else {
		log.InitLog("", "", "", true)
	}
	db.InitDB(Bootstrap.Data.Redis.ADDR, Bootstrap.Data.Redis.Password, Bootstrap.Kafka.Addr, Bootstrap.Kafka.Topics)
	im_auth.InitSecret(Bootstrap.App.Auth.Secret)
}

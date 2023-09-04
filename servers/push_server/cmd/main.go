package main

import (
	"context"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/servers/push_server/internal/conf"
	"github.com/tanjl855/tan_go_im/servers/push_server/internal/controller"
	"github.com/tanjl855/tan_go_im/servers/push_server/internal/db"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := db.DB.KafkaConsumerGroup.Consume(context.Background(), conf.Bootstrap.Kafka.Topics, controller.NewMsgConsumer()); err != nil {
			log.Panic("消费消息服务关闭，错误：", err)
		}
	}()
	log.Info("消费消息服务开启，消费topic:", conf.Bootstrap.Kafka.Topics)
	<-quit
	log.Warn("关闭服务")

	if err := db.DB.KafkaConsumerGroup.Close(); err != nil {
		log.Panic("强制关闭服务：", err)
	}

	log.Info("服务退出")
}

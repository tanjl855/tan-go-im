package db

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"time"
)

var DB *PoolServerDB

type PoolServerDB struct {
	Rdb                *redis.Client
	MongoDB            *qmgo.Database
	KafkaConsumerGroup sarama.ConsumerGroup
}

func InitDB(redisAddr, redisPW string, kafkaAddr string, kafkaGroup string) {
	DB = &PoolServerDB{}
	DB.Rdb = data.NewRedis(redisAddr, redisPW)
	DB.KafkaConsumerGroup = initKafkaConsumer(kafkaAddr, kafkaGroup)
	//DB.MongoDB = data.NewMongoDB(mongoUrl, mongoDB)
}

// 初始化消费者group
func initKafkaConsumer(kafkaAddr string, ConsumerGroup string) sarama.ConsumerGroup {
	cli := data.NewKafka([]string{kafkaAddr})

	kafkaConsumerGroup, err := sarama.NewConsumerGroupFromClient(ConsumerGroup, cli)
	if err != nil {
		log.Panic(err)
	}
	log.Info("消费group：", ConsumerGroup)
	return kafkaConsumerGroup
}

// SetObjectToRedis 把结构体存入redis
func (db *PoolServerDB) SetObjectToRedis(ctx context.Context, key string, data interface{}, ex time.Duration) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return db.Rdb.Set(ctx, key, jsonStr, ex).Err()
}

// GetObjectFromRedis 从redis获取结构体
func (db *PoolServerDB) GetObjectFromRedis(ctx context.Context, key string, out interface{}) error {
	jsonStr, err := db.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), out)
}

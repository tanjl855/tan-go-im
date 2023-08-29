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
	Rdb          *redis.Client
	MongoDB      *qmgo.Database
	KafkaProduct sarama.SyncProducer
}

func InitDB(redisAddr, redisPW string, kafkaAddr string, kafkaTopics []string) {
	DB = &PoolServerDB{}
	DB.Rdb = data.NewRedis(redisAddr, redisPW)
	DB.KafkaProduct = initKafkaProduct(kafkaAddr, kafkaTopics)
}

func initKafkaProduct(kafkaAddr string, kafkaTopics []string) sarama.SyncProducer {
	cli := data.NewKafka([]string{kafkaAddr})
	admin, err := sarama.NewClusterAdminFromClient(cli)
	if err != nil {
		log.Panic(err)
	}
	// 获取已经存在的topic信息
	topicList, err := admin.DescribeTopics(kafkaTopics)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < len(topicList); i++ {
		if topicList[i].Err != 0 {
			// 创建初始化topic
			err = admin.CreateTopic(kafkaTopics[i], &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 1,
			}, false)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	kafkaProduct, err := sarama.NewSyncProducerFromClient(cli)
	if err != nil {
		log.Panic(err)
	}
	return kafkaProduct
}

// SetObjectToRedis 把结构存入redis
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

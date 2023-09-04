package controller

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/proto/pb_msg"
	"github.com/tanjl855/tan_go_im/servers/push_server/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var _ = sarama.ConsumerGroupHandler(&MsgConsumer{})

type MsgConsumer struct{}

func (m *MsgConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m *MsgConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m *MsgConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("消费消息推送出现panic: %+v", err))
			return
		}
	}()
	for message := range claim.Messages() {
		ctx := context.Background()

		msg := &pb_msg.Msg{}
		err := proto.Unmarshal(message.Value, msg)
		if err != nil {
			log.Error(fmt.Sprintf("消费消息出错，消息：%+v，错误：%+v", message, err))
			session.MarkMessage(message, "")
			continue
		}
		log.Info(fmt.Sprintf("消费 msg:%+v", msg.Content))
		//单人推送
		if msg.SessionType == 0 {
			addr, err := db.DB.Rdb.Get(ctx, "user_conn:"+msg.RecvID).Result()
			if err != nil || addr == "" {
				log.Error(fmt.Sprintf("消费消息出错,用户不在线，消息：%+v，错误：%+v", message, err))
				session.MarkMessage(message, "")
				continue
			}
			conn, err := grpc.Dial(addr, grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024*1000), // 最大消息接收大小为100MB
				grpc.MaxCallSendMsgSize(1024*1024*1000), // 最大消息发送大小为100MB
			), grpc.WithInsecure())
			if err != nil {
				log.Error(fmt.Sprintf("创建grpcClient失败，addr：%s,错误：%+v", addr, err))
				session.MarkMessage(message, "")
				continue
			}
			cli := pb_msg.NewMsgServerClient(conn)
			rsp, err := cli.Push(ctx, msg)
			if err != nil || rsp.ErrCode == -1 || rsp.ErrCode == -2 {
				log.Error(fmt.Sprintf("消费消息出错,push错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
				session.MarkMessage(message, "")
				continue
			}
			//其余情况重试一次
			if err == nil && rsp.ErrCode != 0 {
				rsp, err := cli.Push(ctx, msg)
				if err != nil || rsp.ErrCode != 0 {
					log.Error(fmt.Sprintf("消费消息出错,push错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
					session.MarkMessage(message, "")
					continue
				}
			}
			//只消费一次，保证消费
			session.MarkMessage(message, "")
			continue
		}
		//群组推送
		if msg.SessionType == 1 {
			//直接获取群员列表，每个都当单人推送
			userList, err := db.DB.Rdb.SMembers(ctx, "group_members:"+msg.GroupID).Result()
			if err != nil {
				log.Error(fmt.Sprintf("消费消息出错,获取群组链接信息错误，消息：%+v，错误：%+v", message, err))
				session.MarkMessage(message, "")
				continue
			}
			for i := 0; i < len(userList); i++ {
				if userList[i] == msg.SendID {
					continue
				}
				addr, err := db.DB.Rdb.Get(ctx, "user_conn:"+userList[i]).Result()
				if err != nil {
					if err != redis.Nil {
						log.Error(fmt.Sprintf("获取在线用户链接地址错误：%+v", err))
					}
					continue
				}
				conn, err := grpc.Dial(addr, grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(1024*1024*10), // 最大消息接收大小为1MB
					grpc.MaxCallSendMsgSize(1024*1024*10), // 最大消息发送大小为1MB
				), grpc.WithInsecure())
				if err != nil {
					log.Error(fmt.Sprintf("创建grpcClient失败，addr：%s,错误：%+v", addr, err))
					continue
				}
				cli := pb_msg.NewMsgServerClient(conn)
				//更新接收人
				msg.RecvID = userList[i]
				rsp, err := cli.Push(ctx, msg)
				//返回err不为nil或者消息解析错误或者用户不存在，直接跳出
				if err != nil || rsp.ErrCode == -1 || rsp.ErrCode == -2 {
					////用户的ws链接不在该服务器，用户的在线链接，应该全交给pool服务操作
					//if rsp.ErrCode == -1 {
					//	log.Error(fmt.Sprintf("消费消息出错,用户不在线错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
					//	//设置5s后过期
					//	err = db.DB.Rdb.Set(ctx, "user_conn:"+msg.RecvID, addr, time.Second*5).Err()
					//	if err != nil {
					//		log.Error(fmt.Sprintf("设置用户链接为过期链接失败：%+v", err))
					//	}
					//} else {
					//	log.Error(fmt.Sprintf("消费消息出错,push错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
					//}
					log.Error(fmt.Sprintf("消费消息出错,push错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
					continue
				}
				//其余情况重试一次
				if err == nil && rsp.ErrCode != 0 {
					rsp, err := cli.Push(ctx, msg)
					//返回err不为nil
					if err != nil || rsp.ErrCode != 0 {
						log.Error(fmt.Sprintf("消费消息出错,push错误，消息：%+v，响应：%+v，错误：%+v", message, rsp, err))
						continue
					}
				}
			}
			//保证已经消费消息
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func NewMsgConsumer() sarama.ConsumerGroupHandler {
	return &MsgConsumer{}
}

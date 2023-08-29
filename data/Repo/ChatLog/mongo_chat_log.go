package ChatLog

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IMongoChatLogRepo interface {
	GetChatLogIDByChatId(ctx context.Context, chatId string) (string, error)
	GetChatLogById(ctx context.Context, id string) (*MongoModel.ChatLog, error)
	GetChatLogByChatId(ctx context.Context, chatId string) (*MongoModel.ChatLog, error)
	CreateChatLog(ctx context.Context, chat *MongoModel.ChatLog) error
	UpdateChatLogById(ctx context.Context, id string, msg *MongoModel.ChatLog) error
	UpdateChatLogByChatId(ctx context.Context, ChatId string, msg *MongoModel.ChatLog) error
	DelChatLog(ctx context.Context, id string) error
	AppendMsgByChatId(ctx context.Context, ChatId string, msg *MongoModel.MsgInfo) error
	PullMsgByChatId(ctx context.Context, chatId string) error
}

type chatLogRepo struct {
	data.IMdBaseRepo
}

func (c *chatLogRepo) GetChatLogIDByChatId(ctx context.Context, chatId string) (string, error) {

	temp := &MongoModel.ChatLog{}
	filter := bson.M{"chat_id": chatId}
	err := c.QueryFirst(ctx, temp, filter, bson.M{})
	if err != nil {
		return "", err
	}

	return temp.Id.Hex(), nil
}

func (c *chatLogRepo) GetChatLogById(ctx context.Context, id string) (*MongoModel.ChatLog, error) {
	out := &MongoModel.ChatLog{}
	err := c.FindOne(ctx, id, out)
	return out, err
}

func (c *chatLogRepo) GetChatLogByChatId(ctx context.Context, chatId string) (*MongoModel.ChatLog, error) {
	out := &MongoModel.ChatLog{}
	filter := bson.M{"chat_id": chatId}
	err := c.QueryFirst(ctx, out, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return out, err
}

func (c *chatLogRepo) CreateChatLog(ctx context.Context, chat *MongoModel.ChatLog) error {
	return c.Save(ctx, chat)
}

func (c *chatLogRepo) UpdateChatLogById(ctx context.Context, id string, chat *MongoModel.ChatLog) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	return c.FindOneAndUpdate(ctx, chat, filter)
}

func (c *chatLogRepo) UpdateChatLogByChatId(ctx context.Context, ChatId string, chat *MongoModel.ChatLog) error {
	filter := bson.M{"chat_id": ChatId}
	return c.FindOneAndUpdate(ctx, chat, filter)
}

func (c *chatLogRepo) DelChatLog(ctx context.Context, id string) error {
	return c.Delete(ctx, id)
}

func (c *chatLogRepo) AppendMsgByChatId(ctx context.Context, ChatId string, msg *MongoModel.MsgInfo) error {
	filter := bson.M{"chat_id": ChatId}
	update := bson.M{"$push": bson.M{"chat_msg": msg}}
	return c.FindOneAndUpdate(ctx, filter, update)
}

func (c *chatLogRepo) PullMsgByChatId(ctx context.Context, chatId string) error {
	// 设置过滤条件
	filter := bson.M{"chat_id": chatId}
	// 设置更新操作，使用$pull操作符删除7天前数据
	update := bson.M{
		"$pull": bson.M{
			"chat_msg": bson.M{
				"send_time": bson.M{
					"$lt": time.Now().AddDate(0, 0, -7).Unix(),
				},
			},
		},
	}
	return c.FindOneAndUpdate(ctx, filter, update)
}

func NewIChatLogRepo(database *qmgo.Database) IMongoChatLogRepo {
	Msg := MongoModel.ChatLog{}
	IRepo := &chatLogRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: Msg.TableName(),
	}}

	return IRepo
}

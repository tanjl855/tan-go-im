package SaveLog

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaveLogRepo interface {
	CreatSaveLog(ctx context.Context, save *MongoModel.SaveLog) error
	DelSaveLog(ctx context.Context, id string) error
	GetSaveLogList(ctx context.Context, chatId string, startTime, endTime int64, Limit int64) ([]*MongoModel.SaveLog, error)
}

type saveLogRepo struct {
	data.IMdBaseRepo
}

func (s saveLogRepo) CreatSaveLog(ctx context.Context, save *MongoModel.SaveLog) error {
	return s.Save(ctx, save)
}

func (s saveLogRepo) DelSaveLog(ctx context.Context, id string) error {
	return s.Delete(ctx, id)
}

func (s saveLogRepo) GetSaveLogList(ctx context.Context, chatId string, startTime, endTime int64, Limit int64) ([]*MongoModel.SaveLog, error) {
	out := []*MongoModel.SaveLog{}
	filterList := bson.D{
		bson.E{Key: "chat_id", Value: bson.M{operator.Eq: chatId}},
		bson.E{Key: "send_time", Value: bson.M{operator.Gte: startTime, operator.Lte: endTime}},
	}
	//按发送时间排序返回
	err := s.QueryManyAndOrder(ctx, out, filterList, bson.M{}, 0, Limit, "send_time")
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewISaveLogRepo(database *qmgo.Database) ISaveLogRepo {
	save := MongoModel.SaveLog{}
	IRepo := &saveLogRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: save.TableName(),
	}}
	return IRepo
}

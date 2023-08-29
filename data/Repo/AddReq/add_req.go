package AddReq

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type IMongoAddReqRepo interface {
	Create(ctx context.Context, addReq *MongoModel.AddReq) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, addReq *MongoModel.AddReq) error
	GetListByRecvID(ctx context.Context, recvId string) ([]*MongoModel.AddReq, error)
	GetListBySendID(ctx context.Context, sendId string) ([]*MongoModel.AddReq, error)
	GetListByGroupID(ctx context.Context, groupId string) ([]*MongoModel.AddReq, error)
	GetListByGroupIDAndSendID(ctx context.Context, groupId, sendId string) ([]*MongoModel.AddReq, error)
	GetListByRecvIDAndSendID(ctx context.Context, recvId string, sendId string) ([]*MongoModel.AddReq, error)
	GetInfoByID(ctx context.Context, id string) (*MongoModel.AddReq, error)
}

type addReqRepo struct {
	data.IMdBaseRepo
}

func (a addReqRepo) GetInfoByID(ctx context.Context, id string) (*MongoModel.AddReq, error) {
	out := &MongoModel.AddReq{}
	err := a.FindOne(ctx, id, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (a addReqRepo) Create(ctx context.Context, addReq *MongoModel.AddReq) error {
	return a.Save(ctx, addReq)
}

func (a addReqRepo) Update(ctx context.Context, id string, addReq *MongoModel.AddReq) error {
	return a.UpdateById(ctx, id, addReq)
}

func (a addReqRepo) GetListByRecvID(ctx context.Context, recvId string) ([]*MongoModel.AddReq, error) {
	results := []*MongoModel.AddReq{}
	filter := bson.M{"recv_id": recvId}
	err := a.QueryMany(ctx, &results, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (a addReqRepo) GetListBySendID(ctx context.Context, sendId string) ([]*MongoModel.AddReq, error) {
	results := []*MongoModel.AddReq{}
	filter := bson.M{
		"send_id": sendId,
	}
	err := a.QueryMany(ctx, results, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (a addReqRepo) GetListByGroupID(ctx context.Context, groupId string) ([]*MongoModel.AddReq, error) {
	results := []*MongoModel.AddReq{}
	filter := bson.M{
		"group_id": groupId,
	}
	err := a.QueryMany(ctx, &results, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (a addReqRepo) GetListByGroupIDAndSendID(ctx context.Context, groupId string, sendId string) ([]*MongoModel.AddReq, error) {
	results := []*MongoModel.AddReq{}
	filter := bson.M{
		"group_id": groupId,
		"send_id":  sendId,
	}
	err := a.QueryMany(ctx, results, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (a addReqRepo) GetListByRecvIDAndSendID(ctx context.Context, recvId string, sendId string) ([]*MongoModel.AddReq, error) {
	results := []*MongoModel.AddReq{}
	filter := bson.M{
		"recv_id": recvId,
		"send_id": sendId,
	}
	err := a.QueryMany(ctx, results, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func NewIMongoAddReqRepo(database *qmgo.Database) IMongoAddReqRepo {
	addReq := MongoModel.AddReq{}
	Irepo := &addReqRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: addReq.TableName(),
	}}
	return Irepo
}

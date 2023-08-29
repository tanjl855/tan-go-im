package AddReq

import (
	"context"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
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

package Group

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type IMongoUserToGroupRepo interface {
	CreateUserToGroup(ctx context.Context, info *MongoModel.UserToGroup) error
	UpdateUserToGroup(ctx context.Context, info *MongoModel.UserToGroup) error
	DelUserToGroup(ctx context.Context, id string) error
	DelUserToGroupByUIDAndGroupID(ctx context.Context, uid string, groupId string) error
	GetUserToGroupListByUID(ctx context.Context, uid string) ([]*MongoModel.UserToGroup, error)
	GetUserToGroupInfoByObjectId(ctx context.Context, id string) (*MongoModel.UserToGroup, error)
	GetUserToGroupInfoByGroupIdAndUID(ctx context.Context, groupId string, uid string) (*MongoModel.UserToGroup, error)
}

type userToGroupRepo struct {
	data.IMdBaseRepo
}

func (u userToGroupRepo) CreateUserToGroup(ctx context.Context, info *MongoModel.UserToGroup) error {
	return u.Save(ctx, info)
}

func (u userToGroupRepo) UpdateUserToGroup(ctx context.Context, info *MongoModel.UserToGroup) error {
	return u.UpdateById(ctx, info.Id.Hex(), info)
}

func (u userToGroupRepo) DelUserToGroup(ctx context.Context, id string) error {
	return u.Delete(ctx, id)
}

func (u userToGroupRepo) DelUserToGroupByUIDAndGroupID(ctx context.Context, uid string, groupId string) error {
	out := &MongoModel.UserToGroup{}
	filter := bson.M{
		"uid":      uid,
		"group_id": groupId,
	}
	err := u.QueryFirst(ctx, out, filter, bson.M{})
	if err != nil {
		return err
	}
	if !out.Id.IsZero() {
		err = u.Delete(ctx, out.Id.Hex())
		if err != nil {
			return err
		}
	}
	return nil
}

func (u userToGroupRepo) GetUserToGroupListByUID(ctx context.Context, uid string) ([]*MongoModel.UserToGroup, error) {
	out := []*MongoModel.UserToGroup{}
	filter := bson.M{"uid": uid}
	err := u.QueryMany(ctx, &out, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u userToGroupRepo) GetUserToGroupInfoByObjectId(ctx context.Context, id string) (*MongoModel.UserToGroup, error) {
	out := &MongoModel.UserToGroup{}
	err := u.FindOne(ctx, id, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u userToGroupRepo) GetUserToGroupInfoByGroupIdAndUID(ctx context.Context, groupId string, uid string) (*MongoModel.UserToGroup, error) {
	out := &MongoModel.UserToGroup{}
	filter := bson.M{"group_id": groupId, "uid": uid}
	err := u.QueryFirst(ctx, out, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewIUserToGroupRepo(database *qmgo.Database) IMongoUserToGroupRepo {
	temp := MongoModel.UserToGroup{}
	IRepo := &userToGroupRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: temp.TableName(),
	}}

	return IRepo
}

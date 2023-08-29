package Group

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type IMongoGroupRepo interface {
	CreateGroup(ctx context.Context, GroupMembers *MongoModel.Group) error
	UpdateGroup(ctx context.Context, GroupMembers *MongoModel.Group) error
	DelGroup(ctx context.Context, id string) error
	GetGroupInfoByGroupId(ctx context.Context, groupId string) (*MongoModel.Group, error)
	GetGroupInfoByObjectId(ctx context.Context, id string) (*MongoModel.Group, error)
	GetGroupListByGroupName(ctx context.Context, groupName string, page int64, pagesize int64) ([]*MongoModel.Group, int, error)
}

type groupRepo struct {
	data.IMdBaseRepo
}

func (u *groupRepo) CreateGroup(ctx context.Context, GroupMembers *MongoModel.Group) error {
	return u.Save(ctx, GroupMembers)
}

func (u *groupRepo) UpdateGroup(ctx context.Context, GroupMembers *MongoModel.Group) error {
	return u.UpdateById(ctx, GroupMembers.Id.Hex(), GroupMembers)
}

func (u *groupRepo) DelGroup(ctx context.Context, id string) error {
	return u.Delete(ctx, id)
}

func (u *groupRepo) GetGroupInfoByGroupId(ctx context.Context, groupId string) (*MongoModel.Group, error) {
	GroupMembers := &MongoModel.Group{}
	filter := bson.M{"group_id": groupId}
	err := u.QueryFirst(ctx, GroupMembers, filter, bson.M{})
	return GroupMembers, err
}

func (u *groupRepo) GetGroupInfoByObjectId(ctx context.Context, id string) (*MongoModel.Group, error) {
	GroupMembers := &MongoModel.Group{}
	err := u.FindOne(ctx, id, GroupMembers)
	return GroupMembers, err
}

// GetGroupListByGroupName 群名搜索
func (u *groupRepo) GetGroupListByGroupName(ctx context.Context, groupName string, page int64, pagesize int64) ([]*MongoModel.Group, int, error) {
	GroupMembers := []*MongoModel.Group{}
	//filter := bson.E{"group_name", bson.M{operator.Regex: "^" + groupName}}
	filter := bson.M{"group_name": groupName}
	err := u.QueryManyAndOrder(ctx, &GroupMembers, filter, bson.M{}, page*pagesize, pagesize, "-updated_at")
	if err != nil {
		return nil, 0, err
	}
	sum, err := u.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return GroupMembers, int(sum), err
}

func NewIGroupRepo(database *qmgo.Database) IMongoGroupRepo {
	GroupMembers := MongoModel.Group{}
	IRepo := &groupRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: GroupMembers.TableName(),
	}}

	return IRepo
}

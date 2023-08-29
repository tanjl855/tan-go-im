package Friend

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type IMongoFriendsRepo interface {
	CreateFriend(ctx context.Context, firstUID string, secondUID string) error
	UpdateFriend(ctx context.Context, friend *MongoModel.Friend) error
	DelFriend(ctx context.Context, id string) error
	GetFriendListByUID(ctx context.Context, uid string) ([]*MongoModel.Friend, error)
	GetFriend(ctx context.Context, firstUID string, secondUID string) (*MongoModel.Friend, error)
}

type FriendsRepo struct {
	data.IMdBaseRepo
}

func (f FriendsRepo) GetFriend(ctx context.Context, firstUID string, secondUID string) (*MongoModel.Friend, error) {
	if firstUID > secondUID {
		firstUID, secondUID = secondUID, firstUID
	}
	out := &MongoModel.Friend{}
	filter := bson.M{"first_uid": firstUID, "second_uid": secondUID}
	err := f.QueryFirst(ctx, out, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (f FriendsRepo) CreateFriend(ctx context.Context, firstUID string, secondUID string) error {
	if firstUID > secondUID {
		firstUID, secondUID = secondUID, firstUID
	}
	return f.Save(ctx, &MongoModel.Friend{
		FirstUID:  firstUID,
		SecondUID: secondUID,
	})
}

func (f FriendsRepo) UpdateFriend(ctx context.Context, friend *MongoModel.Friend) error {
	return f.UpdateById(ctx, friend.Id.Hex(), friend)
}

func (f FriendsRepo) DelFriend(ctx context.Context, id string) error {
	return f.Delete(ctx, id)
}

func (f FriendsRepo) GetFriendListByUID(ctx context.Context, uid string) ([]*MongoModel.Friend, error) {
	out := []*MongoModel.Friend{}
	filter := bson.M{operator.Or: []bson.M{bson.M{"first_uid": uid}, bson.M{"second_uid": uid}}}
	err := f.QueryMany(ctx, &out, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func NewIFriendRepo(database *qmgo.Database) IMongoFriendsRepo {
	Friends := MongoModel.Friend{}
	IRepo := &FriendsRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: Friends.TableName(),
	}}

	return IRepo
}

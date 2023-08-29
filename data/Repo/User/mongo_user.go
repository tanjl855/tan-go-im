package User

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/tanjl855/tan_go_im/data"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"go.mongodb.org/mongo-driver/bson"
)

type IMongoUserRepo interface {
	CreateUser(ctx context.Context, user *MongoModel.User) error
	UpdateUser(ctx context.Context, id string, user *MongoModel.User) error
	DelUser(ctx context.Context, id string) error
	GetUserInfoByUID(ctx context.Context, uid string) (*MongoModel.User, error)
	GetUserInfoByObjectId(ctx context.Context, id string) (*MongoModel.User, error)
	GetUserInfoByEmail(ctx context.Context, email string) (*MongoModel.User, error)
	GetUserListByNickName(ctx context.Context, nickName string, pageSize int64, page int64) ([]*MongoModel.User, error)
}

type userRepo struct {
	data.IMdBaseRepo
}

func (u userRepo) GetUserListByNickName(ctx context.Context, nickName string, pageSize int64, page int64) ([]*MongoModel.User, error) {
	userList := []*MongoModel.User{}
	filter := bson.M{"nick_name": nickName}
	err := u.QueryManyAndOrder(ctx, &userList, filter, bson.M{}, pageSize*page, pageSize, "-created_at")
	return userList, err
}

func (u userRepo) GetUserInfoByEmail(ctx context.Context, email string) (*MongoModel.User, error) {
	user := &MongoModel.User{}
	filter := bson.M{"email": email}
	err := u.QueryFirst(ctx, user, filter, bson.M{})
	return user, err
}

func (u userRepo) CreateUser(ctx context.Context, user *MongoModel.User) error {
	return u.Save(ctx, user)
}

func (u userRepo) UpdateUser(ctx context.Context, id string, user *MongoModel.User) error {
	return u.UpdateById(ctx, id, user)
}

func (u userRepo) DelUser(ctx context.Context, id string) error {
	return u.Delete(ctx, id)
}

func (u userRepo) GetUserInfoByUID(ctx context.Context, uid string) (*MongoModel.User, error) {
	user := &MongoModel.User{}
	filter := bson.M{"uid": uid}
	err := u.QueryFirst(ctx, user, filter, bson.M{})
	return user, err
}

func (u userRepo) GetUserInfoByObjectId(ctx context.Context, id string) (*MongoModel.User, error) {
	user := &MongoModel.User{}
	err := u.FindOne(ctx, id, user)
	return user, err
}

func NewIUserRepo(database *qmgo.Database) IMongoUserRepo {
	user := MongoModel.User{}
	IRepo := &userRepo{&data.MdBaseRepo{
		Db:         database,
		Collection: user.TableName(),
	}}

	return IRepo
}

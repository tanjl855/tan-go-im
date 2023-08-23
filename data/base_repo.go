package data

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ IMdBaseRepo = (*MdBaseRepo)(nil)

type TransactionFunc func(ctx context.Context) (interface{}, error)
type IMdBaseRepo interface {
	PrintCollection() string
	GetCollection() string
	FindOne(ctx context.Context, id string, entity interface{}) (err error)
	FindAll(ctx context.Context, entity interface{}) (err error)
	Save(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id string) (err error)
	Update(ctx context.Context, filter interface{}, upData interface{}) error
	UpdateById(ctx context.Context, id string, upData interface{}) error
	FindOneAndUpdate(ctx context.Context, filter interface{}, out interface{}) error
	Query(ctx context.Context, out interface{}, data interface{}, args interface{}) (err error)
	QueryFirst(ctx context.Context, out interface{}, data interface{}, args interface{}, sort ...string) (err error)
	QueryManyAndOrder(ctx context.Context, out interface{}, data interface{}, args interface{}, skip int64, size int64, sort ...string) (err error)
	QueryMany(ctx context.Context, out interface{}, data interface{}, args interface{}) (err error)
	QueryManyAndSort(ctx context.Context, out interface{}, data interface{}, args interface{}, sort ...string) (err error)
	FilterAndPage(ctx context.Context, out interface{}, data interface{}, args interface{}, page int64, size int64) (err error)
	Count(ctx context.Context, data interface{}) (c int64, err error)
	structToBson(dat interface{}) (b bson.M, err error)
	Aggregate(ctx context.Context, out, pipeline interface{}, opts ...options.AggregateOptions) error
	AggregateOne(ctx context.Context, out, pipeline interface{}, opts ...options.AggregateOptions) error

	One(ctx context.Context, out, filter interface{}, opts ...QueryOption) error
	List(ctx context.Context, out, filter interface{}, opts ...QueryOption) (err error)
	GetCount(ctx context.Context, filter interface{}, opts ...QueryOption) (c int64, err error)
	Remove(ctx context.Context, filter interface{}) (err error)

	WithSort(fields ...string) QueryOption
	WithSelect(selector interface{}) QueryOption
	WithSkip(n int64) QueryOption
	WithLimit(n int64) QueryOption
	Paginate(page, size int64) QueryOption
}

type MdBaseRepo struct {
	Db         *qmgo.Database
	Collection string
}

func (m *MdBaseRepo) Remove(ctx context.Context, filter interface{}) (err error) {
	return m.Db.Collection(m.Collection).Remove(ctx, filter)
}

func (m *MdBaseRepo) AggregateOne(ctx context.Context, out, pipeline interface{}, opts ...options.AggregateOptions) error {
	err := m.Db.Collection(m.Collection).Aggregate(ctx, pipeline, opts...).One(out)
	return err
}

func (m *MdBaseRepo) Aggregate(ctx context.Context, out, pipeline interface{}, opts ...options.AggregateOptions) error {
	err := m.Db.Collection(m.Collection).Aggregate(ctx, pipeline, opts...).All(out)
	return err
}

type QueryOption func(i qmgo.QueryI) qmgo.QueryI

func (m *MdBaseRepo) GetCount(ctx context.Context, filter interface{}, opts ...QueryOption) (int64, error) {
	i := m.Db.Collection(m.Collection).Find(ctx, filter)
	i = m.optionDB(i, opts...)
	count, err := i.Count()
	if err != nil {
		return 0, err
	}
	return count, err
}

func (m *MdBaseRepo) WithSort(fields ...string) QueryOption {
	return func(i qmgo.QueryI) qmgo.QueryI {
		return i.Sort(fields...)
	}
}

func (m *MdBaseRepo) WithSelect(selector interface{}) QueryOption {
	return func(i qmgo.QueryI) qmgo.QueryI {
		return i.Select(selector)
	}
}
func (m *MdBaseRepo) Paginate(page, size int64) QueryOption {
	return func(i qmgo.QueryI) qmgo.QueryI {
		if page == 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}
		offset := (page - 1) * size
		return i.Skip(offset).Limit(size)
	}
}
func (m *MdBaseRepo) WithSkip(n int64) QueryOption {

	return func(i qmgo.QueryI) qmgo.QueryI {
		return i.Skip(n)
	}
}

func (m *MdBaseRepo) WithLimit(n int64) QueryOption {
	return func(i qmgo.QueryI) qmgo.QueryI {
		return i.Limit(n)
	}
}
func (m *MdBaseRepo) optionDB(db qmgo.QueryI, opts ...QueryOption) qmgo.QueryI {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
func (m *MdBaseRepo) One(ctx context.Context, out, filter interface{}, opts ...QueryOption) (err error) {
	i := m.Db.Collection(m.Collection).Find(ctx, filter)
	i = m.optionDB(i, opts...)
	err = i.One(out)
	return err
}
func (m *MdBaseRepo) List(ctx context.Context, out, filter interface{}, opts ...QueryOption) (err error) {
	i := m.Db.Collection(m.Collection).Find(ctx, filter)
	i = m.optionDB(i, opts...)
	err = i.All(out)
	return err
}
func (m *MdBaseRepo) Save(ctx context.Context, entity interface{}) (err error) {
	_, err = m.Db.Collection(m.Collection).InsertOne(ctx, entity)
	return err
}
func (m *MdBaseRepo) QueryMany(ctx context.Context, out interface{}, data interface{}, args interface{}) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).All(out)
	return err
}
func (m *MdBaseRepo) QueryManyAndSort(ctx context.Context, out interface{}, data interface{}, args interface{}, sort ...string) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).Sort(sort...).All(out)
	return err
}

func (m *MdBaseRepo) Count(ctx context.Context, d interface{}) (c int64, err error) {
	n, err := m.Db.Collection(m.Collection).Find(ctx, d).Count()

	return n, err
}

func (m *MdBaseRepo) QueryManyAndOrder(ctx context.Context, out interface{}, data interface{}, args interface{}, skip int64, size int64, sort ...string) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).Sort(sort...).Skip(skip).Limit(size).All(out)
	return
}

func (m *MdBaseRepo) FilterAndPage(ctx context.Context, out interface{}, data interface{}, args interface{}, page int64, size int64) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).Skip(page).Limit(size).All(out)
	return err
}

func (m *MdBaseRepo) Query(ctx context.Context, out interface{}, data interface{}, args interface{}) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).One(out)
	if err != nil && err == qmgo.ErrNoSuchDocuments {
		return nil
	}
	return err
}

func (m *MdBaseRepo) QueryFirst(ctx context.Context, out interface{}, data interface{}, args interface{}, sort ...string) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, data).Select(args).Sort(sort...).One(out)
	if err != nil && err == qmgo.ErrNoSuchDocuments {
		return nil
	}
	return err
}

func (m *MdBaseRepo) UpdateById(ctx context.Context, id string, upData interface{}) error {
	un := MongoModel.MongoDBModel{}
	ud, err := m.structToBson(upData)
	if err != nil {
		return err
	}
	data := bson.M{"$set": ud}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	data["$set"].(primitive.M)["updatedAt"] = time.Now().Local()
	err = m.Db.Collection(m.Collection).UpdateOne(ctx, filter, data, options.UpdateOptions{UpdateHook: un})
	return err
}

func (m *MdBaseRepo) Update(ctx context.Context, filter interface{}, upData interface{}) error {
	un := MongoModel.MongoDBModel{}
	ud, err := m.structToBson(upData)
	if err != nil {
		return err
	}
	data := bson.M{"$set": ud}
	data["$set"].(primitive.M)["updatedAt"] = time.Now().Local()
	err = m.Db.Collection(m.Collection).UpdateOne(ctx, filter, data, options.UpdateOptions{UpdateHook: un})
	return err
}

func (m *MdBaseRepo) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) error {
	un := MongoModel.MongoDBModel{}
	err := m.Db.Collection(m.Collection).UpdateOne(ctx, filter, update, options.UpdateOptions{
		UpdateHook:    un,
		UpdateOptions: options2.Update().SetUpsert(true),
	})
	if err != nil && err == qmgo.ErrNoSuchDocuments {
		return nil
	}
	return err
}

//还是不用了，主键还不清楚会不会冲突
//func (m *MdBaseRepo) UpdateORCreate(ctx context.Context, id string, upData interface{}) error {
//	un := model.MongoDBModel{}
//	ud, err := m.structToBson(upData)
//	if err != nil {
//		return err
//	}
//	data := bson.M{"$set": ud}
//	objID, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return err
//	}
//	filter := bson.M{"_id": objID}
//	data["$set"].(primitive.M)["updatedAt"] = time.Now().Local()
//	//设置为更新，若无记录则创建
//	updateOrCreate := true
//	err = m.Db.Collection(m.Collection).UpdateOne(ctx, filter, data, options.UpdateOptions{
//		UpdateHook: un,
//		UpdateOptions: &options2.UpdateOptions{
//			ArrayFilters:             nil,
//			BypassDocumentValidation: nil,
//			Collation:                nil,
//			Hint:                     nil,
//			Upsert:                   &updateOrCreate,
//		},
//	})
//	return err
//}

func (m *MdBaseRepo) GetCollection() string {
	return m.Collection
}

func (m *MdBaseRepo) Delete(ctx context.Context, id string) (err error) {

	//pid, _ := primitive.ObjectIDFromHex(id)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	err = m.Db.Collection(m.Collection).Remove(ctx, filter)
	return err
}

func (m *MdBaseRepo) FindAll(ctx context.Context, entity interface{}) (err error) {
	err = m.Db.Collection(m.Collection).Find(ctx, bson.D{}).All(entity)
	return
}

func (m *MdBaseRepo) FindOne(ctx context.Context, id string, entity interface{}) (err error) {

	//pid, _ := primitive.ObjectIDFromHex(id)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	err = m.Db.Collection(m.Collection).Find(ctx, filter).One(entity)

	return
}
func (m *MdBaseRepo) FindOneByUID(ctx context.Context, id string, entity interface{}) (err error) {

	//pid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"UID": id}

	err = m.Db.Collection(m.Collection).Find(ctx, filter).One(entity)

	return
}

func (m *MdBaseRepo) FindOneByUserID(ctx context.Context, id string, entity interface{}) (err error) {

	//pid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"UserID": id}

	err = m.Db.Collection(m.Collection).Find(ctx, filter).One(entity)

	return
}

func (m *MdBaseRepo) PrintCollection() string {
	log.Info(m.Collection)
	return m.Collection
}

func (m *MdBaseRepo) structToBson(dat interface{}) (b bson.M, err error) {
	pByte, err := bson.Marshal(&dat)
	if err != nil {
		return nil, err
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		return nil, err
	}
	return update, nil
}

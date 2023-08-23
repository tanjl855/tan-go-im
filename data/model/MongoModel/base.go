package MongoModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

var nilTime time.Time

// mongodb基础表
type MongoDBModel struct {
	//mongo 原生ID
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"object_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m *MongoDBModel) AfterInsert() error {
	m.CreatedAt = time.Now().Local()
	m.UpdatedAt = time.Now().Local()
	return nil
}

func (m *MongoDBModel) AfterUpdate() error {
	m.UpdatedAt = time.Now().Local()
	return nil
}

// DefaultUpdateAt changes the default updateAt field
func (m *MongoDBModel) DefaultUpdateAt() {
	m.UpdatedAt = time.Now().Local()
}

// DefaultCreateAt changes the default createAt field
func (m *MongoDBModel) DefaultCreateAt() {
	if reflect.DeepEqual(m.CreatedAt, nilTime) {
		m.CreatedAt = time.Now().Local()
	}
}

// DefaultId  changes the default _id field
func (m *MongoDBModel) DefaultId() {
	if m.Id.IsZero() {
		m.Id = primitive.NewObjectID()
	}
}

func (m *MongoDBModel) IsValid() bool {
	return !m.Id.IsZero()
}

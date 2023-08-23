package entity

import (
	"github.com/tanjl855/tan_go_im/data/model/MongoModel"
	"time"
)

type User struct {
	Id          string `json:"id"`
	UID         string `json:"uid"`
	NickName    string `bson:"nick_name" json:"nick_name"`
	FaceURL     string `bson:"face_url" json:"face_url"`
	Gender      int32  `bson:"gender" json:"gender"`
	Email       string `bson:"email" json:"email"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
	Birth       int64  `bson:"birth" json:"birth"`
	Status      int32  `bson:"status" json:"status"`
	//Password    string    `bson:"password" json:"password"`
	Version int `bson:"version" json:"version"`
}

func TransformToModel(user *User) *MongoModel.User {
	return &MongoModel.User{
		UID:         user.UID,
		NickName:    user.NickName,
		FaceURL:     user.FaceURL,
		Gender:      user.Gender,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Birth:       time.Unix(user.Birth, 0),
		Status:      user.Status,
		//Password:    user.Password,
		Version: user.Version,
	}
}

func TransformFromModel(user *MongoModel.User) *User {
	return &User{
		Id:          user.Id.Hex(),
		UID:         user.UID,
		NickName:    user.NickName,
		FaceURL:     user.FaceURL,
		Gender:      user.Gender,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Birth:       user.Birth.Unix(),
		Status:      user.Status,
		//Password:    user.Password,
		Version: user.Version,
	}
}

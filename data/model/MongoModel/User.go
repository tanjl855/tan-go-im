package MongoModel

import "time"

type User struct {
	MongoDBModel `bson:",inline" json:",inline"`
	UID          string    `bson:"uid" json:"uid"`
	NickName     string    `bson:"nick_name" json:"nick_name"`
	FaceURL      string    `bson:"face_url" json:"face_url"`
	Gender       int32     `bson:"gender" json:"gender"`
	Email        string    `bson:"email" json:"email"`
	PhoneNumber  string    `bson:"phone_number" json:"phone_number"`
	Birth        time.Time `bson:"birth" json:"birth"`
	Status       int32     `bson:"status" json:"status"`
	Password     string    `bson:"password" json:"password"`
	Version      int       `bson:"version" json:"version"`
}

func (u *User) TableName() string {
	return "User"
}

package MongoModel

type Friend struct {
	MongoDBModel `bson:",inline" json:",inline"`
	FirstUID     string `bson:"first_uid" json:"first_uid"`
	SecondUID    string `bson:"second_uid" json:"second_uid"`
}

func (f *Friend) TableName() string {
	return "friend"
}

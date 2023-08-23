package MongoModel

type AddReq struct {
	MongoDBModel `bson:",inline" json:",inline"`
	SendID       string `bson:"send_id" json:"send_id""`
	SendNickName string `bson:"send_nick_name" json:"send_nick_name"`
	RecvID       string `bson:"recv_id" json:"recv_id"`
	GroupID      string `bson:"group_id" json:"group_id"` //groupid为空则申请类型为用户-用户，否则为用户-群
	GroupName    string `bson:"group_name" json:"group_name"`
}

func (a *AddReq) TableName() string {
	return "add_req"
}

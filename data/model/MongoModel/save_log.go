package MongoModel

// SaveLog 持久化的库
type SaveLog struct {
	MongoDBModel `bson:",inline" json:",inline"`
	//群聊为“群聊ID”，用户-用户为“用户ID_用户ID”【小的在前面】
	ChatId string `bson:"chat_id" json:"chat_id"`
	//聊天类型，群聊为1，用户-用户为0
	ChatType int `bson:"chat_type" json:"chat_type"`
	//发送时间(查询返回时可能需要)
	SendTime int64 `bson:"send_time" json:"send_time"`
	//每条消息都为一条记录
	Msg *[]byte `bson:"chat_msg" json:"chat_msg"`
}

func (s *SaveLog) TableName() string {
	return "save_log"
}

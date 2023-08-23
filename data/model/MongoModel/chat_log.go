package MongoModel

type MsgInfo struct {
	SendTime int64   `bson:"send_time" json:"send_time"`
	Msg      *[]byte `bson:"msg" json:"msg"`
}

// ChatLog 群或者用户-用户聊天记录(不为永久，为临时，相当于几天聊天的缓存用)
type ChatLog struct {
	MongoDBModel `bson:",inline" json:",inline"`
	//群聊为“群聊ID”，用户-用户为“用户ID_用户ID”【小的在前面】
	ChatId string `bson:"chat_id" json:"chat_id"`
	//聊天类型，群聊为1，用户-用户为0
	ChatType int        `bson:"chat_type" json:"chat_type"`
	Msg      []*MsgInfo `bson:"chat_msg" json:"chat_msg"`
}

func (c *ChatLog) TableName() string {
	return "chat_log"
}

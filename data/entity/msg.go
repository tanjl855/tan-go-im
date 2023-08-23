package entity

type Msg struct {
	SendID           string `protobuf:"bytes,1,opt,name=sendID,proto3" json:"sendID,omitempty"` //发送者id
	RecvID           string `protobuf:"bytes,2,opt,name=recvID,proto3" json:"recvID,omitempty"`
	GroupID          string `protobuf:"bytes,3,opt,name=groupID,proto3" json:"groupID,omitempty"`
	SenderPlatformID int32  `protobuf:"varint,4,opt,name=senderPlatformID,proto3" json:"senderPlatformID,omitempty"`
	SenderNickname   string `protobuf:"bytes,5,opt,name=senderNickname,proto3" json:"senderNickname,omitempty"`
	SenderFaceURL    string `protobuf:"bytes,6,opt,name=senderFaceURL,proto3" json:"senderFaceURL,omitempty"`
	SessionType      int32  `protobuf:"varint,7,opt,name=sessionType,proto3" json:"sessionType,omitempty"`
	MsgFrom          int32  `protobuf:"varint,8,opt,name=msgFrom,proto3" json:"msgFrom,omitempty"`
	ContentType      int32  `protobuf:"varint,9,opt,name=contentType,proto3" json:"contentType,omitempty"`
	Content          []byte `protobuf:"bytes,10,opt,name=content,proto3" json:"content,omitempty"`
	Seq              uint32 `protobuf:"varint,11,opt,name=seq,proto3" json:"seq,omitempty"`
	SendTime         int64  `protobuf:"varint,12,opt,name=sendTime,proto3" json:"sendTime,omitempty"`
	Status           int32  `protobuf:"varint,13,opt,name=status,proto3" json:"status,omitempty"`
}

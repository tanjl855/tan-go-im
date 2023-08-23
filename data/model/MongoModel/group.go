package MongoModel

import "time"

// Group 群
type Group struct {
	MongoDBModel `bson:",inline" json:",inline"`
	GroupID      string   `bson:"group_id" json:"group_id"`
	GroupName    string   `bson:"group_name" json:"group_name"`
	UIDList      []string `bson:"uid_list" json:"uid_list"`
	Admin        []string `bson:"admin" json:"admin"`
	BanUID       []string `bson:"ban_uid" json:"ban_uid"`
	OwnerUID     string   `bson:"owner_uid" json:"owner_uid"`
	MaxMember    int      `bson:"max_member" json:"max_member"`
	State        int      `bson:"state" json:"state"` //0表示正常，1为已解散
}

func (g *Group) TableName() string {
	return "group"
}

type UserToGroup struct {
	MongoDBModel `bson:",inline" json:",inline"`
	GroupID      string    `bson:"group_id" json:"group_id"`
	UID          string    `bson:"uid" json:"uid"`
	Nickname     string    `bson:"nickname" json:"nickname"`
	JoinTime     time.Time `bson:"join_time" json:"join_time"`
}

func (g *UserToGroup) TableName() string {
	return "user_to_group"
}

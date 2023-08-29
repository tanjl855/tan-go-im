package ws

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/tanjl855/tan_go_im/data/Repo/Group"
	"github.com/tanjl855/tan_go_im/pkg/im_auth"
	log "github.com/tanjl855/tan_go_im/pkg/im_log"
	"github.com/tanjl855/tan_go_im/pkg/im_rsp"
	"github.com/tanjl855/tan_go_im/proto/pb_msg"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/db"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	groupRepo       = Group.NewIGroupRepo(db.DB.MongoDB)
	userToGroupRepo = Group.NewIUserToGroupRepo(db.DB.MongoDB)
)

// UserConn 用户ws长连接
type UserConn struct {
	*websocket.Conn
	w   *sync.Mutex
	Uid string
}

// WServer 控制长连接的结构
type WServer struct {
	rwLock       *sync.RWMutex
	wsOutAddr    string
	grpcAddr     string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsUserToConn map[string]*UserConn
}

var WS *WServer

func Init(outAddr string, grpcAddr string, connMaxNum int, timeout int, maxMsgLen int) {
	WS = &WServer{}
	//获取外网地址
	WS.wsOutAddr = outAddr
	//获取内网grpc用ip+port
	WS.grpcAddr = grpcAddr
	WS.rwLock = &sync.RWMutex{}
	WS.wsMaxConnNum = connMaxNum
	WS.wsUserToConn = make(map[string]*UserConn)
	WS.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(timeout) * time.Second,
		ReadBufferSize:   maxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	// 起协程保活ws的外链地址
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(fmt.Sprintf("更新服务器心跳失败，维持在线状态失败：%+v", err))
				return
			}
		}()
		err := db.DB.Rdb.ZAdd(context.Background(), "ws_pool_out_addr_list", &redis.Z{
			Score:  float64(len(WS.wsUserToConn)),
			Member: WS.wsOutAddr,
		}).Err()
		if err != nil {
			log.Panic(fmt.Sprintf("初始化ws连接池服务到缓存失败: %v", err))
		}
		minTick := time.NewTicker(60 * time.Second)
		for {
			err := db.DB.Rdb.Set(context.Background(), "ws_pool_out_addr:"+WS.wsOutAddr, 1, 120*time.Second).Err()
			if err != nil {
				log.Error(fmt.Sprintf("更新服务器心跳失败: %v", err))
			}
			<-minTick.C
		}
	}()
}

type wsHandlerReq struct {
	SendId string `form:"send_id"`
	Token  string `form:"token"`
}

// WsHandler ws连接池入口，用户http协议升级为ws协议
func (ws *WServer) WsHandler(ctx *gin.Context) {
	req := &wsHandlerReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		im_rsp.Failed(ctx, im_rsp.ParameterValidationError, err)
		return
	}
	// 直接拿内网ip做机器号
	operationID := ws.grpcAddr
	log.Debug(fmt.Sprintf("ws pool服务, 机器号: %s, sendId: %s", operationID, req.SendId))
	conn, err := ws.wsUpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Error(fmt.Sprintf("ws pool服务, 机器号: %s, sendId: %s, http协议升级ws失败, err: %v", operationID, req.SendId, err))
		im_rsp.Failed(ctx, im_rsp.ERROR, err)
		return
	}
	newConn := &UserConn{Conn: conn, w: &sync.Mutex{}, Uid: req.SendId}
	err = ws.addUserConn(ctx, newConn)
	if err != nil {
		ws.delUserConn(ctx, newConn)
		log.Error(fmt.Sprintf("ws pool服务, 机器号: %s, sendId: %s, ws链接入池失败,err: %v", operationID, req.SendId, err))
		im_rsp.Failed(ctx, im_rsp.ERROR, err)
		return
	}
	go ws.readMsg(newConn)
}

func (ws *WServer) GetUserConn(uid string) *UserConn {
	ws.rwLock.RLock()
	conn, ok := ws.wsUserToConn[uid]
	ws.rwLock.RUnlock()
	if ok {
		return conn
	}
	return nil
}

// 添加用户ws链接进连接池
func (ws *WServer) addUserConn(ctx context.Context, conn *UserConn) error {
	log.Info(fmt.Sprintf("add用户链接到链接池，uid:%s ip:%s ", conn.Uid, conn.RemoteAddr().String()))
	err := db.DB.Rdb.Set(ctx, "user_conn"+conn.Uid, ws.grpcAddr, 86400*time.Second).Err()
	if err != nil {
		log.Error(fmt.Sprintf("add用户链接到链接池失败，uid:%s ip:%s err:%+v", conn.Uid, conn.RemoteAddr().String()), err)
		return err
	}
	// 更新ws连接池实时连接数，读取map长度不会触发同时读写竞争导致的panic
	err = db.DB.Rdb.ZAdd(ctx, "ws_pool_out_addr_list", &redis.Z{
		Score:  float64(len(ws.wsUserToConn) + 1),
		Member: ws.wsOutAddr,
	}).Err()
	if err != nil {
		log.Error(fmt.Sprintf("更新ws链接池实时链接数失败,ip:%s,err:%+v", ws.wsOutAddr, err))
		return err
	}
	ws.rwLock.Lock()
	ws.wsUserToConn[conn.Uid] = conn
	ws.rwLock.Unlock()
	return nil
}

// 删除用户ws链接
func (ws *WServer) delUserConn(ctx context.Context, conn *UserConn) {
	ws.rwLock.Lock()
	delete(ws.wsUserToConn, conn.Uid)
	ws.rwLock.Unlock()
	err := db.DB.Rdb.Del(ctx, "user_conn:"+conn.Uid).Err()
	if err != nil {
		log.Error(fmt.Sprintf("机器号：%s ,删除用户ws链接失败,uid:%s ,ip:%s", ws.grpcAddr, conn.Uid, conn.RemoteAddr().String()))
	}
	// 更新ws连接池实时连接数，读取map长度不会触发同时读写竞争导致的panic
	err = db.DB.Rdb.ZAdd(ctx, "ws_pool_out_addr_list", &redis.Z{
		Score:  float64(len(ws.wsUserToConn)),
		Member: WS.wsOutAddr,
	}).Err()
	if err != nil {
		log.Error(fmt.Sprintf("更新ws链接池实时链接数失败,ip:%s,err:%+v", ws.wsOutAddr, err))
	}
	err = conn.Close()
	if err != nil {
		log.Error(fmt.Sprintf("机器号：%s ,关闭从池中删除的未关闭用户ws链接失败,uid:%s ,ip:%s", ws.grpcAddr, conn.Uid, conn.RemoteAddr().String()))
	}
}

// 读取ws链接的数据
func (ws *WServer) readMsg(conn *UserConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("ws通道读消息错误：%+v", err))
		}
	}()
	for {
		messageType, msgReader, err := conn.NextReader()
		if messageType == websocket.PingMessage {
			log.Info(fmt.Sprintf("uid:%s pingMessage", conn.Uid))
			continue
		}
		if err != nil {
			log.Error(fmt.Sprintf("从通道读取消息错误，uid:%s ,IP:%s ,err:%+v", conn.Uid, conn.RemoteAddr().String(), err))
			ws.delUserConn(context.Background(), conn)
			return
		}
		ws.msgParse(conn, msgReader)
	}
}

// WriteMsg 往ws链接写数据
func WriteMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(300) * time.Second))
	return conn.WriteMessage(a, msg)
}

func wsVerifyToken(token string, sendId string) (bool, error) {
	userClaims, err := im_auth.GetClaimFromToken(token)
	if err != nil {
		return false, err
	}
	if userClaims.UID != sendId {
		return false, err
	}
	return true, nil
}

// msgParse 解析函数，梭哈到消息队列
func (ws *WServer) msgParse(conn *UserConn, msgReader io.Reader) {
	pbMsg := pb_msg.Msg{}
	msg, err := io.ReadAll(msgReader)
	if err != nil {
		log.Error(fmt.Sprintf("读取失败：%+v", err))
		WriteMsg(conn, websocket.BinaryMessage, msg)
		return
	}
	if len(msg) >= 1024*1024*15 {
		log.Error(fmt.Sprintf("消息过大：%+v", err))
		WriteMsg(conn, websocket.BinaryMessage, msg)
		return
	}
	err = proto.Unmarshal(msg, &pbMsg)
	if err != nil {
		log.Error(fmt.Sprintf("解析失败：%+v", err))
		WriteMsg(conn, websocket.BinaryMessage, msg)
		return
	}
	if pbMsg.SendID != conn.Uid {
		log.Info(fmt.Sprintf("发送者ID不为本人"))
		WriteMsg(conn, websocket.BinaryMessage, msg)
		return
	}
	pbMsg.Status = 1
	sendmsg, err := proto.Marshal(&pbMsg)
	if err != nil {
		WriteMsg(conn, websocket.BinaryMessage, msg)
		return
	}
	_, _, err = db.DB.KafkaProduct.SendMessage(&sarama.ProducerMessage{Topic: "im_msg", Value: sarama.ByteEncoder(sendmsg)})
	if err != nil {
		WriteMsg(conn, websocket.BinaryMessage, msg)
		log.Info(fmt.Sprintf("解析并发送到kafka失败：%s", err))
		return
	}
	log.Info(fmt.Sprintf("解析并发送到kafka成功：%s", err))
	WriteMsg(conn, websocket.BinaryMessage, sendmsg)
}

// Broadcast 广播
func (ws *WServer) Broadcast(conn *UserConn, a int, msg []byte) {
	for _, v := range ws.wsUserToConn {
		if conn.Uid != v.Uid {
			err := WriteMsg(v, a, msg)
			if err != nil {
				ws.delUserConn(context.Background(), v)
			}
		}
	}
}

// ReFlashOneUserGroupConn 更新单个群组的在线状态
func (ws *WServer) ReFlashOneUserGroupConn(uid string, groupId string) error {
	err := db.DB.Rdb.Set(context.Background(), "user_conn:"+uid, ws.grpcAddr, 86400).Err()
	if err != nil {
		log.Error(fmt.Sprintf("更新用户在线状态错误：%+v", err))
		return err
	}
	err = db.DB.Rdb.Set(context.Background(), "group_conn:"+groupId+":"+uid+":"+ws.grpcAddr, 1, 0).Err()
	if err != nil {
		log.Error(fmt.Sprintf("更新用户：%s在群组：%s 的在线缓存失败：%+v", uid, groupId, err))
	}
	return nil
}

// ReFlashAllUserConn 更新多个群组的在线状态
func (ws *WServer) ReFlashAllUserConn(uid string) error {
	err := db.DB.Rdb.Set(context.Background(), "user_conn:"+uid, ws.grpcAddr, 86400).Err()
	if err != nil {
		log.Error(fmt.Sprintf("更新用户在线状态错误：%+v", err))
		return err
	}
	//查用户所有群组列表
	groupList, err := userToGroupRepo.GetUserToGroupListByUID(context.Background(), uid)
	if err != nil {
		log.Error(fmt.Sprintf("获取用户所加群列表错误：%+v", err))
		return err
	}
	for i := 0; i < len(groupList); i++ {
		err = db.DB.Rdb.Set(context.Background(), "group_conn:"+groupList[i].GroupID+":"+uid+":"+ws.grpcAddr, 1, 0).Err()
		if err != nil {
			log.Error(fmt.Sprintf("更新用户：%s在群组：%s 的在线缓存失败：%+v", uid, groupList[i].GroupID, err))
		}
		return err
	}
	return nil
}

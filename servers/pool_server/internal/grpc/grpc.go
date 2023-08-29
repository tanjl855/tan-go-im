package grpc

import (
	"github.com/tanjl855/tan_go_im/proto/pb_msg"
	"github.com/tanjl855/tan_go_im/servers/pool_server/internal/controller/grpc/grpc_msg"
	"google.golang.org/grpc"
)

func RegisterServerList(s *grpc.Server) {
	pb_msg.RegisterMsgServerServer(s, &grpc_msg.MsgServer{})
}

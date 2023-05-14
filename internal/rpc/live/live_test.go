package live

import (
	"Open_IM/pkg/proto/live"
	"context"
	"fmt"
	"testing"
)

func Test_rpcLive_JoinRoom(t *testing.T) {
	srv := NewRpcLiveServer(8000)

	got, err := srv.JoinRoom(context.Background(), &live.JoinRoomReq{UserID: "123", OperationID: "123", ChannelID: "3817720326"})
	fmt.Println(got, err)

}

func Test_rpcLive_LeveRoom(t *testing.T) {
	srv := NewRpcLiveServer(8000)

	got, err := srv.LeveRoom(context.Background(), &live.LeveRoomReq{UserID: "1684026049", OperationID: "1684026049", ChannelID: "3817720326"})
	fmt.Println(got, err)

}
func Test_rpcLive_StartLive(t *testing.T) {
	srv := NewRpcLiveServer(8000)

	got, err := srv.StartLive(context.Background(), &live.StartLiveReq{UserID: "1910360909", OperationID: "1910360909", ChannelID: "3817720326"})
	fmt.Println(got, err)

}
func Test_rpcLive_CloseLive(t *testing.T) {
	srv := NewRpcLiveServer(8000)

	got, err := srv.CloseLive(context.Background(), &live.CloseLiveReq{UserID: "1910360909", OperationID: "1910360909", ChannelID: "3817720326"})
	fmt.Println(got, err)

}

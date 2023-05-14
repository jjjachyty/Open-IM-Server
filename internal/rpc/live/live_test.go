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

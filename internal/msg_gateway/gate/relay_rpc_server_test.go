package gate

import (
	pbRelay "Open_IM/pkg/proto/relay"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"context"
	"fmt"
	"testing"
)

func TestRPCServer_SuperGroupBackgroundOnlinePush(t *testing.T) {

	r := &RPCServer{}
	r.onInit(9999)
	got, err := r.SuperGroupBackgroundOnlinePush(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{OperationID: "1", MsgData: &sdk_ws.MsgData{SendID: "1225925647", SenderNickname: "天道", ServerMsgID: "1", ClientMsgID: "1", RecvID: "2892606321", GroupID: "1763879570", ContentType: 101, Content: []byte("666")}, PushToUserIDList: []string{"2892606321"}})
	fmt.Println(got, err)

}

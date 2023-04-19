/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

import (
	pbPush "Open_IM/pkg/proto/push"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"testing"
)

func TestMsgToUser(t *testing.T) {

	MsgToUser(&pbPush.PushMsgReq{OperationID: "111111", PushToUserID: "2892606321", MsgData: &server_api_params.MsgData{SendID: "1225925647", SenderNickname: "天道", ServerMsgID: "1", ClientMsgID: "1", RecvID: "2892606321", GroupID: "1763879570", ContentType: 101, Content: []byte("666")}})

}

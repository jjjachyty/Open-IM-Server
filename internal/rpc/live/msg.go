package live

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"strings"
)

func (rpc *rpcLive) sendCloseMsg(operationID string, sendID string, channelID string) error {
	// content := make([]byte, 4)
	// binary.LittleEndian.PutUint64(content, uint64(channelID))

	pbData := msg.SendMsgReq{
		OperationID: operationID,
		MsgData: &open_im_sdk.MsgData{
			SendID:           sendID,
			SenderPlatformID: constant.AndroidPlatformID,
			ClientMsgID:      utils.GetMsgID(sendID),
			SessionType:      constant.LiveChatType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      constant.CloseLivingNotification,
			Content:          []byte(channelID),
			CreateTime:       utils.GetCurrentTimestampByMill(),
		},
	}

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationID, errMsg)
		return errors.New(errMsg)
	}
	client := msg.NewMsgClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &pbData)
	if err != nil {
		log.NewError(operationID, "SendMsg rpc failed, ", pbData.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(operationID, "SendMsg rpc failed, ", pbData.String(), reply.ErrCode, reply.ErrMsg)
	}
	return nil
}

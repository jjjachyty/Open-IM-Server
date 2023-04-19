package msg

import (
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"sync"
	"time"
)

func (rpc *rpcChat) fansAutoReply(msgData *sdk_ws.MsgData, m map[string][]string) {
	<-time.NewTimer(time.Second * 2).C
	log.NewInfo("开始粉丝互动>>>>>>>>>>>>>>>>>>>>>>>>>>")
	senderInfo, err := rocksCache.GetGroupMemberInfoFromCache(msgData.GroupID, msgData.SendID)
	if err != nil {
		log.Error("获取群用户信息出错", err)
		return
	}
	if senderInfo.RoleLevel == 3 {
		var newMsg sdk_ws.MsgData
		var sendTag bool
		var wg sync.WaitGroup

		count := len(rpc.robots) / 2
		for k, v := range rpc.robots {

			if count < 0 {
				return
			}

			//只发送在线的
			wg.Add(1)
			newMsg = *msgData
			newMsg.SendID = k
			newMsg.RecvID = k
			newMsg.SenderNickname = v
			newMsg.ServerMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.ClientMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.SendTime++

			go rpc.sendMsgToGroupOptimization(m[constant.OnlineStatus], &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)
			count--
		}
		wg.Wait()
	}
}

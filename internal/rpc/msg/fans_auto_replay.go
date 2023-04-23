package msg

import (
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"math/rand"
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
	recivers := make([]string, 0)
	for _, v := range m[constant.OnlineStatus] {
		if rpc.robots[v] == nil {
			recivers = append(recivers, v)
		}
	}

	if senderInfo.RoleLevel == 3 {
		var sendTag bool
		var wg sync.WaitGroup

		count := rand.Intn(len(rpc.robots))

		for k, v := range rpc.robots {
			var newMsg sdk_ws.MsgData

			if count < 0 {
				return
			}

			//只发送在线的
			wg.Add(1)
			newMsg = *msgData
			newMsg.SendID = k
			newMsg.RecvID = k
			newMsg.SenderNickname = v.Nickname
			newMsg.SenderFaceURL = v.FaceURL
			newMsg.ServerMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.ClientMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.SendTime++

			go rpc.sendMsgToGroupOptimization(recivers, &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)
			count--
		}
		wg.Wait()
	}
}

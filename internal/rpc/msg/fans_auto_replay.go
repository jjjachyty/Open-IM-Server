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
	log.NewInfo("开始粉丝互动>>>>>>>>>>>>>>>>>>>>>>>>>>")
	senderInfo, err := rocksCache.GetGroupMemberInfoFromCache(msgData.GroupID, msgData.SendID)
	if err != nil {
		log.Error("获取群用户信息出错", err)
		return
	}
	robots, err := rocksCache.GetGroupRobotsRoundFromCache(msgData.GroupID)
	if err != nil {
		log.Error("随机获取机器人出错", err)
		return
	}
	recivers := make([]string, 0)
	for _, v := range m[constant.OnlineStatus] {
		recivers = append(recivers, v)
	}

	if senderInfo.RoleLevel == 3 {
		var sendTag bool
		var wg sync.WaitGroup

		for _, v := range robots {
			var newMsg sdk_ws.MsgData

			//只发送在线的
			wg.Add(1)
			newMsg = *msgData
			newMsg.SendID = v.UserID
			newMsg.RecvID = v.UserID
			newMsg.SenderNickname = v.Nickname
			newMsg.SenderFaceURL = v.FaceURL
			newMsg.ServerMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.ClientMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.SendTime++
			rand.Seed(time.Now().UnixNano())
			<-time.NewTimer(time.Duration(rand.Int63n(2000))).C

			go rpc.sendMsgToGroupOptimization(recivers, &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)

		}
		wg.Wait()
	}
}

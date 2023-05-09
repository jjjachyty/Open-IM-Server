package msg

import (
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (rpc *rpcChat) liveFansAutoReply(msgData *sdk_ws.MsgData, members []string) {
	log.NewInfo("开始粉丝互动>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// senderInfo, err := rocksCache.GetUserin(msgData.SendID)
	// if err != nil {
	// 	log.Error("获取用户信息出错", err)
	// 	return
	// }
	sendID, _ := strconv.ParseInt(msgData.SendID, 0, 64)
	channelID, _ := strconv.ParseInt(msgData.GroupID, 0, 64)
	robots, err := rocksCache.GetLiveRobotsFromCache(channelID)
	if err != nil {
		log.Error("随机获取机器人出错", err)
		return
	}

	is, err := rocksCache.IsLiveAtmosphereUser(channelID, sendID)
	if err != nil {
		log.Error("检测是否是氛围组出错", err)
		return
	}
	if is {
		var sendTag bool
		var wg sync.WaitGroup
		count := rand.Intn(len(robots))
		for k, v := range robots {
			count--
			if count <= 0 {
				break
			}

			var newMsg sdk_ws.MsgData
			dt := strings.Split(v, ",")
			//只发送在线的
			wg.Add(1)
			newMsg = *msgData
			newMsg.SendID = k
			newMsg.RecvID = k
			newMsg.SenderNickname = dt[0]
			newMsg.SenderFaceURL = dt[1]
			newMsg.ServerMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.ClientMsgID = utils.GetMsgID(newMsg.SendID)
			newMsg.SendTime++
			rand.Seed(time.Now().UnixNano())
			<-time.NewTimer(time.Second * time.Duration(rand.Intn(5))).C
			// for _, v := range v {
			// 	go rpc.sendMsgToGroupOptimization(recivers, &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)
			// }
			//20 个一批
			memberCount := len(members)
			for i := 0; i < memberCount/20; i++ {
				start := i * 20
				end := (i + 1) * 20
				if end >= memberCount-1 {
					end = memberCount - 1
				}
				go rpc.sendMsgToGroupOptimization(members[start:end], &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)
			}
		}
		wg.Wait()
	}

}

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
			<-time.NewTimer(time.Second * time.Duration(rand.Intn(5))).C
			go rpc.sendMsgToGroupOptimization(recivers, &msg.SendMsgReq{MsgData: &newMsg}, constant.OnlineStatus, &sendTag, &wg)

		}
		wg.Wait()
	}
}

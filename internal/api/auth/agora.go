package apiAuth

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	putils "Open_IM/pkg/utils"
	"net/http"
	"strconv"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"github.com/gin-gonic/gin"
)

type RTCTokenRequest struct {
	Uid_rtc_int  string `json:"uid"`
	Channel_name string `json:"channelName"`
	Role         uint32 `json:"role"`
	OperationID  string `json:"operationID"`
}
type RTCTokenResponse struct {
	api.CommResp
	Data TokenInfo `json:"data"`
}
type TokenInfo struct {
	AppId string `json:"appId"`
	Token string `json:"token"`
	// LeftDuration int32  `json:"leftDuration"`
}

// var rtc_token string
// var int_uid uint32
// var channel_name string

// var role_num uint32
// var role rtctokenbuilder.Role

// 使用 RtcTokenBuilder 来生成 RTC Token。

func RTCToken(c *gin.Context) {
	params := RTCTokenRequest{}
	err := c.BindJSON(&params)
	if err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	uid, err := strconv.ParseInt(params.Uid_rtc_int, 0, 64)
	if err != nil {
		log.NewError(params.OperationID, err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err})
		return
	}
	switch params.Role {
	case 1:
		params.Role = rtctokenbuilder.RolePublisher
	case 2:
		params.Role = rtctokenbuilder.RoleSubscriber
	}
	// liveUser, err := imdb.GetLiveByUserID(params.Uid_rtc_int)
	// if err != nil {
	// 	log.NewError(params.OperationID, err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err})
	// 	return
	// }
	token, err := putils.GenerateRtcToken(uint32(uid), params.Channel_name, 2*60*60, 2*61*60, rtctokenbuilder.Role(params.Role))
	if err != nil {
		log.NewError(params.OperationID, err)
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "获取token出错"})
		return
	}
	c.JSON(http.StatusOK, RTCTokenResponse{CommResp: api.CommResp{}, Data: TokenInfo{AppId: config.Config.Agora.AppID, Token: token}})
}

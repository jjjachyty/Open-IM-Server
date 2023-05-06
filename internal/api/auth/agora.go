package apiAuth

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
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
	AppId        string `json:"appId"`
	Token        string `json:"token"`
	LeftDuration int32  `json:"leftDuration"`
}

// var rtc_token string
// var int_uid uint32
// var channel_name string

// var role_num uint32
// var role rtctokenbuilder.Role

// 使用 RtcTokenBuilder 来生成 RTC Token。
func generateRtcToken(int_uid uint32, channelName string, role rtctokenbuilder.Role) (string, error) {

	appID := config.Config.Agora.AppID
	appCertificate := config.Config.Agora.AppCertificate
	// AccessToken2 过期的时间，单位为秒。
	// 当 AccessToken2 过期但权限未过期时，用户仍在频道里并且可以发流，不会触发 SDK 回调。
	// 但一旦用户和频道断开连接，用户将无法使用该 Token 加入同一频道。请确保 AccessToken2 的过期时间晚于权限过期时间。
	tokenExpireTimeInSeconds := uint32(30)
	// 权限过期的时间，单位为秒。
	// 权限过期30秒前会触发 token-privilege-will-expire 回调。
	// 权限过期时会触发 token-privilege-did-expire 回调。
	// 为作演示，在此将过期时间设为 40 秒。你可以看到客户端自动更新 Token 的过程。
	privilegeExpireTimeInSeconds := uint32(40)

	return rtctokenbuilder.BuildTokenWithUid(appID, appCertificate, channelName, int_uid, role, tokenExpireTimeInSeconds, privilegeExpireTimeInSeconds)

}

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
	liveUser, err := imdb.GetLiveByUserID(params.Uid_rtc_int)
	if err != nil {
		log.NewError(params.OperationID, err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err})
		return
	}
	token, err := generateRtcToken(uint32(uid), params.Channel_name, rtctokenbuilder.Role(params.Role))
	if err != nil {
		log.NewError(params.OperationID, err)
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": "获取token出错"})
		return
	}
	c.JSON(http.StatusOK, RTCTokenResponse{CommResp: api.CommResp{}, Data: TokenInfo{AppId: config.Config.Agora.AppID, Token: token, LeftDuration: liveUser.LeftDuration}})
}

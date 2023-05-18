package live

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/live"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func JoinLiveRoom(c *gin.Context) {
	params := api.JoinRoomReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	if params.UserID == "" {
		params.UserID = utils.Int64ToString(time.Now().Unix())
		params.NickName = "游客"
		params.FaceURL = ""
	}

	log.NewInfo(params.OperationID, "JoinLiveRoom req: ", params)
	req := &rpc.JoinRoomReq{}
	utils.CopyStructFields(req, &params)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImLiveName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewLiveClient(etcdConn)
	RpcResp, err := client.JoinRoom(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "JoinLiveRoom failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.NewInfo(req.OperationID, "JoinLiveRoom api return ", RpcResp)
	data := jsonData.JsonDataList(RpcResp.UserLive)
	data = append(data, jsonData.JsonDataOne(RpcResp.Owner))
	data = append(data, map[string]interface{}{"rtcToken": RpcResp.RtcToken})
	c.JSON(http.StatusOK, api.LiveCommonResp{CommResp: api.CommResp{ErrCode: 0, ErrMsg: ""}, Data: data})
}
func LevelLiveRoom(c *gin.Context) {
	params := api.LevelLiveRoomreq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "LevelLiveRoom req: ", params)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImLiveName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewLiveClient(etcdConn)
	RpcResp, err := client.LeveRoom(context.Background(), &rpc.LeveRoomReq{OperationID: params.OperationID, ChannelID: params.ChannelID, UserID: params.UserID})
	if err != nil {
		log.NewError(params.OperationID, "LevelLiveRoom failed ", err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	commonResp := api.CommResp{}
	utils.CopyStructFields(&commonResp, &RpcResp.CommonResp)
	c.JSON(http.StatusOK, api.LiveCommonResp{CommResp: commonResp})
}
func LiveRoomUsers(c *gin.Context) {
	params := api.LiveRoomUsersReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "JoinLiveRoom req: ", params)
	req := &rpc.GetRoomUserReq{}
	utils.CopyStructFields(req, &params)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImLiveName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewLiveClient(etcdConn)
	RpcResp, err := client.GetRoomUser(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	log.NewInfo(req.OperationID, "GetUserInfo api return ", RpcResp)
	c.JSON(http.StatusOK, RpcResp)
}
func StartLive(c *gin.Context) {
	params := api.StartLiveReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "StartLive req: ", params)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImLiveName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewLiveClient(etcdConn)
	RpcResp, err := client.StartLive(context.Background(), &rpc.StartLiveReq{OperationID: params.OperationID, ChannelID: params.ChannelID, UserID: params.UserID})
	if err != nil {
		log.NewError(params.OperationID, "StartLive failed ", err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	commonResp := api.CommResp{}
	utils.CopyStructFields(&commonResp, &RpcResp.CommonResp)
	c.JSON(http.StatusOK, api.LiveCommonResp{CommResp: commonResp, Data: []map[string]interface{}{{"RtcToken": RpcResp.RtcToken}}})
}
func CLoseLive(c *gin.Context) {
	params := api.CloseLiveReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "CLoseLive req: ", params)

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImLiveName, params.OperationID)
	if etcdConn == nil {
		errMsg := params.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewLiveClient(etcdConn)
	RpcResp, err := client.CloseLive(context.Background(), &rpc.CloseLiveReq{OperationID: params.OperationID, ChannelID: params.ChannelID, UserID: params.UserID})
	if err != nil {
		log.NewError(params.OperationID, "CLoseLive failed ", err.Error(), "")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	commonResp := api.CommResp{}
	utils.CopyStructFields(&commonResp, &RpcResp.CommonResp)
	c.JSON(http.StatusOK, api.LiveCommonResp{CommResp: commonResp})
}

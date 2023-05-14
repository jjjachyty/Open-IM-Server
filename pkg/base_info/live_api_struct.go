package base_info

import (
	pbLive "Open_IM/pkg/proto/live"
)

type JoinRoomReq struct {
	OperationID string `json:"operationID" `
	ChannelID   string `json:"channelID" binding:"required"`
	UserID      string `json:"userID"`
	NickName    string `json:"nickName"`
	FaceURL     string `json:"faceURL"`
}
type LiveCommonResp struct {
	CommResp
	Data []map[string]interface{} `json:"data" swaggerignore:"true"`
}
type LevelLiveRoomreq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   string `json:"channelID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type StartLiveReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   string `json:"channelID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
	Platform    int    `json:"platform"`
}

type CloseLiveReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   string `json:"channelID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

type LiveRoomUsersReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   string `json:"channelID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type LiveRoomUsersResp struct {
	CommResp
	UserInfoList []*pbLive.LiveUserInfo `json:"userInfoList"`
}

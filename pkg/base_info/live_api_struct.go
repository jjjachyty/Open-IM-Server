package base_info

import (
	pbLive "Open_IM/pkg/proto/live"
)

type JoinRoomReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   int64  `json:"channelID" binding:"required"`
	UserID      int64  `json:"userID" binding:"required"`
	NickName    string `json:"nickName" binding:"required"`
	FaceURL     string `json:"faceURL" binding:"required"`
}
type LiveCommonResp struct {
	CommResp
	Data []map[string]interface{} `json:"data" swaggerignore:"true"`
}
type LevelLiveRoomreq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   int64  `json:"channelID" binding:"required"`
	UserID      int64  `json:"userID" binding:"required"`
}

type StartLiveReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   int64  `json:"channelID" binding:"required"`
	UserID      int64  `json:"userID" binding:"required"`
	Platform    int    `json:"platform"`
}

type CloseLiveReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   int64  `json:"channelID" binding:"required"`
	UserID      int64  `json:"userID" binding:"required"`
}

type LiveRoomUsersReq struct {
	OperationID string `json:"operationID" binding:"required"`
	ChannelID   int64  `json:"channelID" binding:"required"`
	UserID      int64  `json:"userID" binding:"required"`
}
type LiveRoomUsersResp struct {
	CommResp
	UserInfoList []*pbLive.LiveUserInfo `json:"userInfoList"`
}

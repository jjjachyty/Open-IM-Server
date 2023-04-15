package base_info

import "Open_IM/pkg/common/db"

type SetClientInitConfigReq struct {
	OperationID     string  `json:"operationID"  binding:"required"`
	DiscoverPageURL *string `json:"discoverPageURL"`
}

type SetClientInitConfigResp struct {
	CommResp
}

type GetClientInitConfigReq struct {
	OperationID string `json:"operationID"  binding:"required"`
}

type GetClientInitConfigResp struct {
	CommResp
	Data db.ClientInitConfig `json:"data"`
}

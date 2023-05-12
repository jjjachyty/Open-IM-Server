package live

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"

	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pblive "Open_IM/pkg/proto/live"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"Open_IM/pkg/common/config"
	sdkws "Open_IM/pkg/proto/sdk_ws"

	"google.golang.org/grpc"
)

func (rpc *rpcLive) JoinRoom(_ context.Context, req *pblive.JoinRoomReq) (*pblive.JoinRoomResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())

	liveInfo, err := rocksCache.GetLiveRoomFromCache(req.ChannelID)
	if err != nil {
		return &pblive.JoinRoomResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "查询直播信息出错"}}, nil
	}

	user, err := rocksCache.GetUserInfoFromCache(fmt.Sprintf("%d", liveInfo.UserID))
	if err != nil {
		return &pblive.JoinRoomResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "查询主播信息出错"}}, nil
	}
	respUserLiveInfo := &sdkws.UserLive{}
	utils.CopyStructFields(&respUserLiveInfo, liveInfo)

	promePkg.PromeInc(promePkg.LiveUserCounter)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ")
	if err = rocksCache.JoinLiveRoom(req.ChannelID, req.UserID, req.NickName, req.FaceURL, false); err != nil {
		return &pblive.JoinRoomResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "加入房间出错"}}, nil
	}

	token, err := utils.GenerateRtcToken(uint32(req.UserID), fmt.Sprintf("%d", req.ChannelID), uint32(2*60*60), uint32(2*62*60), 2) //默认2小时
	if err != nil {
		log.NewError(req.OperationID, err)
		return &pblive.JoinRoomResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	return &pblive.JoinRoomResp{CommonResp: &pblive.CommonResp{}, UserLive: respUserLiveInfo, Owner: &pblive.LiveUserInfo{UserID: liveInfo.UserID, NickName: user.Nickname, FaceURL: user.FaceURL}, RtcToken: token}, nil
}

func (rpc *rpcLive) GetRoomUser(_ context.Context, req *pblive.GetRoomUserReq) (*pblive.GetRoomUserResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())

	liveUsers, err := rocksCache.GetLiveUsersLimitFromCache(req.ChannelID, 100)
	if err != nil {
		return &pblive.GetRoomUserResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "查询直播信息出错"}}, nil
	}

	resp := make([]*pblive.LiveUserInfo, 0)
	for k, v := range liveUsers {
		userID, _ := strconv.ParseInt(k, 0, 64)
		nickName, faceURL := rocksCache.GetLiveUsersValues(v)
		resp = append(resp, &pblive.LiveUserInfo{UserID: userID, NickName: nickName, FaceURL: faceURL})
	}
	//不足100个 查询机器人补充
	if len(liveUsers) < 100 {
		liveUsers, err = rocksCache.GetLiveRobotsLimitFromCache(req.ChannelID, int64(100-len(liveUsers)))
		if err != nil {
			return &pblive.GetRoomUserResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "查询直播信息出错"}}, nil
		}
		for k, v := range liveUsers {
			userID, _ := strconv.ParseInt(k, 0, 64)
			nickName, faceURL := rocksCache.GetLiveUsersValues(v)
			resp = append(resp, &pblive.LiveUserInfo{UserID: userID, NickName: nickName, FaceURL: faceURL})
		}
	}
	promePkg.PromeInc(promePkg.LiveUserCounter)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ")

	return &pblive.GetRoomUserResp{CommonResp: &pblive.CommonResp{}, Users: resp}, nil
}

func (rpc *rpcLive) LeveRoom(_ context.Context, req *pblive.LeveRoomReq) (*pblive.LeveRoomResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())

	err := rocksCache.LevelLiveRoom(req.ChannelID, req.UserID)
	if err != nil {
		log.NewError(req.OperationID, "LevelLiveRoom failed ", err.Error(), req.UserID)
		return &pblive.LeveRoomResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "离开出错"}}, nil
	}
	// promePkg.PromeGaugeDec(promePkg.LiveUserCounter)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ")

	return &pblive.LeveRoomResp{CommonResp: &pblive.CommonResp{}}, nil
}

func (s *rpcLive) GetLiveByUserID(ctx context.Context, req *pblive.GetLiveByUserIDReq) (resp *pblive.GetLiveByUserIDResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())

	userLive, err := imdb.GetLiveByUserID(req.UserID)
	if err != nil {
		return &pblive.GetLiveByUserIDResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}
	var userLiveResp sdkws.UserLive
	utils.CopyStructFields(&userLiveResp, userLive)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", userLiveResp.String())

	return &pblive.GetLiveByUserIDResp{CommonResp: &pblive.CommonResp{}, UserLive: &userLiveResp}, err
}

func (s *rpcLive) StartLive(ctx context.Context, req *pblive.StartLiveReq) (resp *pblive.StartLiveResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	if req.ChannelID == 0 || req.UserID == 0 {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 400, ErrMsg: err.Error()}}, err
	}
	//获取用户信息
	user, err := rocksCache.GetUserInfoFromCache(fmt.Sprintf("%d", req.UserID))
	if err != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	//检查用户是否已有未结束的直播
	live, err := imdb.GetUserLiving(fmt.Sprintf("%d", req.UserID))
	if err != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	token, err := utils.GenerateRtcToken(uint32(req.UserID), fmt.Sprintf("%d", req.ChannelID), uint32(user.LeftDuration*60), uint32(user.LeftDuration*60), 1)
	if err != nil {
		log.NewError(req.OperationID, err)
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}
	if live != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{}, RtcToken: token}, err
	}
	live = &db.UserLive{UserID: req.UserID, GroupID: req.GroupID, ChannelID: req.ChannelID, ChannelName: req.ChannelName, StartAt: time.Now().Unix()}
	if err := imdb.CreateLiveInfo(live); err != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}
	if err = rocksCache.CreateLiveRoom(*live); err != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	if user.LeftDuration <= 1 {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "剩余分钟不足1分钟"}}, err
	}

	in, err := rocksCache.UserInRoom(live.ChannelID, live.UserID)
	if err != nil {
		log.NewError(req.OperationID, err)
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err

	}
	if in {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{}, RtcToken: token}, err
	}
	//加入房间
	if err = rocksCache.JoinLiveRoom(live.ChannelID, live.UserID, user.Nickname, user.FaceURL, false); err != nil {
		return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	return &pblive.StartLiveResp{CommonResp: &pblive.CommonResp{}, RtcToken: token}, err
}

func (s *rpcLive) CloseLive(ctx context.Context, req *pblive.CloseLiveReq) (resp *pblive.CloseLiveResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	if req.ChannelID == 0 || req.UserID == 0 {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 400, ErrMsg: err.Error()}}, err
	}

	//获取用户信息
	liveInfo, err := rocksCache.GetLiveInfoFromCache(req.ChannelID)
	if err != nil {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	if liveInfo.UserID != req.UserID {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: "您暂时不能关闭直播"}}, err
	}

	//清理房间
	if err = rocksCache.CloseLiveFromCache(req.ChannelID); err != nil {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}
	//更新直播计时器
	if err = imdb.UpdateUserLeftDuration(req.UserID, (time.Now().Unix()-liveInfo.StartAt)/60); err != nil {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	if err := s.sendCloseMsg(req.OperationID, fmt.Sprintf("%d", req.UserID), req.ChannelID); err != nil {
		return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, err
	}

	return &pblive.CloseLiveResp{CommonResp: &pblive.CommonResp{}}, err
}

type rpcLive struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewRpcLiveServer(port int) *rpcLive {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcLive{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImLiveName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcLive) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc live start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
	}
	log.NewInfo(operationID, "listen network success, ", address, listener)
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		promePkg.NewUserRegisterCounter()
		promePkg.NewUserLoginCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	//service registers with etcd
	pblive.RegisterLiveServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc live ok")
}

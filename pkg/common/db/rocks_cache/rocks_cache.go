package rocksCache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

const (
	userInfoCache             = "USER_INFO_CACHE:"
	friendRelationCache       = "FRIEND_RELATION_CACHE:"
	blackListCache            = "BLACK_LIST_CACHE:"
	groupCache                = "GROUP_CACHE:"
	liveCache                 = "Live_CACHE:"
	liveMemberCache           = "Live_MEMBER_CACHE:"
	liveRobotCache            = "Live_ROBOT_CACHE:"
	liveAtmosphereCache       = "Live_ATMOSPHERE_CACHE:"
	groupInfoCache            = "GROUP_INFO_CACHE:"
	groupOwnerIDCache         = "GROUP_OWNER_ID:"
	joinedGroupListCache      = "JOINED_GROUP_LIST_CACHE:"
	groupMemberInfoCache      = "GROUP_MEMBER_INFO_CACHE:"
	groupAllMemberInfoCache   = "GROUP_ALL_MEMBER_INFO_CACHE:"
	allFriendInfoCache        = "ALL_FRIEND_INFO_CACHE:"
	allDepartmentCache        = "ALL_DEPARTMENT_CACHE:"
	allDepartmentMemberCache  = "ALL_DEPARTMENT_MEMBER_CACHE:"
	joinedSuperGroupListCache = "JOINED_SUPER_GROUP_LIST_CACHE:"
	groupMemberListHashCache  = "GROUP_MEMBER_LIST_HASH_CACHE:"
	groupMemberNumCache       = "GROUP_MEMBER_NUM_CACHE:"
	conversationCache         = "CONVERSATION_CACHE:"
	conversationIDListCache   = "CONVERSATION_ID_LIST_CACHE:"
	extendMsgSetCache         = "EXTEND_MSG_SET_CACHE:"
	extendMsgCache            = "EXTEND_MSG_CACHE:"
)

func DelKeys() {
	fmt.Println("init to del old keys")
	for _, key := range []string{groupCache, friendRelationCache, blackListCache, userInfoCache, groupInfoCache, groupOwnerIDCache, joinedGroupListCache,
		groupMemberInfoCache, groupAllMemberInfoCache, allFriendInfoCache} {
		fName := utils.GetSelfFuncName()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = db.DB.RDB.Scan(context.Background(), cursor, key+"*", 3000).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = db.DB.RDB.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = db.DB.RDB.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func GetFriendIDListFromCache(userID string) ([]string, error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := imdb.GetFriendIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	friendIDListStr, err := db.DB.Rc.Fetch(friendRelationCache+userID, time.Second*30*60, getFriendIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendIDList []string
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDList)
	return friendIDList, utils.Wrap(err, "")
}

func DelFriendIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(friendRelationCache + userID)
	return err
}

func GetBlackListFromCache(userID string) ([]string, error) {
	getBlackIDList := func() (string, error) {
		blackIDList, err := imdb.GetBlackIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(blackIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	blackIDListStr, err := db.DB.Rc.Fetch(blackListCache+userID, time.Second*30*60, getBlackIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var blackIDList []string
	err = json.Unmarshal([]byte(blackIDListStr), &blackIDList)
	return blackIDList, utils.Wrap(err, "")
}

func DelBlackIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(blackListCache + userID)
}

func GetJoinedGroupIDListFromCache(userID string) ([]string, error) {
	getJoinedGroupIDList := func() (string, error) {
		joinedGroupList, err := imdb.GetJoinedGroupIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(joinedGroupList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	joinedGroupIDListStr, err := db.DB.Rc.Fetch(joinedGroupListCache+userID, time.Second*30*60, getJoinedGroupIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var joinedGroupList []string
	err = json.Unmarshal([]byte(joinedGroupIDListStr), &joinedGroupList)
	return joinedGroupList, utils.Wrap(err, "")
}

func DelJoinedGroupIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(joinedGroupListCache + userID)
}

func DelGroupMemberIDListFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupCache + groupID)
	return err
}

func GetUserInfoFromCache(userID string) (*db.User, error) {
	getUserInfo := func() (string, error) {
		userInfo, err := imdb.GetUserByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(userInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoCache+userID, time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfo := &db.User{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

func DelUserInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(userInfoCache + userID)
}

func GetGroupMemberInfoFromCache(groupID, userID string) (*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMemberInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMemberInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupMemberInfoStr, err := db.DB.Rc.Fetch(groupMemberInfoCache+groupID+"-"+userID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupMember := &db.GroupMember{}
	err = json.Unmarshal([]byte(groupMemberInfoStr), groupMember)
	return groupMember, utils.Wrap(err, "")
}

func DelGroupMemberInfoFromCache(groupID, userID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberInfoCache + groupID + "-" + userID)
}

func GetGroupRobotsRoundFromCache(groupID string) ([]*db.GroupMember, error) {
	groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(groupMemberIDList), func(i, j int) {
		groupMemberIDList[i], groupMemberIDList[j] = groupMemberIDList[j], groupMemberIDList[i]
	})
	roundN := len(groupMemberIDList)
	if roundN > 50 {
		roundN = 50
	}
	count := rand.Intn(roundN)
	var groupMemberList []*db.GroupMember
	for _, userID := range groupMemberIDList {

		if len(groupMemberList) >= count {
			break
		}
		groupMembers, err := GetGroupMemberInfoFromCache(groupID, userID)
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), err.Error(), groupID, userID)
			continue
		}
		if groupMembers.IsRobot == 1 {
			groupMemberList = append(groupMemberList, groupMembers)
		}
	}
	return groupMemberList, err
}

func GetGroupMembersInfoFromCache(count, offset int32, groupID string) ([]*db.GroupMember, error) {
	groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
	if err != nil {
		return nil, err
	}
	if count < 0 || offset < 0 {
		return nil, nil
	}
	var groupMemberList []*db.GroupMember
	var start, stop int32
	start = offset
	stop = offset + count
	l := int32(len(groupMemberIDList))
	if start > stop {
		return nil, nil
	}
	if start >= l {
		return nil, nil
	}
	if count != 0 {
		if stop >= l {
			stop = l
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	} else {
		if l < 1000 {
			stop = l
		} else {
			stop = 1000
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	}
	//log.NewDebug("", utils.GetSelfFuncName(), "ID list: ", groupMemberIDList)
	for _, userID := range groupMemberIDList {
		groupMembers, err := GetGroupMemberInfoFromCache(groupID, userID)
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), err.Error(), groupID, userID)
			continue
		}
		groupMemberList = append(groupMemberList, groupMembers)
	}
	return groupMemberList, nil
}

func GetAllGroupMembersInfoFromCache(groupID string) ([]*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMembers, err := imdb.GetGroupMemberListByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMembers)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupMembersStr, err := db.DB.Rc.Fetch(groupAllMemberInfoCache+groupID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMembers []*db.GroupMember
	err = json.Unmarshal([]byte(groupMembersStr), &groupMembers)
	return groupMembers, utils.Wrap(err, "")
}

func DelAllGroupMembersInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupAllMemberInfoCache + groupID)
}

func GetGroupInfoFromCache(groupID string) (*db.Group, error) {
	getGroupInfo := func() (string, error) {
		groupInfo, err := imdb.GetGroupInfoByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupInfoStr, err := db.DB.Rc.Fetch(groupInfoCache+groupID, time.Second*30*60, getGroupInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupInfo := &db.Group{}
	err = json.Unmarshal([]byte(groupInfoStr), groupInfo)
	return groupInfo, utils.Wrap(err, "")
}

func GetLiveInfoFromCache(channelID string) (*db.UserLive, error) {
	getLiveInfo := func() (string, error) {
		liveInfo, err := imdb.GetLiveByChannelID(channelID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(liveInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupLiveStr, err := db.DB.Rc.Fetch(liveCache+channelID, time.Second*30*60, getLiveInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	liveInfo := &db.UserLive{}
	err = json.Unmarshal([]byte(groupLiveStr), liveInfo)
	return liveInfo, utils.Wrap(err, "")
}

func CheckLiveExits(channelID string) bool {

	count, err := db.DB.RDB.Exists(context.Background(), liveCache+channelID).Result()
	if err != nil {
		log.NewError("checkLiveExits has error", err.Error())
		return false
	}
	return count > 0
}
func GetLiveAtmosphereCache(channelID string) {

}
func GetLiveUsersFromCache(channelID string) (map[string]string, error) {
	userMap, err := db.DB.RDB.HGetAll(context.Background(), liveMemberCache+channelID).Result()
	return userMap, utils.Wrap(err, "")
}

func GetLiveUsersValues(value string) (nickName, faceURL string) {
	data := strings.Split(value, ",")
	return data[0], data[1]
}

func GetLiveUsersLimitFromCache(channelID string, count int64) (map[string]string, error) {
	userMap := make(map[string]string, 0)
	key := liveMemberCache + channelID
	iterator := db.DB.RDB.HScan(context.Background(), key, 0, "*", count).Iterator()
	var values []interface{}
	var err error
	for iterator.Next(context.Background()) {
		userID := iterator.Val()
		values, err = db.DB.RDB.HMGet(context.Background(), key, userID).Result()
		if err != nil {
			return nil, utils.Wrap(err, "")
		}
		userMap[userID] = values[0].(string)
	}
	return userMap, utils.Wrap(err, "")
}

func GetLiveRobotsLimitFromCache(channelID string, count int64) (map[string]string, error) {
	userMap := make(map[string]string, 0)
	key := liveRobotCache + channelID
	iterator := db.DB.RDB.HScan(context.Background(), key, 0, "*", count).Iterator()
	var values []interface{}
	var err error
	for iterator.Next(context.Background()) {
		userID := iterator.Val()
		values, err = db.DB.RDB.HMGet(context.Background(), key, userID).Result()
		if err != nil {
			return nil, utils.Wrap(err, "")
		}
		userMap[userID] = values[0].(string)
	}
	return userMap, utils.Wrap(err, "")
}

func GetLiveRobotsFromCache(channelID string) (map[string]string, error) {
	userMap, err := db.DB.RDB.HGetAll(context.Background(), liveRobotCache+channelID).Result()
	return userMap, utils.Wrap(err, "")
}

func AddLiveAtmosphere(channelID string, userID string) error {
	err := db.DB.RDB.SAdd(context.Background(), liveAtmosphereCache+channelID, userID).Err()
	return utils.Wrap(err, "")
}

// 氛围组
func GetLiveAtmosphereFromCache(channelID string) ([]string, error) {
	userMap, err := db.DB.RDB.SMembers(context.Background(), liveAtmosphereCache+channelID).Result()
	return userMap, utils.Wrap(err, "")
}

func IsLiveAtmosphereUser(channelID string, userID string) (bool, error) {
	has, err := db.DB.RDB.SIsMember(context.Background(), liveAtmosphereCache+channelID, userID).Result()
	return has, utils.Wrap(err, "")
}

func CreateLiveRoom(live db.UserLive) error {
	m := structs.Map(live)
	err := db.DB.RDB.HMSet(context.Background(), liveCache+live.ChannelID, m).Err()
	return utils.Wrap(err, "")
}
func GetLiveRoomFromCache(channelID string) (*db.UserLive, error) {
	var user = &db.UserLive{}
	err := db.DB.RDB.HGetAll(context.Background(), liveCache+channelID).Scan(user)
	if err != nil {
		return nil, err
	}
	return user, utils.Wrap(err, "")
}
func JoinLiveRoom(channelID string, userID string, nickName string, faceURL string, isRoot bool) error {
	prefixKey := liveMemberCache
	if isRoot {
		prefixKey = liveRobotCache
	}
	in, err := UserInRoom(channelID, userID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if in {
		return nil
	}

	err = db.DB.RDB.HMSet(context.Background(), prefixKey+channelID, userID, fmt.Sprintf("%s,%s", nickName, faceURL)).Err()
	if err != nil {
		return utils.Wrap(err, "")
	}

	if err = db.DB.RDB.HIncrBy(context.Background(), liveCache+channelID, "TotalView", 1).Err(); err != nil {
		return utils.Wrap(err, "更新总观看人数失败")
	}
	if err = db.DB.RDB.HIncrBy(context.Background(), liveCache+channelID, "CurrentView", 1).Err(); err != nil {
		return utils.Wrap(err, "更新当前观看人数失败")
	}

	return utils.Wrap(err, "")
}
func UserInRoom(channelID string, userID string) (bool, error) {
	prefixKey := liveMemberCache
	in, err := db.DB.RDB.HExists(context.Background(), prefixKey+channelID, userID).Result()
	if err != nil {
		return false, utils.Wrap(err, "")
	}

	return in, utils.Wrap(err, "")
}
func LevelLiveRoom(channelID string, userID string) error {
	prefixKey := liveMemberCache
	in, err := UserInRoom(channelID, userID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if !in {
		return nil
	}
	err = db.DB.RDB.HDel(context.Background(), prefixKey+channelID, userID).Err()
	if err != nil {
		return utils.Wrap(err, "")
	}

	return utils.Wrap(err, "")
}

func CloseLiveFromCache(channelID string) error {
	liveMemberCacheKey := liveMemberCache + channelID
	liveRobotCacheKey := liveRobotCache + channelID
	liveCacheKey := liveCache + channelID
	return db.DB.RDB.Del(context.Background(), liveMemberCacheKey, liveRobotCacheKey, liveCacheKey).Err()
}

func DelGroupInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupInfoCache + groupID)
}

func GetAllFriendsInfoFromCache(userID string) ([]*db.Friend, error) {
	getAllFriendInfo := func() (string, error) {
		friendInfoList, err := imdb.GetFriendListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendInfoList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allFriendInfoStr, err := db.DB.Rc.Fetch(allFriendInfoCache+userID, time.Second*30*60, getAllFriendInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendInfoList []*db.Friend
	err = json.Unmarshal([]byte(allFriendInfoStr), &friendInfoList)
	return friendInfoList, utils.Wrap(err, "")
}

func DelAllFriendsInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(allFriendInfoCache + userID)
}

func GetAllDepartmentsFromCache() ([]db.Department, error) {
	getAllDepartments := func() (string, error) {
		departmentList, err := imdb.GetSubDepartmentList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allDepartmentsStr, err := db.DB.Rc.Fetch(allDepartmentCache, time.Second*30*60, getAllDepartments)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartments []db.Department
	err = json.Unmarshal([]byte(allDepartmentsStr), &allDepartments)
	return allDepartments, utils.Wrap(err, "")
}

func DelAllDepartmentsFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentCache)
}

func GetAllDepartmentMembersFromCache() ([]db.DepartmentMember, error) {
	getAllDepartmentMembers := func() (string, error) {
		departmentMembers, err := imdb.GetDepartmentMemberList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentMembers)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	allDepartmentMembersStr, err := db.DB.Rc.Fetch(allDepartmentMemberCache, time.Second*30*60, getAllDepartmentMembers)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartmentMembers []db.DepartmentMember
	err = json.Unmarshal([]byte(allDepartmentMembersStr), &allDepartmentMembers)
	return allDepartmentMembers, utils.Wrap(err, "")
}

func DelAllDepartmentMembersFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentMemberCache)
}

func GetJoinedSuperGroupListFromCache(userID string) ([]string, error) {
	getJoinedSuperGroupIDList := func() (string, error) {
		userToSuperGroup, err := db.DB.GetSuperGroupByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(userToSuperGroup.GroupIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	joinedSuperGroupListStr, err := db.DB.Rc.Fetch(joinedSuperGroupListCache+userID, time.Second*30*60, getJoinedSuperGroupIDList)
	if err != nil {
		return nil, err
	}
	var joinedSuperGroupList []string
	err = json.Unmarshal([]byte(joinedSuperGroupListStr), &joinedSuperGroupList)
	return joinedSuperGroupList, utils.Wrap(err, "")
}

func DelJoinedSuperGroupIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(joinedSuperGroupListCache + userID)
	return err
}

func GetGroupMemberListHashFromCache(groupID string) (uint64, error) {
	generateHash := func() (string, error) {
		groupInfo, err := GetGroupInfoFromCache(groupID)
		if err != nil {
			return "0", utils.Wrap(err, "GetGroupInfoFromCache failed")
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return "0", nil
		}
		groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupMemberIDListFromCache failed")
		}
		sort.Strings(groupMemberIDList)
		var all string
		for _, v := range groupMemberIDList {
			all += v
		}
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(all)[0:8], 16)
		return strconv.Itoa(int(bi.Uint64())), nil
	}
	hashCode, err := db.DB.Rc.Fetch(groupMemberListHashCache+groupID, time.Second*30*60, generateHash)
	if err != nil {
		return 0, utils.Wrap(err, "fetch failed")
	}
	hashCodeUint64, err := strconv.Atoi(hashCode)
	return uint64(hashCodeUint64), err
}
func GetGroupMemberIDListFromCache(groupID string) ([]string, error) {
	f := func() (string, error) {
		groupInfo, err := GetGroupInfoFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupInfoFromCache failed")
		}
		var groupMemberIDList []string
		if groupInfo.GroupType == constant.SuperGroup {
			superGroup, err := db.DB.GetSuperGroup(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
			groupMemberIDList = superGroup.MemberIDList
		} else {
			groupMemberIDList, err = imdb.GetGroupMemberIDListByGroupID(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
		}
		bytes, err := json.Marshal(groupMemberIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupIDListStr, err := db.DB.Rc.Fetch(groupCache+groupID, time.Second*30*60, f)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMemberIDList []string
	err = json.Unmarshal([]byte(groupIDListStr), &groupMemberIDList)
	return groupMemberIDList, utils.Wrap(err, "")
}

func DelGroupMemberListHashFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupMemberListHashCache + groupID)
	return err
}

func GetGroupMemberNumFromCache(groupID string) (int64, error) {
	getGroupMemberNum := func() (string, error) {
		num, err := imdb.GetGroupMemberNumByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return strconv.Itoa(int(num)), nil
	}
	groupMember, err := db.DB.Rc.Fetch(groupMemberNumCache+groupID, time.Second*30*60, getGroupMemberNum)
	if err != nil {
		return 0, utils.Wrap(err, "")
	}
	num, err := strconv.Atoi(groupMember)
	return int64(num), err
}

func DelGroupMemberNumFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberNumCache + groupID)
}

func GetUserConversationIDListFromCache(userID string) ([]string, error) {
	getConversationIDList := func() (string, error) {
		conversationIDList, err := imdb.GetConversationIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "getConversationIDList failed")
		}
		log.NewDebug("", utils.GetSelfFuncName(), conversationIDList)
		bytes, err := json.Marshal(conversationIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	conversationIDListStr, err := db.DB.Rc.Fetch(conversationIDListCache+userID, time.Second*30*60, getConversationIDList)
	var conversationIDList []string
	err = json.Unmarshal([]byte(conversationIDListStr), &conversationIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return conversationIDList, nil
}

func DelUserConversationIDListFromCache(userID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationIDListCache+userID), "DelUserConversationIDListFromCache err")
}

func GetConversationFromCache(ownerUserID, conversationID string) (*db.Conversation, error) {
	getConversation := func() (string, error) {
		conversation, err := imdb.GetConversation(ownerUserID, conversationID)
		if err != nil {
			return "", utils.Wrap(err, "get failed")
		}
		bytes, err := json.Marshal(conversation)
		if err != nil {
			return "", utils.Wrap(err, "Marshal failed")
		}
		return string(bytes), nil
	}
	conversationStr, err := db.DB.Rc.Fetch(conversationCache+ownerUserID+":"+conversationID, time.Second*30*60, getConversation)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	conversation := db.Conversation{}
	err = json.Unmarshal([]byte(conversationStr), &conversation)
	if err != nil {
		return nil, utils.Wrap(err, "Unmarshal failed")
	}
	return &conversation, nil
}

func GetConversationsFromCache(ownerUserID string, conversationIDList []string) ([]db.Conversation, error) {
	var conversationList []db.Conversation
	for _, conversationID := range conversationIDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func GetUserAllConversationList(ownerUserID string) ([]db.Conversation, error) {
	IDList, err := GetUserConversationIDListFromCache(ownerUserID)
	if err != nil {
		return nil, err
	}
	var conversationList []db.Conversation
	log.NewDebug("", utils.GetSelfFuncName(), IDList)
	for _, conversationID := range IDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func DelConversationFromCache(ownerUserID, conversationID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationCache+ownerUserID+":"+conversationID), "DelConversationFromCache err")
}

func GetExtendMsg(sourceID string, sessionType int32, clientMsgID string, firstModifyTime int64) (*db.ExtendMsg, error) {
	getExtendMsg := func() (string, error) {
		extendMsg, err := db.DB.GetExtendMsg(sourceID, sessionType, clientMsgID, firstModifyTime)
		if err != nil {
			return "", utils.Wrap(err, "GetExtendMsgList failed")
		}
		bytes, err := json.Marshal(extendMsg)
		if err != nil {
			return "", utils.Wrap(err, "Marshal failed")
		}
		return string(bytes), nil
	}
	extendMsgStr, err := db.DB.Rc.Fetch(extendMsgCache+clientMsgID, time.Second*30*60, getExtendMsg)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	extendMsg := &db.ExtendMsg{}
	err = json.Unmarshal([]byte(extendMsgStr), extendMsg)
	if err != nil {
		return nil, utils.Wrap(err, "Unmarshal failed")
	}
	return extendMsg, nil
}

func DelExtendMsg(ID string, index int32, clientMsgID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(extendMsgCache+clientMsgID), "DelExtendMsg err")
}

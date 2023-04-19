package demo

import (
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/utils"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"
)

var url = "https://api.multiavatar.com/%s.png"

func TestGetRandomName(t *testing.T) {
	for i := 0; i < 1; i++ {
		fmt.Println("i>>>>>>>>>>>>>>>>>>>>>>", i)
		<-time.After(1 * time.Nanosecond)

		name := GetRandomName()

		userID := utils.Md5(name + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(userID[0:8], 16)
		userID = bi.String()
		faceURL := fmt.Sprintf(url, userID)
		err := imdb.UserRegister(db.User{UserID: userID, Nickname: name, FaceURL: faceURL, Gender: 2, AppMangerLevel: 1, CreateTime: time.Now(), IsRobot: 1})
		if err != nil {
			panic(err)
		}
		// if err = imdb.InsertIntoGroupMember(db.GroupMember{GroupID: "1763879570", UserID: userID, Nickname: name, FaceURL: faceURL, RoleLevel: 1, JoinTime: time.Now(), JoinSource: 2, InviterUserID: "2892606321", OperatorUserID: "2892606321"}); err != nil {
		// 	panic(err)
		// }
	}

}

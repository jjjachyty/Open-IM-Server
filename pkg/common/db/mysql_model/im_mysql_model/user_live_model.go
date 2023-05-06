package im_mysql_model

import (
	"Open_IM/pkg/common/db"
)

func init() {

}

func GetLiveByUserID(userID string) (*db.UserLive, error) {
	var live db.UserLive
	err := db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("user_id=?", userID).Find(&live).Error
	if err != nil {
		return nil, err
	}
	return &live, nil
}
func UpdateLiveInfo(user db.UserLive) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("user_id=?", user.UserID).Updates(&user).Error
}

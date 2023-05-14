package im_mysql_model

import (
	"Open_IM/pkg/common/db"
)

func init() {

}

func GetUserLiving(userID string) (*db.UserLive, error) {
	var live db.UserLive
	err := db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("user_id=? and end_at is null", userID).Take(&live).Error
	if err != nil {
		return nil, err
	}
	return &live, nil
}

func GetLiveByUserID(userID string) (*db.UserLive, error) {
	var live db.UserLive
	err := db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("user_id=?", userID).Find(&live).Error
	if err != nil {
		return nil, err
	}
	return &live, nil
}

func GetLiveByChannelID(channelID string) (*db.UserLive, error) {
	var live db.UserLive
	err := db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("channel_id=?", channelID).Find(&live).Error
	if err != nil {
		return nil, err
	}
	return &live, nil
}

func CreateLiveInfo(live *db.UserLive) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Create(live).Error
}

func UpdateLiveInfo(live db.UserLive) error {
	return db.DB.MysqlDB.DefaultGormDB().Model(db.UserLive{}).Where("channel_id=?", live.ChannelID).Updates(&live).Error
}

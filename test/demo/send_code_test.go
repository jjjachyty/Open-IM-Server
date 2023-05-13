package demo

import (
	rocksCache "Open_IM/pkg/common/db/rocks_cache"

	"os"
	"testing"
)

func init() {
	os.Setenv("CONFIG_NAME", "/Users/janly/data/go/src/Open-IM-Server/")
	os.Setenv("USUAL_CONFIG_NAME", "/Users/janly/data/go/src/Open-IM-Server/")
}

func TestSendVerificationCode(t *testing.T) {

	liveInfo, err := rocksCache.GetLiveRoomFromCache(3817720326)
	if err != nil {
		panic(err)
	}
	println(liveInfo.UserID)
}

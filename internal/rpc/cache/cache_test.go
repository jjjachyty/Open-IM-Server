package cache

import (
	pbCache "Open_IM/pkg/proto/cache"
	"context"
	"fmt"

	"testing"
)

func Test_cacheServer_GetLiveMemberIDListFromCache(t *testing.T) {
	cache := NewCacheServer(8080)

	gotResp, err := cache.GetLiveMemberIDListFromCache(context.Background(), &pbCache.GetLiveMemberIDListFromCacheReq{ChannelID: "3817720326"})
	fmt.Println(gotResp, err)

}

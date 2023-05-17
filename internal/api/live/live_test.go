package live

import (
	api "Open_IM/pkg/base_info"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	// 添加 Get 请求路由
	router.POST("/", func(context *gin.Context) {
		JoinLiveRoom(context)
	})
	return router
}

func TestJoinLiveRoom(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	params := api.JoinRoomReq{ChannelID: "3817720326", UserID: "1910360909"}
	paramsData, _ := json.Marshal(params)
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(paramsData))
	router.ServeHTTP(w, req)
	fmt.Println(http.StatusOK, w.Code)
	fmt.Println("hello gin", w.Body.String())
}

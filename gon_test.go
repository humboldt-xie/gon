package gon

import (
	"context"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type myHttp struct {
}

type param struct {
	Param string
}
type result struct {
	Result string
}

func (m *myHttp) Request2(ctx context.Context, h *param) result {
	return result{Result: "201"}
}

func (m *myHttp) Request(h *param) result {
	return result{Result: "200"}
}

func post(engine *gin.Engine, path string, postbody string) (*http.Response, string) {
	//构建返回值
	w := httptest.NewRecorder()
	//构建请求
	r, _ := http.NewRequest("POST", path, strings.NewReader(postbody))
	//调用请求接口
	engine.ServeHTTP(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	return resp, string(body)
}

func newEngine(in interface{}) *gin.Engine {
	gin.SetMode(gin.TestMode)
	eg := gin.New()
	RegisterInterface("path", eg, in)
	return eg
}

func TestMap(t *testing.T) {
	eg := newEngine(&myHttp{})

	resp, body := post(eg, "/path/request2", `{"Param":""}`)

	t.Log(resp.StatusCode)
	//获得结果，并检查
	if string(body) != `{"Result":"201"}` {
		t.Fatal(body)
	}

}

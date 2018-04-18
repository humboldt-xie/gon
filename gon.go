package gon

import (
	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(m reflect.Value) string {
	return runtime.FuncForPC(m.Pointer()).Name()
}

func JsonMap(c interface{}) {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Type().Method(i)
		log.Printf("method:%s\n", method.Name)
		numin := v.Method(i).Type().NumIn()
		in := v.Method(i).Type().In(numin - 1)
		inPut := reflect.New(in.Elem())
		log.Printf("new input:%#v\n", inPut)
	}
}

// input   (context.Context,request) or (request)
// return  Object or (Object,error)

func GinRPC(r gin.IRouter, v reflect.Value, method reflect.Method) {
	log.Printf("gin rpc method:%s", method.Name)
	handler := func(c *gin.Context) {
		num := v.Type().NumIn()
		inValue := make([]reflect.Value, 0)
		if num > 0 {
			in := v.Type().In(num - 1)
			inPut := reflect.New(in.Elem())
			err := c.BindJSON(inPut.Interface())
			if err != nil {
				c.String(401, fmt.Sprintf("%s", err))
				return
			}
			if num > 1 {
				ctx := reflect.ValueOf(context.Background())
				inValue = append(inValue, ctx)
				inValue = append(inValue, inPut)
			} else {
				inValue = append(inValue, inPut)
			}
		}
		//log.Printf("bind input:%#v %s\n", inPut, err)
		response := v.Call(inValue)
		//log.Printf("output :%#v \n", response)
		if len(response) > 1 {
			if response[1].Interface() != nil {
				c.String(500, fmt.Sprintf("%s", response[1].Interface()))
				return
			}
		}
		c.JSON(200, response[0].Interface())
	}
	lowerName := strings.ToLower(method.Name)
	if lowerName != method.Name {
		r.POST(lowerName, handler)
	}
	r.POST(method.Name, handler)
}

func RegisterInterface(name string, r gin.IRouter, c interface{}) {
	group := r.Group("/" + name)
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Type().Method(i)
		GinRPC(group, v.Method(i), method)
	}

}

func NewHttp(name string, i interface{}, listen string) {
	eg := gin.New()
	JsonMap(i)
	RegisterInterface(name, eg, i)
	eg.Run(listen)
}

// SimpleMVC project main.go
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"

	_ "github.com/qqcowboy/dunhuaparts.com/Controllers"
	. "github.com/qqcowboy/dunhuaparts.com/System/Routing"
	. "github.com/qqcowboy/dunhuaparts.com/System/Web"
)

//当前配置文件的端口为6080,输入http://localhost:6080/可查看运行结果
//注册路由
func init() {
	RouteTable.AddRote(&RouteItem{
		Name:     "default",
		Url:      "{controller}/{action}",
		Defaults: map[string]interface{}{"controller": "home", "action": "index", "views": ""},
	})
}
func main() {
	//程序意外退时，记录错误日志
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			App.Log.Add(err.Error() + "\r\n" + string(debug.Stack()))
			fmt.Println(err)
		}
	}()
	//设置最大可同时执行的进程数
	runtime.GOMAXPROCS(runtime.NumCPU()*2 - 1)
	//监听http请求
	err := App.Run()
	fmt.Println(err)
}

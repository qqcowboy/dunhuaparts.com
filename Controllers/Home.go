package Controllers

import (
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
)

type Home struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Home{})
}
func (this *Home) OnLoad() {
	//如果在OnLoad里调用了Response.Write需选设置Content-Type头，设置文档类型
	//this.Response.Header().Add("Content-Type", "text/html;charset=utf-8")
	//this.Response.Write([]byte("在 OnLoad函数里面"))
}

func (this *Home) Index() *Web.ViewResult {
	return this.View()
}

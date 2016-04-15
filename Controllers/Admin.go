package Controllers

import (
	_ "github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
)

type Admin struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Admin{})
}
func (this *Admin) OnLoad() {
	//如果没有登陆则跳转到后台主页，由主页显示登陆页面
	if !this.isLogin() && this.RouteData["action"] != "Index" {
		this.Redirect("/admin")
		this.ResponseEnd()
	}
}
func (this *Admin) isLogin() bool {
	if tmp, ok := this.Session["login"]; ok {
		if login, ok := tmp.(bool); ok && login {
			return true
		}
	}
	return false
}

/*Index
@see : 后台页面
*/
func (this *Admin) Index() *Web.ViewResult {
	if !this.isLogin() {
		this.RouteData["action"] = "Login"
	}
	return this.View()
}

/*Loginout
@see : 退出后台
*/
func (this *Admin) Logout() *Web.ViewResult {
	this.RouteData["action"] = "Login"
	this.Session["login"] = false
	return this.View()
}

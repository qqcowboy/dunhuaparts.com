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

/*Login
@see : 登陆后台
*/
func (this *Admin) Login() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	this.Session["login"] = true
	return this.Json(map[string]interface{}{"code": 1})
}

/*Loginout
@see : 退出后台
*/
func (this *Admin) Logout() *Web.ViewResult {
	this.RouteData["action"] = "Login"
	this.Session["login"] = false
	return this.View()
}

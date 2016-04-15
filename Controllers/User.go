package Controllers

import (
	"strconv"

	"github.com/dchest/captcha"
	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
	"github.com/qqcowboy/lib/myjson"
)

type User struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(User{})
}
func (this *User) OnLoad() {
}

/*AdminLogin
@see : 登陆后台
*/
func (this *User) AdminLogin() *Web.JsonResult {
	if !this.IsPost {
		return this.JsonResult(43002, nil, "需要POST请求")
	}
	user, _ := this.Form["User"]
	psw, _ := this.Form["Psw"]
	code, _ := this.Form["Code"] //验证码
	if len(user) < 1 || len(psw) < 1 {
		return this.JsonResult(43007, nil, "请输入用户名和密码")
	}
	if len(code) < 1 {
		return this.JsonResult(43008, nil, "请输入验证码")
	}

	if !captcha.VerifyString(this.Cookies[VerityCodeCookieName], code) {
		return this.JsonResult(43008, nil, "验证码不正确")
	}
	var userinfo Model.UserInfo
	var err error
	if unum, err1 := strconv.Atoi(user); err1 == nil {
		userinfo, err = Model.UserModel.FindOne("", psw, unum)
	} else {
		userinfo, err = Model.UserModel.FindOne(user, psw, 0)
	}
	if err != nil {
		return this.JsonResult(43006, nil, "用户名或密码不正确")
	}
	if !userinfo.IsAdmin() {
		return this.JsonResult(43009, nil, "该用户无权限")
	}
	this.Session["login"] = true
	this.Session["user"] = myjson.JsonEncode(userinfo)
	return this.JsonResult(1, userinfo)
}

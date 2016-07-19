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
func (this *User) isLogin() bool {
	if tmp, ok := this.Session["login"]; ok {
		if login, ok := tmp.(bool); ok && login {
			return true
		}
	}
	return false
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

/*Info
@see : 当前用户信息
*/
func (this *User) Info() *Web.JsonResult {
	if !this.isLogin() {
		return this.JsonResult(40001, nil, "用户未登陆")
	}
	udata, ok := this.Session["user"]
	if !ok {
		return this.JsonResult(40002, nil, "用户登陆信息出错")
	}
	userinfo := Model.UserInfo{}
	myjson.JsonDecode(udata.(string), &userinfo)

	return this.JsonResult(1, userinfo)
}

/*ChPsw
@see : 更新密码
@param: OP 旧密码
@param: NP 新密码
*/
func (this *User) ChPsw() *Web.JsonResult {
	if !this.IsPost {
		return this.JsonResult(43002, false, "需要POST请求")
	}
	if !this.isLogin() {
		return this.JsonResult(40001, false, "用户未登陆")
	}
	if _, ok := this.Form["OP"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	if _, ok := this.Form["NP"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	udata, ok := this.Session["user"]
	if !ok {
		return this.JsonResult(40002, false, "用户登陆信息出错")
	}
	userinfo := Model.UserInfo{}
	myjson.JsonDecode(udata.(string), &userinfo)
	err := Model.UserModel.ChPsw(userinfo.UID, this.Form["OP"], this.Form["NP"])
	if err != nil {
		return this.JsonResult(43010, false, err.Error())
	}
	return this.JsonResult(1, true)
}

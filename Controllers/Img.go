package Controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
)

type Img struct {
	Web.Controller
}

var VerityCodeCookieName = "vercode"

//注册Controller
func init() {
	Web.App.RegisterController(Img{})
}
func (this *Img) OnLoad() {
}

//输出产品图片
func (this *Img) Product() *Web.ImgResult {
	result := &Web.ImgResult{Response: this.Response}
	if tmp, ok := this.RouteData["views"]; ok {
		if str, ok := tmp.(string); ok && len(str) > 0 {
			views := strings.Split(str, "/")
			i, err := strconv.Atoi(views[0])
			if err == nil {
				result.MaxSize = i * 1024
			}
		}
	}
	params := this.QueryString
	if this.IsPost {
		params = this.Form
	}
	id, _ := params["id"]
	if len(id) < 1 {
		return result
	}
	str, _ := Model.MProducts.GetImage(id)
	result.Base64 = str
	//result.Bytes, _ = Model.AddWaterMark(str)
	return result
}

//输出验证码
func (this *Img) VerityCode() *Web.ImgResult {
	result := &Web.ImgResult{Response: this.Response}
	vkey := captcha.NewLen(5)
	http.SetCookie(this.Response, &http.Cookie{Name: VerityCodeCookieName, Value: vkey, Path: "/", HttpOnly: true})
	captcha.WriteImage(this.Response, vkey, 120, 50)
	return result
}

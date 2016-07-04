package Controllers

import (
	"strings"

	"github.com/dchest/captcha"
	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
	"github.com/qqcowboy/lib/mail"
	"github.com/qqcowboy/lib/myjson"
	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2/bson"
)

type Feedback struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Feedback{})
}
func (this *Feedback) OnLoad() {
}

func (this *Feedback) hasAuth() bool {
	if tmp, ok := this.Session["login"]; ok {
		if login, ok := tmp.(bool); ok && login {
			return true
		}
	}
	return false
}

/*AQuery
@see 管理员查询[post]
@param data : json {Start:float64,Limit:int,Key:string,Mail:string,Hide:0/1/2}
*/
func (this *Feedback) AQuery() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	start := -0.1
	limit := 20
	key := ""
	email := ""
	hide := 0
	if _, ok := this.Form["data"]; ok {
		params := make(map[string]interface{})
		myjson.JsonDecode(this.Form["data"], &params)
		if _, ok := params["Start"]; ok && mystr.IsFloat64(params["Start"]) {
			start = params["Start"].(float64)
		}
		if tmp, ok := params["Limit"]; ok {
			tmp1, err := mystr.ToInt(tmp)
			if err == nil {
				limit = tmp1
			}
		}
		if tmp, ok := params["Hide"]; ok {
			tmp1, err := mystr.ToInt(tmp)
			if err == nil {
				hide = tmp1
			}
		}
		if tmp, ok := params["Mail"]; ok && mystr.IsString(tmp) {
			email = tmp.(string)
		}
		if tmp, ok := params["Key"]; ok && mystr.IsString(tmp) {
			key = tmp.(string)
		}
	}
	exttype := []int{}
	sel := map[string]interface{}{"Content": 0, "Reply": 0}
	count, lists, err := Model.MFeedback.QueryFeedback(start, limit, []string{"-Version"}, hide, key, email, sel, exttype...)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": lists, "count": count})
}

/*Query
@see 查询[post]
@param data : json {Start:float64,Limit:int,Key:string,Mail:string}
*/
func (this *Feedback) Query() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	start := -0.1
	limit := 20
	key := ""
	email := ""
	if _, ok := this.Form["data"]; ok {
		params := make(map[string]interface{})
		myjson.JsonDecode(this.Form["data"], &params)
		if _, ok := params["Start"]; ok && mystr.IsFloat64(params["Start"]) {
			start = params["Start"].(float64)
		}
		if tmp, ok := params["Limit"]; ok {
			tmp1, err := mystr.ToInt(tmp)
			if err == nil {
				limit = tmp1
			}
		}
		if tmp, ok := params["Mail"]; ok && mystr.IsString(tmp) {
			email = tmp.(string)
		}
		if tmp, ok := params["Key"]; ok && mystr.IsString(tmp) {
			key = tmp.(string)
		}
	}
	hide := 2
	if len(email) > 0 {
		hide = 0
	}
	exttype := []int{}
	sel := map[string]interface{}{"Content": 0, "Reply": 0, "Mail": 0, "Phone": 0}
	count, lists, err := Model.MFeedback.QueryFeedback(start, limit, []string{"-Version"}, hide, key, email, sel, exttype...)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": lists, "count": count})
}

/*Create
@see 新增[post]
@param data : json {Title,Content,Mail,UserName,Phone,Hide},Code
*/
func (this *Feedback) Create() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	code, _ := this.Form["Code"] //验证码
	if len(code) < 1 {
		return this.JsonResult(43008, nil, "请输入验证码")
	}

	if !captcha.VerifyString(this.Cookies[VerityCodeCookieName], code) {
		return this.JsonResult(43008, nil, "验证码不正确")
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)

	title := ""
	content := ""
	hide := 0
	user := ""
	phone := ""
	email := ""
	ip := this.Request.RemoteAddr
	if title, ok := params["Title"]; !ok || !mystr.IsString(title) || len(strings.TrimSpace(title.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入标题"})
	}
	title = strings.TrimSpace(params["Title"].(string))

	if tmp, ok := params["Content"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.MAXCONTENTSIZE {
				tmp1 = string([]byte(tmp1)[0:Model.MAXCONTENTSIZE])
			}
			content = tmp1
		}
	}
	if len(content) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入内容"})
	}
	if tmp, ok := params["UserName"]; ok {
		if tmp1, ok := tmp.(string); ok {
			user = tmp1
		}
	}
	if len(user) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入用户名"})
	}
	if tmp, ok := params["Mail"]; ok {
		if tmp1, ok := tmp.(string); ok {
			email = tmp1
		}
	}
	email, err := mail.ParseMailAddr(email)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "邮件地址格式不正确"})
	}
	if tmp, ok := params["Phone"]; ok {
		if tmp1, ok := tmp.(string); ok {
			phone = tmp1
		}
	}
	if tmp, ok := params["Hide"]; ok {
		tmp1, err := mystr.ToInt(tmp)
		if err == nil {
			hide = tmp1
		}
	}
	Feedback, err := Model.MFeedback.CreateFeedback(title, content, 0, hide, user, phone, email, ip)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": Feedback})
}

/*Detail
@see 详情[get]
@param data : json {ID:string}
*/
func (this *Feedback) Detail() *Web.JsonResult {
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43001})
	}
	if _, ok := this.QueryString["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.QueryString["data"], &params)
	if idstr, ok := params["ID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "参数不正确"})
	}
	Feedback, err := Model.MFeedback.Get(strings.TrimSpace(params["ID"].(string)))
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": Feedback})
}

/*Remove
@see 删除[get]
@param data : json {IDs:string}
*/
func (this *Feedback) Remove() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43001})
	}
	if _, ok := this.QueryString["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.QueryString["data"], &params)
	if idstr, ok := params["IDs"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "没有需要删除的留言"})
	}
	idstr := strings.TrimSpace(params["IDs"].(string))
	ids := strings.Split(idstr, ",")
	idsi := []interface{}{}
	for _, v := range ids {
		idsi = append(idsi, bson.ObjectIdHex(v))
	}
	err := Model.MFeedback.Removes(idsi...)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*Update
@see 更新[post]
@param data : json {ID,Data:map[string]interface{}}
*/
func (this *Feedback) Update() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if idstr, ok := params["ID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "ID是必须参数"})
	}
	if _, ok := params["Data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "Data是必须参数"})
	}
	param, ok := params["Data"].(map[string]interface{})
	if !ok || len(param) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "Data格式不正确"})
	}
	ID := bson.ObjectIdHex(strings.TrimSpace(params["ID"].(string)))
	err := Model.MFeedback.UpdateId(ID, param)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*AddReply
@see 新增回复[post]
@param data : json {ID,UserName,Content}
*/
func (this *Feedback) AddReply() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if idstr, ok := params["ID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "ID是必须参数"})
	}
	content := ""
	if tmp, ok := params["Content"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.MAXCONTENTSIZE {
				tmp1 = string([]byte(tmp1)[0:Model.MAXCONTENTSIZE])
			}
			content = tmp1
		}
	}
	if len(content) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入内容"})
	}
	user := ""
	if tmp, ok := params["UserName"]; ok {
		if tmp1, ok := tmp.(string); ok {
			user = tmp1
		}
	}
	if len(user) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入用户名"})
	}
	ID := strings.TrimSpace(params["ID"].(string))
	r, err := Model.MFeedback.AddReply(ID, user, content)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": r})
}

/*DelReply
@see 删除回复[get]
@param data : json {ID,RID}
*/
func (this *Feedback) DelReply() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43001})
	}
	if _, ok := this.QueryString["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.QueryString["data"], &params)
	if idstr, ok := params["ID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "ID是必须参数"})
	}
	if idstr, ok := params["RID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "RID是必须参数"})
	}
	ID := strings.TrimSpace(params["ID"].(string))
	RID := strings.TrimSpace(params["RID"].(string))
	err := Model.MFeedback.DelReply(ID, RID)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*UpdateReply
@see 更新回复[post]
@param data : json {ID,RID,UserName,Content}
*/
func (this *Feedback) UpdateReply() *Web.JsonResult {
	if this.hasAuth() == false {
		return this.Json(map[string]interface{}{"code": 40003, "msg": "无权限操作"})
	}
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if idstr, ok := params["ID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "ID是必须参数"})
	}
	if idstr, ok := params["RID"]; !ok || !mystr.IsString(idstr) || len(strings.TrimSpace(idstr.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "RID是必须参数"})
	}
	content := ""
	if tmp, ok := params["Content"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.MAXCONTENTSIZE {
				tmp1 = string([]byte(tmp1)[0:Model.MAXCONTENTSIZE])
			}
			content = tmp1
		}
	}
	user := ""
	if tmp, ok := params["UserName"]; ok {
		if tmp1, ok := tmp.(string); ok {
			user = tmp1
		}
	}
	ID := strings.TrimSpace(params["ID"].(string))
	RID := strings.TrimSpace(params["RID"].(string))
	err := Model.MFeedback.UpdateReply(ID, RID, user, content)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

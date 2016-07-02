package Controllers

import (
	"strings"

	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
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

/*Query
@see 查询[post]
@param data : json {Start:float64,Limit:int,Key:string,Top:int,Lang:all/cn/en}
*/
func (this *Feedback) Query() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	start := -0.1
	limit := 20
	top := 0
	lang := "all"
	key := ""
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
		if tmp, ok := params["Top"]; ok {
			tmp1, err := mystr.ToInt(tmp)
			if err == nil {
				top = tmp1
			}
		}
		if tmp, ok := params["Lang"]; ok && mystr.IsString(tmp) {
			switch tmp.(string) {
			case "cn":
				lang = "cn"
			case "en":
				lang = "en"
			default:
				lang = "all"
			}
		}
		if tmp, ok := params["Key"]; ok && mystr.IsString(tmp) {
			key = tmp.(string)
		}
	}
	exttype := []int{}
	switch lang {
	case "cn":
		exttype = []int{0}
	case "en":
		exttype = []int{1}
	default:
		exttype = []int{0, 1}
	}
	count, lists, err := Model.MFeedback.QueryFeedback(start, limit, []string{"-Top", "-Version"}, top, key, exttype...)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": lists, "count": count})
}

/*Create
@see 新增新闻[post]
@param data : json {Title,ExtType:0普通,Content,Lead,Top,Author}
*/
func (this *Feedback) Create() *Web.JsonResult {
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
	if title, ok := params["Title"]; !ok || !mystr.IsString(title) || len(strings.TrimSpace(title.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入标题"})
	}
	title := strings.TrimSpace(params["Title"].(string))
	exttype := 0
	if tmp, ok := params["ExtType"]; ok {
		tmp1, err := mystr.ToInt(tmp)
		if err == nil {
			exttype = tmp1
		}
	}
	top := 0
	if tmp, ok := params["Top"]; ok {
		tmp1, err := mystr.ToInt(tmp)
		if err == nil {
			top = tmp1
		}
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
	lead := ""
	if tmp, ok := params["Lead"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.MAXLEADSIZE {
				tmp1 = string([]byte(tmp1)[0:Model.MAXLEADSIZE])
			}
			lead = tmp1
		}
	}
	author := ""
	if tmp, ok := params["Author"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.MAXLEADSIZE {
				tmp1 = string([]byte(tmp1)[0:Model.MAXLEADSIZE])
			}
			author = tmp1
		}
	}
	Feedback, err := Model.MFeedback.CreateFeedback(title, lead, content, exttype, top, author)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": Feedback})
}

/*Detail
@see 新闻详情[get]
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
@see 删除新闻[get]
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
		return this.Json(map[string]interface{}{"code": 43005, "msg": "没有需要删除的新闻"})
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
@see 更新新闻[post]
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

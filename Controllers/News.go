package Controllers

import (
	"strings"

	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
	"github.com/qqcowboy/lib/myjson"
	"github.com/qqcowboy/lib/mystr"
)

type News struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(News{})
}
func (this *News) OnLoad() {
}

/*Create
@see 新增新闻[post]
@param data : json {Title,ExtType:0普通,Content,Lead}
*/
func (this *News) Create() *Web.JsonResult {
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
	news, err := Model.MNews.CreateNews(title, lead, content, exttype)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": news})
}

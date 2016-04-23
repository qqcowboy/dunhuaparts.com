package Controllers

import (
	"fmt"
	"strings"

	"github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
	"github.com/qqcowboy/lib/myjson"
	"github.com/qqcowboy/lib/mystr"
)

type Products struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Products{})
}
func (this *Products) OnLoad() {
}

/*Create
@see 新增产品[post]
@param data : json {Title,TitleCN,CategoryID,Remark,RemarkCN,Images,ExtType:0普通/1首页大图}
*/
func (this *Products) Create() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if title, ok := params["Title"]; !ok || !mystr.IsString(title) || len(strings.TrimSpace(title.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入产品名称"})
	}
	title := strings.TrimSpace(params["Title"].(string))
	if titlecn, ok := params["TitleCN"]; !ok || !mystr.IsString(titlecn) || len(strings.TrimSpace(titlecn.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入产品中文名称"})
	}
	titlecn := strings.TrimSpace(params["TitleCN"].(string))
	if indexid, ok := params["CategoryID"]; !ok || !mystr.IsString(indexid) || len(strings.TrimSpace(indexid.(string))) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请选择产品分类"})
	}
	indexid := strings.TrimSpace(params["CategoryID"].(string))
	if _, ok := params["Images"]; !ok {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请上传产品图片"})
	}
	images, err := mystr.ToStrings(params["Images"])
	if err != nil || len(images) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请上传产品图片"})
	}
	exttype := 0
	if tmp, ok := params["ExtType"]; ok {
		tmp1, err := mystr.ToInt(tmp)
		if err == nil {
			exttype = tmp1
		}
	}
	remark := ""
	if tmp, ok := params["Remark"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			remark = tmp1
		}
	}
	remarkcn := ""
	if tmp, ok := params["RemarkCN"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			remarkcn = tmp1
		}
	}
	if len(title) > Model.TitleMaxLen {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "产品名称过长"})
	}
	product, err := Model.MProducts.CreateProduct(title, titlecn, indexid, remark, remarkcn, images, exttype)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": product})
}

/*Update
@see 更新产品[post]
@param data : json {ID,Title,TitleCN,CategoryID,Remark,RemarkCN,ExtType,Images:[]string}
*/
func (this *Products) Update() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if id, ok := params["ID"]; !ok || !mystr.IsString(id) || len(id.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "要更新的ID不存在"})
	}
	id := params["ID"].(string)

	fields := map[string]interface{}{}
	if title, ok := params["Title"]; ok && mystr.IsString(title) && len(strings.TrimSpace(title.(string))) > 0 {
		fields["Title"] = strings.TrimSpace(title.(string))
		if len(strings.TrimSpace(title.(string))) > Model.TitleMaxLen {
			return this.Json(map[string]interface{}{"code": 43005, "msg": "产品名称过长"})
		}
	}
	if titlecn, ok := params["TitleCN"]; ok && mystr.IsString(titlecn) && len(strings.TrimSpace(titlecn.(string))) > 0 {
		fields["TitleCN"] = strings.TrimSpace(titlecn.(string))
		if len(strings.TrimSpace(titlecn.(string))) > Model.TitleMaxLen {
			return this.Json(map[string]interface{}{"code": 43005, "msg": "产品中文名称过长"})
		}
	}
	if indexid, ok := params["CategoryID"]; ok && mystr.IsString(indexid) && len(indexid.(string)) > 0 {
		fields["IndexID"] = indexid
	}
	if images, ok := params["Images"]; ok && mystr.IsStrings(images) {
		fields["Images"] = images
	}

	if tmp, ok := params["ExtType"]; ok {
		tmp1, err := mystr.ToInt(tmp)
		if err == nil {
			fields["ExtType"] = tmp1
		}
	}
	if tmp, ok := params["Remark"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			fields["Remark"] = tmp1
		}
	}
	if tmp, ok := params["RemarkCN"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			fields["RemarkCN"] = tmp1
		}
	}
	err := Model.MProducts.UpdateId(id, fields)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*UpdateCategory
@see 更新产品分类[post]
@param data : json {ID,Name,Remark,NameCN,RemarkCN}
*/
func (this *Products) UpdateCategory() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if id, ok := params["ID"]; !ok || !mystr.IsString(id) || len(id.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "要更新的ID不存在"})
	}
	id := params["ID"].(string)

	fields := map[string]interface{}{}
	if name, ok := params["Name"]; ok && mystr.IsString(name) && len(strings.TrimSpace(name.(string))) > 0 {
		fields["Name"] = strings.TrimSpace(name.(string))
		if len(strings.TrimSpace(name.(string))) > Model.TitleMaxLen {
			return this.Json(map[string]interface{}{"code": 43005, "msg": "分类名称过长"})
		}
	}
	if name, ok := params["NameCN"]; ok && mystr.IsString(name) && len(strings.TrimSpace(name.(string))) > 0 {
		fields["NameCN"] = strings.TrimSpace(name.(string))
		if len(strings.TrimSpace(name.(string))) > Model.TitleMaxLen {
			return this.Json(map[string]interface{}{"code": 43005, "msg": "分类中文名称过长"})
		}
	}
	if tmp, ok := params["Remark"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			fields["Remark"] = tmp1
		}
	}
	if tmp, ok := params["RemarkCN"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			fields["RemarkCN"] = tmp1
		}
	}
	err := Model.MProducts.UpdateId(id, fields)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*AddProductImg
@see 新增产品图片[post]
@param data : json {ProductID,Image:string}
*/
func (this *Products) AddProductImg() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if id, ok := params["ProductID"]; !ok || !mystr.IsString(id) || len(id.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "参数不完整ProductID"})
	}
	id := params["ProductID"].(string)

	if image, ok := params["Image"]; !ok || !mystr.IsString(image) || len(image.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "参数不完整Image"})
	}
	image := params["Image"].(string)
	err := Model.MProducts.AddImg(id, image)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": true})
}

/*Remove
@see 删除产品[get]
@param data : json {ID:[string]}
*/
func (this *Products) Remove() *Web.JsonResult {
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.QueryString["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.QueryString["data"], &params)

	if _, ok := params["ID"]; !ok {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "选择要删除的产品"})
	}
	ids, err := mystr.ToStrings(params["ID"])
	if err != nil || len(ids) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "选择要删除的产品"})
	}
	err = Model.MProducts.RemoveProducts(ids)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1})
}

/*Query
@see 查询普通产品[post]
@param data : json {Start:float64,Limit:int,CategoryID:string}
*/
func (this *Products) Query() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	start := -0.1
	limit := 20
	indexid := ""
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
		if tmp, ok := params["CategoryID"]; ok && mystr.IsString(tmp) && len(tmp.(string)) > 0 {
			indexid = tmp.(string)
		}
	}
	count, lists, err := Model.MProducts.QueryProducts(start, limit, indexid, []string{"-Version"})
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": lists, "count": count})
}

/*RemoveCategory
@see 删除产品分类[get]
@param data : json {ID:string}
*/
func (this *Products) RemoveCategory() *Web.JsonResult {
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.QueryString["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]string)
	myjson.JsonDecode(this.QueryString["data"], &params)

	if _, ok := params["ID"]; !ok {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "选择要删除的产品分类"})
	}
	id := params["ID"]
	if len(id) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "选择要删除的产品分类"})
	}
	err := Model.MProducts.RemoveCategory(id)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1})
}

/*Categories
@see 查询产品分类[get]
@param data : json {Start:float64,Limit:int}
*/
func (this *Products) Categories() *Web.JsonResult {
	if this.IsPost {
		return this.Json(map[string]interface{}{"code": 43001})
	}
	start := -0.1
	limit := 20
	if _, ok := this.QueryString["data"]; ok {
		params := make(map[string]interface{})
		myjson.JsonDecode(this.QueryString["data"], &params)
		fmt.Println(params, params["Start"], mystr.IsFloat64(params["Start"]))
		if _, ok := params["Start"]; ok && mystr.IsFloat64(params["Start"]) {
			start = params["Start"].(float64)
		}
		if tmp, ok := params["Limit"]; ok {
			tmp1, err := mystr.ToInt(tmp)
			if err == nil {
				limit = tmp1
			}
		}
	}
	count, lists, err := Model.MProducts.QueryCategory(start, limit)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": lists, "count": count})
}

/*CreateCategory
@see 新增新产品分类，暂不支持多级分类[post]
@params Remark,Name,RemarkCN,NameCN
*/
func (this *Products) CreateCategory() *Web.JsonResult {
	if !this.IsPost {
		return this.Json(map[string]interface{}{"code": 43002})
	}
	if _, ok := this.Form["data"]; !ok {
		return this.Json(map[string]interface{}{"code": 43004})
	}
	params := make(map[string]interface{})
	myjson.JsonDecode(this.Form["data"], &params)
	if name, ok := params["Name"]; !ok || !mystr.IsString(name) || len(name.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入产品分类名称"})
	}
	name := params["Name"].(string)
	if len(name) > Model.TitleMaxLen {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "产品类名过长"})
	}
	if namecn, ok := params["NameCN"]; !ok || !mystr.IsString(namecn) || len(namecn.(string)) < 1 {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "请输入产品分类中文名称"})
	}
	namecn := params["NameCN"].(string)

	if len(namecn) > Model.TitleMaxLen {
		return this.Json(map[string]interface{}{"code": 43005, "msg": "产品分类中文名过长"})
	}
	remark := ""
	if tmp, ok := params["Remark"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			remark = tmp1
		}
	}
	remarkcn := ""
	if tmp, ok := params["RemarkCN"]; ok {
		if tmp1, ok := tmp.(string); ok {
			if len(tmp1) > Model.RemarkMaxLen {
				tmp1 = string([]byte(tmp1)[0:Model.RemarkMaxLen])
			}
			remarkcn = tmp1
		}
	}
	category, err := Model.MProducts.CreateCategory(name, namecn, remark, remarkcn)
	if err != nil {
		return this.Json(map[string]interface{}{"code": 40000, "msg": err.Error()})
	}
	return this.Json(map[string]interface{}{"code": 1, "data": category})
}

/*
用户
{_id:用户id,Name:用于登陆的用户名,NickName:昵称,Password:md5,Status:0 new/1 no/2 ok状态,CreateDate,Auth:[]int}
{
UID string `_id`
Name string
ExtType int
NickName string
Password string
Status int
CreateDate string
Auth int
}
留言
{_id,Content:string,ParentID:string 回复的留言,NickName:string ,UID,CreateDate,Version}

新闻
{_id,Content,CreateDate,Version,UID,ExtType:0首页显示/1/2,Score:int 分数判断是否热门，点击数}

产品分类
{_id,Remark,Name,CreateDate,ParentID,UID}
产品
{_id,Remark,CoverData,ClassID,CreateDate,Version,UID,Score:int}

首页图片
{_id,Remark,Base64}

*/

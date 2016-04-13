package Model

import (
	"fmt"
	"sync"

	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2/bson"
)

var MProducts *Products

type ProductsInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	Title      string        `bson:"Title" 产品名称`
	Images     []string      `bson:"Images" json:",omitempty" 图片id`
	CategoryID string        `bson:"IndexID"  json:"CategoryID" 分类`
	Remark     string        `bson:"Remark" 产品说明`
	ExtType    int           `bson:"ExtType" 0普通/1首页大图`
	Score      int           `bson:"Score" 分数 浏览次数 分数多的为热门商品`
	Version    float64       `bson:"Version" 时间戳 用于排序 索引`
	CreateDate string        `bson:"CreateDate"  创建时间`
}
type ProductsCategory struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"Name" 分类名称`
	Remark     string        `bson:"Remark" 说明`
	ExtType    int           `bson:"ExtType" json:",omitempty" 100`
	Version    float64       `bson:"Version" 时间戳　用于排序　索引`
	CreateDate string        `bson:"CreateDate"  创建时间`
}
type Image struct {
	ID         bson.ObjectId `bson:"_id"`
	Base64     string        `bson:"Base64" 数据`
	CategoryID string        `bson:"IndexID" json:"CategoryID,omitempty" 所属类，冗余字段`
	ProductID  string        `bson:"ProductID" 产品ID`
	ExtType    int           `bson:"ExtType" 200`
	Version    float64       `bson:"Version" 时间戳　用于排序　索引`
	CreateDate string        `bson:"CreateDate"  创建时间`
}
type Products struct {
	moduleBase
}

/*创建索引*/
func init() {
	MProducts = new(Products)
	err := MProducts.Init("dunhuaparts", "products")
	if err != nil {
		fmt.Println("Products初始化失败")
		panic(err)
	}
}

var (
	TitleMaxLen  = 120
	RemarkMaxLen = 1024
)

/*QueryProducts
@see :查询产品，默认按时间例序
@params :exttype []int
*/
func (this *Products) QueryProducts(start float64, limit int, indexid string, sorts []string, exttype ...int) (count int, products []ProductsInfo, err error) {
	//{_id,Remark,Name,CreateDate,ParentID,UID}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)

	if start < 0 {
		start = mystr.TimeStamp()
	}
	if limit < 1 {
		limit = 20
	}
	query := bson.M{"ExtType": bson.M{"$lt": 100}, "Version": bson.M{"$lt": start}}
	if len(exttype) > 0 {
		query["ExtType"] = bson.M{"$in": exttype}
	}
	if len(indexid) > 0 {
		query["IndexID"] = indexid
	}
	qs := col.Find(query) //.Select(bson.M{"Images": 0})
	if len(sorts) > 0 {
		qs.Sort(sorts...)
	} else {
		qs.Sort("-Version")
	}
	count, err = qs.Count()
	if err != nil {
		return
	}
	qs.Limit(limit)
	products = make([]ProductsInfo, 0)
	err = qs.All(&products)
	return
}

/*GetImage
@see :获取图片base64数据
@params :id
*/
func (this *Products) GetImage(id string) (base64 string, err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	qs := col.FindId(bson.ObjectIdHex(id))
	img := Image{}
	qs.Select(bson.M{"Base64": 1})
	err = qs.One(&img)
	if err != nil {
		return
	}
	base64 = img.Base64
	return
}

/*CreateProduct
@see :新增产品
@params : Images [base64]
*/
func (this *Products) CreateProduct(Title, IndexID, Remark string, Images []string, ExtType int) (product ProductsInfo, err error) {
	imageids := []string{}
	productid := bson.NewObjectId()
	tmpproduct := bson.M{"_id": productid, "Title": Title, "IndexID": IndexID, "Remark": Remark,
		"ExtType": ExtType, "Images": imageids, "Score": 0, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(),
	}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	objs := []interface{}{tmpproduct}
	for _, image := range Images {
		tmpid := bson.NewObjectId()
		tmp := bson.M{"_id": tmpid, "IndexID": IndexID, "ProductsID": productid.Hex(),
			"ExtType": 200, "Base64": image, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(),
		}
		imageids = append(imageids, tmpid.Hex())
		objs = append(objs, tmp)
	}
	tmpproduct["Images"] = imageids
	err = col.Insert(objs...)
	if err != nil {
		return
	}
	qs := col.FindId(tmpproduct["_id"])
	product = ProductsInfo{}
	err = qs.One(&product)
	return
}

/*AddImg
@see :添加产品的图片 先锁定--取产品信息--插入图片--更新产品信息
@params :
*/
func (this *Products) AddImg(id string, imgdata string) (err error) {
	if len(imgdata) < 1 {
		return nil
	}
	product := ProductsInfo{}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	qs := col.FindId(bson.ObjectIdHex(id))
	err = qs.One(&product)
	if err != nil {
		return
	}
	tmpid := bson.NewObjectId()
	tmp := bson.M{"_id": tmpid, "IndexID": product.CategoryID, "ProductsID": id,
		"ExtType": 200, "Base64": imgdata, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(),
	}
	err = col.Insert(tmp)
	if err != nil {
		return
	}
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$push": bson.M{"Images": tmpid.Hex()}})
	if err != nil {
		col.RemoveId(tmpid)
	}
	return
}

/*UpdateProduct
@see :更新产品
@params :Remark,Title,IndexID,Images...
*/
func (this *Products) UpdateProduct(id string, fields map[string]interface{}) (err error) {
	if len(fields) < 1 {
		return nil
	}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": fields})
	return
}

/*AddScore
@see :增加阅读
@params :
*/
func (this *Products) AddScore(id string) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"Score": 1}})
	return
}

/*CreateCategory
@see :创建分类
@params :
*/
func (this *Products) CreateCategory(Name, Remark string) (category ProductsCategory, err error) {
	//{_id,Remark,Name,CreateDate,ParentID,UID}
	tmp := bson.M{"_id": bson.NewObjectId(), "Name": Name, "Remark": Remark,
		"ExtType": 100, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(),
	}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.Insert(tmp)
	if err != nil {
		return
	}
	qs := col.FindId(tmp["_id"])
	//qs.Select(bson.M{"_id":1,"Title":1,"Remark":1,"ExtType":1,"Score":1,"Version":1})
	qs.Select(bson.M{"ExtType": 0})
	category = ProductsCategory{}
	err = qs.One(&category)
	return
}

/*QueryCategory
@see :查询分类
@params :
*/
func (this *Products) QueryCategory(start float64, limit int) (count int, categorys []ProductsCategory, err error) {
	//{_id,Remark,Name,CreateDate,UID}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)

	if start < 0 {
		start = mystr.TimeStamp()
	}
	if limit < 1 {
		limit = 20
	}
	qs := col.Find(bson.M{"ExtType": 100, "Version": bson.M{"$lt": start}}).Sort("-Version").Select(bson.M{"ExtType": 0})
	count, err = qs.Count()
	if err != nil {
		return
	}
	qs.Limit(limit)
	categorys = make([]ProductsCategory, 0)
	err = qs.All(&categorys)
	return
}

/*RemoveCategory
@see :删除分类及分类的产品和图片
@params :id
*/
func (this *Products) RemoveCategory(id string) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	_, err = col.RemoveAll(bson.M{"$or": []bson.M{
		bson.M{"_id": bson.ObjectIdHex(id), "ExtType": 100},
		bson.M{"IndexID": id},
	}})
	return
}

/*RemoveProducts
@see :删除产品及其图片
@params :id
*/
func (this *Products) RemoveProducts(ids []string) (err error) {
	if len(ids) < 1 {
		return
	}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	objIds := []bson.ObjectId{}
	for _, id := range ids {
		objIds = append(objIds, bson.ObjectIdHex(id))
	}
	_, err = col.RemoveAll(bson.M{"$or": []bson.M{
		bson.M{"_id": bson.M{"$in": objIds}},
		bson.M{"IndexID": bson.M{"$in": ids}, "ExtType": 200},
	}})
	return
}

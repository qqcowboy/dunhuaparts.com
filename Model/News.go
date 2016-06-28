package Model

import (
	"fmt"

	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2/bson"
)

type NewsInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	Title      string        `bson:"Title" 产品名称`
	Content    string        `bson:"Content" json:",omitempty" 内容`
	Lead       string        `bson:"Lead" 导语`
	Author     string        `bson:"Author" 作者`
	ExtType    int           `bson:"ExtType" 0中文/1英文`
	Top        int           `bson:"Top" 0/1置顶`
	Score      int           `bson:"Score" 分数 浏览次数 分数多的为热门`
	Version    float64       `bson:"Version" 时间戳 用于排序 索引`
	CreateDate string        `bson:"CreateDate"  创建时间`
}
type News struct {
	moduleBase
}

const (
	MAXLEADSIZE    = 1024 * 2
	MAXCONTENTSIZE = 1024 * 1024
)

var MNews *News

/*创建索引*/
func init() {
	MNews = new(News)
	err := MNews.Init("dunhuaparts", "news")
	if err != nil {
		fmt.Println("News初始化失败")
		panic(err)
	}
}

/*QueryNews
@see :查询新闻，默认按时间例序
@params :exttype []int
*/
func (this *News) QueryNews(start float64, limit int, sorts []string, top int, key string, exttype ...int) (count int, news []NewsInfo, err error) {
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
	query := bson.M{"Version": bson.M{"$lt": start}}
	if len(exttype) > 0 {
		query["ExtType"] = bson.M{"$in": exttype}
	}
	if top > 0 {
		query["Top"] = 1
	}
	if len(key) > 0 {
		reg := bson.RegEx{fmt.Sprintf("^.*%s.*$", key), "g"}
		query["$or"] = []bson.M{bson.M{"Content": reg}, bson.M{"Title": reg}}
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
	qs.Select(bson.M{"Content": 0})
	news = make([]NewsInfo, 0)
	err = qs.All(&news)
	return
}

/*CreateNews
@see :新增News
*/
func (this *News) CreateNews(Title, Lead, Content string, ExtType, Top int, author string) (news NewsInfo, err error) {
	tmp := bson.M{"_id": bson.NewObjectId(), "Title": Title, "Content": Content, "Lead": Lead, "Top": Top,
		"ExtType": ExtType, "Score": 0, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(), "Author": author,
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
	news = NewsInfo{}
	err = qs.One(&news)
	fmt.Println(tmp, news)
	return
}

/*AddScore
@see :增加阅读
@params :
*/
func (this *News) AddScore(id string) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"Score": 1}})
	return
}

/*Get
@see :新闻详情
@params :id
*/
func (this *News) Get(id string) (news NewsInfo, err error) {
	//{_id,Remark,Name,CreateDate,ParentID,UID}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)

	qs := col.FindId(bson.ObjectIdHex(id))
	err = qs.One(&news)
	this.AddScore(id)
	return
}

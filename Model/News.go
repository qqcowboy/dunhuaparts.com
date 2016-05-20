package Model

import (
	"fmt"

	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2/bson"
)

type NewsInfo struct {
	ID         bson.ObjectId `bson:"_id"`
	Title      string        `bson:"Title" 产品名称`
	Content    string        `bson:"Content" 内容`
	Lead       string        `bson:"Lead" 导语`
	User       string        `bson:"User" 作者`
	ExtType    int           `bson:"ExtType" 0/1/2`
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

/*CreateNews
@see :新增News
*/
func (this *News) CreateNews(Title, Lead, Content string, ExtType int) (news NewsInfo, err error) {
	tmp := bson.M{"_id": bson.NewObjectId(), "Title": Title, "Content": Content, "Lead": Lead,
		"ExtType": ExtType, "Score": 0, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(),
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

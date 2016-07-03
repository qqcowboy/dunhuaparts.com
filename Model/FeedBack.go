package Model

import (
	"fmt"

	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2/bson"
)

type FeedbackInfo struct {
	ID         bson.ObjectId   `bson:"_id"`
	Title      string          `bson:"Title" 标题`
	Content    string          `bson:"Content" json:",omitempty" 内容`
	ExtType    int             `bson:"ExtType" 0`
	Hide       int             `bson:"Hide" 0公开/1隐藏`
	UserName   string          `bson:"UserName" 用户名`
	Mail       string          `bson:"Mail" json:",omitempty" 邮件地址`
	Phone      string          `bson:"Phone" json:",omitempty" 电话号码`
	IP         string          `bson:"IP" IP地址`
	Reply      []FeedbackReply `bson:"Reply" json:",omitempty" 回复`
	Version    float64         `bson:"Version" 时间戳 用于排序 索引`
	CreateDate string          `bson:"CreateDate"  创建时间`
}
type FeedbackReply struct {
	ID         string  `bson:"ID" id`
	UserName   string  `bson:"UserName" 用户名`
	Content    string  `bson:"Content" json:",omitempty" 内容`
	Version    float64 `bson:"Version" 时间戳 用于排序 索引`
	CreateDate string  `bson:"CreateDate"  创建时间`
}
type Feedback struct {
	moduleBase
}

var MFeedback *Feedback

/*创建索引*/
func init() {
	MFeedback = new(Feedback)
	err := MFeedback.Init("dunhuaparts", "feedback")
	if err != nil {
		fmt.Println("Feedback初始化失败")
		panic(err)
	}
}

/*QueryFeedback
@see :查询新闻，默认按时间例序
@params :exttype []int
*/
func (this *Feedback) QueryFeedback(start float64, limit int, sorts []string, hide int, key, email string, selField map[string]interface{}, exttype ...int) (count int, Feedback []FeedbackInfo, err error) {
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
	if hide > 0 {
		if hide > 1 {
			query["Hide"] = 0
		} else {
			query["Hide"] = 1
		}
	}
	if len(key) > 0 {
		reg := bson.RegEx{fmt.Sprintf("^.*%s.*$", key), "g"}
		query["$or"] = []bson.M{bson.M{"Content": reg}, bson.M{"Title": reg}}
	}
	if len(email) > 0 {
		query["Mail"] = email
	}
	fmt.Println(query)
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
	qs.Select(selField)
	Feedback = make([]FeedbackInfo, 0)
	err = qs.All(&Feedback)
	return
}

/*CreateFeedback
@see :新增Feedback
*/
func (this *Feedback) CreateFeedback(Title, Content string, ExtType, Hide int, User, Phone, Mail, IP string) (Feedback FeedbackInfo, err error) {
	tmp := bson.M{"_id": bson.NewObjectId(), "Title": Title, "Content": Content, "Hide": Hide, "User": User, "IP": IP,
		"ExtType": ExtType, "Phone": Phone, "Version": mystr.TimeStamp(), "CreateDate": mystr.Date(), "Mail": Mail,
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
	Feedback = FeedbackInfo{}
	err = qs.One(&Feedback)
	return
}

/*Get
@see :详情
@params :id
*/
func (this *Feedback) Get(id string) (Feedback FeedbackInfo, err error) {
	//{_id,Remark,Name,CreateDate,ParentID,UID}
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)

	qs := col.FindId(bson.ObjectIdHex(id))
	err = qs.One(&Feedback)
	return
}

/*AddReply
@see :新增回复
@params :fdid, username, content
*/
func (this *Feedback) AddReply(fdid, username, content string) (result FeedbackReply, err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	result = FeedbackReply{
		ID:         bson.NewObjectId().Hex(),
		UserName:   username,
		Content:    content,
		Version:    mystr.TimeStamp(),
		CreateDate: mystr.Date(),
	}
	err = col.UpdateId(bson.ObjectIdHex(fdid), bson.M{"$push": result})
	return
}

/*DelReply
@see :删除回复
@params :fdid,rpid
*/
func (this *Feedback) DelReply(fdid, rpid string) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(fdid), bson.M{"$pull": bson.M{"Reply": bson.M{"ID": rpid}}})
	return
}

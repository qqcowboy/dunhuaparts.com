package Model

import (
	"fmt"
	"sync"

	"github.com/qqcowboy/lib/mystr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var UserModel *User

func init() {
	UserModel = new(User)
	UserModel.Init("dunhuaparts", "user")
	UserModel.addlock = new(sync.Mutex)
}

type UserNum struct {
	ID      bson.ObjectId `bson:"_id" json:"ID"`
	ExtType int
	AuthNum int
}
type UserInfo struct {
	UID      bson.ObjectId `bson:"_id" json:"UID"`
	UserNum  int           `bson:"UserNum"`
	UserName string        `bson:"UserName"`
	Password string        `bson:"Password" json:",omitempty"`
	ExtType  int           `bson:"ExtType" json:",omitempty"`
	Enable   bool          `bson:"Enable" json:",omitempty"`
	Auths    []int         `bson:"Auths" json:",omitempty"`
}

type User struct {
	moduleBase
	addlock sync.Locker
}

/*AddUser
@see : 新增用户
*/
func (this *User) AddUser(name, psw string, enable bool, auths []int) error {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	this.addlock.Lock()
	defer this.addlock.Unlock()
	_, err := col.Find(bson.M{"Name": name}).Count()
	if err != mgo.ErrNotFound {
		return fmt.Errorf("该用户已存在")
	}

	autonum := UserNum{}
	_, err := col.Find(bson.M{"ExtType": 100}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"AutoNum": 1}},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}, &autonum)
	if err != nil {
		return fmt.Errorf("生成用户编号失败：%s", err.Error())
	}

	tmp := bson.M{"_id": bson.NewObjectId(), "Name": name, "Password": psw, "ExtType": 0, "UserNum": autonum.AutoNum,
		"Version": mystr.TimeStamp(), "CreateDate": mystr.Date(), "Auths": auths, "Enable": enable,
	}
	return col.Insert(tmp)
}

/*FindId
@see : 删除用户
*/
func (this *User) DelUser(id string) error {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	this.addlock.Lock()
	defer this.addlock.Unlock()
	n, err := col.Find(bson.M{"ExtType": 1}).Count()
	if err != nil {
		return err
	}
	if n == 1 {
		return fmt.Errorf("不能删除最后一个用户")
	}
	return col.RemoveId(bson.ObjectIdHex(id))
}

/*FindId
@see : 查询用户信息
*/
func (this *User) FindId(id string) (user UserInfo, err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	user = UserInfo{}
	err = col.FindId(bson.ObjectIdHex(id)).One(&user)
	return
}

/*AddAuth
@see : 增加权限
*/
func (this *User) AddAuth(id string, auth ...int) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$pushAll": auth})
	return
}

/*DelAuth
@see : 删除权限
*/
func (this *User) AddAuth(id string, auth ...int) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$pullAll": auth})
	return
}

/*Query
@see : 查询用户列表
*/
func (this *User) Query(skip float64, limit int) (users []UserInfo, count int, err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	users = make([]UserInfo, 0)
	if skip < 1 {
		skip = mystr.TimeStamp()
	}
	qs := col.Find(bson.M{"ExtType": 0, "Version": bson.M{"$lt": skip}})
	count, err = qs.Count()
	if err != nil {
		return
	}
	qs = qs.Limit(limit)
	err = qs.All(&users)
	return
}

/*FindOne
@see : 查询用户（一般用于登陆）
*/
func (this *User) FindOne(name, psw string, num int) (user UserInfo, err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	user = UserInfo{}
	var qs mgo.Query
	if len(name) > 0 {
		qs = col.Find(bson.M{"Name": name, "Password": psw, "Enable": true})
	} else {
		qs = col.Find(bson.M{"UserNum": num, "Password": psw, "Enable": true})
	}
	qs.Select(bson.M{"ExtType": 0, "Password": 0, "Enable": 0})
	err := qs.One(&user)
	return
}

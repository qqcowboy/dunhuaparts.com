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
	AutoNum int
}

/*UserInfo
@see : Auths 0超级管理/<10后台管理
*/
type UserInfo struct {
	UID      bson.ObjectId `bson:"_id" json:"UID"`
	UserNum  int           `bson:"UserNum"`
	UserName string        `bson:"UserName"`
	Mail     string        `bson:"Mail" json:",omitempty"`
	Password string        `bson:"Password" json:",omitempty"`
	ExtType  int           `bson:"ExtType" json:",omitempty"`
	Enable   bool          `bson:"Enable" json:",omitempty"`
	Auths    []int         `bson:"Auths" json:",omitempty"`
}

type User struct {
	moduleBase
	addlock sync.Locker
}

func (this UserInfo) IsSupAdmin() bool {
	for _, auth := range this.Auths {
		if auth == 1 {
			return true
		}
	}
	return false
}
func (this UserInfo) IsAdmin() bool {
	for _, auth := range this.Auths {
		if auth < 10 {
			return true
		}
	}
	return false
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
	_, err = col.Find(bson.M{"ExtType": 100}).Apply(mgo.Change{
		Update:    bson.M{"$inc": bson.M{"AutoNum": 1}},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}, &autonum)
	if err != nil {
		return fmt.Errorf("生成用户编号失败：%s", err.Error())
	}
	if autonum.AutoNum < 1000 {
		_, err = col.Find(bson.M{"ExtType": 100}).Apply(mgo.Change{
			Update:    bson.M{"$inc": bson.M{"AutoNum": 1000}},
			Upsert:    true,
			Remove:    false,
			ReturnNew: true,
		}, &autonum)
		if err != nil {
			return fmt.Errorf("生成用户编号失败：%s", err.Error())
		}
	}
	tmp := bson.M{"_id": bson.NewObjectId(), "Name": name, "Password": psw, "ExtType": 0, "UserNum": autonum.AutoNum,
		"Version": mystr.TimeStamp(), "CreateDate": mystr.Date(), "Auths": auths, "Enable": enable,
	}
	return col.Insert(tmp)
}

/*DelUser
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
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$pushAll": bson.M{"Auths": auth}})
	return
}

/*DelAuth
@see : 删除权限
*/
func (this *User) DelAuth(id string, auth ...int) (err error) {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{"$pullAll": bson.M{"Auths": auth}})
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
	var qs *mgo.Query
	if len(name) > 0 {
		qs = col.Find(bson.M{"UserName": name, "Password": psw, "Enable": true})
	} else {
		qs = col.Find(bson.M{"UserNum": num, "Password": psw, "Enable": true})
	}
	qs.Select(bson.M{"ExtType": 0, "Password": 0, "Enable": 0})
	err = qs.One(&user)
	return
}
func (this *User) ChPsw(uid bson.ObjectId, oldpsw, newpsw string) error {
	mongo := this.mSession()
	defer mongo.Close()
	mdb := mongo.DB(this.db)
	col := mdb.C(this.coll)
	var qs *mgo.Query
	qs = col.Find(bson.M{"_id": uid, "Password": oldpsw, "Enable": true})
	c, err := qs.Count()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("密码不正确")
	}
	return col.Update(bson.M{"_id": uid, "Password": oldpsw}, bson.M{"$set": bson.M{"Password": newpsw}})
}

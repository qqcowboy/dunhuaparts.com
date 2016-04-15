package Model

import (
	"fmt"
	"sync"

	"github.com/qqcowboy/lib/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type moduleBase struct {
	mongo *mgo.Session
	db    string
	coll  string
}

func (this *moduleBase) Init(db, collection string) error {
	var err error
	this.mongo, err = GetMgoSession()
	this.db = db
	this.coll = collection
	return err
}
func (this *moduleBase) mSession() *mgo.Session {
	return this.mongo.Clone()
}

/*UpdateId
@see :通用更新
@params :Remark,Title,IndexID,Images...
*/
func (this *moduleBase) UpdateId(id string, fields map[string]interface{}) (err error) {
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

var session *mgo.Session
var lock sync.Locker

func init() {
	lock = new(sync.Mutex)
}
func GetMgoSession() (msession *mgo.Session, err error) {
	lock.Lock()
	defer lock.Unlock()
	if session == nil {
		uri := ""
		uri, err = config.Config.GetValue("mongo", "url")
		if err != nil {
			uri = "mongodb://127.0.0.1:27017"
		}
		session, err = mgo.Dial(uri)
		if err != nil {
			fmt.Println("连接mongodb数据库失败", err)
			return
		}
	}
	msession = session.Clone()
	return
}

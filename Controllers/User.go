package Controllers

import (
	//_ "github.com/qqcowboy/dunhuaparts.com/Model"
	"github.com/qqcowboy/dunhuaparts.com/System/Web"
)

type User struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(User{})
}
func (this *User) OnLoad() {
}

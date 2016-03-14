package Web

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ControllerNotExist error = errors.New("Controller不存在")
var ActionNotExist error = errors.New("Action不存在")

type ControllerInfo struct {
	ControllerName string
	ControllerType reflect.Type
	Methods        map[string]string
}

type ControllersCollection struct {
	Controllers map[string]*ControllerInfo
}

func NewControllerInfo(controllerName string, t reflect.Type) *ControllerInfo {
	c := new(ControllerInfo)
	c.ControllerName = controllerName
	c.ControllerType = t
	c.Methods = make(map[string]string)
	return c
}
func NewControllersCollection() *ControllersCollection {
	c := &ControllersCollection{Controllers: make(map[string]*ControllerInfo)}
	return c
}

func (this *ControllersCollection) Add(c interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	rt := reflect.TypeOf(c)
	/*取类型名和类型*/
	typeName, t := this.getTypeNameAndType(rt)
	ctlinfo := NewControllerInfo(typeName, t)
	//如果原来是变为非指针,变为指针类型，
	if rt.Kind() != reflect.Ptr {
		rt = reflect.PtrTo(rt)
	}
	/*获取该类型下的函数*/
	this.getTypeMethod(rt, ctlinfo)
	this.Controllers[strings.ToLower(typeName)] = ctlinfo
	return nil
}
func (this *ControllersCollection) getTypeMethod(t reflect.Type, ctlinfo *ControllerInfo) {
	for i, j := 0, t.NumMethod(); i < j; i++ {
		m := t.Method(i)
		strMethodName := strings.ToLower(m.Name)
		ctlinfo.Methods[strMethodName] = m.Name
	}
}
func (this *ControllersCollection) getTypeNameAndType(t reflect.Type) (string, reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	strName := t.Name()
	if strings.HasSuffix(strName, "Controller") && len(strName) > len("Controller") {
		strName = strings.TrimRight(strName, "Controller")
	}
	return strName, t
}

func (this *ControllersCollection) GetController(routeData map[string]interface{}) (reflect.Value, error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err := e.(error)
			App.Log.AddErrMsg("In ControllerInfo.GetController:\t" + err.Error())
		}
	}()
	var result reflect.Value
	controllerName := strings.ToLower(fmt.Sprintf("%v", routeData["controller"]))
	actionName := strings.ToLower(fmt.Sprintf("%v", routeData["action"]))
	var controllers = this.Controllers
	ctlinfo, ok := controllers[controllerName]
	if !ok {
		return result, ControllerNotExist
	}
	methodName, ok := ctlinfo.Methods[actionName]
	if !ok {
		methodName = "Default"
		views := fmt.Sprintf("%v", routeData["views"])
		if len(views) > 0 {
			routeData["views"] = fmt.Sprintf("%s/%s", strings.Title(actionName), routeData["views"])
		} else {
			routeData["views"] = strings.Title(actionName)
		}
	}
	routeData["controller"] = ctlinfo.ControllerName
	routeData["action"] = methodName
	result = reflect.New(ctlinfo.ControllerType)
	return result, nil
}

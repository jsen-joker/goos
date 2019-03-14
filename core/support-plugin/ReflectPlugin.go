package support_plugin

import (
	"log"
	"reflect"
)

func ReflectCreatePlugin(funcInter interface{}) *Plugin {
	v := reflect.ValueOf(funcInter)
	if v.Kind() != reflect.Func {
		log.Fatal("funcInter is not func")
	}

	values := v.Call(make([]reflect.Value, 0)) //方法调用并返回值
	return values[0].Interface().(*Plugin)
	//fmt.Print(len(values))
	//for i := range values {
	//	fmt.Println(values[i].Interface().(*Plugin))
	//}
	//return nil
}

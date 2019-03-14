package manager

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jsen-joker/goos/core/support-plugin/manager/utils"
	utils2 "github.com/jsen-joker/goos/core/utils"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

//type Pipline struct {
//	Name string
//	handle interface{}
//}

type Routes map[string] *Route
var routes = make(Routes)
func RegisterRouter(route *Route) {
	routes[route.Name] = route
}

type Pipline struct {
	Name          string
	IgnorePattern []string
	Handle        interface{}
	Priority int
}
type Piplines []Pipline
func (p Piplines) Len() int {
	return len(p)
}
func (p Piplines) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Piplines) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}

//var piplines = make(Piplines)
var piplines Piplines
func AddPipline(name string, priority int, pattern []string, handle interface{})  {
	piplines = append(piplines, Pipline{
		Name:          name,
		IgnorePattern: pattern,
		Handle:        handle,
		Priority:	priority,
	})
}

func CreateRouter() *mux.Router {
	AddPipline("log", 1000, []string {}, utils.HttpPipeLogger)
	AddPipline("rest", 1000, []string {}, utils.HttpPipeRest)

	sort.Sort(sort.Reverse(piplines))

	for _, v := range piplines {
		fmt.Println(v.Name)
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		mapper := fmt.Sprintf("mapper(%20s) pipline(", route.Name)

		for _, pipline := range piplines {
			path := ""
			if !strings.HasSuffix(route.Pattern, "/**") {
				path = route.Pattern
			} else {
				path = route.Pattern[:len(route.Pattern) - 3]
			}
			if !utils2.Matcher(path, pipline.IgnorePattern) {
				v := reflect.ValueOf(pipline.Handle)
				if v.Kind() != reflect.Func {
					log.Fatal("funcInter is not func")
				}
				mapper += fmt.Sprintf("%5s ", pipline.Name)
				in := make([]reflect.Value, 2)
				in[0] = reflect.ValueOf(handler)
				in[1] = reflect.ValueOf(route.Name)
				values := v.Call(in) //方法调用并返回值
				handler = values[0].Interface().(http.Handler)
			} else {
				mapper += fmt.Sprintf("%5s ", "")
			}
		}

		mapper += ") "

		//handler = utils.HttpPipeLogger(handler, route.Name)
		//handler = utils.HttpPipeRest(handler, route.Name)

		rt := router.Methods(route.Method)

		if !strings.HasSuffix(route.Pattern, "/**") {
			rt.Path(route.Pattern)
			mapper += fmt.Sprintf("%6s %6s %s", route.Method, "path", route.Pattern)
		} else {
			prefix := route.Pattern[:len(route.Pattern) - 3]
			mapper += fmt.Sprintf("%6s %6s %s", route.Method, "prefix", prefix)
			rt.PathPrefix(prefix)
		}
		log.Println(mapper)
		rt.Name(route.Name).Handler(handler)
	}
	_ = router.PathPrefix("/").Name("default").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := utils2.Failed("no api")
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}))
	return router
}
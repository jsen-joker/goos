package utils

import (
	"reflect"
	"strings"
)

type Pagination struct {
	Current int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Total int64 `json:"total"`
}

type SimpleResponse map[string]interface{}

//type SimpleResponse struct {
//	Code int `json:"code"`
//	Data interface{} `json:"data,omitempty"`
//	Extra interface{} `json:"extra,omitempty"`
//	List interface{} `json:"list,omitempty"`
//	Pagination *Pagination `json:"pagination,omitempty"`
//	Msg string `json:"msg,omitempty"`
//}

func (r *SimpleResponse) SetCode(code int) *SimpleResponse {
	(*r)["code"] = code
	//r.Code = code
	return r
}

func (r *SimpleResponse) Put(name string, value interface{}) *SimpleResponse {
	(*r)[name] = value
	//r.Code = code
	return r
}

func (r *SimpleResponse) SetData(data interface{}) *SimpleResponse {
	(*r)["data"] = data
	//r.Data = data
	return r
}

func (r *SimpleResponse) SetExtra(extra interface{}) *SimpleResponse {
	(*r)["extra"] = extra

	//r.Extra = extra
	return r
}

func (r *SimpleResponse) SetList(list interface{}) *SimpleResponse {
	(*r)["list"] = list

	//r.List = list
	return r
}

func (r *SimpleResponse) SetPagination(pagination *Pagination) *SimpleResponse {
	(*r)["pagination"] = pagination

	//r.Pagination = pagination
	return r
}

func (r *SimpleResponse) SetMsg(msg string) *SimpleResponse {
	(*r)["msg"] = msg

	//r.Msg = msg
	return r
}

func (r *SimpleResponse) PutAll(obj interface{}) *SimpleResponse {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		jTag := t.Field(i).Tag.Get("json")
		if jTag != "" {
			if jTag != "-" {
				tags := strings.Split(jTag, ",")
				(*r)[tags[0]] = v.Field(i).Interface()
				//if len(tags) > 1 {
				//	if tags[1] == "omitempty" {
				//		if v.Field(i).Interface() != "" {
				//			(*r)[tags[0]] = v.Field(i).Interface()
				//		}
				//	} else {
				//		(*r)[tags[0]] = v.Field(i).Interface()
				//	}
				//} else {
				//	(*r)[tags[0]] = v.Field(i).Interface()
				//}
			}
		} else {
			(*r)[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return r
}

func Succeed(data interface{}) *SimpleResponse {
	sc := &SimpleResponse{}
	(*sc)["code"] = 0
	if data != nil {
		(*sc)["data"] = data
	}
	return sc
	//if data == nil {
	//	return &SimpleResponse{Code: 0}
	//} else {
	//	return &SimpleResponse{Code: 0, Data: data}
	//}
}

func Failed(msg string) *SimpleResponse {
	sc := &SimpleResponse{}
	(*sc)["code"] = 1
	(*sc)["msg"] = msg
	return sc

	//return &SimpleResponse{Code: 1, Msg: msg}
}
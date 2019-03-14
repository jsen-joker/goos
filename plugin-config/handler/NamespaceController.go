package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-config/entity"
	"github.com/jsen-joker/goos/plugin-config/service"
	"io/ioutil"
	"net/http"
)

// 创建更新 namespace
func Namespace(w http.ResponseWriter, r *http.Request)  {

	body, err := ioutil.ReadAll(r.Body)
	//fmt.Println(r.Header.Get("Content-Type"))
	if err != nil {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}
	//else {
	//	fmt.Println(bytes.NewBuffer(body).String())
	//}

	var namespace entity.Namespace
	if err := json.Unmarshal([]byte(bytes.NewBuffer(body).String()), &namespace); err != nil {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		if namespace.Name == "" || namespace.Desc == "" {
			resp := utils.Failed(errors.New("need name and desc").Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}
		if eff, err := service.Namespace(namespace); err == nil {
			resp := utils.Succeed(eff)
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
		} else {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
		}
	}
}

// 获取所有namespace
func GetNamespaceList(w http.ResponseWriter, r *http.Request)  {
	currentPage, err := utils.GetInt64Parameter("currentPage", true, 1, w, r)
	if err != nil {
		return
	}
	pageSize, err := utils.GetInt64Parameter("pageSize", true, 20, w, r)
	if err != nil {
		return
	}


	query, err := utils.GetParameter("query", true, "", w, r)
	if err != nil {
		return
	}

	if list, total, err := service.GetNamespaceList(currentPage, pageSize, query); err == nil {
		resp := utils.Succeed(nil).SetList(list).SetPagination(&utils.Pagination{
			Current: currentPage,
			PageSize: pageSize,
			Total: total,
		})
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}
}

// 删除namespace
func DeleteNamespace(w http.ResponseWriter, r *http.Request) {

	id, err := utils.GetInt64Parameter("id", false, 0, w, r)
	if err != nil {
		return
	}

	if eff, err := service.DeleteNamespace(id); err == nil {
		resp := utils.Succeed(eff)
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	}

}
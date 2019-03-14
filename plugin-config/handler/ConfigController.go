package handler

import (
	"bytes"
	"encoding/json"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-config/entity"
	"github.com/jsen-joker/goos/plugin-config/entity/vo"
	"github.com/jsen-joker/goos/plugin-config/service"
	"io/ioutil"
	"net/http"
)

// 配置文件处理

// 发布配置文件
func Config(w http.ResponseWriter, r *http.Request) {
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

	var config entity.Config
	if err := json.Unmarshal([]byte(bytes.NewBuffer(body).String()), &config); err != nil {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		if eff, err := service.Config(config); err == nil {
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



	//fId := r.PostForm["id"]
	//dataId, err := utils.GetParameter("dataId", w, r)
	//if err != nil {
	//	return
	//}
	//groupId, err := utils.GetParameter("groupId", w, r)
	//if err != nil {
	//	return
	//}
	//content, err := utils.GetParameter("content", w, r)
	//if err != nil {
	//	return
	//}
	//fFType, err := utils.GetParameter("type", w, r)
	//if err != nil {
	//	return
	//}
	//var id int64 = 0
	//if len(fId) > 0 {
	//	id, _ = strconv.ParseInt(fId[0], 10, 64)
	//}


}

// 获取配置文件
func GetConfigList(w http.ResponseWriter, r *http.Request)  {


	currentPage, err := utils.GetInt64Parameter("currentPage", true, 1, w, r)
	if err != nil {
		return
	}
	pageSize, err := utils.GetInt64Parameter("pageSize", true, 20, w, r)
	if err != nil {
		return
	}

	queryDataId, err := utils.GetParameter("queryDataId", true, "", w, r)
	if err != nil {
		return
	}

	queryGroup, err := utils.GetParameter("queryGroup", true, "", w, r)
	if err != nil {
		return
	}

	namespace, err := utils.GetInt64Parameter("namespace", false, 0, w, r)
	if err != nil {
		return
	}


	if list, total, err := service.GetConfigList(currentPage, pageSize, namespace, queryDataId, queryGroup); err == nil {
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

// 获取配置文件的配置信息（详细数据）
func GetConfig(w http.ResponseWriter, r *http.Request)  {

	id, err := utils.GetInt64Parameter("id", false, 0, w, r)
	if err != nil {
		return
	}

	if result, err := service.GetConfig(id); err == nil {
		resp := utils.Succeed(nil).PutAll(result)
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

// 删除配置文件
func DeleteConfig(w http.ResponseWriter, r *http.Request) {

	id, err := utils.GetInt64Parameter("id", false, 0, w, r)
	if err != nil {
		return
	}

	if eff, err := service.DeleteConfig(id); err == nil {
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




///////////////////////////////////   RSA CLIENT API      ///////////////////////////////
// 传入配置文件MD5，比较文件MD5，不一致则返回RSA加密数据


func RsaGetConfig(w http.ResponseWriter, r *http.Request) {
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

	var configQuery vo.ConfigQuery
	if err := json.Unmarshal([]byte(bytes.NewBuffer(body).String()), &configQuery); err != nil {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {
		if len(configQuery.ConfigList) == 0 {
			resp := utils.Succeed(nil)
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
			return
		}
		if result, err := service.RsaGetConfig(configQuery); err == nil {
			resp := utils.Succeed(result)
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

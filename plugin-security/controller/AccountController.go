package controller

import (
	"bytes"
	"encoding/json"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-security/entity"
	"github.com/jsen-joker/goos/plugin-security/entity/vo"
	"github.com/jsen-joker/goos/plugin-security/service"
	utils2 "github.com/jsen-joker/goos/plugin-security/utils"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request)  {


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

	var account vo.AccountQuery
	if err := json.Unmarshal([]byte(bytes.NewBuffer(body).String()), &account); err != nil {
		resp := utils.Failed(err.Error())
		if err := json.NewEncoder(w).Encode(resp); err != nil{
			panic(err)
		}
	} else {

		if account, err := service.Login(account); err == nil {
			token, err := utils2.CreateToken(account.(entity.Account).ID, account.(entity.Account).Username)
			if err != nil {
				resp := utils.Failed(err.Error())
				if err := json.NewEncoder(w).Encode(resp); err != nil{
					panic(err)
				}
			} else {
				resp := utils.Succeed(account).SetExtra(token)
				if err := json.NewEncoder(w).Encode(resp); err != nil{
					panic(err)
				}
			}
		} else {
			resp := utils.Failed(err.Error())
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
		}
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	resp := utils.Failed("no impl")
	if err := json.NewEncoder(w).Encode(resp); err != nil{
		panic(err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resp := utils.Succeed(nil)
	if err := json.NewEncoder(w).Encode(resp); err != nil{
		panic(err)
	}
}

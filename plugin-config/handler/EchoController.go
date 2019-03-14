package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jsen-joker/goos/core/support-db"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-security/entity"
	"net/http"
)

func ConfigEcho(w http.ResponseWriter, r *http.Request)  {


	t, err := support_db.QueryOne((&orm.QueryWrapper{}).Entity(entity.Account{ID:1}))
	if err != nil {
		panic(err)
	}
	fmt.Print(t)


	resp := utils.Succeed(nil)
	if err := json.NewEncoder(w).Encode(resp); err != nil{
		panic(err)
	}
}

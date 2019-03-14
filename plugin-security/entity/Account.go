package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

type Account struct {
	ID int64 `name:"id"PK:"true"auto:"true"json:"id"`
	Username string `name:"username"json:"username"`
	Password string `name:"password"json:"-"`

	TableName orm.TableName `name:"account"json:"-"`
}
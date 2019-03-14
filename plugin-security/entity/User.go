package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

type User struct {
	ID int64 `name:"id"PK:"true"auto:"true"json:"id"`
	Name string `name:"name"json:"name"`
	Phone string `name:"phone"json:"phone"`
	Avatar string `name:"avatar"json:"avatar"`
	Email string `name:"email"json:"email"`
	AccountID int64 `name:"account_id"json:"accountId"`
	SelectedNamespace int64 `name:"selected_namespace"json:"selectedNamespace"`

	TableName orm.TableName `name:"user"json:"-"`
}
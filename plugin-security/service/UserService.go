package service

import (
	"github.com/jsen-joker/goos/core/support-db"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/plugin-security/entity"
)

func GetUserByAccountId(id int64) (result interface{}, err error) {
	return support_db.QueryOne((&orm.QueryWrapper{}).Entity(entity.User{AccountID: id}))
}

func ChangeSelectedNamespace(id int64, selectedNamespace int64) (eff int64, err error) {
	return support_db.UpdateBy(entity.User{SelectedNamespace:selectedNamespace}, entity.User{AccountID:id})
}
package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

func Init() {
	orm.Register(new(Account))
	orm.Register(new(User))
}
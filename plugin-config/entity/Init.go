package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

func Init() {
	orm.Register(new(Config))
	orm.Register(new(Namespace))
}
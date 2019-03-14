package test

import (
	"fmt"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/plugin-config/entity"
	"testing"
)

func TestQueryWrapper(t *testing.T) {
	q := (&orm.QueryWrapper{}).Entity(entity.Config{}).Where("id", 1).And().Where("data_id",
		"lucy").Or().GroupAnd("id", 3).And("data_id", "jack").Build()
	fmt.Println(q.GetSelectSql())
	fmt.Println(q.GetCountSql())
	fmt.Println(q.GetValues())
}
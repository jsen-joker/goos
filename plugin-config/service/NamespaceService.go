package service

import (
	"github.com/jsen-joker/goos/core/support-db"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/plugin-config/entity"
	"github.com/snluu/uuid"
	"time"
)

// 创建或者更新
func Namespace(namespace entity.Namespace) (eff int64, err error)  {

	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)

	uuid := uuid.Rand()

	namespace.GmtModified = tm.Format("2006-01-02 15:04:05")
	namespace.NamespaceID = uuid.Hex()
	namespace.System = 1

	if namespace.Id == 0 {
		// do create
		namespace.GmtCreate = namespace.GmtModified

		return support_db.Insert(namespace)
	} else {
		// do update
		return support_db.Update(namespace)
	}

}

func GetNamespaceList(currentPage int64, pageSize int64, query string) (list []interface{}, total int64, err error) {
	skip := (currentPage - 1) * pageSize
	qc := (&orm.QueryWrapper{}).Entity(entity.Namespace{})
	if query != "" {
		qc.Where("name", "%" + query + "%").Like()
	}
	total, e := support_db.Count(qc)
	if e != nil {
		return nil, 0, e
	}
	q := (&orm.QueryWrapper{}).Select("id", "namespace_id", "name", "desc", "system").Entity(entity.Namespace{}).Skip(skip).Capacity(pageSize)
	if query != "" {
		q.Where("name", "%" + query + "%").Like()
	}
	rs, e := support_db.Query(q)
	if e == nil && rs == nil {
		return []interface{}{}, total, nil
	}
	return rs, total, e
}

func DeleteNamespace(id int64) (eff int64, err error) {
	if eff, err := support_db.Delete(entity.Namespace{Id: id}); err == nil {
		return eff, err
	} else {
		return eff, err
	}
}


package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

type Namespace struct {
	Id int64 `name:"id"PK:"true"auto:"true"json:"id"`
	NamespaceID string `name:"namespace_id"json:"namespaceId"`
	Name string `name:"name"json:"name,omitempty"`
	Desc string `name:"desc"json:"desc,omitempty"`
	GmtCreate string `name:"gmt_create"json:"gmtCreate,omitempty"`
	GmtModified string `name:"gmt_modified"json:"gmtModified,omitempty"`
	System int `name:"system"json:"system,omitempty"`

	TableName orm.TableName `name:"plugin_namespace"json:"-"`
}
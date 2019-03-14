package entity

import "github.com/jsen-joker/goos/core/support-db/orm"

type Config struct {
	Id int64 `name:"id"PK:"true"auto:"true"json:"id"`
	DataID string `name:"data_id"json:"dataId"`
	GroupID string `name:"group_id"json:"groupId"`
	Content string `name:"content"json:"content,omitempty"`
	NamespaceID int64 `name:"namespace_id"json:"namespaceId,omitempty"`
	MD5 string `name:"md5"json:"md5"`
	GmtCreate string `name:"gmt_create"json:"gmtCreate,omitempty"`
	GmtModified string `name:"gmt_modified"json:"gmtModified,omitempty"`
	SrcUser string `name:"src_user"json:"srcUser,omitempty"`
	SrcIP string `name:"src_ip"json:"srcIp,omitempty"`
	Type string `name:"type"json:"type"`

	TableName orm.TableName `name:"plugin_config"json:"-"`
}
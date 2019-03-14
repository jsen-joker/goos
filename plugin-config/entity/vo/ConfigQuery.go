package vo

type ConfigItem struct {
	MD5 string `name:"md5"json:"md5"`
	DataID string `name:"dataId"json:"dataId"`
}
type ConfigQuery struct {
	GroupID string `name:"groupId"json:"groupId"`
	NamespaceID string `name:"namespaceId"json:"namespaceId"`
	ConfigList []ConfigItem `name:"configList"json:"configList"`
}

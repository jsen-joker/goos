package vo

type AccountQuery struct {
	Username string `name:"username"json:"username"`
	Password string `name:"password"json:"password"`
	Type string `name:"type"json:"type"`
}
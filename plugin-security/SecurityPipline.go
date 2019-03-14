package plugin_security

import (
	"encoding/json"
	"github.com/jsen-joker/goos/core/utils"
	utils2 "github.com/jsen-joker/goos/plugin-security/utils"
	"net/http"
	"strings"
)

const TOKEN_PREFIX  = "Bearer "
const AUTHORIZATION_HEADER = "Authorization"
const AUTHORIZATION_TOKEN = "access_token"

// 安全验证
func getToken(r *http.Request) string {
	token := r.Header.Get(AUTHORIZATION_HEADER)
	if utils.NotEmpty(&token) && strings.Index(token, TOKEN_PREFIX) == 0 {
		return token[7 : ]
	}
	tokens := r.URL.Query()[AUTHORIZATION_TOKEN]

	if len(tokens) == 0 {
		return ""
	}

	return tokens[0]
}

// 这个pipline会在所有不忽略的白名单api路径上拦截检查token是否合法
func HttpPipeSecurity(inner http.Handler, name string) http.HandlerFunc  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwt := getToken(r)

		if err := utils2.ValidToken(jwt); err == nil {
			inner.ServeHTTP(w, r)
		} else {
			resp := utils.Failed("Unauthorized").Put("status", 401)
			w.WriteHeader(401)
			if err := json.NewEncoder(w).Encode(resp); err != nil{
				panic(err)
			}
		}
	})
}

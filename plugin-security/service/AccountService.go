package service

import (
	"github.com/jsen-joker/goos/core/support-db"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/core/utils"
	"github.com/jsen-joker/goos/plugin-security/entity"
	"github.com/jsen-joker/goos/plugin-security/entity/vo"
	"github.com/pkg/errors"
)

func Login(account vo.AccountQuery) (result interface{}, err error) {
	acc, err := support_db.QueryOne((&orm.QueryWrapper{}).Entity(entity.Account{Username: account.Username}))

	if err != nil || acc == nil {
		return nil, errors.New("account not exist or some oops")
	} else {
		// check password
		if err := utils.Matches(acc.(entity.Account).Password, account.Password); err != nil {
			// password error
			return nil, errors.New("password error")
		} else {
			// password correct
			return acc, nil
		}
	}
}


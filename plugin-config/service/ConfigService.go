package service

import (
	"crypto/md5"
	"fmt"
	"github.com/jsen-joker/goos/core/support-db"
	"github.com/jsen-joker/goos/core/support-db/orm"
	"github.com/jsen-joker/goos/plugin-config/entity"
	"github.com/jsen-joker/goos/plugin-config/entity/vo"
	"time"
)

// 创建或者更新
func Config(config entity.Config) (eff int64, err error)  {

	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)


	config.GmtModified = tm.Format("2006-01-02 15:04:05")
	config.MD5 = fmt.Sprintf("%x", md5.Sum([]byte(config.Content)))

	//if data, e := alg.RsaEncrypt(config.Content); e != nil {
	//	return 0, e
	//} else {
	//	config.Content = data
	//}

	if config.Id == 0 {
		// do create
		config.GmtCreate = config.GmtModified

		return support_db.Insert(config)
	} else {
		// do update
		return support_db.Update(config)
	}

}

func GetConfigList(currentPage int64, pageSize int64, namespace int64, queryDataId string, queryGroup string) (list []interface{}, total int64, err error) {
	skip := (currentPage - 1) * pageSize
	q := (&orm.QueryWrapper{}).Entity(entity.Config{NamespaceID:namespace})
	if queryDataId != "" {
		q.Where("data_id", "%" + queryDataId + "%").Like()
	}
	if queryGroup != "" {
		q.Where("group_id", "%" + queryGroup + "%").Like()
	}

	total, e := support_db.Count(q)
	if e != nil {
		return nil, 0, e
	}
	q.Select("id", "data_id", "group_id", "md5", "type").Skip(skip).Capacity(pageSize)
	rs, e := support_db.Query(q)
	if e == nil && rs == nil {
		return []interface{}{}, total, nil
	}
	return rs, total, e
}

func GetConfig(id int64) (result interface{}, err error) {
	if result, err := support_db.QueryOne((&orm.QueryWrapper{}).Entity(entity.Config{Id:id})); err == nil {
		//r := result.(entity.Config)
		//if data, e := alg.RsaDecrypt(r.Content); e != nil {
		//	return nil, e
		//} else {
		//	r.Content = data
		//	return r, nil
		//}
		return result, nil
	} else {
		return result, err
	}
}

func DeleteConfig(id int64) (eff int64, err error) {
	if eff, err := support_db.Delete(entity.Config{Id: id}); err == nil {
		return eff, err
	} else {
		return eff, err
	}
}


///////////////////////////////////   RSA CLIENT API      ///////////////////////////////

func RsaGetConfig(query vo.ConfigQuery) (result []interface{}, err error) {
	if result, err := support_db.QueryOne((&orm.QueryWrapper{}).Entity(entity.Namespace{NamespaceID: query.NamespaceID})); err == nil {
		namespaceId := result.(entity.Namespace).Id

		configList := query.ConfigList

		var params []interface{}
		sql := "select * from plugin_config where namespace_id=? and group_id=? and("

		params = append(params, namespaceId, query.GroupID)


		for _, config := range configList  {
			sql += "(md5!=? and data_id=?) or"
			params = append(params, config.MD5, config.DataID)
		}
		sql = sql[ : len(sql) - 3]
		sql += ")"

		fmt.Println(sql)
		fmt.Println(params)

		myrows, err := support_db.SqlQuery(sql, params...)
		if err != nil{
			return nil, err
		}
		list, err :=  myrows.To((&orm.QueryWrapper{}).Entity(entity.Config{}).GetFrom())
		if err != nil{
			return nil, err
		}

		if list == nil {
			list = []interface{} {}
		}
		return list, err
	} else {
		return nil, err
	}

}
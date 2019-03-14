package orm

import (
	"reflect"
	"strconv"
	"strings"
)

/*
  参数字段
*/
type ParamField struct{
	context *QueryWrapper
	name string
	value interface{}
	// and or
	operation string
	// > >= < <= = like !=
	compare string
	// 儿子
	children *[]*ParamField
	lastField *ParamField // 最后一项
}
func (p *ParamField) Eq() *ParamField {
	return p.cp("=")
}
func (p *ParamField) Gt() *ParamField {
	return p.cp(">")
}
func (p *ParamField) Gte() *ParamField {
	return p.cp(">=")
}
func (p *ParamField) Lt() *ParamField {
	return p.cp("<")
}
func (p *ParamField) Lte() *ParamField {
	return p.cp("<=")
}
func (p *ParamField) Ne() *ParamField {
	return p.cp("!=")
}
func (p *ParamField) Like() *ParamField {
	return p.cp("like")
}
func (p *ParamField) cp(o string) *ParamField {
	p.compare = o
	return p
}
func (p *ParamField) op(o string) *ParamField {
	p.operation = o
	return p
}

func (p *ParamField) And(name string, val interface{}) *ParamField {
	return p.condition(name, val, "and", "=")
}
func (p *ParamField) Or(name string, val interface{}) *ParamField {
	return p.condition(name, val, "or", "=")
}
func (p *ParamField) condition(name string, val interface{}, operation string, compare string) *ParamField {
	//founded := false
	//for _, item := range *p.children {
	//	if item.name == name {
	//		founded = true
	//		break
	//	}
	//}
	//if !founded {
	//	*p.children = append(*p.children, ParamField{name:name, value:val, operation: "or", children: &[]ParamField{}})
	//}
	p.lastField = &ParamField{context: p.context, name:name, value:val, operation: operation, children: &[]*ParamField{}, compare: compare}
	*p.children = append(*p.children, p.lastField)
	return p
}
//
//func (p *ParamField) Group() *ParamField {
//	return p.lastField
//}


func (p *ParamField) Build() *QueryWrapper {
	return p.context
}


type QueryWrapper struct {
	from string  //表名
	fields []string  //select的字段名
	where []*ParamField // where条件，暂时只支持and连接
	nowField *ParamField // 当前编辑的项
	values []interface{} //查询的参数值
	entity interface{}
	// limit
	skip int64
	capacity int64
}

//设置表名
func (q *QueryWrapper) From(from string) *QueryWrapper{
	q.from = from
	return q
}

func (q *QueryWrapper) Entity(entity interface{}) *QueryWrapper {

	q.entity = entity
	ti, err := getTableInfo(entity)

	if err == nil {
		q.From(ti.Name)


		rt := reflect.TypeOf(entity)
		rv := reflect.ValueOf(entity)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
			rv = rv.Elem()
		}

		for k, v := range ti.MapperXML {
			rV := rv.FieldByName(v).Interface()
			if rV != nil {
				tp := reflect.TypeOf(rV).String()
				if tp == "string" && rV == "" {
					continue
				} else if tp == "int" && rV.(int) == 0 {
					continue
				} else if tp == "int64" && rV.(int64) == 0 {
					continue
				}
				q.Where(k, rV)
			}
		}
	}
	return q
}

func (q *QueryWrapper) Skip(skip int64) *QueryWrapper {
	q.skip = skip
	return q
}

func (q *QueryWrapper) Capacity(capacity int64) *QueryWrapper {
	q.capacity = capacity
	return q
}

//设置查询条件
func (q *QueryWrapper) Where(name string, val interface{}) *QueryWrapper {
	return q.condition(name, val, "and", "=")
}
func (q *QueryWrapper) GroupAnd(name string, val interface{}) *ParamField {
	return q.gCondition(name, val, "and", "=")
}
func (q *QueryWrapper) GroupOr(name string, val interface{}) *ParamField {
	return q.gCondition(name, val, "or", "=")
}

//func (q *QueryWrapper) Group() *ParamField {
//	return q.nowField
//}

// operation
func (q *QueryWrapper) And() *QueryWrapper {
	return q.op("and")
}
func (q *QueryWrapper) Or() *QueryWrapper {
	return q.op("or")
}

// compare
func (q *QueryWrapper) Eq() *QueryWrapper {
	return q.cp("=")
}
func (q *QueryWrapper) Gt() *QueryWrapper {
	return q.cp(">")
}
func (q *QueryWrapper) Gte() *QueryWrapper {
	return q.cp(">=")
}
func (q *QueryWrapper) Lt() *QueryWrapper {
	return q.cp("<")
}
func (q *QueryWrapper) Lte() *QueryWrapper {
	return q.cp("<=")
}
func (q *QueryWrapper) Ne() *QueryWrapper {
	return q.cp("!=")
}
func (q *QueryWrapper) Like() *QueryWrapper {
	return q.cp("like")
}
func (q *QueryWrapper) cp(o string) *QueryWrapper {
	if q.nowField != nil {
		q.nowField.cp(o)
	}
	return q
}
func (q *QueryWrapper) op(o string) *QueryWrapper {
	if q.nowField != nil {
		q.nowField.op(o)
	}
	return q
}

func (q *QueryWrapper) condition(name string, val interface{}, operation string, compare string) *QueryWrapper {
	//founded := false
	//for _, item := range q.where {
	//	if item.name == name {
	//		founded = true
	//		break
	//	}
	//}
	//if !founded {
	//	q.where = append(q.where, ParamField{context: q, name:name, value:val, operation: operation, children: &[]ParamField{}})
	//}
	p := &ParamField{context: q, name:name, value:val, operation: operation, children: &[]*ParamField{}, compare:compare}
	p.lastField = p
	q.nowField = p
	q.where = append(q.where, p)
	return q
}
func (q *QueryWrapper) gCondition(name string, val interface{}, operation string, compare string) *ParamField {
	//for _, item := range q.where {
	//	if item.name == name {
	//		return &item
	//	}
	//}
	p := &ParamField{context: q, name:name, value:val, operation: operation, children: &[]*ParamField{}, compare:compare}
	q.nowField = p
	q.where = append(q.where, p)
	return p
}

func (q *QueryWrapper) Build() *QueryWrapper {
	return q
}

//设置select字段
func (q *QueryWrapper) Select(args ...string) *QueryWrapper{
	for _, v := range args{
		q.fields = append(q.fields, v)
	}
	return q
}

//获取参数
func (q QueryWrapper)GetValues()[]interface{}{
	return q.values
}
//获取参数
func (q QueryWrapper)GetFrom()string{
	return q.from
}

func (q QueryWrapper)getPk(defaultPk string) string {
	tbInfo, err := getTableInfo(q.entity)

	if err == nil {
		for _, v := range tbInfo.Fields {
			if v.IsPrimaryKey {
				return v.Name
			}
		}
	}
	return defaultPk
}

// 处理pf的循环嵌套问题
func (q *QueryWrapper) handleGroup(v *ParamField) (condition string, values []interface{}) {
	c := ""
	cp := v.compare
	if cp == "" {
		cp = "="
	}
	var vals []interface{}
	vals = append(vals, v.value)
	if len(*v.children) == 0 {
		c += v.operation + " `" + v.name + "` " + cp + " ? "
	} else {
		// group
		c += v.operation + " (`" + v.name + "` " + cp + " ? "

		for _, vi := range *v.children {
			ci, valsi := q.handleGroup(vi)
			c += ci
			vals = append(vals, valsi...)
		}
		c += ") "
	}

	return c, vals
}

func (q *QueryWrapper) withSelectTemplate(out string) string {
	q.values = q.values[0:0]
	strSql := "select " + out + " from " + q.from
	if q.where != nil && len(q.where) > 0{
		strSql += " where "
		condition := ""
		for _, v := range q.where{
			c, vals := q.handleGroup(v)
			condition += c
			q.values = append(q.values, vals...)
		}
		condition = strings.TrimLeft(strings.TrimLeft(condition, " and "), " or ")
		strSql += condition
	}

	if q.capacity > 0 {
		strSql += " limit " + strconv.FormatInt(q.skip, 10) + "," + strconv.FormatInt(q.capacity, 10)
	}
	return strSql
}

//获取查询的sql语句
func (q *QueryWrapper) GetSelectSql() string {
	strSql := ""
	if q.fields == nil || len(q.fields) < 1{
		strSql += "*"
	}else{
		for _, v := range q.fields{
			strSql += "`" + v + "`,"
		}
		strSql = strings.TrimRight(strSql, ",")
	}

	return q.withSelectTemplate(strSql)
}
func (q *QueryWrapper) GetCountSql() string {
	return q.withSelectTemplate("count(" + q.getPk("*") + ") total")
}
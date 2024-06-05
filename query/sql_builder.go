package query

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/errors"
)

type Option struct {
	DisableLimit bool `form:"disableLimit"`
	IsActive     bool
	IsInactive   bool
}

type BuildQueryOption struct {
	isLike        bool
	isMany        bool
	isSQLNull     bool
	paramTagValue string
	dbTagValue    string
	fieldValue    any
}

type sqlBuilder struct {
	dbTag         string
	paramTag      string
	rawQuery      *bytes.Buffer
	disableLimit  bool
	param         reflect.Value
	fieldValues   []any
	sortValue     []string
	pageValue     int64
	limitValue    int64
	mapDBTagExist map[string]bool
	option        *Option
}

func NewSQLQueryBuilder(paramTag, dbTag string, option *Option) *sqlBuilder {
	qb := sqlBuilder{
		rawQuery:      bytes.NewBufferString(" WHERE 1=1"),
		dbTag:         dbTag,
		paramTag:      paramTag,
		mapDBTagExist: map[string]bool{},
		option:        option,
	}

	if option != nil {
		if option.DisableLimit {
			qb.disableLimit = true
		}
		if option.IsActive {
			_, _ = qb.rawQuery.WriteString(" AND status=1")
		}
		if option.IsInactive {
			_, _ = qb.rawQuery.WriteString(" AND status=-1")
		}
	}

	return &qb
}

func (s *sqlBuilder) Build(param interface{}) (string, []interface{}, string, []interface{}, error) {
	var (
		newQuery      string
		newArgs       []interface{}
		newCountQuery string
		newCountArgs  []interface{}
	)

	paramReflectVal := reflect.ValueOf(param)
	if paramReflectVal.Kind() != reflect.Ptr || paramReflectVal.IsNil() {
		return newQuery, newArgs, newCountQuery, newCountArgs, errors.NewWithCode(codes.CodeInvalidValue, "passed param should be a pointer and cannot be nil")
	}

	s.param = paramReflectVal

	s.processParam(paramReflectVal, "", "")

	countQuery := s.rawQuery.Bytes()

	s.processSort()

	s.processPagination()

	newQuery, newArgs, err := sqlx.In(s.rawQuery.String()+";", s.fieldValues...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newQuery = sqlx.Rebind(sqlx.QUESTION, newQuery)

	newCountQuery, newCountArgs, err = sqlx.In(string(countQuery)+";", s.fieldValues...)
	if err != nil {
		return "", nil, "", nil, err
	}
	newCountQuery = sqlx.Rebind(sqlx.QUESTION, newCountQuery)

	s.restoreStruct()

	return newQuery, newArgs, newCountQuery, newCountArgs, nil
}

func (s *sqlBuilder) buildQuery(buildOption BuildQueryOption) {
	s.mapDBTagExist[buildOption.dbTagValue] = true

	if buildOption.fieldValue == nil {
		return
	}

	if isSortBy(buildOption.paramTagValue) {
		s.sortValue = buildOption.fieldValue.([]string)
		return
	}

	if isPage(buildOption.paramTagValue) {
		s.pageValue = validatePage(buildOption.fieldValue.(int64))
		return
	}

	if isLimit(buildOption.paramTagValue) {
		s.limitValue = validateLimit(buildOption.fieldValue.(int64))
		return
	}

	// write logical operator first
	if strings.Contains(buildOption.paramTagValue, "__opt") {
		s.rawQuery.WriteString(" OR")
	} else {
		s.rawQuery.WriteString(" AND")
	}

	// write condition clause if value is not slices
	if !buildOption.isMany {
		if buildOption.isLike {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + " LIKE " + s.getBindVar())
		} else if strings.Contains(buildOption.paramTagValue, "__gte") {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + ">=" + s.getBindVar())
		} else if strings.Contains(buildOption.paramTagValue, "__lte") {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + "<=" + s.getBindVar())
		} else if strings.Contains(buildOption.paramTagValue, "__lt") {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + "<" + s.getBindVar())
		} else if strings.Contains(buildOption.paramTagValue, "__gt") {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + ">" + s.getBindVar())
		} else if strings.Contains(buildOption.paramTagValue, "__ne") {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + "<>" + s.getBindVar())
		} else {
			s.rawQuery.WriteString(" " + buildOption.dbTagValue + "=" + s.getBindVar())
		}

		s.fieldValues = append(s.fieldValues, buildOption.fieldValue)
		return
	}

	// write condition clause if value is slices
	if strings.Contains(buildOption.paramTagValue, "__nin") {
		s.rawQuery.WriteString(" " + buildOption.dbTagValue + " NOT IN (" + s.getBindVar() + ")")
		s.fieldValues = append(s.fieldValues, buildOption.fieldValue)
		return
	}

	s.rawQuery.WriteString(" " + buildOption.dbTagValue + " IN (" + s.getBindVar() + ")")
	s.fieldValues = append(s.fieldValues, buildOption.fieldValue)
}

func (s *sqlBuilder) getBindVar() string {
	return "?"
}

func (s *sqlBuilder) restoreStruct() {
	*s = sqlBuilder{
		rawQuery:      bytes.NewBufferString(" WHERE 1=1"),
		dbTag:         s.dbTag,
		paramTag:      s.paramTag,
		mapDBTagExist: map[string]bool{},
		disableLimit:  s.disableLimit,
	}

	if s.option != nil {
		if s.option.DisableLimit {
			s.disableLimit = true
		}
		if s.option.IsActive {
			s.rawQuery.WriteString(" AND status=1")
		}
		if s.option.IsInactive {
			s.rawQuery.WriteString(" AND status=-1")
		}
	}
}

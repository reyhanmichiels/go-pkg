package query

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/reyhanmichiels/go-pkg/null"
)

/*
Determine specific action for each element in param
*/
func (s *sqlBuilder) processParam(param reflect.Value, paramTagValue string, dbTagValue string, isUpdate bool) {
	switch {

	case param.Kind() == reflect.Pointer || param.Kind() == reflect.Interface:
		if !param.IsValid() && param.IsNil() {
			return
		}

		s.processParam(param.Elem(), "", "", isUpdate)

	case param.Kind() == reflect.Struct && !isNullType(param) && !isTimeType(param):
		for i := 0; i < param.NumField(); i++ {
			if param.Field(i).CanSet() {
				paramTagValue := param.Type().Field(i).Tag.Get(s.paramTag)
				dbTagValue := param.Type().Field(i).Tag.Get(s.dbTag)

				if dbTagValue == "-" || dbTagValue == "" && param.Field(i).Kind() != reflect.Struct {
					continue
				}

				s.processParam(param.Field(i), paramTagValue, dbTagValue, isUpdate)
			}
		}

	default:
		s.processElem(param, paramTagValue, dbTagValue, isUpdate)
	}
}

// collect element build option and build the query
func (s *sqlBuilder) processElem(element reflect.Value, paramTagValue string, dbTagValue string, isUpdate bool) {
	buildOption := BuildQueryOption{
		paramTagValue: paramTagValue,
		dbTagValue:    dbTagValue,
	}

	if element.Kind() == reflect.String && strings.Contains(element.String(), "%") {
		buildOption.isLike = true
	}

	if element.Kind() == reflect.Slice {
		buildOption.isMany = true
	}

	buildOption = s.setBuildOption(element, buildOption)

	if isUpdate {
		s.buildQueryUpdate(buildOption)
		return
	}

	s.buildQuery(buildOption)
}

func (s *sqlBuilder) setBuildOption(element reflect.Value, buildOption BuildQueryOption) BuildQueryOption {
	switch convertedElement := element.Interface().(type) {
	case null.String:
		buildOption.isSQLNull = convertedElement.SqlNull
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.String
			if strings.Contains(convertedElement.String, "%") {
				buildOption.isLike = true
			}
		}
	case []null.String:
		var temp []string
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid && len(v.String) > 0 {
					temp = append(temp, v.String)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.String:
		var temp []string
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid && len(v.String) > 0 {
					temp = append(temp, v.String)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case null.Int64:
		buildOption.isSQLNull = convertedElement.SqlNull
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.Int64
		}
	case []null.Int64:
		var temp []int64
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid {
					temp = append(temp, v.Int64)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.Int64:
		var temp []int64
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid {
					temp = append(temp, v.Int64)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case null.Float64:
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.Float64
		}
	case []null.Float64:
		var temp []float64
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid {
					temp = append(temp, v.Float64)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.Float64:
		var temp []float64
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid {
					temp = append(temp, v.Float64)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case null.Bool:
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.Bool
		}
	case []null.Bool:
		var temp []bool
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid {
					temp = append(temp, v.Bool)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.Bool:
		var temp []bool
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid {
					temp = append(temp, v.Bool)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case null.Time:
		buildOption.isSQLNull = convertedElement.SqlNull
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.Time
		}
	case []null.Time:
		var temp []time.Time
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid {
					temp = append(temp, v.Time)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.Time:
		var temp []time.Time
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid {
					temp = append(temp, v.Time)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case null.Date:
		buildOption.isSQLNull = convertedElement.SqlNull
		if convertedElement.Valid {
			buildOption.fieldValue = convertedElement.Time
		}
	case []null.Date:
		var temp []time.Time
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v.Valid {
					temp = append(temp, v.Time)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	case []*null.Date:
		var temp []time.Time
		if len(convertedElement) > 0 {
			for _, v := range convertedElement {
				if v != nil && v.Valid {
					temp = append(temp, v.Time)
				}
			}
			buildOption.isMany = true
			buildOption.fieldValue = temp
		}
	default:
		if !element.IsZero() {
			buildOption.fieldValue = element.Interface()
		}
	}
	return buildOption
}

func (s *sqlBuilder) processSort() {
	sortValue := []string{}
	for _, v := range s.sortValue {
		sortOrder := "ASC"
		if !regexp.MustCompile(`(?P<sign>-)?(?P<col>[a-zA-Z_\\.0-9]+),?`).MatchString(v) {
			continue
		}

		if strings.Contains(v, "-") {
			sortOrder = "DESC"
			v = strings.Split(v, "-")[1]
		}

		if s.mapDBTagExist[v] {
			sortValue = append(sortValue, fmt.Sprintf("%v %v", v, sortOrder))
		}
	}

	if len(sortValue) > 0 {
		s.rawQuery.WriteString(" ORDER BY " + strings.Join(sortValue, ", "))
	}
}

func (s *sqlBuilder) processPagination() {
	if s.pageValue > 0 || s.limitValue > 0 {
		offset := getOffset(s.pageValue, s.limitValue)
		s.rawQuery.WriteString(fmt.Sprintf(" LIMIT %d, %d", offset, s.limitValue))
	}
}

func getOffset(p, l int64) int64 {
	if p > 0 {
		return (p - 1) * l
	}
	return 0
}

func isTimeType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct && (e.Type().String() == "null.Time" || e.Type().String() == "null.Date" || e.Type().String() == "time.Time")
}

func isNullType(e reflect.Value) bool {
	return e.Kind() == reflect.Struct &&
		(e.Type().String() == "null.String" ||
			e.Type().String() == "null.Bool" ||
			e.Type().String() == "null.Float64" ||
			e.Type().String() == "null.Int64" ||
			e.Type().String() == "null.Time" ||
			e.Type().String() == "null.Date")
}

func isPage(paramTagValue string) bool {
	return paramTagValue == "page"
}

func isLimit(paramTagValue string) bool {
	return paramTagValue == "limit"
}

func isSortBy(paramTagValue string) bool {
	return paramTagValue == "sort-by" || paramTagValue == "sort_by" || paramTagValue == "sortBy" || paramTagValue == "sortby"
}

func validateLimit(l int64) int64 {
	if l < 1 {
		return 10
	}
	return l
}

func validatePage(p int64) int64 {
	if p < 1 {
		return 1
	}
	return p
}

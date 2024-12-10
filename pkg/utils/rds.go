package golibutils

import (
	"reflect"
	"strings"
	"sync"
)

var fieldToColumnCache = sync.Map{}
var fieldToJsonBCache = sync.Map{}

// GetTableColumnFromModel returns the column name of the field in the rds model
func GetTableColumnFromModel(model interface{}, fieldName string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	// Handle slice of pointers to structs
	if modelType.Kind() == reflect.Slice && modelType.Elem().Kind() == reflect.Struct {
		modelType = modelType.Elem()
	}

	cacheKey := modelType.String() + "." + fieldName
	if cachedColumn, ok := fieldToColumnCache.Load(cacheKey); ok {
		return cachedColumn.(string)
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			columnName := GetTableColumnFromField(field)
			fieldToColumnCache.Store(cacheKey, columnName)
			return columnName
		}
	}

	return ""
}

func GetTableColumnFromField(field reflect.StructField) string {
	tag := field.Tag.Get("gorm")

	tagParts := strings.Split(tag, ";")

	for _, part := range tagParts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}

	return ""
}

func IsJsonBField(model interface{}, fieldName string) bool {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	// Handle slice of pointers to structs
	if modelType.Kind() == reflect.Slice && modelType.Elem().Kind() == reflect.Struct {
		modelType = modelType.Elem()
	}

	cacheKey := modelType.String() + "." + fieldName
	if cachedColumn, ok := fieldToJsonBCache.Load(cacheKey); ok {
		return cachedColumn.(bool)
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			if strings.Contains(field.Tag.Get("gorm"), "type:jsonb") {
				fieldToJsonBCache.Store(cacheKey, true)
				return true
			}
		}
	}

	return false
}

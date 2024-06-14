package core_database

import (
	"encoding/hex"
	"fmt"
	"reflect"
)

var (
	nilIntField      = 0
	nilStringField   = ""
	NIL_INT_FIELD    = CreateField(&nilIntField)
	NIL_STRING_FIELD = CreateField(&nilStringField)
)

func CreateField[T any](value T) *T {
	return &value
}

func CreateColumnsQuery(fieldsMap any) string {
	query := ""
	valueIterator := reflect.Indirect(reflect.ValueOf(fieldsMap))
	fieldsNum := valueIterator.NumField()
	typeIterator := valueIterator.Type()

	for i := 0; i < fieldsNum; i++ {
		var v string

		tag := typeIterator.Field(i).Tag.Get("db")
		kind := valueIterator.Field(i).Kind()

		if isNil := valueIterator.Field(i).IsNil(); !isNil {
			switch kind {
			case reflect.String:
				v = fmt.Sprintf("'%s'", valueIterator.Field(i).String())
			case reflect.Pointer:
				switch valueIterator.Field(i).Interface() {
				case NIL_INT_FIELD:
					v = "NULL"
				case NIL_STRING_FIELD:
					v = "NULL"
				default:
					rv := reflect.Indirect(valueIterator.Field(i).Elem())

					switch rv.Kind() {
					case reflect.String:
						v = fmt.Sprintf("'%s'", rv)
					case reflect.Bool:
						if rv.Bool() {
							v = "TRUE"
						} else {
							v = "FALSE"
						}
					case reflect.Slice:
						buf := []byte{}
						buf = append(buf, `\x`...)
						buf = append(buf, hex.EncodeToString(rv.Bytes())...)
						v = fmt.Sprintf("'%s'", buf)
					default:
						v = fmt.Sprintf("%v", rv)
					}
				}
			default:
				v = fmt.Sprintf("%v", valueIterator.Field(i).Interface())
			}

			if query == "" {
				query += fmt.Sprintf("%s = %v", tag, v)
			} else {
				query += fmt.Sprintf(", %s = %v", tag, v)
			}
		}
	}

	return query
}

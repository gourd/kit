package httpservice

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// EnforceCreate enforce some field on creation,
// as defined in struct field tags
func EnforceCreate(payload interface{}) (err error) {

	// find payload pointer
	ptr := reflect.ValueOf(payload)
	if !ptr.IsValid() {
		err = fmt.Errorf("reflect.ValueOf(payload) is of zero Value")
		return
	}
	if ptr.Type().Kind() != reflect.Ptr {
		err = fmt.Errorf("payload is not pointer but %s", ptr.Type())
		return
	}
	val := ptr.Elem()
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		err = fmt.Errorf("payload is not struct but %s", typ.Kind())
	}

	// pre-allocate for gourdcreate:"now"
	now := reflect.ValueOf(time.Now())
	emptyTime := time.Time{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		createTag := field.Tag.Get("gourdcreate")
		if createTag == "" {
			continue // skip all fields without gourdcreate tag
		}

		tagContents := strings.Split(createTag, ",")
		if tagContents[0] == "now" {
			fieldVal := val.Field(i)
			if field.Type.PkgPath() == "time" && field.Type.Name() == "Time" {
				currentValue := fieldVal.Addr().Interface().(*time.Time)
				if len(tagContents) > 1 && tagContents[1] == "omitnotempty" && (*currentValue) != emptyTime {
					continue // omit if not empty
				}
				fieldVal.Set(now)
			} else {
				err = fmt.Errorf("field %s is tagged with \"now\" but not of type time.Time", field.Name)
				return
			}
		}
	}
	return
}

// EnforceUpdate enforce some fields value on update,
// as defined in struct field tags
func EnforceUpdate(original, update interface{}) (err error) {

	// find update pointer, value and type
	ptr := reflect.ValueOf(update)
	if !ptr.IsValid() {
		err = fmt.Errorf("reflect.ValueOf(update) is of zero Value")
		return
	}
	if ptr.Type().Kind() != reflect.Ptr {
		err = fmt.Errorf("update is not pointer but %s", ptr.Type())
		return
	}
	val := ptr.Elem()
	typ := val.Type()

	// test if it is struct type
	if typ.Kind() != reflect.Struct {
		err = fmt.Errorf("update is not a struct, but %s", typ.Kind())
		return
	}

	// find original pointer, value and type
	optr := reflect.ValueOf(original)
	if ptr.Type().Kind() != reflect.Ptr {
		err = fmt.Errorf("original is not pointer but %s", ptr.Type())
		return
	}
	oval := optr.Elem()
	otyp := oval.Type()

	// test if original and update are of the same type
	if otyp != typ {
		err = fmt.Errorf("*original (%s) is not of same type of *update (%s)",
			otyp, typ)
		return
	}

	// pre-allocate for gourdcreate:"now"
	now := reflect.ValueOf(time.Now())
	emptyTime := time.Time{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		updateTag := field.Tag.Get("gourdupdate")
		if updateTag == "" {
			continue // skip all fields without gourdcreate tag
		}

		tagContents := strings.Split(updateTag, ",")
		if tagContents[0] == "preserve" {
			currentValue := oval.Field(i)
			val.Field(i).Set(currentValue)
		} else if tagContents[0] == "now" {
			fieldVal := val.Field(i)
			if field.Type.PkgPath() == "time" && field.Type.Name() == "Time" {
				currentValue := fieldVal.Addr().Interface().(*time.Time)
				if len(tagContents) > 1 && tagContents[1] == "omitnotempty" && (*currentValue) != emptyTime {
					continue // omit if not empty
				}
				fieldVal.Set(now)
			} else {
				err = fmt.Errorf("field %s is tagged with \"now\" but not of type time.Time", field.Name)
				return
			}
		}
	}
	return
}

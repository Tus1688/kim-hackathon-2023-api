package jsonutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// ShouldBind is used to bind the request body to the struct and check the required fields
func ShouldBind(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	return checkRequiredFields(v)
}

func checkRequiredFields(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get("binding")
		if tag == "required" && isZeroOfUnderlyingType(field.Interface()) {
			return fmt.Errorf("'%s' field is required", val.Type().Field(i).Name)
		}
	}
	return nil
}

func isZeroOfUnderlyingType(x any) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

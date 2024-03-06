package validatate

import (
	"errors"
	"fmt"
	"reflect"
)

//Ух ля, ну и надоело мне этот валидатор писать под тесты.

func ChecStructForNil(s interface{}) error {
	if s == nil {
		return errors.New("validator.ChecStructForNil: struct is nil")
	}
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct && v.Kind() != reflect.Ptr {
		return errors.New("validator.ChecStructForNil: argument must be a struct or pointer to a struct")
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return errors.New("validator: struct is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic caught:", r)
		}
	}()
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
			return errors.New(fmt.Sprintf("validator: field %s is nil", v.Field(i).Type().Name()))
		}
	}
	return nil
}

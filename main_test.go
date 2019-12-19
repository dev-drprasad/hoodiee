package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	d := getSiteDefinition("1337x")
	v := reflect.ValueOf(d.Search.List.Fields)
	typeOfDefinition := v.Type()

	for i := 0; i < v.NumField(); i++ {
		s := v.Field(i).Interface().(Selector)
		n := typeOfDefinition.Field(i).Name
		fmt.Printf("%v", s)
		fmt.Printf("%s", n)
	}
}

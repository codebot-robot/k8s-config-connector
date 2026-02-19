package main

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/bigtable"
)

func main() {
	var conf bigtable.UpdateSchemaBundleConf
	fmt.Println("UpdateSchemaBundleConf fields:")
	printFields(reflect.TypeOf(conf))
}

func printFields(t reflect.Type) {
	if t.Kind() != reflect.Struct {
		fmt.Printf("Not a struct: %s\n", t.Kind())
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("  %s: %s\n", field.Name, field.Type)
	}
}

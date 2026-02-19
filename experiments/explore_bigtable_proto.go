package main

import (
	"fmt"
	"reflect"

	pb "cloud.google.com/go/bigtable/admin/apiv2/adminpb"
)

func main() {
	var bundle pb.SchemaBundle
	fmt.Println("SchemaBundle fields:")
	printFields(reflect.TypeOf(bundle))
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

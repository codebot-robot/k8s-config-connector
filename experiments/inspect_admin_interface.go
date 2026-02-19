package main

import (
	"fmt"
	"reflect"

	pb "cloud.google.com/go/bigtable/admin/apiv2/adminpb"
)

func main() {
	iface := reflect.TypeOf((*pb.BigtableTableAdminServer)(nil)).Elem()
	fmt.Println("BigtableTableAdminServer methods:")
	for i := 0; i < iface.NumMethod(); i++ {
		method := iface.Method(i)
		fmt.Printf("  %s\n", method.Name)
	}
}

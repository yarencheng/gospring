package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yarencheng/gospring"
)

func main() {

	//
	// Configuration of beans
	//

	type MyStruct struct {
		I int
		S string
	}

	beans := gospring.Beans(
		gospring.Bean(MyStruct{}).
			Id("a unique of this bean").
			PropertyValue("I", 12345).
			PropertyValue("S", "a string value"),
	)

	//
	// Create a application contex
	//

	ctx, e1 := gospring.ApplicationContext(beans)

	if e1 != nil {
		fmt.Println("Create application context failed. Cuased by: ", e1)
		os.Exit(1)
	}

	// Use the context to get bean

	myStruct, e2 := ctx.GetBean("a unique of this bean")
	if e2 != nil {
		fmt.Println("Can't get bean. Cuased by: ", e2)
		os.Exit(1)
	}

	print("myStruct", myStruct)
}

func print(name string, i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name, " = ", string(b))
}

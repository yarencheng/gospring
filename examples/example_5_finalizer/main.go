package main

import (
	"encoding/json"
	"fmt"

	"github.com/yarencheng/gospring"
)

type MyStruct struct{}

func (m *MyStruct) Finalize() {
	fmt.Println("Finalized by (*MyStruct).Finalize()")
}

func (m *MyStruct) CostumeFunction() {
	fmt.Println("Finalized by (*MyStruct).CostumeFunction()")
}

func main() {

	fmt.Println("Finalizer can be (1) a member function named \"Finalize\" or a costumized function")

	beans := gospring.Beans(

		gospring.Bean(MyStruct{}).
			Id("id_1"),

		gospring.Bean(MyStruct{}).
			Id("id_2").
			Finalize("CostumeFunction"),
	)

	ctx, _ := gospring.ApplicationContext(beans)

	ctx.GetBean("id_1")
	ctx.GetBean("id_2")

	ctx.Finalize()
}

func print(name string, i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name, " = ", string(b))
}

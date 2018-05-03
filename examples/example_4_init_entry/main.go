package main

import (
	"encoding/json"
	"fmt"

	"github.com/yarencheng/gospring"
)

type MyStruct struct{ MyValue string }

func (m *MyStruct) Init() {
	m.MyValue = "Init by (*MyStruct).Init()"
}

func (m *MyStruct) CostumeFunction() {
	m.MyValue = "Init by (*MyStruct).CostumeFunction()"
}

func main() {

	fmt.Println("Initailizer can be (1) a member function named \"Init\" or a costumized function")

	beans := gospring.Beans(

		gospring.Bean(MyStruct{}).
			Id("id_1"),

		gospring.Bean(MyStruct{}).
			Id("id_2").
			Init("CostumeFunction"),
	)

	ctx, _ := gospring.ApplicationContext(beans)

	bean_id_1, _ := ctx.GetBean("id_1")
	bean_id_2, _ := ctx.GetBean("id_2")

	print("bean_id_1", bean_id_1)
	print("bean_id_2", bean_id_2)
}

func print(name string, i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name, " = ", string(b))
}

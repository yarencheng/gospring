package main

import (
	"encoding/json"
	"fmt"

	"github.com/yarencheng/gospring"
)

func main() {

	type ChildStruct1 struct{ MyValue string }
	type ChildStruct2 struct{ MyValue string }
	type ChildStruct3 struct{ MyValue string }
	type ParentStruct struct {
		Child1       ChildStruct1
		Child1_p     *ChildStruct1
		Child2       ChildStruct2
		Child2_p     *ChildStruct2
		Child3_local ChildStruct3
	}

	fmt.Println("The order of beans defination dose not matter.")

	beans := gospring.Beans(

		gospring.Bean(ChildStruct1{}).
			Id("child1_id").
			PropertyValue("MyValue", "it is child 1"),

		gospring.Bean(ParentStruct{}).
			Id("parent_id").
			PropertyRef("Child1", "child1_id").
			PropertyRef("Child1_p", "child1_id").
			PropertyRef("Child2", "child2_id").
			PropertyRef("Child2_p", "child2_id").
			PropertyBean("Child3_local",
				gospring.Bean(ChildStruct3{}).
					PropertyValue("MyValue", "it is child 3")),

		gospring.Bean(ChildStruct2{}).
			Id("child2_id").
			PropertyValue("MyValue", "it is child 2"),
	)

	ctx, _ := gospring.ApplicationContext(beans)

	parent, _ := ctx.GetBean("parent_id")

	print(parent)
}

func print(i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

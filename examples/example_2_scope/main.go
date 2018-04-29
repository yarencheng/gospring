package main

import (
	"encoding/json"
	"fmt"

	"github.com/yarencheng/gospring"
)

func main() {

	type MyStruct struct {
		MyValue string
	}

	beans := gospring.Beans(
		gospring.Bean(MyStruct{}).Id("default_id").PropertyValue("MyValue", "AAA"),
		gospring.Bean(MyStruct{}).Id("sigleton_id").Singleton().PropertyValue("MyValue", "BBB"),
		gospring.Bean(MyStruct{}).Id("prototype_id").Prototype().PropertyValue("MyValue", "BBB"),
	)

	ctx, _ := gospring.ApplicationContext(beans)

	{
		//
		// Singleton scope
		//

		bean1, _ := ctx.GetBean("sigleton_id")
		bean2, _ := ctx.GetBean("sigleton_id")

		fmt.Printf("\n___ Singleton ___\n\n")
		fmt.Println("=== bean1 and bean2 are pointers and point to a same bean")

		print("bean1", bean1)
		print("bean2", bean2)

		fmt.Println("=== After modify bean1.MyValue, bean2.MyValue is changed as well")
		bean1.(*MyStruct).MyValue = "CCC"
		print("bean1", bean1)
		print("bean2", bean2)
	}

	{
		//
		// Prototype scope
		//

		bean1, _ := ctx.GetBean("prototype_id")
		bean2, _ := ctx.GetBean("prototype_id")

		fmt.Printf("\n___ Prototype ___\n\n")
		fmt.Println("=== bean1 and bean2 are pointers and point to diffrent beans")

		print("bean1", bean1)
		print("bean2", bean2)

		fmt.Println("=== After modify bean1.MyValue, bean2.MyValue remains unchanged")
		bean1.(*MyStruct).MyValue = "CCC"
		print("bean1", bean1)
		print("bean2", bean2)
	}

	{
		//
		// Default scope
		//

		bean1, _ := ctx.GetBean("default_id")
		bean2, _ := ctx.GetBean("default_id")

		fmt.Printf("\n___ Default ___\n\n")
		fmt.Println("=== The default scope of bean is singleton")
		fmt.Println("=== bean1 and bean2 are pointers and point to a same bean")

		print("bean1", bean1)
		print("bean2", bean2)

		fmt.Println("=== After modify bean1.MyValue, bean2.MyValue is changed as well")
		bean1.(*MyStruct).MyValue = "CCC"
		print("bean1", bean1)
		print("bean2", bean2)
	}
}

func print(name string, i interface{}) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(name, " = ", string(b))
}

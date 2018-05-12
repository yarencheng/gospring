package main

import (
	"encoding/json"
	"fmt"
	"os"

	gs "github.com/yarencheng/gospring"
)

type Astruct struct {
	// native type
	IntValue    int
	StringValue string

	// struct type
	Asingleton *Bstruct
	Aprototype *Cstruct
}

type Bstruct struct {
	Name string
}
type Cstruct struct {
	Name string
}

func main() {

	beans := gs.Beans(
		gs.Bean(Astruct{}).
			ID("a_id").
			Property(
				"IntValue",
				123,
			).
			Property(
				"StringValue",
				"a string",
			).
			Property(
				"Asingleton",
				gs.Ref("b_id"),
			).
			Property(
				"Aprototype",
				gs.Ref("c_id"),
			),
		gs.Bean(Bstruct{}).
			ID("b_id").
			Property(
				"Name",
				"It's B.",
			),
		gs.Bean(Cstruct{}).
			ID("c_id").
			Property(
				"Name",
				"It's C.",
			),
	)

	ctx, e := gs.NewApplicationContext(beans...)

	if e != nil {
		fmt.Printf("Something goes wrong while parsing bean definitions. Caused by: %v\n", e)
		os.Exit(1)
	}

	bean, e := ctx.GetBean("a_id")

	if e != nil {
		fmt.Printf("Something goes wrong while getting a bean with ID a_id. Caused by: %v\n", e)
		os.Exit(1)
	}

	print(bean)
}

func print(v interface{}) {
	json, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(json))
}

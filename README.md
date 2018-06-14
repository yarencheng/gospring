# gospring [![Build Status](https://travis-ci.org/yarencheng/gospring.svg?branch=develop-v2.0)](https://travis-ci.org/yarencheng/gospring) [![codecov](https://codecov.io/gh/yarencheng/gospring/branch/develop-v2.0/graph/badge.svg)](https://codecov.io/gh/yarencheng/gospring/branch/develop-v2.0) [![Go Report Card](https://goreportcard.com/badge/github.com/yarencheng/gospring)](https://goreportcard.com/report/github.com/yarencheng/gospring)
A Go injection library inspired by Java Spring Framework.

Java Spring Framework is a famous IoC framework. This library try to imitate spring's function.

## Quick start

1. Define your go structs

    ```go
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
    ```
2. Define the beans
    ```go
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
    ```
3. Put it into a context
    ```go
    ctx, e := gs.NewApplicationContext(beans...)
    ```

4. Create a bean from context
    ```go
    bean, e := ctx.GetBean("a_id")
    ```
5. Use the bean
    ```go
    json, _ := json.MarshalIndent(bean, "", "    ")
    fmt.Println(string(json))
    ```
    ```json
    {
        "IntValue": 123,
        "StringValue": "a string",
        "Asingleton": {
            "Name": "It's B."
        },
        "Aprototype": {
            "Name": "It's C."
        }
    }
    ```
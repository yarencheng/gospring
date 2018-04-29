# gospring [![Build Status](https://travis-ci.org/yarencheng/gospring.svg?branch=master)](https://travis-ci.org/yarencheng/gospring)
A injection Go library inspired by Java Spring Framework.

Java Spring Framework is a famous IoC framework. This library try to imitate spring's function. Every thing is still buggy.

# Example

## Basic


1. Configuration of beans

    ```go
    type MyStruct struct {
		I int
		S string
	}

	beans := config.Beans(
		config.Bean(MyStruct{}).
			Id("a unique of this bean").
			PropertyValue("I", 12345).
			PropertyValue("S", "a string value"),
	)
    ```
2. Create a application contex
    ```go
    ctx, e1 := config.ApplicationContext(beans)
    ```
3. Use the context to get bean
    ```go
	myStruct, e2 := ctx.GetBean("a unique of this bean")
    ```
    ```json
    myStruct  =  {
        "I": 12345,
        "S": "a string value"
    }
    ```

## Scope
```go
type MyStruct struct {
    MyValue string
}

beans := config.Beans(
    config.Bean(MyStruct{}).
        Id("default_id").
        PropertyValue("MyValue", "AAA"),
    config.Bean(MyStruct{}).
        Id("sigleton_id").
        Singleton().
        PropertyValue("MyValue", "BBB"),
    config.Bean(MyStruct{}).
        Id("prototype_id").
        Prototype().
        PropertyValue("MyValue", "BBB"),
)
```

1. Singleton scope

    ```go
    bean1, _ := ctx.GetBean("sigleton_id")
    bean2, _ := ctx.GetBean("sigleton_id")
    ```
    
    bean1 and bean2 are pointers and point to a same bean

    ```
    bean1  =  {
        "MyValue": "BBB"
    }
    bean2  =  {
        "MyValue": "BBB"
    }
    ```

    After modify bean1.MyValue, bean2.MyValue is changed as well

    ```go
    bean1.(*MyStruct).MyValue = "CCC"
    ```
    ```json
    bean1  =  {
        "MyValue": "CCC"
    }
    bean2  =  {
        "MyValue": "CCC"
    }
    ```

2. Prototype scope
3. 
    ```go
    bean1, _ := ctx.GetBean("prototype_id")
    bean2, _ := ctx.GetBean("prototype_id")
    ```

    bean1 and bean2 are pointers and point to diffrent beansAfter modify bean1.MyValue, bean2.MyValue remains unchanged

    ```go
    bean1.(*MyStruct).MyValue = "CCC"
    ```

    ```json
    bean1  =  {
        "MyValue": "CCC"
    }
    bean2  =  {
        "MyValue": "BBB"
    }
    ```

3. Default scope

    The default scope of bean is singleton

    ```go
    bean1, _ := ctx.GetBean("default_id")
    bean2, _ := ctx.GetBean("default_id")
    ```

    bean1 and bean2 are pointers and point to a same bean
    After modify bean1.MyValue, bean2.MyValue is changed as well

    ```go
    bean1.(*MyStruct).MyValue = "CCC"
    ```

    ```json
    bean1  =  {
        "MyValue": "CCC"
    }
    bean2  =  {
        "MyValue": "CCC"
    }
    ```

## Relation

```go
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
```

* The order of beans defination dose not matter.
* A bean with a ID can be referenced from other beans.
* A bean defined inside another bean called loacl bean.

```go
beans := config.Beans(

config.Bean(ChildStruct1{}).
    Id("child1_id").
    PropertyValue("MyValue", "it is child 1"),

config.Bean(ParentStruct{}).
    Id("parent_id").
    PropertyRef("Child1", "child1_id").
    PropertyRef("Child1_p", "child1_id").
    PropertyRef("Child2", "child2_id").
    PropertyRef("Child2_p", "child2_id").
    PropertyBean("Child3_local",
        config.Bean(ChildStruct3{}).
            PropertyValue("MyValue", "it is child 3")),

config.Bean(ChildStruct2{}).
    Id("child2_id").
    PropertyValue("MyValue", "it is child 2"),
)

ctx, _ := config.ApplicationContext(beans)

parent, _ := ctx.GetBean("parent_id")
```

```json
parent = {
  "Child1": {
    "MyValue": "it is child 1"
  },
  "Child1_p": {
    "MyValue": "it is child 1"
  },
  "Child2": {
    "MyValue": "it is child 2"
  },
  "Child2_p": {
    "MyValue": "it is child 2"
  },
  "Child3_local": {
    "MyValue": "it is child 3"
  }
}
```
# Bean

Index

* [Create a simple bean](#create-a-simple-bean)
* [ID of beans](#ID)
* [Scope](#scope)

## aaaaa

## Create a simple bean

* Define your object

    ```go
    type MyObject struct{
        ...
    }
    ```

* Define a bean with type ```MyObject``` and ID ```my_object```

    ```go
    b := Bean(MyObject{}).ID("my_object")
    ```

* Create an application context with the definition

    ```go
    ctx, e := NewApplicationContext(b)
    ```

* Get the bean instance from the application context by its ID

    ```go
    bean, e := ctx.GetBean("my_object")
    myObject, ok := bean.(*MyObject)
    ```

* Create mutiple beans at one time

    ```go
    bs := Beans(
        Bean(...).ID("id_1"),
        Bean(...).ID("id_2"),
        ...
    )
    ctx, e := NewApplicationContext(bs)
    ```

## ID

Bean can be difined with an ID or not. If a bean have an ID, it can be acquired by ```ApplicationContextI.GetBean("a_unique_id")```, or it can't. Please make sure each ID of beans is unique. A bean without ID is possible and is called local bean.

* Good

    ```go
    Beans(
        Bean(...).ID("id_1"),
        Bean(...).ID("id_2"),
        ...
    )
    ```

* Bad

    ```go
    Beans(
        Bean(...).ID("id_1"),
        Bean(...).ID("id_1"),
        ...
    )
    ```

## Scope

There are 2 type of scopes:

1. Singleton

    A instance is singleton means that the instance is shared globally. A singleton bean is created at the first time when some one try to acquired it. After the first time, the ```ApplicationContectI``` always return the same instance which is already created. See below example, ```bean1``` and ```bean2``` are two pointer point to the same instance.

    ```go
    type MyObject struct {}
    bs := Beans(
        Bean(MyObject{}).
            ID("id_1").
            Singleton(),
    )

    ctx, e := NewApplicationContext(bs)

    // bean1 == bean2
    bean1, _ := ctx.GetBean("id_1")
    bean2, _ := ctx.GetBean("id_1")
    ```

2. Prototype

    A prototype bean is contrary to a singleton bean. It always be created when someone try to acquired it. See below example, ```bean1``` and ```bean2``` are two pointer point to the different instance.

    ```go
    type MyObject struct {}
    bs := Beans(
        Bean(MyObject{}).
            ID("id_1").
            Prototype(),
    )

    ctx, e := NewApplicationContext(bs)

    // bean1 != bean2
    bean1, _ := ctx.GetBean("id_1")
    bean2, _ := ctx.GetBean("id_1")
    ```

## Type

A minimal configration of a bean is its type. Type could be a native type, a struct, or a slice. A pointer type dose not allow.

* Native type - OK

    ```go
    Bean(int(123))
    Bean(int64(456))
    Bean(float32(123.456))
    Bean(string("hello gospring"))
    ```
* ```struct``` - OK

    ```go
    type Foo struct { ... }
    Bean(Foo{})
    ```

    Be noticed that function ```Bean(type insterface{})``` is convenient comvention for gospring to know which **type of a struct** should be created. gospring dose not **copy** the value, so don't initialize the struct by ```Bean(type insterface{})```. Please use the ```(*BeanI) Factory(fninterface{})``` to do the initialization.

    ```go
    type Foo struct { Value int }

    ctx, _ := NewApplicationContext(
        Bean(Foo{ Value: 123 }).
            ID("foo"),
    )

    bean, _ := ctx.GetBean("foo")
    foo := bean.(*Foo)

    if foo.Value == 0 { // true
        ...
    }
    ```

* slice - OK

    ```go
    Bean([]int)
    Bean([]string)

    type Foo struct { ... }
    Bean([]Foo)
    ```

* ```map``` - not support
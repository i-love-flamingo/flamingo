# Dingo

Dependency injection for go


## Examples:

### standard struct injection

Dingo can inject structs automatically.
Its prefered to inject pointer to struct types:

```

type (
  MyStruct1 struct {
    // may have other dependencies that should be injected here
  }

  
  MyService struct {
    MyStruct1  *MyStruct1 `inject:""`
  }
)

```

### binding and injection of Interfaces

A very common use case is to the binding of interfaces types. 
(e.g. for ports and adapters)


```

type (
  MyInterface interface {
   MyDefinedBehaviour()
  }
  
  MyFakeImplementation struct {}
  
  MyService struct {
    MyInterface  MyInterface `inject:""`
  }
)

var (
 //Check interface
 var _ MyInterface = &MyFakeImplementation{}
)

func (m *MyFakeImplementation) MyDefinedBehaviour() {}


```

Bind in your module the implementation of the interface that Dingo should inject:

```
injector.Bind((*MyInterface)(nil)).To(MyFakeImplementation{})
```


### Using of Provider
Using of Provider for lazy binding or if you need new instances on demand

```
MyStructProvider func() *Service
MyStruct         struct {}

MyService struct {
  MyStructProvider MyStructProvider `inject:""`
}

```


### Usage of annotations

Default annotation is empty
Use .AnnotatedWith("my-annotation") Inject via
```
MyIface Iface `inject:"my-annotation"`
```
Used for example by configuration for configuration values

```
MyConfig string `inject:"config:myconfig,optional"`
```


### Bind multi
Allows to bind multiple instances/providers/types for one type
Injected by requesting a slice:

```
MyService struct {
  Ifaces []Iface `inject:""`
}
```


In module / registration:
```
injector.BindMulti((*Iface)(nil)).To(IfaceImpl{})
injector.BindMulti((*Iface)(nil)).To(IfaceImpl2{})
```


### Bind maps

Similiar to Multibindings, but with a key instead of a list
```
MyService struct {
  Ifaces map[string]Iface `inject:""`
}
```

```
injector.BindMap("impl1", (*Iface)(nil)).To(IfaceImpl{}) injector.BindMap("impl2", (*Iface)(nil)).To(IfaceImpl2{})

```

### Singleton scope

If really necessary it is possible to use singletons
``` 
.AsEagerSingleton() binds as a singleton, and loads it when the application is initialized
.In(dingo.Singleton) makes it a global singleton
.In(dingo.ChildSingleton) makes it a singleton limited to the config area
```

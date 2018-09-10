# Dingo

Dependency injection for go


## Hello Dingo

Dingo works very very similiar to [Guice](github.com/google/guice/wiki/GettingStarted)

Basically one binds implementations/factories to interfaces, which are then resolved by Dingo.

Given that Dingo's idea is based on Guice we use similar examples in this documentation:

The following example shows a BillingService with two injected dependencies. Please note
that Go's nature does not allow contructors, and does not allow decorations/annotations
beside struct-tags, thus, we only use struct tags (and later arguments for providers).

Also Go does not have a way to reference types (like Java's `Something.class`) we use `nil`
and cast it to a pointer to the interface we want to specify: `(*Something)(nil)`.
Dingo then knowns how to dereference it properly and derive the correct type `Something`.
This is not necessary for structs, where we can just use the null value via `Something{}`.

```go
type BillingService struct {
	Processor CreditCardProcessor `inject:""`
	TransactionLog TransactionLog `inject:""`
}

func (billingservice *BillingService) ChargeOrder(order PizzaOrder, creditCard CreditCard) Receipt {
	// ...
}
```

We want the BillingService to get certain dependencies, and configure this in a `BillingModule`
which implements `dingo.Module`:

```go
type BillingModule struct {}

func (module *BillingModule) Configure(injector *dingo.Injector) {
     /*
      * This tells Dingo that whenever it sees a dependency on a TransactionLog,
      * it should satisfy the dependency using a DatabaseTransactionLog.
      */
    injector.Bind((*TransactionLog)(nil)).To(DatabaseTransactionLog{})

     /*
      * Similarly, this binding tells Dingo that when CreditCardProcessor is used in
      * a dependency, that should be satisfied with a PaypalCreditCardProcessor.
      */
    injector.Bind((*CreditCardProcessor)(nil)).To(PaypalCreditCardProcessor{})
  }
}
```


## Requesting injection

Every instance that is created through the container can use injection. 

Dingo supports two ways of requesting dependencies that should be injected:

* usage of struct tags to allow structs to request injection into fields. This should be used for public fields.
* implement a public Inject(...) method to request injections of private fields. Dingo calls this method automatically and passes the requested injections.

For every requested injection (unless an exception applies) Dingo does the following:

- Is there a binding? If so: delegate to the binding
    - Is the binding in a certain scope (Singleton)? If so, delegate to scope (might result in a new loop)
    - Binding is bound to an instance: inject instance
    - Binding is bound to a provider: call provider
    - Binding is bound to a type: request injection of this type (might return in a new loop to resolve the binding)
- No binding? Try to create (only possible for concrete types, not interfaces or functions)


*Example:*
Here is another example using the Inject method for private fields
```go
type MyBillingService struct {
	processor CreditCardProcessor `inject:""`
	accountId string
}

func (m *MyBillingService) Inject(
	processor CreditCardProcessor,
	config *struct {
		AccountId  string `inject:"config:myModule.myBillingService.accountId"`
	},
) {
   m.processor = CreditCardProcessor
   m.accountId = config.AccountId
   
```

### Usage of Providers

Dingo allows to request the injection of provider instead of instances.
A "Provider" for dingo is a function that return a new Instance of a certain type.

```go
struct {
    PizzaProvider func() Pizza `inject:""`
}
```

If there is no concrete binding to the type `func() Pizza`, then instead of constructing one `Pizza` instance
Dingo will create a new function which, on every call, will return a new instance of `Pizza`.

The type must be of `func() T`, a function without any arguments which returns a type, which again has a binding.

This allows to lazily create new objects whenever needed, instead of requesting the Dingo injector itself.


You can use Providers and call them to always get a new instance.
Dingo will provide you with an automatic implementation of a Provider if you did not bind a specific one.

*Use a Provider instead of requesting the Type directly when*:

* for lazy binding 
* if you need new instances on demand
* In general it is best practice to use a Provider for everything that has a state that might be changed. This way you will avoid undesired side effects. That is especially important for dependencies in objects that are shared between requests - for example a controller!

*Example 1:*
This is the only code required to request a Provider as a dependency:
```
MyStructProvider func() *MyStruct
MyStruct         struct {}

MyService struct {
  MyStructProvider MyStructProvider `inject:""`
}

```

*Example 2:*
```go
func createSomething(factoryDependency SomethingElse) Something{
    return &MyType{}
}

injector.Bind((*Something)(nil)).ToProvider(createSomething)

struct {
    SomethingProvider func() Something `inject:""`
}
```

will essentially call `createSomething(new(SomethingElse))` everytime `SomethingProvider()` is called,
passing the resulting instance thru the injection to finalize uninjected fields. 


### Optional injection

An injection struct tag can be marked as optional by adding the suffix `,optional` to it.
This means that for interfaces, slices, pointers etc where dingo can not resolve a concrete type, the `nil`-type is injected.

You can check via `if my.Prop == nil` if this is nil.

## Providing Default Configurations
In your modul you can provide default configuraton values by implementing the Method "DefaultConfig".

Example:
```go

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"session.backend":                "memory",
		"session.secret":                 "flamingosecret",
		...
	}
}

```


## Bindings

Dingo uses bindings to express dependencies resolutions, and will panic if there is more than one
binding for a type with the same name (or unnamed), unless you use multibindings.

### Bind

Bind creates a new binding, and tells Dingo how to resolve the type when it encounters a request for this type.
Bindings can chain, but need to implement the correct interfaces.

```go
injector.Bind((*Something)(nil))
```

### AnnotatedWith

By default a binding is unnamend, and thus requested with the `inject:""` tag.

However you can name bindings to have more concrete kinds of it. Using `AnnotatedWith` you can specify the name:

```go
injector.Bind((*Something)(nil)).AnnotatedWith("myAnnotation")
```

It is requested via the `inject:"myAnnotation"` tag. For example:

```go
struct {
    PaypalPaymentProcessor PaymentProcessor `inject:"Paypal"`
}
```

### To

To defines which type should be created when this type is requested.
This can be an Interface which implements to one it is bound to, or a concrete type.
The type is then created via `reflect.New`.

```go
injector.Bind((*Something)(nil)).To(MyType{})
```

### ToProvider

If you want a factory to create your types then you rather use `ToProvider` instead of `To`.

`ToProvider` is a function which returns an instance (which again will go thru Dingo to fill dependencies).

Also the provider can request arguments from Dingo which are necessary to construct the bounded type.
If you need named arguments (e.g. a string instance annotated with a configuration value) you need to request
an instance of an object with these annotations, because Go does not allow to pass any meta-information on function
arguments.

```go
func MyTypeProvider(se SomethingElse) *MyType {
    return &MyType{
        Special: se.DoSomething(),
    }
}

injector.Bind((*Something)(nil)).ToProvider(MyTypeProvider)
```

This example will make Dingo call `MyTypeProvider` and pass in an instance of `SomethingElse` as it's first argument,
then take the result of `*MyType` as the value for `Something`.

`ToProvider` takes precedence over `To`.

### ToInstance

For situations where you have one, and only one, concrete instance you can use `ToInstance` to bind
something to the concrete instance. This is not the same as a Singleton!
(Even though the resuting behaviour is very similar.)

```go
var myInstance = new(MyType)
myInstance.Connect(somewhere)
injector.Bind((*Something)(nil)).ToInstance(myInstance)
```

You can also bind an instance it to a struct obviously, not only to interfaces.

`ToInstance` takes precedence over both `To` and `ToProvider`.

### In (Singleton scopes)

If really necessary it is possible to use singletons
``` 
.AsEagerSingleton() binds as a singleton, and loads it when the application is initialized
.In(dingo.Singleton) makes it a global singleton
.In(dingo.ChildSingleton) makes it a singleton limited to the config area
```

`In` allows us to bind in a scope, making the created instances scoped in a certain way.

Currently Dingo only allows to bind to `dingo.Singleton` and `dingo.ChildSingleton`.

```go
injector.Bind((*Something)(nil)).In(dingo.Singleton).To(MyType{})
```

#### dingo.Singleton

The `dingo.Singleton` scope makes sure a dependency is only resolved once, and the result is
reused. Because the Singleton needs synchronisation for types over multiple concurrent
goroutines and make sure that a Singleton is only created once, the initial creation
can be costly and also the injection of a Singleton is always taking more resources than creation
of an immutable new object.

The synchronisation is done on multiple levels, a first test tries to find the singleton,
if that is not possible a lock-mechanism via a scoped Mutex takes care of delegating
the concrete creation to one goroutine via a scope+type specific Mutex which then generates
the Singleton and makes it available to other currently waiting injection requests, as
well as future injection requests.

By default it is advised to not use Singletons whenever possible, and rather use
immutable objects you inject whenever you need them.

#### dingo.ChildSingleton

The ChildSingleton is just another Singleton (actually of the same type), but dingo will create a new one
for every derived child injector.

This allows frameworks like Flamingo to distinguish at a root level between singleton scopes, e.g. for
multi-page setups where we need a wide scope for routers.

Since ChildSingleton is very similar to Singleton you should only use it with care.

#### AsEagerSingleton

Singleton creation is always costly due to synchronisation overhead, therefore
Dingo bindings allow to mark a binding `AsEagerSingleton`.

This makes sure the Singleton is created as soon as possible, before the rest of the Application
runs. `AsEagerSingleton` implies `In(dingo.Singleton)`.

```go
injector.Bind((*Something)(nil)).To(MyType{}).AsEagerSingleton()
```

It is also possible to bind a concrete type without `To`:

```go
injector.Bind(MyType{}).AsEagerSingleton()
```

Binding this type as an eager singleton inject the singleton instance whenever `MyType` is requested. `MyType` is a concrete type (struct) here, so we can use this mechanism to create an instance explicitly before the application is run.

### Override

In rare cases you might have to override an existing binding, which can be done with `Override`:

```go
injector.Override((*Something)(nil), "").To(MyBetterType{})
```

`Override` also returns a binding such as `Bind`, but removes the original binding.

The second argument sets the annotation if you want to override a named binding.

### MultiBindings

MultiBindings provide a way of binding multiple implementations of a type to a type,
making the injection a list.

Essentially this means that multiple modules are able to register for a type, and a user of this
type can request an injection of a slice `[]T` to get a list of all registered bindings.

```go
injector.BindMulti((*Something)(nil)).To(MyType1{})
injector.BindMulti((*Something)(nil)).To(MyType2{})

struct {
    List []Something `inject:""`  // List is a slice of []Something{MyType1{}, MyType2{}}
}
```

MultiBindings are used to allow multiple modules to register for a certain type, such as a list of
encoders, subscribers, etc.

Please not that MultiBindings are not always a clear pattern, as it might hide certain complexity.

Usually it is easier to request some kind of a registry in your module, and then register explicitly.


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

### Binding basic types

Dingo allows binding values to `int`, `string` etc., such as with any other type.

This can be used to inject configuration values.

Flamingo makes an annotated binding of every configuration value in the form of:

```go
var Configuration map[string]interface{}

for k, v := range Configuration {
	injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
}
```

In this case Dingo learns the actual type of `v` (such as string, bool, int) and provides the annotated injection.

Later this can be used via

```go
struct {
	ConfigParam string `inject:"config:myconfigParam"`
}
```



## Dingo Interception

Dingo allows modules to bind interceptors for interfaces.

Essentially this means that whenever the injection of a certain type is happening,
the interceptor is injected instead with the actual injection injected into the interceptor's
first field. This mechanism can only work for interface interception.

Multiple interceptors stack upon each other.

Interception should be used with care!

```go
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindInterceptor((*template.Engine)(nil), TplInterceptor{})
	injector.BindInterceptor((*template.Function)(nil), FunctionInterceptor{})
}

type (
	TplInterceptor struct {
		template.Engine
	}

	FunctionInterceptor struct {
		template.Function
	}
)

func (t *TplInterceptor) Render(context web.Context, name string, data interface{}) io.Reader {
	log.Println("Before Rendering", name)
	start := time.Now()
	r := t.Engine.Render(context, name, data)
	log.Println("After Rendering", time.Since(start))
	return r
}

func (f *FunctionInterceptor) Name() string {
	funcname := f.Function.Name()
	log.Println("Function", funcname, "used")
	return funcname
}
```



## Initializing Dingo
At the topmost level the injector is created and used in the following way:

```go
package main

import "go.aoe.com/flamingo/framework/dingo"

func main() {
  var injector = dingo.NewInjector()
  
  // The injector can be initialized by modules:
  injector.InitModules(new(BillingModule))
  
  /*
   * Now that we've got the injector, we can build objects.
   * We get a new instance, and cast it accordingly:
   */
  var billingService = injector.GetInstance((*BillingService)(nil)).(BillingService)
  //...
}
```

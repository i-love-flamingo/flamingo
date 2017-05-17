# Dingo Dependency Injection

## About dependency injection

In general we suggest using Dependency Injection Pattern - this Patterns means nothing more than
if your object requires collaboration with others, then expect the user (or client)
of your object to set this dependencies from outside.

To use this pattern you don't need a seperate dependency injection container. 
But using this approach gives you higher testability and often leads to a cleaner and flexible architecture.
Typical "things" that can be injected are services, repositories or factories. If your object just expects a certain
"interface" the user/client of your object can decide what concrete object it wants your object to use.

It might sound like a "hen <-> egg" problem - because someone has to decide on the concrete instance that should
be injected.
 
So somewhere it need to start - and someone needs to inject the correct dependencies to your object - right?
This can be for example:
- the orchestration logic (normaly in the application layer) deciding which instance(s) to inject.
You can achive this without any framework support.
- a dependency registration concept - where you allow also other packages to influence which object should be injected.
This normaly requires a dependency injection container in the framework.

## DI Container in Flamingo

Flamingo Framework comes with a DI Container called Dingo.

The Container acts as kind of registry for services (objects of any type), factories and parameters.
The container can then return (or resolve) objects and can inject depenendcies automatically with some magic involved.

It is mainly used in the core for:
- managing different contexts and stateful objects (like routing) in the contexts
- registering ports and adapters
- ...

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

The modules provide information to Dingo on how to build the injection tree, and at the topmost level
the injector is created and used in the following way:

```go
package main

import "flamingo/framework/dingo"

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

A module itself can ask for dependencies, e.g.

```go
type Module struct {
    RouterRegistry *router.Registry `inject:""`
}

func (m *Module) Configure(injector *dingo.Inject) {
    m.RouterRegistry.Register(new(SomethingRegisterable))
}
```

# Bindings

Dingo uses bindings to express dependencies resolutions, and will panic if there is more than one
binding for a type with the same name (or unnamed), unless you use multibindings.

## Bind

Bind creates a new binding, and tells Dingo how to resolve the type when it encounters a request for this type.
Bindings can chain, but need to implement the correct interfaces.

```go
injector.Bind((*Something)(nil))
```

## AnnotatedWith

By default a binding is unnamend, and thus requested with the `inject:""` tag.

However you can name bindings to have more concrete kinds of it. Using `AnnotatedWith` you can specify the name:

```go
injector.Bind((*Something)(nil)).AnnotatedWith("myAnnotation")
```

It is requested via the `inject:"myAnnotation"` tag.

## To

To defines which type should be created when this type is requested.
This can be an Interface which implements to one it is bound to, or a concrete type.
The type is then created via `reflect.New`.

```go
injector.Bind((*Something)(nil)).To(MyType{})
```

## ToProvider

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

## ToInstance

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

## In

`In` allows us to bind in a scope, making the created instances scoped in a certain way.

Currently Dingo only allows to bind to `dingo.Singleton`, no other scopes exist.

```go
injector.Bind((*Something)(nil)).In(dingo.Singleton).To(MyType{})
```

### dingo.Singleton

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

## AsEagerSingleton

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

## Override

In rare cases you might have to override an existing binding, which can be done with `Override`:

```go
injector.Override((*Something)(nil), "").To(MyBetterType{})
```

`Override` also returns a binding such as `Bind`, but removes the original binding.

The second argument sets the annotation if you want to override a named binding.

# MultiBindings

MultiBindings provide a way of binding multiple implementations of a type to a type,
making the injection a list.

Essentially this means that multiple modules are able to register for a type, and a user of this
type can request an injection of `[]T` to get a list of all registered bindings.

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

# Requesting injection

Dingo uses struct tags to allow structs to request injection into fields.

As stated earlier this is the only Go-way to make annotations, everything else is not supported!

For every requested injection (unless an exception applies) Dingo does the following:

- Is there a binding? If so: delegate to the binding
  - Is the binding in a certain scope (Singleton)? If so, delegate to scope (might result in a new loop)
  - Binding is bound to an instance: inject instance
  - Binding is bound to a provider: call provider
  - Binding is bound to a type: request injection of this type (might return in a new loop to resolve the binding)
- No binding? Try to create (only possible for concrete types, not interfaces or functions)

## MultiBindings

Injection of multibindings:

```go
struct {
    ListOfTypes []Type `inject:""`
}
```

## AnnotatedWith

Injection of annotated values:

```go
struct {
    PaypalPaymentProcessor PaymentProcessor `inject:"Paypal"`
}
```

## Provider

Dingo allows to request the injection of provider instead of instances.

```go
struct {
    PizzaProvider func() Pizza `inject:""`
}
```

If there is no concrete binding to the type `func() Pizza`, then instead of constructing one `Pizza` instance
Dingo will create a new function which, on every call, will return a new instance of `Pizza`.

The type must be of `func() T`, a function without any arguments which returns a type, which again has a binding.

This allows to lazily create new objects whenever needed, instead of requesting the Dingo injector itself.

For example:

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

# How and Where "injection" can be used

Every instance that is created through the container can use injection. 
To get instances injected you just have to use the "inject" annotation inside of structs like this:

```
type RenderAware struct {
	Router *router.Router  `inject:""`
	Engine template.Engine `inject:""`
}
```

## When to use the depencency injection container in flamingo


- It is ok to not use the dependency injection container. In fact overusing the container adds unneccessary complexity.
When writing a package you should think of beeing able to also use it without the container
 So it is ok to:
	 - Explicitly initialize your object yourself and decide in the application layer what to inject (if you use dependency injection)
	 - Explicitly use your own factory directly 
- Every object that has a state that is depending on the running configuration-context, e.g. in a project where multiple configuration-contexts exist, should be injected,
because every configuration-context has its own initialized container the di container takes care of giving you
the correct initialized instance.
	 - For example the Router ( ```*responder.RenderAware `inject:""` ``` )
	 - Also for settings/parameters/configurations  
- Also the DI Container is used get the "right" interface implementation - in order to implement a flexible
"Ports and Adapters" concept (see below)

## Ports and Adapters with the Container

Basti: Wouln't it be cool 

# Binding basic types

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

# Dingo Interception

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

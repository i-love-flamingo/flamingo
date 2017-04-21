# Dependency Injection

## About dependency injection

In general we suggest using Dependency Injection Pattern - this Patterns means nothing more then if your object requires collaboration with others, then expect the user (or client) of your object to set this dependencies from outside.
To use this pattern you don't need a seperate dependency injection container. 
But using this approach gives you higher testability and often leads to a cleaner and flexible architecture.
Typical "things" that can be injected are services, repositories or factories. If your object just expects a certain "interface" the user/client of your object can decide what concrete object it wants your object to use.

It might sound like a "hen <-> egg" problem - because someone has to decide on the concrete instance that should be injected. 
So somewhere it need to start - and someone needs to inject the correct dependencies to your object - right?
This can be for example:
* the orchestration logic (normaly in the application layer) deciding which instance(s) to inject. You can achive this without any framework support.
* a dependency registration concept - where you allow also other packages to influence which object should be injected. This normaly requires a dependency injection container in the framework.

## Dingo AOP

```go
package main

func (m *Module) Configure(injector *dingo.Injector) {
	// ...
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
	r := f.Function.Name()
	log.Println("Function", r, "used")
	return r
}
```

## DI Container in Flamingo

Flamingo Framework comes with a DI Container. 
The Container acts as kind of registry for services (objects of any type), factories and parameters.
The container can then return (or resolve) objects and can inject depenendcies automatically with some magic involved.

It is mainly used in the core for:
* managing different contexts and stateful objects (like routing) in the contexts
* registering ports and adapters
* ...


### Register Services (unnamed)

Given "container" is an instace of the Container you can:

Register any Type by calling Register:

```
EventRouter := new EventRouter()
container.Register(EventRouter)
```
Type can be an Pointer or Struct or any other type.
If you register a Pointer the DI Container would inject that Pointer - and so you will have something like singletons.

Now you can get this object injected like this:

```
EventRouter      *EventRouter `inject:""`
```
Please note that this injection only works for Pointers, Interface Types (or Functions). (@Basti - why? and should we disallow registration of other types already?)

In addition you can assign tags while registering objects:

```
container.Register(Widget, "tag1","tag2")
```
Certain tags are used by convention in the Core - e.g. "template.func" - they are retrieved with container.GetTagged


-----

@Basti to discuss. Its a wiring through implicit conventions.. Maybe instead we should have speaking interfaces to some objects hadling thins - e.g.
  EventPublisher.Subscribe()  or
  Template.RegisterTemplateFunc()
  
Just discussion: Maybe another bootstrap could solve this? 
  Create Context
  Initialize context and its Containers (Calls Register funcs)
  Then call a PackageInit interface in every package that can use DI already to initialize something on the Framework Objects (like routing, templating...)
  Then load contexts configuration ?  
  
-----

### Register Named Services
Basti: Where is it used? Should we use it for interface types and ports and adapters?

### Parameter (aka Settings) injection

Parameters can be added to the Container and Parameters can then be injected or used in service Factories (see below)

Parameters can be set like this:
```
container.SetParameter("somepackage.apiclient.masterdataportal.uri","https://example.org")
```
But in fact they are most likely be set in the configuration (see [structure of projects](structure,md))


Get access to a parameter is as simple as:
```
type Some struct {
  ApiURI string `inject:"param:somepackage.apiclient.masterdataportal.uri"`
```
ToDo - support Default Value? (Basti?))

Todo - How to get access if no injection is used (see below) @Basti?

### Register Factorys

Factories return a new instance of an object each time you call them. 
You can register Factories that will be called everytime the Container needs to inject a object of the type:

```
 container.RegisterFactory(func() EventRouter { return new(EventRouterImpl) })
```

Now when someone expects a "EventRouter" - the container will call this factory and passes the instance to the object.

A factory can also be parameterized. Use this in cases where you depend on some context specific settings, because the container will make sure that the factory is called with the correct Parameters (aka Config)
```
container.RegisterFactory(func(bar int, foo string) *Something { return &Something{Foo: foo, Bar: bar} }).AddParameters("bar", "foo")

```

(@Basti - why RegisterFactory does not check for func type? Also we should avoid mixing DDD Factory and the Factory in a container.. (Maybe name it serviceCreator?))

### Using functions in the container

Registering functions can be used also:

Use Cases are:
* Register Routes on the RoutingService in the Context


(Basti: Should we use it this way - are there other use cases? For routing there may be a way to make it more speaking? 

Also whats the point of adding it to the cache? Isnt it overriden each time anyhow? and its probably not used for injection? 
When is this executed. Also seems to be executed twice?
Maybe the "registerfunc" get container and routing object? Pro and Con?  + Same question like for tags..  



Example Usage:
```
container.Register(func(r *router.Router) {
		fmt.Println("is called now 2")
		// default handlers
		//r.Handle("cart.page.view", new(application.CartController))

		// default routes
		r.Route("/page/{name}", "cms.page.view")
	}, "router.register")
	
	
		for _, registerFunc := range container.GetTagged("router.register") {
  		registerFunc.Value.(func(r *Router))(router)   //Basti? Why the . after Value? And why is it wrapped?
  	}
```




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


* its ok to not use the dependency injection container.  In fact overusing the container adds unneccessary complexity. When writing a package you should think of beeing able to also use it without the container
 So it is ok to:
   * Explicitly initialize your object yourself and decide in the application layer what to inject (if you use dependency injection)
   * Explicitly use your own factory directly
   * 
   
* Every object that has a state that is depending on the "context" and the "bootstraping" should be injected, because every context has its own initialized container the di container takes care of giving you the correct initialized instance.
   *  For example the Router ( *responder.RenderAware `inject:""` ) - (BASTI: btw - why does RenderAware / Controller need the Router ? Maybe we just need a Dispatcher?)
   *  Also for settings/parameters/configurations 
   
* Also the DI Container is used get the "right" interface implementation - in order to implement a flexible "Ports and Adapters" concept (see below)
 
## Ports and Adapters with the Container

Basti: Wouln't it be cool 

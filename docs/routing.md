# Routing

-----
@Basti 
Normaly 

 METHOD:PATH -> Controller.Action
 
Here:

 PATH -> Controller  / Action Func
 
    - Method indirect
    - No Action
    
    
Also support redirects (No need for custom package)?

/ -> RedirectController(cms.page.view(home))


Also assigning multiple routes to one handler is not working (should we support it?)

( Maybe ) Also allow some generated documentation (e.g. for cart API payload etc)
I Read: hub.io/golang/urlrouter/vestigo/2015/09/22/vesigo.html
------

# Changelog

## v3

- "locale" package:
  - the templatefunc __(key) is now returning a Label and instead additional parameters you need to use the label setters (see doc)
- Depricated Features are removed:
  - `flamingo.me/dingo` need to be used now
  - support for responder.*Aware types is removed
- `framework/web.Response` is now `framework/web.Result`
- `web.Request` is heavily condensed
  - Access to Params has changed
- `web.Session` does not expose `.GS()` for the gorilla session anymore
- `event.Subscriber` changes:
  - is getting `context.Context` as the first argument: `Notify(ctx context.Context, e flamingo.Event)`
  - `event.Subscriber` are registered via `framework/flamingo.BindEventSubscriber(injector).To(...)`
  - There is no SubscriberWithContext anymore!
- several other Modules have been moved out of flamingo and exist as seperate modules:
  - **For all the stuff in this section:** you may use the script `docs/updates/v3/renameimports.sh` for autoupdate the import paths in your project and to do some first replacements.
  - moved modules outside of flamingo:
    - flamingo/core/redirects
    - flamingo/core/pugtemplate
    - flamingo/core/form2
    - flamingo/core/form (removed)
    - flamingo/core/csrf
    - flamingo/core/csp
    - flamingo/core/captcha
    
  - restructures inside flamingo:
    - `core/requestTask` is renamed to `core/requesttask`
    - `core/canonicalUrl` is renamed to `core/canonicalurl`
    - `core/cmd` package moved to `framework/cmd` and the cmd have been moved to the packages they belong to
    - `framework/router` package merged into `framework/web`
    - `framework/event` package merged into `framework/flamingo`
    - `framework/template` package merged into `framework/flamingo`:
      - instead of `template.BindFunc` and `template.BindCtxFunc` you have to use `flamingo.BindTemplateFunc`
    - `framework/session` package merged into `framework/flamingo`:
      - instead of using the module `session.Module` use `flamingo.SessionMdule`
  
- flamingo now uses go mod - we encourage to use it also in the projects:
  - Init the project
    ```
    go mod init YOURMODULEPATH
    ```
  - If you want to link the flamingo core to your project (because you are working on the core also)
    - Option 1:
      edit "go.mod" and add this lines (make sure to not commit them to git)
      ```
      replace (
        flamingo.me/flamingo/v3 => ../../PATHTOFLAMINGO
        flamingo.me/flamingo-commerce/v3 => ../../PATHTOFLAMINGO
      )
      ```
    - Option 2:
      use `go mod vendor` and link the modules after this from vendor folder
- Prefixrouter configuration: rename *prefixrouter.baseurl* in *flamingo.router.path*

## v2

- `web.Responder` is now used
  - instead of injecting 
    ```
     responder.JSONAware
     responder.RenderAware
     responder.RedirectAware
     ``` 
     in a controller you need to inject 
     ```
     responder *web.Responder
     ```
     And use the Methods of the Responder:
     `c.responder.Data()` `c.responder.Render()`  `c.responder.Redirect()`
- `dingo` is moved out to `flamingo.me/dingo` and we recommend to use the Inject() methods instead of public properties.

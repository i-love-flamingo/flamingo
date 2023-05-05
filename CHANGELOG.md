# Changelog

## Version v3.6.1 (2023-05-04)

### Fixes

- **web:** handle reverse router being nil (2b79ba7d)

### Chores and tidying

- **deps:** update module github.com/spf13/cobra to v1.7.0 (#252) (5ac51fbd)

## Version v3.6.0 (2023-04-19)

### Features

- **zap:** silent zap logging (#308) (c91a24b3)

### Chores and tidying

- **deps:** update module github.com/go-redis/redis/v8 to v9 (#333) (d5741fd9)
- add regex to detect go run/install commands (dc92fff0)

## Version v3.5.1 (2023-03-10)

### Fixes

- **systemendpoint:** prevent data race on server (#324) (0d094a24)
- updating community join url to discord (#320) (7947dfd5)
- **core/auth/oauth:** Properly mark OIDC callback error handler as optional (#318) (f45a90c1)

### Documentation

- introduce slack channel (#323) (23a68460)

### Chores and tidying

- **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 (59000bde)
- **deps:** update mockery golangci-lint and  (#325) (4fd9102d)

## Version v3.5.0 (2023-01-26)

### Features

- **core/auth/oauth:** Add optional OIDC callback error handler (4fe0f89e)

### Fixes

- **cache:** do not rewrite trace id (#306) (b1a5918a)
- **framework/opencensus:** x-correlation-id condition (#310) (5909a729)

### Chores and tidying

- **deps:** update module github.com/coreos/go-oidc/v3 to v3.5.0 (#313) (372b7b45)
- **deps:** update module golang.org/x/oauth2 to v0.4.0 (bd3ae5d3)
- **deps:** update module github.com/hashicorp/golang-lru to v0.6.0 (9c62154e)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v4.4.3 (7189aed3)
- **deps:** update module contrib.go.opencensus.io/exporter/prometheus to v0.4.2 (cb8d1a5d)
- make linter happy (3017d8ab)
- **deps:** update module go.opencensus.io to v0.24.0 (#302) (6e915e2a)
- **deps:** update module flamingo.me/dingo to v0.2.10 (#299) (0aace54a)
- **deps:** update module github.com/gofrs/uuid to v4.3.1+incompatible (#295) (1fe7e9c4)
- **deps:** update module golang.org/x/oauth2 to v0.2.0 (#296) (0ad70578)

## Version v3.4.0 (2022-11-03)

### Features

- **session:** add timeout for redis connection (#277) (eb0fe55f)
- **framework/config:** module dump command (#273) (8a1b9154)
- **auth:** add session refresh interface (#269) (3eba3475)
- **core/oauth:** support issuer URL overriding (#227) (fa5bd34b)
- **oauth:** add oauth identifier (#220) (9883a4dc)

### Fixes

- **framework/web:** Avoid error log pollution, switch context cancelled log level to debug (#294) (c1d6dc87)
- **framework/flamingo/log:** added nil check for StdLogger and Logger (f636b7ec)
- **framework:** Add missing scheme property to the router cue config (#274) (3ee35cd8)
- **router:** allow config for public endpoint (3f47d251)
- **framework/systemendpoint:** use real port information in systemendpoint (4f59dc4a)
- **framework/prefixrouter:** use real port information in ServerStartEvent (79ae6f95)
- **servemodule:** use real port information in ServerStartEvent (c5209de3)
- **oauth:** correctly map access-token claims (5a7331f3)
- **auth:** add missing auth.Identity interface (#216) (27b93c16)
- **deps:** exclude unmaintained redigo (#218) (6061f4ab)
- fix missing gob register (0c488981)

### Tests

- **internalauth:** add unittests (#258) (46341b4d)
- **requestlogger:** add unittests (c8c18474)
- **framework/opencensus:** add tests (932299a5)

### Ops and CI/CD

- adjust gloangci-lint config for github actions (bbabfc08)
- make "new-from-rev" work for golangci-lint (258a7b50)
- remove now unnecessary steps from main pipeline (dc2df28f)
- fix git rev (84ebd56d)
- add golangci-lint to pipeline (bf3eaeab)
- **semanticore:** add semanticore (a741f30d)

### Documentation

- **core/gotemplate:** Enhance documentation (#291) (8b5b8a97)
- typos and wording (#290) (9745b18e)
- **flamingo:** update logos (4227815b)

### Chores and tidying

- **deps:** update module github.com/coreos/go-oidc/v3 to v3.4.0 (#293) (f5a791f1)
- **deps:** update module github.com/google/go-cmp to v0.5.9 (#282) (8d546ca7)
- **deps:** update module github.com/openzipkin/zipkin-go to v0.4.1 (#286) (bf432c9f)
- **deps:** update module github.com/stretchr/testify to v1.8.1 (#292) (56ddf827)
- **deps:** update irongut/codecoveragesummary action to v1.3.0 (#278) (2399c75f)
- bump to go 1.19 (#279) (da07d03e)
- **deps:** update module github.com/stretchr/testify to v1.8.0 (#271) (c81b3226)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v4.4.2 (#272) (7dd92887)
- **deps:** update module github.com/stretchr/testify to v1.7.2 (#270) (fb271199)
- **deps:** update module go.uber.org/automaxprocs to v1.5.1 (2b868ee5)
- **deps:** update module github.com/openzipkin/zipkin-go to v0.4.0 (38acb057)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v4.4.1 (daf8fd7a)
- **deps:** update module github.com/hashicorp/golang-lru to v0.5.4 (2b8bd64c)
- **gomod:** go mod tidy (3e911268)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.2.0 (0e52661e)
- **deps:** update dependency quay.io/dexidp/dex to v2.28.1 (2643bd00)
- **deps:** update module contrib.go.opencensus.io/exporter/prometheus to v0.4.1 (5bd0be1a)
- **deps:** update module github.com/stretchr/testify to v1.7.1 (d643c92c)
- **deps:** update module github.com/google/go-cmp to v0.5.8 (559dea58)
- **deps:** update module contrib.go.opencensus.io/exporter/zipkin to v0.1.2 (11533458)
- **deps:** update actions/setup-go action to v3 (3592ea01)
- **deps:** update actions/checkout action to v3 (c9d2f124)
- **deps:** update module github.com/go-redis/redis/v8 to v8.11.5 (2b6581a8)
- **deps:** update module contrib.go.opencensus.io/exporter/jaeger to v0.2.1 (f7285f09)
- **deps:** update dependency quay.io/keycloak/keycloak to v8.0.2 (58f1a27b)
- **deps:** update golang.org/x/oauth2 digest to 9780585 (bb60190b)
- **deps:** update github.com/golang/groupcache digest to 41bb18b (68a4f7e2)
- Update renovate.json (ca2c97b9)
- **deps:** add renovate.json (e0edaf27)
- bump go version to 1.17, replace golint with staticcheck (#222) (ae2b39e8)
- **auth:** switch to github.com/gofrs/uuid (1854abc6)

### Other

- continue-on-error for coverage pr comment (#288) (21b26df3)
- framework/flamingo: replace redis session backend (#219) (8451ed0b)
- core/auth: update to recent go-oidc v3, allow oidc issuer URL override (#212) (86076485)
- add comment to StateEntry (a8be5d77)
- allow multiple parallel state responses (d7b30a06)

## Important Notes

- core/internalauth:
  - switched from `github.com/dgrijalva/jwt-go` to `github.com/golang-jwt/jwt/v4`. this is a drop-in replacement
    
    use search and replace to change the import path or add a replace statement to your go.mod:
    ```
    replace (
        github.com/dgrijalva/jwt-go v3.2.0+incompatible => github.com/golang-jwt/jwt/v4 v4.1.0
    )
    ```
    More details can be found here: https://github.com/golang-jwt/jwt/blob/main/MIGRATION_GUIDE.md
- core/auth:
  - oauth.Identity includes Identity. This is a backwards-incompatible break

## v3.2.0

- license:
  - Flamingo now uses the MIT license. The CLA has been removed.
- core/auth:
  - Flamingo v3.2.0 provides a new auth package which makes authentication easier and more canonical.
  - the old core/oauth is deprecated and provides a compatibility layer for core/auth.
- sessions:
  - web.SessionStore provides programmatic access to web.Session
  - flamingo.session.saveMode allows to define a more granular session save behaviour
- config loading:
  - both routes.yml and routes.yaml are now supported
- framework/web:
  - the framework router got a couple of stability updates.
  - the web responder and responses don't fail anymore for uninitialized responses.
  - error responses are wrapped with a http error message
  - the flamingo.static.file controller needs a dir to not serve from root.
- errors:
  - all errors are handled via Go's error package
- go 1.13/1.14:
  - support for 1.12 has been dropped

## v3

- "locale" package:
  - the templatefunc __(key) is now returning a Label and instead additional parameters you need to use the label setters (see doc)
- Deprecated Features are removed:
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
- several other Modules have been moved out of flamingo and exist as separate modules:
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

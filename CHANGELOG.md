# Changelog

## Version v3.17.1 (2025-11-27)

### Chores and tidying

- **deps:** bump golang.org/x/crypto from 0.41.0 to 0.45.0 (#547) (53b2d8d8)
- **deps:** update module golang.org/x/sync to v0.18.0 (#544) (c38d2e63)
- **deps:** update module github.com/redis/go-redis/v9 to v9.16.0 (#540) (3196a62c)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.16.0 (#536) (a2c11be8)
- **deps:** update module golang.org/x/oauth2 to v0.32.0 (#538) (d1264567)
- **deps:** update dependency go to v1.25.3 (#537) (39128908)
- **deps:** update actions/setup-go action to v6 (#534) (ff47d8ca)

## Version v3.17.0 (2025-09-16)

### Features

- **session:** add custom session store backend using binding to CustomSessionBackend (#527) (5cb20714)

### Chores and tidying

- **deps:** update linter and mockery (#535) (8bfa5dae)
- **deps:** bump github.com/go-viper/mapstructure/v2 (#528) (6e59a9b6)
- **deps:** update module github.com/golang-jwt/jwt/v5 to v5.3.0 (#520) (ede1e7b8)
- **deps:** update module github.com/redis/go-redis/v9 to v9.14.0 (#516) (3c6c2096)
- **deps:** update module golang.org/x/sync to v0.17.0 (#533) (66fb1c3d)
- **deps:** update module github.com/stretchr/testify to v1.11.1 (#529) (ceebbff0)
- **deps:** update module golang.org/x/oauth2 to v0.31.0 (#532) (10c0db59)
- **deps:** update module github.com/spf13/cobra to v1.10.1 (#531) (4b520bbb)
- **deps:** update module github.com/spf13/pflag to v1.0.10 (#530) (284b112e)
- **deps:** update actions/checkout action to v5 (#526) (e3334ca1)
- **deps:** update module github.com/spf13/pflag to v1.0.7 (#521) (c78205f1)

## Version v3.16.0 (2025-08-06)

### Features

- **session:** add username support for redis (#524) (d503f108)

### Chores and tidying

- **deps:** bump github.com/go-viper/mapstructure/v2 (#517) (168a4e63)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.15.0 (#523) (d68f6482)
- **deps:** update module golang.org/x/sync to v0.16.0 (#519) (a6d3300e)

## Version v3.15.0 (2025-06-13)

### Features

- zap logger caller skip configurable (#512) (6494dc09)

### Chores and tidying

- **deps:** update module github.com/redis/go-redis/v9 to v9.10.0 (#501) (56776d36)
- **deps:** update dependency go to v1.24.4 (#496) (3dd1f7d5)
- **deps:** update module github.com/vektra/mockery/v3 to v3.3.6 (#513) (d6f1605b)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.14.1 (#497) (886be3a6)
- **deps:** update module github.com/vektra/mockery/v2 to v3 (#508) (60afe694)
- **deps:** update module golang.org/x/sync to v0.15.0 (#509) (dc9f9f88)
- migrate golangci-lint to v2 (#506) (8e42afee)

## Version v3.14.0 (2025-05-19)

### Features

- **healthcheck:** add measured status healthcheck interface (#503) (69049df2)
- **web:** support Flush on http responses (#494) (e83a8c62)

### Chores and tidying

- **deps:** update module github.com/gorilla/sessions to v1.4.0 (#428) (3cc8e206)
- **deps:** update module golang.org/x/oauth2 to v0.30.0 (#498) (b3368bea)
- **deps:** update module golang.org/x/sync to v0.14.0 (#499) (e017aa5b)
- **deps:** bump golang.org/x/net from 0.37.0 to 0.38.0 (#500) (57ea554d)
- **deps:** update dependency go to v1.24.1 (#486) (712321fc)
- **deps:** update module github.com/google/go-cmp to v0.7.0 (#482) (85d60c90)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.13.0 (#489) (e5313988)
- **deps:** update module github.com/redis/go-redis/v9 to v9.7.3 [security] (#491) (cbc03189)
- **deps:** bump github.com/golang-jwt/jwt/v5 from 5.2.1 to 5.2.2 (#493) (eabcfb91)
- **deps:** bump github.com/go-jose/go-jose/v4 from 4.0.2 to 4.0.5 (#484) (a7e5d1d4)
- make post run script available for both Run and RunE functions (#478) (7bbfff73)

### Other

- Add tests for event dispatch when running commands and graceful shutdown (#485) (875204d2)

## Version v3.13.0 (2025-02-17)

### Features

- **core/zap:** Make caller encoder configurable, add additional one (#301) (fc037fce)

### Chores and tidying

- **deps:** update module github.com/spf13/cobra to v1.9.1 (#477) (5fad27fd)
- **deps:** update module golang.org/x/sync to v0.11.0 (#474) (5800854f)
- **deps:** update module golang.org/x/oauth2 to v0.26.0 (#473) (48e6b955)
- **deps:** update module github.com/spf13/pflag to v1.0.6 (#471) (5667c5d3)
- **deps:** bump golang.org/x/net from 0.27.0 to 0.33.0 (#467) (45c1368c)
- bump go version with mockery and linter (#476) (98b92505)
- edit description of sessions module (#475) (fe338e85)
- **deps:** update module github.com/vektra/mockery/v2 to v2.51.1 (#468) (caa5d915)

## Version v3.12.0 (2025-01-09)

### Features

- **app:** add dingo-trace-injections flag, which switches injection tracing (#455) (c25a60ee)
- **web/handler:** add http method to opencensus tag for flamingo_router_controller (#432) (7597e699)

### Chores and tidying

- **deps:** update module github.com/coreos/go-oidc/v3 to v3.12.0 (#465) (e40d9eda)
- **deps:** update module github.com/vektra/mockery/v2 to v2.50.4 (#464) (7d2d7f05)
- **deps:** update module golang.org/x/oauth2 to v0.25.0 (#466) (27262c8a)
- **docs:** add troubleshooting information (#457) (8426b77e)
- edit mockery configuration, remove scattered go:generate annotations, edited linter config (#463) (586b32ba)
- **deps:** update dependency go to v1.23.4 (#459) (a307b167)
- **deps:** bump golang.org/x/crypto from 0.25.0 to 0.31.0 (#462) (ad862f2d)
- **deps:** update module golang.org/x/sync to v0.10.0 (#461) (e1b39ea6)
- **deps:** update module golang.org/x/oauth2 to v0.24.0 (#452) (a90e42b2)
- **deps:** update module golang.org/x/sync to v0.9.0 (#453) (acde3c2b)
- **deps:** update module github.com/vektra/mockery/v2 to v2.49.1 (#454) (af56c07a)
- **deps:** update dependency go to v1.23.3 (#451) (77cf423b)
- **deps:** update module github.com/redis/go-redis/v9 to v9.7.0 (#450) (0569eadc)
- **deps:** update dependency go to v1.23.2 (#447) (7d64c4e4)
- **deps:** update module github.com/vektra/mockery/v2 to v2.46.3 (#448) (2871fb36)
- **deps:** update quay.io/keycloak/keycloak docker tag to v25.0.6 (#443) (be287200)
- **deps:** update module go.uber.org/automaxprocs to v1.6.0 (#444) (9de9569f)
- **deps:** update module github.com/vektra/mockery/v2 to v2.46.1 (#445) (5158f5a6)

### Other

- Add version command to the default framework CLI commands (#458) (7903a65b)

## Version v3.11.0 (2024-09-19)

### Features

- decouple opencensus from core (#442) (95e83305)

### Fixes

- **cache:** pass deadline to context in load function (#440) (41c5fe24)

### Chores and tidying

- **deps:** update quay.io/keycloak/keycloak docker tag to v25.0.5 (#439) (634092be)
- **deps:** update dependency go to v1.23.1 (#437) (e9e669e0)
- **deps:** update module github.com/vektra/mockery/v2 to v2.46.0 (#438) (39733f67)

## Version v3.10.1 (2024-09-05)

### Fixes

- fix file session store creation, if directory exists (#433) (6469f207)

### Chores and tidying

- **deps:** update module golang.org/x/oauth2 to v0.23.0 (#435) (c4a1fe51)

## Version v3.10.0 (2024-09-02)

### Features

- **flamingo:** add optional attribute to configure a session key prefix for REDIS (flamingo.session.redis.keyPrefix) (#424) (70a73751)

### Fixes

- **framework:** fix attachment content disposition typo (#431) (e78acf00)

### Chores and tidying

- bump minimum Go version to 1.22 (#429) (5d6d1a0a)
- **deps:** update quay.io/keycloak/keycloak docker tag to v25.0.4 (#427) (8cb9ebe7)
- **deps:** update module github.com/vektra/mockery/v2 to v2.45.0 (#426) (9f959d7a)
- **deps:** update module golang.org/x/oauth2 to v0.22.0 (#422) (727144f4)
- **deps:** update module github.com/vektra/mockery/v2 to v2.44.1 (#421) (651a4b78)
- **deps:** update module golang.org/x/sync to v0.8.0 (#423) (b9d1a7b3)
- **deps:** update module github.com/redis/go-redis/v9 to v9.6.1 (#419) (2554e116)

## Version v3.9.0 (2024-07-18)

### Features

- **log:** Introduce Trace log level (#412) (9946e28b)
- **auth/oidc:** Use BadRequest response for state mismatches, enhance logging with context (#403) (85bcd6fe)

### Fixes

- **opencensus:** rename to takeParentDecision (78fd979e)

### Documentation

- add missing headline (58b1962e)
- add more information about the logger setup/usage (#414) (35f70ad6)

### Chores and tidying

- **testutil:** deprecate PACT support (#417) (5fd3f2c3)
- **deps:** update quay.io/keycloak/keycloak docker tag to v25.0.2 (#416) (7afaa6a0)
- **deps:** update module github.com/openzipkin/zipkin-go to v0.4.3 (#401) (6c780922)
- disable mockery version string (#415) (e26c5f08)
- **deps:** update module github.com/vektra/mockery/v2 to v2.43.2 (#402) (cd8f56ad)
- **deps:** update module github.com/redis/go-redis/v9 to v9.5.2 (#408) (bf44c971)
- **deps:** update quay.io/keycloak/keycloak docker tag to v25 (#409) (ec72db70)
- **deps:** update module github.com/gorilla/sessions to v1.3.0 (#411) (9574cf75)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.11.0 (#413) (585ad1fa)
- **deps:** update module golang.org/x/oauth2 to v0.21.0 (#405) (e3c774fc)
- **deps:** update module github.com/spf13/cobra to v1.8.1 (#410) (87a97fc0)
- **deps:** update golangci/golangci-lint-action action to v6 (#406) (3840ede5)

## Version v3.8.1 (2024-04-26)

### Ops and CI/CD

- fix version matrix (#398) (9a3f8230)

### Chores and tidying

- **deps:** update quay.io/keycloak/keycloak docker tag to v24 (#397) (e769812c)
- **deps:** update module github.com/vektra/mockery/v2 to v2.42.3 (#399) (ca1ae56b)
- **deps:** update golangci/golangci-lint-action action to v5 (#400) (9b17e2f9)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.10.0 (#393) (b00c8146)
- **dep:** update go to 1.22.0 (#379) (64cddb8a)
- **deps:** update module golang.org/x/sync to v0.7.0 (#396) (ed02827e)
- **deps:** bump golang.org/x/net from 0.22.0 to 0.23.0 (#395) (4cb80063)
- **deps:** update module golang.org/x/oauth2 to v0.19.0 (#394) (bcc53161)
- **deps:** update module github.com/golang-jwt/jwt/v5 to v5.2.1 (#390) (9037a18f)
- **deps:** update quay.io/keycloak/keycloak docker tag to v23.0.7 (#388) (157b02d5)
- **deps:** update module go.uber.org/zap to v1.27.0 (#387) (6f7f4fff)
- **deps:** update module github.com/redis/go-redis/v9 to v9.5.1 (#386) (4f32b254)
- **deps:** update module github.com/vektra/mockery/v2 to v2.42.2 (#384) (81b575b3)
- **deps:** update golangci/golangci-lint-action action to v4 (#385) (28304163)
- **deps:** update module golang.org/x/oauth2 to v0.18.0 (#383) (b2a64b45)
- **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.1 to 3.0.3 (#391) (268a07bd)
- **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 (#392) (0533567e)
- **deps:** update module github.com/stretchr/testify to v1.9.0 (#389) (352dec7c)
- fix linter errors (#381) (f2405d33)

## Version v3.8.0 (2024-02-07)

### Features

- **oidc:** Make state timeout duration configurable to support long taking sign ins (#362) (f4107584)

### Chores and tidying

- **dep:** update go to 1.21 (#380) (65f3f924)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v5 (#377) (f9ef8155)
- **deps:** update quay.io/keycloak/keycloak docker tag to v23 (#378) (79074598)
- **deps:** update module golang.org/x/sync to v0.6.0 (#375) (72d07ff3)
- **deps:** update actions/setup-go action to v5 (#376) (09552898)
- **deps:** update module go.uber.org/zap to v1.26.0 (#374) (ffa968ef)
- **deps:** update module github.com/vektra/mockery/v2 to v2.40.1 (#354) (1acc174f)
- **deps:** update module github.com/google/go-cmp to v0.6.0 (#372) (7f0687d8)
- **deps:** update module github.com/spf13/cobra to v1.8.0 (#373) (710a409c)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v5 (#345) (867154a6)
- **deps:** update module github.com/redis/go-redis/v9 to v9.4.0 (#351) (8271981c)
- **deps:** update module github.com/gofrs/uuid to v4.4.0+incompatible (#317) (14bc5586)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.9.0 (#371) (8f993a76)
- **deps:** update module github.com/gorilla/sessions to v1.2.2 (#370) (3e62d69d)
- **deps:** update module golang.org/x/oauth2 to v0.16.0 (#358) (5162059d)
- **deps:** update quay.io/keycloak/keycloak docker tag to v22.0.5 (#359) (5067e9a9)
- **deps:** update module github.com/hashicorp/golang-lru/v2 to v2.0.7 (#360) (bf357e3b)
- **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.1 (#363) (de2f7138)
- **deps:** bump golang.org/x/crypto from 0.12.0 to 0.17.0 (#367) (f57d291c)
- **docs:** fix link to opentelemetry (#366) (9e7dbb47)
- **docs:** Small documentation updates for docs.flamingo.me (#365) (9eb5e0fe)
- **opencensus:** deprecate opencensus module (#364) (8a2b30c4)
- **deps:** update actions/checkout action to v4 (#357) (258bd2d6)
- **deps:** update module github.com/hashicorp/golang-lru/v2 to v2.0.6 (#356) (d56a090a)
- **deps:** update module github.com/hashicorp/golang-lru to v2 (#314) (29db3cf4)

## Version v3.7.0 (2023-08-21)

### Features

- **log:** enhance memory efficiency for the zap logger with fields (#340) (f977f5f7)
- **core/auth/oidc:** Add support of `AuthCodeOption` in the token exchange during callback (#338) (b2453bab)

### Fixes

- **silentlogger:** remove unnecessary fields (#347) (63a53c88)
- **zap:** only record log metrics when we actually log something (#341) (b0dca2b3)

### Refactoring

- **core/zap:** clean up zap module and add tests (a7f94766)

### Ops and CI/CD

- force on latest go version &amp; cleanup jobs (#353) (69507dd2)

### Chores and tidying

- **deps:** update module golang.org/x/oauth2 to v0.11.0 (#331) (9d3adbd6)
- **deps:** update quay.io/keycloak/keycloak docker tag to v22 (#352) (98f5a027)
- **deps:** update module github.com/leekchan/accounting to v0.3.1 (#250) (63c33cd2)
- **deps:** update module go.uber.org/automaxprocs to v1.5.3 (#350) (7a7bf310)
- **deps:** update module github.com/openzipkin/zipkin-go to v0.4.2 (#348) (9cea11d2)
- **deps:** update module go.uber.org/zap to v1.25.0 (2ab0354b)
- **go:** bump go version to 1.20 (#349) (d63caa90)
- **deps:** update module github.com/vektra/mockery/v2 to v2.32.4 (#328) (38460154)
- **session:** health check config cleanup (#342) (41f0c480)
- **deps:** update module github.com/redis/go-redis/v9 to v9.0.5 (#343) (f616b5ea)
- **deps:** update module github.com/coreos/go-oidc/v3 to v3.6.0 (#344) (94a60a28)
- **deps:** update module github.com/golang-jwt/jwt/v4 to v4.5.0 (#322) (e6d8895d)
- **deps:** update module go.uber.org/automaxprocs to v1.5.2 (#330) (62c67564)
- **deps:** update actions/setup-go action to v4 (#332) (bec3e4cc)
- **deps:** update module github.com/stretchr/testify to v1.8.4 (#329) (063732d9)

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

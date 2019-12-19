module flamingo.me/flamingo/v3

go 1.13

require (
	contrib.go.opencensus.io/exporter/jaeger v0.1.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	cuelang.org/go v0.0.15
	flamingo.me/dingo v0.2.6
	github.com/boj/redistore v0.0.0-20180917114910-cd5dcc76aeff
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.3.1
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gorilla/sessions v1.1.3
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hashicorp/logutils v0.0.0-20150609070431-0dc08b1671f3 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/labstack/gommon v0.0.0-20180613044413-d6898124de91
	github.com/leekchan/accounting v0.0.0-20191104051123-0b9b0bd19c36
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/nicksnyder/go-i18n v0.0.0-20180814031359-04f547cc50da
	github.com/openzipkin/zipkin-go v0.2.0
	github.com/pact-foundation/pact-go v0.0.13
	github.com/pkg/errors v0.8.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5 // indirect
	github.com/zemirco/memorystore v0.0.0-20160308183530-ecd57e5134f6
	go.opencensus.io v0.22.2-0.20191001044506-fa651b05963c
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/automaxprocs v1.2.0
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190812073006-9eafafc0a87e // indirect
	google.golang.org/api v0.8.0 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/grpc v1.22.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/square/go-jose.v2 v2.1.9 // indirect
)

replace golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190121141151-b76268579942

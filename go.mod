module flamingo.me/flamingo/v3

require (
	flamingo.me/dingo v0.1.3
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/boj/redistore v0.0.0-20160128113310-fc113767cd6b
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gorilla/sessions v1.1.3
	github.com/hashicorp/golang-lru v0.5.0
	github.com/hashicorp/logutils v0.0.0-20150609070431-0dc08b1671f3 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/gommon v0.0.0-20180613044413-d6898124de91
	github.com/leekchan/accounting v0.0.0-20161211142212-a35854c07593
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/nicksnyder/go-i18n v0.0.0-20180814031359-04f547cc50da
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/pact-foundation/pact-go v0.0.13
	github.com/pkg/errors v0.8.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.2.2
	github.com/zemirco/memorystore v0.0.0-20160308183530-ecd57e5134f6
	go.opencensus.io v0.19.1
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/oauth2 v0.0.0-20190212230446-3e8b2be13635
	gopkg.in/square/go-jose.v2 v2.1.9 // indirect
)

replace (
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)

module flamingo.me/flamingo/v3

require (
	flamingo.me/dingo v0.1.3
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/boj/redistore v0.0.0-20160128113310-fc113767cd6b
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/corpix/uarand v0.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/captcha v0.0.0-20170622155422-6a29415a8364
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/etgryphon/stringUp v0.0.0-20121020160746-31534ccd8cac // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/ghodss/yaml v0.0.0-20180820084758-c7ce16629ff4
	github.com/go-playground/form v3.1.3+incompatible
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0
	github.com/go-test/deep v1.0.1
	github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
	github.com/google/uuid v1.1.0
	github.com/gorilla/sessions v1.1.3
	github.com/hashicorp/golang-lru v0.0.0-20180201235237-0fb14efe8c47
	github.com/hashicorp/logutils v0.0.0-20150609070431-0dc08b1671f3 // indirect
	github.com/icrowley/fake v0.0.0-20180203215853-4178557ae428 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/labstack/gommon v0.0.0-20180613044413-d6898124de91
	github.com/leebenson/conform v0.0.0-20180615210222-bc2e0311fd85
	github.com/leekchan/accounting v0.0.0-20161211142212-a35854c07593
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/nicksnyder/go-i18n v0.0.0-20180814031359-04f547cc50da
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pact-foundation/pact-go v0.0.13
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.0.6
	github.com/sony/gobreaker v0.0.0-20170530031423-e9556a45379e
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/zemirco/memorystore v0.0.0-20160308183530-ecd57e5134f6
	go.opencensus.io v0.0.0-20180823191657-71e2e3e3082a
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20190108225652-1e06a53dbb7e
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890
	golang.org/x/sys v0.0.0-20181122145206-62eef0e2fa9b // indirect
	google.golang.org/api v0.0.0-20180824000442-943e5aafc110 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.21.1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.1.9 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace (
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b762685799422ab779adefde348535e7a204c363
)

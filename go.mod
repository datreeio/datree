module github.com/datreeio/datree

go 1.19

require (
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/briandowns/spinner v1.12.0
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/kyokomi/emoji v2.2.4+incompatible
	github.com/lithammer/shortuuid v3.0.0+incompatible
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.8.1
	github.com/xeipuuv/gojsonschema v1.2.0

	// 24.1.23: A new version of kubeconform will be released soon, using a different json-schema library.
	// This affects our output when resources fail schema validation. Do not update kubeconform version
	// without consulting the product team
	github.com/yannh/kubeconform v0.4.14
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/apimachinery v0.23.5
	sigs.k8s.io/yaml v1.2.0
)

require github.com/open-policy-agent/opa v0.49.2

require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/a8m/envsubst v1.3.0 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alecthomas/participle/v2 v2.0.0-beta.5 // indirect
	github.com/elliotchance/orderedmap v1.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/goccy/go-yaml v1.9.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/jinzhu/copier v0.3.5 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.8.1 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tklauser/numcpus v0.3.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
)

require (
	github.com/itchyny/gojq v0.12.10
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mikefarah/yq/v4 v4.27.3
	github.com/owenrumney/go-sarif/v2 v2.1.2
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shirou/gopsutil/v3 v3.22.2
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	k8s.io/utils v0.0.0-20220823124924-e9cbc92d1a73
)

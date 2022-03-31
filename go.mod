module alicefr/kubevirt-test

go 1.17

replace (
	github.com/openshift/api => github.com/openshift/api v0.0.0-20191219222812-2987a591a72c
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.1
	k8s.io/client-go => k8s.io/client-go v0.23.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.23.1
	kubevirt.io/api => kubevirt.io/api v0.51.0
	kubevirt.io/client-go => kubevirt.io/client-go v0.51.0
	kubevirt.io/containerized-data-importer-api => kubevirt.io/containerized-data-importer-api v1.42.0
)

require (
	github.com/spf13/cobra v1.4.0
	k8s.io/client-go v12.0.0+incompatible
	kubevirt.io/client-go v0.51.0
	kubevirt.io/kubevirt v0.51.0

)

require (
	github.com/coreos/prometheus-operator v0.38.1-0.20200424145508-7e176fda06cc // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.5.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20191119172530-79f836b90111 // indirect
	github.com/kubernetes-csi/external-snapshotter/v2 v2.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/openshift/api v0.0.0 // indirect
	github.com/openshift/client-go v0.0.0 // indirect
	github.com/openshift/custom-resource-status v0.0.0-20200602122900-c002fd1547ca // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20211209124913-491a49abca63 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.23.1 // indirect
	k8s.io/apiextensions-apiserver v0.23.1 // indirect
	k8s.io/apimachinery v0.23.1 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kube-aggregator v0.23.1 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b // indirect
	kubevirt.io/api v0.0.0-20220308214750-7df276b474a4 // indirect
	kubevirt.io/containerized-data-importer-api v1.42.0 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk v0.2.2 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

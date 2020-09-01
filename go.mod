module github.com/presslabs/controller-util

go 1.14

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/blendle/zapdriver v1.3.1
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/go-logr/zapr v0.2.0
	github.com/go-test/deep v1.0.7
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.11
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200831180312-196b9ba8737a // indirect

	// kubernetes-1.18
	k8s.io/api v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
	sigs.k8s.io/controller-runtime v0.6.2
)

replace gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7

module github.com/presslabs/controller-util

go 1.12

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.0
	github.com/go-test/deep v1.0.2
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.6
	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.4.2
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20180906233101-161cd47e91fd
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0-rc.0
)

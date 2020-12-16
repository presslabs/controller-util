module github.com/presslabs/controller-util

go 1.14

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/blendle/zapdriver v1.3.1
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/go-test/deep v1.0.7
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.11
	github.com/nxadm/tail v1.4.5 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.2
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
)

replace gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7

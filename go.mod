module github.com/presslabs/controller-util

go 1.16

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/blendle/zapdriver v1.3.1
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/go-test/deep v1.0.7
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.11
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	golang.org/x/tools v0.1.0 // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.8.3
)

replace gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7

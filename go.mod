module github.com/presslabs/controller-util

go 1.12

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/go-test/deep v1.0.2
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.8
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.7.1
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1 // indirect
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/oauth2 v0.0.0-20191122200657-5d9234df094c // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect

	// kubernetes-1.16.0
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/utils v0.0.0-20191114200735-6ca3b61696b6 // indirect
	sigs.k8s.io/controller-runtime v0.4.0
)

replace gopkg.in/fsnotify.v1 => gopkg.in/fsnotify.v1 v1.4.7

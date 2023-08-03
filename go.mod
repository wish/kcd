module github.com/wish/kcd

go 1.12

require (
	github.com/DataDog/datadog-go v2.2.0+incompatible
	github.com/aws/aws-sdk-go v1.21.8
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/golang/glog v1.0.0
	github.com/heroku/docker-registry-client v0.0.0-20181004091502-47ecf50fd8d4
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/myesui/uuid v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/twinj/uuid v1.0.0
	goji.io v2.0.2+incompatible
	gopkg.in/stretchr/testify.v1 v1.2.2 // indirect
	k8s.io/api v0.25.0 //kubernetes-1.25
	k8s.io/apiextensions-apiserver v0.25.0 //kubernetes-1.25
	k8s.io/apimachinery v0.25.0 //kubernetes-1.25
	k8s.io/apiserver v0.25.0 //kubernetes-1.25
	k8s.io/client-go v0.25.0 //kubernetes-1.25
)

replace bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d

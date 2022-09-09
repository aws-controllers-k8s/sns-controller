module github.com/aws-controllers-k8s/sns-controller

go 1.14

require (
	github.com/aws/aws-controllers-k8s v0.0.2
	github.com/aws/aws-sdk-go v1.35.5
	github.com/go-logr/logr v1.2.0
	github.com/google/go-cmp v0.5.5
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
)

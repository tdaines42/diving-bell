module github.com/tdaines42/diving-bell

require (
	github.com/mitchellh/go-homedir v1.1.0
	k8s.io/klog v1.0.0
	k8s.io/api v0.0.0-20191121175643-4ed536977f46
	k8s.io/apimachinery v0.0.0-20191121175448-79c2a76c473a
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191121175643-4ed536977f46
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191121175448-79c2a76c473a
)

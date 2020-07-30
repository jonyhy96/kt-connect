module github.com/jonyhy96/kt-connect

go 1.12

require (
	github.com/deckarep/golang-set v1.7.1
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/mock v1.4.1
	github.com/gorilla/websocket v1.4.1
	github.com/kubernetes/dashboard v1.10.1
	github.com/lextoumbourou/goodhosts v2.1.0+incompatible
	github.com/miekg/dns v0.0.0-20190106042521-5beb9624161b
	github.com/rs/zerolog v0.0.0-20190704061603-77a169535877
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/urfave/cli v0.0.0-20190203184040-693af58b4d51
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d
	golang.org/x/sys v0.0.0-20200413165638-669c56c373c4 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	istio.io/api v0.0.0-20200221025927-228308df3f1b
	istio.io/client-go v0.0.0-20200221055756-736d3076b458
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/cli-runtime v0.17.2
	k8s.io/client-go v0.17.2
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43

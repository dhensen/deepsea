package k8s

import (
	"flag"
	"time"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K8SClientSet *kubernetes.Clientset
var K8SDiscoveryClient *discovery.DiscoveryClient

func init() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	var config *rest.Config
	var err error
	if *kubeconfig == "" {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		// creates the out of cluster config based on .kube/config file
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
	}

	// set the http client timeout to half a second
	config.Timeout = 500 * time.Millisecond
	K8SDiscoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	K8SClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

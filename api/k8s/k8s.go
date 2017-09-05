package k8s

import (
	"flag"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8SClientSet *kubernetes.Clientset
var K8SDiscoveryClient *discovery.DiscoveryClient

func init() {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	if *kubeconfig == "" {
		panic("-kubeconfig not specified")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	K8SDiscoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	K8SClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

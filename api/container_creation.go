package main

import (
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/google/uuid"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func AddContainer(w http.ResponseWriter, r *http.Request) {

	presetId := r.FormValue("preset_id")
	imagePreset := imagePresets[presetId]
	dnsLabel := fmt.Sprintf("wordpress-%s", uuid.New().String())

	// I'm using the dnsLabel for:
	// 1. .metadata.name
	// 2. .metadata.labels.app
	// 3. .spec.template.spec.containers.name
	// TODO: properly differentiate these three vars

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dnsLabel,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"run":     dnsLabel,
						"deepsea": "v0",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dnsLabel,
							Image: imagePreset.Image,
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "WORDPRESS_TABLE_PREFIX",
									Value: strings.Replace(dnsLabel, "-", "_", -1),
								},
								{
									Name:  "WORDPRESS_DB_HOST",
									Value: "mysql-1",
								},
								{
									Name:  "WORDPRESS_DB_USER",
									Value: "root",
								},
								{
									Name:  "WORDPRESS_DB_PASSWORD",
									Value: "root",
								},
							},
						},
					},
				},
			},
		},
	}

	clientset := getClientSet()
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Creating deployment...")
	_, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}

	// TODO: expose deployment => create service

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dnsLabel,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: util.FromInt(80),
				},
			},
			Selector: map[string]string{
				"run": dnsLabel,
			},
		},
	}

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	fmt.Println("Creating Service...")
	_, err = serviceClient.Create(service)
	if err != nil {
		panic(err)
	}

	ingress := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: dnsLabel,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: fmt.Sprintf("%s.local", dnsLabel),
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: dnsLabel,
										ServicePort: util.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	ingressClient := clientset.ExtensionsV1beta1().Ingresses(apiv1.NamespaceDefault)
	fmt.Println("Creating Ingress...")
	_, err = ingressClient.Create(ingress)
	if err != nil {
		panic(err)
	}
}

func getClientSet() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func int32Ptr(i int32) *int32 { return &i }

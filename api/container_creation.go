package main

import (
	"fmt"

	uuid "github.com/google/uuid"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/rest"
)

func AddContainer(id string) {

	imagePreset := imagePresets[id]
	dns_label := fmt.Sprintf("wordpress-%s", uuid.New().String())

	// I'm using the dns_label for:
	// 1. .metadata.name
	// 2. .metadata.labels.app
	// 3. .spec.template.spec.containers.name
	// TODO: properly differentiate these three vars

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: dns_label,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": dns_label,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  dns_label,
							Image: imagePreset.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }

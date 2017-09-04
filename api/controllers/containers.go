package controllers

import (
	"encoding/json"
	"flag"
	"fmt"
	"local/deepsea/api/models"
	"log"
	"net/http"
	"os/exec"
	"strings"

	uuid "github.com/google/uuid"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string

func init() {
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	if *kubeconfig == "" {
		panic("-kubeconfig not specified")
	}
}

func ListContainers(w http.ResponseWriter, r *http.Request) {

}

func AddContainer(w http.ResponseWriter, r *http.Request) {

	// TODO:  preset_id must be valid
	presetID := r.FormValue("preset_id")
	if presetID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "preset_id is required",
		})
		return
	}
	domainName := r.FormValue("domain_name")
	if domainName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "domain_name is required",
		})
		return
	}

	domain := models.Domain{
		UUID:     uuid.New().String(),
		Name:     domainName,
		Provider: "k8s-dev",
	}

	imagePreset := imagePresets[presetID]
	dnsLabel := strings.Replace(domainName, ".", "-", -1)

	// I'm using the dnsLabel for:
	// 1. .metadata.name
	// 2. .metadata.labels.app
	// 3. .spec.template.spec.containers.name
	// TODO: properly differentiate these three vars

	clientset := getClientSet()

	persistentVolume := CreatePersistentVolume(dnsLabel)
	volumeClient := clientset.CoreV1().PersistentVolumes()
	_, err := volumeClient.Create(persistentVolume)
	if err != nil {
		panic(err)
	}

	persistentVolumeClaim := CreatePersistentVolumeClaim(dnsLabel)
	volumeClaimClient := clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)
	_, err = volumeClaimClient.Create(persistentVolumeClaim)
	if err != nil {
		panic(err)
	}

	deployment := CreateDeployment(imagePreset.Image, dnsLabel, dnsLabel, dnsLabel, dnsLabel)

	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	log.Println("Creating deployment...")
	_, err = deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	log.Println("Deployment created succesfully!")

	service := CreateService(dnsLabel, dnsLabel)
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	log.Println("Creating Service...")
	_, err = serviceClient.Create(service)
	if err != nil {
		panic(err)
	}
	log.Println("Service created succesfully!")

	ingress := CreateIngress(dnsLabel, domainName, dnsLabel)
	ingressClient := clientset.ExtensionsV1beta1().Ingresses(apiv1.NamespaceDefault)
	log.Println("Creating Ingress...")
	_, err = ingressClient.Create(ingress)
	if err != nil {
		panic(err)
	}
	log.Println("Ingress created succesfully!")

	json.NewEncoder(w).Encode(domain)
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

func CreatePersistentVolume(volumeName string) *apiv1.PersistentVolume {
	hostPath := fmt.Sprintf("/data/volumes/%s", volumeName)

	// Call mkdir hostPath via ssh in minikube
	cmd := exec.Command("/usr/local/bin/minikube", "ssh", fmt.Sprintf("mkdir -pv %s", hostPath))
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	log.Printf("%s", out)
	// TODO: handle failure

	return &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volumeName,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			StorageClassName: "manual",
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Capacity: apiv1.ResourceList{
				apiv1.ResourceStorage: resource.MustParse("500Mi"),
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: hostPath,
				},
			},
		},
	}
}

func CreatePersistentVolumeClaim(pvcName string) *apiv1.PersistentVolumeClaim {
	storageClassName := "manual"
	storageClassPointer := &storageClassName
	return &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvcName,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			StorageClassName: storageClassPointer,
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("500Mi"),
				},
			},
		},
	}
}

func CreateDeployment(image string, metadataName string, runLabelValue string, containerName string, pvcName string) *appsv1beta1.Deployment {
	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: metadataName,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"run":     runLabelValue,
						"deepsea": "v0",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  containerName,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "WORDPRESS_TABLE_PREFIX",
									Value: strings.Replace(metadataName, "-", "_", -1),
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
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/var/www/html",
									Name:      pvcName,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: pvcName,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func CreateService(metadataName string, runSelectorValue string) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: metadataName,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"run": runSelectorValue,
			},
		},
	}
}

func CreateIngress(metadataName string, domainName string, serviceName string) *extensionsv1beta1.Ingress {
	return &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: metadataName,
			Labels: map[string]string{
				"deepsea": "v0",
			},
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: domainName,
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

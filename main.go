package main

import (
	"context"
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// getAllResource 获取资源
func getAllResource(clientset *kubernetes.Clientset) {
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodeList.Items {
		fmt.Println(node.Name)

		for _, nodeCondition := range node.Status.Conditions {
			if nodeCondition.Type == "Ready" && nodeCondition.Status == "True" {
				fmt.Println("Ready")
			}
		}

		// jsonBytes, err := json.Marshal(node.Status)
		// if err != nil {
		// 	panic(err.Error())
		// }

		// fmt.Println(string(jsonBytes))

		for _, nodeAddress := range node.Status.Addresses {
			fmt.Println(nodeAddress.Type)
			fmt.Println(nodeAddress.Address)
		}

		fmt.Println(node.Status.NodeInfo.OSImage)
		fmt.Println(node.Status.NodeInfo.ContainerRuntimeVersion)
		fmt.Println(node.Status.NodeInfo.KubeletVersion)
		fmt.Println(node.Status.NodeInfo.Architecture)
		fmt.Println(node.Status.NodeInfo.OperatingSystem)
	}

	nameSpaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	for _, nameSpace := range nameSpaceList.Items {

		fmt.Println(nameSpace.Name)
		fmt.Println(nameSpace.CreationTimestamp)
		fmt.Println(nameSpace.Status.Phase)

		// podList, err := clientset.CoreV1().Pods(nameSpace.Name).List(context.TODO(), metav1.ListOptions{})

		// if err != nil {
		// 	panic(err.Error())
		// }

		// for _, pod := range podList.Items {
		// 	fmt.Println(pod.Name)
		// }

	}

	svcList, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, svc := range svcList.Items {
		fmt.Println(svc.Name)
		fmt.Println(svc.Namespace)
		fmt.Println(svc.Spec.Type)
		fmt.Println(svc.CreationTimestamp)
		fmt.Println(svc.Spec.ClusterIP)
		fmt.Println(svc.Spec.Ports)

		for _, svcePort := range svc.Spec.Ports {
			fmt.Println(svcePort.Name)
			fmt.Println(svcePort.Protocol)
			fmt.Println(svcePort.AppProtocol)
			fmt.Println(svcePort.Port)
			fmt.Println(svcePort.TargetPort)
			fmt.Println(svcePort.NodePort)
		}
	}

	deploymentList, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, deployment := range deploymentList.Items {
		fmt.Println(deployment.Name)
		fmt.Println(deployment.Namespace)
		for key, value := range deployment.Labels {
			fmt.Println(key, value)
		}
		for key, value := range deployment.Spec.Selector.MatchLabels {
			fmt.Println(key, value)
		}
		fmt.Println(deployment.Status.Replicas)
		fmt.Println(deployment.Status.AvailableReplicas)
		fmt.Println(deployment.CreationTimestamp)
	}

	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, pod := range podList.Items {

		// jsonBytes, err := json.Marshal(pod.Status)
		// if err != nil {
		// 	log.Fatal(err)
		// 	return
		// }

		// fmt.Println(string(jsonBytes))

		fmt.Printf("podName: %s \t podStatus: %s \t podMsg: %s \t podIP: %s \t nodeName: %s \t namespace: %s\n",
			pod.Name,
			pod.Status.Phase,
			pod.Status.Message,
			pod.Status.PodIP,
			pod.Spec.NodeName,
			pod.Namespace,
		)
	}
}

// createDeployment 创建deployment
func createDeployment(clientset *kubernetes.Clientset) {

	nameSpace := "default"

	var replicas int32 = 5

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "go-client-nginx",
			// deployment 的标签
			Labels: map[string]string{
				"app": "go-client-nginx",
				"env": "go-client-dev",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				// deployment选择器选择的标签，对应pod的标签
				MatchLabels: map[string]string{
					"app": "go-client-nginx",
					"env": "go-client-dev",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "go-client-nginx",
					// pod的标签
					Labels: map[string]string{
						"app": "go-client-nginx",
						"env": "go-client-dev",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							// 容器名称
							Name:  "go-client-nginx-container",
							Image: "hub.atguigu.com/library/nginx:latest",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// deployment.ObjectMeta.Name = "go-client-nginx"
	// deployment.ObjectMeta.Labels = map[string]string{
	// 	"app": "go-client-nginx",
	// 	"env": "go-client-dev",
	// }

	// deployment.Spec.Replicas = &replicas

	// deployment.Spec.Selector = &metav1.LabelSelector{
	// 	MatchLabels: map[string]string{
	// 		"app": "go-client-nginx",
	// 		"env": "go-client-dev",
	// 	},
	// }

	// deployment.Spec.Template = corev1.PodTemplateSpec{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "go-client-nginx",
	// 		Labels: map[string]string{
	// 			"app": "go-client-nginx",
	// 			"env": "go-client-dev",
	// 		},
	// 	},
	// 	Spec: corev1.PodSpec{
	// 		Containers: []corev1.Container{
	// 			{
	// 				// 容器名称
	// 				Name:  "go-client-nginx-container",
	// 				Image: "hub.atguigu.com/library/nginx:latest",
	// 				Ports: []corev1.ContainerPort{
	// 					{
	// 						Name:          "http",
	// 						Protocol:      corev1.ProtocolTCP,
	// 						ContainerPort: 80,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	deployment, err := clientset.AppsV1().Deployments(nameSpace).Create(context.TODO(), deployment, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	jsonBytes, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(jsonBytes))

}

// createService 创建service
func createService(clientset *kubernetes.Clientset) {

	nameSpace := "default"

	service := &corev1.Service{}

	service.ObjectMeta.Name = "client-go-service-v3"
	service.Spec.Type = corev1.ServiceTypeNodePort
	// service选择器，选择pod
	service.Spec.Selector = map[string]string{
		"app": "go-client-nginx",
		"env": "go-client-dev",
	}

	service.Spec.Ports = []corev1.ServicePort{}

	servicePort := corev1.ServicePort{}

	servicePort.Name = "http"
	servicePort.Port = 80
	servicePort.Protocol = corev1.ProtocolTCP
	servicePort.TargetPort = intstr.Parse("80")

	service.Spec.Ports = append(service.Spec.Ports, servicePort)

	service, err := clientset.CoreV1().Services(nameSpace).Create(context.TODO(), service, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	jsonBytes, err := json.Marshal(service)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(jsonBytes))

}

func deleteService(clientset *kubernetes.Clientset) {
	nameSpace := "default"
	ServiceName := "client-go-service"

	err := clientset.CoreV1().Services(nameSpace).Delete(context.TODO(), ServiceName, metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}

}

func updateDeployment(clientset *kubernetes.Clientset) {

	nameSpace := "default"

	deploymentName := "go-client-nginx"

	deployment, err := clientset.AppsV1().Deployments(nameSpace).Get(context.TODO(), deploymentName, metav1.GetOptions{})

	var replicas int32 = 5

	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:v1"
	deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = corev1.PullIfNotPresent

	deployment, err = clientset.AppsV1().Deployments(nameSpace).Update(context.TODO(), deployment, metav1.UpdateOptions{})

	jsonBytes, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(jsonBytes))

}

func updateDeploymentScale(clientset *kubernetes.Clientset) {

	nameSpace := "default"

	deploymentName := "go-client-nginx"

	scale, err := clientset.AppsV1().Deployments(nameSpace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})

	scale.Spec.Replicas = 1

	scale, err = clientset.AppsV1().Deployments(nameSpace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})

	jsonBytes, err := json.Marshal(scale)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(jsonBytes))

}

func main() {
	var configPath = "config"
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// getAllResource(clientset)

	// createDeployment(clientset)

	// createService(clientset)

	deleteService(clientset)

	// updateDeployment(clientset)

	// updateDeploymentScale(clientset)

}

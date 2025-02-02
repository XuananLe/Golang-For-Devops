package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var deploymentName = "k8s-go-demo"

func getClient() (*kubernetes.Clientset, error) {
	var kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Fuck")
		return nil, err
	}
	return clientset, nil
}

func delete(c *kubernetes.Clientset, deploymentName string) error {
	return c.AppsV1().Deployments("default").Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
}

func checkDeployment(c *kubernetes.Clientset) (bool, error) {
	deployment, err := c.AppsV1().Deployments(apiv1.NamespaceDefault).Get(context.Background(), deploymentName, metav1.GetOptions{});
	if err != nil {
		return false, errors.New(err.Error());
	}
	if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
		return true, nil
	} 
	return false, nil;
}

func deploy(c *kubernetes.Clientset) (map[string]string, error) {
	deployment := &v1.Deployment{}
	appFile, err := os.ReadFile("deployment.yaml")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	switch obj := obj.(type) {
	case *v1.Deployment:
		deployment = obj
	default:
		return nil, fmt.Errorf("UNRECOGNIZED TYPE: %+v", groupVersionKind)
	}

	_, err = c.AppsV1().Deployments("default").Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err == nil {
		log.Println("Deployment k8s-go-demo exists, We have to delete it LMAO")
		delete(c, deploymentName)
	} else if k8sErrors.IsNotFound(err) {
		log.Println("Deployment k8s-go-demo Not exists, gotta create it first LMAO")
	}

	deploymentResponse, err := c.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for {
		ready, _ := checkDeployment(c);
		if ready {
			fmt.Println("All pods are ready! ✅")
			break
		} else {
			fmt.Println("Pods are not ready yet, waiting...")
			time.Sleep(5 * time.Second)
		}
	}
	return deploymentResponse.Spec.Template.Labels, nil
}

func main() {
	var (
		client *kubernetes.Clientset
		err    error
	)
	if client, err = getClient(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(deploy(client))
}

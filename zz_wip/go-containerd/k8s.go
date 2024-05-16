package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

func k8s(kc string) error {

	// kubeconfig := *kc

	// if kubeconfig == "" {
	// 	log.Printf("using in-cluster configuration")
	// 	config, err = rest.InClusterConfig()
	// } else {
	// 	log.Printf("using configuration from '%s'", kubeconfig)
	// 	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	// }

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	crdConfig := *config
	crdGroup := "openbayes.com"
	crdVersion := "v1alpha1"
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: crdGroup, Version: crdVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	crdClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println(clientset)
	fmt.Println(crdClient)

	// result := &BayesJobList{}
	// ctx := context.Background()
	// crdClient.Get().Namespace("default").Resource("bayesjob").Do(ctx).Into(result)

	// for {
	// 	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// 	// Examples for error handling:
	// 	// - Use helper functions like e.g. errors.IsNotFound()
	// 	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	// 	namespace := "default"
	// 	pod := "example-xxxxx"
	// 	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	// 	if errors.IsNotFound(err) {
	// 		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	// 	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	// 		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
	// 			pod, namespace, statusError.ErrStatus.Message)
	// 	} else if err != nil {
	// 		panic(err.Error())
	// 	} else {
	// 		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	// 	}

	// 	time.Sleep(10 * time.Second)
	// }

	return nil
}

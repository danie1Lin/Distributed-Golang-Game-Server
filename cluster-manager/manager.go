package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {

	//config, err := clientcmd.
	//Check kubernetes status
	config, err := rest.InClusterConfig()
	fmt.Print(config)

	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	log.Print(clientset)
	if err != nil {
		log.Fatal(err)
	}

	api := clientset.CoreV1()

	log.Print(api)
}

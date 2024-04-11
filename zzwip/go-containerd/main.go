package main

import (
	"log"
)

func main() {

	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()

	// if err := k8s(kubeconfig); err != nil {
	// 	log.Fatal(err)
	// }

	if err := k8s(""); err != nil {
		log.Fatal(err)
	}
}

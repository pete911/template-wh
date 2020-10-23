package main

import (
	"github.com/pete911/template-wh/pkg/k8s"
	"github.com/pete911/template-wh/pkg/server"
	"log"
)

func main() {

	flags := ParseFlags()
	log.Printf("starting template admission webhook with flags: %+v", flags)

	values, err := getConfigMapData(flags.Kubeconfig, flags.ConfigmapNamespace, flags.ConfigmapName)
	if err != nil {
		log.Fatalf("get configmap values: %v", err)
	}

	var mutateFn = func(body []byte) ([]byte, error) {
		return k8s.Mutate(body, values)
	}
	log.Fatal(server.ListenAndServeTLS(mutateFn, flags.ServerCertFile, flags.ServerKeyFile))
}

func getConfigMapData(kubeconfigPath, namespace, name string) (map[string]string, error) {

	kubeconfig, err := k8s.LoadKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("get kubeconfig: %v", err)
	}
	return k8s.NewClient(kubeconfig.Clientset).GetConfigMapData(namespace, name)
}

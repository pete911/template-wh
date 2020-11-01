package main

import (
	"github.com/pete911/template-wh/pkg/k8s"
	"github.com/pete911/template-wh/pkg/server"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	tlsCertFile string
	tlsKeyFile  string
)

func init() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("get user home dir: %v", err)
	}
	tlsCertFile = filepath.Join(homeDir, "tls.crt")
	tlsKeyFile = filepath.Join(homeDir, "tls.key")
}

func main() {

	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("parse flags: %v", err)
	}
	log.Printf("starting template admission webhook with flags: %s", flags)

	if err := ioutil.WriteFile(tlsCertFile, []byte(flags.TLSCrt), 0640); err != nil {
		log.Fatalf("write tls.crt: %v", err)
	}
	if err := ioutil.WriteFile(tlsKeyFile, []byte(flags.TLSKey), 0600); err != nil {
		log.Fatalf("write tls.key: %v", err)
	}

	k8sClient := getK8sClient(flags.Kubeconfig)
	var mutateFn = func(body []byte) ([]byte, error) {
		values := getValues(k8sClient, flags.ConfigmapNamespace, flags.ConfigmapName)
		return k8s.Mutate(body, values)
	}
	log.Fatal(server.ListenAndServeTLS(mutateFn, tlsCertFile, tlsKeyFile))
}

func getValues(k8sClient k8s.Client, namespace, name string) map[string]string {

	values, err := k8sClient.GetConfigMapData(namespace, name)
	if err != nil {
		log.Fatalf("get configmap values: %v", err)
	}
	return values
}

func getK8sClient(kubeconfigPath string) k8s.Client {

	log.Print("loading kubeconfig")
	kubeconfig, err := k8s.LoadKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("get kubeconfig: %v", err)
	}
	return k8s.NewClient(kubeconfig.Clientset)
}

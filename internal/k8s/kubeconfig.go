package k8s

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type Kubeconfig struct {
	Host      string
	CA        []byte
	Clientset *kubernetes.Clientset
}

func LoadKubeconfig(kubeconfigPath string) (Kubeconfig, error) {

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return Kubeconfig{}, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return Kubeconfig{}, err
	}

	ca, err := getCA(restConfig.TLSClientConfig)
	if err != nil {
		return Kubeconfig{}, err
	}

	return Kubeconfig{
		Host:      restConfig.Host,
		CA:        ca,
		Clientset: clientset,
	}, nil
}

func getCA(tls rest.TLSClientConfig) ([]byte, error) {

	if tls.CAFile != "" {
		return os.ReadFile(tls.CAFile)
	}
	if len(tls.CAData) != 0 {
		return tls.CAData, nil
	}
	return nil, errors.New("cannot find CA file or CA data in tls client config")
}

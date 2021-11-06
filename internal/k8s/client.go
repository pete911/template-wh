package k8s

import (
	"context"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	certificates "k8s.io/client-go/kubernetes/typed/certificates/v1"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	clientSet        *kubernetes.Clientset
	configMapsGetter core.ConfigMapsGetter
	csr              certificates.CertificateSigningRequestInterface
}

func NewClient(clientSet *kubernetes.Clientset) Client {

	clientSet.CertificatesV1()
	return Client{
		clientSet:        clientSet,
		configMapsGetter: clientSet.CoreV1(),
		csr:              clientSet.CertificatesV1().CertificateSigningRequests(),
	}
}

func (c Client) GetConfigMapData(namespace, name string) (map[string]string, error) {

	cm, err := c.clientSet.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, meta.GetOptions{})
	if err != nil || cm == nil {
		return nil, err
	}
	return cm.Data, nil
}

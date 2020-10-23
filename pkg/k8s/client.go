package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
)

// --- stripped down kubernetes interfaces to simplify testing ---

type configMapsInterface interface {
	Get(ctx context.Context, name string, opts meta.GetOptions) (*v1.ConfigMap, error)
}

type configMapsGetter interface {
	ConfigMaps(namespace string) configMapsInterface
}

type configMaps struct {
	getter core.ConfigMapsGetter
}

func (c configMaps) ConfigMaps(namespace string) configMapsInterface {
	return c.getter.ConfigMaps(namespace)
}

// --- ------------------------------------------------------- ---

type Client struct {
	configMapsGetter      configMapsGetter
}

func NewClient(clientSet *kubernetes.Clientset) Client {

	return Client{
		configMapsGetter: configMaps{getter: clientSet.CoreV1()},
	}
}

func (c Client) GetConfigMapData(namespace, name string) (map[string]string, error) {

	cm, err := c.configMapsGetter.ConfigMaps(namespace).Get(context.Background(), name, meta.GetOptions{})
	if err != nil || cm == nil {
		return nil, err
	}
	return cm.Data, nil
}

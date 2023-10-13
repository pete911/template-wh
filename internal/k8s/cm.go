package k8s

import (
	"context"
	"fmt"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetConfigMapData(clientSet *kubernetes.Clientset, namespace, name string) (map[string]string, error) {

	cm, err := clientSet.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, meta.GetOptions{})
	if err != nil {
		return nil, err
	}
	if cm == nil {
		return nil, fmt.Errorf("nil %s configmap ", name)
	}
	return cm.Data, nil
}

package main

import (
	"errors"
	"fmt"
	"github.com/pete911/template-wh/internal/k8s"
	"github.com/pete911/template-wh/internal/server"
	"log/slog"
	"net/http"
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
		slog.Error(fmt.Sprintf("get user home dir: %v", err))
		os.Exit(1)
	}
	tlsCertFile = filepath.Join(homeDir, "tls.crt")
	tlsKeyFile = filepath.Join(homeDir, "tls.key")
}

func main() {

	flags, err := ParseFlags()
	if err != nil {
		slog.Error(fmt.Sprintf("parse flags: %v", err))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("starting template admission webhook with flags: %s", flags))

	if err := os.WriteFile(tlsCertFile, []byte(flags.TLSCrt), 0640); err != nil {
		slog.Error(fmt.Sprintf("write tls.crt: %v", err))
		os.Exit(1)
	}
	if err := os.WriteFile(tlsKeyFile, []byte(flags.TLSKey), 0600); err != nil {
		slog.Error(fmt.Sprintf("write tls.key: %v", err))
		os.Exit(1)
	}

	k8sClient := getK8sClient(flags.Kubeconfig)
	var mutateFn = func(body []byte) ([]byte, error) {
		values := getValues(k8sClient, flags.ConfigmapNamespace, flags.ConfigmapName)
		return k8s.Mutate(body, values)
	}

	if err := server.ListenAndServeTLS(mutateFn, tlsCertFile, tlsKeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("server closed")
}

func getValues(k8sClient k8s.Client, namespace, name string) map[string]string {

	values, err := k8sClient.GetConfigMapData(namespace, name)
	if err != nil {
		slog.Error(fmt.Sprintf("get configmap values: %v", err))
		os.Exit(1)
	}
	return values
}

func getK8sClient(kubeconfigPath string) k8s.Client {

	slog.Info("loading kubeconfig")
	kubeconfig, err := k8s.LoadKubeconfig(kubeconfigPath)
	if err != nil {
		slog.Error(fmt.Sprintf("get kubeconfig: %v", err))
		os.Exit(1)
	}
	return k8s.NewClient(kubeconfig.Clientset)
}

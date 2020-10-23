package main

import (
	"flag"
	"os"
)

type Flags struct {
	Kubeconfig         string
	ConfigmapName      string
	ConfigmapNamespace string
	ServerCertFile     string
	ServerKeyFile      string
}

func ParseFlags() Flags {

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	kubeconfig := f.String("kubeconfig", getStringEnv("KUBECONFIG", ""), "path to kubeconfig file, or empty for in-cluster kubeconfig")
	configmapName := f.String("configmap-name", getStringEnv("TWH_CONFIGMAP_NAME", "template-wh"), "name of the configmap that has template values")
	configmapNamespace := f.String("configmap-namespace", getStringEnv("TWH_CONFIGMAP_NAMESPACE", "kube-system"), "namespace of the configmap that has template values")
	serverCertFile := f.String("server-cert-file", getStringEnv("TWH_SERVER_CERT_FILE", "/etc/template-wh/ssl/cert.pem"), "cert file used by template admission webhook server")
	serverKeyFile := f.String("server-key-file", getStringEnv("TWH_SERVER_KEY_FILE", "/etc/template-wh/ssl/key.pem"), "key file used by template admission webhook server")
	f.Parse(os.Args[1:])

	return Flags{
		Kubeconfig:         stringValue(kubeconfig),
		ConfigmapName:      stringValue(configmapName),
		ConfigmapNamespace: stringValue(configmapNamespace),
		ServerCertFile:     stringValue(serverCertFile),
		ServerKeyFile:      stringValue(serverKeyFile),
	}
}

func getStringEnv(envName string, defaultValue string) string {

	env, ok := os.LookupEnv(envName)
	if !ok {
		return defaultValue
	}
	return env
}

func stringValue(v *string) string {

	if v == nil {
		return ""
	}
	return *v
}

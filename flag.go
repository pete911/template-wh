package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	Kubeconfig         string
	ConfigmapName      string
	ConfigmapNamespace string
	TLSCrt             string
	TLSKey             string
}

func (f Flags) String() string {

	return fmt.Sprintf("kubeconfig: %q configmap-name: %q configmap-namespace: %q tls-crt: %q tls-key: %q",
		f.Kubeconfig, f.ConfigmapName, f.ConfigmapNamespace, "****", "****")
}

func ParseFlags() (Flags, error) {

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	kubeconfig := f.String("kubeconfig", getStringEnv("KUBECONFIG", ""),
		"path to kubeconfig file, or empty for in-cluster kubeconfig")
	configmapName := f.String("configmap-name", getStringEnv("TWH_CONFIGMAP_NAME", "template-wh"),
		"name of the configmap that has template values")
	configmapNamespace := f.String("configmap-namespace", getStringEnv("TWH_CONFIGMAP_NAMESPACE", "kube-system"),
		"namespace of the configmap that has template values")
	tlsCrt := f.String("tls-crt", getStringEnv("TWH_TLS_CRT", ""),
		"tls certificate to be used by this service")
	tlsKey := f.String("tls-key", getStringEnv("TWH_TLS_KEY", ""),
		"tls key to be used by this service")
	f.Parse(os.Args[1:])

	flags := Flags{
		Kubeconfig:         stringValue(kubeconfig),
		ConfigmapName:      stringValue(configmapName),
		ConfigmapNamespace: stringValue(configmapNamespace),
		TLSCrt:             stringValue(tlsCrt),
		TLSKey:             stringValue(tlsKey),
	}

	if flags.TLSKey == "" || flags.TLSCrt == "" {
		return Flags{}, errors.New("required flags tls-crt and tls-key are not set")
	}
	return flags, nil
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

package kube

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var kubeConfigPath string

type KubeConfig struct {
	ConfigPath string
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	kubeConfigPath = filepath.Join(home, ".kube", "config")
}

func (kc *KubeConfig) getConfig() *api.Config {
	if kc.ConfigPath != "" {
		kubeConfigPath = kc.ConfigPath
	}
	config, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatalf("The config could not be loaded from %s. Confirm the file exists.", kubeConfigPath)
	}
	return config
}

func (kc *KubeConfig) GetContexts() []string {
	var contexts []string

	config := kc.getConfig()

	if len(config.Contexts) == 0 {
		log.Fatalf("There are no contexts defined for this config")
	}

	for k, v := range config.Contexts {
		log.Printf("Context %s - cluster: %s", k, v.Cluster)
		contexts = append(contexts, k)
	}

	return contexts
}

func (kc *KubeConfig) GetCurrentContext() string {
	config := kc.getConfig()

	if config.CurrentContext == "" {
		log.Fatalf("CurrentConfig is not configured on config")
	}

	return config.CurrentContext
}

func (kc *KubeConfig) SetCurrentContext(context string) {
	config := kc.getConfig()
	_, ok := config.Contexts[context]
	if !ok {
		log.Fatalf("The context %s does not exist in the config", context)
	}
	config.CurrentContext = context
}

func GetClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to build config from %q", kubeConfigPath)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client")
	}

	return client
}

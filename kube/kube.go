package kube

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfig struct {
	ConfigPath string
	Config     *api.Config
}

func NewKubeConfig(configPath string) *KubeConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kubeConfigPath := filepath.Join(home, ".kube", "config")
	if configPath != "" {
		kubeConfigPath = configPath
	}

	config, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatalf("The config could not be loaded from %s. Confirm the file exists.", kubeConfigPath)
	}

	return &KubeConfig{ConfigPath: kubeConfigPath, Config: config}
}

func (kc *KubeConfig) GetContexts() []string {
	var contexts []string

	if len(kc.Config.Contexts) == 0 {
		log.Fatalf("There are no contexts defined for this config")
	}

	for k := range kc.Config.Contexts {
		contexts = append(contexts, k)
	}

	return contexts
}

func (kc *KubeConfig) GetCurrentContext() string {
	if kc.Config.CurrentContext == "" {
		log.Fatalf("CurrentConfig is not configured on config")
	}

	return kc.Config.CurrentContext
}

func (kc *KubeConfig) SetCurrentContext(context string) {
	_, ok := kc.Config.Contexts[context]
	if !ok {
		log.Fatalf("The context %s does not exist in the config", context)
	}
	kc.Config.CurrentContext = context
	err := clientcmd.WriteToFile(*kc.Config, kc.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to write configuration to %s", kc.ConfigPath)
	}
}

func GetClient(kubeConfig *KubeConfig) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to build config from %q", kubeConfig.ConfigPath)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client")
	}

	return client
}

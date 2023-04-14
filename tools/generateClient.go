package tools

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

type ClientManager struct {
	Client        kubernetes.Interface
	DynamicClient dynamic.Interface
}

func AuthToClusterByKubeConfigFile(kubeConfigFilePath string) (*rest.Config, error) {
	cfg := new(rest.Config)
	var err error
	cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfigFilePath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func AuthToClusterByKubeConfigContent(kubeConfigContent []byte) (cfg *rest.Config, err error) {
	var cfgGetter clientcmd.KubeconfigGetter = func() (*clientcmdapi.Config, error) {
		var tmp = new(clientcmdapi.Config)
		err = yaml.Unmarshal(kubeConfigContent, tmp)
		if err != nil {
			return nil, err
		}
		return tmp, nil
	}
	return clientcmd.BuildConfigFromKubeconfigGetter("", cfgGetter)
}

func GenerateClient(cfg *rest.Config) (cm *ClientManager, err error) {
	cm = &ClientManager{}
	cm.Client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	cm.DynamicClient, err = dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cm, err
}

package kubeconfig

import (
	"backend/Models"
	"encoding/json"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

type KubeCluster struct {
	Name    string      `json:"name"`
	Cluster api.Cluster `json:"cluster"`
}

type KubeContext struct {
	Name    string      `json:"name"`
	Context api.Context `json:"context"`
}

type KubeUser struct {
	Name string       `json:"name"`
	User api.AuthInfo `json:"user"`
}

type KubeConfig struct {
	ApiVersion     string        `json:"apiVersion"`
	Kind           string        `json:"kind"`
	CurrentContext string        `json:"current-context"`
	Clusters       []KubeCluster `json:"clusters"`
	Contexts       []KubeContext `json:"contexts"`
	Users          []KubeUser    `json:"users"`
}

// LoadKubeConfigFileToJSON 加载kubeconfig文件,转换为KubeConfig结构体
func (kc *KubeConfig) LoadKubeConfigFileToJSON(kubeConfigPath string) (*KubeConfig, error) {
	// 读取kubeconfig yaml文件
	f, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	// 将yaml转为json
	j, err := yaml.YAMLToJSON(f)
	if err != nil {
		return nil, err
	}
	// 将json转换为kubeConfig结构体
	err = json.Unmarshal(j, kc)
	if err != nil {
		return nil, err
	}
	return kc, err
}

// ParseKubeConfigToRestConfig 通过KubeConfig结构体生成rest.Config（自定义currentContext）
func (kc *KubeConfig) ParseKubeConfigToRestConfig(currentContext string) (*rest.Config, error) {
	var cfg api.Config
	var configAccess clientcmd.ConfigAccess

	// 将*kubeConfig转换为api.Config
	cfg.Kind = kc.Kind
	cfg.APIVersion = kc.ApiVersion
	cfg.CurrentContext = kc.CurrentContext
	cfg.Contexts = make(map[string]*api.Context)
	cfg.Clusters = make(map[string]*api.Cluster)
	cfg.AuthInfos = make(map[string]*api.AuthInfo)
	for _, item := range kc.Contexts {
		cfg.Contexts[item.Name] = &item.Context
	}
	for _, item := range kc.Clusters {
		cfg.Clusters[item.Name] = &item.Cluster
	}
	for _, item := range kc.Users {
		cfg.AuthInfos[item.Name] = &item.User
	}
	// !!!!!!!参考client-go中的clientcmd包，还有NewNonInteractiveClientConfig和NewInteractiveClientConfig支持传入更多参数，比如currentcontext等
	//return clientcmd.NewDefaultClientConfig(cfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	return clientcmd.NewNonInteractiveClientConfig(cfg, currentContext, &clientcmd.ConfigOverrides{}, configAccess).ClientConfig()
}

// ParseKubeConfigToCluster 转换KubeConfig结构体为Models.Cluster结构体类型，用于存储到数据库中
func (kc *KubeConfig) ParseKubeConfigToCluster() (clusters []Models.Cluster) {
	for _, ctx := range kc.Contexts {
		var cluster Models.Cluster
		cluster.APIVersion = kc.ApiVersion
		cluster.Kind = kc.Kind
		cluster.ContextName = ctx.Name
		cluster.Cluster = ctx.Context.Cluster
		cluster.AuthInfo = ctx.Context.AuthInfo
		for _, clt := range kc.Clusters {
			if clt.Name == cluster.Cluster {
				cluster.Server = clt.Cluster.Server
				cluster.InsecureSkipTLSVerify = clt.Cluster.InsecureSkipTLSVerify
				cluster.CertificateAuthorityData = string(clt.Cluster.CertificateAuthorityData)
				break
			}
		}
		for _, user := range kc.Users {
			if user.Name == cluster.AuthInfo {
				cluster.ClientCertificateData = string(user.User.ClientCertificateData)
				cluster.ClientKeyData = string(user.User.ClientKeyData)
				cluster.Token = user.User.Token
				break
			}
		}
		clusters = append(clusters, cluster)
	}

	return clusters
}

// ParseByteToKubeConfig 解析kubeconfig文件的[]byte为KubeConfig结构体
func ParseByteToKubeConfig(kubeconfig []byte) (kc *KubeConfig, err error) {
	kc = new(KubeConfig)
	err = json.Unmarshal(kubeconfig, kc)
	if err != nil {
		return nil, err
	}
	return kc, nil
}

// ParseClusterToKubeConfig 转换Models.Cluster结构体为KubeConfig结构体
func ParseClusterToKubeConfig(cluster Models.Cluster) (kc KubeConfig) {
	var kubeCluster KubeCluster
	var kubeContext KubeContext
	var kubeUser KubeUser

	kc.ApiVersion = cluster.APIVersion
	kc.Kind = cluster.Kind
	kc.CurrentContext = cluster.ContextName

	kubeCluster.Name = cluster.Cluster
	kubeCluster.Cluster.Server = cluster.Server
	kubeCluster.Cluster.InsecureSkipTLSVerify = cluster.InsecureSkipTLSVerify
	kubeCluster.Cluster.CertificateAuthorityData = []byte(cluster.CertificateAuthorityData)

	kubeContext.Name = cluster.ContextName
	kubeContext.Context.Cluster = cluster.Cluster
	kubeContext.Context.AuthInfo = cluster.AuthInfo

	kubeUser.Name = cluster.AuthInfo
	kubeUser.User.ClientKeyData = []byte(cluster.ClientKeyData)
	kubeUser.User.ClientCertificateData = []byte(cluster.ClientCertificateData)
	kubeUser.User.Token = cluster.Token

	kc.Clusters = append(kc.Clusters, kubeCluster)
	kc.Contexts = append(kc.Contexts, kubeContext)
	kc.Users = append(kc.Users, kubeUser)
	
	return kc
}

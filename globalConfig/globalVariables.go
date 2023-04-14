package globalConfig

import (
	"backend/tools"
	"k8s.io/client-go/rest"
)

// 初始化全局变量
var MyClient = new(tools.ClientManager)
var MyCfg = new(rest.Config)

//
//func init() {
//	var err error
//	MyCfg, err = clientcmd.BuildConfigFromFlags("", "./test")
//	if err != nil {
//		log.Panicf("加载kubeconfig文件失败: %s", err.Error())
//	}
//	MyClient.Client, err = kubernetes.NewForConfig(MyCfg)
//	if err != nil {
//		log.Panicf("生成kubernetes.Interface失败: %s", err.Error())
//	}
//	MyClient.DynamicClient, err = dynamic.NewForConfig(MyCfg)
//	if err != nil {
//		log.Panicf("生成dynamic.Interface失败: %s", err.Error())
//	}
//}

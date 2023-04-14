package resourcesCreat

import (
	"backend/globalConfig"
	"backend/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"net/http"
	"strings"
)

// AppDeploymentFromFileSpec is a specification for deployment from file
type AppDeploymentFromFileSpec struct {
	// Name of the file
	Name string `json:"name"`

	// 这个namespace是前端页面上下拉框选择的当前namespace
	Namespace string `json:"namespace"`

	// File content
	Content string `json:"content"`

	// Whether validate content before creation or not
	Validate bool `json:"validate"`
}

func DeployResourcesFromFile(c *gin.Context) (bool, error) {
	spec := new(AppDeploymentFromFileSpec)
	err := c.BindJSON(spec)
	fmt.Println(err)
	if err != nil {
		return false, err
	}
	reader := strings.NewReader(spec.Content)

	d := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	// 为什么用for循环？？？？？？？？？？
	for {
		data := unstructured.Unstructured{}
		if err := d.Decode(&data); err != nil {
			if err == io.EOF {
				return true, nil
			}
			return false, err
		}
		// 获取yaml中version，kind的字符串
		version := data.GetAPIVersion()
		kind := data.GetKind()

		// 将groupversion字符串解析为GroupVersion结构体
		gv, err := schema.ParseGroupVersion(version)
		fmt.Println("gv:", gv)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		//获取k8s资源管理客户端实例
		discoveryClient := globalConfig.MyClient.Client.Discovery()

		// 反应一个groupversion支持的所有资源类型，以及每个资源类型的名称，verbs（动作）,缩写，是否是namespace部署等
		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return false, err
		}

		//aa,_ := json.Marshal(*apiResourceList)
		//fmt.Println("apiResourceList", string(aa))
		//输出： {"kind":"APIResourceList","groupVersion":"v1","resources":[{"name":"bindings","singularName":"","namespaced":true,"kind":"Binding","verbs":["create"]},{"name":"componentstatuses","singularName":"","namespaced":false,"kind":"ComponentStatus","verbs":["get","list"],"shortNames":["cs"]},{"name":"configmaps","singularName":"","namespaced":true,"kind":"ConfigMap","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["cm"],"storageVersionHash":"qFsyl6wFWjQ="},{"name":"endpoints","singularName":"","namespaced":true,"kind":"Endpoints","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["ep"],"storageVersionHash":"fWeeMqaN/OA="},{"name":"events","singularName":"","namespaced":true,"kind":"Event","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["ev"],"storageVersionHash":"r2yiGXH7wu8="},{"name":"limitranges","singularName":"","namespaced":true,"kind":"LimitRange","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["limits"],"storageVersionHash":"EBKMFVe6cwo="},{"name":"namespaces","singularName":"","namespaced":false,"kind":"Namespace","verbs":["create","delete","get","list","patch","update","watch"],"shortNames":["ns"],"storageVersionHash":"Q3oi5N2YM8M="},{"name":"namespaces/finalize","singularName":"","namespaced":false,"kind":"Namespace","verbs":["update"]},{"name":"namespaces/status","singularName":"","namespaced":false,"kind":"Namespace","verbs":["get","patch","update"]},{"name":"nodes","singularName":"","namespaced":false,"kind":"Node","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["no"],"storageVersionHash":"XwShjMxG9Fs="},{"name":"nodes/proxy","singularName":"","namespaced":false,"kind":"NodeProxyOptions","verbs":["create","delete","get","patch","update"]},{"name":"nodes/status","singularName":"","namespaced":false,"kind":"Node","verbs":["get","patch","update"]},{"name":"persistentvolumeclaims","singularName":"","namespaced":true,"kind":"PersistentVolumeClaim","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["pvc"],"storageVersionHash":"QWTyNDq0dC4="},{"name":"persistentvolumeclaims/status","singularName":"","namespaced":true,"kind":"PersistentVolumeClaim","verbs":["get","patch","update"]},{"name":"persistentvolumes","singularName":"","namespaced":false,"kind":"PersistentVolume","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["pv"],"storageVersionHash":"HN/zwEC+JgM="},{"name":"persistentvolumes/status","singularName":"","namespaced":false,"kind":"PersistentVolume","verbs":["get","patch","update"]},{"name":"pods","singularName":"","namespaced":true,"kind":"Pod","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["po"],"categories":["all"],"storageVersionHash":"xPOwRZ+Yhw8="},{"name":"pods/attach","singularName":"","namespaced":true,"kind":"PodAttachOptions","verbs":["create","get"]},{"name":"pods/binding","singularName":"","namespaced":true,"kind":"Binding","verbs":["create"]},{"name":"pods/eviction","singularName":"","namespaced":true,"group":"policy","version":"v1beta1","kind":"Eviction","verbs":["create"]},{"name":"pods/exec","singularName":"","namespaced":true,"kind":"PodExecOptions","verbs":["create","get"]},{"name":"pods/log","singularName":"","namespaced":true,"kind":"Pod","verbs":["get"]},{"name":"pods/portforward","singularName":"","namespaced":true,"kind":"PodPortForwardOptions","verbs":["create","get"]},{"name":"pods/proxy","singularName":"","namespaced":true,"kind":"PodProxyOptions","verbs":["create","delete","get","patch","update"]},{"name":"pods/status","singularName":"","namespaced":true,"kind":"Pod","verbs":["get","patch","update"]},{"name":"podtemplates","singularName":"","namespaced":true,"kind":"PodTemplate","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"storageVersionHash":"LIXB2x4IFpk="},{"name":"replicationcontrollers","singularName":"","namespaced":true,"kind":"ReplicationController","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["rc"],"categories":["all"],"storageVersionHash":"Jond2If31h0="},{"name":"replicationcontrollers/scale","singularName":"","namespaced":true,"group":"autoscaling","version":"v1","kind":"Scale","verbs":["get","patch","update"]},{"name":"replicationcontrollers/status","singularName":"","namespaced":true,"kind":"ReplicationController","verbs":["get","patch","update"]},{"name":"resourcequotas","singularName":"","namespaced":true,"kind":"ResourceQuota","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["quota"],"storageVersionHash":"8uhSgffRX6w="},{"name":"resourcequotas/status","singularName":"","namespaced":true,"kind":"ResourceQuota","verbs":["get","patch","update"]},{"name":"secrets","singularName":"","namespaced":true,"kind":"Secret","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"storageVersionHash":"S6u1pOWzb84="},{"name":"serviceaccounts","singularName":"","namespaced":true,"kind":"ServiceAccount","verbs":["create","delete","deletecollection","get","list","patch","update","watch"],"shortNames":["sa"],"storageVersionHash":"pbx9ZvyFpBE="},{"name":"services","singularName":"","namespaced":true,"kind":"Service","verbs":["create","delete","get","list","patch","update","watch"],"shortNames":["svc"],"categories":["all"],"storageVersionHash":"0/CO1lhkEBI="},{"name":"services/proxy","singularName":"","namespaced":true,"kind":"ServiceProxyOptions","verbs":["create","delete","get","patch","update"]},{"name":"services/status","singularName":"","namespaced":true,"kind":"Service","verbs":["get","patch","update"]}]}

		// 获取groupversion中支持的所有资源信息
		apiResources := apiResourceList.APIResources

		// 定义一个APIResource变量
		var resource v1.APIResource
		for _, apiResource := range apiResources {
			// 排除name中包含"/"的资源
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = apiResource
				break
			}
		}

		if tools.IsStructEmpty(resource, v1.APIResource{}) {
			return false, fmt.Errorf("Unknown resource kind: %s", kind)
		}

		// ！！！重要Resource: resource.Name,否则会报错
		groupVersionResource := schema.GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource.Name}

		// 从yaml中获取资源类型，通过动态client获取资源类型对应的dynamicClient
		dynamicClientPool := globalConfig.MyClient.DynamicClient
		dynamicClient := dynamicClientPool.Resource(groupVersionResource)
		fmt.Println("groupversionresource", groupVersionResource)
		fmt.Println("errrrrrrrrrrrrrrrrrrrr:", err)
		// 判断前端当前namespace是否选择的是allnamespace,前端如果没有选择allnamespace，则判断，如果新建资源的yaml文件中写了namespace，就以yaml文件里为准，如果没有写namespace，就以当前namespace为准
		if strings.Compare(spec.Namespace, "all") == 0 || data.GetNamespace() != "" {
			_, err = dynamicClient.Namespace(data.GetNamespace()).Create(&data, v1.CreateOptions{})
		} else {
			_, err = dynamicClient.Namespace(spec.Namespace).Create(&data, v1.CreateOptions{})
		}

		if err != nil {
			return false, err
		}
	}
}

func ResourcesCreate(c *gin.Context) {
	_, err := DeployResourcesFromFile(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	}
}

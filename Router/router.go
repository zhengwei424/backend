package Router

import (
	"backend/Controllers/clusterManagement"
	myClusterRole "backend/Controllers/k8s/RBAC/clusterRole"
	myClusterRoleBinding "backend/Controllers/k8s/RBAC/clusterRoleBinding"
	myRole "backend/Controllers/k8s/RBAC/role"
	myRoleBinding "backend/Controllers/k8s/RBAC/roleBinding"
	myServiceAccount "backend/Controllers/k8s/RBAC/serviceAccount"
	myEndpoint "backend/Controllers/k8s/SVC/endpoint"
	myIngress "backend/Controllers/k8s/SVC/ingress"
	myNetworkPolicy "backend/Controllers/k8s/SVC/networkPolicy"
	myService "backend/Controllers/k8s/SVC/service"
	myConfigMap "backend/Controllers/k8s/config/configMap"
	mySecret "backend/Controllers/k8s/config/secret"
	myNamespace "backend/Controllers/k8s/namespace"
	myNode "backend/Controllers/k8s/node"
	myEvent "backend/Controllers/k8s/other/event"
	myTemplate "backend/Controllers/k8s/other/resourceTemplate"
	myResourcesCreate "backend/Controllers/k8s/other/resourcesCreat"
	myPersistentVolume "backend/Controllers/k8s/storage/persistentVolume"
	myPersistentVolumeClaim "backend/Controllers/k8s/storage/persistentVolumeClaim"
	myStorageClass "backend/Controllers/k8s/storage/storageClass"
	myCronJob "backend/Controllers/k8s/workload/cronJob"
	myDaemonSet "backend/Controllers/k8s/workload/daemonSet"
	myDeployment "backend/Controllers/k8s/workload/deployment"
	myJob "backend/Controllers/k8s/workload/job"
	myPod "backend/Controllers/k8s/workload/pod"
	myReplicaSet "backend/Controllers/k8s/workload/replicaSet"
	myReplicationController "backend/Controllers/k8s/workload/replicationController"
	myStatefulSet "backend/Controllers/k8s/workload/statefulSet"
	"backend/Controllers/login"
	"backend/Controllers/register"
	"backend/Middlewares"
	"backend/deprecated"
	"backend/tools"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func InitRouter() {
	var f1 = new(os.File)
	var f2 = new(os.File)
	var err error

	// 运行模式：debug、release、test
	gin.SetMode(gin.DebugMode)

	//gin.DisableConsoleColor() // 禁用日志颜色
	//gin.ForceConsoleColor() // 强制开启日志颜色

	// 记录标准输出（io.MultiWriter支持传入多个Writer，可以只输出到文件，也可以同时输出到控制台等）
	f1, _ = os.Create("logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f1, os.Stdout)

	//记录错误输出
	f2, _ = os.Create("logs/gin_err.log")
	gin.DefaultErrorWriter = io.MultiWriter(f2, os.Stderr)

	// 创建gin实例
	r := gin.Default()
	//r := gin.New()
	r.GET("/workload/exec", deprecated.WsHandler)
	// websocket  <*path>很重要，能匹配/api/sockjs/和/api/sockjs/info等
	// 不经过自定义的跨域中间件，sockJS自己默认有。。。。
	r.Any("/api/sockjs/*path", gin.WrapH(myPod.CreateAttachHandler("/api/sockjs")))

	// 添加中间件
	r.Use(Middlewares.Cors(), gin.LoggerWithFormatter(tools.MyLogFormatter))

	// 添加session

	// 登录
	r.POST("/login", login.Login)

	// 注册
	r.POST("/register", register.Register)

	// 获取登录验证码
	r.GET("/captcha", tools.GenerateCaptcha)

	// 路由组
	r.RouterGroup.Use(Middlewares.AuthMiddleware())

	// 上传kubeconfig文件
	cm := r.Group("clusterManagement")
	{
		cm.POST("/upload", clusterManagement.UpLoadKubeConfigFile)
		cm.GET("/clusters", clusterManagement.GetClusterKubeConfigPartInfo)
		cm.POST("/connectCluster", clusterManagement.ConnectCluster)
	}

	// 设置cluster资源路由组
	cluster := r.Group("cluster")
	{
		// namespace Info
		cluster.GET("/namespaces", myNamespace.GetNamespacesInfo)
		// namespace CRUD
		cluster.POST("/namespace/delete", myNamespace.DeleteNamespace)
		cluster.POST("/namespace/update", myNamespace.UpdateNamespace)
		cluster.GET("/namespace/get", myNamespace.GetNamespace)
		// node Info
		cluster.GET("/nodes", myNode.GetNodesInfo)
		// node CRUD
		cluster.POST("/node/delete", myNode.DeleteNode)
		cluster.POST("/node/update", myNode.UpdateNode)
		cluster.GET("/node/get", myNode.GetNode)
	}

	// 设置config资源路由组
	config := r.Group("config")
	{
		// configMap Info
		config.GET("/configMaps", myConfigMap.GetConfigMapsInfo)
		// configMap CRUD
		config.POST("/configMap/delete", myConfigMap.DeleteConfigMap)
		config.POST("/configMap/update", myConfigMap.UpdateConfigMap)
		config.GET("/configMap/get", myConfigMap.GetConfigMap)
		// secret Info
		config.GET("/secrets", mySecret.GetSecretsInfo)
		// secret CRUD
		config.POST("/secret/delete", mySecret.DeleteSecret)
		config.POST("/secret/update", mySecret.UpdateSecret)
		config.GET("/secret/get", mySecret.GetSecret)
	}

	// 设置other资源路由组
	other := r.Group("other")
	{
		// 新建资源统一入口
		other.POST("/resourcesCreate", myResourcesCreate.ResourcesCreate)
		// 新建资源模板文件
		other.GET("/template", myTemplate.GetResourceTemplate)
		// 事件信息
		other.GET("/events", myEvent.GetEventsInfo)
	}

	// 设置rbac资源路由组
	rbac := r.Group("rbac")
	{
		// clusterRole Info
		rbac.GET("/clusterRoles", myClusterRole.GetClusterRolesInfo)
		// clusterRole CRUD
		rbac.POST("/clusterRole/delete", myClusterRole.DeleteClusterRole)
		rbac.POST("/clusterRole/update", myClusterRole.UpdateClusterRole)
		rbac.GET("/clusterRole/get", myClusterRole.GetClusterRole)
		// clusterRoleBinding Info
		rbac.GET("/clusterRoleBindings", myClusterRoleBinding.GetClusterRoleBindingsInfo)
		// clusterRoleBinding CRUD
		rbac.POST("/clusterRoleBinding/delete", myClusterRoleBinding.DeleteClusterRoleBinding)
		rbac.POST("/clusterRoleBinding/update", myClusterRoleBinding.UpdateClusterRoleBinding)
		rbac.GET("/clusterRoleBinding/get", myClusterRoleBinding.GetClusterRoleBinding)
		// role Info
		rbac.GET("/roles", myRole.GetRolesInfo)
		// role CRUD
		rbac.POST("/role/delete", myRole.DeleteRole)
		rbac.POST("/role/update", myRole.UpdateRole)
		rbac.GET("/role/get", myRole.GetRole)
		// roleBinding Info
		rbac.GET("/roleBindings", myRoleBinding.GetRoleBindingsInfo)
		// roleBinding CRUD
		rbac.POST("/roleBinding/delete", myRoleBinding.DeleteRoleBinding)
		rbac.POST("/roleBinding/update", myRoleBinding.UpdateRoleBinding)
		rbac.GET("/roleBinding/get", myRoleBinding.GetRoleBinding)
		// serviceAccount Info
		rbac.GET("/serviceAccounts", myServiceAccount.GetServiceAccountsInfo)
		// serviceAccount CRUD
		rbac.POST("/serviceAccount/delete", myServiceAccount.DeleteServiceAccount)
		rbac.POST("/serviceAccount/update", myServiceAccount.UpdateServiceAccount)
		rbac.GET("/serviceAccount/get", myServiceAccount.GetServiceAccount)
	}
	// 设置storage资源路由组
	storage := r.Group("storage")
	{
		// persistentVolume Info
		storage.GET("/persistentVolumes", myPersistentVolume.GetPersistentVolumesInfo)
		// persistentVolume CRUD
		storage.POST("/persistentVolume/delete", myPersistentVolume.DeletePersistentVolume)
		storage.POST("/persistentVolume/update", myPersistentVolume.UpdatePersistentVolume)
		storage.GET("/persistentVolume/get", myPersistentVolume.GetPersistentVolume)
		// persistentVolumeClaim Info
		storage.GET("/persistentVolumeClaims", myPersistentVolumeClaim.GetPersistentVolumeClaimsInfo)
		// persistentVolumeClaim CRUD
		storage.POST("/persistentVolumeClaim/delete", myPersistentVolumeClaim.DeletePersistentVolumeClaim)
		storage.POST("/persistentVolumeClaim/update", myPersistentVolumeClaim.UpdatePersistentVolumeClaim)
		storage.GET("/persistentVolumeClaim/get", myPersistentVolumeClaim.GetPersistentVolumeClaim)
		// storageClass Info
		storage.GET("/storageClasses", myStorageClass.GetStorageClassesInfo)
		// storageClass CRUD
		storage.POST("/storageClass/delete", myStorageClass.DeletePersistentVolume)
		storage.POST("/storageClass/update", myStorageClass.UpdatePersistentVolume)
		storage.GET("/storageClass/get", myStorageClass.GetPersistentVolume)
	}
	// 设置svc资源路由组
	svc := r.Group("svc")
	{
		// ingress Info
		svc.GET("/ingresses", myIngress.GetIngressesInfo)
		// ingress CRUD
		svc.POST("/ingress/delete", myIngress.DeleteIngress)
		svc.POST("/ingress/update", myIngress.UpdateIngress)
		svc.GET("/ingress/get", myIngress.GetIngress)
		// networkPolicy Info
		svc.GET("/networkPolicies", myNetworkPolicy.GetNetworkPoliciesInfo)
		// networkPolicy CRUD
		svc.POST("/networkPolicy/delete", myNetworkPolicy.DeleteNetworkPolicy)
		svc.POST("/networkPolicy/update", myNetworkPolicy.UpdateNetworkPolicy)
		svc.GET("/networkPolicy/get", myNetworkPolicy.GetNetworkPolicy)
		// service Info
		svc.GET("/services", myService.GetServicesInfo)
		// service CRUD
		svc.POST("/service/delete", myService.DeleteService)
		svc.POST("/service/update", myService.UpdateService)
		svc.GET("/service/get", myService.GetService)
		// endpoint Info
		svc.GET("/endpoints", myEndpoint.GetEndpointsInfo)
		// endpoint CRUD
		svc.POST("/endpoint/delete", myEndpoint.Deleteendpoint)
		svc.POST("/endpoint/update", myEndpoint.Updateendpoint)
		svc.GET("/endpoint/get", myEndpoint.GetEndpoint)
	}
	// 设置workload资源路由组
	workload := r.Group("workload")
	{
		// pod Info
		workload.GET("/pods", myPod.GetPodsInfo)
		// pod Exec
		workload.GET("/:namespace/:pod/:container/exec", myPod.HandleExecShell)
		// pod Log
		workload.GET("/log/:namespace/:pod/:container", myPod.HandleLogs)
		// pod LogFile
		workload.GET("/logfile/:namespace/:pod/:container", myPod.HandleLogFile)
		// pod CRUD
		workload.POST("/pod/delete", myPod.DeletePod)
		workload.POST("/pod/update", myPod.UpdatePod)
		workload.GET("/pod/get", myPod.GetPod)
		// deployment Info
		workload.GET("/deployments", myDeployment.GetDeploymentsInfo)
		// deployment CRUD
		workload.POST("/deployment/delete", myDeployment.DeleteDeployment)
		workload.POST("/deployment/update", myDeployment.UpdateDeployment)
		workload.GET("/deployment/get", myDeployment.GetDeployment)
		// daemonSet Info
		workload.GET("/daemonSets", myDaemonSet.GetDaemonSetsInfo)
		// daemonSet CRUD
		workload.POST("/daemonSet/delete", myDaemonSet.DeleteDaemonSet)
		workload.POST("/daemonSet/update", myDaemonSet.UpdateDaemonSet)
		workload.GET("/daemonSet/get", myDaemonSet.GetDaemonSet)
		// statefulSet Info
		workload.GET("/statefulSets", myStatefulSet.GetStatefulSetsInfo)
		// statefulSet CRUD
		workload.POST("/statefulSet/delete", myStatefulSet.DeleteStatefulSet)
		workload.POST("/statefulSet/update", myStatefulSet.UpdateStatefulSet)
		workload.GET("/statefulSet/get", myStatefulSet.GetStatefulSet)
		// replicaSet Info
		workload.GET("/replicaSets", myReplicaSet.GetReplicaSetsInfo)
		// replicaSet CRUD
		workload.POST("/replicaSet/delete", myReplicaSet.DeleteReplicaSet)
		workload.POST("/replicaSet/update", myReplicaSet.UpdateReplicaSet)
		workload.GET("/replicaSet/get", myReplicaSet.GetReplicaSet)
		// replicationController Info
		workload.GET("/replicationControllers", myReplicationController.GetReplicationControllersInfo)
		// replicationController CRUD
		workload.POST("/replicationController/delete", myReplicationController.DeleteReplicationController)
		workload.POST("/replicationController/update", myReplicationController.UpdateReplicationController)
		workload.GET("/replicationController/get", myReplicationController.GetReplicationController)
		// cronJob Info
		workload.GET("/cronJobs", myCronJob.GetCronJobsInfo)
		// cronJob CRUD
		workload.POST("/cronJob/delete", myCronJob.DeleteJob)
		workload.POST("/cronJob/update", myCronJob.UpdateCronJob)
		workload.GET("/cronJob/get", myCronJob.GetCronJob)
		// job Info
		workload.GET("/jobs", myJob.GetJobsInfo)
		// job CRUD
		workload.POST("/job/delete", myJob.DeleteJob)
		workload.POST("/job/update", myJob.UpdateJob)
		workload.GET("/job/get", myJob.GetJob)
	}

	// 启动
	err = r.RunTLS(":9090", "ssl/cert.pem", "ssl/key.pem")
	if err != nil {
		panic(err)
	}
}

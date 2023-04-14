package Router

import (
	"backend/Controllers/clusterManagement"
	myClusterRole "backend/Controllers/k8s/RBAC/clusterRole"
	myClusterRoleBinding "backend/Controllers/k8s/RBAC/clusterRoleBinding"
	myRole "backend/Controllers/k8s/RBAC/role"
	myRoleBinding "backend/Controllers/k8s/RBAC/roleBinding"
	myServiceAccount "backend/Controllers/k8s/RBAC/serviceAccount"
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
		cluster.GET("/namespaces", myNamespace.GetNamespacesInfo)
		cluster.GET("/nodes", myNode.GetNodesInfo)
	}

	// 设置config资源路由组
	config := r.Group("globalConfig")
	{
		config.GET("/configMaps", myConfigMap.GetConfigMapsInfo)
		config.GET("/secrets", mySecret.GetSecretsInfo)
	}

	// 设置other资源路由组
	other := r.Group("other")
	{
		// 资源新建
		other.POST("/resourcesCreate", myResourcesCreate.ResourcesCreate)
		other.GET("/events", myEvent.GetEventsInfo)
		other.GET("/template", myTemplate.GetResourceTemplate)
	}

	// 设置rbac资源路由组
	rbac := r.Group("rbac")
	{
		rbac.GET("/clusterRoles", myClusterRole.GetClusterRolesInfo)
		rbac.GET("/clusterRoleBindings", myClusterRoleBinding.GetClusterRoleBindingsInfo)
		rbac.GET("/roles", myRole.GetRolesInfo)
		rbac.GET("/roleBindings", myRoleBinding.GetRoleBindingsInfo)
		rbac.GET("/serviceAccounts", myServiceAccount.GetServiceAccountsInfo)
	}
	// 设置storage资源路由组
	storage := r.Group("storage")
	{
		storage.GET("/persistentVolumes", myPersistentVolume.GetPersistentVolumesInfo)
		storage.GET("/persistentVolumeClaims", myPersistentVolumeClaim.GetPersistentVolumeClaimsInfo)
		storage.GET("/storageClasses", myStorageClass.GetStorageClassesInfo)
	}
	// 设置svc资源路由组
	svc := r.Group("svc")
	{
		svc.GET("/ingresses", myIngress.GetIngressesInfo)
		svc.GET("/networkPolicies", myNetworkPolicy.GetNetworkPoliciesInfo)
		svc.GET("/services", myService.GetServicesInfo)
	}
	// 设置workload资源路由组
	workload := r.Group("workload")
	{
		workload.POST("/pod/create", myPod.CreatePod)
		workload.POST("/pod/delete", myPod.DeletePod)
		workload.POST("/pod/update", myPod.UpdatePod)
		workload.GET("/pod/get", myPod.GetPod)
		workload.GET("/pods", myPod.GetPodsInfo)
		workload.GET("/deployments", myDeployment.GetDeploymentsInfo)
		workload.GET("/daemonSets", myDaemonSet.GetDaemonSetsInfo)
		workload.GET("/statefulSets", myStatefulSet.GetStatefulSetsInfo)
		workload.GET("/replicaSets", myReplicaSet.GetReplicaSetsInfo)
		workload.GET("/replicationControllers", myReplicationController.GetReplicationControllersInfo)
		workload.GET("/cronJobs", myCronJob.GetCronJobsInfo)
		workload.GET("/jobs", myJob.GetJobsInfo)
		//workload.GET("/exec", myPod.WsHandler)
		workload.GET("/:namespace/:pod/:container/exec", myPod.HandleExecShell)
		//workload.GET("/log", deprecated.GetContainerLog)
		workload.GET("/log/:namespace/:pod/:container", myPod.HandleLogs)
		workload.GET("/logfile/:namespace/:pod/:container", myPod.HandleLogFile)
	}

	// 启动
	err = r.RunTLS(":9090", "ssl/cert.pem", "ssl/key.pem")
	if err != nil {
		panic(err)
	}
}

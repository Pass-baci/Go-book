package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)

type KitDiscoverClient struct {
	Host         string
	Port         int
	client       consul.Client
	config       *api.Config
	mutex        sync.Mutex
	instancesMap sync.Map // 使用map来储存服务实例列表，减少与consul交互
}

// 初始化KitDiscoveryClient
func NewKitDiscoveryClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	consulClient := consul.NewClient(apiClient)
	return &KitDiscoverClient{
		Host:   consulHost,
		Port:   consulPort,
		config: consulConfig,
		client: consulClient,
	}, nil
}

// 服务注册接口
func (consulClient *KitDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,   //服务实例ID
		Name:    serviceName,  // 服务名称
		Address: instanceHost, // 服务地址
		Port:    instancePort, // 服务端口
		Meta:    meta,         // 服务元数据
		Check: &api.AgentServiceCheck{ // 服务健康检查地址
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}
	// 向consul中发起服务注册
	if err := consulClient.client.Register(serviceRegistration); err != nil {
		logger.Println("Register Service Error")
		return false
	}
	logger.Println("DeRegister Service Success")
	return true
}

// 服务注销接口
func (consulClient *KitDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId, //服务实例ID
	}
	if err := consulClient.client.Deregister(serviceRegistration); err != nil {
		logger.Println("DeRegister Service Error")
		return false
	}
	logger.Println("DeRegister Service Success")
	return true
}

// 发现服务实例接口
func (consulClient *KitDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	// 该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}
	// 申请锁
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 再次检查是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	} else {
		// 开启一个goroutine进行consul服务实例监控
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			// 向consul注册Service类型的Watch监控机制
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return
				}
				// 记录服务实例列表
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []interface{}{})
				}
				var healthServices []interface{}
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, service.Service)
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		logger.Println("Discover Services Error")
		return nil
	}
	instances := make([]interface{}, len(entries))
	for i := 0; i < len(entries); i++ {
		instances[i] = entries[i].Service
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

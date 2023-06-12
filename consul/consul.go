package consul

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/JohnSalazar/microservices-go-common/config"

	consul "github.com/hashicorp/consul/api"
)

func NewConsulClient(
	config *config.Config,
) (*consul.Client, string, error) {

	consulConfig := consul.DefaultConfig()
	consulConfig.Address = config.Consul.Host

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, "", err
	}

	serviceID, err := register(config, consulClient)
	if err != nil {
		return nil, "", err
	}

	return consulClient, serviceID, nil
}

func register(config *config.Config, client *consul.Client) (string, error) {

	var check_port int
	address := hostname()

	k8s, _ := strconv.ParseBool(os.Getenv("kubernetes"))
	if k8s {
		address = fmt.Sprintf("%s-%s", config.AppName, config.KubernetesServiceNameSuffix)
	}

	port, err := strconv.Atoi(strings.Split(config.ListenPort, ":")[1])
	if port == 0 || err != nil {
		return "", err
	}

	check_port = port

	if len(strings.TrimSpace(config.GrpcServer.Port)) > 0 {
		port, err = strconv.Atoi(strings.Split(config.GrpcServer.Port, ":")[1])
		if err != nil {
			return "", err
		}
	}

	serviceID := fmt.Sprintf("%s-%s:%v", config.AppName, address, port)

	httpCheck := fmt.Sprintf("https://%s:%v/healthy", address, check_port)
	fmt.Println(httpCheck)

	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    config.AppName,
		Port:    port,
		Address: address,
		Check: &consul.AgentServiceCheck{
			CheckID:                        serviceID,
			Name:                           fmt.Sprintf("Service %s check", config.AppName),
			HTTP:                           httpCheck,
			TLSSkipVerify:                  true,
			Interval:                       "10s",
			Timeout:                        "30s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Printf("failed consul to register service: %s:%v ", address, port)
		return "", err
	}

	log.Printf("successfully consul register service: %s:%v", address, port)

	return serviceID, nil
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	return hostname
}

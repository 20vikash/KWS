package config

const (
	CORE_IMAGE_NAME      = "core_ubuntu:latest"
	CORE_NETWORK_NAME    = "kws_instance"
	CORE_NETWORK_SUBNET  = "172.35.0.0/24"
	CORE_NETWORK_GATEWAY = "172.35.0.1"

	SERVICES_NETWORK_NAME = "kws_services"

	MAIN_INSTANCE_QUEUE = "instance_queue"
	RETRY_QUEUE         = "retry_queue"

	DEPLOY = "deploy"
	STOP   = "stop"
	KILL   = "kill"

	INTERFACE_NAME    = "wg0"
	INTERFACE_ADDRESS = "10.0.0.1/24"
	CIDR              = 24

	STACK_KEY = "ip_stack"

	MAX_WG_DEVICES_PER_USER = 3

	MAX_SERVICE_DB_USERS = 5
	MAX_SERVICE_DB_DB    = 10
)

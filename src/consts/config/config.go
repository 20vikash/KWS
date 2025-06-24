package config

const (
	CORE_IMAGE_NAME   = "core_ubuntu:latest"
	CORE_NETWORK_NAME = "kws_ins_s"

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
)

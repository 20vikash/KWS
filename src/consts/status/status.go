package status

const (
	USER_NAME_INVALID = "user_name_invalid"
	USER_NOT_VERIFIED = "user_not_verified"
	WRONG_CREDENTIALS = "wrong_credentials"

	CONTAINER_ALREADY_RUNNING     = "container_already_running"
	CONTAINER_ALREADY_EXISTS      = "container_already_exists"
	CONTAINER_NOT_FOUND_TO_STOP   = "container_not_found_to_stop"
	CONTAINER_NOT_FOUND_TO_DELETE = "container_not_found_to_delete"
	VOLUME_NOT_FOUND              = "volume_not_found"

	// DB errors
	CONTAINER_START_FAILED  = "container_start_failed"
	CONTAINER_STOP_FAILED   = "container_stop_failed"
	CONTAINER_DELETE_FAILED = "container_delete_failed"

	INTERFACE_ALREADY_EXISTS = "interface_already_exists"

	PEER_DOES_NOT_EXIST = "peer_does_not_exist"

	INVALID_CIDR    = "invalid_cidr"
	HOST_EXHAUSTION = "host_exhaustion"

	EMPTY_IP_STACK = "empty_ip_stack"
)

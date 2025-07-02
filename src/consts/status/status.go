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

	WG_DEVICE_LIMIT = "wg_device_limit"

	PG_MAX_USER_LIMIT      = "pg_max_user_limit"
	PG_MAX_DB_LIMIT        = "pg_max_db_limit"
	PG_USER_ALREDAY_EXISTS = "pg_user_already_exists"
	PG_DB_ALREDAY_EXISTS   = "pg_db_already_exists"
	PG_USER_NOT_FOUND      = "pg_user_not_found"
	PG_DB_NOT_FOUND        = "pg_db_not_found"

	DOMAIN_ALREADY_EXISTS = "domain_already_exists"
	DOMAIN_LIMIT_EXCEEDED = "domain_limit_exceeded"
)

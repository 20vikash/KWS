package config

import "kws/kws/internal/kwsconfig"

// --- Truly constant values (protocol/application level, never change per deployment) ---

const (
	CORE_IMAGE_NAME       = "core_ubuntu:latest"
	CORE_NETWORK_NAME     = "kws_instance"
	SERVICES_NETWORK_NAME = "kws_kws_services"

	MAIN_INSTANCE_QUEUE  = "instance_queue"
	INSTANCE_RETRY_QUEUE = "instance_retry_queue"
	MAIN_TUNNEL_QUEUE    = "tunnel_queue"
	TUNNEL_RETRY_QUEUE   = "tunnel_retry_queue"
	USER_DOMAIN_QUEUE    = "user_domain_queue"
	DOMAIN_RETRY_QUEUE   = "domain_retry_queue"

	DEPLOY = "deploy"
	STOP   = "stop"
	KILL   = "kill"

	ADD_USER_DOMAIN    = "add_user_domain"
	REMOVE_USER_DOMAIN = "remove_user_domain"

	STACK_KEY = "ip_stack"
	LXC_IP    = "lxc_ip"

	NGINX_CONTAINER = "nginx_proxy"

	INSTANCE_START = "start"
	INSTANCE_STOP  = "stop"

	INSTANCE_TEMPLATE = "instance_template"
	DOMAIN_TEMPLATE   = "domain_template"

	NO_DOMAIN_FOR_TUNNEL = "no_domain_for_tunnel"
	X_RETRY_COUNTER      = "x-retry-counter"
)

// --- Deployment-specific values (read from kws_config.yaml) ---

func CORE_NETWORK_SUBNET() string  { return kwsconfig.Get().Network.CoreNetworkSubnet }
func CORE_NETWORK_GATEWAY() string { return kwsconfig.Get().Network.CoreNetworkGateway }

func INTERFACE_NAME() string    { return kwsconfig.Get().Wireguard.InterfaceName }
func INTERFACE_ADDRESS() string { return kwsconfig.Get().Wireguard.Address }
func CIDR() int                 { return kwsconfig.Get().Wireguard.CIDR }
func WG_LISTEN_PORT() int       { return kwsconfig.Get().Wireguard.ListenPort }
func WG_KEEPALIVE_SEC() int     { return kwsconfig.Get().Wireguard.KeepAliveSec }

func MAX_WG_DEVICES_PER_USER() int { return kwsconfig.Get().Limits.MaxWGDevicesPerUser }
func MAX_SERVICE_DB_USERS() int    { return kwsconfig.Get().Limits.MaxServiceDBUsers }
func MAX_SERVICE_DB_DB() int       { return kwsconfig.Get().Limits.MaxServiceDBDatabases }
func USER_DOMAIN_LIMIT() int       { return kwsconfig.Get().Limits.UserDomainLimit }

func LXC_UBUNTU_ALIAS() string { return kwsconfig.Get().Instance.UbuntuAlias }
func LXD_BRIDGE() string       { return kwsconfig.Get().Network.LXDBridgeName }
func STORAGE_POOL() string     { return kwsconfig.Get().Instance.StoragePool }
func DNS_IP() string           { return kwsconfig.Get().Network.DNSIP }

func DOMAIN() string                { return kwsconfig.Get().Domain }
func GATEWAY_PORT() int             { return kwsconfig.Get().Server.GatewayPort }
func CODE_SERVER_PORT() int         { return kwsconfig.Get().Server.CodeServerPort }
func TUNNEL_PROXY_PORT() int        { return kwsconfig.Get().Server.TunnelProxyPort }
func INSTANCE_MEMORY_LIMIT() string { return kwsconfig.Get().Instance.MemoryLimit }

func SSL_CERT_PATH() string { return kwsconfig.Get().SSL.CertPath }
func SSL_KEY_PATH() string  { return kwsconfig.Get().SSL.KeyPath }

func LXD_BRIDGE_SUBNET() string  { return kwsconfig.Get().Network.LXDBridgeSubnet }
func LXD_BRIDGE_GATEWAY() string { return kwsconfig.Get().Network.LXDBridgeGateway }
func LXD_SOCKET_PATH() string    { return kwsconfig.Get().Paths.LXDSocket }
func NGINX_CONF_DIR() string     { return kwsconfig.Get().Paths.NginxConfDir }
func LXC_IP_START() int          { return kwsconfig.Get().Network.LXCIPStart }

func PG_SERVICE_PORT() int { return kwsconfig.Get().Services.PostgresPort }
func ADMINER_PORT() int    { return kwsconfig.Get().Services.AdminerPort }

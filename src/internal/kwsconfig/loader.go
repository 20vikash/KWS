package kwsconfig

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// KWSConfig is the top-level configuration structure for a KWS deployment.
type KWSConfig struct {
	Domain   string `yaml:"domain"`
	PublicIP string `yaml:"public_ip"`

	Server  ServerConfig  `yaml:"server"`
	SSL     SSLConfig     `yaml:"ssl"`
	Network NetworkConfig `yaml:"network"`

	Wireguard WireguardConfig `yaml:"wireguard"`
	Services  ServicesConfig  `yaml:"services"`
	Instance  InstanceConfig  `yaml:"instance"`
	Limits    LimitsConfig    `yaml:"limits"`
	Paths     PathsConfig     `yaml:"paths"`
}

type ServerConfig struct {
	GatewayPort     int `yaml:"gateway_port"`
	CodeServerPort  int `yaml:"code_server_port"`
	TunnelProxyPort int `yaml:"tunnel_proxy_port"`
}

type SSLConfig struct {
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
}

type NetworkConfig struct {
	LXDBridgeName      string `yaml:"lxd_bridge_name"`
	LXDBridgeSubnet    string `yaml:"lxd_bridge_subnet"`
	LXDBridgeGateway   string `yaml:"lxd_bridge_gateway"`
	ServicesSubnet     string `yaml:"services_subnet"`
	ServicesGateway    string `yaml:"services_gateway"`
	CoreNetworkSubnet  string `yaml:"core_network_subnet"`
	CoreNetworkGateway string `yaml:"core_network_gateway"`
	DNSIP              string `yaml:"dns_ip"`
}

type WireguardConfig struct {
	InterfaceName string `yaml:"interface_name"`
	Address       string `yaml:"address"`
	CIDR          int    `yaml:"cidr"`
	ListenPort    int    `yaml:"listen_port"`
	KeepAliveSec  int    `yaml:"keepalive_sec"`
}

type ServicesConfig struct {
	PostgresIP       string         `yaml:"postgres_ip"`
	PostgresHostname string         `yaml:"postgres_hostname"`
	PostgresPort     int            `yaml:"postgres_port"`
	AdminerIP        string         `yaml:"adminer_ip"`
	AdminerHostname  string         `yaml:"adminer_hostname"`
	AdminerPort      int            `yaml:"adminer_port"`
	BridgeAttach     []BridgeAttach `yaml:"bridge_attach"`
}

type BridgeAttach struct {
	Container string `yaml:"container"`
	Bridge    string `yaml:"bridge"`
	IPCIDR    string `yaml:"ip_cidr"`
}

type InstanceConfig struct {
	MemoryLimit  string `yaml:"memory_limit"`
	StoragePool  string `yaml:"storage_pool"`
	UbuntuAlias  string `yaml:"ubuntu_alias"`
}

type LimitsConfig struct {
	MaxWGDevicesPerUser int `yaml:"max_wg_devices_per_user"`
	MaxServiceDBUsers   int `yaml:"max_service_db_users"`
	MaxServiceDBDatabases int `yaml:"max_service_db_databases"`
	UserDomainLimit     int `yaml:"user_domain_limit"`
}

type PathsConfig struct {
	LXDSocket    string `yaml:"lxd_socket"`
	NginxConfDir string `yaml:"nginx_conf_dir"`
}

// Global config singleton
var (
	globalConfig *KWSConfig
	configOnce   sync.Once
)

// Init loads the config from the given path. Must be called once at startup.
func Init(path string) error {
	var initErr error
	configOnce.Do(func() {
		cfg, err := Load(path)
		if err != nil {
			initErr = err
			return
		}
		globalConfig = cfg
	})
	return initErr
}

// Get returns the global config singleton. Panics if Init was not called.
func Get() *KWSConfig {
	if globalConfig == nil {
		log.Fatal("kwsconfig.Init() must be called before kwsconfig.Get()")
	}
	return globalConfig
}

// Load reads and parses a kws_config.yaml file from the given path.
func Load(path string) (*KWSConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg KWSConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	applyDefaults(&cfg)

	return &cfg, nil
}

// applyDefaults fills in zero-value fields with sensible defaults.
func applyDefaults(cfg *KWSConfig) {
	// Server defaults
	if cfg.Server.GatewayPort == 0 {
		cfg.Server.GatewayPort = 8080
	}
	if cfg.Server.CodeServerPort == 0 {
		cfg.Server.CodeServerPort = 8099
	}
	if cfg.Server.TunnelProxyPort == 0 {
		cfg.Server.TunnelProxyPort = 8081
	}

	// Network defaults
	if cfg.Network.LXDBridgeName == "" {
		cfg.Network.LXDBridgeName = "lxdbr0"
	}
	if cfg.Network.LXDBridgeSubnet == "" {
		cfg.Network.LXDBridgeSubnet = "172.30.0.0/24"
	}
	if cfg.Network.LXDBridgeGateway == "" {
		cfg.Network.LXDBridgeGateway = "172.30.0.1"
	}
	if cfg.Network.ServicesSubnet == "" {
		cfg.Network.ServicesSubnet = "172.25.0.0/24"
	}
	if cfg.Network.ServicesGateway == "" {
		cfg.Network.ServicesGateway = "172.25.0.1"
	}
	if cfg.Network.CoreNetworkSubnet == "" {
		cfg.Network.CoreNetworkSubnet = "172.35.0.0/24"
	}
	if cfg.Network.CoreNetworkGateway == "" {
		cfg.Network.CoreNetworkGateway = "172.35.0.1"
	}
	if cfg.Network.DNSIP == "" {
		cfg.Network.DNSIP = "172.30.0.102"
	}

	// Wireguard defaults
	if cfg.Wireguard.InterfaceName == "" {
		cfg.Wireguard.InterfaceName = "wg0"
	}
	if cfg.Wireguard.Address == "" {
		cfg.Wireguard.Address = "10.0.0.1/24"
	}
	if cfg.Wireguard.CIDR == 0 {
		cfg.Wireguard.CIDR = 24
	}
	if cfg.Wireguard.ListenPort == 0 {
		cfg.Wireguard.ListenPort = 51820
	}
	if cfg.Wireguard.KeepAliveSec == 0 {
		cfg.Wireguard.KeepAliveSec = 25
	}

	// Services defaults
	if cfg.Services.PostgresIP == "" {
		cfg.Services.PostgresIP = "172.25.0.2"
	}
	if cfg.Services.PostgresHostname == "" {
		cfg.Services.PostgresHostname = "postgres.kws.services"
	}
	if cfg.Services.PostgresPort == 0 {
		cfg.Services.PostgresPort = 5432
	}
	if cfg.Services.AdminerIP == "" {
		cfg.Services.AdminerIP = "172.25.0.4"
	}
	if cfg.Services.AdminerHostname == "" {
		cfg.Services.AdminerHostname = "adminer.kws.services"
	}
	if cfg.Services.AdminerPort == 0 {
		cfg.Services.AdminerPort = 8080
	}

	// Instance defaults
	if cfg.Instance.MemoryLimit == "" {
		cfg.Instance.MemoryLimit = "1500MB"
	}
	if cfg.Instance.StoragePool == "" {
		cfg.Instance.StoragePool = "kws"
	}
	if cfg.Instance.UbuntuAlias == "" {
		cfg.Instance.UbuntuAlias = "ubuntu-22:04"
	}

	// Limits defaults
	if cfg.Limits.MaxWGDevicesPerUser == 0 {
		cfg.Limits.MaxWGDevicesPerUser = 3
	}
	if cfg.Limits.MaxServiceDBUsers == 0 {
		cfg.Limits.MaxServiceDBUsers = 5
	}
	if cfg.Limits.MaxServiceDBDatabases == 0 {
		cfg.Limits.MaxServiceDBDatabases = 10
	}
	if cfg.Limits.UserDomainLimit == 0 {
		cfg.Limits.UserDomainLimit = 3
	}

	// Paths defaults
	if cfg.Paths.LXDSocket == "" {
		cfg.Paths.LXDSocket = "/var/snap/lxd/common/lxd/unix.socket"
	}
	if cfg.Paths.NginxConfDir == "" {
		cfg.Paths.NginxConfDir = "/app/nginx_conf"
	}
}

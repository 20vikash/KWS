package main

// KWSConfig is the top-level config structure written to kws_config.yaml
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
	MemoryLimit string `yaml:"memory_limit"`
	StoragePool string `yaml:"storage_pool"`
	UbuntuAlias string `yaml:"ubuntu_alias"`
}

type LimitsConfig struct {
	MaxWGDevicesPerUser   int `yaml:"max_wg_devices_per_user"`
	MaxServiceDBUsers     int `yaml:"max_service_db_users"`
	MaxServiceDBDatabases int `yaml:"max_service_db_databases"`
	UserDomainLimit       int `yaml:"user_domain_limit"`
}

type PathsConfig struct {
	LXDSocket    string `yaml:"lxd_socket"`
	NginxConfDir string `yaml:"nginx_conf_dir"`
}

// EnvConfig holds all values that go into the .env file
type EnvConfig struct {
	// Main DB
	DBUsername string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Gmail
	GmailAppPassword string
	GmailAddress     string
	SMTPHost         string
	SMTPPort         string

	// Environment
	Env string

	// RabbitMQ
	MQUser       string
	MQPassword   string
	MQServerPort string
	MQUIPort     string
	MQHost       string

	// Wireguard
	WGPrivateKey string

	// PG Service
	PGServiceUsername string
	PGServicePassword string
	PGServiceHost     string
	PGServicePort     string
	PGServiceDB       string

	// Compose network overrides
	ServicesSubnet  string
	ServicesGateway string
	PGServiceIP     string
	AdminerIP       string

	// Bridge attach list (space-separated container:bridge:ipcidr triples)
	AttachServices string
}

type AnsibleConfig struct {
	TargetIP   string
	SSHUser    string
	SSHKeyPath string
	WGPrivKey  string
	WGPubKey   string
}

// DefaultKWSConfig returns a config with sensible defaults pre-populated.
func DefaultKWSConfig() *KWSConfig {
	return &KWSConfig{
		Server: ServerConfig{
			GatewayPort:     8080,
			CodeServerPort:  8099,
			TunnelProxyPort: 8081,
		},
		Network: NetworkConfig{
			LXDBridgeName:      "lxdbr0",
			LXDBridgeSubnet:    "172.30.0.0/24",
			LXDBridgeGateway:   "172.30.0.1",
			ServicesSubnet:     "172.25.0.0/24",
			ServicesGateway:    "172.25.0.1",
			CoreNetworkSubnet:  "172.35.0.0/24",
			CoreNetworkGateway: "172.35.0.1",
			DNSIP:              "172.30.0.102",
		},
		Wireguard: WireguardConfig{
			InterfaceName: "wg0",
			Address:       "10.0.0.1/24",
			CIDR:          24,
			ListenPort:    51820,
			KeepAliveSec:  25,
		},
		Services: ServicesConfig{
			PostgresIP:       "172.25.0.2",
			PostgresHostname: "postgres.kws.services",
			PostgresPort:     5432,
			AdminerIP:        "172.25.0.4",
			AdminerHostname:  "adminer.kws.services",
			AdminerPort:      8080,
			BridgeAttach: []BridgeAttach{
				{Container: "postgres.kws.services", Bridge: "lxdbr0", IPCIDR: "172.30.0.100/24"},
				{Container: "adminer.kws.services", Bridge: "lxdbr0", IPCIDR: "172.30.0.101/24"},
				{Container: "dnsmasq_kws", Bridge: "lxdbr0", IPCIDR: "172.30.0.102/24"},
			},
		},
		Instance: InstanceConfig{
			MemoryLimit: "1500MB",
			StoragePool: "kws",
			UbuntuAlias: "ubuntu-22:04",
		},
		Limits: LimitsConfig{
			MaxWGDevicesPerUser:   3,
			MaxServiceDBUsers:     5,
			MaxServiceDBDatabases: 10,
			UserDomainLimit:       3,
		},
		Paths: PathsConfig{
			LXDSocket:    "/var/snap/lxd/common/lxd/unix.socket",
			NginxConfDir: "/app/nginx_conf",
		},
	}
}

// DefaultEnvConfig returns env config with sensible defaults.
func DefaultEnvConfig() *EnvConfig {
	return &EnvConfig{
		DBUsername:         "kws",
		DBPassword:         "",
		DBName:             "kws_db",
		DBHost:             "localhost",
		DBPort:             "5432",
		RedisHost:          "localhost",
		RedisPort:          "6379",
		RedisPassword:      "",
		Env:                "production",
		MQUser:             "mq_user",
		MQPassword:         "",
		MQServerPort:       "5672",
		MQUIPort:           "15672",
		MQHost:             "localhost",
		WGPrivateKey:       "",
		PGServiceUsername:  "pgadmin",
		PGServicePassword:  "",
		PGServiceHost:      "postgres.kws.services",
		PGServicePort:      "5433",
		PGServiceDB:        "pg_service",
		SMTPHost:           "smtp.gmail.com",
		SMTPPort:           "587",
		ServicesSubnet:     "172.25.0.0/24",
		ServicesGateway:    "172.25.0.1",
		PGServiceIP:        "172.25.0.2",
		AdminerIP:          "172.25.0.4",
		AttachServices:     "postgres.kws.services:lxdbr0:172.30.0.100/24 adminer.kws.services:lxdbr0:172.30.0.101/24 dnsmasq_kws:lxdbr0:172.30.0.102/24",
	}
}

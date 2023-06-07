package config

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Config struct {
	Production                  bool                  `json:"production"`
	AppName                     string                `json:"appName"`
	ApiVersion                  string                `json:"apiVersion"`
	KubernetesServiceNameSuffix string                `json:"kubernetesServiceNameSuffix"`
	ListenPort                  string                `json:"listenPort"`
	Folders                     []string              `json:"folders"`
	SecondsToReloadServicesName int                   `json:"secondsToReloadServicesName"`
	MongoDB                     MongoDBConfig         `json:"mongodb"`
	Certificates                CertificatesConfig    `json:"certificates"`
	Token                       TokenConfig           `json:"token"`
	SecurityKeys                SecurityKeysConfig    `json:"securityKeys"`
	SecurityRSAKeys             SecurityRSAKeysConfig `json:"securityRSAKeys"`
	SMTPServer                  SMTPConfig            `json:"smtpServer"`
	Company                     CompanyConfig         `json:"company"`
	Prometheus                  PrometheusConfig      `json:"prometheus"`
	MongoDbExporter             MongoDbExporterConfig `json:"mongoDbExporter"`
	Nats                        NatsConfig            `json:"nats"`
	Jaeger                      JaegerConfig          `json:"jaeger"`
	GrpcServer                  GrpcServerConfig      `json:"grpcServer"`
	EmailService                EmailServiceConfig    `json:"emailService"`
	Postgres                    PostgresConfig        `json:"postgres"`
	Redis                       RedisConfig           `json:"redis"`
	Consul                      ConsulConfig          `json:"consul"`
}

type MongoDBConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	MaxPoolSize int    `json:"maxPoolSize"`
}

type CertificatesConfig struct {
	FolderName                    string `json:"foldername"`
	FileNameCert                  string `json:"filenamecert"`
	FileNameKey                   string `json:"filenamekey"`
	HashPermissionEndPoint        string `json:"hashPermissionEndPoint"`
	PasswordPermissionEndPoint    string `json:"passwordPermissionEndPoint"`
	ServiceName                   string `json:"serviceName"`
	APIPathCertificateCA          string `json:"apiPathCertificateCA"`
	EndPointGetCertificateCA      string `json:"endPointGetCertificateCA"`
	APIPathCertificateHost        string `json:"apiPathCertificateHost"`
	EndPointGetCertificateHost    string `json:"endPointGetCertificateHost"`
	APIPathCertificateHostKey     string `json:"apiPathCertificateHostKey"`
	EndPointGetCertificateHostKey string `json:"endPointGetCertificateHostKey"`
	MinutesToReloadCertificate    int    `json:"minutesToReloadCertificate"`
}

type TokenConfig struct {
	Issuer                    string `json:"issuer"`
	MinutesToExpireToken      int    `json:"minutesToExpireToken"`
	HoursToExpireRefreshToken int    `json:"hoursToExpireRefreshToken"`
}

type SecurityKeysConfig struct {
	DaysToExpireKeys            int    `json:"daysToExpireKeys"`
	MinutesToRefreshPrivateKeys int    `json:"minutesToRefreshPrivateKeys"`
	MinutesToRefreshPublicKeys  int    `json:"minutesToRefreshPublicKeys"`
	SavePublicKeyToFile         bool   `json:"savePublicKeyToFile"`
	FileECPPublicKey            string `json:"fileECPPublicKey"`
	ServiceName                 string `json:"serviceName"`
	APIPathPublicKeys           string `json:"apiPathPublicKeys"`
	EndPointGetPublicKeys       string `json:"endPointGetPublicKeys"`
}

type SecurityRSAKeysConfig struct {
	DaysToExpireRSAKeys            int    `json:"daysToExpireRSAKeys"`
	MinutesToRefreshRSAPrivateKeys int    `json:"minutesToRefreshRSAPrivateKeys"`
	MinutesToRefreshRSAPublicKeys  int    `json:"minutesToRefreshRSAPublicKeys"`
	ServiceName                    string `json:"serviceName"`
	APIPathRSAPublicKeys           string `json:"apiPathRSAPublicKeys"`
	EndPointGetRSAPublicKeys       string `json:"endPointGetRSAPublicKeys"`
}

type CompanyConfig struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	AddressNumber     string `json:"addressNumber"`
	AddressComplement string `json:"addressComplement"`
	Locality          string `json:"locality"`
	Country           string `json:"country"`
	PostalCode        string `json:"postalCode"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
}

type SMTPConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	TLS          bool   `json:"tls"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SupportEmail string `json:"supportEmail"`
}

type PrometheusConfig struct {
	PROMETHEUS_PUSHGATEWAY string `json:"prometheus_pushgateway"`
}

type MongoDbExporterConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type NatsConfig struct {
	Url         string `json:"url"`
	ClusterId   string `json:"clusterId"`
	ClientId    string `json:"clientId"`
	ConnectWait int    `json:"connectWait"`
	PubAckWait  int    `json:"pubAckWait"`
	Interval    int    `json:"interval"`
	MaxOut      int    `json:"maxOut"`
}

type JaegerConfig struct {
	JaegerEndpoint string `json:"jaegerEndpoint"`
	ServiceName    string `json:"serviceName"`
	ServiceVersion string `json:"serviceVersion"`
}

type GrpcServerConfig struct {
	Port              string `json:"port"`
	MaxConnectionIdle int    `json:"maxConnectionIdle"`
	MaxConnectionAge  int    `json:"maxConnectionAge"`
	Timeout           int    `json:"timeout"`
}

type EmailServiceConfig struct {
	ServiceName string `json:"serviceName"`
	Host        string `json:"host"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslMode"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	PoolSize int    `json:"poolSize"`
}

type ConsulConfig struct {
	Host string `json:"host"`
}

func LoadConfig(production bool, path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config-dev")
	if production {
		viper.SetConfigName("config-prod")
	}
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %s", err))
	}

	config.Production = production

	if config.Production {
		HASH := "HASHPERMISSIONENDPOINT"
		PASSWORD := "PASSWORDPERMISSIONENDPOINT"
		MONGO_USER := "MONGO_USER"
		MONGO_PASSWORD := "MONGO_PASSWORD"
		POSTGRES_USER := "POSTGRES_USER"
		POSTGRES_PASSWORD := "POSTGRES_PASSWORD"
		REDIS_PASSWORD := "REDIS_PASSWORD"
		MONGO_EXPORTER_USER := "MONGO_EXPORTER_USER"
		MONGO_EXPORTER_PASSWORD := "MONGO_EXPORTER_PASSWORD"
		SMTP_SERVER_USER := "SMTP_SERVER_USER"
		SMTP_SERVER_PASSWORD := "SMTP_SERVER_PASSWORD"
		PROMETHEUS_PUSHGATEWAY := "PROMETHEUS_PUSHGATEWAY"
		if checkEnvFile() {
			viper.Reset()
			viper.SetConfigFile(".env")
			viper.ReadInConfig()
			config.Certificates.HashPermissionEndPoint = viper.GetString(HASH)
			config.Certificates.PasswordPermissionEndPoint = viper.GetString(PASSWORD)
			config.MongoDB.User = viper.GetString(MONGO_USER)
			config.MongoDB.Password = viper.GetString(MONGO_PASSWORD)
			config.Postgres.User = viper.GetString(POSTGRES_USER)
			config.Postgres.Password = viper.GetString(POSTGRES_PASSWORD)
			config.Redis.Password = viper.GetString(REDIS_PASSWORD)
			config.MongoDbExporter.User = viper.GetString(MONGO_EXPORTER_USER)
			config.MongoDbExporter.Password = viper.GetString(MONGO_EXPORTER_PASSWORD)
			config.SMTPServer.Username = viper.GetString(SMTP_SERVER_USER)
			config.SMTPServer.Password = viper.GetString(SMTP_SERVER_PASSWORD)
			config.Prometheus.PROMETHEUS_PUSHGATEWAY = viper.GetString(PROMETHEUS_PUSHGATEWAY)
		} else {
			config.Certificates.HashPermissionEndPoint = os.Getenv(HASH)
			config.Certificates.PasswordPermissionEndPoint = os.Getenv(PASSWORD)
			config.MongoDB.User = os.Getenv(MONGO_USER)
			config.MongoDB.Password = os.Getenv(MONGO_PASSWORD)
			config.Postgres.User = os.Getenv(POSTGRES_USER)
			config.Postgres.Password = os.Getenv(POSTGRES_PASSWORD)
			config.Redis.Password = os.Getenv(REDIS_PASSWORD)
			config.MongoDbExporter.User = os.Getenv(MONGO_EXPORTER_USER)
			config.MongoDbExporter.Password = os.Getenv(MONGO_EXPORTER_PASSWORD)
			config.SMTPServer.Username = os.Getenv(SMTP_SERVER_USER)
			config.SMTPServer.Password = os.Getenv(SMTP_SERVER_PASSWORD)
			config.Prometheus.PROMETHEUS_PUSHGATEWAY = os.Getenv(PROMETHEUS_PUSHGATEWAY)
		}

		config.Nats.ClientId += "_" + uuid.New().String()

		fmt.Printf("ENVIRONMENT: production\n")
	}

	if len(string(config.Certificates.PasswordPermissionEndPoint)) == 0 {
		log.Fatal("PasswordPermissionEndPoint cannot be a null or empty value")
	}

	fmt.Printf("MONGO_HOST: %s\nMONGO_PORT: %s\n", config.MongoDB.Host, config.MongoDB.Port)
	fmt.Printf("NATS_URL: %s\n", config.Nats.Url)

	return config
}

func checkEnvFile() bool {
	info, err := os.Stat(".env")
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

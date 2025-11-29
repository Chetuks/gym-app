package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

const ProfileEnvKey string = "ACTIVE_PROFILE"
const PageSize int64 = 50

var configs *Configurations

/* The Configuration struct captrues the application configurations
This struct can be extended to capture any application specific configurations
A pointer to the configs can be obtained using the Getconfigs() method
*/

type Configurations struct {
	Server    *ServerConfigurations
	AppConfig *appConfigurations
}

type TopCategoryConfig struct {
	AutoAppointmentEligible []string `yaml:"autoAppointmentEligible"`
}

type AuthConfiguration struct {
	TokenUrl                  string
	ExpiryLimit               float64
	ClientName                string
	DCConfigClientId          string
	OutboundRetryClientId     string
	DCConfigClientSecret      string
	OutboundRetryClientSecret string
	JwksUrl                   string
	Issuer                    string
}

//	type AuthConfig struct {
//		JwksUrl                       string
//		PermissionUrl                 string
//		JWTValidationCacheSizeKey     string
//		JWTValidationCacheTTLKey      string
//		DefaultJWTValidationCacheSize string
//		DefaultJWTValidationCacheTTL  string
//	}
type EmailNotificationConfigurations struct {
	EmailApiKey string
	ToName      string
	ToAddress   string
	FromName    string
	FromAddress string
	Text        string
	IsEnabled   bool
	Subjects    Subjects
}

type Subjects struct {
	CreateSubject  string
	UpdateSubject  string
	DeletedSubject string
}
type appConfigurations struct {
	PermissionModPanic bool
	MaxFilterLimit     int16
	ScopeAuthSecret    string
	ScopeAuthId        string
}

type ServerConfigurations struct {
	Name      string
	Port      int
	PreFork   bool
	RateLimit int
	ExpSecs   int
}
type ClientConfigurations struct {
	//ExtUserInfoService       *ServiceConfigurations
	PolicyService            *ServiceConfigurations
	OutButRetryService       *ServiceConfigurations
	EmailNotificationService *EmailServiceConfigurations
	ChatNotificationService  *ChatNotificationConfigurations
	AppointmentConfigService *ServiceConfigurations
	MDMService               *ServiceConfigurations
	PdfGeneratorService      *ServiceConfigurations
	ASNService               *ServiceConfigurations
}

type ServiceConfigurations struct {
	ServiceName string
	Host        string
	BaseURL     string
}

type EmailServiceConfigurations struct {
	ServiceName     string
	Host            string
	BaseURL         string
	EmailRecipients string
	EmailVersion    string
	Profile         string
}

type ChatNotificationConfigurations struct {
	ServiceName string
	Host        string
	BaseURL     string
	AppName     string
	Application string
	Profile     string
}

/*
Redis configurations
*/
type redisConfigurations struct {
	Host       string
	Port       int
	Sentinels  []string
	Type       string
	PoolSize   int
	User       string
	Password   string
	Namespaces map[string]Storage
}

type Storage struct {
	Name string
	Ttl  string
}

/*
Mongo configurations
*/
type mongoConfigurations struct {
	Url           string
	DbName        string
	CreateIndexes bool
	Indexes       []IndexModel
}

type IndexModel struct {
	Name       string
	Collection string
	Fields     []Field
	Unique     bool
}

type Field struct {
	Name string
	Type int
}

type cryptoConfigurations struct {
	KeyRotator KeyRotator
	Jwt        Jwt
}

type KeyRotator struct {
	Enabled bool
	Cron    string
}

type Jwt struct {
	Issuer      string
	AtExpiryTtl string
	RtExpiryTtl string
}

type authConfigurations struct {
	MaxAttempts int
	MockOtp     bool
}

type kafkaConfigurations struct {
	BootstrapServers        string
	ConsumerCreatePO        string
	ConsumerUpdatePO        string
	ConsumerCreateGrn       string
	ConsumerDcpUpdatePO     string
	ProducerUpdateASN       string
	ProducerUpdate          string
	ProducerAutoAppointment string
	Events                  string
	ProducerUpdatePR        string
	AcknowledgmentPO        string
	ConsumerGroup           string
	SslKeyStoreLocation     string
	SaslUsername            string
	SaslPassword            string

	//ServiceCertLocation string
	//ServiceKeyLocation  string
	CaCertLocation  string
	Protocol        string
	Mechanism       string
	PurchaseOrderBQ *PurchaseOrderBQ
}

type PurchaseOrderBQ struct {
	Topic             string
	DataType          string
	DataTypeArticle   string
	DataTypeHeader    string
	DataTypePOArticle string
}

type RetryConfigurations struct {
	MaxAttempts   int
	RetryInterval int
}

type CronJobScheduleConfigurations struct {
	JobScheduleTimeInterval string
	RedisLockKey            string
	RedisLockTTL            int
}

type PdfConfigurations struct {
	Interval     string
	TemplateId   string
	BucketConfig BucketConfigurations
	RedisLockKey string
	RedisLockTTL int
}

type BucketConfigurations struct {
	Bucket    string
	SrcFolder string
}

func init() {

	configFileName := "config"
	env, profileSet := os.LookupEnv(ProfileEnvKey)
	if profileSet {
		log.Info().Str("Setting active profile to ", env).Send()
		configFileName = configFileName + "-" + env
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		//set default to trace
		lvl = "INFO"

	}

	switch lvl {
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Info().Msg("Completed setting log level to " + lvl)

	if len(env) == 0 {
		log.Info().Msg("Setting Active Profile to NONE")
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// viper.SetConfigFile("/Users/chethan/Dlabs_Workspace/Projects/dlnk-api-po/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Error().AnErr("Error reading config file", err).Send()
		panic(-1)
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			val := getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"))
			viper.Set(k, val)
		}
	}
	err := viper.Unmarshal(&configs)
	if err != nil {
		log.Error().AnErr("Error creating configs", err).Send()
		panic(-1)
	}

	log.Info().Msg("Configurations loaded successfully")

}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		log.Panic().Msg("Config load failed")
		panic("Mandatory env variable not found:" + env)
	}
	return res
}

func GetConfig() *Configurations {
	return configs
}

func GetLoginDetails(c *fiber.Ctx) error {
	log.Info().Msg("GetLoginDetails called")
	sampleJson := `
	{
		"waterGoal": 164,
        "waterConsumed": 118,
        "stepsGoal": 17400,
        "stepsDone": 13572,
        "exercisesCals": 0,
        "exercisesHours": 0,
        "sleepGoal": 7,
        "sleepDone": 6
	}`
	log.Info().Msg("Login details fetched: " + sampleJson)
	return c.Status(fiber.StatusOK).SendString(sampleJson)
}

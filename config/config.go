package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	ServiceHost string
	HTTPPort    string
	HTTPScheme  string

	InfluxURL   string
	InfluxToken string

	InfluxOrg    string
	InfluxBucket string

	SecretKey string

	PasscodePool   string
	PasscodeLength int

	DefaultOffset string
	DefaultLimit  string

	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

// Load ...
func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

    config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "BeatSync"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	// config.ServiceHost = cast.ToString(getOrReturnDefaultValue("SERVICE_HOST", "39.101.69.244"))
	config.ServiceHost = cast.ToString(getOrReturnDefaultValue("SERVICE_HOST", "localhost"))

	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8080"))
	config.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_SCHEME", "http"))

	config.InfluxURL = cast.ToString(getOrReturnDefaultValue("INFLUX_URL", "http://localhost:8086"))
	config.InfluxToken = cast.ToString(getOrReturnDefaultValue("INFLUX_TOKEN", "WtNUJsTl0qDSOzQa1OO4mnlaPh1tGOtSC34-LvdAoKzNZNXk9OgkjbyU3W4QMCTxogmsTYk3FLX0JjtBSwqE3A=="))
	config.InfluxOrg = cast.ToString(getOrReturnDefaultValue("INFLUX_ORG", "beatsync"))
	config.InfluxBucket = cast.ToString(getOrReturnDefaultValue("INFLUX_BUCKET", "beatsync-bucket"))
	
	config.SecretKey = cast.ToString(getOrReturnDefaultValue("SECRET_KEY", "Here$houldBe$ome$ecretKey"))

	config.PasscodePool = cast.ToString(getOrReturnDefaultValue("PASSCODE_POOL", "0123456789"))
	config.PasscodeLength = cast.ToInt(getOrReturnDefaultValue("PASSCODE_LENGTH", "6"))


	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10"))

	// Load email configuration
	config.SMTPHost = cast.ToString(getOrReturnDefaultValue("SMTP_HOST", "smtp.gmail.com"))
	config.SMTPPort = cast.ToInt(getOrReturnDefaultValue("SMTP_PORT", 587))
	config.SMTPUsername = cast.ToString(getOrReturnDefaultValue("SMTP_USERNAME", ""))
	config.SMTPPassword = cast.ToString(getOrReturnDefaultValue("SMTP_PASSWORD", ""))
	config.SMTPFrom = cast.ToString(getOrReturnDefaultValue("SMTP_FROM", ""))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}

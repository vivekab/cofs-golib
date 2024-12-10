package golibconfig

import (
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
)

var envVars []string

var skipList = []string{"CONNECT_ENV", "LOCAL_RUNNING_SERVICES"}

func SetupConfig(envs ...string) {
	viper.AutomaticEnv()
	viper.SetDefault("API_ENV", "local")
	viper.SetConfigType("yml")
	viper.SetConfigName("application-" + viper.GetString("API_ENV"))
	viper.AddConfigPath("common-config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read common config file, err = %s", err.Error())
	}
	viper.SetConfigFile("config/application-" + viper.GetString("API_ENV") + ".yml")
	err = viper.MergeInConfig()
	if err != nil {
		log.Fatalf("failed to read config file, err = %s", err.Error())
	}
	CheckEnv(envs...)
}

func SetupTestConfig(envs ...string) {
	viper.AutomaticEnv()
	viper.SetDefault("API_ENV", "local")
	CheckEnv(envs...)
}

func CheckEnv(envs ...string) {
	for _, key := range envs {
		envVars = append(envVars, key)
		if GetString(key) == "" {
			log.Fatalf("Env variable %s is not set", key)
		}
	}
}

func isLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func validateEnv(env string) {
	// skip the config variables defined in files
	if strings.Contains(env, ".") || isLower(env) {
		return
	}
	for _, key := range skipList {
		if key == env {
			return
		}
	}
	for _, key := range envVars {
		if key == env {
			return
		}
	}
	log.Fatalf("Env variable %s is not added to checkEnv list", env)
}

// returns the value associated with the key as string
func GetString(key string) string {
	validateEnv(key)
	return viper.GetString(key)
}

// returns the value associated with the key as integer
func GetInt(key string) int {
	validateEnv(key)
	return viper.GetInt(key)
}

// returns the value associated with the key as bool
func GetBool(key string) bool {
	validateEnv(key)
	return viper.GetBool(key)
}

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 {
	validateEnv(key)
	return viper.GetInt32(key)
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 {
	validateEnv(key)
	return viper.GetInt64(key)
}

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint {
	validateEnv(key)
	return viper.GetUint(key)
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 {
	validateEnv(key)
	return viper.GetUint32(key)
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 {
	validateEnv(key)
	return viper.GetUint64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {
	validateEnv(key)
	return viper.GetFloat64(key)
}

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time {
	validateEnv(key)
	return viper.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration {
	validateEnv(key)
	return viper.GetDuration(key)
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func GetIntSlice(key string) []int {
	validateEnv(key)
	return viper.GetIntSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string {
	validateEnv(key)
	return viper.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} {
	validateEnv(key)
	return viper.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string {
	validateEnv(key)
	return viper.GetStringMapString(key)
}

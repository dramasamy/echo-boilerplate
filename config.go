package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	AppName             string
	Address             net.IP
	Port                uint
	LogFile             string
	LogFormat           string
	LogLevel            string
	LogRequestsDisabled bool
	MongoDBUri          string
	SysAdminEmail       string
	SysAdminPassword    string
	SysAdminUsername    string
	TokenSecret         string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		AppName:             "boilerplate",
		Address:             net.ParseIP("127.0.0.1"),
		Port:                1323,
		LogFile:             "stdout",
		LogFormat:           "text",
		LogLevel:            "info",
		LogRequestsDisabled: false,
		MongoDBUri:          "mongodb://",
		SysAdminEmail:       "admin@example.com",
		SysAdminPassword:    "",
		SysAdminUsername:    "admin",
		TokenSecret:         "changeme",
	}
	return &cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.AppName, "app-name", cnf.AppName, "The name of the application.")
	fs.IPVar(&cnf.Address, "address", cnf.Address, "The IP address to listen at.")
	fs.UintVar(&cnf.Port, "port", cnf.Port, "The port to listen at.")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr and 'null' means discard log messages.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat,
		"The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid log levels: debug, info, warning, error and critical.")
	fs.BoolVar(&cnf.LogRequestsDisabled, "log-requests-disabled", cnf.LogRequestsDisabled,
		"Disables HTTP requests logging.")
	fs.StringVar(&cnf.MongoDBUri, "mongodb-uri", cnf.MongoDBUri, "MongoDB URI")
	fs.StringVar(&cnf.SysAdminEmail, "sysadmin-email", cnf.SysAdminEmail,
		"Email address for the system administrator user.")
	fs.StringVar(&cnf.SysAdminPassword, "sysadmin-password", cnf.SysAdminPassword,
		"Password for the system administrator user. Empty means one will be generated automatically.")
	fs.StringVar(&cnf.SysAdminUsername, "sysadmin-username", cnf.SysAdminUsername,
		"Username for the system administrator user.")
	fs.StringVar(&cnf.TokenSecret, "token-secret", cnf.TokenSecret, "Secret to sign access tokens with.")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.addFlags(pflag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		m := fmt.Sprintf("Error binding flags: %v", err)
		log.Panic(m)
		panic(m)
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	n := viper.GetString("app-name")
	if len(n) < 1 {
		m := fmt.Sprint("Application name cannot be empty!")
		log.Panic(m)
		panic(m)
	}

	viper.SetEnvPrefix(n)
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

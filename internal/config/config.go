package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host               string                `yaml:"host"`
	Port               int                   `yaml:"port"`
	CertFile           string                `yaml:"cert_file,omitempty"`
	KeyFile            string                `yaml:"key_file,omitempty"`
	CORS               *CORSConfig           `yaml:"cors"`
	DefaultRedirectUrl string                `yaml:"default_redirect_url"`
	StaticDirs         []*StaticDirConfig    `yaml:"static_dirs"`
	ReverseProxies     []*ReverseProxyConfig `yaml:"reverse_proxies"`
	Redirects          []*RedirectConfig     `yaml:"redirects"`
}

type CORSConfig struct {
	AllowOrigin string `yaml:"allow_origin"`
}

type StaticDirConfig struct {
	RequestPathPrefix string `yaml:"request_path_prefix"`
	LocalPath         string `yaml:"local_path"`
	MimeType          string `yaml:"mime_type"`
}

type ReverseProxyConfig struct {
	RequestPathPrefix string `yaml:"request_path_prefix"`
	TargetUrl         string `yaml:"target_url"`
	TimeoutSec        int    `yaml:",omitempty"`
}

type RedirectConfig struct {
	RequestPathPrefix string `yaml:"request_path_prefix"`
	TargetUrl         string `yaml:"target_url"`
	PassSubPath       bool   `yaml:",omitempty"`
	PassQuery         bool   `yaml:",omitempty"`
}

var value atomic.Pointer[Config]

func V() *Config {
	return value.Load()
}

func ParseFlag() Stopper {
	createOnly := flag.Bool("c", false, "create config file and exit")
	filePath := flag.String("f", "./config.yaml", "run with config file path")
	flag.Parse()

	if len(flag.Args()) > 0 {
		log.Println("invalid parameter")
		flag.Usage()
		os.Exit(1)
	}

	if *createOnly {
		createConfigFile(*filePath)
		panic("never be reached") // because os.Exit() is called in createConfigFile()
	} else {
		modTime, err := getModTime(*filePath)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}

		newCfg, err := readConfigFile(*filePath)
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		setV(newCfg)

		return startWatch(*filePath, modTime)
	}
}

func setV(newCfg *Config) {
	if value.CompareAndSwap(nil, newCfg) {
		return
	}

	cfg := V()

	// The following variables cannot update during execution.
	newCfg.Host = cfg.Host
	newCfg.Port = cfg.Port
	newCfg.CertFile = cfg.CertFile
	newCfg.KeyFile = cfg.KeyFile

	value.Store(newCfg)
}

func getModTime(cfgFilePath string) (time.Time, error) {
	cfgStat, err := os.Stat(cfgFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, fmt.Errorf("config file not found: %w", err)
		}
		return time.Time{}, fmt.Errorf("failed to stat config file: %w", err)
	}
	return cfgStat.ModTime(), nil
}

func readConfigFile(cfgFilePath string) (*Config, error) {
	newCfg := &Config{}

	cfgFile, err := os.Open(cfgFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer cfgFile.Close()

	err = yaml.NewDecoder(cfgFile).Decode(newCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return newCfg, nil
}

func createConfigFile(cfgFilePath string) {
	emptyConfig := Config{
		Host:     "localhost",
		Port:     9090,
		CertFile: "./cert.pem",
		KeyFile:  "./key.pem",
		CORS: &CORSConfig{
			AllowOrigin: "*",
		},
		DefaultRedirectUrl: "https://www.example.com",
		StaticDirs: []*StaticDirConfig{
			{
				RequestPathPrefix: "/file",
				LocalPath:         "./static",
				MimeType:          "text/plain",
			},
		},
		ReverseProxies: []*ReverseProxyConfig{
			{
				RequestPathPrefix: "/db/get",
				TargetUrl:         "http://localhost:10002/get",
				TimeoutSec:        3,
			},
		},
		Redirects: []*RedirectConfig{
			{
				RequestPathPrefix: "/abc",
				TargetUrl:         "https://other.example.com/def",
				PassSubPath:       true,
				PassQuery:         true,
			},
		},
	}

	cfgFile, err := os.Create(cfgFilePath)
	if err != nil {
		log.Fatalf("failed to create config file: %+v\n", err)
	}
	defer cfgFile.Close()

	enc := yaml.NewEncoder(cfgFile)
	enc.SetIndent(2)
	err = enc.Encode(&emptyConfig)
	if err != nil {
		log.Fatalf("failed to encode config file: %+v\n", err)
	}

	log.Printf("config file created: %s\n", cfgFilePath)
	os.Exit(0)
}

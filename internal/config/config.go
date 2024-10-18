package config

import (
	"flag"
	"log"
	"os"

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
	AllowOrigin  string `yaml:"allow_origin"`
	AllowHeaders string `yaml:"allow_headers"`
}

type StaticDirConfig struct {
	RequestHosts      []string `yaml:"request_hosts"`
	RequestPathPrefix string   `yaml:"request_path_prefix"`
	LocalPath         string   `yaml:"local_path"`
	MimeType          string   `yaml:"mime_type"`
}

type ReverseProxyConfig struct {
	RequestUrl string `yaml:"request_url"`
	TargetUrl  string `yaml:"target_url"`
	TimeoutSec int    `yaml:",omitempty"`
}

type RedirectConfig struct {
	RequestUrl  string `yaml:"request_url"`
	TargetUrl   string `yaml:"target_url"`
	PassSubPath bool   `yaml:",omitempty"`
	PassQuery   bool   `yaml:",omitempty"`
}

var Value Config

func ParseFlag() {
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
	} else {
		readConfigFile(*filePath)
	}
}

func readConfigFile(cfgFilePath string) {
	cfgFile, err := os.Open(cfgFilePath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer cfgFile.Close()

	err = yaml.NewDecoder(cfgFile).Decode(&Value)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
}

func createConfigFile(cfgFilePath string) {
	emptyConfig := Config{
		Host:     "localhost",
		Port:     9090,
		CertFile: "./cert.pem",
		KeyFile:  "./key.pem",
		CORS: &CORSConfig{
			AllowOrigin:  "*",
			AllowHeaders: "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		},
		DefaultRedirectUrl: "https://www.example.com",
		StaticDirs: []*StaticDirConfig{
			{
				RequestHosts:      []string{"https://static.example.com:9091"},
				RequestPathPrefix: "/file",
				LocalPath:         "./static",
				MimeType:          "text/plain",
			},
		},
		ReverseProxies: []*ReverseProxyConfig{
			{
				RequestUrl: "https://api.example.com:9092/db/get",
				TargetUrl:  "http://localhost:10002/get",
				TimeoutSec: 3,
			},
		},
		Redirects: []*RedirectConfig{
			{
				RequestUrl:  "https://some.example.com/abc",
				TargetUrl:   "https://other.example.com/def",
				PassSubPath: true,
				PassQuery:   true,
			},
		},
	}

	cfgFile, err := os.Create(cfgFilePath)
	if err != nil {
		log.Fatalf("failed to create config file: %+v\n", err)
	}
	defer cfgFile.Close()

	err = yaml.NewEncoder(cfgFile).Encode(&emptyConfig)
	if err != nil {
		log.Fatalf("failed to encode config file: %+v\n", err)
	}

	log.Printf("config file created: %s\n", cfgFilePath)
	os.Exit(0)
}

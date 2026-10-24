package bunny

import (
	"os"
	"strings"
)

type Config struct {
	APIKey        string
	StorageZone   string
	Region        string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	AccountAPIKey string
	Local         bool
}

func LoadConfig() Config {
	region := strings.ToLower(strings.TrimSpace(os.Getenv("BUNNY_REGION")))
	if region == "" {
		region = "sg"
	}
	driver := strings.ToLower(strings.TrimSpace(os.Getenv("STORAGE_DRIVER")))
	return Config{
		APIKey:        strings.TrimSpace(os.Getenv("BUNNY_API_KEY")),
		StorageZone:   strings.TrimSpace(os.Getenv("BUNNY_STORAGE_ZONE")),
		Region:        region,
		S3Endpoint:    strings.TrimSpace(os.Getenv("BUNNY_S3_ENDPOINT")),
		S3AccessKey:   strings.TrimSpace(os.Getenv("BUNNY_S3_ACCESS_KEY")),
		S3SecretKey:   strings.TrimSpace(os.Getenv("BUNNY_S3_SECRET_KEY")),
		AccountAPIKey: strings.TrimSpace(os.Getenv("BUNNY_ACCOUNT_API_KEY")),
		Local:         driver == "local",
	}
}

func (c Config) StorageReady() bool {
	if c.Local {
		return true
	}
	return c.APIKey != "" && c.StorageZone != ""
}

func (c Config) S3Ready() bool {
	return c.StorageReady() && !c.Local
}

func (c Config) StatsReady() bool {
	return c.AccountAPIKey != ""
}

func (c Config) storageBaseURL() string {
	if c.Region == "" || c.Region == "de" {
		return "https://storage.bunnycdn.com"
	}
	return "https://" + c.Region + ".storage.bunnycdn.com"
}

func (c Config) s3Endpoint() string {
	if c.S3Endpoint != "" {
		return strings.TrimRight(c.S3Endpoint, "/")
	}
	return "https://" + c.Region + "-s3.storage.bunnycdn.com"
}

func (c Config) s3Bucket() string {
	if c.S3AccessKey != "" {
		return c.S3AccessKey
	}
	return c.StorageZone
}

func (c Config) s3Secret() string {
	if c.S3SecretKey != "" {
		return c.S3SecretKey
	}
	return c.APIKey
}

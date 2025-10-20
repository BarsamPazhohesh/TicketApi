// Package config
package config

import (
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port               int   `yaml:"port"`                  // Application HTTP port
		MaxJsonRequestSize int64 `yaml:"max_json_request_size"` // Max Json Size Can User Send via api Request body KB
		MaxUploadFilesSize int64 `yaml:"max_upload_files_size"` // Max files Size Can User Send via api KB
	} `yaml:"app"`

	CORSConfig struct {
		AllowOrigins     []string `yaml:"allow_origins"`
		AllowMethods     []string `yaml:"allow_methods"`
		AllowHeaders     []string `yaml:"allow_headers"`
		ExposeHeaders    []string `yaml:"expose_headers"`
		AllowCredentials bool     `yaml:"allow_credentials"`
		MaxAgeHours      int      `yaml:"max_age_hours"`
	} `yaml:"cors"`

	Captcha struct {
		Length           int    `yaml:"length"`             // Number of characters in the captcha
		TimeoutMinutes   int    `yaml:"timeout_minutes"`    // In-memory captcha TTL (minutes)
		ExpiredTimeToken int    `yaml:"expired_time_token"` // Captcha JWT token TTL (minutes)
		CookieName       string `yaml:"cookie_name"`        // Cookie name for captcha JWT token
		CleanupInterval  int    `yaml:"cleanup_interval"`   // Cleanup interval for expired captchas (minutes)
		ImageWidth       int    `yaml:"image_width"`        // Captcha image width in pixels
		ImageHeight      int    `yaml:"image_height"`       // Captcha image height in pixels
		MaxCachedCaptcha int    `yaml:"max_cashed_captcha"` // Maximum number of captchas stored in cache
		ValidateIP       bool   `yaml:"validate_ip"`        // Check that user IP matches the one stored in the CAPTCHA token
	} `yaml:"captcha"`

	Token struct {
		HTTPOnly bool `yaml:"httponly"` // Use HttpOnly flag for cookies
		Secure   bool `yaml:"secure"`   // Use Secure flag for cookies
	} `yaml:"token"`

	APIKey struct {
		Size int `yaml:"size"`
	} `yaml:"api_key"`

	OneTimeToken struct {
		CleanupInterval  int `yaml:"cleanup_interval"`   // Cleanup interval for expired one-time tokens (minutes)
		MaxCachedTokens  int `yaml:"max_cashed_tokens"`  // Maximum number of one-time tokens stored in cache
		ExpiredTimeToken int `yaml:"expired_time_token"` // One-time JWT token TTL (minutes)
	} `yaml:"one_time_token"`

	Cache struct {
		TicketTypeTTL   int64 `yaml:"ticket_type_ttl_minutes"`
		DepartmentTTL   int64 `yaml:"department_ttl_minutes"`
		TicketStatusTTL int64 `yaml:"ticket_status_ttl_minutes"`
	} `yaml:"cache"`

	Auth struct {
		ExpiredTimeToken int    `yaml:"expired_time_token"` // Auth JWT token TTL (minutes)
		CookieName       string `yaml:"cookie_name"`        // Cookie name for auth JWT token
	} `yaml:"auth"`

	Mongo struct {
		Enable               bool   `yaml:"enable"`                  // Enable MongoDB integration
		DBName               string `yaml:"db_name"`                 // MongoDB database name
		TicketCollectionName string `yaml:"ticket_collocation_name"` // MongoDB collection name for tickets
	} `yaml:"mongo"`

	Redis struct {
		Enable bool   `yaml:"enable"` // Enable Redis integration
		Host   string `yaml:"host"`   // Redis host
		Port   int    `yaml:"port"`   // Redis port
		DB     int    `yaml:"db"`     // Redis logical database (integer 0 - 15)
	} `yaml:"redis"`

	Minio struct {
		Enable bool   `yaml:"enable"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Bucket string `yaml:"bucket"`
		UseSSL bool   `yaml:"use_ssl"`
	} `yaml:"minio"`

	TicketConfig struct {
		MaxPagingSize            int      `yaml:"max_paging_size"`
		MinPagingSize            int      `yaml:"min_paging_size"`
		DefaultPagingSize        int      `yaml:"default_paging_size"`
		MaxCountingItem          int64    `yaml:"max_counting_item"`
		MaxTicketUploadFile      int      `yaml:"max_ticket_upload_file"`
		MaxTicketUploadFileSize  int64    `yaml:"max_ticket_upload_file_size"`
		AcceptableFilesForUpload []string `yaml:"acceptable_files_for_upload"`
	} `yaml:"ticket"`
}

var (
	cfg  *Config
	mu   sync.RWMutex
	path string
)

// Load reads config file initially.
func Load(p string) *Config {
	path = p
	cfg = &Config{}
	reload()
	watchFile()
	return cfg
}

// reload reads the file and updates cfg.
func reload() {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config: read file error: %v", err)
		return
	}

	var newCfg Config
	if err := yaml.Unmarshal(file, &newCfg); err != nil {
		log.Printf("config: yaml unmarshal error: %v", err)
		return
	}
	cfg = &newCfg
	log.Println("config reloaded")
}

// Get returns a copy-safe config.
func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return *cfg
}

// watchFile watches for changes.
func watchFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					reload()
				}
			case err := <-watcher.Errors:
				log.Println("config watch error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
}

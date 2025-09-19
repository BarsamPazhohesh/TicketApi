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
		Port int `yaml:"port"` // Application HTTP port
	} `yaml:"app"`

	Captcha struct {
		Length           int    `yaml:"length"`             // Number of characters in the captcha
		TimeoutMinutes   int    `yaml:"timeout_minutes"`    // In-memory captcha TTL (minutes)
		ExpiredTimeToken int    `yaml:"expired_time_token"` // Captcha JWT token TTL (minutes)
		CookieName       string `yaml:"cookie_name"`        // Cookie name for captcha JWT token
		CleanupInterval  int    `yaml:"cleanup_interval"`   // Cleanup interval for expired captchas (minutes)
		ImageWidth       int    `yaml:"image_width"`        // Captcha image width in pixels
		ImageHeight      int    `yaml:"image_height"`       // Captcha image height in pixels
		MaxCachedCaptcha int    `yaml:"max_cashed_captcha"` // Maximum number of captchas stored in cache
	} `yaml:"captcha"`

	Token struct {
		HttpOnly bool `yaml:"httponly"` // Use HttpOnly flag for cookies
		Secure   bool `yaml:"secure"`   // Use Secure flag for cookies
	} `yaml:"token"`

	OneTimeToken struct {
		CleanupInterval  int `yaml:"cleanup_interval"`   // Cleanup interval for expired one-time tokens (minutes)
		MaxCachedTokens  int `yaml:"max_cashed_tokens"`  // Maximum number of one-time tokens stored in cache
		ExpiredTimeToken int `yaml:"expired_time_token"` // One-time JWT token TTL (minutes)
	} `yaml:"one_time_token"`

	Auth struct {
		ExpiredTimeToken int    `yaml:"expired_time_token"` // Auth JWT token TTL (minutes)
		CookieName       string `yaml:"cookie_name"`        // Cookie name for auth JWT token
	} `yaml:"auth"`

	Mongo struct {
		Enable               bool   `yaml:"enable"`                  // Enable MongoDB integration
		DBName               string `yaml:"db_name"`                 // MongoDB database name
		TicketCollectionName string `yaml:"ticket_collocation_name"` // MongoDB collection name for tickets
	} `yaml:"mongo"`

	TicketConfig struct {
		MaxPagingSize     int   `yaml:"max_paging_size"`
		MinPagingSize     int   `yaml:"min_paging_size"`
		DefaultPagingSize int   `yaml:"default_paging_size"`
		MaxCountingItem   int64 `yaml:"max_counting_item"`
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

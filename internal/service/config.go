package service

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// EditorConfig represents the editor configuration
type EditorConfig struct {
	Editor struct {
		Theme         string `json:"theme" mapstructure:"theme"`
		FontSize      int    `json:"fontSize" mapstructure:"fontSize"`
		TabSize       int    `json:"tabSize" mapstructure:"tabSize"`
		WordWrap      bool   `json:"wordWrap" mapstructure:"wordWrap"`
		LineNumbers   bool   `json:"lineNumbers" mapstructure:"lineNumbers"`
		RelativeLines bool   `json:"relativeLines" mapstructure:"relativeLines"`
		Minimap       bool   `json:"minimap" mapstructure:"minimap"`
		StickyScroll  bool   `json:"stickyScroll" mapstructure:"stickyScroll"`
		Vim           struct {
			Enabled     bool   `json:"enabled" mapstructure:"enabled"`
			DefaultMode string `json:"defaultMode" mapstructure:"defaultMode"`
		} `json:"vim" mapstructure:"vim"`
	} `json:"editor" mapstructure:"editor"`
	Terminal struct {
		DefaultShell string `json:"defaultShell" mapstructure:"defaultShell"`
		FontSize     int    `json:"fontSize" mapstructure:"fontSize"`
		FontFamily   string `json:"fontFamily" mapstructure:"fontFamily"`
		Theme        struct {
			Background          string `json:"background" mapstructure:"background"`
			Foreground          string `json:"foreground" mapstructure:"foreground"`
			Cursor              string `json:"cursor" mapstructure:"cursor"`
			SelectionBackground string `json:"selectionBackground" mapstructure:"selectionBackground"`
			SelectionForeground string `json:"selectionForeground" mapstructure:"selectionForeground"`
		} `json:"theme" mapstructure:"theme"`
	} `json:"terminal" mapstructure:"terminal"`
	Keyboard struct {
		CustomBindings map[string]KeyBinding `json:"customBindings" mapstructure:"customBindings"`
	} `json:"keyboard" mapstructure:"keyboard"`
}

// KeyBinding represents a keyboard shortcut configuration
type KeyBinding struct {
	Key       string   `json:"key" mapstructure:"key"`
	Modifiers []string `json:"modifiers" mapstructure:"modifiers"`
}

// ConfigService handles editor configuration
type ConfigService struct {
	config     *EditorConfig
	configPath string
}

// NewConfigService creates a new config service instance
func NewConfigService() (*ConfigService, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".edit4i")
	configPath := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return nil, err
		}
	}

	// Initialize viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config EditorConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &ConfigService{
		config:     &config,
		configPath: configPath,
	}, nil
}

// GetConfig returns the current editor configuration
func (s *ConfigService) GetConfig() *EditorConfig {
	return s.config
}

// OpenConfigFile opens the config file in the editor
func (s *ConfigService) OpenConfigFile() string {
	return s.configPath
}

func createDefaultConfig(path string) error {
	defaultConfig := `editor:
  theme: vs-dark
  fontSize: 14
  tabSize: 4
  wordWrap: true
  lineNumbers: true
  relativeLines: false
  minimap: false
  stickyScroll: false
  vim:
    enabled: false
    defaultMode: normal

terminal:
  defaultShell: ""  # Empty means use system default shell
  fontSize: 14
  fontFamily: "monospace"
  theme:
    background: "#181818"
    foreground: "#c5c8c6"
    cursor: "#528bff"
    selectionBackground: "#3e4451"
    selectionForeground: "#d1d5db"

keyboard:
  customBindings: {}`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

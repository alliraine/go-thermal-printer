package service

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/jonasclaes/go-thermal-printer/pkg/model"
	"github.com/pelletier/go-toml/v2"
)

type ConfigService struct {
	config *model.AppConfig
}

func NewConfigService() (*ConfigService, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml" // Default config file path
	}

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &ConfigService{
		config: config,
	}, nil
}

func (cs *ConfigService) GetConfig() *model.AppConfig {
	return cs.config
}

func (cs *ConfigService) GetServerConfig() *model.ServerConfig {
	return &cs.config.Server
}

func (cs *ConfigService) GetPrinterConfig() *model.PrinterConfig {
	return &cs.config.Printer
}

func loadConfig(path string) (*model.AppConfig, error) {
	// Create config with default values first
	config := &model.AppConfig{}
	setDefaultValues(config)

	// If config file exists, load and override defaults
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config, nil
		}

		return nil, fmt.Errorf("failed to access config file %s: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := toml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML config: %w", err)
	}

	return config, nil
}

func setDefaultValues(config *model.AppConfig) {
	setStructDefaults(reflect.ValueOf(config).Elem())
}

func setStructDefaults(v reflect.Value) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			setStructDefaults(field)
			continue
		}

		// Get default value from tag
		defaultValue := fieldType.Tag.Get("default")
		if defaultValue == "" {
			continue
		}

		// Set the default value based on field type
		switch field.Kind() {
		case reflect.String:
			field.SetString(defaultValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
				field.SetInt(val)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val, err := strconv.ParseUint(defaultValue, 10, 64); err == nil {
				field.SetUint(val)
			}
		case reflect.Bool:
			if val, err := strconv.ParseBool(defaultValue); err == nil {
				field.SetBool(val)
			}
		case reflect.Float32, reflect.Float64:
			if val, err := strconv.ParseFloat(defaultValue, 64); err == nil {
				field.SetFloat(val)
			}
		}
	}
}

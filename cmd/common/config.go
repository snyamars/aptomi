package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Aptomi/aptomi/pkg/config"
	"github.com/Aptomi/aptomi/pkg/lang/yaml"
	log "github.com/sirupsen/logrus"
	vp "github.com/spf13/viper"
)

// ReadConfig reads configuration from CLI flags, default or specified file path into provided config object using the
// provided Viper instance and default configuration directory. It'll be checking if --config provided first and there
// are supported config file types in it if it's directory.
func ReadConfig(viper *vp.Viper, cfg config.Base, defaultConfigDir string) error {
	configFilePath := viper.GetString("config")

	if configFilePath != "" { // if config path provided, use it and don't look for default locations
		configAbsPath, err := filepath.Abs(configFilePath)
		if err != nil {
			return fmt.Errorf("error getting abs path for %s: %s", configFilePath, err)
		}

		err = processConfigAbsPath(viper, configAbsPath)
		if err != nil {
			return err
		}
	} else { // if no config path available, search in default places
		// if there is no default config file - just skip config parsing
		if !isConfigExists(defaultConfigDir) {
			return fmt.Errorf("can't find config file in default config dir: %s", defaultConfigDir)
		}

		viper.AddConfigPath(defaultConfigDir)
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("can't read config: %s", err)
	}

	err := viper.Unmarshal(cfg)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config: %s", err)
	}

	if cfg.IsDebug() {
		log.SetLevel(log.DebugLevel)
	}

	val := config.NewValidator(cfg)
	errValidation := val.Validate()
	if errValidation != nil {
		return fmt.Errorf("error while validating config: %s", errValidation)
	}

	log.Debugf("Config:\n%s", yaml.SerializeObject(cfg))

	return nil
}

func processConfigAbsPath(viper *vp.Viper, path string) error {
	if stat, err := os.Stat(path); err == nil {
		if stat.IsDir() { // if dir provided, use only it
			viper.AddConfigPath(path)
		} else { // if specific file provided, use only it
			viper.SetConfigFile(path)
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("specified config path %s doesn't exist: %s", path, err)
	} else {
		return fmt.Errorf("error while processing specified config path %s: %s", path, err)
	}

	return nil
}

func isConfigExists(configDir string) bool {
	exists := false

	// check all supported config types
	for _, supportedType := range vp.SupportedExts {
		defaultConfigFile := path.Join(configDir, "config."+supportedType)

		if _, err := os.Stat(defaultConfigFile); err == nil {
			exists = true
			break
		}
	}

	return exists
}

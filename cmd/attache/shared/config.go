package shared

import "github.com/spf13/viper"

var config *viper.Viper

func GetConfig() (*viper.Viper, error) {
	if config == nil {
		v, err := parseConfig()
		if err != nil {
			return nil, err
		}
		config = v
	}

	return config, nil
}

func parseConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("attache.json")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pdxfixit/awscollector"
	"github.com/spf13/viper"
)

// config is available for global use within the pkg
var config Config

// CollectorConfig relates to the collector app proper
type CollectorConfig struct {
	Debug      bool `json:"debug" mapstructure:"debug"`
	SampleData bool `json:"sample_data" mapstructure:"sample_data"`
}

// Config is the unmarshaled configuration file
type Config struct {
	AccountConfigs []*awscollector.AccountConfig `json:"accounts" mapstructure:"accounts"`
	Collector      CollectorConfig               `json:"collector" mapstructure:"collector"`
}

func loadConfig() {

	log.Println("Loading configuration...")

	// load the config
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/hostdb-collector-aws")
	viper.AddConfigPath(".")

	// load env vars
	viper.SetEnvPrefix("hostdb_collector_aws")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// read the config file, and handle any errors
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %s", err))
	}

	// unmarshal into our struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(fmt.Errorf("unable to decode into struct, %v", err))
	}

	// debug
	if config.Collector.Debug {
		log.Println(fmt.Sprintf("%+v", config.Collector))
		for _, ac := range config.AccountConfigs {
			log.Println(fmt.Sprintf("%+v", *ac))
		}
	}

	// create session objects per region
	for _, ac := range config.AccountConfigs {
		ac.RegionalSessions = []*awscollector.RegionalSession{}
		ac.RegionalClients = []*awscollector.RegionalClient{}
		for _, region := range ac.Regions {
			ac.RegionalSessions = append(ac.RegionalSessions, awscollector.NewRegionalSession(region, ac.AssumeRole))
			ac.RegionalClients = append(ac.RegionalClients, &awscollector.RegionalClient{Region: region})
		}
	}

}

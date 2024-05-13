package cmd

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type senderConfig struct {
	Broker   string   `mapstructure:"broker"`
	Port     int      `mapstructure:"port"`
	ClientID string   `mapstructure:"clientId"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Topic    string   `mapstructure:"topic"`
	Data     emqxBody `mapstructure:"data"`
}

type emqxBody struct {
}

func runSender(cmd *cobra.Command, args []string) {
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var config senderConfig
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", config)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port))
	opts.SetClientID(config.ClientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(0)
	bytes, _ := json.Marshal(config.Data)
	token := client.Publish(config.Topic, 2, false, bytes)
	token.Wait()
	fmt.Println("Published")
}

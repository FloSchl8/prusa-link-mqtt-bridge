package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"prusa-link-mqtt-bridge/pkg/prusalink"
)

func NewMqttClient(broker, username, password string, port int) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(username)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func PublishStatus(client mqtt.Client, topic string, status *prusalink.PrinterStatus) error {
	payload, err := json.Marshal(status)
	if err != nil {
		return err
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

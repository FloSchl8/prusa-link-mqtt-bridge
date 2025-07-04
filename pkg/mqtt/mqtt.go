package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"prusa-link-mqtt-bridge/pkg/prusalink"
)

func NewMqttClient(broker, username, password string, port int, availabilityTopic string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetWill(availabilityTopic, "offline", 1, true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func Publish(client mqtt.Client, topic string, payload string) error {
	token := client.Publish(topic, 1, true, payload)
	token.Wait()
	return token.Error()
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

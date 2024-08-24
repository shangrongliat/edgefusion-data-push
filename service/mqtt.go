package service

import (
	"fmt"

	"edgefusion-data-push/message"
	"edgefusion-data-push/plugin/config"
	log "edgefusion-data-push/plugin/logs"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

type IMqttService interface {
	Start()
	Close()
}

type MqttService struct {
	client     MQTT.Client
	storage    StorageService
	config     *config.Config
	subscribed bool
}

func (m *MqttService) Close() {
	m.client.Disconnect(250)
}

func NewMqttClient(cfg *config.Config) IMqttService {
	storage, err := NewStorageService(cfg)
	if err != nil {
		log.L().Error("storage service 初始化失败.", log.Error(err))
	}
	return &MqttService{
		storage: storage,
		config:  cfg,
	}
}

func (m *MqttService) Start() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.config.Mqtt.Address)
	opts.SetClientID(m.config.Mqtt.ClientID)
	if m.config.Mqtt.Username != "" && m.config.Mqtt.Password != "" {
		opts.SetUsername(m.config.Mqtt.Username)
		opts.SetPassword(m.config.Mqtt.Password)
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.L().Error("Could not connect to MQTT broker: %v", log.Error(token.Error()))
	}
	m.client = client
	m.subscribeToTopics()
}

func (m *MqttService) subscribeToTopics() {
	/*if m.subscribed {
		return
	}*/
	// 主题订阅
	token := m.client.Subscribe("/ef/msg/ir", 0, func(c MQTT.Client, msg MQTT.Message) {
		// 在这里处理消息
		fmt.Printf("Received from topic: %s\n", msg.Topic())
		var data message.Message
		if err := proto.Unmarshal(msg.Payload(), &data); err != nil {
			fmt.Printf("Unmarshal Error: %s\n", err)
		}
		var da message.InferenceResult
		if err := proto.Unmarshal(data.GetData(), &da); err != nil {
			fmt.Printf("Unmarshal Error: %s\n", err)
		}
		//strings.ReplaceAll(data.Metadata["EF_APP_NAME"], "-", "")
		if len(da.GetImageFrame()) > 0 {
			m.storage.ImageStorage(data.Metadata["EF_NODE_ID"], data.Metadata["EF_APP_NAME"], da.GetImageFrame())
		}
		targets := da.Targets
		for _, target := range targets {
			// 进行消息存储
			m.storage.TimeSeriesStorage(data.Metadata["EF_NODE_ID"], data.Metadata["EF_APP_NAME"], target)
			if len(target.Image) > 0 {
				// 如果长度大于0，则说明有图片
				// 进行图片存储
				m.storage.ImageStorage(data.Metadata["EF_NODE_ID"], data.Metadata["EF_APP_NAME"], target.Image)
			}
		}
	})
	if token.Wait() && token.Error() != nil {
		log.L().Error("Failed to subscribe: %v", log.Error(token.Error()))
	}
	m.subscribed = true
}

// 发布消息
func (m *MqttService) Publish(topic string, qos byte, retained bool, payload interface{}) {
	token := m.client.Publish(topic, qos, retained, payload)
	token.Wait()
}

// 断开连接
func (m *MqttService) Disconnect(quiesce uint) {
	m.client.Disconnect(quiesce)
}

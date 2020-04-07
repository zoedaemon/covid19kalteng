package logs

import (
	"encoding/json"
	"log"

	"covid19kalteng/components/basemodel"

	"github.com/Shopify/sarama"
)

type (
	// NorthstarLib main type
	NorthstarLib struct {
		Host         string
		Topic        string
		Secret       string
		Send         bool
		SaramaConfig *sarama.Config
	}
	// Log struct for client template
	Log struct {
		basemodel.BaseModel
		Tag      string `json:"tag"`
		Note     string `json:"note"`
		UID      string `json:"uid"`
		Username string `json:"username"`
		Level    string `json:"level"`
		Messages string `json:"messages"`
	}
	// Audittrail struct for client template
	Audittrail struct {
		basemodel.BaseModel
		Client   string `json:"client"`
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Roles    string `json:"roles"`
		Entity   string `json:"entity"`
		EntityID string `json:"entity_id"`
		Action   string `json:"action"`
		Original string `json:"original"`
		New      string `json:"new"`
	}
)

// SubmitKafkaLog func
func (n *NorthstarLib) SubmitKafkaLog(l interface{}, model string) (err error) {
	if !n.Send {
		return nil
	}
	if len(model) < 1 {
		model = "log"
	}
	build := kafkaLogBuilder(l, model)

	jMarshal, _ := json.Marshal(build)

	//CONCERN: for performance it's safely shared among goroutine
	kafkaProducer, err := sarama.NewAsyncProducer([]string{n.Host}, n.SaramaConfig)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: n.Topic,
		Value: sarama.StringEncoder(n.Secret + ":" + model + ":" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", n.Topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", n.Topic, err)
	}

	return err
}

func kafkaLogBuilder(l interface{}, model string) (payload map[string]interface{}) {
	inrec, _ := json.Marshal(l)
	json.Unmarshal(inrec, &payload)

	return payload
}

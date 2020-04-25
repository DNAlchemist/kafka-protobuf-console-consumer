package consumer

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/DNAlchemist/kafka-protobuf-console-consumer/protobuf_decoder"
	"strings"
)

type SimpleConsumerGroupHandler struct {
	protobufJSONStringify *protobuf_decoder.ProtobufJSONStringify
	prettyJson            bool
	fromBeginning         bool
	withSeparator         bool
}

func NewSimpleConsumerGroupHandler(protobufJSONStringify *protobuf_decoder.ProtobufJSONStringify, prettyJson bool, fromBeginning bool, withSeparator bool) *SimpleConsumerGroupHandler {
	return &SimpleConsumerGroupHandler{protobufJSONStringify: protobufJSONStringify, prettyJson: prettyJson, fromBeginning: fromBeginning, withSeparator: withSeparator}
}

func (SimpleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (SimpleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h SimpleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message strings.Builder

		message.WriteString("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		jsonString, e := h.protobufJSONStringify.JsonString(msg.Value, h.prettyJson)
		if e != nil {
			fmt.Println(e)
		}

		message.WriteString(jsonString + "\n")
		if h.withSeparator {
			message.WriteString("--------------------------------- end message -----------------------------------------\n")
		}
		fmt.Print(message.String())
		sess.MarkMessage(msg, "")
	}
	return nil
}


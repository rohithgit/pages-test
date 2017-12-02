package kafkaclient

import (
	//"fmt"
	"bytes"
	"fmt"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/require"
)

// Test a failed scenario with broker passed a nil
func TestNewProducer(t *testing.T) {
	require := require.New(t)
	_, err := New_kafka_producer(nil)
	require.Nil(err)
}

// Test a failed scenario with broker passed a nil
func TestNewProducera(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	require := require.New(t)
	_, err := New_kafka_producer(kafkacluster)
	require.Nil(err)
}

func TestNewConsumer(t *testing.T) {
	require := require.New(t)
	_, err := New_kafka_consumer(nil)
	require.Nil(err)
}

func TestNewConsumera(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	require := require.New(t)
	_, err := New_kafka_consumer(kafkacluster)
	require.Nil(err)
}

func TestNewMsgProcessor(t *testing.T) {
	require := require.New(t)
	_, err := NewMessageProcessor(nil)
	require.Nil(err)
}

func TestNewMsgProcessora(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	require := require.New(t)
	_, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
}

func TestNewPublishMessage(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = ""
	var Msg = ""
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err1 := MsgProcessor.NewPublishMsg(topic, Msg)
	require.Nil(err1)
}

func TestNewPublishMessagea(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	var Msg = ""
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err1 := MsgProcessor.NewPublishMsg(topic, Msg)
	require.Nil(err1)
}

func TestNewPublishMessageb(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	var Msg = "Publishing the first message:372016 1122PM"
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err1 := MsgProcessor.NewPublishMsg(topic, Msg)
	require.Nil(err1)
}

func TestNewConsumeMessage(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = ""
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	_, err1 := MsgProcessor.NewConsumeMsg(topic)
	require.Nil(err1)
}

func TestNewConsumeMessagea(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	chanMsg, err1 := MsgProcessor.NewConsumeMsg(topic)
	require.Nil(err1)
	fmt.Println("hello: %s", <-chanMsg)
}

func TestNewPublishMessageAvro(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = ""
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err4 := MsgProcessor.NewPublishMsgAvro(topic, nil)
	require.Nil(err4)
}

func TestNewPublishMessageAvroa(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err4 := MsgProcessor.NewPublishMsgAvro(topic, nil)
	require.Nil(err4)
}

func TestNewPublishMessageAvrob(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	recordSchemaJSON := `
{
  "type": "record",
  "name": "comments",
  "doc:": "A basic schema for storing blog comments",
  "namespace": "com.example",
  "fields": [
    {
      "doc": "Name of user",
      "type": "string",
      "name": "username"
    },
    {
      "doc": "The content of the user's message",
      "type": "string",
      "name": "comment"
    },
    {
      "doc": "Unix epoch time in milliseconds",
      "type": "long",
      "name": "timestamp"
    }
  ]
}
`
	require := require.New(t)
	someRecord, err1 := goavro.NewRecord(goavro.RecordSchema(recordSchemaJSON))
	require.Nil(err1)
	someRecord.Set("username", "Aquaman")
	someRecord.Set("comment", "The Atlantic is oddly cold this morning!")
	someRecord.Set("com.example.timestamp", int64(1082196484))
	codec, err2 := goavro.NewCodec(recordSchemaJSON)
	require.Nil(err2)
	bb := new(bytes.Buffer)
	err3 := codec.Encode(bb, someRecord)
	require.Nil(err3)
	Msg := bb.Bytes()
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	err4 := MsgProcessor.NewPublishMsgAvro(topic, Msg)
	fmt.Printf("Message: %#v", Msg)
	require.Nil(err4)
}

func TestNewConsumeMessagec(t *testing.T) {
	var kafkacluster = []string{"128.107.35.115:9092"}
	var topic = "TestSpectre01"
	recordSchemaJSON := `
{
  "type": "record",
  "name": "comments",
  "doc:": "A basic schema for storing blog comments",
  "namespace": "com.example",
  "fields": [
    {
      "doc": "Name of user",
      "type": "string",
      "name": "username"
    },
    {
      "doc": "The content of the user's message",
      "type": "string",
      "name": "comment"
    },
    {
      "doc": "Unix epoch time in milliseconds",
      "type": "long",
      "name": "timestamp"
    }
  ]
}
`
	require := require.New(t)
	MsgProcessor, err := NewMessageProcessor(kafkacluster)
	require.Nil(err)
	chanMsg, err1 := MsgProcessor.NewConsumeMsg(topic)
	require.Nil(err1)
	codec, err2 := goavro.NewCodec(recordSchemaJSON)
	require.Nil(err2)
	fmt.Println("hello: %s", <-chanMsg)
	encoded := []byte(<-chanMsg)
	bb := bytes.NewBuffer(encoded)
	decoded, err3 := codec.Decode(bb)
	require.Nil(err3)
	fmt.Println(decoded) // default String() representation is JSON
	// but direct access to data is provided
	record := decoded.(*goavro.Record)
	fmt.Println("Record Name:", record.Name)
	fmt.Println("Record Fields:")
	for i, field := range record.Fields {
		fmt.Println(" field", i, field.Name, ":", field.Datum)
	}
}

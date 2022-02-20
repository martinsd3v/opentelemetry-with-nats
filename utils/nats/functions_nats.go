package nats

import (
	"encoding/json"

	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type DataTransferWithTrace struct {
	SpanContext *provider.SpanContext `json:"spanContext"`
	Data        interface{}           `data:"data"`
}

func DataToByte(spanContext *provider.SpanContext, data interface{}) []byte {
	dataToByte := DataTransferWithTrace{
		SpanContext: spanContext,
		Data:        data,
	}
	repondByte, err := json.Marshal(dataToByte)
	if err == nil {
		return repondByte
	}
	return nil
}

func ByteToData(btes []byte, data interface{}) (*provider.SpanContext, error) {
	var ByteToData DataTransferWithTrace
	err := json.Unmarshal(btes, &ByteToData)
	if err == nil {
		btesData, err := json.Marshal(ByteToData.Data)
		if err == nil {
			json.Unmarshal(btesData, &data)
		}
	}
	return ByteToData.SpanContext, err
}

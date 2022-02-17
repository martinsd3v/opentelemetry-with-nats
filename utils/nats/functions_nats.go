package nats

import (
	"encoding/json"
	"opentelemetry/utils/tracer"
)

type DataTransferWithTrace struct {
	SpanContext *tracer.SpanContext `json:"spanContext"`
	Data        interface{}         `data:"data"`
}

func DataToByte(spanContext *tracer.SpanContext, data interface{}) []byte {
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

func ByteToData(btes []byte, data interface{}) (*tracer.SpanContext, error) {
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

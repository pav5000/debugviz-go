package request

import "encoding/json"

type Block struct {
	serviceName  string
	handlerName  string
	requestJSON  json.RawMessage
	responseJSON json.RawMessage
	err          error
}

func New(serviceName, handlerName string) *Block {
	block := &Block{
		serviceName: serviceName,
		handlerName: handlerName,
	}
	return block
}

func (b *Block) LogRequest(request interface{}) *Block {
	if b == nil {
		return nil
	}
	b.requestJSON = marshalJSON(request)
	return b
}

func (b *Block) LogResponse(response interface{}) *Block {
	if b == nil {
		return nil
	}
	b.responseJSON = marshalJSON(response)
	return b
}

func (b *Block) LogError(err error) *Block {
	if b == nil {
		return nil
	}
	b.err = err
	return b
}

func (b *Block) Type() string {
	return "external_request"
}

func (b *Block) Data() interface{} {
	if b == nil {
		return nil
	}
	var errText string
	if b.err != nil {
		errText = b.err.Error()
		if errText == "" {
			errText = "empty error"
		}
	}
	return struct {
		ServiceName  string
		HandlerName  string
		RequestJSON  json.RawMessage
		ResponseJSON json.RawMessage
		Error        string
	}{
		ServiceName:  b.serviceName,
		HandlerName:  b.handlerName,
		RequestJSON:  b.requestJSON,
		ResponseJSON: b.responseJSON,
		Error:        errText,
	}
}

func marshalJSON(v interface{}) json.RawMessage {
	rawJSON, err := json.Marshal(v)
	if err != nil {
		rawJSON, _ = json.Marshal(map[string]string{
			"error": "JSON marshal error: " + err.Error(),
		})
		return rawJSON
	}
	return rawJSON
}

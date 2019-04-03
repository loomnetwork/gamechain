package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/pkg/errors"
	amino "github.com/tendermint/go-amino"
	rpcclient "github.com/tendermint/tendermint/rpc/lib/client"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
)

// We need to have our own JSONRPCClient because the deault JSONRPCClient do not allow to have path in url.
// In our case, the url, for example, is http://plasma-beta.dappchains.com:80/query. The deafult JSONRPCClient
// will turn it into http://plasma-beta.dappchains.com:80.query, which is wrong.

// JSONRPCClient takes params as a slice
type JSONRPCClient struct {
	url    string
	client *http.Client
	cdc    *amino.Codec
}

// make sure it consistent with tenderint http client
var _ rpcclient.HTTPClient = &JSONRPCClient{}

// NewJSONRPCClient returns a JSONRPCClient pointed at the given address.
func NewJSONRPCClient(url string) *JSONRPCClient {
	return &JSONRPCClient{
		url:    url,
		client: http.DefaultClient,
		cdc:    amino.NewCodec(),
	}
}

func (c *JSONRPCClient) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	request, err := rpctypes.MapToRequest(c.cdc, rpctypes.JSONRPCStringID("jsonrpc-client"), method, params)
	if err != nil {
		return nil, err
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	// log.Info(string(requestBytes))
	requestBuf := bytes.NewBuffer(requestBytes)
	// log.Info(Fmt("RPC request to %v (%v): %v", c.remote, method, string(requestBytes)))

	httpResponse, err := c.client.Post(c.url, "text/json", requestBuf)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close() // nolint: errcheck

	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalResponseBytes(c.cdc, responseBytes, result)
}

func (c *JSONRPCClient) Codec() *amino.Codec {
	return c.cdc
}

func (c *JSONRPCClient) SetCodec(cdc *amino.Codec) {
	c.cdc = cdc
}

//------------------------------------------------

func unmarshalResponseBytes(cdc *amino.Codec, responseBytes []byte, result interface{}) (interface{}, error) {
	// Read response.  If rpc/core/types is imported, the result will unmarshal
	// into the correct type.
	// log.Notice("response", "response", string(responseBytes))
	var err error
	response := rpctypes.RPCResponse{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return nil, errors.Errorf("Error unmarshalling rpc response: %v", err)
	}
	if response.Error != nil {
		return nil, errors.Errorf("Response error: %v", response.Error)
	}
	// Unmarshal the RawMessage into the result.
	err = cdc.UnmarshalJSON(response.Result, result)
	if err != nil {
		return nil, errors.Errorf("Error unmarshalling rpc response result: %v", err)
	}
	return result, nil
}

func argsToURLValues(cdc *amino.Codec, args map[string]interface{}) (url.Values, error) {
	values := make(url.Values)
	if len(args) == 0 {
		return values, nil
	}
	err := argsToJSON(cdc, args)
	if err != nil {
		return nil, err
	}
	for key, val := range args {
		values.Set(key, val.(string))
	}
	return values, nil
}

func argsToJSON(cdc *amino.Codec, args map[string]interface{}) error {
	for k, v := range args {
		rt := reflect.TypeOf(v)
		isByteSlice := rt.Kind() == reflect.Slice && rt.Elem().Kind() == reflect.Uint8
		if isByteSlice {
			bytes := reflect.ValueOf(v).Bytes()
			args[k] = fmt.Sprintf("0x%X", bytes)
			continue
		}

		data, err := cdc.MarshalJSON(v)
		if err != nil {
			return err
		}
		args[k] = string(data)
	}
	return nil
}

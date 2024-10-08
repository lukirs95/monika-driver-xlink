package xlinkclient

import (
	"context"
	"encoding/json"

	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
)

type subscribeMessage struct {
	System    xlink.Id          `json:"sysid"`
	EnDecoder xlink.EnDecoderId `json:"id"`
}

type SubscribeResponse struct {
	EnDecoderId xlink.EncoderId `json:"id"`
	Data        struct {
		Values struct {
			VideoEnabled bool `json:"v2110NetPriEnabled"`
			AudioEnabled bool `json:"a2110NetPriEnabled"`
		} `json:"values"`
	} `json:"data"`
}

func (client *Client) Subscribe(ctx context.Context, endecoder xlink.EnDecoderId) (*SubscribeResponse, error) {
	if response, err := client.jrpc.SendRequest(ctx, "state.subscribe", subscribeMessage{System: client.systemId, EnDecoder: endecoder}); err != nil {
		return nil, err
	} else {
		subscribeResponse := SubscribeResponse{}
		if err := json.Unmarshal(response, &subscribeResponse); err != nil {
			return nil, err
		}
		return &subscribeResponse, nil
	}
}

func (client *Client) Unsubscribe(ctx context.Context, endecoder xlink.EnDecoderId) error {
	_, err := client.jrpc.SendRequest(ctx, "state.unsubscribe", subscribeMessage{System: client.systemId, EnDecoder: endecoder})
	return err
}

package xlinkclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
)

type configureMessage struct {
	System    xlink.Id          `json:"sysid"`
	EnDecoder xlink.EnDecoderId `json:"id"`
	Values    configureValues   `json:"values"`
}

type configureValues struct {
	VEnabled bool `json:"v2110NetPriEnabled"`
}

type ConfigureResponse struct {
	Response bool `json:"response"`
}

func (client *Client) EnableVideo(ctx context.Context, endecoder xlink.EnDecoderId) error {
	if res, err := client.jrpc.SendRequest(ctx, "config", configureMessage{
		System:    client.systemId,
		EnDecoder: endecoder,
		Values: configureValues{
			VEnabled: true,
		},
	}); err != nil {
		return err
	} else {
		var result ConfigureResponse
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		} else if !result.Response {
			return fmt.Errorf("could not enable video in `%s`", endecoder)
		} else {
			return nil
		}
	}
}

func (client *Client) DisableVideo(ctx context.Context, endecoder xlink.EnDecoderId) error {
	if res, err := client.jrpc.SendRequest(ctx, "config", configureMessage{
		System:    client.systemId,
		EnDecoder: endecoder,
		Values: configureValues{
			VEnabled: false,
		},
	}); err != nil {
		return err
	} else {
		var result ConfigureResponse
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		} else if !result.Response {
			return fmt.Errorf("could not disable video in `%s`", endecoder)
		} else {
			return nil
		}
	}
}

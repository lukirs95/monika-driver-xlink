package xlinkclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type actionMessage struct {
	System    xlink.Id          `json:"sysid"`
	EnDecoder xlink.EnDecoderId `json:"id"`
}

type ActionResponse struct {
	Response bool `json:"response"`
}

func (client *Client) startEnDecoder(ctx context.Context, endecoder xlink.EnDecoderId) error {
	if res, err := client.jrpc.SendRequest(ctx, "start", actionMessage{client.systemId, endecoder}); err != nil {
		return err
	} else {
		var result ActionResponse
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		} else if !result.Response {
			return fmt.Errorf("could not start `%s`", endecoder)
		} else {
			return nil
		}
	}
}

func (client *Client) stopEnDecoder(ctx context.Context, endecoder xlink.EnDecoderId) error {
	if res, err := client.jrpc.SendRequest(ctx, "stop", actionMessage{
		System:    client.systemId,
		EnDecoder: endecoder,
	}); err != nil {
		return err
	} else {
		var result ActionResponse
		if err := json.Unmarshal(res, &result); err != nil {
			return err
		} else if !result.Response {
			return fmt.Errorf("could not stop `%s`", endecoder)
		} else {
			return nil
		}
	}
}

func (client *Client) startModule(ctx context.Context, module types.Module) error {
	var videoIn xlink.EnDecoderId
	var videoOut xlink.EnDecoderId
	for _, iolet := range module.GetIOlets() {
		if before, found := strings.CutSuffix(iolet.GetName(), ":vi"); found {
			videoIn = xlink.EnDecoderId(before)
		}
		if before, found := strings.CutSuffix(iolet.GetName(), ":vo"); found {
			videoOut = xlink.EnDecoderId(before)
		}
	}
	var finalErr error
	if err := client.startEnDecoder(ctx, videoOut); err != nil {
		finalErr = err
	}
	if err := client.startEnDecoder(ctx, videoIn); err != nil {
		finalErr = err
	}
	return finalErr
}

func (client *Client) stopModule(ctx context.Context, module types.Module) error {
	var videoIn xlink.EnDecoderId
	var videoOut xlink.EnDecoderId
	for _, iolet := range module.GetIOlets() {
		if before, found := strings.CutSuffix(iolet.GetName(), ":vi"); found {
			videoIn = xlink.EnDecoderId(before)
		}
		if before, found := strings.CutSuffix(iolet.GetName(), ":vo"); found {
			videoOut = xlink.EnDecoderId(before)
		}
	}
	var finalErr error
	if err := client.stopEnDecoder(ctx, videoIn); err != nil {
		finalErr = err
	}
	if err := client.stopEnDecoder(ctx, videoOut); err != nil {
		finalErr = err
	}
	return finalErr
}

func (client *Client) restartModule(ctx context.Context, module types.Module) error {
	var finalErr error
	if err := client.stopModule(ctx, module); err != nil {
		finalErr = err
	}
	if err := client.startModule(ctx, module); err != nil {
		finalErr = err
	}
	return finalErr
}

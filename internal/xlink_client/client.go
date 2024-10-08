package xlinkclient

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	jsonrpc "github.com/lukirs95/gojsonrpc"
	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type Client struct {
	device   types.Device
	jrpc     *jsonrpc.JsonRPC
	authKey  string
	systemId xlink.Id
	ready    atomic.Bool
}

func NewClient(device types.Device, password string) *Client {
	AuthMessage.Password = password
	jrpc := jsonrpc.NewJsonRPC()
	client := &Client{
		device:  device,
		jrpc:    jrpc,
		authKey: "",
		ready:   atomic.Bool{},
	}
	jrpc.OnDisconnect = client.onDisconnect
	return client
}

type Context string

const ContextID Context = "CTX"

func (client *Client) Connect(ctx context.Context, full jsonrpc.Subscription, fullupdate jsonrpc.Subscription, reconnect chan *Client, errChan chan error) {
	client.jrpc.SubscribeMethod(ctx, "systems.full", full)
	client.jrpc.SubscribeMethod(ctx, "systems.update", fullupdate)
	notifyAuth := make(jsonrpc.Subscription)
	client.jrpc.SubscribeMethod(ctx, "notify.auth", notifyAuth)

	withTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	go client.asyncAuthenticate(withTimeout, notifyAuth)

	endpoint := fmt.Sprintf("ws://%s/jsonrpc", client.device.GetControlIP())
	err := client.jrpc.Listen(ctx, endpoint)

	client.jrpc.UnsubscribeMethod("systems.full")
	client.jrpc.UnsubscribeMethod("systems.update")
	client.jrpc.UnsubscribeMethod("notify.auth")

	if err != nil {
		errChan <- err
		reconnect <- client
	}
}

func (client *Client) Ready() bool {
	return client.ready.Load()
}

func (client *Client) onDisconnect() {
	client.authKey = ""
	client.ready.Store(false)
}

func (client *Client) Device() types.Device {
	return client.device
}

package controller

import (
	"context"
	"encoding/json"
	"log"
	"time"

	jsonrpc "github.com/lukirs95/gojsonrpc"
	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
	xlinkclient "github.com/lukirs95/monika-driver-xlink/internal/xlink_client"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

type Controller struct {
	updater  chan types.Device
	clients  []*xlinkclient.Client
	stateMap xlinkclient.StateMap
}

func NewController(devices []types.Device, updater chan types.Device, xlinkPassword string) *Controller {
	clients := make([]*xlinkclient.Client, 0)
	for _, device := range devices {
		clients = append(clients, xlinkclient.NewClient(device, xlinkPassword))
	}
	return &Controller{
		updater:  updater,
		clients:  clients,
		stateMap: xlinkclient.NewStateMap(),
	}
}

func (controller *Controller) Listen(ctx context.Context, logger *log.Logger) {
	fullMsg := make(jsonrpc.Subscription)
	fullupdateMsg := make(jsonrpc.Subscription)
	reconnectChan := make(chan *xlinkclient.Client)
	errChan := make(chan error)

	for index, client := range controller.clients {
		withValue := context.WithValue(ctx, xlinkclient.ContextID, controller.clients[index])
		logger.Printf("connect to %s", client.Device().GetControlIP())
		go client.Connect(withValue, fullMsg, fullupdateMsg, reconnectChan, errChan)
	}

	for {
		select {
		case full := <-fullMsg:
			if err := controller.parseFullMsg(logger, full); err != nil {
				logger.Println(err)
			}
		case update := <-fullupdateMsg:
			if err := controller.parseUpdateMsg(logger, update); err != nil {
				logger.Println(err)
			}
		case client := <-reconnectChan:
			withValue := context.WithValue(ctx, xlinkclient.ContextID, client)
			logger.Printf("reconnect %s", client.Device().GetControlIP())
			go reconnect(withValue, client, time.Second*10, fullMsg, fullupdateMsg, reconnectChan, errChan)
		case err := <-errChan:
			logger.Println(err)
		}
	}
}

func (controller *Controller) parseFullMsg(logger *log.Logger, msg jsonrpc.Notification) error {
	var fullXLink xlink.XLink
	if err := json.Unmarshal(msg.Params, &fullXLink); err != nil {
		return err
	}

	client := msg.Ctx.Value(xlinkclient.ContextID).(*xlinkclient.Client)
	logger.Printf("Full system received: %s", client.Device().GetControlIP())
	client.ParseSystem(controller.stateMap, &fullXLink)
	xlinkclient.ParseSystemUpdate(controller.stateMap, &fullXLink)
	controller.updater <- client.Device()

	return nil
}

func (controller *Controller) parseUpdateMsg(logger *log.Logger, msg jsonrpc.Notification) error {
	var update xlink.XLink
	if err := json.Unmarshal(msg.Params, &update); err != nil {
		return err
	}

	client := msg.Ctx.Value(xlinkclient.ContextID).(*xlinkclient.Client)
	logger.Printf("System update received: %s", client.Device().GetControlIP())
	xlinkclient.ParseSystemUpdate(controller.stateMap, &update)
	controller.updater <- client.Device()

	return nil
}

func reconnect(ctx context.Context, client *xlinkclient.Client, after time.Duration, full jsonrpc.Subscription, fullupdate jsonrpc.Subscription, reconnect chan *xlinkclient.Client, errChan chan error) {
	time.Sleep(after)
	go client.Connect(ctx, full, fullupdate, reconnect, errChan)
}

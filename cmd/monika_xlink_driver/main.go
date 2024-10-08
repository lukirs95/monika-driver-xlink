package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lukirs95/monika-driver-xlink/internal/controller"
	"github.com/lukirs95/monika-gosdk/pkg/driver"
	"github.com/lukirs95/monika-gosdk/pkg/provider"
	"github.com/lukirs95/monika-gosdk/pkg/types"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	endpoint := os.Getenv("NETBOX_ADDRESS")
	apiKey := os.Getenv("NETBOX_APIKEY")
	deviceType := types.DeviceType(os.Getenv("MONIKA_DEVICETYPE"))
	deviceTypeId, err := strconv.Atoi(os.Getenv("NETBOX_DEVICETYPE_ID"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()

	provider := provider.NewDeviceProviderNetbox(endpoint, apiKey, deviceType, deviceTypeId)
	if err := provider.FetchDevices(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monikadriver, err := driver.NewDriver(provider)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	service := driver.NewService(os.Getenv("MONIKA_ENDPOINT"), monikadriver, log.Default())
	updateChan := make(chan types.Device, 100)
	controller := controller.NewController(provider.GetDevices(), updateChan, os.Getenv("XLINK_PASSWORD"))

	go controller.Listen(ctx, log.Default())

	fmt.Print(service.Listen(ctx, 8090, updateChan))
}

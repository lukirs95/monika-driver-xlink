package xlinkclient

import (
	"context"
	"encoding/json"

	jsonrpc "github.com/lukirs95/gojsonrpc"
	"github.com/lukirs95/monika-driver-xlink/internal/xlink"
)

const (
	METHOD_AUTH jsonrpc.Method = "auth"
)

var AuthMessage struct {
	Auth     bool   `json:"auth"`
	UserId   string `json:"userid"`
	Password string `json:"pass"`
} = struct {
	Auth     bool   `json:"auth"`
	UserId   string `json:"userid"`
	Password string `json:"pass"`
}{Auth: true, UserId: "admin"}

type authResponse struct {
	AuthKey string `json:"authKey"`
}

type authAdvise struct {
	SystemId xlink.Id `json:"sysid"`
}

func (client *Client) asyncAuthenticate(ctx context.Context, adviseChan jsonrpc.Subscription) {
	select {
	case <-ctx.Done():
		return
	case rawAdvise := <-adviseChan:
		advise := authAdvise{}
		if err := json.Unmarshal(rawAdvise.Params, &advise); err != nil {
			return
		}
		client.systemId = advise.SystemId

		response, err := client.jrpc.SendRequest(context.Background(), METHOD_AUTH, AuthMessage)

		if err != nil {
			client.authKey = ""
			return
		}

		var authResponse authResponse
		err = json.Unmarshal(response, &authResponse)
		if err != nil {
			client.authKey = ""
		}
		client.authKey = authResponse.AuthKey
		client.ready.Store(true)
	}
}

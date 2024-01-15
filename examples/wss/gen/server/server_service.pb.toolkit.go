// Code generated by protoc-gen-toolkit. DO NOT EDIT.
// source: github.com/lastbackend/toolkit/examples/wss/apis/server.proto

package serverpb

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	toolkit "github.com/lastbackend/toolkit"
	"github.com/lastbackend/toolkit-plugins/redis"
	"github.com/lastbackend/toolkit/examples/helloworld/gen"
	client "github.com/lastbackend/toolkit/pkg/client"
	runtime "github.com/lastbackend/toolkit/pkg/runtime"
	controller "github.com/lastbackend/toolkit/pkg/runtime/controller"
	tk_http "github.com/lastbackend/toolkit/pkg/server/http"
	errors "github.com/lastbackend/toolkit/pkg/server/http/errors"
	tk_ws "github.com/lastbackend/toolkit/pkg/server/http/websockets"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the toolkit package it is being compiled against and
// suppress "imported and not used" errors
var (
	_ context.Context
	_ emptypb.Empty
	_ http.Handler
	_ errors.Err
	_ io.Reader
	_ json.Marshaler
	_ tk_ws.Client
	_ tk_http.Handler
	_ client.GRPCClient
)

// Definitions

type Redis1Plugin interface {
	redis.Plugin
}

// Service Router define
type serviceRouter struct {
	runtime runtime.Runtime
}

func NewRouterService(name string, opts ...runtime.Option) (_ toolkit.Service, err error) {
	app := new(serviceRouter)

	app.runtime, err = controller.NewRuntime(context.Background(), name, opts...)
	if err != nil {
		return nil, err
	}

	// loop over plugins and initialize plugin instance
	plugin_redis1 := redis.NewPlugin(app.runtime, &redis.Options{Name: "redis1"})

	// loop over plugins and register plugin in toolkit
	app.runtime.Plugin().Provide(func() Redis1Plugin { return plugin_redis1 })

	// create new Router HTTP server
	app.runtime.Server().HTTPNew(name, nil)

	app.runtime.Server().HTTP().UseMiddleware("example")
	app.runtime.Server().HTTP().AddHandler(http.MethodGet, "/events", app.runtime.Server().HTTP().ServerWS)
	app.runtime.Server().HTTP().Subscribe("SayHello", app.handlerWSProxyRouterSayHello)

	return app.runtime.Service(), nil
}

// Define HTTP handlers for Router HTTP server

func (s *serviceRouter) handlerWSProxyRouterSayHello(ctx context.Context, event tk_ws.Event, c *tk_ws.Client) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var protoRequest servicepb.HelloRequest
	var protoResponse servicepb.HelloReply

	if err := json.Unmarshal(event.Payload, &protoRequest); err != nil {
		return err
	}

	callOpts := make([]client.GRPCCallOption, 0)

	if headers := ctx.Value(tk_ws.RequestHeaders); headers != nil {
		if v, ok := headers.(map[string]string); ok {
			callOpts = append(callOpts, client.GRPCOptionHeaders(v))
		}
	}

	if err := s.runtime.Client().GRPC().Call(ctx, "helloworld", "/helloworld.Greeter/SayHello", &protoRequest, &protoResponse, callOpts...); err != nil {
		return err
	}

	return c.WriteJSON(protoResponse)
}

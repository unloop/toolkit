// Code generated by protoc-gen-toolkit. DO NOT EDIT.
// source: github.com/lastbackend/toolkit/examples/gateway/apis/server.proto

package serverpb

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	toolkit "github.com/lastbackend/toolkit"
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

// Service ProxyGateway define
type serviceProxyGateway struct {
	runtime runtime.Runtime
}

func NewProxyGatewayService(name string, opts ...runtime.Option) (_ toolkit.Service, err error) {
	app := new(serviceProxyGateway)

	app.runtime, err = controller.NewRuntime(context.Background(), name, opts...)
	if err != nil {
		return nil, err
	}

	// loop over plugins and initialize plugin instance

	// loop over plugins and register plugin in toolkit

	// create new ProxyGateway GRPC server
	app.runtime.Server().GRPCNew(name, nil)

	// set descriptor to ProxyGateway GRPC server
	app.runtime.Server().GRPC().SetDescriptor(ProxyGateway_ServiceDesc)
	app.runtime.Server().GRPC().SetConstructor(registerProxyGatewayGRPCServer)

	// create new ProxyGateway HTTP server
	app.runtime.Server().HTTPNew(name, nil)

	app.runtime.Server().HTTP().AddHandler(http.MethodPost, "/hello", app.handlerHTTPProxyGatewayHelloWorld)

	return app.runtime.Service(), nil
}

// Define GRPC services for ProxyGateway GRPC server

type ProxyGatewayRpcServer interface {
	HelloWorld(ctx context.Context, req *servicepb.HelloRequest) (*servicepb.HelloReply, error)
}

type proxygatewayGrpcRpcServer struct {
	ProxyGatewayRpcServer
}

func (h *proxygatewayGrpcRpcServer) HelloWorld(ctx context.Context, req *servicepb.HelloRequest) (*servicepb.HelloReply, error) {
	return h.ProxyGatewayRpcServer.HelloWorld(ctx, req)
}

func (proxygatewayGrpcRpcServer) mustEmbedUnimplementedProxyGatewayServer() {}

func registerProxyGatewayGRPCServer(runtime runtime.Runtime, srv ProxyGatewayRpcServer) error {
	runtime.Server().GRPC().RegisterService(&proxygatewayGrpcRpcServer{srv})
	return nil
}

// Define services for ProxyGateway HTTP server

type ProxyGatewayHTTPService interface {
	HelloWorld(ctx context.Context, req *servicepb.HelloRequest) (*servicepb.HelloReply, error)
}

// Define HTTP handlers for Router HTTP server

func (s *serviceProxyGateway) handlerHTTPProxyGatewayHelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var protoRequest servicepb.HelloRequest
	var protoResponse servicepb.HelloReply

	im, om := tk_http.GetMarshaler(s.runtime.Server().HTTP(), r)

	reader, err := tk_http.NewReader(r.Body)
	if err != nil {
		errors.HTTP.InternalServerError(w)
		return
	}

	if err := im.NewDecoder(reader).Decode(&protoRequest); err != nil && err != io.EOF {
		errors.HTTP.InternalServerError(w)
		return
	}

	headers, err := tk_http.PrepareHeaderFromRequest(r)
	if err != nil {
		errors.HTTP.InternalServerError(w)
		return
	}

	callOpts := make([]client.GRPCCallOption, 0)
	callOpts = append(callOpts, client.GRPCOptionHeaders(headers))

	if err := s.runtime.Client().GRPC().Call(ctx, "helloworld", "/helloworld.Greeter/SayHello", &protoRequest, &protoResponse, callOpts...); err != nil {
		errors.GrpcErrorHandlerFunc(w, err)
		return
	}

	buf, err := om.Marshal(protoResponse)
	if err != nil {
		errors.HTTP.InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", om.ContentType())
	if proceed, err := tk_http.HandleGRPCResponse(w, r, headers); err != nil || !proceed {
		return
	}

	if _, err = w.Write(buf); err != nil {
		return
	}
}

// Code generated by protoc-gen-toolkit. DO NOT EDIT.
// source: github.com/lastbackend/toolkit/examples/http/apis/server.proto

package serverpb

import (
	context "context"

	"github.com/lastbackend/toolkit/examples/http/gen/server"
	client "github.com/lastbackend/toolkit/pkg/client"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Suppress "imported and not used" errors
var _ context.Context
var _ emptypb.Empty

// Client gRPC API for Http service
func NewHttpRPCClient(service string, c client.GRPCClient) HttpRPCClient {
	return &httpGrpcRPCClient{service, c}
}

// Client gRPC API for Http service
type HttpRPCClient interface {
	HelloWorld(ctx context.Context, req *serverpb.HelloRequest, opts ...client.GRPCCallOption) (*serverpb.HelloResponse, error)
}

type httpGrpcRPCClient struct {
	service string
	cli     client.GRPCClient
}

func (c *httpGrpcRPCClient) HelloWorld(ctx context.Context, req *serverpb.HelloRequest, opts ...client.GRPCCallOption) (*serverpb.HelloResponse, error) {
	resp := new(serverpb.HelloResponse)
	if err := c.cli.Call(ctx, c.service, Http_HelloWorldMethod, req, resp, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

func (httpGrpcRPCClient) mustEmbedUnimplementedHttpClient() {}

// Client methods for Http service
const (
	Http_HelloWorldMethod = "/http.Http/HelloWorld"
)
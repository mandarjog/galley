// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package server provides HTTP open service galley API server bindings.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	galleypb "istio.io/galley/api/galley/v1"
	"istio.io/galley/pkg/store"
	"istio.io/galley/pkg/store/inventory"
)

// TODO: allow customization
const maxMessageSize uint = 1024 * 1024

// Server data
type Server struct {
	c *GalleyService
	// TODO: watcher should be added.
}

// CreateServer creates a galley server.
func CreateServer(url string) (*Server, error) {
	kvs, err := store.NewRegistrar(inventory.NewInventory()).NewStore(url)
	if err != nil {
		return nil, err
	}
	c, err := NewGalleyService(kvs)
	if err != nil {
		return nil, err
	}
	return &Server{c}, nil
}

func (s *Server) startGateway(grpcPort, restPort uint16) error {
	ctx := context.Background()

	jsonpb := &runtime.JSONPb{}
	ctypes := []string{"text/plain", "text/json", "text/yaml", "application/octet-stream"}
	mopts := make([]runtime.ServeMuxOption, 0, len(ctypes))
	for _, ctype := range ctypes {
		mopts = append(mopts, runtime.WithMarshalerOption(ctype, &RawMarshaler{jsonpb, ctype}))
	}
	mux := runtime.NewServeMux(mopts...)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithMaxMsgSize(int(maxMessageSize)),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
	}
	err := galleypb.RegisterGalleyHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcPort), opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", restPort), mux)
}

// Start runs the server and listen on port.
// TODO(https://github.com/istio/galley/issues/16)
func (s *Server) Start(grpcPort, restPort uint16) error {
	// get the network stuff setup
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("unable to listen on socket: %v", err)
	}

	if restPort != 0 {
		go func() {
			if err := s.startGateway(grpcPort, restPort); err != nil {
				glog.Errorf("failed to start up the gateway: %v", err)
			}
		}()
	}

	grpcOptions := []grpc.ServerOption{
		grpc.MaxMsgSize(int(maxMessageSize)),
		grpc.RPCCompressor(grpc.NewGZIPCompressor()),
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
	}

	// TODO: cert

	// TODO: tracing
	// if enableTracing {
	// 	tracer := bt.New(tracing.IORecorder(os.Stdout))
	// 	ot.InitGlobalTracer(tracer)
	// 	grpcOptions = append(grpcOptions, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	// }
	gs := grpc.NewServer(grpcOptions...)
	galleypb.RegisterGalleyServer(gs, s.c)

	return gs.Serve(listener)
}

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

// Server data
type Server struct {
}

// CreateServer creates a galley server.
func CreateServer() (*Server, error) {
	return &Server{}, nil
}

// Start runs the server and listen on port.
// TODO(https://github.com/istio/galley/issues/16)
func (s *Server) Start(port uint16) {
	select {}
}

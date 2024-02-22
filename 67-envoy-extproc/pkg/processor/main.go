// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	extproc "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var logger *zap.SugaredLogger = nil

func main() {
	var help bool
	var debug bool
	var port int
	var err error
	var sock string

	flag.IntVar(&port, "p", -1, "Listen port; only one of -p and -s may be specified")
	flag.BoolVar(&debug, "d", false, "Enable debug logging")
	flag.BoolVar(&help, "h", false, "Print help")
	flag.StringVar(&sock, "s", "", "Listen socket; only one of -p and -s may be specified")
	flag.Parse()
	if !flag.Parsed() || help || (port < 0 && sock == "") || (port > 0 && sock != "") {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var zapLogger *zap.Logger
	if debug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Sprintf("Can't initialize logger: %s", err))
	}
	logger = zapLogger.Sugar()

	var listener net.Listener
	if port > 0 {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	} else if sock != "" {
		listener, err = net.Listen("unix", sock)
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt)
		go func() {
			select {
			case <-sigChan:
				fmt.Printf("\nreceived interrupt; removing %s\n", sock)
				if err := os.Remove(sock); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				os.Exit(0)
			}
		}()
	}
	if err != nil {
		logger.Fatalf("Can't listen on socket: %s", err)
		os.Exit(3)
	}

	server := grpc.NewServer()
	service := processorService{}
	extproc.RegisterExternalProcessorServer(server, &service)

	logger.Infof("Listening on %s", listener.Addr())

	server.Serve(listener)
}

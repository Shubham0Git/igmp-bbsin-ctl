/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"github.com/opencord/bbsim/api/bbsim"
	"github.com/opencord/bbsim/internal/bbsim/api"
	"github.com/opencord/bbsim/internal/bbsim/devices"
	bbsimLogger "github.com/opencord/bbsim/internal/bbsim/logger"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"runtime/pprof"
	"sync"
)

type CliOptions struct {
	OltID        int
	NumNniPerOlt int
	NumPonPerOlt int
	NumOnuPerPon int
	STag         int
	CTagInit     int
	profileCpu   *string
	logLevel     string
	logCaller    bool
}

func getOpts() *CliOptions {

	olt_id := flag.Int("olt_id", 0, "Number of OLT devices to be emulated (default is 1)")
	nni := flag.Int("nni", 1, "Number of NNI ports per OLT device to be emulated (default is 1)")
	pon := flag.Int("pon", 1, "Number of PON ports per OLT device to be emulated (default is 1)")
	onu := flag.Int("onu", 1, "Number of ONU devices per PON port to be emulated (default is 1)")

	s_tag := flag.Int("s_tag", 900, "S-Tag value (default is 900)")
	c_tag_init := flag.Int("c_tag", 900, "C-Tag starting value (default is 900), each ONU will get a sequentail one (targeting 1024 ONUs per BBSim instance the range is bug enough)")

	profileCpu := flag.String("cpuprofile", "", "write cpu profile to file")

	logLevel := flag.String("logLevel", "debug", "Set the log level (trace, debug, info, warn, error)")
	logCaller := flag.Bool("logCaller", false, "Whether to print the caller filename or not")

	flag.Parse()

	o := new(CliOptions)

	o.OltID = int(*olt_id)
	o.NumNniPerOlt = int(*nni)
	o.NumPonPerOlt = int(*pon)
	o.NumOnuPerPon = int(*onu)
	o.STag = int(*s_tag)
	o.CTagInit = int(*c_tag_init)
	o.profileCpu = profileCpu
	o.logLevel = *logLevel
	o.logCaller = *logCaller

	return o
}

func startApiServer(channel chan bool, group *sync.WaitGroup) {
	// TODO make configurable
	address := "0.0.0.0:50070"
	log.Debugf("APIServer Listening on: %v", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("APIServer failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bbsim.RegisterBBSimServer(grpcServer, api.BBSimServer{})

	reflection.Register(grpcServer)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go grpcServer.Serve(lis)

	for {
		_, ok := <-channel
		if !ok {
			// if the olt channel is closed, stop the gRPC server
			log.Warnf("Stopping API gRPC server")
			grpcServer.Stop()
			wg.Done()
			break
		}
	}

	wg.Wait()
	group.Done()
	return
}

func main() {
	options := getOpts()

	bbsimLogger.SetLogLevel(log.StandardLogger(), options.logLevel, options.logCaller)

	if *options.profileCpu != "" {
		// start profiling
		log.Infof("Creating profile file at: %s", *options.profileCpu)
		f, err := os.Create(*options.profileCpu)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}

	log.WithFields(log.Fields{
		"OltID":        options.OltID,
		"NumNniPerOlt": options.NumNniPerOlt,
		"NumPonPerOlt": options.NumPonPerOlt,
		"NumOnuPerPon": options.NumOnuPerPon,
	}).Info("BroadBand Simulator is on")

	// control channels, they are only closed when the goroutine needs to be terminated
	oltDoneChannel := make(chan bool)
	apiDoneChannel := make(chan bool)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go devices.CreateOLT(options.OltID, options.NumNniPerOlt, options.NumPonPerOlt, options.NumOnuPerPon, options.STag, options.CTagInit, &oltDoneChannel, &apiDoneChannel, &wg)
	log.Debugf("Created OLT with id: %d", options.OltID)
	go startApiServer(apiDoneChannel, &wg)
	log.Debugf("Started APIService")

	wg.Wait()

	defer func() {
		log.Info("BroadBand Simulator is off")
		if *options.profileCpu != "" {
			log.Info("Stopping profiler")
			pprof.StopCPUProfile()
		}
	}()
}

// Copyright 2017 Ryszard Szopa <ryszard.szopa@gmail.com>
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

// sds011 is a simple reader for the SDS011 Air Quality Sensor. It
// outputs data in TSV to standard output (timestamp formatted
// according to RFC3339, PM2.5 levels, PM10 levels).
package sds011ctl

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/elipavlov/sds011"
)

var (
	rootCmd = &cobra.Command{
		Use:   "sds011",
		Short: "sds011 reads data from the SDS011 sensor and sends them to stdout as CSV",
		Long: `sds011 reads data from the SDS011 sensor and sends them to stdout as CSV.

The columns are: an RFC3339 timestamp, the PM2.5 level, the PM10 level.`,

		SilenceUsage: true,
		RunE:         rootCmdRunE,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := flag.CommandLine.Parse([]string{}); err != nil {
				panic(err)
			}
		},
	}

	portPath string
)

func rootCmdRunE(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var active bool
	var sensor *sds011.Sensor
	for i := 0; i < 10; i++ {
		log.V(6).Infof("connect to sensor: %s, attempt: %d ...", portPath, i)
		sensor, err = sds011.New(portPath)
		if err != nil {
			log.Fatal(err)
		}

		active, err = sensor.ReportMode()
		if err == nil {
			break
		}
		if errors.Is(err, sds011.ErrMalformedRead) {
			sensor.Close()
			continue
		}
	}
	defer sensor.Close()
	if err != nil {
		return err
	}

	if !active {
		interval, err := cmd.Flags().GetUint8("interval")
		if err != nil {
			return err
		}

		go func() {
			readWithTimer(ctx, sensor, interval)
		}()
	} else {
		go func() {
			readFromStream(ctx, sensor)
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.V(6).Infof("sig: %s", sig)

	return err
}

func readFromStream(ctx context.Context, sensor *sds011.Sensor) {
	log.V(6).Info("readFromStream...")
	pointCh := make(chan *sds011.Point, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			point, err := sensor.Get()
			if err != nil {
				log.Warningf("ERROR: sensor.Get: %v", err)
				continue
			}

			select {
			case <-ctx.Done():
				return
			case pointCh <- point:
			default:
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			close(pointCh)
			return
		case point := <-pointCh:
			printPoint(point)
		}
	}
}

func readWithTimer(ctx context.Context, sensor *sds011.Sensor, interval uint8) {
	log.V(6).Info("readWithTimer...")
	timer := time.NewTicker(time.Millisecond * 10)

	var once sync.Once
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			once.Do(func() { timer.Reset(time.Second * time.Duration(interval)) })
			point, err := sensor.Query()
			if err != nil {
				log.Warningf("ERROR: sensor.Get: %v", err)
				continue
			}
			printPoint(point)
		}
	}
}

func printPoint(point *sds011.Point) {
	fmt.Fprintf(os.Stdout, "%v,%v,%v\n", point.Timestamp.Format(time.RFC3339), point.PM25, point.PM10)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint8P("interval", "i", 10,
		"sets interval in seconds between queries (work only with the passive mode)")

	rootCmd.PersistentFlags().StringVarP(&portPath, "port_path", "p",
		"/dev/ttyUSB0", "serial port path")

	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

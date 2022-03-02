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
package sds011ctl

import (
	"fmt"
	"strconv"

	log "github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/elipavlov/sds011"
)

var (
	cycleCmd = &cobra.Command{
		Use:   "cycle <[0,30]>",
		Short: "cycle read or set value of the Cycle",
		Long: `cycle read or set value of the Cycle

		Cycle mode is a mode when sensor will made measurement once per 
		provided amount of the minutes.

		0 - for constant measurement (approximately once per second)
		1 - 30 - for one measurement per N minute, where N is chosen value

		Don't pass any value to read current value recorded to the sensor.
		`,
		Args: cobra.MaximumNArgs(1),

		RunE:         cycleCmdRunE,
		SilenceUsage: true,
	}
)

func cycleCmdRunE(cmd *cobra.Command, args []string) error {
	sensor, err := sds011.New(portPath)
	if err != nil {
		return err
	}
	defer sensor.Close()

	switch len(args) {
	case 1:
		v, err := strconv.ParseUint(args[0], 10, 8)
		if err != nil {
			log.Fatalf("bad unsigned number: %v", args[0])
		}
		if err := sensor.SetCycle(uint8(v)); err != nil {
			return err
		}

	default:
		dutyCycle, err := sensor.Cycle()
		if err != nil {
			return err
		}
		fmt.Println(dutyCycle)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cycleCmd)
}

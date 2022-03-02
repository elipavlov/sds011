// Copyright 2022 Ilya Pavlov <eli.pavlov.vn@gmail.com>
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

	"github.com/spf13/cobra"

	"github.com/elipavlov/sds011"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Info shows information about sensor firware and other",
		Args:  cobra.NoArgs,

		RunE:         infoCmdRunE,
		SilenceUsage: true,
	}
)

func infoCmdRunE(cmd *cobra.Command, args []string) error {
	sensor, err := sds011.New(portPath)
	if err != nil {
		return err
	}
	defer sensor.Close()

	deviceID, err := sensor.DeviceID()
	if err != nil {
		return err
	}
	firmware, err := sensor.Firmware()
	if err != nil {
		return err
	}
	awake, err := sensor.IsAwake()
	if err != nil {
		return err
	}
	active, err := sensor.ReportMode()
	if err != nil {
		return err
	}
	cycle, err := sensor.Cycle()
	if err != nil {
		return err
	}

	mode := modeActive
	if !active {
		mode = modePassive
	}

	fmt.Printf(`sds011 sensor information:
device id: %s
firmware:  %s
is awake:  %t
mode:      %s
cycle:     %d
`,
		deviceID, firmware, awake, mode, cycle)

	return nil
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

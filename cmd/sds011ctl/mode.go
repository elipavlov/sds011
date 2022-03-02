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
	"strings"

	"github.com/spf13/cobra"

	"github.com/elipavlov/sds011"
)

var (
	modeCmd = &cobra.Command{
		Use:   "mode <[active,passive]>",
		Short: "mode reads or sets value of the current activity mode",
		Long: `mode reads or sets value of the current activity mode

active  - means sensor will measure constantly regardless whether someone is reading it or not
passive - means sensor will be asleep till someone query the data
		`,
		Args: cobra.MaximumNArgs(1),

		RunE:         modeCmdRunE,
		SilenceUsage: true,
	}
)

const (
	modeActive  = "active"
	modePassive = "passive"
)

func modeCmdRunE(cmd *cobra.Command, args []string) error {
	sensor, err := sds011.New(portPath)
	if err != nil {
		return err
	}
	defer sensor.Close()

	switch len(args) {
	case 1:
		mode := strings.TrimSpace(strings.ToLower(args[0]))
		switch mode {
		case modeActive:
			err = sensor.MakeActive()
		case modePassive:
			err = sensor.MakePassive()
		default:
			err = fmt.Errorf("malformed mode value provided: %q", mode)
		}
		return err

	default:
		mode, err := sensor.ReportMode()
		if err != nil {
			return err
		}
		modeString := modeActive
		if !mode {
			modeString = modePassive
		}
		fmt.Printf("sds011 mode:\n%s\n", modeString)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(modeCmd)
}

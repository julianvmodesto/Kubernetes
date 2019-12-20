/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package genericclioptions

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/klog"
)

type DryRunStrategy int

const (
	DryRunNone DryRunStrategy = iota
	DryRunClient
	DryRunServer
)

func (d DryRunStrategy) Client() bool {
	return d == DryRunClient
}

func (d DryRunStrategy) Server() bool {
	return d == DryRunServer
}

func (d DryRunStrategy) None() bool {
	return d == DryRunNone
}

// DryRunFlags contains all flags associated with the "--dry-run" operation
type DryRunFlags struct {
	// DryRunStrategy indicates the state of the dry-run flag
	DryRunStrategy DryRunStrategy
	DryRunVerifier *resource.DryRunVerifier
}

// NewDryRunFlags provides a DryRunFlags with reasonable default values set for use
func NewDryRunFlags() *DryRunFlags {
	return &DryRunFlags{
		DryRunStrategy: DryRunNone,
	}
}

// Complete is called before the command is run, but after it is invoked to finish the state of the struct before use.
func (f *DryRunFlags) Complete(cmd *cobra.Command) error {
	var err error
	f.DryRunStrategy, err = getDryRunFlag(cmd)
	if err != nil {
		return err
	}
	return nil
}

// AddFlags binds the requested flags to the provided flagset
func (f *DryRunFlags) AddFlags(cmd *cobra.Command, printFlags *PrintFlags) {
	cmd.Flags().String(
		"dry-run",
		"none",
		`Must be "none", "server", or "client". If client strategy, only print the object that would be sent, without sending it. If server strategy, submit server-side request without persisting the resource.`,
	)
}

func (f *DryRunFlags) WithVerifier(dryRunVerifier *resource.DryRunVerifier) {
	f.DryRunVerifier = dryRunVerifier
}

func (f *DryRunFlags) GetStrategy() DryRunStrategy {
	return f.DryRunStrategy
}

func (f *DryRunFlags) GetVerifier() *resource.DryRunVerifier {
	return f.DryRunVerifier
}

func getDryRunFlag(cmd *cobra.Command) (DryRunStrategy, error) {
	var dryRunFlag, err = cmd.Flags().GetString("dry-run")
	if err != nil {
		return DryRunNone, fmt.Errorf("Error accessing --dry-run")
	}
	if dryRunFlag == "" && !cmd.Flags().Changed("dry-run") {
		klog.Warning(`The unset value for --dry-run is deprecated and a value will be required in a future version. Must be "none", "server", or "client".`)
		return DryRunClient, nil
	}
	b, err := strconv.ParseBool(dryRunFlag)
	// The flag is not a boolean
	if err != nil {
		switch dryRunFlag {
		case "client":
			return DryRunClient, nil
		case "server":
			return DryRunServer, nil
		case "none":
			return DryRunNone, nil
		default:
			return DryRunNone, fmt.Errorf(`Invalid dry-run value (%v). Must be "none", "server", or "client".`, dryRunFlag)
		}
	}
	// The flag was a boolean, and indicates true, run client-side
	if b {
		klog.Warning(`Boolean values for --dry-run are deprecated and will be removed in a future version. Must be "none", "server", or "client".`)
		return DryRunClient, nil
	}
	return DryRunNone, nil
}

// printFlagsWithDryRunStrategy sets a success message at print time for the dry run strategy
func printFlagsWithDryRunStrategy(printFlags *PrintFlags, dryRunStrategy DryRunStrategy) *PrintFlags {
	switch dryRunStrategy {
	case DryRunClient:
		printFlags.Complete("%s (dry run)")
	case DryRunServer:
		printFlags.Complete("%s (server dry run)")
	}
	return printFlags
}

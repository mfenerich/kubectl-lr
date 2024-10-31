/*
Copyright 2024 Marcel Fenerich.

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

// Package main is the entry point for the kubectl-lr plugin.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/sample-cli-plugin/pkg/cmd"
)

// main initializes and executes the kubectl-lr plugin.
func main() {
    // Check if the executable is running as a kubectl plugin
/*     if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
        fmt.Println("Running as a kubectl plugin.")
    } else {
        fmt.Println("Running as a standalone command.")
    } */

    // Initialize the flag set
    flags := pflag.NewFlagSet("kubectl-lr", pflag.ExitOnError)
    pflag.CommandLine = flags

    // Create the root command for the plugin
    root := cmd.NewCmdLimit(genericiooptions.IOStreams{
        In:     os.Stdin,
        Out:    os.Stdout,
        ErrOut: os.Stderr,
    })

    // Execute the root command and handle any errors gracefully
    if err := root.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error executing kubectl-lr: %v\n", err)
        os.Exit(1)
    }
}

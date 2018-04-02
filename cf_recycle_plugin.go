/**
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

// constants
const (
	CfRecyclePluginHelpText = "Recycle CF Application Instances"
	PluginName              = "cf-recycle-plugin"
)

// CfRecycleCmd - struct to initialize.
type CfRecycleCmd struct{}

//GetMetadata - required method to implement plugin
func (CfRecycleCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: PluginName,
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "recycle",
				HelpText: CfRecyclePluginHelpText,
			},
		},
	}
}

// main - entry point to the plugin
func main() {
	plugin.Start(new(CfRecycleCmd))
}

// Run - required method to implement plugin.
func (cmd *CfRecycleCmd) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "recycle" {
		cmd.RecycleCommand(cliConnection, args)
	}
}

// RecycleCommand - recycles AI's
func (cmd *CfRecycleCmd) RecycleCommand(cliConnection plugin.CliConnection, args []string) (err error) {

	var (
		totalInstances int
	)
	//Get app status from cf cli
	appArgs := append([]string{"app"}, args[1:]...)

	if appStatus, err := cliConnection.CliCommandWithoutTerminalOutput(appArgs...); err == nil {

		for _, v := range appStatus {
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "instances: ") {
				instances := strings.Split(strings.TrimPrefix(v, "instances: "), "/")
				totalInstances, _ = strconv.Atoi(instances[1])
			}
		}

		for i := 0; i < totalInstances; i++ {
			var state = cmd.getInstanceStatus(cliConnection, i, args[1])
			if state == "running" {
				cmd.restartInstance(cliConnection, args, i)
			}
		}
	}
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	return
}

func (cmd *CfRecycleCmd) restartInstance(cliConnection plugin.CliConnection, args []string, i int) {
	var state string
	restartArgs := []string{"restart-app-instance", args[1], strconv.Itoa(i)}
	fmt.Printf("Restarting %s instance: %v\n", args[1], i)
	if _, err := cliConnection.CliCommandWithoutTerminalOutput(restartArgs...); err == nil {
		for state != "running" {
			state = cmd.getInstanceStatus(cliConnection, i, args[1])
			time.Sleep(10 * time.Second)
		}
	}
}

func (cmd *CfRecycleCmd) getInstanceStatus(cliConnection plugin.CliConnection, instance int, appName string) (status string) {

	//Get app status from cf cli
	instanceArgs := []string{"app", appName}
	instanceStatus, _ := cliConnection.CliCommandWithoutTerminalOutput(instanceArgs...)

	for _, v := range instanceStatus {
		v = strings.TrimSpace(v)
		fields := strings.Fields(v)
		if len(fields) > 0 && fields[0] == fmt.Sprintf("#%v", instance) {
			status = strings.Fields(v)[1]
		}
	}
	fmt.Printf("Instance %v Status: %s\n", instance, status)
	return
}

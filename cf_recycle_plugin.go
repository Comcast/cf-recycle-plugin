/**
* Copyright 2020 Comcast Cable Communications Management, LLC
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
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

// constants
const (
	CfRecyclePluginHelpText = "Recycle CF Application Instances"
	PluginName              = "cf-recycle-plugin"

	// Cloudfoundry states
	RUNNING = "running"
	STARTED = "started"
)

// Version build flags passed in at time of build
var (
	Major string
	Minor string
	Patch string
)

// CfRecycleCmd - struct to initialize.
type CfRecycleCmd struct {
	startTime time.Time
	appName   string
}

//GetMetadata - required method to implement plugin
func (CfRecycleCmd) GetMetadata() plugin.PluginMetadata {

	major, _ := strconv.Atoi(Major)
	minor, _ := strconv.Atoi(Minor)
	patch, _ := strconv.Atoi(Patch)

	return plugin.PluginMetadata{
		Name: PluginName,
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: patch,
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
	if args[0] == "recycle" && len(args) == 2 {
		cmd.RecycleCommand(cliConnection, args)
	}
}

// RecycleCommand - recycles AI's
func (cmd *CfRecycleCmd) RecycleCommand(cliConnection plugin.CliConnection, args []string) (err error) {

	var (
		guid     string
		appState string
		name     string
	)

	fmt.Printf("Restarting %s...\n", args[1])
	cmd.startTime = time.Now()

	// Get the app guid from the cliConnection
	apps, err := cliConnection.GetApps()

	if err != nil {
		fmt.Printf("Error getting apps: %s\n", err.Error())
		return err
	}

	for _, app := range apps {
		if app.Name == args[1] {
			guid = app.Guid
			appState = app.State
			name = app.Name
		}
	}

	if guid == "" {
		return fmt.Errorf("Unable to find application %s", args[1])
	}
	if appState != STARTED {
		return fmt.Errorf("Application %s is not running", args[1])
	}

	//Get app status from cf cli
	app, err := cliConnection.GetApp(name)

	if err != nil {
		fmt.Printf("Error getting app status: %s\n", err.Error())
		return err
	}

	for i, instance := range app.Instances {
		if instance.State == RUNNING && instance.Since.Before(cmd.startTime) {
			cmd.restartInstance(cliConnection, app.Name, i)
		}
	}
	return
}

func (cmd *CfRecycleCmd) restartInstance(cliConnection plugin.CliConnection, appName string, i int) {
	var (
		state string
		since time.Time
	)
	restartArgs := []string{"restart-app-instance", appName, strconv.Itoa(i)}
	fmt.Printf("Restarting %s instance: %v\n", appName, i)

	if _, err := cliConnection.CliCommandWithoutTerminalOutput(restartArgs...); err == nil {
		for state != "running" || since.Before(cmd.startTime) {
			time.Sleep(5 * time.Second)
			state, since = cmd.getInstanceStatus(cliConnection, i, appName)
		}
	}
}

func (cmd *CfRecycleCmd) getInstanceStatus(cliConnection plugin.CliConnection, instance int, appGUID string) (status string, since time.Time) {

	//Get app status from cf cli
	if app, err := cliConnection.GetApp(appGUID); err == nil {
		if len(app.Instances) > instance {
			inst := app.Instances[instance]

			status = inst.State
			since = inst.Since
			fmt.Printf("Instance %v Status: %s Since: %v\n", instance, status, since)
		}
	}
	return
}

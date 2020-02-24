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
	"errors"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CfRecyclePlugin", func() {
	Describe(".GetMetaData", func() {
		var pluginMetadata plugin.PluginMetadata
		Context("when calling GetMetaData", func() {
			BeforeEach(func() {
				pluginMetadata = new(CfRecycleCmd).GetMetadata()
			})
			It("should return the correct help text", func() {
				success := false
				for _, v := range pluginMetadata.Commands {
					if v.HelpText == CfRecyclePluginHelpText {
						success = true
					}
				}
				Ω(success).Should(BeTrue())
			})
		})
	})
	Describe(".RecycleCommand", func() {
		var recycleCmd *CfRecycleCmd
		var ctrlAppName string
		var fakeConnection *pluginfakes.FakeCliConnection
		var ctrlArgs []string
		var err error

		BeforeEach(func() {
			ctrlAppName = "myTestApp#1.2.3-abcde"
			ctrlArgs = []string{"recycle", ctrlAppName}
			recycleCmd = &CfRecycleCmd{}
			fakeConnection = &pluginfakes.FakeCliConnection{}
		})
		Context("when called and unable to retrieve a list of applications from cloudfoundry", func() {
			BeforeEach(func() {
				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{}, errors.New("unable to get apps"))
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})
			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
		Context("when called and the application is not in a started state", func() {
			BeforeEach(func() {

				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{
					plugin_models.GetAppsModel{
						Name:  "firstAppName",
						State: "started",
						Guid:  "abcd38586erewrwe",
					},
					plugin_models.GetAppsModel{
						Name:  ctrlAppName,
						State: "stopped",
						Guid:  "dflgkjdlgjdfkg6567575",
					},
				}, nil)
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})
			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
		Context("when called and unable to get the app status", func() {
			BeforeEach(func() {

				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{
					plugin_models.GetAppsModel{
						Name:  "firstAppName",
						State: "started",
						Guid:  "abcd38586erewrwe",
					},
					plugin_models.GetAppsModel{
						Name:  ctrlAppName,
						State: "started",
						Guid:  "dflgkjdlgjdfkg6567575",
					},
				}, nil)
				fakeConnection.GetAppReturns(plugin_models.GetAppModel{}, errors.New("unable to get app status"))
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})
			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
		Context("when called with a valid connection and valid application name", func() {
			BeforeEach(func() {
				then := time.Now().AddDate(0, -1, 0)

				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{
					plugin_models.GetAppsModel{
						Name:  "firstAppName",
						State: "started",
						Guid:  "abcd38586erewrwe",
					},
					plugin_models.GetAppsModel{
						Name:  ctrlAppName,
						State: "started",
						Guid:  "dflgkjdlgjdfkg6567575",
					},
				}, nil)
				fakeConnection.GetAppReturnsOnCall(0, plugin_models.GetAppModel{
					Name: ctrlAppName,
					Guid: "dflgkjdlgjdfkg6567575",

					Instances: []plugin_models.GetApp_AppInstanceFields{
						{
							State: "running",
							Since: then,
						},
						{
							State: "running",
							Since: then,
						},
					},
				}, nil)
				fakeConnection.GetAppReturnsOnCall(1, plugin_models.GetAppModel{
					Name: ctrlAppName,
					Guid: "dflgkjdlgjdfkg6567575",

					Instances: []plugin_models.GetApp_AppInstanceFields{
						{
							State: "running",
							Since: time.Now(),
						},
						{
							State: "running",
							Since: then,
						},
					},
				}, nil)
				fakeConnection.GetAppReturnsOnCall(2, plugin_models.GetAppModel{
					Name: ctrlAppName,
					Guid: "dflgkjdlgjdfkg6567575",

					Instances: []plugin_models.GetApp_AppInstanceFields{
						{
							State: "running",
							Since: time.Now(),
						},
						{
							State: "running",
							Since: time.Now(),
						},
					},
				}, nil)
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})

			It("should recycle the application", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when called with a valid connection and an invalid application name", func() {
			BeforeEach(func() {
				fakeConnection.GetAppReturns(plugin_models.GetAppModel{
					Name: "asasa",
				}, errors.New("Failed to find app"))
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})
			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
	})
})

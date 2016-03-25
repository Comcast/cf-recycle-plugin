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
	"errors"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/fakes"
	"github.com/cloudfoundry/cli/plugin/models"
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
		var fakeConnection *fakes.FakeCliConnection
		var ctrlArgs []string
		var err error

		BeforeEach(func() {
			ctrlAppName = "myTestApp#1.2.3-abcde"
			ctrlArgs = []string{"recycle-app", ctrlAppName}
			recycleCmd = &CfRecycleCmd{}
			fakeConnection = &fakes.FakeCliConnection{}
		})

		Context("when called with a valid connection and valid application name", func() {
			BeforeEach(func() {
				fakeConnection.GetAppReturns(plugin_models.GetAppModel{
					Name: ctrlAppName,
				}, errors.New("Failed to find app"))
				err = recycleCmd.RecycleCommand(fakeConnection, ctrlArgs)
			})

			It("should recycle the application", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
		XContext("when called with a valid connection and an invalid application name", func() {
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

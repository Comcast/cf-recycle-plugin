# cf-recycle-plugin

This Cloudfoundry cli plugin is to allow the recycling of application instances without interruption to the application availability.

The plugin works by restarting individual Application Instances(AI's) waiting for one to fully restart before moving on to the next.

### Prerequisites
1. Golang 1.5+
2. CF CLI 6.13+

### Installation from Source
```sh
git clone git@github.com:comcast/cf-recycle-plugin.git
go get github.com/cloudfoundry/cli
godep go build -o deploy/cf-recycle-plugin
cf install-plugin deploy/cf-recycle-plugin -f
```
### Download
Binaries are available in the releases section.

### Usage
```sh
cf recycle-app <APP NAME>
```

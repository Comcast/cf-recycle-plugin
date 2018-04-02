# cf-recycle-plugin

This Cloudfoundry cli plugin is to allow the recycling of application instances without interruption to the application availability.

The plugin works by restarting individual Application Instances(AI's) waiting for one to fully restart before moving on to the next.

### Prerequisites
The plugin was built and tested using the below version
1. Golang 1.9.2
2. CloudFoundry CLI 6.33.1

### Installation from Source
```sh
git clone git@github.com:comcast/cf-recycle-plugin.git
go get github.com/cloudfoundry/cli
go build -o deploy/cf-recycle-plugin
cf install-plugin deploy/cf-recycle-plugin -f
```
### Download
Binaries are available in the releases section.

### Usage
```sh
cf recycle <APP NAME>
```

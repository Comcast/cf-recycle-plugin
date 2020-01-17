# cf-recycle-plugin

This Cloudfoundry cli plugin is to allow the recycling of application instances without interruption to the application availability.

The plugin works by restarting individual Application Instances(AI's) waiting for one to fully restart before moving on to the next.

### Prerequisites
The plugin was built and tested using the below versions
1. Golang 1.13.5
2. CloudFoundry CLI 6.48.0

### Installation from Source
Using your favorite versioning system, set variables for the major, minor, and patch versions.
```sh
git clone git@github.com:comcast/cf-recycle-plugin.git
go build -ldflags "-X main.Major=${major} -X main.Minor=${minor} -X main.Patch=${patch}" -o out/cf-recycle-plugin cf_recycle_plugin.go
cf install-plugin out/cf-recycle-plugin -f
```
### Download
Binaries are available in the releases section.

### Usage
```sh
cf recycle <APP NAME>
```

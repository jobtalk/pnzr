# thor
ecs deploy docker container

[![CircleCI](https://circleci.com/gh/jobtalk/thor.svg?style=shield)](https://circleci.com/gh/jobtalk/thor)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

## Support Go version
* Go 1.8

## Installation
Can be installed in either way.

### Use install script
On macOS, or Linux run the following:
```
$ curl https://raw.githubusercontent.com/jobtalk/thor/master/install.sh | sh
```

Note that you may need to run the sudo version below, or alternatively chown /usr/local:
```
$ curl https://raw.githubusercontent.com/jobtalk/thor/master/install.sh | sudo sh
```

### Use Go get
```
$ go get -u github.com/jobtalk/thor
```

### Use Homebrew
```
$ brew tap ieee0824/thor
$ brew update
$ brew install thor
```
## Detailed instructions
Please read the [wiki](https://github.com/jobtalk/thor/wiki).


## Deploy
```
$ thor deploy -f config.json
```

## Option
```
$ go run main.go -h

Usage: thor [--version] [--help] <command> [<args>]

Available commands are:
    deploy    usage: thor deploy [options ...]
options:
    -f thor_setting.json

    -profile=${aws profile name}
        -profile option is arbitrary parameter.
    -region
        aws region
    -vars_path
        setting external values path file
    -V
        setting outer values
===================================================

    mkelb     usage: thor mkelb [options ...]
options:
    -f thor_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
===================================================

    vault     usage: thor vault [options ...]
options:
    -key_id
        set kms key id
    -encrypt
        use encrypt mode
    -decrypt
        use decrypt mode
    -file
        setting target file
    -f        setting target file
    -profile
        aws credential name
    -region
        aws region name
===================================================
```

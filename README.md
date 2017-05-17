# pnzr
ecs deploy docker container

[![CircleCI](https://circleci.com/gh/jobtalk/pnzr.svg?style=shield)](https://circleci.com/gh/jobtalk/pnzr)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

## Support Go version
* Go 1.8

## Installation
Can be installed in either way.

### Use install script
On macOS, or Linux run the following:
```
$ curl https://raw.githubusercontent.com/jobtalk/pnzr/master/install.sh | sh
```

Note that you may need to run the sudo version below, or alternatively chown /usr/local:
```
$ curl https://raw.githubusercontent.com/jobtalk/pnzr/master/install.sh | sudo sh
```

### Use Go get
```
$ go get -u github.com/jobtalk/pnzr
```

## Detailed instructions
Please read the [wiki](https://github.com/jobtalk/pnzr/wiki).


## Deploy
```
$ pnzr deploy -f config.json
```

## Option
```
$ go run main.go -h

Usage: pnzr [--version] [--help] <command> [<args>]

Available commands are:
    deploy    usage: pnzr deploy [options ...]
options:
    -f pnzr_setting.json

    -profile=${aws profile name}
        -profile option is arbitrary parameter.
    -region
        aws region
    -vars_path
        setting external values path file
    -V
        setting outer values
===================================================

    mkelb     usage: pnzr mkelb [options ...]
options:
    -f pnzr_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
===================================================

    vault     usage: pnzr vault [options ...]
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

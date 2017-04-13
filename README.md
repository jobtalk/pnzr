# thor
ecs deploy docker container

[![CircleCI](https://circleci.com/gh/jobtalk/thor.svg?style=shield)](https://circleci.com/gh/jobtalk/thor)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)


## install
```
$ go get -u github.com/jobtalk/thor
```

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

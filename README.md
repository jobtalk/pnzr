# thor
ecs deploy docker container

[![CircleCI](https://circleci.com/gh/ieee0824/thor.svg?style=shield)](https://circleci.com/gh/ieee0824/thor)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)


<img src="http://img.ieee.moe/thor.png" width="150px">

## install
```
$ go get -u github.com/ieee0824/thor/cmd/thor
```

## Deploy
```
$ thor deploy -f config.json --profile credential-name
```

## Option
```
Usage: thor [--version] [--help] <command> [<args>]

Available commands are:
    deploy    usage: thor deploy [options ...]
options:
    -f thor_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
    --vault-password-file=${vault pass file}
    --ask-vault-pass=${vault pass string}
===================================================

    mkelb     usage: thor mkelb [options ...]
options:
    -f thor_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
===================================================

    vault     usage: thor vault [options ...]
options:
    -f vault target json

    -p vault password
===================================================
```

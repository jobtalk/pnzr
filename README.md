# thor
ecs deploy docker container

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
$ thor --help

Usage: thor [--version] [--help] <command> [<args>]

Available commands are:
    deploy    usage: thor deploy [options ...]
options:
    -f thor_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
===================================================

    mkelb     usage: thor mkelb [options ...]
options:
    -f thor_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
===================================================
```

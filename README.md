# thor
ecs deploy docker container

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
    deploy    usage: deploy [options ...]
options:
    -f deploy_setting.json

    --profile=${aws profile name}
        --profile option is arbitrary parameter.
```

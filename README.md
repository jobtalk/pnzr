# pnzr
ecs deploy docker container

```
　　　　　　　 　 　 　 　 　 　 　 　 r/j=＝＝＝＝＝､､
　 　 　 　 　 　 　 　 =========〔' / ,i!uij　　 i　7!
　　　　　　　　　　　　　　　　　　_f/_ ＿＿{日}__＿!__i,j
　　　　　　　　　　　　 -r｡' 二〔´/´￣￣｀,___,￣￣｀￢f｀{!､
　　　　　　　　　　　／´￣￣￣￣￣f´〔ロ〕 j!￣｀'''￢─--'' ､j､
　　　　　　　　 　 ,r‐j=＝＝＝!------､ 'ｰ---‐'r============ヽj
　　　　　　　　　 (´tj〕ｿv'v'､ftj!ｿv'!ｿv'v'v'v'､ftj!ｿv'v'v'（◎）
　　　　　　　　　　'､ヾ'´､＿,fjj､＿＿＿,fjj､＿＿,,fjj､＿＿_ｿ´7
　　　　　　　　　　　＼　 ,(◎X◎)　　 ,(◎X◎)　　,(◎X◎)　　'／
　　　　　　　　　　　　 ｀ﾞ"""""""""""""""""""""""""""""ﾞ´
```

[![CircleCI](https://circleci.com/gh/jobtalk/pnzr.svg?style=shield)](https://circleci.com/gh/jobtalk/pnzr)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

The pnzr package works on Go versions:
* 1.8.x and greater
* 1.9.x and greater
* 1.10.x and greater

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

## Update latest version

```
$ pnzr update
```

## Examples

### Deploy ecs

```
$ pnzr deploy -f setting.json
$ pnzr deploy -profile aws/profile -f setting.json
```

### Encrypt setting

```
$ pnzr vault encrypt -f target.json
$ pnzr vault encrypt -key_id ${KMS_KEY_ID} -f target.json
```

### Decrypt setting

```
$ pnzr vault decrypt -f target.json
```

### Viewer mode of encrypted setting file

#### latest config version
```
$ pnzr vault view -f target.json
```

#### choose config version
```
$ pnzr vault view -v prototype -а target.json
```

#### check default config version
```
$pnzr vault view -h
  -v string
    	config version (default "1.0")
```

### Edit mode of encrypted file

```
$ pnzr vault edit -f target.json
```

### Support Multi-Factor Authentication(MFA)

```
$ pnzr vault view -profile use-mfa-user -f deploy-config.json 
Assume Role MFA token code: ******
```

# use any editor in edit mode
It run when assigning editor name to EDITOR.

## vim
```
$ EDITOR=vim pnzr vault edit -f /path/to/target
```

## vscode
```
$ EDITOR="code --wait" pnzr vault edit -f /path/to/target
```

## atom
```
$ EDITOR="atom --wait" pnzr vault edit -f /path/to/target
```

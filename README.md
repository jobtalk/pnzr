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

## Show options

```
$ pnzr -h
```

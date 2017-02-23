# gump

[![travis Status](https://travis-ci.org/mh-cbon/gump.svg?branch=master)](https://travis-ci.org/mh-cbon/gump)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/gump?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/gump)
[![GoDoc](https://godoc.org/github.com/mh-cbon/gump?status.svg)](http://godoc.org/github.com/mh-cbon/gump)


Gump is an utility to bump your package using semver.


This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

Check the [release page](https://github.com/mh-cbon/gump/releases)!

#### Glide

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/gump
cd $GOPATH/src/github.com/mh-cbon/gump
git clone https://github.com/mh-cbon/gump.git .
glide install
go install
```


#### Chocolatey
```sh
choco install gump
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gump sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh \
| GH=mh-cbon/gump sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gump sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/gump sh -xe
```

# Usage

__$ gump -help__
```sh
Gump - Bump your package

Usage:
  gump prerelease [-b|--beta] [-a|--alpha] [-d|--dry] [-m <message>]
  gump patch [-d|--dry] [-m <message>]
  gump minor [-d|--dry] [-m <message>]
  gump major [-d|--dry] [-m <message>]
  gump -h | --help
  gump -v | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -d --dry              Only display the new version number.
  -b --beta             Update last beta version.
  -a --alpha            Update last alpha version.
  -m                    Set tag message.

Examples
  # Bump patch with a message
  gump patch -m "tag message"
  # Bump major with a message
  gump major -m "tag message"
```

# Cli examples

```sh
gump patch -d
gump prerelease -a -d
gump prerelease -b -d
gump minor -d
gump major -d
```

## Pre/Post version scripts

Gump can detect, parse and execute `pre/post` version scripts.

They are numerous hooks executed in this order :

- __prebump__ : Runs in first for any type of update, it does not receive `!newversion!` tag.
It should be used to synchronize your local with remote.
- __prepatch__ : Runs for a `patch` update.
- __preminor__ : Runs for a `minor` update.
- __premajor__ : Runs for a `major` update.
- __preversion__ : Runs for a any type of update.
- __version is set here on your vcs__
- __postversion__ : Runs for a any type of update.
- __postmajor__ : Runs for a `major` update.
- __postminor__ : Runs for a `minor` update.
- __postpatch__ : Runs for a `patch` update.
- __postbump__ : Runs for a any type of update.

If any `pre` script returns an exit code different than `0`,
 the execution will stop and the version will remain untouched.

If `post` script fails, the version has already changed,
`gump` will return an exit code = 1.

Scripts can use few special tags
- `!newversion!` will be replaced by the value of the new version (not available in `prebump`)
- `!tagmessage!` will be replaced by the value of the tag message
- `!isprerelease!` will be replaced by `yes` when `--beta` or `--alpha`
arguments are present, otherwise it is replaced by the value `no`
- `!isprerelease_int!` will be replaced by `1` when `--beta` or `--alpha`
arguments are present, otherwise it is replaced by the value `0`
- `!isprerelease_bool!` will be replaced by `true` when `--beta` or `--alpha`
arguments are present, otherwise it is replaced by the value `false`

#### Using a .version.sh file

Drop a file named `.version.sh` on your root such


__> .version.sh__
```sh
PREBUMP=
  666 git fetch --tags origin master
  666 git pull origin master

PREVERSION=
  philea -s "666 go vet %s" "666 go-fmt-fail %s"
  666 go run gump.go -v
  666 changelog finalize --version !newversion!
  666 commit -q -m "changelog: !newversion!" -f change.log

POSTVERSION=
  666 changelog md -o CHANGELOG.md --vars='{"name":"gump"}'
  666 commit -q -m "changelog: !newversion!" -f CHANGELOG.md
  666 git push
  666 git push --tags
  666 gh-api-cli create-release -n release -o mh-cbon -r gump \
   --ver !newversion! -c "changelog ghrelease --version !newversion!" \
  --draft !isprerelease!
  666 go install --ldflags "-X main.VERSION=!newversion!"
```

#### Using a .version file

Drop a file named `.version` on your root such


__> .version-demo__
```version-demo
# demo of .version file
prebump: git fetch --tags
preversion: go vet && go fmt
postversion: git push && git push --tags
```

#### Using a glide.yaml file

Add a new key `scripts` to your glide file such


__> demo-glide.yaml__
```yaml
package: github.com/mh-cbon/gump
scripts:
  prebump: echo "pre bump"
  postbump: echo "post bump"
import:
- package: gopkg.in/yaml.v2
```

If you have longer scripts to run you can write it like this:


__> demo-long-glide.yaml__
```yaml
scripts:
  preversion: |
    echo "hello" \
    && echo "how are you ?"
  postversion: |
    echo "I m fine"
    && echo "goodbye"
```

## debug

Declare a new env `VERBOSE=gump` or `VERBOSE=*` to get more information.

```sh
VERBOSE=gump gump patch -d
VERBOSE=* gump patch -d
```

## todo

- at some point, probably, move to https://github.com/Masterminds/vcs

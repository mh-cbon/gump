# gump - Bump it!

Bin util to bump your package using semver

## Install

```sh
mkdir -p $GOPATH/github.com/mh-cbon
cd $GOPATH/github.com/mh-cbon
git clone https://github.com/mh-cbon/gump.git
cd gump
glide install
go install
```

## Usage

```sh
Gump - Bump your package

Usage:
  gump prerelease [-b|--beta] [-a|--alpha] [-d|--dry]
  gump patch [-d|--dry]
  gump minor [-d|--dry]
  gump major [-d|--dry]
  gump -h | --help
  gump -v | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -d --dry              Only display the new version number.
  -b --beta             Update last beta version.
  -a --alpha            Update last alpha version.
```

## Pre/Post version scripts

Gump can detect, parse and execute pre/post version scripts.

If `preversion` script fails (return code != 0), then the execution will stop and the version will remain untouched.

If `postversion` script fails, the version has already changed, `gump` will return a return code = 1

#### Using a .version file

Drop a file named `.version` on your root such

```sh
# demo of .version file
preversion: git fetch --tags && philea "go vet %s" "go fmt %s"
postversion: git push && git push --tags \
# multiline works too
&& echo "and also with comments in the middle"
```

#### Using a glide.yaml file

Add a new key `scripts` to your glide file such

```yml
scripts:
  preversion: echo "hello"
  postversion: echo "goodbye"
```

## todo

- at some point, probably, move to https://github.com/Masterminds/vcs

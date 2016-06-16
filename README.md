# gump - Bump it!

Bin util to bump your package using semver

## Install

You can grab a pre-built binary file in the [releases page](https://github.com/mh-cbon/gump/releases)

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

## Pre/Post version scripts

Gump can detect, parse and execute pre/post version scripts.

If `preversion` script fails (return code != 0), then the execution will stop and the version will remain untouched.

If `postversion` script fails, the version has already changed, `gump` will return an exit code = 1.

Scripts can use two special tags
- `!newversion!` will be replaced by the value of the new version
- `!tagmessage!` will be replaced by the value of the tag message

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

If you have longer scripts to run you can write it like this:

```yml
scripts:
  preversion: |
    echo "hello" \
    && echo "how are you ?"
  postversion: |
    echo "I m fine"
    && echo "goodbye"
```

## changelog

- 0.0.11: Add YAML multiline scripts support. Updated release scripts.
- 0.0.10: Add `-m` argument, add `!newversion!` and `!tagmessage!` tags.
- up to 0.0.7: various minor improvements.
- 0.0.1: Initial release

## todo

- at some point, probably, move to https://github.com/Masterminds/vcs

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

#### Using a .version file

Drop a file named `.version` on your root such

```sh
# demo of .version file
prebump: git fetch --tags
preversion: philea "go vet %s" "go fmt %s"
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

## debug

Declare a new env `VERBOSE=gump` or `VERBOSE=*` to get more information.

```sh
VERBOSE=gump gump patch -d
VERBOSE=* gump patch -d
```

## todo

- at some point, probably, move to https://github.com/Masterminds/vcs

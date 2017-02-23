# {{.Name}}

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}

{{pkgdoc "gump.go"}}

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

# Install

{{template "gh/releases" .}}

#### Glide
{{template "glide/install" .}}

#### Chocolatey
{{template "choco/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

# Usage
{{cli "gump" "-help"}}

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

{{file ".version.sh"}}

#### Using a .version file

Drop a file named `.version` on your root such

{{file ".version-demo"}}

#### Using a glide.yaml file

Add a new key `scripts` to your glide file such

{{file "demo-glide.yaml"}}

If you have longer scripts to run you can write it like this:

{{file "demo-long-glide.yaml"}}

## debug

Declare a new env `VERBOSE=gump` or `VERBOSE=*` to get more information.

```sh
VERBOSE=gump gump patch -d
VERBOSE=* gump patch -d
```

## todo

- at some point, probably, move to https://github.com/Masterminds/vcs

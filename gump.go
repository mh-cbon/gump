// Gump is an utility to bump your package using semver.

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt.go"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/gump/config"
	"github.com/mh-cbon/gump/gump"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// VERSION contains the build version number.
var VERSION = "0.0.0"

var bumpPatch = "patch"
var bumpMinor = "minor"
var bumpMajor = "major"
var bumpPrerelease = "prerelease"
var bumps = []string{bumpPatch, bumpMinor, bumpMajor}

func main() {

	usage := `Gump - Bump your package

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
`

	arguments, err := docopt.Parse(usage, nil, true, "Gump - "+VERSION, false)

	logger.Println(arguments)
	exitWithError(err)

	isDry := isDry(arguments)
	logger.Println("isDry=", isDry)
	message := getMessage(arguments)
	logger.Println("message=" + message)

	path, err := os.Getwd()
	logger.Println("path=" + path)
	exitWithError(err)

	vcs, err := repoutils.WhichVcs(path)
	exitWithError(err)
	if vcs == "svn" {
		exitWithError(errors.New("Sorry ! Subversion is not supported !"))
	}

	finder := config.NewFinder()

	conf, err := finder.Load(path)
	if conf != nil {
		exitWithError(err)
	}

	cmd := getCommand(arguments)
	logger.Println("cmd=" + cmd)
	if cmd != bumpPrerelease && cmd != bumpPatch && cmd != bumpMinor && cmd != bumpMajor {
		fmt.Println("Wrong usage: Missing command")
		fmt.Println("")
		fmt.Println(usage)
		os.Exit(1)
	}

	isPreRelease := isBeta(arguments) || isAlpha(arguments)

	// prebump: sync the repo
	executeScript("prebump", conf, path, "", isPreRelease, message, isDry)

	// determine the next version
	newVersion, err := gump.DetermineTheNewTag(path, cmd, isBeta(arguments), isAlpha(arguments))
	logger.Println("newVersion=" + newVersion)
	exitWithError(err)

	// execute one of pre{patch,minor,major}
	for _, bump := range bumps {
		if bump == cmd {
			executeScript("pre"+bump, conf, path, newVersion, isPreRelease, message, isDry)
		}
	}
	// execute preversion, always
	executeScript("preversion", conf, path, newVersion, isPreRelease, message, isDry)

	// create the new tag
	if isDry {
		fmt.Println("The new tag to create is: " + newVersion + "\n")
	} else {
		out, err := applyVersionUpgrade(vcs, path, newVersion, message)
		fmt.Println(out)
		exitWithError(err)
		fmt.Println("Created new tag " + newVersion)
	}

	// execute postversion, always
	executeScript("postversion", conf, path, newVersion, isPreRelease, message, isDry)

	// execute one of post{patch,minor,major}
	for _, bump := range bumps {
		if bump == cmd {
			executeScript("post"+bump, conf, path, newVersion, isPreRelease, message, isDry)
		}
	}

	// execute postbump, always
	executeScript("postbump", conf, path, newVersion, isPreRelease, message, isDry)
}

// Check the vcs is clean, then create new tag according to the given version number
func applyVersionUpgrade(vcs string, path string, newVersion string, message string) (string, error) {
	ok, err := repoutils.IsClean(vcs, path)
	if ok == false {
		return "", errors.New(`Your local copy contains uncommited changes!

gump will not proceed further until you committed all your changes.

Once the repository has a clean state, please run gump again.
`)
	}
	if err != nil {
		return "", err
	}
	if len(message) == 0 {
		message = "tag: " + newVersion
	}
	ok, out, err := repoutils.CreateTag(vcs, path, newVersion, message)
	logger.Printf("ok=%t\n", ok)
	if err == nil && ok != true {
		err = errors.New("Something gone wrong!")
	}
	return out, err
}

// executes the preversion script of given config if it is not empty
func executeScript(which string, conf config.Configured, path string, newVersion string, isPreRelease bool, message string, dry bool) {
	if conf != nil {
		script, _ := conf.GetScript(which)
		if script != "" {
			script = strings.Replace(script, "!newversion!", newVersion, -1)
			script = strings.Replace(script, "!tagmessage!", message, -1)
			if isPreRelease {
				script = strings.Replace(script, "!isprerelease!", "yes", -1)
				script = strings.Replace(script, "!isprerelease_int!", "1", -1)
				script = strings.Replace(script, "!isprerelease_bool!", "true", -1)
			} else {
				script = strings.Replace(script, "!isprerelease!", "no", -1)
				script = strings.Replace(script, "!isprerelease_int!", "0", -1)
				script = strings.Replace(script, "!isprerelease_bool!", "false", -1)
			}
			if dry {
				fmt.Println(which + ":\n" + script + "\n")
			} else {
				logger.Println(which + ":\n" + script + "\n")
				err := gump.ExecScript(path, script)
				if err != nil {
					fmt.Println("An has error occured while executing " + which + " script!")
				}
				exitWithError(err)
			}
		}
	}
}

// exits current program if error is not nil,
// prints error on stdout
func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// helper to get the next type of desired version
func getCommand(arguments map[string]interface{}) string {
	x := []string{bumpPatch, bumpMinor, bumpMajor, bumpPrerelease}
	for _, v := range x {
		if p, ok := arguments[v]; ok {
			if b, ok := p.(bool); ok && b {
				return v
			}
		}
	}
	return ""
}

// helper to get the value of the message of the command line
func getMessage(arguments map[string]interface{}) string {
	message := ""
	if m, ok := arguments["<message>"].(string); ok {
		message = m
	}
	return message
}

// helper to get the value --dry on the command line
func isDry(arguments map[string]interface{}) bool {
	dry := false
	if isDry, ok := arguments["--dry"].(bool); ok {
		dry = isDry
	} else {
		if isD, ok := arguments["-d"].(bool); ok {
			dry = isD
		}
	}
	return dry
}

// helper to get the value of --beta on the command line
func isBeta(arguments map[string]interface{}) bool {
	beta := false
	if isBeta, ok := arguments["--beta"].(bool); ok {
		beta = isBeta
	} else {
		if isB, ok := arguments["-a"].(bool); ok {
			beta = isB
		}
	}
	return beta
}

// helper to get the value of --alpha on the command line
func isAlpha(arguments map[string]interface{}) bool {
	alpha := false
	if isAlpha, ok := arguments["--alpha"].(bool); ok {
		alpha = isAlpha
	} else {
		if isA, ok := arguments["-a"].(bool); ok {
			alpha = isA
		}
	}
	return alpha
}

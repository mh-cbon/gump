package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt.go"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/gump/config"
	"github.com/mh-cbon/gump/gump"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

var VERSION = "0.0.0"

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

	hasConfig := config.Exists(path)
	conf, err := config.Load(path)
	if hasConfig {
		exitWithError(err)
	}

	cmd := getCommand(arguments)
	logger.Println("cmd=" + cmd)

	isPreRelease := isBeta(arguments) || isAlpha(arguments)

	if cmd == "prerelease" || cmd == "patch" || cmd == "minor" || cmd == "major" {

		if hasConfig {
			executeScript("prebump", conf, path, "", isPreRelease, message, isDry)
		}

		newVersion, err := gump.DetermineTheNewTag(path, cmd, isBeta(arguments), isAlpha(arguments))
		logger.Println("newVersion=" + newVersion)
		exitWithError(err)

		if hasConfig && cmd == "patch" {
			executeScript("prepatch", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig && cmd == "minor" {
			executeScript("preminor", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig && cmd == "major" {
			executeScript("premajor", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig {
			executeScript("preversion", conf, path, newVersion, isPreRelease, message, isDry)
		}

		if isDry {
			fmt.Println("The new tag to create is: " + newVersion)
		} else {
			out, err := applyVersionUpgrade(vcs, path, newVersion, message)
			fmt.Println(out)
			exitWithError(err)
			fmt.Println("Created new tag " + newVersion)
		}

		if hasConfig {
			executeScript("postversion", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig && cmd == "major" {
			executeScript("postmajor", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig && cmd == "minor" {
			executeScript("postminor", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig && cmd == "patch" {
			executeScript("postpatch", conf, path, newVersion, isPreRelease, message, isDry)
		}
		if hasConfig {
			executeScript("postbump", conf, path, newVersion, isPreRelease, message, isDry)
		}

	} else if cmd == "" {
		fmt.Println("Wrong usage: Missing command")
		fmt.Println("")
		fmt.Println(usage)
		os.Exit(1)

	} else {
		log.Println("Unknown command: '" + cmd + "'")
		os.Exit(1)
	}
}

// Check the vcs is clean, then create new tag according to the given version number
func applyVersionUpgrade(vcs string, path string, newVersion string, message string) (string, error) {
	ok, err := repoutils.IsClean(vcs, path)
	if ok == false {
		return "", errors.New("Your local copy contains uncommited changes!")
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
	script := config.GetScript(which, conf)
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
			fmt.Println(which + ":" + script)
		} else {
			logger.Println(which + "=" + script)
			err := gump.ExecScript(path, script)
			if err != nil {
				fmt.Println("An has error occured while executing " + which + " script!")
			}
			exitWithError(err)
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
	p, ok := arguments["patch"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "patch"
		}
	}
	p, ok = arguments["minor"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "minor"
		}
	}
	p, ok = arguments["major"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "major"
		}
	}
	p, ok = arguments["prerelease"]
	if ok {
		if b, ok := p.(bool); ok && b {
			return "prerelease"
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

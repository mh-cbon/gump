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

	arguments, err := docopt.Parse(usage, nil, true, "Gump 0.0.4", false)

	logger.Println(arguments)
	if err != nil {
		exitWithError(err)
	}

	isDry := isDry(arguments)
	logger.Println("isDry=", isDry)
	message := getMessage(arguments)
	logger.Println("message=" + message)

	path, err := os.Getwd()
	logger.Println("path=" + path)
	if err != nil {
		exitWithError(err)
	}

	vcs, err := repoutils.WhichVcs(path)
	if err != nil {
		exitWithError(err)
	}
	if vcs == "svn" {
		exitWithError(errors.New("Sorry ! Subversion is not supported !"))
	}

	hasConfig := config.Exists(path)
	conf, err := config.Load(path)
	if hasConfig && err != nil {
		exitWithError(err)
	}

	cmd := getCommand(arguments)
	logger.Println("cmd=" + cmd)

	if cmd == "prerelease" || cmd == "patch" || cmd == "minor" || cmd == "major" {

		newVersion, err := gump.DetermineTheNewTag(path, cmd, isBeta(arguments), isAlpha(arguments))
		logger.Println("newVersion=" + newVersion)
		if err != nil {
			exitWithError(err)
		}

		if hasConfig {
			script := conf.GetPreVersion()
			if script != "" {
				script = strings.Replace(script, "!newversion!", newVersion, -1)
				script = strings.Replace(script, "!tagmessage!", message, -1)
				if isDry {
					fmt.Println("preversion:" + script)
				} else {
					logger.Println("preversion=" + script)
					out, err := gump.ExecScript(script)
					if err != nil {
						fmt.Println("An has error occured while executing preversion script!")
						fmt.Println("script: " + script)
						fmt.Println(out)
						exitWithError(err)
					}
					fmt.Println(out)
				}
			}
		}

		if isDry {
			fmt.Println("The new tag to create is: " + newVersion)
		} else {
			ok, err := repoutils.IsClean(vcs, path)
			if ok == false {
				exitWithError(errors.New("Your local copy contains uncommited changes!"))
			}
			if err != nil {
				exitWithError(err)
			}
			ok, out, err := repoutils.CreateTag(vcs, path, newVersion, message)
			logger.Printf("ok=%t\n", ok)
			if err != nil {
				fmt.Println(out)
				exitWithError(err)
			}
			if ok != true {
				fmt.Println(out)
				exitWithError(errors.New("Something gone wrong!"))
			}
			fmt.Println("Created new tag " + newVersion)
		}

		if hasConfig {
			script := conf.GetPostVersion()
			if script != "" {
				script = strings.Replace(script, "!newversion!", newVersion, -1)
				script = strings.Replace(script, "!tagmessage!", message, -1)
				if isDry {
					fmt.Println("postversion:" + script)
				} else {
					logger.Println("postversion=" + script)
					out, err := gump.ExecScript(script)
					if err != nil {
						fmt.Println("An has error occured while executing postversion script!")
						fmt.Println("script: " + script)
						fmt.Println(out)
						exitWithError(err)
					}
					fmt.Println(out)
				}
			}
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

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

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

func getMessage(arguments map[string]interface{}) string {
	message := ""
	if mess, ok := arguments["-m"].(string); ok {
		message = mess
	}
	return message
}

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

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
`

	arguments, err := docopt.Parse(usage, nil, true, "Gump 0.0.4", false)

	logger.Println(arguments)
	if err != nil {
		exitWithError(err)
	}

	isDry := isDry(arguments)

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
	ok, err := repoutils.IsClean(vcs, path)
	if ok == false {
		exitWithError(errors.New("Your local copy contains uncommited changes!"))
	}
	if err != nil {
		exitWithError(err)
	}

	hasConfig := config.Exists(path)
	conf, err := config.Load(path)
	if hasConfig && err != nil {
		exitWithError(err)
	}

	cmd := getCommand(arguments)
	logger.Println("cmd=" + cmd)

	if cmd == "prerelease" || cmd == "patch" || cmd == "minor" || cmd == "major" {

		if isDry == false && hasConfig {
			script := conf.GetPreVersion()
			if script != "" {
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

		newVersion, err := gump.DetermineTheNewTag(path, cmd, isBeta(arguments), isAlpha(arguments))
		logger.Println("newVersion=" + newVersion)
		if err != nil {
			exitWithError(err)
		}

		if isDry {
			fmt.Println("The new tag to create is: " + newVersion)
		}

		if isDry == false {
			ok, out, err := repoutils.CreateTag(vcs, path, newVersion)
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

		if isDry == false && hasConfig {
			script := conf.GetPostVersion()
			logger.Println("postversion=" + script)
			if script != "" {
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

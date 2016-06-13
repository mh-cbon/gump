package gump

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/gump/stringexec"
	"github.com/mh-cbon/semver"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

func GetTags(path string) ([]string, error) {
	vcs, err := repoutils.WhichVcs(path)
	logger.Println("vcs=" + vcs)
	if err != nil {
		return make([]string, 0), err
	}

	tags, err := repoutils.List(vcs, path)
	if err != nil {
		return make([]string, 0), err
	}
	tags = repoutils.FilterSemverTags(tags)
	tags = repoutils.SortSemverTags(tags)
	logger.Println("tags=")
	logger.Println(tags)
	return tags, nil
}

func GetMostRecentTag(tags []string) string {
	mostRecentTag := ""
	if len(tags) > 0 {
		mostRecentTag = tags[len(tags)-1]
	}
	logger.Println("mostRecentTag=" + mostRecentTag)

	if mostRecentTag == "" {
		mostRecentTag = "0.0.0"
	}
	return mostRecentTag
}

func CreateTheNewTag(how string, mostRecentTag string, beta bool, alpha bool) (string, error) {
	currentVersion, err := semver.NewVersion(mostRecentTag)
	if err != nil {
		return "", err
	}
	if how == "prerelease" {
		return IncrementPrerelease(currentVersion, beta, alpha)
	} else if how == "patch" {
		ok := currentVersion.IncPatch()
		if ok == false {
			return "", errors.New("Failed to increment the patch number in " + mostRecentTag)
		}

	} else if how == "minor" {
		ok := currentVersion.IncMinor()
		if ok == false {
			return "", errors.New("Failed to increment the minor number in " + mostRecentTag)
		}

	} else if how == "major" {
		ok := currentVersion.IncMajor()
		if ok == false {
			return "", errors.New("Failed to increment the major number in " + mostRecentTag)
		}

	}

	return currentVersion.String(), nil
}

func IncrementPrerelease(currentVersion *semver.Version, beta bool, alpha bool) (string, error) {
	if currentVersion.Prerelease() == "" {
		currentVersion.IncPatch()
		if alpha {
			currentVersion.SetPrerelease("alpha")
		} else if beta {
			currentVersion.SetPrerelease("beta")
		} else {
			currentVersion.SetPrerelease("alpha")
		}
	} else {
		re := regexp.MustCompile(`(alpha|beta)(-?\.?[0-9]+)?`)
		if re.MatchString(currentVersion.Prerelease()) == false {
			return "", errors.New("Cannot handle " + currentVersion.Prerelease())
		} else {
			parts := re.FindAllStringSubmatch(currentVersion.Prerelease(), -1)
			name := parts[0][1]
			sid := parts[0][2]
			if name == "alpha" && beta {
				currentVersion.SetPrerelease("beta")
			} else if name == "beta" && alpha {
				currentVersion.IncPatch() // downgrade from beta to alpha not possible without change patch number
				currentVersion.SetPrerelease("alpha")
			} else {
				d := ""
				id := 0
				var err error
				if sid != "" {
					p := sid[0:1]
					if p != "-" && p != "." {
						d = p
						id, err = strconv.Atoi(sid[1:])
						if err != nil {
							return "", err
						}
					} else {
						d = ""
						id, err = strconv.Atoi(sid)
						if err != nil {
							return "", err
						}
					}
				} else {
					d = ""
				}
				currentVersion.SetPrerelease(name + d + strconv.Itoa(id+1))
			}
		}
	}

	return currentVersion.String(), nil
}

func DetermineTheNewTag(path string, how string, beta bool, alpha bool) (string, error) {
	tags, err := GetTags(path)
	if err != nil {
		return "", err
	}

	mostRecentTag := GetMostRecentTag(tags)
	newVersion, err := CreateTheNewTag(how, mostRecentTag, beta, alpha)

	return newVersion, err
}

func ExecScript(cwd string, script string) (string, error) {
	cmd, err := stringexec.Command(cwd, script)
	if err != nil {
		return "", err
	}
	out, err := cmd.CombinedOutput()
	sOut := strings.TrimSuffix(string(out), "\n")
	sOut = strings.TrimSuffix(sOut, "\r")
	return sOut, err
}

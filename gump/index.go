package gump

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/mh-cbon/go-repo-utils/repoutils"
	"github.com/mh-cbon/gump/stringexec"
	"github.com/mh-cbon/semver"
	"github.com/mh-cbon/verbose"
)

var logger = verbose.Auto()

// GetTags return the list of tags contained in the underlying vcs
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

// GetMostRecentTag returns the most recent tag, could be prerelease,
// contained in the underlying vcs
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

// CreateTheNewTag Create the new version string
func CreateTheNewTag(how string, mostRecentTag string, beta bool, alpha bool) (string, error) {
	var newVersion semver.Version
	currentVersion, err := semver.NewVersion(mostRecentTag)
	if err != nil {
		logger.Println("Failed to parse tag=" + mostRecentTag)
		return "", err
	}

	logger.Println("Upgrading version to=" + how)

	if how == "prerelease" {
		return IncrementPrerelease(currentVersion, beta, alpha)

	} else if how == "patch" {
		newVersion = currentVersion.IncPatch()

	} else if how == "minor" {
		newVersion = currentVersion.IncMinor()

	} else if how == "major" {
		newVersion = currentVersion.IncMajor()

	} else {
		return "", errors.New("Unknown verb '" + how + "'")
	}

	return newVersion.String(), nil
}

// IncrementPrerelease Given a version, increment it to reach the next prerelease value
func IncrementPrerelease(currentVersion *semver.Version, beta bool, alpha bool) (string, error) {
	var err error
	var newVersion semver.Version
	if currentVersion.Prerelease() == "" {
		newVersion = currentVersion.IncPatch()
		if alpha {
			newVersion, err = newVersion.SetPrerelease("alpha")
			if err != nil {
				return "", err
			}
		} else if beta {
			newVersion, err = newVersion.SetPrerelease("beta")
			if err != nil {
				return "", err
			}
		} else {
			newVersion, err = newVersion.SetPrerelease("alpha")
			if err != nil {
				return "", err
			}
		}
	} else {
		re := regexp.MustCompile(`(alpha|beta)([-.])?([0-9]+)?`)
		if re.MatchString(currentVersion.Prerelease()) == false {
			return "", errors.New("Cannot handle " + currentVersion.Prerelease())
		}

		parts := re.FindAllStringSubmatch(currentVersion.Prerelease(), -1)
		name := parts[0][1]
		sep := parts[0][2]
		sid := parts[0][3]

		logger.Printf("prerelease parts=%s\n", parts)

		if name == "alpha" && beta {
			newVersion, err = currentVersion.SetPrerelease("beta")
			if err != nil {
				return "", err
			}
		} else if name == "beta" && alpha {
			// downgrade from beta to alpha not possible without change patch number
			newVersion = currentVersion.IncPatch() // as a beta has lower precedence than patch, 1.0.0-beta => 1.0.0
			newVersion = newVersion.IncPatch()     // need to update again to move to 1.0.1
			newVersion, err = newVersion.SetPrerelease("alpha")
			if err != nil {
				return "", err
			}
		} else {
			if sid == "" {
				sid = "0"
			}
			id, err := strconv.Atoi(sid)
			if err != nil {
				logger.Printf("failed to strconv.Atoi x=%s\n", sid)
				return "", err
			}
			newVersion, err = currentVersion.SetPrerelease(name + sep + strconv.Itoa(id+1))
			if err != nil {
				return "", err
			}
		}
	}

	return newVersion.String(), nil
}

// DetermineTheNewTag Given a path managed by a vcs, create the new version string
func DetermineTheNewTag(path string, how string, beta bool, alpha bool) (string, error) {
	tags, err := GetTags(path)
	if err != nil {
		return "", err
	}

	mostRecentTag := GetMostRecentTag(tags)

	return CreateTheNewTag(how, mostRecentTag, beta, alpha)
}

// ExecScript execute a command string on the underlying command system (bash or cmd)
func ExecScript(cwd string, script string) error {
	cmd, err := stringexec.Command(cwd, script)
	if err != nil {
		return err
	}
	return cmd.Run()
}

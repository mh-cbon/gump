package main

import (
	"fmt"
	"github.com/mh-cbon/go-repo-utils/repoutils"
	"os/exec"
	"testing"
)

func TestGumpIncPatch(t *testing.T) {
	args := []string{"patch"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "1.0.3") == false {
		t.Errorf("Expected tags to contain 1.0.3, but it was not found in %s\n", tags)
	}
}

func TestGumpIncMinor(t *testing.T) {
	args := []string{"minor"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "1.1.0") == false {
		t.Errorf("Expected tags to contain 1.1.0, but it was not found in %s\n", tags)
	}
}

func TestGumpIncMajor(t *testing.T) {
	args := []string{"major"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.0") == false {
		t.Errorf("Expected tags to contain 2.0.0, but it was not found in %s\n", tags)
	}
}

func TestGumpIncPrereleaseAlpha(t *testing.T) {
	args := []string{"prerelease", "-a"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.1-alpha") == false {
		t.Errorf("Expected tags to contain 2.0.1-alpha, but it was not found in %s\n", tags)
	}
}

func TestGumpIncPrereleaseAlphaTwice(t *testing.T) {
	args := []string{"prerelease", "-a"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.1-alpha1") == false {
		t.Errorf("Expected tags to contain 2.0.1-alpha1, but it was not found in %s\n", tags)
	}
}

func TestGumpIncPrereleaseAlphaToBeta(t *testing.T) {
	args := []string{"prerelease", "-b"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.1-beta") == false {
		t.Errorf("Expected tags to contain 2.0.1-beta, but it was not found in %s\n", tags)
	}
}

func TestGumpIncPrereleaseBetaTwice(t *testing.T) {
	args := []string{"prerelease", "-b"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.1-beta1") == false {
		t.Errorf("Expected tags to contain 2.0.1-beta1, but it was not found in %s\n", tags)
	}
}

func TestGumpIncPrereleaseBetaToAlpha(t *testing.T) {
	args := []string{"prerelease", "-a"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "2.0.2-alpha") == false {
		t.Errorf("Expected tags to contain 2.0.2-alpha, but it was not found in %s\n", tags)
	}
}

func TestGumpWithOkVersionScripts(t *testing.T) {
	args := []string{"prerelease", "-b"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git_with_ok_scripts"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git_with_ok_scripts")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "0.0.1-beta") == false {
		t.Errorf("Expected tags to contain 0.0.1-beta, but it was not found in %s\n", tags)
	}
}

func TestGumpWithOkGlideScripts(t *testing.T) {
	args := []string{"prerelease", "-b"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git_with_glide_scripts"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == false {
		fmt.Println(string(out))
		t.Errorf("Expected success=true, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git_with_glide_scripts")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if contains(tags, "0.0.1-beta") == false {
		t.Errorf("Expected tags to contain 0.0.1-beta, but it was not found in %s\n", tags)
	}
}

func TestGumpWithKoScripts(t *testing.T) {
	args := []string{"prerelease", "-b"}
	bin := "/vagrant/build/gump"
	cmd := exec.Command(bin, args...)
	cmd.Dir = "/home/vagrant/git_with_ko_scripts"
	fmt.Printf("%s: %s %s\n", cmd.Dir, bin, args)

	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
	}
	if cmd.ProcessState.Success() == true {
		t.Errorf("Expected success=false, got success=%t\n", true)
	}
	tags, err := repoutils.List("git", "/home/vagrant/git_with_ko_scripts")
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s\n", err)
	}
	if len(tags) > 0 {
		t.Errorf("Expected tags to be empty, but it was found %s\n", tags)
	}
	if contains(tags, "0.0.1-beta") == true {
		t.Errorf("Expected tags to not contain 0.0.1-beta, but it was found in %s\n", tags)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

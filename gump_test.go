package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mh-cbon/go-repo-utils/repoutils"
)

var isWindows = runtime.GOOS == "windows"

type TestingStub struct{}

func (t *TestingStub) Errorf(s string, a ...interface{}) {
	log.Fatalf(s+"\n", a...)
}

type TestingExiter struct{ t *testing.T }

func (t *TestingExiter) Errorf(s string, a ...interface{}) {
	t.t.Fatalf(s, a...)
}

var gumpPath = "git_test/gump"

func init() {

	t := &TestingStub{}
	mustFileExists(t, "gump.go")

	os.RemoveAll("git_test")
	os.Mkdir("git_test", os.ModePerm)

	if runtime.GOOS == "windows" {
		gumpPath += ".exe"
	}
	os.Remove(gumpPath)
	cmd := makeCmd(".", "go", "build", "-o", gumpPath, "gump.go")
	mustExecOk(t, cmd)

	var err error
	gumpPath, err = filepath.Abs(gumpPath)
	mustNotErr(t, err)
}

func initGitDir(t Errorer, dir, file, script string, tags ...string) {
	os.MkdirAll(dir, os.ModePerm)
	mustExecOk(t, makeCmd(dir, "git", "init"))
	mustExecOk(t, makeCmd(dir, "git", "config", "user.email", "john@doe.com"))
	mustExecOk(t, makeCmd(dir, "git", "config", "user.name", "John Doe"))
	mustExecOk(t, makeCmd(dir, "touch", "tomate"))
	if file != "" {
		mustWriteFile(t, filepath.Join(dir, file), script)
	}
	mustExecOk(t, makeCmd(dir, "git", "add", "-A"))
	mustExecOk(t, makeCmd(dir, "git", "commit", "-m", "rev 1"))
	for _, tag := range tags {
		mustExecOk(t, makeCmd(dir, "git", "tag", tag))
	}
}

func TestGumpIncPatch(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "patch"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "1.0.3")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncMinor(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	cmd := makeCmd(dir, gumpPath, "minor")

	mustExecOk(tt, cmd)

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "1.1.0")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncMajor(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	cmd := makeCmd(dir, gumpPath, "major")

	mustExecOk(tt, cmd)

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.0")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPrereleaseAlpha(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.1-alpha")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPrereleaseAlphaTwice(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.1-alpha1")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPrereleaseAlphaToBeta(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPrereleaseBetaTwice(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.1-beta1")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPrereleaseBetaToAlpha(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.2-alpha")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpIncPatchWithMessage(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "notsemvertag", "v1.0.2", "v1.0.0")

	mustExecOk(tt, makeCmd(dir, gumpPath, "major"))                  // 2.0.0
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))       // 2.0.0-alpha
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))       // 2.0.0-alpha1
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))       // 2.0.0-beta
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))       // 2.0.0-beta1
	mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-a"))       // 2.0.1-alpha
	mustExecOk(tt, makeCmd(dir, gumpPath, "patch", "-m", "message")) // 2.0.2

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "2.0.2")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithOkVersionScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_ok_scripts"
	initGitDir(t, dir, ".version", `preversion: echo "hello"
  postversion: echo "goodbye"`)

	out := mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustContain(t, out, "hello")
	execOutMustContain(t, out, "Created new tag 0.0.1-beta")
	execOutMustContain(t, out, "goodbye")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithKoScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_ko_scripts"
	initGitDir(t, dir, ".version", `preversion: eccho "hello" \
		&& echo "mustnotdisplay1"
  postversion: eccho "goodbye" && echo "mustnotdisplay2"`)

	out := mustNotExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustNotContain(t, out, "mustnotdisplay1")
	execOutMustNotContain(t, out, "mustnotdisplay2")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustEmpty(tt, tags)
	mustNotContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithOkGlideScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_glide_scripts"
	initGitDir(t, dir, "glide.yaml", `scripts:
    preversion: echo "hello"
    postversion: echo "goodbye"`)

	out := mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustContain(t, out, "hello")
	execOutMustContain(t, out, "Created new tag 0.0.1-beta")
	execOutMustContain(t, out, "goodbye")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithKoGlideScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_ko_glide_scripts"
	initGitDir(t, dir, "glide.yaml", `scripts:
  preversion: |
    eccho "hello" \
    && echo "mustnotdisplay1"
  postversion: eccho "goodbye" && echo "mustnotdisplay2"
`)

	out := mustNotExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustNotContain(t, out, "mustnotdisplay1")
	execOutMustNotContain(t, out, "mustnotdisplay2")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustEmpty(tt, tags)
	mustNotContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithOkShScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_ko_sh_scripts"
	initGitDir(t, dir, ".version.sh", `preversion=
      echo "hello"
    postversion=
      echo "goodbye"`)

	out := mustExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustContain(t, out, "hello")
	execOutMustContain(t, out, "Created new tag 0.0.1-beta")
	execOutMustContain(t, out, "goodbye")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpWithKoShScripts(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git_with_sh_scripts"
	initGitDir(t, dir, ".version.sh", `preversion=
      eccho "hello"
			echo "mustnotdisplay1"
    postversion=
      eccho "goodbye"
		echo "mustnotdisplay2"`)

	out := mustNotExecOk(tt, makeCmd(dir, gumpPath, "prerelease", "-b"))
	execOutMustNotContain(t, out, "mustnotdisplay1")
	execOutMustNotContain(t, out, "mustnotdisplay2")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustEmpty(tt, tags)
	mustNotContain(tt, tags, "0.0.1-beta")

	mustNotErr(t, os.RemoveAll(dir))
}

func TestGumpMultipleMajor(t *testing.T) {
	tt := &TestingExiter{t}

	dir := "git_test/git"
	initGitDir(t, dir, "", "", "2.0.0", "v1.0.2", "v1.0.0")

	// on master last tag is 2.0.0, patch => 2.0.1
	mustExecOk(tt, makeCmd(dir, "git", "checkout", "master"))
	out := mustExecOk(tt, makeCmd(dir, gumpPath, "patch", "-d"))
	execOutMustContain(tt, out, "The new tag to create is: 2.0.1")

	// on 1.0.2, patch => 1.0.3
	mustExecOk(tt, makeCmd(dir, "git", "checkout", "v1.0.2"))
	out = mustExecOk(tt, makeCmd(dir, gumpPath, "patch", "-d"))
	execOutMustContain(tt, out, "The new tag to create is: 1.0.3")

	tags, err := repoutils.List("git", dir)
	mustNotErr(tt, err)
	mustContain(tt, tags, "1.0.3")

	mustNotErr(t, os.RemoveAll(dir))
	t.Errorf("what")
}

type Errorer interface {
	Errorf(string, ...interface{})
}

func execOutMustContain(t Errorer, out string, s string) bool {
	if strings.Index(out, s) == -1 {
		t.Errorf("Output does not match expected to contain %q\n%v\n", s, out)
		return false
	}
	return true
}

func execOutMustNotContain(t Errorer, out string, s string) bool {
	if !isWindows && strings.Index(out, s) > -1 {
		t.Errorf("Output does not match expected to NOT contain %q\n%v\n", s, out)
		return false
	} else if isWindows && strings.Index(out, "An has error occured while executing") == -1 {
		t.Errorf("Output does not match expected to NOT contain %q\n%v\n", s, out)
		return false
	}
	return true
}

func mustExecOk(t Errorer, cmd *exec.Cmd) string {
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("mustExecOk failed, out=\n%v\n-------------", string(out))
	}
	mustNotErr(t, err)
	mustSucceed(t, cmd)
	return string(out)
}
func mustNotExecOk(t Errorer, cmd *exec.Cmd) string {
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Println(string(out))
	}
	mustErr(t, err)
	if cmd != nil {
		mustNotSucceed(t, cmd)
	}
	return string(out)
}

func makeCmd(dir string, bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	fmt.Printf("%s: %s %s\n", dir, bin, args)
	return cmd
}
func mustSucceed(t Errorer, cmd *exec.Cmd) bool {
	if cmd.ProcessState.Success() == false {
		t.Errorf("Expected success=true, got success=%t\n", false)
		return false
	}
	return true
}
func mustNotSucceed(t Errorer, cmd *exec.Cmd) bool {
	if cmd != nil && cmd.ProcessState != nil && cmd.ProcessState.Success() {
		t.Errorf("Expected success=false, got success=%t\n", true)
		return false
	}
	return true
}
func mustErr(t Errorer, err error) bool {
	if err == nil {
		t.Errorf("Expected err!=nil, got err=%s\n", err)
		return false
	}
	return true
}
func mustNotErr(t Errorer, err error) bool {
	if err != nil {
		fmt.Println(err)
		t.Errorf("Expected err=nil, got err=%s\n", err)
		return false
	}
	return true
}

func mustNotContain(t Errorer, tags []string, tag string) bool {
	if contains(tags, tag) {
		t.Errorf("Expected tags to NOT contain %q, but it WAS found in %s\n", tag, tags)
		return false
	}
	return true
}

func mustEmpty(t Errorer, tags []string) bool {
	if len(tags) > 0 {
		t.Errorf("Expected tags to be empty, but it was found %s\n", tags)
		return false
	}
	return true
}

func mustContain(t Errorer, tags []string, tag string) bool {
	if contains(tags, tag) == false {
		t.Errorf("Expected tags to contain %q, but it was not found in %s\n", tag, tags)
		return false
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func mustFileExists(t Errorer, p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Errorf("file mut exists %q", p)
		return false
	}
	return true
}
func mustWriteFile(t Errorer, p string, c string) bool {
	if err := ioutil.WriteFile(p, []byte(c), os.ModePerm); err != nil {
		t.Errorf("file not written %q", p)
		return false
	}
	return true
}

// vcprompt is a simple Go program that prints version control system
// informations. It is designed to be used by shell prompts.
//
// You can customize the output of vcprompt using format strings:
//
//   vcprompt -f="%b"
//
// Format strings use printf-like "%" escape sequences:
//
// %n  current vcs name
// %b  current branch name
// %r  current revision
// %m  + if there are any uncommitted changes (added, modified, or
//     removed files)
//
// All other characters are expanded as-is.
//
// The default format string is
//
//	 "%n:%b"
//
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	githead       = ".git/HEAD"
	refPrefix     = "ref: refs/heads/"
	defaultFormat = `%n:%b`
)

var (
	debug  = flag.Bool("d", false, "debug")
	format = flag.String("f", defaultFormat, "format")
)

// vcs represents a version-control-system state through a user perspective.
type vcs struct {
	available bool

	name       string
	branch     string
	revision   string
	isModified bool
}

func (v vcs) String() string {
	if !v.available {
		return ""
	}

	var buf bytes.Buffer
	var eof rune = 0

	reader := bufio.NewReader(strings.NewReader(*format))

	for {
		r, _, _ := reader.ReadRune()
		if r == eof {
			break
		}

		// write ordinary characters.
		if r != '%' {
			buf.WriteString(string(r))
			continue
		}

		// we have format string
		next, _, _ := reader.ReadRune()
		switch next {
		case 'n': // version control system name
			buf.WriteString(v.name)
		case 'b': // branch name
			buf.WriteString(v.branch)
		case 'r': // revision number
			buf.WriteString(v.revision)
		case 'm': // is modified flag
			if v.isModified {
				buf.WriteString("+")
			}
		default:
			buf.WriteString(string(next))
		}
	}

	return buf.String()
}

// gitInfo checks for a git project and extracts several states of it, such as
// branch, revision and etc.
func gitInfo() vcs {
	v := vcs{name: "git", available: true}

	cwd := probeParent()
	if cwd == "" {
		printdebug("no .git/ directory found")
		v.available = false
		return v
	}

	line, err := readFirstLine(path.Join(cwd, githead))
	if err != nil {
		printdebug(err.Error())
		return v
	}

	// if refPrefix is not found on HEAD, assume it is a revision
	if strings.HasPrefix(line, refPrefix) {
		v.branch = line[len(refPrefix):]
	} else {
		v.revision = line
	}

	v.isModified = isModified()

	return v
}

// isModified reports whether there are things that are modified.
func isModified() bool {
	cmd := exec.Command("git", "diff", "--no-ext-diff", "--quiet", "--exit-code")
	if err := cmd.Run(); err != nil {
		// ExitError indicates there is a change
		if _, ok := err.(*exec.ExitError); ok {
			return true
		}
	}

	return false
}

// probeParent tries to find a ".git" directory until it hits root directory.
func probeParent() string {
	var cwd string
	for {
		cwd, _ = os.Getwd()
		if pathExists(".git") {
			return cwd
		}

		if cwd == "/" {
			return ""
		}

		os.Chdir("..")
	}
}

// readFirstLine reads the first line of the given filename.
func readFirstLine(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("unable to read first line of %s", filename)
	}
	return strings.TrimSpace(line), nil
}

func pathExists(dir string) bool {
	f, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if !f.IsDir() {
		return false
	}

	return true
}

func printdebug(format string, a ...interface{}) {
	if *debug {
		fmt.Printf(format, a...)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: vcprompt [options]")
	fmt.Fprintln(os.Stderr, "options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "formats:")
	fmt.Fprintln(os.Stderr, `  %n show vcs name`)
	fmt.Fprintln(os.Stderr, `  %b show branch`)
	fmt.Fprintln(os.Stderr, `  %r show revision`)
	fmt.Fprintln(os.Stderr, `  %m show modified`)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	fmt.Print(gitInfo())
}

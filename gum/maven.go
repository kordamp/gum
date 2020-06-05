package gum

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type mavenCommand struct {
	quiet             bool
	executable        string
	args              []string
	buildFile         string
	explicitBuildFile string
	rootBuildFile     string
}

func (c mavenCommand) Execute() {
	args := make([]string, 0)

	banner := make([]string, 0)
	banner = append(banner, "Using maven at '"+c.executable+"'")
	nearest, nargs := GrabFlag("-gn", c.args)
	debug, nargs := GrabFlag("-gd", nargs)

	if debug {
		fmt.Println("nearest            = ", nearest)
		fmt.Println("args               = ", nargs)
		fmt.Println("rootBuildFile      = ", c.rootBuildFile)
		fmt.Println("buildFile          = ", c.buildFile)
		fmt.Println("explicitBuildFile  = ", c.explicitBuildFile)
		fmt.Println("")
	}

	if len(c.explicitBuildFile) > 0 {
		banner = append(banner, "to run buildFile '"+c.explicitBuildFile+"':")
	} else if nearest && len(c.buildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.buildFile)
		banner = append(banner, "to run buildFile '"+c.buildFile+"':")
	} else if len(c.rootBuildFile) > 0 {
		args = append(args, "-f")
		args = append(args, c.rootBuildFile)
		banner = append(banner, "to run buildFile '"+c.rootBuildFile+"':")
	}

	for i := range nargs {
		args = append(args, nargs[i])
	}

	if !c.quiet {
		fmt.Println(strings.Join(banner, " "))
	}

	cmd := exec.Command(c.executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// FindMaven finds and executes mvnw/mvn
func FindMaven(quiet bool, explicit bool, args []string) Command {
	pwd := getWorkingDir()

	mvnw, noWrapper := findMavenWrapperExec(pwd)
	mvn, noMaven := findMavenExec()
	explicitBuildFileSet, explicitBuildFile := findExplicitMavenBuildFile(args)

	var executable string
	if noWrapper == nil {
		executable = mvnw
	} else if noMaven == nil {
		warnNoMavenWrapper(quiet, explicit)
		executable = mvn
	} else {
		warnNoMaven(quiet, explicit)

		if explicit {
			os.Exit(-1)
		}
		return nil
	}

	if explicitBuildFileSet {
		return mavenCommand{
			quiet:             quiet,
			executable:        executable,
			args:              args,
			explicitBuildFile: explicitBuildFile}
	}

	rootBuildFile, noRootBuildFile := findMavenRootFile(filepath.Join(pwd, ".."), args)
	buildFile, noBuildFile := findMavenBuildFile(pwd, args)

	if noRootBuildFile != nil {
		rootBuildFile = buildFile
	}

	if noBuildFile != nil {
		if explicit {
			fmt.Println("No Maven project found")
			fmt.Println()
			os.Exit(-1)
		}
		return nil
	}

	return mavenCommand{
		quiet:         quiet,
		executable:    executable,
		args:          args,
		rootBuildFile: rootBuildFile,
		buildFile:     buildFile}
}

func warnNoMavenWrapper(quiet bool, explicit bool) {
	if !quiet && explicit {
		fmt.Printf("No %s set up for this project. ", resolveMavenWrapperExec())
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://maven.apache.org/)")
		fmt.Println()
	}
}

func warnNoMaven(quiet bool, explicit bool) {
	if !quiet && explicit {
		fmt.Printf("No %s found in path. Please install Maven.", resolveMavenExec())
		fmt.Println("(https://maven.apache.org/download.cgi)")
		fmt.Println()
	}
}

// Finds the maven executable
func findMavenExec() (string, error) {
	maven := resolveMavenExec()
	paths := getPaths()

	for i := range paths {
		name := filepath.Join(paths[i], maven)
		if fileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(maven + " not found")
}

// Finds the Maven wrapper (if it exists)
func findMavenWrapperExec(dir string) (string, error) {
	wrapper := resolveMavenWrapperExec()
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if fileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenWrapperExec(parentdir)
}

func findExplicitMavenBuildFile(args []string) (bool, string) {
	found, file := findFlag("-f", args)
	if !found {
		found, file = findFlag("--file", args)
	}

	if found {
		file, _ = filepath.Abs(file)
		return true, file
	}

	return false, ""
}

// Finds the nearest pom.xml
func findMavenBuildFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if fileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenBuildFile(parentdir, args)
}

// Finds the root pom.xml
func findMavenRootFile(dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if fileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenRootFile(parentdir, args)
}

// Resolves the mvnw executable (OS dependent)
func resolveMavenWrapperExec() string {
	if isWindows() {
		return "mvnw.bat"
	}
	return "mvnw"
}

// Resolves the mvn executable (OS dependent)
func resolveMavenExec() string {
	if isWindows() {
		return "mvn.bat"
	}
	return "mvn"
}

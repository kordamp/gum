package gum

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MavenCommand defines an executable Maven command
type MavenCommand struct {
	context           Context
	executable        string
	args              []string
	buildFile         string
	explicitBuildFile string
	rootBuildFile     string
}

// Execute executes the given command
func (c MavenCommand) Execute() {
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

	if !c.context.IsQuiet() {
		fmt.Println(strings.Join(banner, " "))
	}

	cmd := exec.Command(c.executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// FindMaven finds and executes mvnw/mvn
func FindMaven(context Context, args []string) *MavenCommand {
	pwd := context.GetWorkingDir()

	mvnw, noWrapper := findMavenWrapperExec(context, pwd)
	mvn, noMaven := findMavenExec(context)
	explicitBuildFileSet, explicitBuildFile := findExplicitMavenBuildFile(args)

	var executable string
	if noWrapper == nil {
		executable = mvnw
	} else if noMaven == nil {
		warnNoMavenWrapper(context)
		executable = mvn
	} else {
		warnNoMaven(context)

		if context.IsExplicit() {
			context.Exit(-1)
		}
		return nil
	}

	if explicitBuildFileSet {
		return &MavenCommand{
			context:           context,
			executable:        executable,
			args:              args,
			explicitBuildFile: explicitBuildFile}
	}

	rootBuildFile, noRootBuildFile := findMavenRootFile(context, filepath.Join(pwd, ".."), args)
	buildFile, noBuildFile := findMavenBuildFile(context, pwd, args)

	if noRootBuildFile != nil {
		rootBuildFile = buildFile
	}

	if noBuildFile != nil {
		if context.IsExplicit() {
			fmt.Println("No Maven project found")
			fmt.Println()
			context.Exit(-1)
		}
		return nil
	}

	return &MavenCommand{
		context:       context,
		executable:    executable,
		args:          args,
		rootBuildFile: rootBuildFile,
		buildFile:     buildFile}
}

func warnNoMavenWrapper(context Context) {
	if !context.IsQuiet() && context.IsExplicit() {
		fmt.Printf("No %s set up for this project. ", resolveMavenWrapperExec(context))
		fmt.Println("Please consider setting one up.")
		fmt.Println("(https://maven.apache.org/)")
		fmt.Println()
	}
}

func warnNoMaven(context Context) {
	if !context.IsQuiet() && context.IsExplicit() {
		fmt.Printf("No %s found in path. Please install Maven.", resolveMavenExec(context))
		fmt.Println("(https://maven.apache.org/download.cgi)")
		fmt.Println()
	}
}

// Finds the maven executable
func findMavenExec(context Context) (string, error) {
	maven := resolveMavenExec(context)
	paths := context.GetPaths()

	for i := range paths {
		name := filepath.Join(paths[i], maven)
		if context.FileExists(name) {
			return filepath.Abs(name)
		}
	}

	return "", errors.New(maven + " not found")
}

// Finds the Maven wrapper (if it exists)
func findMavenWrapperExec(context Context, dir string) (string, error) {
	wrapper := resolveMavenWrapperExec(context)
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New(wrapper + " not found")
	}

	path := filepath.Join(dir, wrapper)
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenWrapperExec(context, parentdir)
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
func findMavenBuildFile(context Context, dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenBuildFile(context, parentdir, args)
}

// Finds the root pom.xml
func findMavenRootFile(context Context, dir string, args []string) (string, error) {
	parentdir := filepath.Join(dir, "..")

	if parentdir == dir {
		return "", errors.New("Did not find root pom.xml")
	}

	path := filepath.Join(dir, "pom.xml")
	if context.FileExists(path) {
		return filepath.Abs(path)
	}

	return findMavenRootFile(context, parentdir, args)
}

// Resolves the mvnw executable (OS dependent)
func resolveMavenWrapperExec(context Context) string {
	if context.IsWindows() {
		return "mvnw.bat"
	}
	return "mvnw"
}

// Resolves the mvn executable (OS dependent)
func resolveMavenExec(context Context) string {
	if context.IsWindows() {
		return "mvn.bat"
	}
	return "mvn"
}

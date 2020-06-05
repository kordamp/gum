package gum

import (
	"path/filepath"
	"testing"
)

func TestGradleSingleWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-with-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleSingleWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "gradle")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleParentWithWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleParentWithoutWrapper(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-without-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(bin, "gradle")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithExplicitBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-explicit", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{"-b", filepath.Join(pwd, "explicit.gradle")})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, ""},
		{"BuildFile", cmd.buildFile, ""},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, filepath.Join(pwd, "explicit.gradle")},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithExplicitSettingsFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{"-c", filepath.Join(pwd, "..", "settings.gradle")})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "build.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithExplicitProjectDir(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-wrapper", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{"-p", filepath.Join(pwd, "..")})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"Executable", cmd.executable, filepath.Join(pwd, "..", "gradlew")},
		{"RootBuildFile", cmd.rootBuildFile, ""},
		{"BuildFile", cmd.buildFile, ""},
		{"SettingsFile", cmd.settingsFile, ""},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, filepath.Join(pwd, "..")},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithNearestBuildFile(t *testing.T) {
	// given:
	bin, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "bin"))
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "parent-with-conventional-child", "child"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{bin}}

	// when:
	cmd := FindGradle(context, []string{"-gn"})

	// then:
	if cmd == nil {
		t.Error("Expected a command but got nil")
	}

	var checks = []struct {
		title, actual, expected string
	}{
		{"RootBuildFile", cmd.rootBuildFile, filepath.Join(pwd, "..", "build.gradle")},
		{"BuildFile", cmd.buildFile, filepath.Join(pwd, "child.gradle")},
		{"SettingsFile", cmd.settingsFile, filepath.Join(pwd, "..", "settings.gradle")},
		{"ExplicitBuildFile", cmd.explicitBuildFile, ""},
		{"ExplicitSettingsFile", cmd.explicitSettingsFile, ""},
		{"ExplicitProjectDir", cmd.explicitProjectDir, ""},
	}

	for _, check := range checks {
		if check.actual != check.expected {
			t.Errorf("%s: got %s, want %s", check.title, check.actual, check.expected)
		}
	}
}

func TestGradleWithoutExecutables(t *testing.T) {
	// given:
	pwd, _ := filepath.Abs(filepath.Join("..", "tests", "gradle", "single-without-wrapper"))

	context := testContext{
		quiet:      true,
		explicit:   true,
		windows:    false,
		workingDir: pwd,
		paths:      []string{}}

	// when:
	cmd := FindGradle(context, []string{})

	// then:
	if cmd != nil {
		t.Error("Expected a nil command but got something")
	}
}

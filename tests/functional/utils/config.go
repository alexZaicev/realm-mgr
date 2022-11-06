package utils

import (
	"errors"
	"flag"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/alexZaicev/realm-mgr/internal/drivers/config"
)

const (
	defaultProjectRoot = "realm-mgr"
	defaultConfigFile  = "./tests/functional/bench-config.yml"
)

var configFile string
var projectRootDirName string

// Regex for the project directory root
var re *regexp.Regexp

// init sets up flags used to determine the config for running functional tests
func init() { //nolint:gochecknoinits // init is the most appropriate place to setup flags
	flag.StringVar(
		&configFile,
		"config-file",
		defaultConfigFile,
		"path to config file for testing (relative to project root), defaults to bench config (use \"none\" for no config)",
	)
	flag.StringVar(
		&projectRootDirName,
		"project-root-dir",
		defaultProjectRoot,
		"directory name of the project root, only needs to be set if you changed the dir name when cloning the project",
	)
}

func LoadConfig() (config.Config, error) {
	configFilePath := ""
	if configFile != "none" {
		var err error
		configFilePath, err = getConfigFilePath()
		if err != nil {
			return nil, err
		}
	}

	cfg, err := config.NewKoanfConfig(".").
		WithYAML(configFilePath, false /*optional*/).
		TreatInt64AsInt(true).
		Load()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	projectRoot, err := getProjectDirRoot()
	if err != nil {
		return "", err
	}

	if configFile == "" {
		// If flag provided with empty value, use default config file
		configFile = defaultConfigFile
	}
	if !strings.HasPrefix(configFile, "/") {
		configFile = "/" + configFile
	}

	return projectRoot + configFile, nil
}

func getProjectDirRoot() (string, error) {
	// Regex for the project directory root (assuming that no subdirectories have the same name)
	// Stop this running more than once
	if re == nil {
		re = regexp.MustCompile(`^.*/` + projectRootDirName + `(_(master|PR-\d+))?$`)
	}
	// Get the current working directory and work upwards checking for a match with the above regex to find the project directory root
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var rootPath []byte
	for workingDir != "/" {
		rootPath = re.Find([]byte(workingDir))
		if rootPath != nil {
			return string(rootPath), nil
		}
		workingDir = path.Dir(workingDir)
	}
	return "", errors.New("could not find project root")
}

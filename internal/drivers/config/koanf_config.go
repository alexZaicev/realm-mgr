package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

// Basic enums used for switching file reading logic.
const (
	tomlFile = "toml"
	yamlFile = "yaml"
	jsonFile = "json"
)

// koanfSource provides an abstraction over the loading of config
// to allow for precedence to be established over different
// types of sources (files, configmap volumes, env vars)
type koanfSource interface {
	Load(k *KoanfConfig) error
}

// koanfConfigFile models a config file that should be loaded
// by KoanfConfig.
type koanfConfigFile struct {
	path     string
	fileType string
	optional bool
}

// Load implements the koanfSource interface using the underlying
// koanf logic to load a config file.
func (configFile koanfConfigFile) Load(k *KoanfConfig) error {
	if configFile.optional {
		if _, err := fs.Stat(k.statFS, configFile.path); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
	}

	var parser koanf.Parser
	switch configFile.fileType {
	case tomlFile:
		parser = toml.Parser()
	case yamlFile:
		parser = yaml.Parser()
	case jsonFile:
		parser = json.Parser()
	default:
		// The fileType var value is controlled internally to this module,
		// but this accounts for future changes
		return fmt.Errorf("unexpected file type: %s", configFile.fileType)
	}

	return k.config.Load(file.Provider(configFile.path), parser)
}

// koanfConfigmapVolume models a K8S configmap that has been
// mounted as a volume that should be loaded by KoanfConfig.
type koanfConfigmapVolume struct {
	path           string
	optional       bool
	valueConverter func(key string, value string) any
}

// Load implements the koanfSource interface using the underlying
// koanf logic to load a configmap volume.
func (volume koanfConfigmapVolume) Load(k *KoanfConfig) error {
	if volume.optional {
		if _, err := fs.Stat(k.statFS, volume.path); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
	}

	return k.config.Load(ConfigmapVolumeProvider(volume.path, k.delim, volume.valueConverter), nil)
}

// koanfEnvVars model the logic required to load config from
// environment variables.
type koanfEnvVars struct {
	prefix       string
	envConverter func(key string, value string) (string, any)
}

// Load implements the koanfSource interface using the underlying
// koanf logic to load env vars.
func (envVars koanfEnvVars) Load(k *KoanfConfig) error {
	return k.config.Load(env.ProviderWithValue(envVars.prefix, k.delim, envVars.envConverter), nil)
}

// KoanfConfig implements that Config and MustConfig interfaces
// using the Koanf package. It allows for config to be loaded from
// any combination of YAML, JSON or TOML files, as well as from env vars
// and K8S config maps mounted as volumes. Config sources are loaded
// in the order that they are specified, with duplicate keys in the later
// sources overriding the previous values.
type KoanfConfig struct {
	config          *koanf.Koanf
	delim           string
	sources         []koanfSource
	treatInt64AsInt bool

	// This var is used internally to perform the optional
	// checks on files and volumes
	statFS fs.FS
}

// NewKoanfConfig initialises a new Koanf based config loader, using
// the specified deliminator for nested config keys. E.g. if delim is
// specified as '.' then a nested config key would look like 'database.host',
// where as if '/' was specified it would be 'database/host'.
// To add sources of config, use the WithX functions on the returned
// struct.
//
//nolint:misspell // This is England
func NewKoanfConfig(delim string) *KoanfConfig {
	return &KoanfConfig{
		config:          koanf.New(delim),
		delim:           delim,
		treatInt64AsInt: false,
		statFS:          os.DirFS("/"),
	}
}

// withFile provides an internal, generic version of the public WithX file
// functions.
func (k *KoanfConfig) withFile(path, fileType string, optional bool) *KoanfConfig {
	k.sources = append(
		k.sources,
		koanfConfigFile{
			path:     path,
			fileType: fileType,
			optional: optional,
		},
	)
	return k
}

// WithTOML adds the TOML file specified (as an absolute path) by the file path.
// The optional flag can be used to indicate that the file not existing should
// not be considered an error.
func (k *KoanfConfig) WithTOML(path string, optional bool) *KoanfConfig {
	return k.withFile(path, tomlFile, optional)
}

// WithYAML adds the YAML file specified (as an absolute path) by the file path.
// The optional flag can be used to indicate that the file not existing should
// not be considered an error.
func (k *KoanfConfig) WithYAML(path string, optional bool) *KoanfConfig {
	return k.withFile(path, yamlFile, optional)
}

// WithJSON adds the JSON file specified (as an absolute path) by the file path.
// The optional flag can be used to indicate that the file not existing should
// not be considered an error.
func (k *KoanfConfig) WithJSON(path string, optional bool) *KoanfConfig {
	return k.withFile(path, jsonFile, optional)
}

// WithConfigmapVolume adds the K8S config map mounted as a volume specified
// (as an absolute path) by the path provided. By default the config values
// will be loaded as strings.  A value conversion function can be provided
// to convert specified keys to the desired data types.
func (k *KoanfConfig) WithConfigmapVolume(
	path string,
	optional bool,
	valueConverter func(key string, value string) any,
) *KoanfConfig {
	k.sources = append(
		k.sources,
		koanfConfigmapVolume{
			path:           path,
			optional:       optional,
			valueConverter: valueConverter,
		},
	)
	return k
}

// WithEnvVars adds env var loading to the KoanfConfig logic. The env vars
// are expected to be prefixed to provide namespacing safely (e.g. APP_), and only
// env vars with the specified namespace will be loaded. An envConverter function
// can be provided to not only convert the config values to the desired data type
// (by default they will be strings), but also to convert the environment variable name
// to the desired key. E.g. to convert 'APP_DATABASE_HOST' to 'database.host'.
func (k *KoanfConfig) WithEnvVars(prefix string, envConverter func(string, string) (string, any)) *KoanfConfig {
	k.sources = append(
		k.sources,
		koanfEnvVars{
			prefix:       prefix,
			envConverter: envConverter,
		},
	)
	return k
}

// EnvConverterFromKeyMap provides a utility function to generate an envConverter function
// from an existing valueConverter function and a map translating env var keys to the
// desired config keys. This allows the same valueConverter function to be re-used
// between ConfigmapVolume and EnvVar loading logic.
func EnvConverterFromKeyMap(keyMap map[string]string, valueConverter func(string, string) any) func(string, string) (string, any) {
	return func(key string, value string) (string, any) {
		outKey := key
		if convertedKey, exists := keyMap[key]; exists {
			outKey = convertedKey
		}

		var outValue any = value
		if valueConverter != nil {
			outValue = valueConverter(outKey, value)
		}

		return outKey, outValue
	}
}

// Load loads in the config values from the registered config files, volume mounts and
// env vars. The precedence for loading config are files (in the order they are added),
// then K8S configmap volumes followed by environment variables, with duplicate keys in
// the later sources overriding the previous value.
func (k *KoanfConfig) Load() (*KoanfConfig, error) {
	for _, source := range k.sources {
		if err := source.Load(k); err != nil {
			return nil, err
		}
	}
	return k, nil
}

// TreatInt64AsInt indicates if the config struct should return
// values that are stored in the inner koanf.Koanf as int64 as
// ints. This is useful if you want to avoid casting int64 values
// within your own code.
func (k *KoanfConfig) TreatInt64AsInt(convertInt64 bool) *KoanfConfig {
	k.treatInt64AsInt = convertInt64
	return k
}

// Get fetches the config value using the specified key. If no entry
// exists with the specified key, nil is returned.
func (k *KoanfConfig) Get(key string) any {
	val := k.config.Get(key)
	if val == nil {
		return nil
	}

	if k.treatInt64AsInt {
		if castVal, isInt64 := val.(int64); isInt64 {
			return int(castVal)
		}
	}

	return val
}

// Keys returns a slice of all the keys for which there are config entries.
func (k *KoanfConfig) Keys() []string {
	return k.config.Keys()
}

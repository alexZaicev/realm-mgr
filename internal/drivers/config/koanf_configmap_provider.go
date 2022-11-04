package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// KoanfConfigmapVolumeProvider implements the koanf package's Provider
// interface, allowing the loading of config maps that are mounted as
// volumes in K8S, where each configmap entry is presented as a file.
type KoanfConfigmapVolumeProvider struct {
	dir            string
	delim          string
	valueConverter func(string, string) any
	getDataDirFS   func(string) (fs.FS, error)
}

// ConfigmapVolumeProvider creates the provider for the specified volume,
// using the defined deliminator to add nested keys. As all values are
// strings by default, given that they are file contents, a valueConverter
// function can be provided to convert given keys to the correct data type.
// By default the provider uses a KIND aware function for getting the
// directory with config data in it.
func ConfigmapVolumeProvider(dir, delim string, valueConverter func(string, string) any) *KoanfConfigmapVolumeProvider {
	return &KoanfConfigmapVolumeProvider{
		dir:            dir,
		valueConverter: valueConverter,
		delim:          delim,
		getDataDirFS:   kindAwareGetDataDirFS,
	}
}

// SetDataDirFSFetcher allows for the overriding of the function that the
// KoanfConfigmapVolumeProvider uses to open the directory that it reads the
// config entry files from. This is due to the fact that some on bench targeted
// implementations of K8S (e.g. KIND) perform symlinking and other tricks
// that might need to be accounted for.
func (k *KoanfConfigmapVolumeProvider) SetDataDirFSFetcher(inFunc func(string) (fs.FS, error)) {
	k.getDataDirFS = inFunc
}

// ReadBytes is not implemented for the KoanfConfigmapVolumeProvider.
// Instead Read() is used.
func (k *KoanfConfigmapVolumeProvider) ReadBytes() ([]byte, error) {
	return nil, errors.New("ReadBytes not implemented")
}

// kindAwareGetDataDirFS accounts for the slight strangeness
// reading mounted config files in KIND - due to the fact that
// Go's FS module doesn't follow symlinks by default.
func kindAwareGetDataDirFS(volume string) (fs.FS, error) {
	// Looks like KIND has some dodgy symlinking going on when mounting
	// config maps, so we account for that
	dataDir := filepath.Join(volume, "..data")
	if _, err := os.Stat(dataDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return os.DirFS(volume), nil
		}
		return nil, err
	}

	actualDir, err := filepath.EvalSymlinks(dataDir)
	if err != nil {
		return nil, err
	}

	return os.DirFS(actualDir), nil
}

// Read implements the Koanf Provider interface and reads in the volume's files, adding
// them as entries in Koanf's map.
func (k *KoanfConfigmapVolumeProvider) Read() (map[string]any, error) {
	fileSystem, err := k.getDataDirFS(k.dir)
	if err != nil {
		return nil, err
	}

	return k.read(fileSystem)
}

// read performs the actual logic of reading in the files and nested dirs inside the volume,
// converting values if required and adding them to the config map.
func (k *KoanfConfigmapVolumeProvider) read(fileSystem fs.FS) (map[string]any, error) {
	values := make(map[string]any)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		// We can assume that we have full permissions to read
		// all files mounted from configMaps, so treat any error
		// as a termination
		if err != nil {
			return err
		}

		// We can skip to top level dir as that's just the top level config map
		if path == "." {
			return nil
		}

		// Nested config (e.g. database.host) are modeled as nested maps.
		// Use the nested file names to drill down the map we need to add
		// the actual key/value to.
		mapToAddValueTo := values
		parts := strings.Split(path, string(os.PathSeparator))
		for i := 0; i < len(parts)-1; i++ {
			var ok bool
			mapToAddValueTo, ok = mapToAddValueTo[parts[i]].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to extract nested config key from path: %s", path)
			}
		}

		// If we come across a dir it means that we have nested config, so model that
		// as a nested map.
		if d.IsDir() {
			mapToAddValueTo[d.Name()] = make(map[string]any)
			return nil
		}

		value, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return err
		}

		var val any = string(value)
		if k.valueConverter != nil {
			key := strings.Join(parts, k.delim)
			val = k.valueConverter(key, string(value))
		}
		mapToAddValueTo[d.Name()] = val
		return nil
	})

	if err != nil {
		return nil, err
	}
	return values, nil
}

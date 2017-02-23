package config

import "fmt"

// Finder is a plugable config loader.
type Finder struct {
	loaders []Loader
}

// NewFinder initialize a config loader with preconfigured loader.
func NewFinder() *Finder {
	return &Finder{
		loaders: []Loader{&ShConfig{}, &SimpleConfig{}, &GlideConfig{}},
	}
}

// Load identify, read and parse the version script from givne directory.
func (c *Finder) Load(wd string) (Configured, error) {
	for _, loader := range c.loaders {
		if loader.Exists(wd) {
			return loader, loader.LoadDefault(wd)
		}
	}
	return nil, fmt.Errorf("Scripts not found in %q", wd)
}

// Loader can tell if a config exists, and load it.
type Loader interface {
	Configured
	Exists(wd string) bool
	LoadDefault(wd string) error
	Parse(data []byte) error
	Load(path string) error
}

// Configured is a facade to a glide.yaml / . version script file
type Configured interface {
	GetScript(name string) (string, bool)
}

package main

type IndexPath struct {
	// Tilde expansion of homedir is supported
	Path string

	Included bool
}

type Config struct {
	IndexPaths []IndexPath
}

func (c *Config) Validate() error {
	return nil
}

package config

type Overrides struct {
	Root            *string
	Host            *string
	Port            *int
	Theme           *string
	Write           *bool
	ShowHidden      *bool
	MaxPreviewBytes *int64
	Columns         ColumnOverrides
}

func (c *Config) apply(overrides Overrides) {
	if overrides.Root != nil {
		c.Root = *overrides.Root
	}
	if overrides.Host != nil {
		c.Host = *overrides.Host
	}
	if overrides.Port != nil {
		c.Port = *overrides.Port
	}
	if overrides.Theme != nil {
		c.Theme = *overrides.Theme
	}
	if overrides.Write != nil {
		c.Write = *overrides.Write
	}
	if overrides.ShowHidden != nil {
		c.ShowHidden = *overrides.ShowHidden
	}
	if overrides.MaxPreviewBytes != nil {
		c.MaxPreviewBytes = *overrides.MaxPreviewBytes
	}
	c.Columns.apply(overrides.Columns)
}

func (c *Config) applyRoot(overrides Overrides) {
	if overrides.Root != nil {
		c.Root = *overrides.Root
	}
}

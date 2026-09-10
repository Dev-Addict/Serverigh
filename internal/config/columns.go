package config

type Columns struct {
	Size     bool
	Modified bool
	Created  bool
	Mode     bool
}

type ColumnOverrides struct {
	Size     *bool
	Modified *bool
	Created  *bool
	Mode     *bool
}

func DefaultColumns() Columns {
	return Columns{
		Size:     true,
		Modified: true,
		Created:  true,
		Mode:     true,
	}
}

func (c *Columns) apply(overrides ColumnOverrides) {
	if overrides.Size != nil {
		c.Size = *overrides.Size
	}
	if overrides.Modified != nil {
		c.Modified = *overrides.Modified
	}
	if overrides.Created != nil {
		c.Created = *overrides.Created
	}
	if overrides.Mode != nil {
		c.Mode = *overrides.Mode
	}
}

package config

func LoadColumns(root string, flagOverrides ColumnOverrides) (Columns, error) {
	cfg, err := Default()
	if err != nil {
		return Columns{}, err
	}
	cfg.Root = root

	loaded, err := Load(cfg, Overrides{Columns: flagOverrides})
	if err != nil {
		return Columns{}, err
	}

	return loaded.Columns, nil
}

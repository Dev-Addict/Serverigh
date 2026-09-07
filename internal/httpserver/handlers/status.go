package handlers

type StatusView struct {
	Label      string
	ItemCount  int
	Mode       string
	Truncated  bool
	EntryLimit int
	OOB        bool
}

func modeLabel(writeEnabled bool) string {
	if writeEnabled {
		return "write-enabled"
	}

	return "read-only"
}

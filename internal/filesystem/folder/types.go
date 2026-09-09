package folder

const maxTreeEntries = 500

type Tree struct {
	Entries   []Entry
	Truncated bool
}

type Entry struct {
	Name        string
	Path        string
	Depth       int
	Expanded    bool
	HasChildren bool
	Selected    bool
}

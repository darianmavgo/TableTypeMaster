package tabletypes

// TableType represents a tabular data format and its attributes.
type TableType struct {
	Name            string
	FileExt         string
	MimeType        string
	MagicNumber     string
	ListTables      string
	ListColumns     string
	ListColumnTypes string
}

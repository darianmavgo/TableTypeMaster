package tabletypes

import (
	"html/template"
	"os"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Table Types</title>
	<style>
		table { border-collapse: collapse; width: 100%; }
		th, td { border: 1px solid black; padding: 8px; text-align: left; }
		th { background-color: #f2f2f2; }
	</style>
</head>
<body>
	<h1>Table Types</h1>
	<table>
		<tr>
			<th>Name</th>
			<th>File Ext</th>
			<th>Mime Type</th>
			<th>Magic Number</th>
			<th>List Tables</th>
			<th>List Columns</th>
			<th>List Column Types</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.Name}}</td>
			<td>{{.FileExt}}</td>
			<td>{{.MimeType}}</td>
			<td>{{.MagicNumber}}</td>
			<td>{{.ListTables}}</td>
			<td>{{.ListColumns}}</td>
			<td>{{.ListColumnTypes}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`

// ExportToHTML exports the table types data to an HTML file.
func ExportToHTML(filename string) error {
	t, err := template.New("tabletypes").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, AllTableTypes)
}

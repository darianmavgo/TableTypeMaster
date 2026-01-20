package main

import (
	"fmt"
	"log"

	"github.com/darianmavgo/TableTypeMaster/pkg/tabletypes"
)

func main() {
	fmt.Println("Exporting to HTML...")
	err := tabletypes.ExportToHTML("tabletypes.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Exported tabletypes.html")

	fmt.Println("Exporting to SQLite...")
	err = tabletypes.ExportToSQLite("tabletypes.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Exported tabletypes.sqlite")
}

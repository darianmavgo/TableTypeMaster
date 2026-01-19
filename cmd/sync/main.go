package main

import (
	"fmt"
	"log"
	"os"
	"tabletypes/pkg/tabletypes"
)

func main() {
	readmePath := "README.md"
	goDataPath := "pkg/tabletypes/data.go"

	// 1. Read README
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		log.Fatalf("Failed to read README.md: %v", err)
	}

	// 2. Parse README table
	readmeTable, err := tabletypes.ParseReadme(string(readmeContent))
	if err != nil {
		log.Fatalf("Failed to parse README table: %v", err)
	}
	fmt.Printf("Parsed %d rows from README.md\n", len(readmeTable))

	// 3. Get Go Data
	goTable := tabletypes.AllTableTypes
	fmt.Printf("Loaded %d rows from data.go\n", len(goTable))

	// 4. Sync
	syncedTable := tabletypes.Sync(readmeTable, goTable)
	fmt.Printf("Synced table has %d rows\n", len(syncedTable))

	// 5. Update README if needed
	// We re-generate the string to compare, or just blindly write.
	// Comparing is better to avoid git noise.
	// However, Sync logic preserves order of list1 (README) then appends list2 (Go).
	// If README was [A, B] and Go was [B, A], Sync(README, Go) -> [A, B].
	// Sync(Go, README) -> [B, A].
	// The requirement is "whichever table is larger catch the other table up".
	// If sizes are different, we definitely sync.
	// If sizes are same but content different? "Catch up" implies strict addition?
	// The Sync function I wrote does a Union.

	// Let's generate the strings and compare.

	// We need to extract the old table string from README to compare, or just overwrite.
	// Overwriting is safer to ensure consistency.
	err = tabletypes.UpdateReadmeFile(readmePath, syncedTable)
	if err != nil {
		log.Fatalf("Failed to update README.md: %v", err)
	}
	fmt.Println("Updated README.md")

	// 6. Update data.go if needed
	// We should compare the syncedTable with goTable.
	// Using reflect.DeepEqual or just generating the content.
	newGoDataContent, err := tabletypes.GenerateGoDataContent(syncedTable)
	if err != nil {
		log.Fatalf("Failed to generate Go data content: %v", err)
	}

	// Read old file to compare (optional, but good)
	oldGoDataBytes, err := os.ReadFile(goDataPath)
	if err == nil {
		if string(oldGoDataBytes) == newGoDataContent {
			fmt.Println("pkg/tabletypes/data.go is already up to date")
		} else {
			if err := os.WriteFile(goDataPath, []byte(newGoDataContent), 0644); err != nil {
				log.Fatalf("Failed to write data.go: %v", err)
			}
			fmt.Println("Updated pkg/tabletypes/data.go")
		}
	} else {
		// File might not exist or error
		if err := os.WriteFile(goDataPath, []byte(newGoDataContent), 0644); err != nil {
			log.Fatalf("Failed to write data.go: %v", err)
		}
		fmt.Println("Created pkg/tabletypes/data.go")
	}

	// Verify sizes match now (conceptually)
	if len(syncedTable) != len(readmeTable) || len(syncedTable) != len(goTable) {
		fmt.Println("Sync completed. Run again to verify convergence if changes were made.")
	}
}

//go:build pure

package sqlite_test

import (
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/drivers/sqlite"
)

func ExampleNew_unimported() {
	db, err := sqlite.New("file::memory:")
	fmt.Println(db)
	fmt.Println(err)

	// Output:
	// <nil>
	// failed to open database: sql: unknown driver "sqlite" (forgotten import?)
}

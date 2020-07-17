package tape

import (
	"fmt"
)

func getDataPath(tableName string) string {
	return fmt.Sprintf("db/%s.Data.husk", tableName)
}

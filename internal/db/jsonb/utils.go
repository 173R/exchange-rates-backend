package jsonb

import (
	"encoding/json"
	"fmt"
)

// ScanTo является shorthand-функцией для конвертации массива байт
// в указанную структуру.
func ScanTo(value any, to any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("некорректное значения для JSON: %s", value)
	}
	return json.Unmarshal(bytes, &to)
}

package functional

import (
	"fmt"
)

// ErrMissingColumn will be returned when requested by the consumer
type ErrMissingColumn struct {
	col string
}

func (emc ErrMissingColumn) Error() string {
	return fmt.Sprintf("[%q] is not a valid column", emc.col)
}

// GroupBy returns a grouped result of an array of a map indexed by string
// this deliberately does not implement
func GroupBy[T comparable](rows []map[string]T, keys ...string) (map[string][]map[string]T, error) {
	result := make(map[string][]map[string]T)
	for _, row := range rows {
		groupName, err := generateGroupNameFromRowData(row, keys...)
		if err != nil {
			return nil, err
		}
		result[groupName] = append(result[groupName], row)
	}
	return result, nil
}

// generateGroupNameFromRowData is a helper function to create group values for GroupBy
func generateGroupNameFromRowData[T comparable](row map[string]T, keys ...string) (string, error) {
	gn := ""
	for _, key := range keys {
		if val, ok := row[key]; ok {
			if gn == "" {
				gn = fmt.Sprintf("%s", val)
				continue
			}
			gn = fmt.Sprintf("%s_%s", gn, val)
		} else {
			return "", ErrMissingColumn{col: key}
		}
	}
	return gn, nil
}

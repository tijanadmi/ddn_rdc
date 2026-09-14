package oraclerepo

import (

	"database/sql"

)

// ScanSafe пресреће NULL вредности за све string показиваче и уписује ""
func ScanSafe(rows *sql.Rows, dest ...interface{}) error {
	scanners := make([]interface{}, len(dest))
	nullStrings := make([]sql.NullString, len(dest))

	for i, d := range dest {
		// Ако је одредиште показивач на string (*string)
		if _, ok := d.(*string); ok {
			scanners[i] = &nullStrings[i]
		} else {
			scanners[i] = d
		}
	}

	// Скенирамо кроз sql.NullString за сва string поља
	if err := rows.Scan(scanners...); err != nil {
		return err
	}

	// Враћамо вредности: ако је био NULL, nullStrings[i].String је већ ""
	for i, d := range dest {
		if strPtr, ok := d.(*string); ok {
			*strPtr = nullStrings[i].String
		}
	}

	return nil
}

// ScanRowSafe пресреће NULL вредности за pojedinaчни *sql.Row
func ScanRowSafe(row *sql.Row, dest ...interface{}) error {
	scanners := make([]interface{}, len(dest))
	nullStrings := make([]sql.NullString, len(dest))

	for i, d := range dest {
		if _, ok := d.(*string); ok {
			scanners[i] = &nullStrings[i]
		} else {
			scanners[i] = d
		}
	}

	// Скенирамо директно из *sql.Row
	if err := row.Scan(scanners...); err != nil {
		return err
	}

	for i, d := range dest {
		if strPtr, ok := d.(*string); ok {
			*strPtr = nullStrings[i].String
		}
	}

	return nil
}
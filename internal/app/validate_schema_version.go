package app

import (
	"fmt"
	"slices"
)

const SCHEMA_VERSION_WANT int32 = 7

func ValidateSchemaVersion(version int32) {
	delta := version - SCHEMA_VERSION_WANT
	if delta != 0 {
		versions := []int32{SCHEMA_VERSION_WANT, version}
		slices.Sort(versions)
		message := "Newer"
		if delta < 0 {
			message = "Old"
		}
		panic(fmt.Sprintf(
			"%s db schema version: want %d, current %d. Need to run migration(s) %d...%d (check also const SCHEMA_VERSION_WANT)",
			message,
			SCHEMA_VERSION_WANT, version,
			versions[0], versions[1],
		))
	}
}

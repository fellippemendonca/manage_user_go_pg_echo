package common

import (
	"github.com/google/uuid"
)

func Paginate[R any](results []R, limit int32, idFromResult func(R) uuid.UUID) ([]R, string) {
	if len(results) > int(limit) {
		last := results[limit]
		return results[:limit], EncodeUUIDToBase64(idFromResult(last))
	}
	return results, ""
}

package common

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
)

func Paginate(results []*models.User, limit int) ([]*models.User, string) {
	if len(results) > int(limit) {
		last := results[limit]
		return results[:limit], EncodeUUIDToBase64(last.ID)
	}
	return results, ""
}

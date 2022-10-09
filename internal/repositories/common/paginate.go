package common

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
)

// Paginate function simply takes a list of Users and removes the last element while transform removed last element ID in a page-token
func Paginate(results []*models.User, limit int) ([]*models.User, string) {
	if len(results) > int(limit) {
		last := results[limit]
		return results[:limit], EncodeUUIDToBase64(last.ID)
	}
	return results, ""
}

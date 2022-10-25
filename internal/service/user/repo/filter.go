package repo

import (
	"strings"

	"github.com/gocomerse/internal/service/user/model"
)

func buildUpdate(user model.User) string {
	var query string
	if user.FirstName != "" {
		query += createQuery("first_name", user.FirstName)
	}
	if user.LastName != "" {
		query += createQuery("last_name", user.LastName)

	}
	if user.Phone != "" {
		query += createQuery("phone", user.Phone)
	}
	query = strings.TrimSuffix(query, ",")
	return query
}

func createQuery(tag, value string) string {
	return tag + `=` + "'" + value + "'" + ","
}

package repo

import (
	"strconv"
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
func createSearchQuery(
	firstName string,
	lastName string,
	email string,
	sort string,
	order string,
	limit int,
	page int,
	pass bool,

) string {

	var arr []string
	var where string
	var orderBy string
	var password string
	if pass {
		password = ", password "
	}

	if firstName != "" {
		arr = append(arr, " first_name LIKE '%"+firstName+"%'")

	}
	if lastName != "" {
		arr = append(arr, " last_name LIKE '%"+lastName+"%'")

	}
	if email != "" {
		arr = append(arr, " email LIKE '%"+email+"%'")

	}
	if len(arr) > 0 {

		where = ` WHERE`

	} else {
		where = ""
	}
	if order == "" {
		order = "ASC"
	}
	if sort != "" {
		orderBy = "ORDER BY " + sort + " " + strings.ToUpper(order)
	}
	if limit == 0 {
		limit = 5
	}
	if page == 0 {
		page = 1
	}
	query := `SELECT 
			 id,
			 first_name,
			 last_name,
			 email,
			 phone ` + password + `
			 FROM "user"` + where + " " + strings.Join(arr, " AND ") + " " + orderBy +
		" OFFSET " + strconv.Itoa((page-1)*limit) + " LIMIT " + strconv.Itoa(limit)
	return query
}

package repo

const (
	getUser     = `SELECT id,first_name,last_name,email,phone from "user"`
	getUserByID = `SELECT id,first_name,last_name,email,phone from "user" where id =$1`
	insertUser  = `Insert into "user" (first_name,last_name,email,phone) Values($1,$2,$3,$4) RETURNING id,first_name,last_name,email,phone`
	deleteUser  = `DELETE FROM "user" where id =$1`
)

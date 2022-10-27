package repo

const (
	getUserByID = `SELECT id,first_name,last_name,email,phone from "user" where id =$1`
	insertUser  = `Insert into "user" 
				   (first_name,last_name,email,phone,password)
				   Values($1,$2,$3,$4,$5) 
				   RETURNING id,first_name,last_name,email,phone
				   `
	deleteUser = `DELETE FROM "user" where id =$1`
)

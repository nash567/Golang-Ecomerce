package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	logModel "github.com/gocomerse/internal/logger/model"
	"github.com/gocomerse/internal/service/user/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}

}

func (r *Repository) Get(ctx context.Context, log logModel.Logger, queryParams model.QueryParams, pass bool) ([]*model.User, error) {
	var Users []*model.User
	query := createSearchQuery(

		queryParams.FirstName,
		queryParams.LastName,
		queryParams.Email,
		queryParams.Sort,
		queryParams.Order,
		queryParams.Limit,
		queryParams.Page,
		pass,
	)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.WithError(err).Error("failed to prepare context with query")
		return nil, fmt.Errorf("failed to prepare context for user: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get users")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	for rows.Next() {
		var user model.User
		if pass {
			err = rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password)
		} else {
			err = rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Phone)

		}
		if err != nil {
			log.WithError(err).Error("failed to get users:")

			return nil, fmt.Errorf("failed to get users: %w", err)
		}
		Users = append(Users, &user)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("failed to get all the users")

		return Users, fmt.Errorf("failed to get all users: %w", err)
	}
	if len(Users) == 0 {
		log.WithError(err).Error("no user found")

		return nil, fmt.Errorf("%w", model.ErrNoRecordFound)
	}
	defer rows.Close()

	return Users, nil
}
func (r *Repository) GetByID(ctx context.Context, log logModel.Logger, id int) (*model.User, error) {
	var User model.User
	stmt, err := r.db.PrepareContext(ctx, getUserByID)
	if err != nil {
		log.WithError(err).Error("failed to prepare context with query")
		return nil, fmt.Errorf("failed to prepare context for user: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(&User.UserID, &User.FirstName, &User.LastName, &User.Email, &User.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.WithError(err).Error("user with this id do not exist")
			return nil, fmt.Errorf("user do not exist: %w", err)
		}
		log.WithError(err).Error("failed to get user by id")
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &User, nil
}

func (r *Repository) Create(ctx context.Context, log logModel.Logger, user model.User) (*model.User, error) {
	var User model.User
	stmt, err := r.db.PrepareContext(ctx, insertUser)
	if err != nil {
		log.WithError(err).Error("failed to prepare context with query")
		return nil, fmt.Errorf("failed to prepare context for Addiying user: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.FirstName, user.LastName, user.Email, user.Phone, user.Password)
	err = row.Scan(&User.UserID, &User.FirstName, &User.LastName, &User.Email, &User.Phone)
	if err != nil {
		log.WithError(err).Error("failed to scan while adding user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User, nil
}

func (r *Repository) Update(ctx context.Context, log logModel.Logger, user model.User) (*model.User, error) {

	var User model.User
	query := `Update "user" set ` + buildUpdate(user) + ` where id=$1 RETURNING id,first_name,last_name,email,phone`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.WithError(err).Error("failed to prepare context for update user")
		return nil, fmt.Errorf("failed to prepare context for update user: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.UserID)
	err = row.Scan(&User.UserID, &User.FirstName, &User.LastName, &User.Email, &User.Phone)
	if err != nil {
		log.WithError(err).Error("failed to scan while adding user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User, nil
}

func (r *Repository) Delete(ctx context.Context, log logModel.Logger, id int) error {

	res, err := r.db.ExecContext(ctx, deleteUser, id)

	if err != nil {
		log.WithError(err).Error("failed to delete user with id")
		return fmt.Errorf("failed to delete user :%w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.WithError(err).Error("failed to delete user with id")
		return fmt.Errorf("failed to delete user :%w", err)
	}

	if count == 1 {
		return nil
	} else if count >= 0 {
		return fmt.Errorf("error occurred: %w", err)
	}
	return nil
}

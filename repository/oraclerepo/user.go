package oraclerepo

import (
	"context"
	"errors"
	"time"

	"github.com/tijanadmi/ddn_rdc/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate authenticates a user
func (m *OracleDBRepo) Authenticate(username, testPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	/*var id int
	var hashedPassword string*/

	var user models.User

	query := `select id, username, password from tis_services_users where username = :1`

	row := m.DB.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("incorrect password")
	} else if err != nil {
		return err
	}
	return nil
}

func (m *OracleDBRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, username, password from tis_services_users where username = :1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}
	query = `select ru.id,RU.ID_USER, RU.ID_ROLE, R.CODE, R.NAME
	from tis_services_role_user ru, tis_services_roles r
	where RU.ID_USER =:1
	and ru.id_role = r.id
	`

	var roles []models.UserRole
	rows, _ := m.DB.QueryContext(ctx, query, user.ID)
	defer rows.Close()

	for rows.Next() {
		var r models.UserRole
		err := rows.Scan(
			&r.ID,
			&r.IdUser,
			&r.IdRole,
			&r.RoleCode,
			&r.RoleName,
		)

		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	user.UserRole = roles

	return &user, nil
}

func (m *OracleDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, username, password from tis_services_users where id = :1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	query = `select ru.id,RU.ID_USER, RU.ID_ROLE, R.CODE, R.NAME
	from tis_services_role_user ru, tis_services_roles r
	where RU.ID_USER =:1
	and ru.id_role = r.id
	`

	var roles []models.UserRole
	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	for rows.Next() {
		var r models.UserRole
		err := rows.Scan(
			&r.ID,
			&r.IdUser,
			&r.IdRole,
			&r.RoleCode,
			&r.RoleName,
		)

		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	user.UserRole = roles

	return &user, nil
}
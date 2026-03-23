package oraclerepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/tijanadmi/ddn_rdc/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate authenticates a user
func (m *OracleDBRepo) Authenticate(ctx context.Context, username, testPassword string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	/*var id int
	var hashedPassword string*/

	var user models.User

	query := `select id, username, password, full_name from tis_services_users where username = :1`

	row := m.DB.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.FullName)
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

func (m *OracleDBRepo) getUserOptimized(ctx context.Context, field string, value any) (*models.User, error) {
	// 1. Dohvati osnovne podatke o korisniku
	query := fmt.Sprintf(`select id, username, password, RTRIM(full_name)
						  from tis_services_users
						  where %s = :1`, field)

	var user models.User
	var fullName sql.NullString
	row := m.DB.QueryRowContext(ctx, query, value)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &fullName); err != nil {
		return nil, err
	}
	if fullName.Valid {
		user.FullName = strings.TrimSpace(fullName.String)
	}

	// 2. Dohvati role korisnika
	roleQuery := `select RU.ID_ROLE, R.CODE, R.NAME
				  from tis_services_role_user ru
				  join tis_services_roles r on ru.id_role = r.id
				  where ru.ID_USER = :1`

	rows, err := m.DB.QueryContext(ctx, roleQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	var ddnRoles, tdnRoles, pgiRoles []int // id_role liste po tipu
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Code, &r.Name); err != nil {
			return nil, err
		}

		roles = append(roles, r)

		switch r.Code {
		case "DDN":
			ddnRoles = append(ddnRoles, r.ID)
		case "TDN":
			tdnRoles = append(tdnRoles, r.ID)
		case "PGI":
			pgiRoles = append(pgiRoles, r.ID)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 3. Dohvati sve DDN podatke u jednom query-ju
	if len(ddnRoles) > 0 {
		var ddnData models.DDNData
		err := m.DB.QueryRowContext(ctx,
			`select tip_priv_prip, id_s_mrc
			 from ddn_kor
			 where sifra_dd = UPPER(:1)`, user.Username,
		).Scan(&ddnData.TipPrivPrip, &ddnData.IdSMrc)
		if err == nil {
			for i := range roles {
				if roles[i].Code == "DDN" {
					roles[i].DDN = ddnData
				}
			}
		}
	}

	// 4. Dohvati sve TDN podatke u jednom query-ju
	if len(tdnRoles) > 0 {
		var tdnData models.TDNData
		err := m.DB.QueryRowContext(ctx,
			`select tip_priv_prip, id_p0_traf
			 from tdn_kor
			 where sifra_dd = UPPER(:1)`, user.Username,
		).Scan(&tdnData.TipPrivPrip, &tdnData.IdP0Traf)
		if err == nil {
			for i := range roles {
				if roles[i].Code == "TDN" {
					roles[i].TDN = tdnData
				}
			}
		}
	}

	// 5. Dohvati sve PGI podatke u jednom query-ju
	if len(pgiRoles) > 0 {
		var pgiData models.PGIData
		err := m.DB.QueryRowContext(ctx,
			`select TIP_PRIV_PRIP, TIP_IZV_P, TIP_IZV_D, TIP_IZV_M,
			        ID_S_KOM, TIP_AKC, ID_S_MRC,
			        T1,T2,T3,T4,T5,T6,T7,T8,T9,T10
			 from pgi_kor
			 where sifra_pi = UPPER(:1)`, user.Username,
		).Scan(
			&pgiData.TipPrivPrip, &pgiData.TipIzvP, &pgiData.TipIzvD, &pgiData.TipIzvM,
			&pgiData.IdSKom, &pgiData.TipAkc, &pgiData.IdSMrc,
			&pgiData.T1, &pgiData.T2, &pgiData.T3, &pgiData.T4, &pgiData.T5,
			&pgiData.T6, &pgiData.T7, &pgiData.T8, &pgiData.T9, &pgiData.T10,
		)
		if err == nil {
			for i := range roles {
				if roles[i].Code == "PGI" {
					roles[i].PGI = pgiData
				}
			}
		}
	}

	user.Roles = roles

	// fmt.Printf("User: %+v\n", user)
	return &user, nil
}

func (m *OracleDBRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return m.getUserOptimized(ctx, "id", id)
}

func (m *OracleDBRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.getUserOptimized(ctx, "username", username)
}

// func (m *OracleDBRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	// defer cancel()

// 	query := `select id,
// 					username,
// 					password,
// 					RTRIM(full_name)
// 					from tis_services_users where username = :1`

// 	var user models.User
// 	var fullName sql.NullString
// 	row := m.DB.QueryRowContext(ctx, query, username)

// 	err := row.Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Password,
// 		&fullName,
// 	)

// 	if fullName.Valid {
// 		user.FullName = strings.TrimSpace(fullName.String)
// 	} else {
// 		user.FullName = ""
// 	}

// 	if err != nil {

// 		return nil, err
// 	}
// 	query = `select RU.ID_ROLE, R.CODE, R.NAME
// 	from tis_services_role_user ru, tis_services_roles r
// 	where RU.ID_USER =:1
// 	and ru.id_role = r.id
// 	`

// 	// var roles []string
// 	var roles []models.Role

// 	rows, err := m.DB.QueryContext(ctx, query, user.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var r models.Role
// 		err := rows.Scan(
// 			&r.ID,
// 			&r.Code,
// 			&r.Name,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		/***** dodatak za DDN, TDN, PGI *****/
// 		switch r.Code {

// 		case "DDN":
// 			data := models.DDNData{}
// 			err := m.DB.QueryRowContext(ctx,
// 				`select tip_priv_prip, id_s_mrc
// 			 from ddn_kor
// 			 where sifra_dd = :1`,
// 				user.Username,
// 			).Scan(&data.TipPrivPrip, &data.IdSMrc)

// 			if err == nil {
// 				r.DDN = &data
// 			}

// 		case "TDN":
// 			data := models.TDNData{}
// 			err := m.DB.QueryRowContext(ctx,
// 				`select tip_priv_prip, id_p0_traf
// 			 from tdn_kor
// 			 where sifra_dd = :1`,
// 				user.Username,
// 			).Scan(&data.TipPrivPrip, &data.IdP0Traf)

// 			if err == nil {
// 				r.TDN = &data
// 			}

// 		case "PGI":
// 			data := models.PGIData{}
// 			err := m.DB.QueryRowContext(ctx,
// 				`select TIP_PRIV_PRIP, TIP_IZV_P, TIP_IZV_D, TIP_IZV_M,
// 			        ID_S_KOM, TIP_AKC, ID_S_MRC,
// 			        T1,T2,T3,T4,T5,T6,T7,T8,T9,T10
// 			 from pgi_kor
// 			 where sifra_pi = :1`,
// 				user.Username,
// 			).Scan(
// 				&data.TipPrivPrip,
// 				&data.TipIzvP,
// 				&data.TipIzvD,
// 				&data.TipIzvM,
// 				&data.IdSKom,
// 				&data.TipAkc,
// 				&data.IdSMrc,
// 				&data.T1, &data.T2, &data.T3, &data.T4, &data.T5,
// 				&data.T6, &data.T7, &data.T8, &data.T9, &data.T10,
// 			)

// 			if err == nil {
// 				r.PGI = &data
// 			}
// 		}
// 		/****** kraj dodatka *****/
// 		// roles = append(roles, r.Code)
// 		roles = append(roles, r)
// 	}
// 	user.Roles = roles

// 	fmt.Printf("User: %+v\n", user)

// 	return &user, nil
// }

// func (m *OracleDBRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	// defer cancel()

// 	query := `select id, username, password, full_name from tis_services_users where id = :1`

// 	var user models.User
// 	row := m.DB.QueryRowContext(ctx, query, id)

// 	err := row.Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Password,
// 		&user.FullName,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	query = `select  R.code
// 	from tis_services_role_user ru, tis_services_roles r
// 	where RU.ID_USER =:1
// 	and ru.id_role = r.id
// 	`

// 	var roles []string
// 	rows, _ := m.DB.QueryContext(ctx, query, id)
// 	defer rows.Close()

// 	for rows.Next() {
// 		var r models.Role
// 		err := rows.Scan(
// 			&r.Code,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}
// 		roles = append(roles, r.Code)
// 	}
// 	user.Role = roles

// 	return &user, nil
// }

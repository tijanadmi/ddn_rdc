package oraclerepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) getUserOptimized(ctx context.Context, field string, value any) (*models.User, error) {
	// 1. Dohvati osnovne podatke o korisniku
	query := fmt.Sprintf(`select id, ad_username,  RTRIM(full_name)
						  from tis_users
						  where %s = :1`, field)

	var user models.User
	var fullName sql.NullString
	row := m.DB.QueryRowContext(ctx, query, value)
	if err := ScanRowSafe(row, &user.ID, &user.Username, &fullName); err != nil {
		return nil, err
	}
	if fullName.Valid {
		user.FullName = strings.TrimSpace(fullName.String)
	}

	// 2. Dohvati role korisnika
	roleQuery := `select RU.ID_ROLE, R.CODE, R.NAME
				  from tis_role_user ru
				  join tis_roles r on ru.id_role = r.id
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
		if err := ScanSafe(rows, &r.ID, &r.Code, &r.Name); err != nil {
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

		row := m.DB.QueryRowContext(ctx,
			`select tip_priv_prip, id_s_mrc
			 from ddn_kor
			 where sifra_dd = UPPER(:1)`, user.Username,
		)
		err := ScanRowSafe(row, &ddnData.TipPrivPrip, &ddnData.IdSMrc)
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
		row := m.DB.QueryRowContext(ctx,
			`select tip_priv_prip, id_p0_traf
			 from tdn_kor
			 where sifra_dd = UPPER(:1)`, user.Username,
		)

		err := ScanRowSafe(row, &tdnData.TipPrivPrip, &tdnData.IdP0Traf)
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
		row := m.DB.QueryRowContext(ctx,
			`select TIP_PRIV_PRIP, TIP_IZV_P, TIP_IZV_D, TIP_IZV_M,
			        ID_S_KOM, TIP_AKC, ID_S_MRC,
			        T1,T2,T3,T4,T5,T6,T7,T8,T9,T10
			 from pgi_kor
			 where sifra_pi = UPPER(:1)`, user.Username,
		)
		err := ScanRowSafe(row, &pgiData.TipPrivPrip, &pgiData.TipIzvP, &pgiData.TipIzvD, &pgiData.TipIzvM,
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
	return m.getUserOptimized(ctx, "ad_username", username)
}

package oraclerepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetPiPIT4ByParams(ctx context.Context, arg models.ListPiDDT4Params) ([]*models.PiMMT4, int, error) {

	var query string
	var rows *sql.Rows
	var err error

	if strings.ToUpper(arg.IdSMrc) == "ALL" {
		query = `SELECT ROWNUM AS id,
			CASE
  				WHEN kom1 = 1 OR kom2 = 1 OR kom3 = 1 OR kom4 = 1 OR
       				kom5 = 1 OR kom6 = 1 OR kom7 = 1 OR kom8 = 1 THEN '1'
  				WHEN kom1 = 0 OR kom2 = 0 OR kom3 = 0 OR kom4 = 0 OR
       				kom5 = 0 OR kom6 = 0 OR kom7 = 0 OR kom8 = 0 THEN '0'
  				ELSE ''
			END AS stav,
			to_char(DATIZV,'dd.mm.yyyy'),
			to_char(id1),
			mrc,
			tekst,
			COALESCE(kom1, ''),
			COALESCE(kom2, ''),
			COALESCE(kom3, ''),
			COALESCE(kom4, ''),
			COALESCE(kom5, ''),
			COALESCE(kom6, ''),
			COALESCE(kom7, ''),
			COALESCE(kom8, ''),
			COALESCE(opist4, ''),
			COUNT(*) OVER () AS TOTAL_COUNT
		FROM PI_T4_MM_V
		WHERE DATIZV = to_date(:1,'dd.mm.yyyy') 
		ORDER BY datizv, id_s_tipd, id1`

		rows, err = m.DB.QueryContext(ctx, query, arg.Datizv)
	} else {
		query = `SELECT ROWNUM AS id,
			'' AS stav,
			to_char(DATIZV,'dd.mm.yyyy'),
			to_char(id1),
			ted.td_nazivi.td_daj_sif ('S_MRC', 'NAZIV', 'ID', PI_DD.ID_S_MRC, 'Q') as mrc,
			tac4 AS tekst,
			'' AS kom1,
			'' AS kom2,
			'' AS kom3,
			'' AS kom4,
			'' AS kom5,
			'' AS kom6,
			'' AS kom7,
			'' AS kom8,
			'' AS opist4,
			COUNT(*) OVER () AS TOTAL_COUNT
		FROM PI_PI
		WHERE DATIZV = to_date(:1,'dd.mm.yyyy') 
			AND id_s_tipd = 4
			AND id_s_mrc = :2
		ORDER BY datizv, id_s_tipd, id1`

		rows, err = m.DB.QueryContext(ctx, query, arg.Datizv, arg.IdSMrc)
	}

	if err != nil {
		fmt.Println("Greška prilikom izvršavanja upita:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.PiMMT4
	var totalCount int

	for rows.Next() {
		var ue models.PiMMT4
		var count int
		err := rows.Scan(
			&ue.Id,
			&ue.Stav,
			&ue.Datizv,
			&ue.Id1,
			&ue.Mrc,
			&ue.Tekst,
			&ue.Kom1,
			&ue.Kom2,
			&ue.Kom3,
			&ue.Kom4,
			&ue.Kom5,
			&ue.Kom6,
			&ue.Kom7,
			&ue.Kom8,
			&ue.Opist4,
			&count,
		)

		if err != nil {
			fmt.Println("Greška pri čitanju reda:", err)
			return nil, 0, err
		}

		ues = append(ues, &ue)
		totalCount = count
	}

	return ues, totalCount, nil
}

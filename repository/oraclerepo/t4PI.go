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
            '' AS stav,
            to_char(DATIZV,'dd.mm.yyyy'),
            to_char(id1),
            ted.td_nazivi.td_daj_sif ('S_MRC', 'NAZIV', 'ID', PI_PI.ID_S_MRC, 'Q') as mrc,
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
		ORDER BY datizv, id_s_tipd, id1`

		rows, err = m.DB.QueryContext(ctx, query, arg.Datizv)
	} else {
		query = `SELECT ROWNUM AS id,
            '' AS stav,
            to_char(DATIZV,'dd.mm.yyyy'),
            to_char(id1),
            ted.td_nazivi.td_daj_sif ('S_MRC', 'NAZIV', 'ID', PI_PI.ID_S_MRC, 'Q') as mrc,
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

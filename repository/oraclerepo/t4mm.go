package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetPiMMT4ByParams(ctx context.Context, arg models.ListPiMMT4Params) ([]*models.PiMMT4, int, error) {

	query := `SELECT ROWNUM AS id,
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
      where DATIZV  BETWEEN to_date(:1,'dd.mm.yyyy') AND to_date(:2,'dd.mm.yyyy') 
      order by datizv,id_s_tipd,id1
	  OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.StartDate, arg.EndDate, arg.Offset, arg.Limit)
	//fmt.Println("Upit za PiMM T4:", query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
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
			return nil, 0, err
		}

		ues = append(ues, &ue)
		totalCount = count
	}

	return ues, totalCount, nil
}

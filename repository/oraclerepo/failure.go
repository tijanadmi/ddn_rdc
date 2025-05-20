package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetPiMMByParams(ctx context.Context, arg models.ListPiMMParams) ([]*models.PiMM, int, error) {

	query := `select ROWNUM AS id,
					TIPD,
					to_char(DATIZV,'dd.mm.yyyy'),
					COALESCE(to_char(ID1), ''),
					COALESCE(to_char(ID2), ''),
   					VREPOC_EXCEL,
   					VREZAV_EXCEL,
   					TRAJ_EXCEL,
   					OB_ID,
    				TIPOB,
   					OB_SIF,
   					NAZOB,
					COALESCE(to_char(VRPD), ''),
   					NAZVRPD,
   					GRUZR1 ||'/'||
   					UZROK1,
					COALESCE(VREM_USL, ''),
      				TEKST_EXCEL,
					COALESCE(to_char(SNAGA), ''),
					COUNT(*) OVER () AS TOTAL_COUNT  
					from pgi.pi_mm_v
					where DATIZV  BETWEEN to_date(:1,'dd.mm.yyyy') AND to_date(:2,'dd.mm.yyyy') 
   					AND TIPD LIKE UPPER(:3)
   					and kom2=1
   					order by DATIZV,id1,id2`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.StartDate, arg.EndDate, arg.Tipd)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.PiMM
	var totalCount int

	for rows.Next() {
		var ue models.PiMM
		var count int
		err := rows.Scan(
			&ue.Id,
			&ue.Tipd,
			&ue.Datizv,
			&ue.Id1,
			&ue.Id2,
			&ue.Vrepoc,
			&ue.Vrezav,
			&ue.Traj,
			&ue.ObId,
			&ue.TipOb,
			&ue.ObSif,
			&ue.NazOb,
			&ue.Vrpd,
			&ue.Nazvrpd,
			&ue.Uzrok,
			&ue.VrmUsl,
			&ue.Tekst,
			&ue.Snaga,
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

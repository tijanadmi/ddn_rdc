package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

/****  pomocna funkcija ***/
func sameDay(a, b string) bool {
	return a == b
}

/*** Funkcija koja izvlaci podatke za mesecni izvestaj PI MM po grupama***/
func (m *OracleDBRepo) GetPiMMReportByParams(ctx context.Context, arg models.ListPiMMByParam) (*models.Report, error) {

	query := `
SELECT
    PI_DDD.ID_S_TIPD        AS TIPD,
    ted.td_nazivi.td_daj_sif('S_TIPD','NAZIV','ID',PI_DDD.ID_S_TIPD,'Q') AS NAZTIPD,
    to_char(PI_DDD.DATIZV,'dd.mm.yyyy') AS DATIZV,
    COALESCE(to_char(PI_DDD.ID1), '') ID1,
    PI_DDD.TEKST,
    to_char(PI_DD.VREPOC,'dd.mm.yyyy HH24:MI') AS VREPOC,
    to_char(PI_DD.VREZAV,'dd.mm.yyyy HH24:MI') AS VREZAV,
    PI_TRAJANJE.VREME_DHM(PI_DD.VREZAV, PI_DD.VREPOC, PI_DDD.DATIZV) AS TRAJ,
    UPPER(
        CASE
            WHEN V_S_OB.TIPOB IN (8,9,4,10)
                THEN V_S_OB.OB_SIF || ' ' || V_S_OB.NAZOB
            ELSE V_S_OB.OPIS
        END
    ) AS OBJEKAT, 
    COALESCE(ted.TD_NAZIVI.TD_DAJ_SIF('V_S_POLJE_SVA_AP','POLJE','P2_TRAF_ID',PI_DD.P2_TRAF_ID,'Q'), '') POLJE, 
    COALESCE(ted.TD_NAZIVI.TD_DAJ_SIF('V_S_POLJE_SVA_AP','IME_PO','P2_TRAF_ID',PI_DD.P2_TRAF_ID,'Q'), '') IME_PO, 
    COALESCE(to_char(PI_DD.SNAGA), '') SNAGA,
    ted.td_nazivi.td_daj_sif('S_VRPD','NAZIV','ID',PI_DD.ID_S_VRPD,'Q') AS VRSTA_DOG,
    ted.TD_NAZIVI.TD_DAJ_SIF('S_GRUZR','NAZIV','ID',PI_DD.ID1_S_GRUZR,'Q') AS GRUZR,
    ted.TD_NAZIVI.TD_DAJ_SIF('S_UZROK','NAZIV','ID',PI_DD.ID1_S_UZROK,'Q') AS UZROK
FROM PI_DDD
JOIN PI_DD
  ON PI_DDD.DATIZV = PI_DD.DATIZV
 AND PI_DDD.ID_S_TIPD = PI_DD.ID_S_TIPD
 AND PI_DDD.ID1 = PI_DD.ID1
LEFT JOIN V_S_OB
  ON PI_DD.OB_ID = V_S_OB.OB_ID
 AND PI_DD.ID_TIPOB = V_S_OB.TIPOB
WHERE PI_DDD.DATIZV BETWEEN TO_DATE(:1,'dd.mm.yyyy') AND TO_DATE(:2,'dd.mm.yyyy')
  AND PI_DDD.ID_S_TIPD LIKE :3
  AND (
        (:4 = '0' AND (
            PI_DDD.KOM1='1' OR PI_DDD.KOM2='1' OR PI_DDD.KOM3='1' OR PI_DDD.KOM4='1' OR
            PI_DDD.KOM5='1' OR PI_DDD.KOM6='1' OR PI_DDD.KOM7='1' OR PI_DDD.KOM8='1'
        ))
        OR (:5 <> '0' AND
            DECODE(:6,
                '1', PI_DDD.KOM1,
                '2', PI_DDD.KOM2,
                '3', PI_DDD.KOM3,
                '4', PI_DDD.KOM4,
                '5', PI_DDD.KOM5,
                '6', PI_DDD.KOM6,
                '7', PI_DDD.KOM7,
                '8', PI_DDD.KOM8
            ) = '1'
        )
  )
ORDER BY TIPD, PI_DDD.DATIZV, PI_DDD.ID1, PI_DD.VREPOC
`

	tipdParam := "%"
	if arg.Tipd != 0 {
		tipdParam = fmt.Sprintf("%d", arg.Tipd)
	}

	rows, err := m.DB.QueryContext(
		ctx,
		query,
		arg.StartDate,
		arg.EndDate,
		tipdParam,
		arg.Kom,
		arg.Kom,
		arg.Kom,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report models.Report

	var currentTipd *models.TipdGroup
	var currentDay *models.DayGroup
	var currentEvent *models.EventGroup

	for rows.Next() {

		var (
			tipd    string
			naztipd string
			datizv  string
			id1     int
			tekst   string
			row     models.DetailRow
		)

		err := rows.Scan(
			&tipd,
			&naztipd,
			&datizv,
			&id1,
			&tekst,
			&row.Vrepoc,
			&row.Vrezav,
			&row.Traj,
			&row.Objekat,
			&row.Polje,
			&row.ImePolja,
			&row.Snaga,
			&row.VrstaDogadjaja,
			&row.GrupaUzroka,
			&row.Uzrok,
		)
		if err != nil {
			return nil, err
		}

		// TIPD GROUP
		if currentTipd == nil || currentTipd.Tipd != tipd {
			report.TipdGroups = append(report.TipdGroups, models.TipdGroup{
				Tipd:  tipd,
				Naziv: naztipd,
			})
			currentTipd = &report.TipdGroups[len(report.TipdGroups)-1]
			currentDay = nil
			currentEvent = nil
		}

		// DAY GROUP
		if currentDay == nil || !sameDay(currentDay.Date, datizv) {
			currentTipd.Days = append(currentTipd.Days, models.DayGroup{
				Date: datizv,
			})
			currentDay = &currentTipd.Days[len(currentTipd.Days)-1]
			currentEvent = nil
		}

		// EVENT GROUP (ID1)
		if currentEvent == nil || currentEvent.ID1 != id1 {
			currentDay.Events = append(currentDay.Events, models.EventGroup{
				ID1:   id1,
				Tekst: tekst,
			})
			currentEvent = &currentDay.Events[len(currentDay.Events)-1]
		}

		currentEvent.Rows = append(currentEvent.Rows, row)
	}

	return &report, nil
}

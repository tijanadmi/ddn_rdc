package oraclerepo

import (
	"context"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetShemeByOrg(ctx context.Context, idOrg int, tipDok string) ([]models.Shema, error) {

	query := `
		SELECT
			ID,
			IME_DOK,
			PUTANJA,
			TIP_DOK,
			ID_S_ORG,
			TO_CHAR(DATPRI,'dd.mm.yyyy') DATPRI
		FROM TED.TIS_DOK
		WHERE ID_S_ORG = :1
		  AND TIP_DOK = :2
		ORDER BY IME_DOK
	`

	rows, err := m.DB.QueryContext(ctx, query, idOrg, tipDok)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sheme []models.Shema

	for rows.Next() {

		var s models.Shema

		err := rows.Scan(
			&s.ID,
			&s.ImeDok,
			&s.Putanja,
			&s.TipDok,
			&s.IdSOrg,
			&s.Datpri,
		)

		if err != nil {
			return nil, err
		}

		sheme = append(sheme, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sheme, nil
}

func (m *OracleDBRepo) GetShemaPutanjaByID(ctx context.Context, id int) (string, string, error) {

	query := `
		SELECT PUTANJA, IME_DOK
		FROM TED.TIS_DOK
		WHERE ID = :1
	`

	var putanja string
	var imeDok string

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&putanja, &imeDok)
	if err != nil {
		return "", "", err
	}

	return putanja, imeDok, nil
}

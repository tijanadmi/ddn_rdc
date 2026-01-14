package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetMrcById(ctx context.Context, id int) (*models.SMrc, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status, naziv_cir
			  from s_mrc
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var mrc models.SMrc

	err := row.Scan(
		&mrc.ID,
		&mrc.Code,
		&mrc.Name,
		&mrc.Status,
		&mrc.NameCir,
	)

	if err != nil {
		return nil, err
	}

	return &mrc, err
}

// type ListLimitOffsetParams struct {
// 	Limit  int32  `json:"limit"`
// 	Offset int32  `json:"offset"`
// }

// Get returns all s_mrc and error, if any
func (m *OracleDBRepo) GetSMrc(ctx context.Context) ([]*models.SMrc, error) {
	query := `select id, sifra, naziv, status, naziv_cir
			  from s_mrc
			  where id not in (5,7)
			  ORDER BY id desc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.SMrc

	for rows.Next() {
		var mrc models.SMrc
		if err := rows.Scan(
			&mrc.ID,
			&mrc.Code,
			&mrc.Name,
			&mrc.Status,
			&mrc.NameCir,
		); err != nil {
			return nil, err
		}
		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}

// Get returns  s_mrc for Lov for create/update forms and error, if any
func (m *OracleDBRepo) GetSMrcForInsert(ctx context.Context) ([]*models.SMrc, error) {
	query := `select id, sifra, naziv, status, naziv_cir
			  from s_mrc
			  where id not in (5,7,9)
			  ORDER BY id desc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.SMrc

	for rows.Next() {
		var mrc models.SMrc
		if err := rows.Scan(
			&mrc.ID,
			&mrc.Code,
			&mrc.Name,
			&mrc.Status,
			&mrc.NameCir,
		); err != nil {
			return nil, err
		}
		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}

func (m *OracleDBRepo) GetSTipPrekById(ctx context.Context, id int) (*models.STipPrek, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status
			  from s_tip_prek
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var tip models.STipPrek

	err := row.Scan(
		&tip.ID,
		&tip.Code,
		&tip.Name,
		&tip.Status,
	)

	if err != nil {
		return nil, err
	}

	return &tip, err
}

// Get returns all s_tip_prek and error, if any
func (m *OracleDBRepo) GetSTipPrek(ctx context.Context) ([]*models.STipPrek, error) {
	query := `select id, sifra, naziv, status
			  from s_tip_prek
			  ORDER BY id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tips []*models.STipPrek

	for rows.Next() {
		var tip models.STipPrek
		if err := rows.Scan(
			&tip.ID,
			&tip.Code,
			&tip.Name,
			&tip.Status,
		); err != nil {
			return nil, err
		}
		tips = append(tips, &tip)
	}

	return tips, nil
}

func (m *OracleDBRepo) GetSVrPrekById(ctx context.Context, id int) (*models.SVrPrek, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status
			  from s_vr_prek
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var vr models.SVrPrek

	err := row.Scan(
		&vr.ID,
		&vr.Code,
		&vr.Name,
		&vr.Status,
	)

	if err != nil {
		return nil, err
	}

	return &vr, err
}

// Get returns all s_vr_prek and error, if any
func (m *OracleDBRepo) GetSVrPrek(ctx context.Context) ([]*models.SVrPrek, error) {
	query := `select id, sifra, naziv, status
			  from s_vr_prek
			  ORDER BY id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vrps []*models.SVrPrek

	for rows.Next() {
		var vr models.SVrPrek
		if err := rows.Scan(
			&vr.ID,
			&vr.Code,
			&vr.Name,
			&vr.Status,
		); err != nil {
			return nil, err
		}
		vrps = append(vrps, &vr)
	}

	return vrps, nil
}

// Get returns all S_VRSTA_PREKIDA_P_GEN_V and error, if any
func (m *OracleDBRepo) GetSPodVrPrek(ctx context.Context) ([]*models.SPodVrPrek, error) {
	query := `select   a.ID_TIP_OBJEKTA,
   			  a.OPIS,
   			  a.ID_TIP_DOGADJAJA, a.id_s_vr_prek, B.NAZIV
			  from S_VRSTA_PREKIDA_P_GEN_V a, s_vr_prek b
			  where a.id_s_vr_prek=b.id
			  and pasiva is null`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vrps []*models.SPodVrPrek

	for rows.Next() {
		var vr models.SPodVrPrek
		if err := rows.Scan(
			&vr.IdTipObjekta,
			&vr.Opis,
			&vr.IdTipDogadjaja,
			&vr.IdSVrPrek,
			&vr.Naziv,
		); err != nil {
			return nil, err
		}
		vrps = append(vrps, &vr)
	}

	return vrps, nil
}

func (m *OracleDBRepo) GetSUzrokPrekById(ctx context.Context, id int) (*models.SUzrokPrek, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status
			  from s_uzrok_prek
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var uzrok models.SUzrokPrek

	err := row.Scan(
		&uzrok.ID,
		&uzrok.Code,
		&uzrok.Name,
		&uzrok.Status,
	)

	if err != nil {
		return nil, err
	}

	return &uzrok, err
}

// Get returns all s_uzrok_prek and error, if any
func (m *OracleDBRepo) GetSUzrokPrek(ctx context.Context) ([]*models.SUzrokPrek, error) {
	query := `select id, sifra, naziv, status
			  from s_uzrok_prek
			  ORDER BY id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uzroks []*models.SUzrokPrek

	for rows.Next() {
		var uzrok models.SUzrokPrek
		if err := rows.Scan(
			&uzrok.ID,
			&uzrok.Code,
			&uzrok.Name,
			&uzrok.Status,
		); err != nil {
			return nil, err
		}
		uzroks = append(uzroks, &uzrok)
	}

	return uzroks, nil
}

func (m *OracleDBRepo) GetSPoduzrokPrekById(ctx context.Context, id int) (*models.SPoduzrokPrek, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status
			  from s_poduzrok_prek
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var poduzrok models.SPoduzrokPrek

	err := row.Scan(
		&poduzrok.ID,
		&poduzrok.Code,
		&poduzrok.Name,
		&poduzrok.Status,
	)

	if err != nil {
		return nil, err
	}

	return &poduzrok, err
}

// Get returns all s_poduzrok_prek and error, if any
func (m *OracleDBRepo) GetSPoduzrokPrek(ctx context.Context) ([]*models.SPoduzrokPrek, error) {
	query := `select id, sifra, naziv, status
			  from s_poduzrok_prek
			  ORDER BY id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var poduzroks []*models.SPoduzrokPrek

	for rows.Next() {
		var poduzrok models.SPoduzrokPrek
		if err := rows.Scan(
			&poduzrok.ID,
			&poduzrok.Code,
			&poduzrok.Name,
			&poduzrok.Status,
		); err != nil {
			return nil, err
		}
		poduzroks = append(poduzroks, &poduzrok)
	}

	return poduzroks, nil
}

func (m *OracleDBRepo) GetSMernaMestaById(ctx context.Context, id int) (*models.SMernaMesta, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select id, sifra,naziv, status
			  from s_merna_mesta
			  where id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var mm models.SMernaMesta

	err := row.Scan(
		&mm.ID,
		&mm.Code,
		&mm.Name,
		&mm.Status,
	)

	if err != nil {
		return nil, err
	}

	return &mm, err
}

// Get returns all s_poduzrok_prek and error, if any
func (m *OracleDBRepo) GetSMernaMesta(ctx context.Context) ([]*models.SMernaMesta, error) {
	query := `select id, sifra, naziv, status
			  from s_merna_mesta
			  ORDER BY id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mms []*models.SMernaMesta

	for rows.Next() {
		var mm models.SMernaMesta
		if err := rows.Scan(
			&mm.ID,
			&mm.Code,
			&mm.Name,
			&mm.Status,
		); err != nil {
			return nil, err
		}
		mms = append(mms, &mm)
	}

	return mms, nil
}

func (m *OracleDBRepo) GetObjById(ctx context.Context, id int) (*models.ObjLOV, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select OPIS skr,
                id_s_mrc1 id_s_mrc1,
                ID,
                ID_S_TIPOB TIPOB,
                ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',ID_S_TIPOB,'Q') SIF_TIPOB,
                SIFRA
                from P0_TRAF
                WHERE ID=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var mrc models.ObjLOV

	err := row.Scan(
		&mrc.Opis,
		&mrc.IdSMrc1,
		&mrc.ObId,
		&mrc.Tipob,
		&mrc.SifTipob,
		&mrc.ObSif,
	)

	if err != nil {
		return nil, err
	}

	return &mrc, err
}

func (m *OracleDBRepo) GetPoljeGEById(ctx context.Context, id int) (*models.PoljaLOV, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select P2_TRAF_ID,
       			polje,
       			IME_PO,
       			S_NAP.ID NN_ID,
       			s_nap.naziv||' '||s_nap.jedinica nn_naziv
 				from V_s_polje_SVA,s_nap,S_FUP
				where NN_ID = s_nap.id
				AND FUP_ID=S_FUP.ID
				AND S_FUP.SIFRA='20'
                AND P2_TRAF_ID=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var mrc models.PoljaLOV

	err := row.Scan(
		&mrc.Id,
		&mrc.Polje,
		&mrc.PoljeNaziv,
		&mrc.NNId,
		&mrc.NNNaziv,
	)

	if err != nil {
		return nil, err
	}

	return &mrc, err
}

// type ListLimitOffsetParams struct {
// 	Limit  int32  `json:"limit"`
// 	Offset int32  `json:"offset"`
// }

// Get returns all s_mrc and error, if any
func (m *OracleDBRepo) GetObjTSRP(ctx context.Context, arg models.ListObjectLimitOffsetParams) ([]*models.ObjLOV, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select OPIS skr,
				id_s_mrc1 id_s_mrc1,
				OB_ID,
				TIPOB,
				ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',TIPOB,'Q') SIF_TIPOB,
				OB_SIF
				from V_s_ob
				where   ( :1=8 OR (id_s_mrc1 = :2  or NVL(id_s_mrc2,0) = :3 )  )
				and upper(status) = 'A'
				AND ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',TIPOB,'Q')  IN ('TS','TT','RP','TP')
			  ORDER BY OPIS
			  OFFSET :4 ROWS FETCH NEXT :5 ROWS ONLY`

	rows, err := m.DB.QueryContext(ctx, query, arg.Mrc, arg.Mrc, arg.Mrc, arg.Offset, arg.Limit)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.ObjLOV

	for rows.Next() {
		var mrc models.ObjLOV
		err := rows.Scan(
			&mrc.Opis,
			&mrc.IdSMrc1,
			&mrc.ObId,
			&mrc.Tipob,
			&mrc.SifTipob,
			&mrc.ObSif,
		)

		if err != nil {
			return nil, err
		}

		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}

// Get returns all s_mrc and error, if any
func (m *OracleDBRepo) GetObjHETEVE(ctx context.Context, arg models.ListObjectParams) ([]*models.ObjLOV, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select OPIS skr,
				id_s_mrc1 id_s_mrc1,
				OB_ID,
				TIPOB,
				ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',TIPOB,'Q') SIF_TIPOB,
				OB_SIF
				from V_s_ob
				where   ( :1=8 OR (id_s_mrc1 = :2  or NVL(id_s_mrc2,0) = :3 )  )
				and upper(status) = 'A'
				AND  SUBSTR(OPIS,1,2) IN ('HE','TE','VE','RH')
			  ORDER BY OPIS`

	rows, err := m.DB.QueryContext(ctx, query, arg.Mrc, arg.Mrc, arg.Mrc)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.ObjLOV

	for rows.Next() {
		var mrc models.ObjLOV
		err := rows.Scan(
			&mrc.Opis,
			&mrc.IdSMrc1,
			&mrc.ObId,
			&mrc.Tipob,
			&mrc.SifTipob,
			&mrc.ObSif,
		)

		if err != nil {
			return nil, err
		}

		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}

func (m *OracleDBRepo) GetPoljaGE(ctx context.Context, arg models.ListPoljaLimitOffsetParams) ([]*models.PoljaLOV, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `select P2_TRAF_ID,
       			polje,
       			IME_PO,
       			S_NAP.ID NN_ID,
       			s_nap.naziv||' '||s_nap.jedinica nn_naziv
 				from V_s_polje_SVA,s_nap,S_FUP
				where ob_id = :1
				and NN_ID = s_nap.id
				AND FUP_ID=S_FUP.ID
				AND S_FUP.SIFRA='20'
				AND upper(V_S_POLJE_SVA.STATUS) = 'A'
			  	ORDER BY polje`

	rows, err := m.DB.QueryContext(ctx, query, arg.ObjId)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.PoljaLOV

	for rows.Next() {
		var mrc models.PoljaLOV
		err := rows.Scan(
			&mrc.Id,
			&mrc.Polje,
			&mrc.PoljeNaziv,
			&mrc.NNId,
			&mrc.NNNaziv,
		)

		if err != nil {
			return nil, err
		}

		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}

package oraclerepo

import (
	"context"
	"fmt"
	"time"

	"github.com/tijanadmi/ddn_rdc/models"
)




func (m *OracleDBRepo) DeleteDDNInterruptionOfDelivery(Id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from ddn_prekid_isp where id = :1`

	_, err := m.DB.ExecContext(ctx, stmt, Id)
	if err != nil {
		return err
	}

	return nil
}


// func (m *OracleDBRepo) GetDDNInterruptionOfDeliveryById(id int) (*models.DDNInterruptionOfDelivery, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	query := `SELECT PI.ID,
// 			PI.ID_S_MRC,
//   			MR.naziv,
//   			PI.ID_S_TIPD,
//   			PI.ID_S_VRPD,
//   			O.TIPOB,
//   			O.OB_ID,
//   			O.OB_SIF,
//   			O.NAZOB,
//   			O.OPIS,
//   			O.SKR,
//   			O.ID_S_MRC1,
//   			O.MRC1,
//   			O.ID_S_MRC2,
//   			O.MRC2,
//   			O.ID_S_ORG1,
//   			O.ID_S_ORG2,
//   			PI.VREPOC,
//   			PI.VREZAV,
//   			PI.ID_S_VR_PREK,
//   			V.NAZIV,
//   			PI.ID_S_UZROK_PREK,
//   			U.NAZIV,
//   			PI.SNAGA,
//   			PI.OPIS,
//   			PI.DDN_KOR,
//   			PI.ID_DOG_SMENE,
//   			PI.ID_STAVKE,
//   			PI.MOD,
//   			PI.ID_S_MERNA_MESTA,
//   			M.NAZIV,
//   			PI.BROJ_MMESTA,
//   			PI.IND,
//   			PI.ID_P2_TRAF,
//   			PO.IME_PO,
//   			PO.OPIS,
//   			PI.BI,
//   			PI.ID_S_PODUZROK_PREK,
//   			PI.ID_DOG_PREKID_P,
//   			PI.ID_TIP_OBJEKTA_NDC,
//   			PI.ID_TIP_DOGADJAJA_NDC,
// 			PI.SYNSOFT_ID
//   			FROM DDN_PREKID_ISP PI
//   			INNER JOIN  S_MRC MR ON PI.ID_S_MRC=MR.ID
//   			INNER JOIN  V_S_OB O ON PI.OB_ID=O.OB_ID AND PI.ID_TIPOB=O.TIPOB
//   			LEFT JOIN S_VR_PREK V ON PI.ID_S_VR_PREK=V.ID
//   			LEFT JOIN S_UZROK_PREK U ON PI.ID_S_UZROK_PREK=U.ID
//   			LEFT JOIN S_MERNA_MESTA M ON PI.ID_S_MERNA_MESTA=M.ID
//   			LEFT JOIN V_S_POLJE_SVA_AP PO ON PI.ID_P2_TRAF=PO.P2_TRAF_ID
//   			where id=:1`

// 	row, err := m.DB.QueryContext(ctx, query,id)
	
	
// 	var i models.DDNInterruptionOfDelivery
// 	err := rows.Scan(
// 			&i.Id,
// 			&i.IdSMrc,
// 			&i.SMrc.Name,
// 			&i.IdSTipd,
// 			&i.IdSVrpd,
// 			&i.VSOb.Tipob,
// 			&i.VSOb.ObId,
// 			&i.VSOb.ObSif,
// 			&i.VSOb.Nazob,
// 			&i.VSOb.Opis,
// 			&i.VSOb.Skt,
// 			&i.VSOb.IdSMrc1,
// 	        &i.VSOb.Mrc1,
// 			&i.VSOb.IdSMrc2,
// 			&i.VSOb.Mrc2,
// 			&i.VSOb.IdSOrg1,
// 			&i.VSOb.IdSOrg2,
// 			&i.Vrepoc,
// 			&i.Vrezav,
// 			&i.SVrPrek.ID,
// 			&i.SVrPrek.Name,
// 			&i.SUzrokPrek.ID,
// 			&i.SUzrokPrek.Name,
// 			&i.Snaga,
// 			&i.Opis,
// 			&i.KorUneo,
// 			&i.IdDogSmene,
// 			&i.IdStavke,
// 			&i.Mod,
// 			&i.IdSMernaMesta,
// 			&i.SMernaMesta.Name,
// 			&i.BrojMesta,
// 			&i.Ind,
// 			&i.P2TrafId,
// 			&i.VSPoljeSvaAP.ImePo,
// 			&i.VSPoljeSvaAP.Opis,
// 			&i.Bi,
// 			&i.IdSPoduzrokPrek,
// 			&i.IdDogPrekidP,
// 			&i.IdTipObjektaNdc,
// 			&i.IdTipDogadjajaNdc,
// 			&i.SynsoftId,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 	return &i, nil
// }



// func (m *OracleDBRepo) GetDDNInterruptionOfDeliveryById(id int) ([]*models.DDNInterruptionOfDelivery, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	query := `select ID_S_MRC,
// 	COALESCE(to_char(ID_S_TIPD), ''),
// 	COALESCE(to_char(ID_S_VRPD), ''),
// 	COALESCE(to_char(ID_TIPOB), ''),
// 	COALESCE(to_char(OB_ID), ''),
//    to_char(VREPOC, 'dd.mm.yyyy HH24:MI:SS'),
//    to_char(VREZAV, 'dd.mm.yyyy HH24:MI:SS'),
//    COALESCE(to_char(ID_S_VR_PREK), ''),
//    COALESCE(to_char(ID_S_UZROK_PREK), ''),
//    COALESCE(to_char(SNAGA), ''),
//    COALESCE(OPIS, ''),
//    COALESCE(DDN_KOR, ''),
//    COALESCE(to_char(ID_S_MERNA_MESTA), ''),
//    COALESCE(to_char(BROJ_MMESTA), ''),
//    COALESCE(IND, ''),
//    COALESCE(to_char(ID_P2_TRAF), ''),
//    COALESCE(to_char(BI), ''),
//    COALESCE(to_char(ID_S_PODUZROK_PREK), ''),
//    COALESCE(to_char(ID_DOG_PREKID_P), ''),
//    COALESCE(to_char(ID_TIP_OBJEKTA_NDC), ''),
//    COALESCE(to_char(ID_TIP_DOGADJAJA_NDC), ''),
//    COALESCE(SYNSOFT_ID, '')
//    from ddn_prekid_isp
//    where id=:1`

// 	rows, err := m.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		fmt.Println("Pogresan upit ili nema rezultata upita")
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var p []*models.DDNInterruptionOfDelivery

// 	for rows.Next() {
// 		var ue models.DDNInterruptionOfDelivery
// 		err := rows.Scan(
// 			&ue.IdSMrc,
// 			&ue.IdSTipd,
// 			&ue.IdSVrpd,
// 			&ue.IdTipob,
// 			&ue.ObId,
// 			&ue.Vrepoc,
// 			&ue.Vrezav,
// 			&ue.IdSVrPrek,
// 			&ue.IdSUzrokPrek,
// 			&ue.Snaga,
// 			&ue.Opis,
// 			&ue.KorUneo,
// 			&ue.IdSMernaMesta,
// 			&ue.BrojMesta,
// 			&ue.Ind,
// 			&ue.P2TrafId,
// 			&ue.Bi,
// 			&ue.IdSPoduzrokPrek,
// 			&ue.IdDogPrekidP,
// 			&ue.IdTipObjektaNdc,
// 			&ue.IdTipDogadjajaNdc,
// 			&ue.SynsoftId,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		p = append(p, &ue)
// 	}

// 	return p, nil
// }


func (m *OracleDBRepo) GetMrcById(id int) (*models.SMrc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

// Get returns all s_mrc and error, if any
func (m *OracleDBRepo) GetSMrc() ([]*models.SMrc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status, naziv_cir
			  from s_mrc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var mrcs []*models.SMrc

	for rows.Next() {
		var mrc models.SMrc
		err := rows.Scan(
			&mrc.ID,
			&mrc.Code,
			&mrc.Name,
			&mrc.Status,
			&mrc.NameCir,
		)

		if err != nil {
			return nil, err
		}

		mrcs = append(mrcs, &mrc)
	}

	return mrcs, nil
}


func (m *OracleDBRepo) GetSTipPrekById(id int) (*models.STipPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
func (m *OracleDBRepo) GetSTipPrek() ([]*models.STipPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status
			  from s_tip_prek`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var tips []*models.STipPrek

	for rows.Next() {
		var tip models.STipPrek
		err := rows.Scan(
			&tip.ID,
			&tip.Code,
			&tip.Name,
			&tip.Status,
		)

		if err != nil {
			return nil, err
		}

		tips = append(tips, &tip)
	}

	return tips, nil
}


func (m *OracleDBRepo) GetSVrPrekById(id int) (*models.SVrPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
func (m *OracleDBRepo) GetSVrPrek() ([]*models.SVrPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status
			  from s_vr_prek`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var vrps []*models.SVrPrek

	for rows.Next() {
		var vr models.SVrPrek
		err := rows.Scan(
			&vr.ID,
			&vr.Code,
			&vr.Name,
			&vr.Status,
		)

		if err != nil {
			return nil, err
		}

		vrps = append(vrps, &vr)
	}

	return vrps, nil
}

func (m *OracleDBRepo) GetSUzrokPrekById(id int) (*models.SUzrokPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
func (m *OracleDBRepo) GetSUzrokPrek() ([]*models.SUzrokPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status
			  from s_uzrok_prek`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var uzroks []*models.SUzrokPrek

	for rows.Next() {
		var uzrok models.SUzrokPrek
		err := rows.Scan(
			&uzrok.ID,
			&uzrok.Code,
			&uzrok.Name,
			&uzrok.Status,
		)

		if err != nil {
			return nil, err
		}

		uzroks = append(uzroks, &uzrok)
	}

	return uzroks, nil
}

func (m *OracleDBRepo) GetSPoduzrokPrekById(id int) (*models.SPoduzrokPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
func (m *OracleDBRepo) GetSPoduzrokPrek() ([]*models.SPoduzrokPrek, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status
			  from s_poduzrok_prek`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var poduzroks []*models.SPoduzrokPrek

	for rows.Next() {
		var poduzrok models.SPoduzrokPrek
		err := rows.Scan(
			&poduzrok.ID,
			&poduzrok.Code,
			&poduzrok.Name,
			&poduzrok.Status,
		)

		if err != nil {
			return nil, err
		}

		poduzroks = append(poduzroks, &poduzrok)
	}

	return poduzroks, nil
}

func (m *OracleDBRepo) GetSMernaMestaById(id int) (*models.SMernaMesta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
func (m *OracleDBRepo) GetSMernaMesta() ([]*models.SMernaMesta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sifra,naziv, status
			  from s_merna_mesta`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var mms []*models.SMernaMesta

	for rows.Next() {
		var mm models.SMernaMesta
		err := rows.Scan(
			&mm.ID,
			&mm.Code,
			&mm.Name,
			&mm.Status,
		)

		if err != nil {
			return nil, err
		}

		mms = append(mms, &mm)
	}

	return mms, nil
}
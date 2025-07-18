package oraclerepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) DeleteDDNInterruptionOfDelivery(ctx context.Context, Id string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	stmt := `delete from ddn_prekid_isp where id = :1`

	_, err := m.DB.ExecContext(ctx, stmt, Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *OracleDBRepo) GetDDNInterruptionOfDeliveryById(ctx context.Context, id int) (*models.DDNInterruptionOfDelivery, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	query := `select PI.ID,PI.ID_S_MRC,MR.naziv,
    COALESCE(to_char(PI.ID_S_TIPD), ''),
    COALESCE(to_char(PI.ID_S_VRPD), ''),
    COALESCE(to_char(PI.ID_TIPOB), ''),
    COALESCE(to_char(PI.OB_ID), ''),
    COALESCE(to_char(O.NAZOB), ''),
    COALESCE(to_char(O.OPIS), ''),
   to_char(PI.VREPOC,'dd.mm.yyyy HH24:MI:SS'),
   to_char(PI.VREZAV,'dd.mm.yyyy HH24:MI:SS'),
   case
   when PI.VREZAV is not null
   then
   pgi.PI_TRAJANJE.VREME_DHM(PI.VREZAV,PI.VREPOC,PI.VREZAV)
   else
   pgi.PI_TRAJANJE.VREME_DHM(sysdate,PI.VREPOC,PI.VREZAV)
   end
   as trajanje,
   COALESCE(to_char(PI.ID_S_VR_PREK), ''),
    COALESCE(to_char(V.NAZIV), ''),
    COALESCE(to_char(vp.OPIS), ''), 
   COALESCE(to_char(PI.ID_S_UZROK_PREK), ''),
    COALESCE(to_char(U.NAZIV), ''),
   COALESCE(to_char(PI.SNAGA), ''),
   COALESCE(PI.OPIS, ''),
   COALESCE(PI.DDN_KOR, ''),
   COALESCE(to_char(PI.ID_S_MERNA_MESTA), ''),
   COALESCE(to_char(M.NAZIV), ''),
   COALESCE(to_char(PI.BROJ_MMESTA), ''),
   COALESCE(PI.IND, ''),
   COALESCE(to_char(PI.ID_P2_TRAF), ''),
   COALESCE(to_char(PO.IME_PO), ''),
   COALESCE(to_char(PO.OPIS), ''),
   COALESCE(to_char(PI.BI), ''),
   COALESCE(to_char(PI.ID_S_PODUZROK_PREK), ''),
   COALESCE(to_char(PU.NAZIV), ''),
   COALESCE(to_char(PI.ID_DOG_PREKID_P), ''),
   COALESCE(to_char(PI.ID_TIP_OBJEKTA_NDC), ''),
   COALESCE(to_char(PI.ID_TIP_DOGADJAJA_NDC), ''),
   COALESCE(PI.SYNSOFT_ID, ''),
   COALESCE(PI.version, 0)
   from ddn_prekid_isp PI
   INNER JOIN  S_MRC MR ON PI.ID_S_MRC=MR.ID
   INNER JOIN  V_S_OB O ON PI.OB_ID=O.OB_ID AND PI.ID_TIPOB=O.TIPOB
   LEFT JOIN S_VR_PREK V ON PI.ID_S_VR_PREK=V.ID
   LEFT JOIN S_UZROK_PREK U ON PI.ID_S_UZROK_PREK=U.ID
   LEFT JOIN S_PODUZROK_PREK PU ON PI.ID_S_PODUZROK_PREK=PU.ID
   LEFT JOIN S_MERNA_MESTA M ON PI.ID_S_MERNA_MESTA=M.ID
   LEFT JOIN V_S_POLJE_SVA_AP PO ON PI.ID_P2_TRAF=PO.P2_TRAF_ID
   LEFT JOIN S_VRSTA_PREKIDA_P_GEN_V vp ON PI.ID_TIP_OBJEKTA_NDC=vp.ID_TIP_OBJEKTA AND PI.ID_TIP_DOGADJAJA_NDC=vp.ID_TIP_DOGADJAJA AND PI.ID_S_VR_PREK=vp.ID_S_VR_PREK
	  WHERE PI.id=:1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var ue models.DDNInterruptionOfDelivery
	err := row.Scan(
		&ue.Id,
		&ue.IdSMrc,
		&ue.Mrc,
		&ue.IdSTipd,
		&ue.IdSVrpd,
		&ue.IdTipob,
		&ue.ObId,
		&ue.ObNaziv,
		&ue.ObOpis,
		&ue.Vrepoc,
		&ue.Vrezav,
		&ue.Trajanje,
		&ue.IdSVrPrek,
		&ue.VrstaPrek,
		&ue.PodvrstaPrek,
		&ue.IdSUzrokPrek,
		&ue.Uzrok,
		&ue.Snaga,
		&ue.Opis,
		&ue.KorUneo,
		&ue.IdSMernaMesta,
		&ue.MernaMesta,
		&ue.BrojMesta,
		&ue.Ind,
		&ue.P2TrafId,
		&ue.PoljeNaziv,
		&ue.PoljeOpis,
		&ue.Bi,
		&ue.IdSPoduzrokPrek,
		&ue.PoduzrokPrek,
		&ue.IdDogPrekidP,
		&ue.IdTipObjektaNdc,
		&ue.IdTipDogadjajaNdc,
		&ue.SynsoftId,
		&ue.Version,
	)
	if err != nil {
		return nil, err
	}

	return &ue, nil
}

func (m *OracleDBRepo) GetDDNInterruptionOfDeliveryByPage(ctx context.Context, arg models.ListInterruptionWithPaginationParams) ([]*models.DDNInterruptionOfDelivery, int, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	mrcParam := "%"
	if strings.ToUpper(arg.Mrc) != "ALL" {
		mrcParam = strings.ToUpper(arg.Mrc)
	}

	query := `select PI.ID,PI.ID_S_MRC,MR.naziv,
    COALESCE(to_char(PI.ID_S_TIPD), ''),
    COALESCE(to_char(PI.ID_S_VRPD), ''),
    COALESCE(to_char(PI.ID_TIPOB), ''),
    COALESCE(to_char(PI.OB_ID), ''),
    COALESCE(to_char(O.NAZOB), ''),
    COALESCE(to_char(O.OPIS), ''),
   to_char(PI.VREPOC,'dd.mm.yyyy HH24:MI:SS'),
   to_char(PI.VREZAV,'dd.mm.yyyy HH24:MI:SS'),
   case
   when PI.VREZAV is not null
   then
   pgi.PI_TRAJANJE.VREME_DHM(PI.VREZAV,PI.VREPOC,PI.VREZAV)
   else
   pgi.PI_TRAJANJE.VREME_DHM(sysdate,PI.VREPOC,PI.VREZAV)
   end
   as trajanje,
   COALESCE(to_char(PI.ID_S_VR_PREK), ''),
    COALESCE(to_char(V.NAZIV), ''),
    COALESCE(to_char(vp.OPIS), ''), 
   COALESCE(to_char(PI.ID_S_UZROK_PREK), ''),
    COALESCE(to_char(U.NAZIV), ''),
   COALESCE(to_char(PI.SNAGA), ''),
   COALESCE(PI.OPIS, ''),
   COALESCE(PI.DDN_KOR, ''),
   COALESCE(to_char(PI.ID_S_MERNA_MESTA), ''),
   COALESCE(to_char(M.NAZIV), ''),
   COALESCE(to_char(PI.BROJ_MMESTA), ''),
   COALESCE(PI.IND, ''),
   COALESCE(to_char(PI.ID_P2_TRAF), ''),
   COALESCE(to_char(PO.IME_PO), ''),
   COALESCE(to_char(PO.OPIS), ''),
   COALESCE(to_char(PI.BI), ''),
   COALESCE(to_char(PI.ID_S_PODUZROK_PREK), ''),
   COALESCE(to_char(PU.NAZIV), ''),
   COALESCE(to_char(PI.ID_DOG_PREKID_P), ''),
   COALESCE(to_char(PI.ID_TIP_OBJEKTA_NDC), ''),
   COALESCE(to_char(PI.ID_TIP_DOGADJAJA_NDC), ''),
   COALESCE(PI.SYNSOFT_ID, ''),
   COALESCE(PI.version, 0),
   COUNT(*) OVER () AS TOTAL_COUNT
   from ddn_prekid_isp PI
   INNER JOIN  S_MRC MR ON PI.ID_S_MRC=MR.ID
   INNER JOIN  V_S_OB O ON PI.OB_ID=O.OB_ID AND PI.ID_TIPOB=O.TIPOB
   LEFT JOIN S_VR_PREK V ON PI.ID_S_VR_PREK=V.ID
   LEFT JOIN S_UZROK_PREK U ON PI.ID_S_UZROK_PREK=U.ID
   LEFT JOIN S_PODUZROK_PREK PU ON PI.ID_S_PODUZROK_PREK=PU.ID
   LEFT JOIN S_MERNA_MESTA M ON PI.ID_S_MERNA_MESTA=M.ID
   LEFT JOIN V_S_POLJE_SVA_AP PO ON PI.ID_P2_TRAF=PO.P2_TRAF_ID
   LEFT JOIN S_VRSTA_PREKIDA_P_GEN_V vp ON PI.ID_TIP_OBJEKTA_NDC=vp.ID_TIP_OBJEKTA AND PI.ID_TIP_DOGADJAJA_NDC=vp.ID_TIP_DOGADJAJA AND PI.ID_S_VR_PREK=vp.ID_S_VR_PREK
	  WHERE PI.IND=:1   AND PI.ID_S_MRC like (:2) AND  
	(PI.VREPOC >= to_date(:3,'dd.mm.yyyy HH24:MI:SS') AND PI.VREPOC<= to_date(:4,'dd.mm.yyyy HH24:MI:SS'))
   ORDER BY id
			  OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.Ind, mrcParam, arg.StartDate, arg.EndDate, arg.Offset, arg.Limit)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.DDNInterruptionOfDelivery
	var totalCount int

	for rows.Next() {
		var ue models.DDNInterruptionOfDelivery
		var count int

		err := rows.Scan(
			&ue.Id,
			&ue.IdSMrc,
			&ue.Mrc,
			&ue.IdSTipd,
			&ue.IdSVrpd,
			&ue.IdTipob,
			&ue.ObId,
			&ue.ObNaziv,
			&ue.ObOpis,
			&ue.Vrepoc,
			&ue.Vrezav,
			&ue.Trajanje,
			&ue.IdSVrPrek,
			&ue.VrstaPrek,
			&ue.PodvrstaPrek,
			&ue.IdSUzrokPrek,
			&ue.Uzrok,
			&ue.Snaga,
			&ue.Opis,
			&ue.KorUneo,
			&ue.IdSMernaMesta,
			&ue.MernaMesta,
			&ue.BrojMesta,
			&ue.Ind,
			&ue.P2TrafId,
			&ue.PoljeNaziv,
			&ue.PoljeOpis,
			&ue.Bi,
			&ue.IdSPoduzrokPrek,
			&ue.PoduzrokPrek,
			&ue.IdDogPrekidP,
			&ue.IdTipObjektaNdc,
			&ue.IdTipDogadjajaNdc,
			&ue.SynsoftId,
			&ue.Version,
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

func (m *OracleDBRepo) GetAllDDNInterruptionOfDelivery(ctx context.Context, arg models.ListInterruptionParams) ([]*models.DDNInterruptionOfDelivery, int, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	mrcParam := "%"
	if strings.ToUpper(arg.Mrc) != "ALL" {
		mrcParam = strings.ToUpper(arg.Mrc)
	}

	query := `select PI.ID,PI.ID_S_MRC,MR.naziv,
    COALESCE(to_char(PI.ID_S_TIPD), ''),
    COALESCE(to_char(PI.ID_S_VRPD), ''),
    COALESCE(to_char(PI.ID_TIPOB), ''),
    COALESCE(to_char(PI.OB_ID), ''),
    COALESCE(to_char(O.NAZOB), ''),
    COALESCE(to_char(O.OPIS), ''),
   to_char(PI.VREPOC,'dd.mm.yyyy HH24:MI:SS'),
   to_char(PI.VREZAV,'dd.mm.yyyy HH24:MI:SS'),
   case
   when PI.VREZAV is not null
   then
   pgi.PI_TRAJANJE.VREME_DHM(PI.VREZAV,PI.VREPOC,PI.VREZAV)
   else
   pgi.PI_TRAJANJE.VREME_DHM(sysdate,PI.VREPOC,PI.VREZAV)
   end
   as trajanje,
   COALESCE(to_char(PI.ID_S_VR_PREK), ''),
    COALESCE(to_char(V.NAZIV), ''),
    COALESCE(to_char(vp.OPIS), ''), 
   COALESCE(to_char(PI.ID_S_UZROK_PREK), ''),
    COALESCE(to_char(U.NAZIV), ''),
   COALESCE(to_char(PI.SNAGA), ''),
   COALESCE(PI.OPIS, ''),
   COALESCE(PI.DDN_KOR, ''),
   COALESCE(to_char(PI.ID_S_MERNA_MESTA), ''),
   COALESCE(to_char(M.NAZIV), ''),
   COALESCE(to_char(PI.BROJ_MMESTA), ''),
   COALESCE(PI.IND, ''),
   COALESCE(to_char(PI.ID_P2_TRAF), ''),
   COALESCE(to_char(PO.IME_PO), ''),
   COALESCE(to_char(PO.OPIS), ''),
   COALESCE(to_char(PI.BI), ''),
   COALESCE(to_char(PI.ID_S_PODUZROK_PREK), ''),
   COALESCE(to_char(PU.NAZIV), ''),
   COALESCE(to_char(PI.ID_DOG_PREKID_P), ''),
   COALESCE(to_char(PI.ID_TIP_OBJEKTA_NDC), ''),
   COALESCE(to_char(PI.ID_TIP_DOGADJAJA_NDC), ''),
   COALESCE(PI.SYNSOFT_ID, ''),
   COALESCE(PI.version, 0),
   COUNT(*) OVER () AS TOTAL_COUNT
   from ddn_prekid_isp PI
   INNER JOIN  S_MRC MR ON PI.ID_S_MRC=MR.ID
   INNER JOIN  V_S_OB O ON PI.OB_ID=O.OB_ID AND PI.ID_TIPOB=O.TIPOB
   LEFT JOIN S_VR_PREK V ON PI.ID_S_VR_PREK=V.ID
   LEFT JOIN S_UZROK_PREK U ON PI.ID_S_UZROK_PREK=U.ID
   LEFT JOIN S_PODUZROK_PREK PU ON PI.ID_S_PODUZROK_PREK=PU.ID
   LEFT JOIN S_MERNA_MESTA M ON PI.ID_S_MERNA_MESTA=M.ID
   LEFT JOIN V_S_POLJE_SVA_AP PO ON PI.ID_P2_TRAF=PO.P2_TRAF_ID
   LEFT JOIN S_VRSTA_PREKIDA_P_GEN_V vp ON PI.ID_TIP_OBJEKTA_NDC=vp.ID_TIP_OBJEKTA AND PI.ID_TIP_DOGADJAJA_NDC=vp.ID_TIP_DOGADJAJA AND PI.ID_S_VR_PREK=vp.ID_S_VR_PREK
	  WHERE PI.IND=:1   AND PI.ID_S_MRC like (:2) AND  
	(PI.VREPOC >= to_date(:3,'dd.mm.yyyy HH24:MI:SS') AND PI.VREPOC<= to_date(:4,'dd.mm.yyyy HH24:MI:SS'))
   ORDER BY id`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.Ind, mrcParam, arg.StartDate, arg.EndDate)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.DDNInterruptionOfDelivery
	var totalCount int

	for rows.Next() {
		var ue models.DDNInterruptionOfDelivery
		var count int

		err := rows.Scan(
			&ue.Id,
			&ue.IdSMrc,
			&ue.Mrc,
			&ue.IdSTipd,
			&ue.IdSVrpd,
			&ue.IdTipob,
			&ue.ObId,
			&ue.ObNaziv,
			&ue.ObOpis,
			&ue.Vrepoc,
			&ue.Vrezav,
			&ue.Trajanje,
			&ue.IdSVrPrek,
			&ue.VrstaPrek,
			&ue.PodvrstaPrek,
			&ue.IdSUzrokPrek,
			&ue.Uzrok,
			&ue.Snaga,
			&ue.Opis,
			&ue.KorUneo,
			&ue.IdSMernaMesta,
			&ue.MernaMesta,
			&ue.BrojMesta,
			&ue.Ind,
			&ue.P2TrafId,
			&ue.PoljeNaziv,
			&ue.PoljeOpis,
			&ue.Bi,
			&ue.IdSPoduzrokPrek,
			&ue.PoduzrokPrek,
			&ue.IdDogPrekidP,
			&ue.IdTipObjektaNdc,
			&ue.IdTipDogadjajaNdc,
			&ue.SynsoftId,
			&ue.Version,
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

// func (m *OracleDBRepo) GetDDNInterruptionOfDeliveryById(ctx context.Context,id int) (*models.DDNInterruptionOfDelivery, error) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	// defer cancel()

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

// 	row := m.DB.QueryRowContext(ctx, query,id)

// 	var  vrPrekName, uzrokPrekName, sMernaMestaName, imePo, opisPo sql.NullString
// 	var idSMernaMesta sql.NullInt64

// 	var i models.DDNInterruptionOfDelivery
// 	err := row.Scan(
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
// 			&vrPrekName,
// 			&i.SUzrokPrek.ID,
// 			&uzrokPrekName,
// 			&i.Snaga,
// 			&i.Opis,
// 			&i.KorUneo,
// 			&i.IdDogSmene,
// 			&i.IdStavke,
// 			&i.Mod,
// 			&idSMernaMesta,
// 			&sMernaMestaName,
// 			&i.BrojMesta,
// 			&i.Ind,
// 			&i.P2TrafId,
// 			&imePo,
// 			&opisPo,
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

// 		// Assign the null strings to the struct fields

// 	if vrPrekName.Valid {
// 		if i.SVrPrek == nil {
// 			i.SVrPrek = &models.SVrPrek{}
// 		}
// 		i.SVrPrek.Name = vrPrekName.String
// 	}
// 	if uzrokPrekName.Valid {
// 		if i.SUzrokPrek == nil {
// 			i.SUzrokPrek = &models.SUzrokPrek{}
// 		}
// 		i.SUzrokPrek.Name = uzrokPrekName.String
// 	}
// 	if sMernaMestaName.Valid {
// 		if i.SMernaMesta == nil {
// 			i.SMernaMesta = &models.SMernaMesta{}
// 		}
// 		i.SMernaMesta.Name = sMernaMestaName.String
// 	}
// 	// if i.IdSMernaMesta.Valid {
//     //     i.IdSMernaMesta = i.IdSMernaMesta.Int64
//     // } else {
//     //     i.IdSMernaMesta = 0
//     // }
// 	if idSMernaMesta.Valid {
// 		i.IdSMernaMesta = idSMernaMesta
// 	} else {
// 		i.IdSMernaMesta = sql.NullInt64{Int64: 0, Valid: false}
// 	}
// 	if imePo.Valid {
// 		if i.VSPoljeSvaAP == nil {
// 			i.VSPoljeSvaAP = &models.VSPoljeSvaAP{}
// 		}
// 		i.VSPoljeSvaAP.ImePo = imePo.String
// 	}
// 	if opisPo.Valid {
// 		if i.VSPoljeSvaAP == nil {
// 			i.VSPoljeSvaAP = &models.VSPoljeSvaAP{}
// 		}
// 		i.VSPoljeSvaAP.Opis = opisPo.String
// 	}

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

//		return p, nil
//	}

func (m *OracleDBRepo) InsertDDNInterruptionOfDeliveryP(ctx context.Context, ddnintd models.CreateDDNInterruptionOfDeliveryPParams) (int, error) {

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	query := `INSERT INTO DDN_PREKID_ISP (
		ID_S_MRC,
		ID_S_TIPD,
		ID_TIPOB,
		OB_ID,
		VREPOC,
		VREZAV,
		ID_S_VR_PREK,
		ID_S_UZROK_PREK,
		SNAGA,
		OPIS,
		DDN_KOR,
		DATPRI,
		IND,
		ID_P2_TRAF,
		ID_S_PODUZROK_PREK,
		VERSION
	) VALUES (
		:1, :2, :3, :4,
		to_date(:5, 'dd.mm.yyyy HH24:MI:SS'),
		to_date(:6, 'dd.mm.yyyy HH24:MI:SS'),
		:7, :8, :9, :10, :11,
		SYSDATE,
		:12, :13, :14, 0
	) RETURNING id INTO :15`

	PIdSTipd := 12
	PInd := "P"

	var id int

	_, err := m.DB.ExecContext(ctx, query,
		ddnintd.IdSMrc,
		PIdSTipd,
		ddnintd.IdTipob,
		ddnintd.ObId,
		ddnintd.Vrepoc,
		ddnintd.Vrezav,
		ddnintd.IdSVrPrek,
		ddnintd.IdSUzrokPrek,
		ddnintd.Snaga,
		ddnintd.Opis,
		ddnintd.KorUneo,
		PInd,
		ddnintd.P2TrafId,
		ddnintd.IdSPoduzrokPrek,
		sql.Out{Dest: &id},
	)

	//fmt.Println(query)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *OracleDBRepo) UpdateDDNInterruptionOfDeliveryP(ctx context.Context, id int, version int, ddnintd models.CreateDDNInterruptionOfDeliveryPParams) error {
	query := `UPDATE DDN_PREKID_ISP SET
		ID_S_MRC = $1,
		ID_S_TIPD = $2,
		ID_TIPOB = $3,
		OB_ID = $4,
		VREPOC = $5,
		VREZAV = $6,
		ID_S_VR_PREK = $7,
		ID_S_UZROK_PREK = $8,
		SNAGA = $9,
		OPIS = $10,
		DDN_KOR = $11,
		IND = $12,
		ID_P2_TRAF = $13,
		ID_S_PODUZROK_PREK = $14,
		VERSION = VERSION + 1
	WHERE id = $15 AND VERSION = $16`

	PIdSTipd := 12
	PInd := "P"

	result, err := m.DB.ExecContext(ctx, query,
		ddnintd.IdSMrc,
		PIdSTipd,
		ddnintd.IdTipob,
		ddnintd.ObId,
		ddnintd.Vrepoc,
		ddnintd.Vrezav,
		ddnintd.IdSVrPrek,
		ddnintd.IdSUzrokPrek,
		ddnintd.Snaga,
		ddnintd.Opis,
		ddnintd.KorUneo,
		PInd,
		ddnintd.P2TrafId,
		ddnintd.IdSPoduzrokPrek,
		id,
		version,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("optimistic lock failed: object may have been updated by another transaction")
	}

	return nil
}

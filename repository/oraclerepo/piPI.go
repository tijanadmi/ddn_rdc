package oraclerepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetPiPIByParams(ctx context.Context, arg models.ListPiDDParams) ([]*models.PiDD, int, error) {

	mrcParam := "%"
	if strings.ToUpper(arg.IdSMrc) != "ALL" {
		mrcParam = strings.ToUpper(arg.IdSMrc)
	}

	fupParam := "%"
	if strings.ToUpper(arg.Fup) != "ALL" {
		fupParam = strings.ToUpper(arg.Fup)
	}

	query := `select ROWNUM AS id,
                	id_s_mrc,
                	mrc,
                    TIPD,
                    to_char(DATIZV,'dd.mm.yyyy'),
                    COALESCE(to_char(ID1), ''),
                    COALESCE(to_char(ID2), ''),
                    VREPOC_EXCEL,
                    TO_CHAR (vrepoc, 'hh24:mi') vrepoc_hhmi_excel,
                    COALESCE(POC_PP, '') as POC_PP,
                    VREZAV_EXCEL,
                    TO_CHAR (vrezav, 'hh24:mi') vrezav_hhmi_excel,
                    COALESCE(ZAV_PP, '') as ZAV_PP,
                    TRAJ_EXCEL,
                    OB_ID,
                    TIPOB,
                    OB_SIF,
                    NAZOB,
                    COALESCE(POLJE_TRAFO, '') as POLJE_TRAFO,
                    COALESCE(to_char(id_s_nap), '') as id_s_nap,
                    COALESCE(to_char(TRAFO_ID), '') as TRAFO_ID,
                    COALESCE(to_char(P2_TRAF_ID), '') as P2_TRAF_ID,
                    COALESCE(NAPON, '') as NAPON,
                    COALESCE(POLJE, '') as POLJE,
                    COALESCE(IME_PO, '') as IME_PO,
                    COALESCE(to_char(FUNKC), '') as FUNKC,
                    COALESCE(to_char(VRPD), ''),
                	NAZVRPD,
                    COALESCE(GRUZR1, '') as   GRUZR ,
                    COALESCE(UZROK1, '') as   UZROK,
                    COALESCE(GRRAZ, '') as GRRAZ ,
                    COALESCE(RAZLOG, '') as   RAZLOG,
                    COALESCE(VREM_USL, ''),
					COALESCE(opis, ''),
                    COALESCE(to_char(SNAGA), ''),
					COALESCE(NAZSOP, '') as NAZSOP,
   					COALESCE(SOP_NAZIV, '') as SOP_NAZIV,
   					COALESCE(to_char(ID_S_SOP), '') as ID_S_SOP,
   					COALESCE(to_char(ID_SOP), '') as ID_SOP,
                    COALESCE(TO_CHAR(ID_Z_DSDF_GL1), '') AS ID_Z_DSDF_GL1,
                    COALESCE(Z_DSDF_GL1, '') AS Z_DSDF_GL1,
                    COALESCE(TO_CHAR(ID_Z_KVAR_GL1), '') AS ID_Z_KVAR_GL1,
                    COALESCE(Z_KVAR_GL1, '') AS Z_KVAR_GL1,
                    COALESCE(TO_CHAR(ID_Z_RAPU_GL1), '') AS ID_Z_RAPU_GL1,
                    COALESCE(Z_RAPU_GL1, '') AS Z_RAPU_GL1,
                    COALESCE(TO_CHAR(ID_Z_PRST_GL1), '') AS ID_Z_PRST_GL1,
                    COALESCE(Z_PRST_GL1, '') AS Z_PRST_GL1,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_GL1), '') AS ID_Z_ZMSP_GL1,
                    COALESCE(Z_ZMSP_GL1, '') AS Z_ZMSP_GL1,
                    COALESCE(TO_CHAR(ID_Z_UZMS_GL1), '') AS ID_Z_UZMS_GL1,
                    COALESCE(Z_UZMS_GL1, '') AS Z_UZMS_GL1,
                    COALESCE(TO_CHAR(Z_LOKK_GL1), '') AS Z_LOKK_GL1,
                    COALESCE(TO_CHAR(ID_Z_DSDF_GL2), '') AS ID_Z_DSDF_GL2,
                    COALESCE(Z_DSDF_GL2, '') AS Z_DSDF_GL2,
                    COALESCE(TO_CHAR(ID_Z_KVAR_GL2), '') AS ID_Z_KVAR_GL2,
                    COALESCE(Z_KVAR_GL2, '') AS Z_KVAR_GL2,
                    COALESCE(TO_CHAR(ID_Z_RAPU_GL2), '') AS ID_Z_RAPU_GL2,
                    COALESCE(Z_RAPU_GL2, '') AS Z_RAPU_GL2,
                    COALESCE(TO_CHAR(ID_Z_PRST_GL2), '') AS ID_Z_PRST_GL2,
                    COALESCE(Z_PRST_GL2, '') AS Z_PRST_GL2,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_GL2), '') AS ID_Z_ZMSP_GL2,
                    COALESCE(Z_ZMSP_GL2, '') AS Z_ZMSP_GL2,
                    COALESCE(TO_CHAR(ID_Z_UZMS_GL2), '') AS ID_Z_UZMS_GL2,
                    COALESCE(Z_UZMS_GL2, '') AS Z_UZMS_GL2,
                    COALESCE(TO_CHAR(Z_LOKK_GL2), '') AS Z_LOKK_GL2,
                    COALESCE(TO_CHAR(ID_Z_DIS_REZ), '') AS ID_Z_DIS_REZ,
                    COALESCE(Z_DIS_REZ, '') AS Z_DIS_REZ,
                    COALESCE(TO_CHAR(ID_Z_KVAR_REZ), '') AS ID_Z_KVAR_REZ,
                    COALESCE(Z_KVAR_REZ, '') AS Z_KVAR_REZ,
                    COALESCE(TO_CHAR(ID_Z_PRST_REZ), '') AS ID_Z_PRST_REZ,
                    COALESCE(Z_PRST_REZ, '') AS Z_PRST_REZ,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_REZ), '') AS ID_Z_ZMSP_REZ,
                    COALESCE(Z_ZMSP_REZ, '') AS Z_ZMSP_REZ,
                    COALESCE(TO_CHAR(ID_Z_PREK_VN), '') AS ID_Z_PREK_VN,
                    COALESCE(Z_PREK_VN, '') AS Z_PREK_VN,
                    COALESCE(TO_CHAR(ID_Z_PREK_NN), '') AS ID_Z_PREK_NN,
                    COALESCE(Z_PREK_NN, '') AS Z_PREK_NN,
                    COALESCE(TO_CHAR(ID_Z_NEL1), '') AS ID_Z_NEL1,
                    COALESCE(Z_NEL1, '') AS Z_NEL1,
                    COALESCE(TO_CHAR(ID_Z_NEL2), '') AS ID_Z_NEL2,
                    COALESCE(Z_NEL2, '') AS Z_NEL2,
                    COALESCE(TO_CHAR(ID_Z_NEL3), '') AS ID_Z_NEL3,
                    COALESCE(Z_NEL3, '') AS Z_NEL3,
                    COALESCE(TO_CHAR(ID_Z_SABZ_SAB), '') AS ID_Z_SABZ_SAB,
                    COALESCE(Z_SABZ_SAB, '') AS Z_SABZ_SAB,
                    COALESCE(TO_CHAR(ID_Z_OTPR_SAB), '') AS ID_Z_OTPR_SAB,
                    COALESCE(Z_OTPR_SAB, '') AS Z_OTPR_SAB,
                    COALESCE(TO_CHAR(ID_Z_JPS_VN), '') AS ID_Z_JPS_VN,
                    COALESCE(Z_JPS_VN, '') AS Z_JPS_VN,
                    COALESCE(TO_CHAR(ID_Z_JPS_NN), '') AS ID_Z_JPS_NN,
                    COALESCE(Z_JPS_NN, '') AS Z_JPS_NN,
                    COALESCE(TO_CHAR(ID_Z_TELE_POC_GL1), '') AS ID_Z_TELE_POC_GL1,
                    COALESCE(Z_TELE_POC_GL1, '') AS Z_TELE_POC_GL1,
                    COALESCE(TO_CHAR(ID_Z_TELE_KRAJ_GL1), '') AS ID_Z_TELE_KRAJ_GL1,
                    COALESCE(Z_TELE_KRAJ_GL1, '') AS Z_TELE_KRAJ_GL1,
                    COALESCE(TO_CHAR(ID_Z_TELE_POC_GL2), '') AS ID_Z_TELE_POC_GL2,
                    COALESCE(Z_TELE_POC_GL2, '') AS Z_TELE_POC_GL2,
                    COALESCE(TO_CHAR(ID_Z_TELE_KRAJ_GL2), '') AS ID_Z_TELE_KRAJ_GL2,
                    COALESCE(Z_TELE_KRAJ_GL2, '') AS Z_TELE_KRAJ_GL2,
                    FUP,
                    COUNT(*) OVER () AS TOTAL_COUNT 
                    from pgi.pi_pi_v
                    where DATIZV  = to_date(:1,'dd.mm.yyyy') 
                       AND TIPD LIKE UPPER(:2)
                    and id_s_mrc like (:3)
					AND FUP LIKE UPPER(:4)
                       order by DATIZV,vrepoc`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.Datizv, arg.Tipd, mrcParam, fupParam)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.PiDD
	var totalCount int

	for rows.Next() {
		var ue models.PiDD
		var count int
		err := rows.Scan(
			&ue.Id,
			&ue.IdSMrc,
			&ue.Mrc,
			&ue.Tipd,
			&ue.Datizv,
			&ue.Id1,
			&ue.Id2,
			&ue.Vrepoc,
			&ue.VrepocHHMI,
			&ue.PocPP,
			&ue.Vrezav,
			&ue.VrezavHHMI,
			&ue.ZavPP,
			&ue.Traj,
			&ue.ObId,
			&ue.TipOb,
			&ue.ObSif,
			&ue.NazOb,
			&ue.PoljeTrafo,
			&ue.IdSNap,
			&ue.TrafoId,
			&ue.P2TrafId,
			&ue.Napon,
			&ue.Polje,
			&ue.ImePo,
			&ue.Funkc,
			&ue.Vrpd,
			&ue.Nazvrpd,
			&ue.GrUzrok,
			&ue.Uzrok,
			&ue.GrRazlog,
			&ue.Razlog,
			&ue.VrmUsl,
			&ue.Opis,
			&ue.Snaga,
			&ue.NazSop,
			&ue.SopNaziv,
			&ue.IdSSop,
			&ue.IdSop,
			&ue.IdZDsdfGl1,
			&ue.ZDsdfGl1,
			&ue.IdZKvarGl1,
			&ue.ZKvarGl1,
			&ue.IdZRapuGl1,
			&ue.ZRapuGl1,
			&ue.IdZPrstGl1,
			&ue.ZPrstGl1,
			&ue.IdZZmspGl1,
			&ue.ZZmspGl1,
			&ue.IdZUzmsGl1,
			&ue.ZUzmsGl1,
			&ue.ZLokkGl1,

			&ue.IdZDsdfGl2,
			&ue.ZDsdfGl2,
			&ue.IdZKvarGl2,
			&ue.ZKvarGl2,
			&ue.IdZRapuGl2,
			&ue.ZRapuGl2,
			&ue.IdZPrstGl2,
			&ue.ZPrstGl2,
			&ue.IdZZmspGl2,
			&ue.ZZmspGl2,
			&ue.IdZUzmsGl2,
			&ue.ZUzmsGl2,
			&ue.ZLokkGl2,

			&ue.IdZDisRez,
			&ue.ZDisRez,
			&ue.IdZKvarRez,
			&ue.ZKvarRez,
			&ue.IdZPrstRez,
			&ue.ZPrstRez,
			&ue.IdZZmspRez,
			&ue.ZZmspRez,
			&ue.IdZPrekVn,
			&ue.ZPrekVn,
			&ue.IdZPrekNn,
			&ue.ZPrekNn,
			&ue.IdZNel1,
			&ue.ZNel1,
			&ue.IdZNel2,
			&ue.ZNel2,
			&ue.IdZNel3,
			&ue.ZNel3,
			&ue.IdZSabzSab,
			&ue.ZSabzSab,
			&ue.IdZOtprSab,
			&ue.ZOtprSab,
			&ue.IdZJpsVn,
			&ue.ZJpsVn,
			&ue.IdZJpsNn,
			&ue.ZJpsNn,
			&ue.IdZTelePocGl1,
			&ue.ZTelePocGl1,
			&ue.IdZTeleKrajGl1,
			&ue.ZTeleKrajGl1,
			&ue.IdZTelePocGl2,
			&ue.ZTelePocGl2,
			&ue.IdZTeleKrajGl2,
			&ue.ZTeleKrajGl2,
			&ue.Fup,
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

func (m *OracleDBRepo) GetPiPIByParamsByPage(ctx context.Context, arg models.ListPiDDParamsByPage) ([]*models.PiDD, int, error) {

	fupParam := "%"
	if strings.ToUpper(arg.Fup) != "ALL" {
		fupParam = strings.ToUpper(arg.Fup)
	}

	mrcParam := "%"
	if strings.ToUpper(arg.IdSMrc) != "ALL" {
		mrcParam = strings.ToUpper(arg.IdSMrc)
	}

	query := `select ROWNUM AS id,
                	id_s_mrc,
                	mrc,
                    TIPD,
                    to_char(DATIZV,'dd.mm.yyyy'),
                    COALESCE(to_char(ID1), ''),
                    COALESCE(to_char(ID2), ''),
                    VREPOC_EXCEL,
                    TO_CHAR (vrepoc, 'hh24:mi') vrepoc_hhmi_excel,
                    COALESCE(POC_PP, '') as POC_PP,
                    VREZAV_EXCEL,
                    TO_CHAR (vrezav, 'hh24:mi') vrezav_hhmi_excel,
                    COALESCE(ZAV_PP, '') as ZAV_PP,
                    TRAJ_EXCEL,
                    OB_ID,
                    TIPOB,
                    OB_SIF,
                    NAZOB,
                    COALESCE(POLJE_TRAFO, '') as POLJE_TRAFO,
                    COALESCE(to_char(id_s_nap), '') as id_s_nap,
                    COALESCE(to_char(TRAFO_ID), '') as TRAFO_ID,
                    COALESCE(to_char(P2_TRAF_ID), '') as P2_TRAF_ID,
                    COALESCE(NAPON, '') as NAPON,
                    COALESCE(POLJE, '') as POLJE,
                    COALESCE(IME_PO, '') as IME_PO,
                    COALESCE(to_char(FUNKC), '') as FUNKC,
                    COALESCE(to_char(VRPD), ''),
                	NAZVRPD,
                    COALESCE(GRUZR1, '') as   GRUZR ,
                    COALESCE(UZROK1, '') as   UZROK,
                    COALESCE(GRRAZ, '') as GRRAZ ,
                    COALESCE(RAZLOG, '') as   RAZLOG,
                    COALESCE(VREM_USL, ''),
					COALESCE(opis, ''),
                    COALESCE(to_char(SNAGA), ''),
					COALESCE(NAZSOP, '') as NAZSOP,
   					COALESCE(SOP_NAZIV, '') as SOP_NAZIV,
   					COALESCE(to_char(ID_S_SOP), '') as ID_S_SOP,
   					COALESCE(to_char(ID_SOP), '') as ID_SOP,
                    COALESCE(TO_CHAR(ID_Z_DSDF_GL1), '') AS ID_Z_DSDF_GL1,
                    COALESCE(Z_DSDF_GL1, '') AS Z_DSDF_GL1,
                    COALESCE(TO_CHAR(ID_Z_KVAR_GL1), '') AS ID_Z_KVAR_GL1,
                    COALESCE(Z_KVAR_GL1, '') AS Z_KVAR_GL1,
                    COALESCE(TO_CHAR(ID_Z_RAPU_GL1), '') AS ID_Z_RAPU_GL1,
                    COALESCE(Z_RAPU_GL1, '') AS Z_RAPU_GL1,
                    COALESCE(TO_CHAR(ID_Z_PRST_GL1), '') AS ID_Z_PRST_GL1,
                    COALESCE(Z_PRST_GL1, '') AS Z_PRST_GL1,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_GL1), '') AS ID_Z_ZMSP_GL1,
                    COALESCE(Z_ZMSP_GL1, '') AS Z_ZMSP_GL1,
                    COALESCE(TO_CHAR(ID_Z_UZMS_GL1), '') AS ID_Z_UZMS_GL1,
                    COALESCE(Z_UZMS_GL1, '') AS Z_UZMS_GL1,
                    COALESCE(TO_CHAR(Z_LOKK_GL1), '') AS Z_LOKK_GL1,
                    COALESCE(TO_CHAR(ID_Z_DSDF_GL2), '') AS ID_Z_DSDF_GL2,
                    COALESCE(Z_DSDF_GL2, '') AS Z_DSDF_GL2,
                    COALESCE(TO_CHAR(ID_Z_KVAR_GL2), '') AS ID_Z_KVAR_GL2,
                    COALESCE(Z_KVAR_GL2, '') AS Z_KVAR_GL2,
                    COALESCE(TO_CHAR(ID_Z_RAPU_GL2), '') AS ID_Z_RAPU_GL2,
                    COALESCE(Z_RAPU_GL2, '') AS Z_RAPU_GL2,
                    COALESCE(TO_CHAR(ID_Z_PRST_GL2), '') AS ID_Z_PRST_GL2,
                    COALESCE(Z_PRST_GL2, '') AS Z_PRST_GL2,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_GL2), '') AS ID_Z_ZMSP_GL2,
                    COALESCE(Z_ZMSP_GL2, '') AS Z_ZMSP_GL2,
                    COALESCE(TO_CHAR(ID_Z_UZMS_GL2), '') AS ID_Z_UZMS_GL2,
                    COALESCE(Z_UZMS_GL2, '') AS Z_UZMS_GL2,
                    COALESCE(TO_CHAR(Z_LOKK_GL2), '') AS Z_LOKK_GL2,
                    COALESCE(TO_CHAR(ID_Z_DIS_REZ), '') AS ID_Z_DIS_REZ,
                    COALESCE(Z_DIS_REZ, '') AS Z_DIS_REZ,
                    COALESCE(TO_CHAR(ID_Z_KVAR_REZ), '') AS ID_Z_KVAR_REZ,
                    COALESCE(Z_KVAR_REZ, '') AS Z_KVAR_REZ,
                    COALESCE(TO_CHAR(ID_Z_PRST_REZ), '') AS ID_Z_PRST_REZ,
                    COALESCE(Z_PRST_REZ, '') AS Z_PRST_REZ,
                    COALESCE(TO_CHAR(ID_Z_ZMSP_REZ), '') AS ID_Z_ZMSP_REZ,
                    COALESCE(Z_ZMSP_REZ, '') AS Z_ZMSP_REZ,
                    COALESCE(TO_CHAR(ID_Z_PREK_VN), '') AS ID_Z_PREK_VN,
                    COALESCE(Z_PREK_VN, '') AS Z_PREK_VN,
                    COALESCE(TO_CHAR(ID_Z_PREK_NN), '') AS ID_Z_PREK_NN,
                    COALESCE(Z_PREK_NN, '') AS Z_PREK_NN,
                    COALESCE(TO_CHAR(ID_Z_NEL1), '') AS ID_Z_NEL1,
                    COALESCE(Z_NEL1, '') AS Z_NEL1,
                    COALESCE(TO_CHAR(ID_Z_NEL2), '') AS ID_Z_NEL2,
                    COALESCE(Z_NEL2, '') AS Z_NEL2,
                    COALESCE(TO_CHAR(ID_Z_NEL3), '') AS ID_Z_NEL3,
                    COALESCE(Z_NEL3, '') AS Z_NEL3,
                    COALESCE(TO_CHAR(ID_Z_SABZ_SAB), '') AS ID_Z_SABZ_SAB,
                    COALESCE(Z_SABZ_SAB, '') AS Z_SABZ_SAB,
                    COALESCE(TO_CHAR(ID_Z_OTPR_SAB), '') AS ID_Z_OTPR_SAB,
                    COALESCE(Z_OTPR_SAB, '') AS Z_OTPR_SAB,
                    COALESCE(TO_CHAR(ID_Z_JPS_VN), '') AS ID_Z_JPS_VN,
                    COALESCE(Z_JPS_VN, '') AS Z_JPS_VN,
                    COALESCE(TO_CHAR(ID_Z_JPS_NN), '') AS ID_Z_JPS_NN,
                    COALESCE(Z_JPS_NN, '') AS Z_JPS_NN,
                    COALESCE(TO_CHAR(ID_Z_TELE_POC_GL1), '') AS ID_Z_TELE_POC_GL1,
                    COALESCE(Z_TELE_POC_GL1, '') AS Z_TELE_POC_GL1,
                    COALESCE(TO_CHAR(ID_Z_TELE_KRAJ_GL1), '') AS ID_Z_TELE_KRAJ_GL1,
                    COALESCE(Z_TELE_KRAJ_GL1, '') AS Z_TELE_KRAJ_GL1,
                    COALESCE(TO_CHAR(ID_Z_TELE_POC_GL2), '') AS ID_Z_TELE_POC_GL2,
                    COALESCE(Z_TELE_POC_GL2, '') AS Z_TELE_POC_GL2,
                    COALESCE(TO_CHAR(ID_Z_TELE_KRAJ_GL2), '') AS ID_Z_TELE_KRAJ_GL2,
                    COALESCE(Z_TELE_KRAJ_GL2, '') AS Z_TELE_KRAJ_GL2,
                    FUP,
                    COUNT(*) OVER () AS TOTAL_COUNT  
                    from pgi.pi_pi_v
                    where DATIZV  = to_date(:1,'dd.mm.yyyy') 
                       AND TIPD LIKE UPPER(:2)
                    AND FUP LIKE UPPER(:3)
                    and id_s_mrc like (:4)
                       order by DATIZV,id1,id2,vrepoc
					OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.Datizv, arg.Tipd, fupParam, mrcParam, arg.Offset, arg.Limit)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, 0, err
	}
	defer rows.Close()

	var ues []*models.PiDD
	var totalCount int

	for rows.Next() {
		var ue models.PiDD
		var count int
		err := rows.Scan(
			&ue.Id,
			&ue.IdSMrc,
			&ue.Mrc,
			&ue.Tipd,
			&ue.Datizv,
			&ue.Id1,
			&ue.Id2,
			&ue.Vrepoc,
			&ue.VrepocHHMI,
			&ue.PocPP,
			&ue.Vrezav,
			&ue.VrezavHHMI,
			&ue.ZavPP,
			&ue.Traj,
			&ue.ObId,
			&ue.TipOb,
			&ue.ObSif,
			&ue.NazOb,
			&ue.PoljeTrafo,
			&ue.IdSNap,
			&ue.TrafoId,
			&ue.P2TrafId,
			&ue.Napon,
			&ue.Polje,
			&ue.ImePo,
			&ue.Funkc,
			&ue.Vrpd,
			&ue.Nazvrpd,
			&ue.GrUzrok,
			&ue.Uzrok,
			&ue.GrRazlog,
			&ue.Razlog,
			&ue.VrmUsl,
			&ue.Opis,
			&ue.Snaga,
			&ue.NazSop,
			&ue.SopNaziv,
			&ue.IdSSop,
			&ue.IdSop,
			&ue.IdZDsdfGl1,
			&ue.ZDsdfGl1,
			&ue.IdZKvarGl1,
			&ue.ZKvarGl1,
			&ue.IdZRapuGl1,
			&ue.ZRapuGl1,
			&ue.IdZPrstGl1,
			&ue.ZPrstGl1,
			&ue.IdZZmspGl1,
			&ue.ZZmspGl1,
			&ue.IdZUzmsGl1,
			&ue.ZUzmsGl1,
			&ue.ZLokkGl1,

			&ue.IdZDsdfGl2,
			&ue.ZDsdfGl2,
			&ue.IdZKvarGl2,
			&ue.ZKvarGl2,
			&ue.IdZRapuGl2,
			&ue.ZRapuGl2,
			&ue.IdZPrstGl2,
			&ue.ZPrstGl2,
			&ue.IdZZmspGl2,
			&ue.ZZmspGl2,
			&ue.IdZUzmsGl2,
			&ue.ZUzmsGl2,
			&ue.ZLokkGl2,

			&ue.IdZDisRez,
			&ue.ZDisRez,
			&ue.IdZKvarRez,
			&ue.ZKvarRez,
			&ue.IdZPrstRez,
			&ue.ZPrstRez,
			&ue.IdZZmspRez,
			&ue.ZZmspRez,
			&ue.IdZPrekVn,
			&ue.ZPrekVn,
			&ue.IdZPrekNn,
			&ue.ZPrekNn,
			&ue.IdZNel1,
			&ue.ZNel1,
			&ue.IdZNel2,
			&ue.ZNel2,
			&ue.IdZNel3,
			&ue.ZNel3,
			&ue.IdZSabzSab,
			&ue.ZSabzSab,
			&ue.IdZOtprSab,
			&ue.ZOtprSab,
			&ue.IdZJpsVn,
			&ue.ZJpsVn,
			&ue.IdZJpsNn,
			&ue.ZJpsNn,
			&ue.IdZTelePocGl1,
			&ue.ZTelePocGl1,
			&ue.IdZTeleKrajGl1,
			&ue.ZTeleKrajGl1,
			&ue.IdZTelePocGl2,
			&ue.ZTelePocGl2,
			&ue.IdZTeleKrajGl2,
			&ue.ZTeleKrajGl2,
			&ue.Fup,
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

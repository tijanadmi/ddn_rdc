package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetOtvoreneSmene(ctx context.Context) ([]models.Smena, error) {

	query := `
select 
  smena.id,
  smena.datdnev,
  smena.id_s_mrc,
  ted.TD_NAZIVI.TD_DAJ_SIF('S_MRC','NAZIV_CIR','ID',SMENA.ID_S_MRC,'Q') RDC,
  smena.dez_disp1,
  ek1.ime || ' ' || ek1.prezime,
  smena.dez_disp2,
  ek2.ime || ' ' || ek2.prezime,
  smena.dez_disp3,
  ek3.ime || ' ' || ek3.prezime,
  smena.dez_disp4,
  ek4.ime || ' ' || ek4.prezime,
  smena.ID_TIP_SMENA,
  ted.TD_NAZIVI.TD_DAJ_SIF('DDN.TIP_SMENA','SKR_NAZ','ID',SMENA.ID_TIP_SMENA,'Q'),
  ted.TD_NAZIVI.TD_DAJ_SIF('DDN.TIP_SMENA','NAZIV','ID',SMENA.ID_TIP_SMENA,'Q'),
  ted.TD_NAZIVI.TD_DAJ_SIF('DDN.TIP_SMENA','INTERVAL','ID',SMENA.ID_TIP_SMENA,'Q'),
  ted.TD_NAZIVI.TD_DAJ_SIF('ted.S_MRC','SIFRA','ID',smena.ID_S_MRC,'Q'),
  SMENA.PREDAO_DISP1,
  SMENA.PREDAO_DISP2,
  SMENA.PRIM_DISP1,
  SMENA.PRIM_DISP2,
  SMENA.PREDAO_DISP3,
  SMENA.PRIM_DISP3,
  smena.koment_zat,
  smena.otv_spec,
  smena.zat_spec,
  SMENA.ID_KAT_DOK
from ddn.smena smena
left join tis_kor_v tk1 on tk1.sifra = smena.dez_disp1
left join ems_kadar ek1 on ek1.id = tk1.id_hr_kadar
left join tis_kor_v tk2 on tk2.sifra = smena.dez_disp2
left join ems_kadar ek2 on ek2.id = tk2.id_hr_kadar
left join tis_kor_v tk3 on tk3.sifra = smena.dez_disp3
left join ems_kadar ek3 on ek3.id = tk3.id_hr_kadar
left join tis_kor_v tk4 on tk4.sifra = smena.dez_disp4
left join ems_kadar ek4 on ek4.id = tk4.id_hr_kadar
where smena.stat_smene = 0
`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var smene []models.Smena

	for rows.Next() {
		var s models.Smena

		err := rows.Scan(
			&s.IdSmene,
			&s.DatDnev,

			&s.IdSMRC,
			&s.Rdc,

			&s.DezDisp1,
			&s.DezDisp1Ime,

			&s.DezDisp2,
			&s.DezDisp2Ime,

			&s.DezDisp3,
			&s.DezDisp3Ime,

			&s.DezDisp4,
			&s.DezDisp4Ime,

			&s.IDTipSmena,
			&s.TipSkr,
			&s.TipSmena,
			&s.IntSmena,
			&s.MrcSif,

			&s.PredaoDisp1,
			&s.PredaoDisp2,
			&s.PrimDisp1,
			&s.PrimDisp2,
			&s.PredaoDisp3,
			&s.PrimDisp3,
			&s.KomentZat,
			&s.OtvSpec,
			&s.ZatSpec,
			&s.IDKatDok,
		)
		if err != nil {
			fmt.Printf("Greška prilikom skeniranja smene: %v\n", err)
			return nil, err
		}
		// fmt.Printf("Pročitana smena: %+v\n", s.IdSmene)

		smene = append(smene, s)
	}
	// fmt.Println("Posle iscitanjih smena")

	if err := rows.Err(); err != nil {
		fmt.Println("Greška prilikom iteracije kroz smene:", err)
		return nil, err
	}

	// 👉 DRUGI KORAK – učitavanje događaja po smeni

	dogQuery := `
SELECT 
  dog_smene.id,
  dog_smene.id_smena,
  dog_smene.rb_dog,
  dog_smene.naslov,
  dog_smene.id_tip_dog,
  tip_dog.naziv,
  tip_dog.naziv_cir,
  tip_dog.tip,
  dog_smene.dopuna,
  dog_smene.status
FROM ddn.dog_smene, ddn.tip_dog
WHERE dog_smene.id_smena = :1
and ddn.dog_smene.id_tip_dog = tip_dog.id
ORDER BY dog_smene.rb_dog
`

	for i := range smene {
		rowsDog, err := m.DB.QueryContext(ctx, dogQuery, smene[i].IdSmene)
		if err != nil {
			fmt.Printf("Greška prilikom učitavanja događaja za smenu %d: %v\n", smene[i].IdSmene, err)
			return nil, err

		}

		var dogadjaji []models.Dogadjaj

		for rowsDog.Next() {
			var d models.Dogadjaj

			err := rowsDog.Scan(
				&d.ID,
				&d.IDSmena,
				&d.RbDog,
				&d.Naslov,
				&d.IDTipDog,
				&d.TipDog,
				&d.TipDogCir,
				&d.Tip,
				&d.Dopuna,
				&d.Status,
			)
			if err != nil {
				rowsDog.Close()
				fmt.Printf("Greška prilikom skeniranja događaja za smenu %d: %v\n", smene[i].IdSmene, err)
				return nil, err
			}
			// fmt.Printf("Pročitani događaj za smenu %d: %+v\n", smene[i].IdSmene, d)

			dogadjaji = append(dogadjaji, d)
		}

		if err := rowsDog.Err(); err != nil {
			rowsDog.Close()
			fmt.Printf("Greška prilikom iteracije kroz događaje za smenu %d: %v\n", smene[i].IdSmene, err)
			return nil, err
		}

		rowsDog.Close()

		smene[i].Dogadjaji = dogadjaji
	}

	return smene, nil
}

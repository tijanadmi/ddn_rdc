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

	//  DRUGI KORAK – učitavanje događaja po smeni

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

/****** Funkcija vraca tip dogadjaja Iskljucenje/ukljucenje sa manipulacijama *************/

func buildRecenica(d *models.DogadjajIskljucenje) string {
	format := "02.01.2006"

	if d.Dopuna != nil {
		switch *d.Dopuna {
		case "1":
			if d.TipSmeneDopune != nil && *d.TipSmeneDopune == "D" {
				return "Dopuna dnevne smene od: " + d.DatumDopune.Format(format) + "\n"
			} else if d.DatumDopune != nil {
				return "Dopuna noćne smene od: " +
					d.DatumDopune.Format(format) + "/" +
					d.DatumDopune.AddDate(0, 0, 1).Format(format) + "\n"
			}

		case "2":
			if d.VezaSa != nil && d.DatumDopune != nil {
				if d.TipSmeneDopune != nil && *d.TipSmeneDopune == "D" {
					return fmt.Sprintf("Dopuna događaja br. %s dnevne smene od %s\n",
						*d.RbDogVezaSa, d.DatumDopune.Format(format))
				} else {
					return fmt.Sprintf("Dopuna događaja br. %s noćne smene od %s/%s\n",
						*d.RbDogVezaSa,
						d.DatumDopune.Format(format),
						d.DatumDopune.AddDate(0, 0, 1).Format(format))
				}
			}
		}
	}

	// default redovan unos
	datum := d.DatumSmene
	if d.TipSmene == "N" {
		return fmt.Sprintf("Redovan unos za smenu od: %s / %s",
			datum.Format(format),
			datum.AddDate(0, 0, 1).Format(format))
	}

	return fmt.Sprintf("Redovan unos za smenu od: %s", datum.Format(format))
}

func (m *OracleDBRepo) GetIskljucenjeById(ctx context.Context, id int) (*models.DogadjajIskljucenje, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.rb_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  ted.TD_NAZIVI.TD_DAJ_SIF('DDN.DOG_SMENE','RB_DOG','ID',d.ID_DOG_SMENE,'Q') AS RB_DOG_VEZA_SA,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
WHERE d.id = :1
`

	var d models.DogadjajIskljucenje

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.RbDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.DatumSmene,
		&d.DatumDopune,
		&d.TipSmene,
		&d.TipSmeneDopune,
	)

	if err != nil {
		return nil, err
	}

	/**** Izmena naslova podnaslova ****/
	d.Podnaslov = buildRecenica(&d)

	//  DETAIL QUERY
	detailQuery := `
select 
A.ID_DOG_SMENE,a.rb,A.VREPOC,A.VREZAV,a.manipulacija,a.tekst_man,a.ees,a.tekst_ees,napomena,
case when tip_ob in ('DV','TD','KB','TK') THEN decode(p2_traf_id,null,'',TD_NAZIVI.COL_V_S_OB_P2_TRAF(P2_TRAF_ID,'OPIS','Q'))
     else OBJEKAT_NAZIV
     end objekat,
case when tip_ob in ('DV','TD','KB','TK') THEN OBJEKAT_SIFRA
     else case when TRAFO_ID IS NOT NULL THEN ted.td_nazivi.td_daj_sif('V_S_TR','TRAFO_NAZ','P0_TRAN_ID',TRAFO_ID,'Q')||
     NVL2(NAPON,' na '||napon||' kV ','')
     else case when (funkc is not null and napon is not null) then 
                     CASE WHEN ted.TD_NAZIVI.TD_DAJ_SIF('s_fup','SIFRA','ID',funkc,'Q') IN ('02','14','16','17') THEN
                      ted.TD_NAZIVI.TD_DAJ_SIF('s_fup','NAZIV','ID',funkc,'Q')||  
                                       case when napon<110 then ' POLJE br. ' || POLJE || ' (' || IME_PO || ') na ' || napon||' kV '
                                            else ' POLJE (' || IME_PO || ') na ' || napon||' kV '
                                            end 
                      ELSE  
                           case when napon<110 then IME_PO ||  ' POLJE br. ' || POLJE || ' na ' ||  napon||' kV '
                           else IME_PO ||  ' POLJE ' || ' na ' ||  napon||' kV '
                           end
                      END
         else ' '                                              
     end end end dv_trafo_polje, A.STATUS STATUS_MAN,A.DOPUNA DOPUNA_MAN, ID_STAVKE_M
from (
    select ID_DOG_SMENE,
    rb,
    to_char(vrepoc,'hh24:mi') vrepoc,
    to_char(vrezav,'hh24:mi') vrezav,
    ted.TD_NAZIVI.TD_DAJ_SIF('DDN.TIP_MAN','NAZIV','ID',DOG_MAN.ID_TIP_MAN,'Q') manipulacija,
    tekst_man,
    ted.TD_NAZIVI.TD_DAJ_SIF('DDN.TIP_EES','NAZIV','ID',DOG_MAN.ID_TIP_EES,'Q') ees,
    tekst_ees,
    napomena,
    ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',DOG_MAN.ID_TIPOB,'Q') tip_ob,
    pgi.PI_NAZIVI.NAZ_OBJ(DOG_MAN.ID_TIPOB,DOG_MAN.OB_ID,'OB_SIF','Q') OBJEKAT_SIFRA,
    pgi.PI_NAZIVI.NAZ_OBJ(DOG_MAN.ID_TIPOB,DOG_MAN.OB_ID,'OPIS','Q') OBJEKAT_NAZIV,
    ted.TD_NAZIVI.TD_DAJ_SIF('S_NAP','NAZIV','ID',DOG_MAN.ID_S_NAP,'Q') NAPON, 
    ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','POLJE','P2_TRAF_ID',DOG_MAN.P2_TRAF_ID,'Q') POLJE, 
    ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','IME_PO','P2_TRAF_ID',DOG_MAN.P2_TRAF_ID,'Q') IME_PO, 
    ted.TD_NAZIVI.TD_DAJ_ID('V_S_POLJE_SVA_AP','FUP_ID','P2_TRAF_ID',DOG_MAN.P2_TRAF_ID,'Q') FUNKC,
    ob_id,trafo_id,P2_TRAF_ID,STATUS,DOPUNA,DOG_MAN.ID ID_STAVKE_M
    from dog_man 
    where id_dog_smene = :1
) a
ORDER BY rb
`

	rows, err := m.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manipulacije []models.Manipulacija

	for rows.Next() {
		var mnp models.Manipulacija

		err := rows.Scan(
			&mnp.IDDogSmene,
			&mnp.Rb,
			&mnp.Vrepoc,
			&mnp.Vrezav,
			&mnp.Manipulacija,
			&mnp.TekstMan,
			&mnp.Ees,
			&mnp.TekstEes,
			&mnp.Napomena,
			&mnp.Objekat,
			&mnp.DvTrafoPolje,
			&mnp.StatusMan,
			&mnp.DopunaMan,
			&mnp.IDStavkeM,
		)
		if err != nil {
			return nil, err
		}

		manipulacije = append(manipulacije, mnp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	d.Manipulacije = manipulacije

	return &d, nil
}

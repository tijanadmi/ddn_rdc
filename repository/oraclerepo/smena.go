package oraclerepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
  d.id,
  d.id_smena,
  d.rb_dog,
  d.naslov,
  d.id_tip_dog,
  td.naziv,
  td.naziv_cir,
  td.tip,
   ted.TD_NAZIVI.TD_DAJ_SIF('TIP_OBV','SIFRA','ID',o.id_tip_obv,'Q') AS tip_obv,
  d.dopuna,
  d.status

FROM ddn.dog_smene d

JOIN ddn.tip_dog td 
  ON d.id_tip_dog = td.id

LEFT JOIN ddn.dog_obav o 
  ON d.id = o.id_dog_smene

WHERE d.id_smena = :1
ORDER BY d.rb_dog
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
				&d.TipObav,
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

/****** F-ja vraca sve zatvorene smene sa pripadajucim dogadjajima za izabrani RDC i interval	 *************/

func (m *OracleDBRepo) GetZatvoreneSmene(ctx context.Context, arg models.ListShiftsWithPaginationParams) ([]models.Smena, int, error) {

	mrcParam := "%"
	if strings.ToUpper(arg.Mrc) != "ALL" {
		mrcParam = strings.ToUpper(arg.Mrc)
	}

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
  SMENA.ID_KAT_DOK,
  COUNT(*) OVER () AS TOTAL_COUNT
from ddn.smena smena
left join tis_kor_v tk1 on tk1.sifra = smena.dez_disp1
left join ems_kadar ek1 on ek1.id = tk1.id_hr_kadar
left join tis_kor_v tk2 on tk2.sifra = smena.dez_disp2
left join ems_kadar ek2 on ek2.id = tk2.id_hr_kadar
left join tis_kor_v tk3 on tk3.sifra = smena.dez_disp3
left join ems_kadar ek3 on ek3.id = tk3.id_hr_kadar
left join tis_kor_v tk4 on tk4.sifra = smena.dez_disp4
left join ems_kadar ek4 on ek4.id = tk4.id_hr_kadar
where smena.stat_smene = 1
and smena.ID_S_MRC like (:1) AND  
    (TRUNC(smena.datdnev) BETWEEN TO_DATE(:2, 'dd.mm.yyyy')
                           AND TO_DATE(:3, 'dd.mm.yyyy'))
   ORDER BY id
              OFFSET :4 ROWS FETCH NEXT :5 ROWS ONLY
`

	rows, err := m.DB.QueryContext(ctx, query, mrcParam, arg.StartDate, arg.EndDate, arg.Offset, arg.Limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var smene []models.Smena
	var totalCount int

	for rows.Next() {
		var s models.Smena
		var count int

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
			&count,
		)
		if err != nil {
			fmt.Printf("Greška prilikom skeniranja smene: %v\n", err)
			return nil, 0, err
		}
		// fmt.Printf("Pročitana smena: %+v\n", s.IdSmene)

		smene = append(smene, s)
		totalCount = count
	}
	// fmt.Println("Posle iscitanjih smena")

	if err := rows.Err(); err != nil {
		fmt.Println("Greška prilikom iteracije kroz smene:", err)
		return nil, 0, err
	}

	//  DRUGI KORAK – učitavanje događaja po smeni

	dogQuery := `
SELECT 
  d.id,
  d.id_smena,
  d.rb_dog,
  d.naslov,
  d.id_tip_dog,
  td.naziv,
  td.naziv_cir,
  td.tip,
   ted.TD_NAZIVI.TD_DAJ_SIF('TIP_OBV','SIFRA','ID',o.id_tip_obv,'Q') AS tip_obv,
  d.dopuna,
  d.status

FROM ddn.dog_smene d

JOIN ddn.tip_dog td 
  ON d.id_tip_dog = td.id

LEFT JOIN ddn.dog_obav o 
  ON d.id = o.id_dog_smene

WHERE d.id_smena = :1
ORDER BY d.rb_dog
`

	for i := range smene {
		rowsDog, err := m.DB.QueryContext(ctx, dogQuery, smene[i].IdSmene)
		if err != nil {
			fmt.Printf("Greška prilikom učitavanja događaja za smenu %d: %v\n", smene[i].IdSmene, err)
			return nil, 0, err

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
				&d.TipObav,
				&d.Dopuna,
				&d.Status,
			)
			if err != nil {
				rowsDog.Close()
				fmt.Printf("Greška prilikom skeniranja događaja za smenu %d: %v\n", smene[i].IdSmene, err)
				return nil, 0, err
			}
			// fmt.Printf("Pročitani događaj za smenu %d: %+v\n", smene[i].IdSmene, d)

			dogadjaji = append(dogadjaji, d)
		}

		if err := rowsDog.Err(); err != nil {
			rowsDog.Close()
			fmt.Printf("Greška prilikom iteracije kroz događaje za smenu %d: %v\n", smene[i].IdSmene, err)
			return nil, 0, err
		}

		rowsDog.Close()

		smene[i].Dogadjaji = dogadjaji
	}

	return smene, totalCount, nil
}

/****** Funkcija vraca tip dogadjaja Iskljucenje/ukljucenje sa manipulacijama *************/

func buildRecenica(d *models.DogadjajDetaljno) string {
	format := "02.01.2006"

	var dop string

	// =========================
	// 1. DOPUNA LOGIKA
	// =========================
	if d.Dopuna != nil {
		switch *d.Dopuna {

		case "1":
			if d.DatumDopune != nil {
				if d.TipSmeneDopune != nil && *d.TipSmeneDopune == "D" {
					dop = "Dopuna dnevne smene od: " +
						d.DatumDopune.Format(format) + "\n"
				} else {
					dop = "Dopuna noćne smene od: " +
						d.DatumDopune.Format(format) + "/" +
						d.DatumDopune.AddDate(0, 0, 1).Format(format) + "\n"
				}
			}

		case "2":
			if d.DatumDopune != nil && d.RbDogVezaSa != nil {
				if d.TipSmeneDopune != nil && *d.TipSmeneDopune == "D" {
					dop = fmt.Sprintf(
						"Dopuna događaja br. %s dnevne smene od %s\n",
						*d.RbDogVezaSa,
						d.DatumDopune.Format(format),
					)
				} else {
					dop = fmt.Sprintf(
						"Dopuna događaja br. %s noćne smene od %s/%s\n",
						*d.RbDogVezaSa,
						d.DatumDopune.Format(format),
						d.DatumDopune.AddDate(0, 0, 1).Format(format),
					)
				}
			}
		}
	}

	// =========================
	// 2. DEFAULT (REDOVAN UNOS)
	// =========================
	if dop == "" {
		datum := d.DatumSmene

		if d.TipSmene == "N" {
			dop = fmt.Sprintf(
				"Redovan unos za smenu od: %s / %s",
				datum.Format(format),
				datum.AddDate(0, 0, 1).Format(format),
			)
		} else {
			dop = fmt.Sprintf(
				"Redovan unos za smenu od: %s",
				datum.Format(format),
			)
		}
	}

	// =========================
	// 3. VEZA LOGIKA (NAJBITNIJE)
	// =========================
	if (d.Dopuna == nil || *d.Dopuna != "2") && d.VezaSa != nil && d.RbDogVezaSa != nil && d.DatumVezaSa != nil {

		datum := *d.DatumVezaSa

		if d.TipSmenaVezaSa != nil && *d.TipSmenaVezaSa == "D" {
			return fmt.Sprintf(
				"Veza sa događajem br. %s od dana: %s - %s",
				*d.RbDogVezaSa,
				datum.Format(format),
				dop,
			)
		} else {
			return fmt.Sprintf(
				"Veza sa događajem br. %s od dana: %s/%s - %s",
				*d.RbDogVezaSa,
				datum.Format(format),
				datum.AddDate(0, 0, 1).Format(format),
				dop,
			)
		}
	}

	// =========================
	// 4. BEZ VEZE
	// =========================
	return fmt.Sprintf("%s", dop)
}

func (m *OracleDBRepo) GetIskljucenjeById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
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
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
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

	if len(manipulacije) > 0 {
		d.Manipulacije = manipulacije
	}

	return &d, nil
}

func (m *OracleDBRepo) GetObavBeleskaById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	// =========================
	// MASTER QUERY
	// =========================
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,

  d.DOPUNA,
  d.ID_SMENA_D,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.DatumSmene,
		&d.DatumDopune,
		&d.TipSmene,
		&d.TipSmeneDopune,
	)

	if err != nil {
		return nil, err
	}

	// ✅ podnaslov (isti buildRecenica)
	d.Podnaslov = buildRecenica(&d)

	// =========================
	// DETAIL QUERY (DOG_OBAV)
	// =========================
	detailQuery := `
SELECT 
  o.NAPOMENA,
  o.TEKST_OBV,
  t.SIFRA,
  o.DOPUNA
FROM dog_obav o
JOIN TIP_OBV t ON o.ID_TIP_OBV = t.ID
WHERE o.ID_DOG_SMENE = :1
AND t.SIFRA != 'F'
ORDER BY o.RB
FETCH FIRST 1 ROWS ONLY
`

	var obav models.ObavBeleska
	var tekstObv sql.NullString

	err = m.DB.QueryRowContext(ctx, detailQuery, id).Scan(
		&obav.Napomena,
		&tekstObv,
		&obav.TipObv,
		&obav.Dopuna,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Nema obavestenja, ostaje nil
			d.ObavBeleske = nil
		} else {
			return nil, err
		}
	} else {
		if tekstObv.Valid {
			obav.TekstObv = tekstObv.String
		}
		d.ObavBeleske = &obav
	}

	return &d, nil
}

func (m *OracleDBRepo) GetRadTKById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,
  d.POSL_TEKST,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.Posledice,
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
SELECT 
    d.ID,
    TO_CHAR(d.vrepoc,'hh24:mi') vrepoc,
    TO_CHAR(d.vrezav,'hh24:mi') vrezav,

    ted.TD_NAZIVI.TD_DAJ_SIF('S_VRPD','NAZIV','ID', d.ID_S_VRPD,'Q') VRSTADOG,

    d.OPIS,
    d.ob_id,

    v1.nazob AS OBJEKAT_NAZIV,

    d.OB_ID2,

    v2.nazob AS OBJEKAT_NAZIV2,

    TD_NAZIVI.TD_DAJ_SIF('S_VROPR','NAZIV','ID', d.ID_S_VROPR,'Q') VROPR,

    d.id_s_vrpd AS VRPD,

    d.STATUS AS STATUS_ISPKV5678,
    d.DOPUNA AS DOPUNA_ISPKV5678,

    d.vrepoc AS VREPOC_SORT_ISPKV5678

FROM DOG_ISPKV d

-- objekat 1
LEFT JOIN V_s_ob v1 
    ON v1.ob_id = d.ob_id
   AND v1.tipob = d.id_tipob

-- objekat 2
LEFT JOIN V_s_ob v2 
    ON v2.ob_id = d.ob_id2
   AND v2.tipob = d.id2_tipob

WHERE d.id_dog_smene = :1

ORDER BY d.vrepoc
`

	rows, err := m.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tks []models.TK

	for rows.Next() {
		var tk models.TK

		err := rows.Scan(
			&tk.ID,
			&tk.Vrepoc,
			&tk.Vrezav,
			&tk.VrstaDog,
			&tk.Opis,
			&tk.ObID,
			&tk.ObjekatNaziv,
			&tk.ObID2,
			&tk.ObjekatNaziv2,
			&tk.Vropr,
			&tk.Vrpd,
			&tk.Status,
			&tk.Dopuna,
			&tk.VrepocSort,
		)
		if err != nil {
			return nil, err
		}

		tks = append(tks, tk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tks) > 0 {
		d.TK = &tks
	}

	return &d, nil
}

func (m *OracleDBRepo) GetRadTSUById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,
  d.POSL_TEKST,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.Posledice,
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
SELECT 
    d.ID,
    TO_CHAR(d.vrepoc,'hh24:mi') vrepoc,
    TO_CHAR(d.vrezav,'hh24:mi') vrezav,

    ted.TD_NAZIVI.TD_DAJ_SIF('S_VRPD','NAZIV','ID', d.ID_S_VRPD,'Q') VRSTADOG,

    d.OPIS,
    d.ob_id,

    v1.nazob AS OBJEKAT_NAZIV,

    TD_NAZIVI.TD_DAJ_SIF('S_VROPR','NAZIV','ID', d.ID_S_VROPR,'Q') VROPR,

    d.id_s_vrpd AS VRPD,

    d.STATUS AS STATUS_ISPKV5678,
    d.DOPUNA AS DOPUNA_ISPKV5678,

    d.vrepoc AS VREPOC_SORT_ISPKV5678

FROM DOG_ISPKV d

-- objekat 1
LEFT JOIN V_s_ob v1 
    ON v1.ob_id = d.ob_id
   AND v1.tipob = d.id_tipob

WHERE d.id_dog_smene = :1

ORDER BY d.vrepoc
`

	rows, err := m.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tsus []models.TSU

	for rows.Next() {
		var tsu models.TSU

		err := rows.Scan(
			&tsu.ID,
			&tsu.Vrepoc,
			&tsu.Vrezav,
			&tsu.VrstaDog,
			&tsu.Opis,
			&tsu.ObID,
			&tsu.ObjekatNaziv,
			&tsu.Vropr,
			&tsu.Vrpd,
			&tsu.Status,
			&tsu.Dopuna,
			&tsu.VrepocSort,
		)
		if err != nil {
			return nil, err
		}

		tsus = append(tsus, tsu)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tsus) > 0 {
		d.TSU = &tsus
	}

	return &d, nil
}

func (m *OracleDBRepo) GetRadSOPById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,
  d.POSL_TEKST,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.Posledice,
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
SELECT 
    d.ID,
    TO_CHAR(d.vrepoc,'hh24:mi') vrepoc,
    TO_CHAR(d.vrezav,'hh24:mi') vrezav,

    ted.TD_NAZIVI.TD_DAJ_SIF('S_VRPD','NAZIV','ID', d.ID_S_VRPD,'Q') VRSTADOG,

    d.OPIS,
    d.ob_id,

    v1.nazob AS OBJEKAT_NAZIV,

    TD_NAZIVI.TD_DAJ_SIF('S_SOP','NAZIV','ID',d.ID_S_SOP,'Q') NAZSOP,
    s1.r_br,
    d.ID_S_SOP,
    d.ID_SOP,

    d.id_s_vrpd AS VRPD,

    d.STATUS AS STATUS_ISPKV5678,
    d.DOPUNA AS DOPUNA_ISPKV5678,

    d.vrepoc AS VREPOC_SORT_ISPKV5678

FROM DOG_ISPKV d

-- objekat 1
LEFT JOIN V_s_ob v1 
    ON v1.ob_id = d.ob_id
   AND v1.tipob = d.id_tipob

-- sop
LEFT JOIN v_s_spo s1 
    ON s1.tipsp = d.ID_S_SOP
   AND s1.sop_id = d.ID_SOP

WHERE d.id_dog_smene = :1

ORDER BY d.vrepoc
`

	rows, err := m.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sops []models.SOP

	for rows.Next() {
		var sop models.SOP

		err := rows.Scan(
			&sop.ID,
			&sop.Vrepoc,
			&sop.Vrezav,
			&sop.VrstaDog,
			&sop.Opis,
			&sop.ObID,
			&sop.ObjekatNaziv,
			&sop.NazSop,
			&sop.RBrSop,
			&sop.IdSSop,
			&sop.IdSop,
			&sop.Vrpd,
			&sop.Status,
			&sop.Dopuna,
			&sop.VrepocSort,
		)
		if err != nil {
			return nil, err
		}

		sops = append(sops, sop)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(sops) > 0 {
		d.SOP = &sops
	}

	return &d, nil
}

func (m *OracleDBRepo) GetIspadById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,
  d.POSL_TEKST,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.Posledice,
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

	/**** Deo za ispad ****/

	//  DETAIL QUERY
	detailQueryI := `
SELECT 
    A.VREPOC_SORT_ISPKV,
    A.ID,
    A.RB,
    A.VREPOC,
    A.VREZAV,
    A.VRSTADOG,
    A.VR_DOG_SIF,
    A.GRUZR1,
    A.UZROK1,

    CASE 
        WHEN A.tip_ob IN ('DV','TD','KB','TK') 
        THEN DECODE(A.p2_traf_id, NULL, '', TD_NAZIVI.COL_V_S_OB_P2_TRAF(P2_TRAF_ID,'NAZOB','Q'))
        ELSE OBJEKAT_NAZIV
    END objekat,

    TIP_OB,
    NAPON,

    CASE 
        WHEN A.tip_ob IN ('DV','TD','KB','TK') THEN OBJEKAT_SIFRA
        WHEN TRAFO_ID IS NOT NULL THEN 
            ted.td_nazivi.td_daj_sif('V_S_TR','TRAFO_NAZ','P0_TRAN_ID',TRAFO_ID,'Q') ||
            NVL2(NAPON,' na '||napon||' kV ','')
        WHEN (funkc IS NOT NULL AND napon IS NOT NULL) THEN 
            CASE 
                WHEN ted.TD_NAZIVI.TD_DAJ_SIF('s_fup','SIFRA','ID',funkc,'Q') IN ('02','14','16','17') THEN
                    ted.TD_NAZIVI.TD_DAJ_SIF('s_fup','NAZIV','ID',funkc,'Q') ||
                    CASE 
                        WHEN napon < 110 THEN ' POLJE br. ' || POLJE || ' (' || IME_PO || ') na ' || napon||' kV '
                        ELSE ' POLJE (' || IME_PO || ') na ' || napon||' kV '
                    END
                ELSE
                    CASE 
                        WHEN napon < 110 THEN IME_PO || ' POLJE br. ' || POLJE || ' na ' || napon||' kV '
                        ELSE IME_PO || ' POLJE na ' || napon||' kV '
                    END
            END
        ELSE ' '
    END dv_trafo_polje,

    ted.TD_NAZIVI.TD_DAJ_SIF('s_fup','SIFRA','ID',funkc,'Q') fup,

    A.OPIS,
    A.SM_PK,
    A.RAD_APU,
    A.VREM_USL,
    A.STATUS_ISPKV1,
    DOPUNA_ISPKV1,
    ID_STAVKE_I,
    Z_DSDF_GL1,Z_KVAR_GL1,Z_PRST_GL1,Z_ZMSP_GL1,Z_UZMS_GL1,
Z_RAPU_GL1,Z_LOKK_GL1,Z_DSDF_GL2,Z_KVAR_GL2,Z_PRST_GL2,Z_ZMSP_GL2,Z_UZMS_GL2,Z_RAPU_GL2,Z_LOKK_GL2,
Z_DSDF_GL3,Z_KVAR_GL3,Z_PRST_GL3,Z_ZMSP_GL3,Z_UZMS_GL3,Z_RAPU_GL3,Z_LOKK_GL3,
Z_DIS_REZ,Z_KVAR_REZ,Z_PRST_REZ,Z_ZMSP_REZ,
Z_DIS_REZ2,Z_KVAR_REZ2,Z_PRST_REZ2,Z_ZMSP_REZ2,
Z_PREK_VN,Z_PREK_NN,Z_NEL1,Z_NEL2, Z_NEL3, Z_SABZ_SAB,Z_SABZ_SAB2, Z_OTPR_SAB, Z_OTPR_SAB2,
Z_JPS_VN, Z_JPS_NN,
Z_JPS_VN2,Z_JPS_NN2,
ID_Z_TELE_POC_GL1, ID_Z_TELE_KRAJ_GL1, ID_Z_TELE_POC_GL2, ID_Z_TELE_KRAJ_GL2,
ID_Z_TELE_POC_GL3, ID_Z_TELE_KRAJ_GL3,
SNAGA snaga
FROM (
    SELECT 
        DOG_ISPKV.ID,
        DOG_ISPKV.RB,
        TO_CHAR(DOG_ISPKV.VREPOC,'hh24:mi') VREPOC,
        TO_CHAR(DOG_ISPKV.VREZAV,'hh24:mi') VREZAV,

        ted.TD_NAZIVI.TD_DAJ_SIF('S_VRPD','NAZIV','ID',DOG_ISPKV.ID_S_VRPD,'Q') VRSTADOG,
        ted.TD_NAZIVI.TD_DAJ_SIF('S_VRPD','SIFRA','ID',DOG_ISPKV.ID_S_VRPD,'Q') VR_DOG_SIF,

        ted.TD_NAZIVI.TD_DAJ_SIF('S_GRUZR','SKR','ID',DOG_ISPKV.ID1_S_GRUZR,'Q') GRUZR1,
        ted.TD_NAZIVI.TD_DAJ_SIF('S_UZROK','SKR','ID',DOG_ISPKV.ID1_S_UZROK,'Q') UZROK1,

        DOG_ISPKV.OPIS,

        ted.TD_NAZIVI.TD_DAJ_SIF('S_SM_PK','NAZIV','ID',DOG_ISPKV.ID_S_SM_PK,'Q') SM_PK,
        ted.TD_NAZIVI.TD_DAJ_SIF('S_RAPU','SKR_NAZ','ID',DOG_ISPKV.ID_S_RAPU,'Q') RAD_APU,
        ted.TD_NAZIVI.TD_DAJ_SIF('S_VREM_USL','NAZIV','ID',DOG_ISPKV.ID_S_VREM_USL,'Q') VREM_USL,

        ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',DOG_ISPKV.ID_TIPOB,'Q') tip_ob,

        pgi.PI_NAZIVI.NAZ_OBJ(DOG_ISPKV.ID_TIPOB,DOG_ISPKV.OB_ID,'OB_SIF','Q') OBJEKAT_SIFRA,
        pgi.PI_NAZIVI.NAZ_OBJ(DOG_ISPKV.ID_TIPOB,DOG_ISPKV.OB_ID,'NAZOB','Q') OBJEKAT_NAZIV,

        ted.TD_NAZIVI.TD_DAJ_SIF('S_NAP','NAZIV','ID',DOG_ISPKV.ID_S_NAP,'Q') NAPON,

        ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','POLJE','P2_TRAF_ID',DOG_ISPKV.P2_TRAF_ID,'Q') POLJE,
        ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','IME_PO','P2_TRAF_ID',DOG_ISPKV.P2_TRAF_ID,'Q') IME_PO,
        ted.TD_NAZIVI.TD_DAJ_ID('V_S_POLJE_SVA_AP','FUP_ID','P2_TRAF_ID',DOG_ISPKV.P2_TRAF_ID,'Q') FUNKC,

        DOG_ISPKV.OB_ID,
        DOG_ISPKV.TRAFO_ID,
        DOG_ISPKV.P2_TRAF_ID,

        DOG_ISPKV.STATUS STATUS_ISPKV1,
        DOG_ISPKV.DOPUNA DOPUNA_ISPKV1,
        DOG_ISPKV.VREPOC VREPOC_SORT_ISPKV,

        DOG_ISPKV.ID ID_STAVKE_I,
        ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_DSDF_GL1,'Q') Z_DSDF_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('PD_KK','NAZIV_EDD','ID',DOG_ISPKV.ID_Z_KVAR_GL1,'Q') Z_KVAR_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PRST_GL1,'Q') Z_PRST_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_ZMSP_GL1,'Q') Z_ZMSP_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_UZMS_GL1,'Q') Z_UZMS_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_RAPU','SKR_NAZ','ID',DOG_ISPKV.ID_Z_RAPU_GL1,'Q') Z_RAPU_GL1,
       DOG_ISPKV.Z_LOKK_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_DSDF_GL2,'Q') Z_DSDF_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('PD_KK','NAZIV_EDD','ID',DOG_ISPKV.ID_Z_KVAR_GL2,'Q') Z_KVAR_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PRST_GL2,'Q') Z_PRST_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_ZMSP_GL2,'Q') Z_ZMSP_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_UZMS_GL2,'Q') Z_UZMS_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_RAPU','SKR_NAZ','ID',DOG_ISPKV.ID_Z_RAPU_GL2,'Q') Z_RAPU_GL2,
       DOG_ISPKV.Z_LOKK_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_DSDF_GL3,'Q') Z_DSDF_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('PD_KK','NAZIV_EDD','ID',DOG_ISPKV.ID_Z_KVAR_GL3,'Q') Z_KVAR_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PRST_GL3,'Q') Z_PRST_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_ZMSP_GL3,'Q') Z_ZMSP_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_UZMS_GL3,'Q') Z_UZMS_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_RAPU','SKR_NAZ','ID',DOG_ISPKV.ID_Z_RAPU_GL3,'Q') Z_RAPU_GL3,
       DOG_ISPKV.Z_LOKK_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_DIS_REZ,'Q') Z_DIS_REZ,
       ted.TD_NAZIVI.TD_DAJ_SIF('PD_KK','NAZIV_EDD','ID',DOG_ISPKV.ID_Z_KVAR_REZ,'Q') Z_KVAR_REZ,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PRST_REZ,'Q') Z_PRST_REZ,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_ZMSP_REZ,'Q') Z_ZMSP_REZ,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_DIS_REZ2,'Q') Z_DIS_REZ2,
       ted.TD_NAZIVI.TD_DAJ_SIF('PD_KK','NAZIV_EDD','ID',DOG_ISPKV.ID_Z_KVAR_REZ2,'Q') Z_KVAR_REZ2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PRST_REZ2,'Q') Z_PRST_REZ2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_ZMSP_REZ2,'Q') Z_ZMSP_REZ2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PREK_VN,'Q') Z_PREK_VN,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_PREK_NN,'Q') Z_PREK_NN,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_NEL1,'Q') Z_NEL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_NEL2,'Q') Z_NEL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_NEL3,'Q') Z_NEL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_SABZ_SAB,'Q') Z_SABZ_SAB,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_SABZ_SAB2,'Q') Z_SABZ_SAB2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_OTPR_SAB,'Q') Z_OTPR_SAB,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_OTPR_SAB2,'Q') Z_OTPR_SAB2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_JPS_VN,'Q') Z_JPS_VN,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_JPS_NN,'Q') Z_JPS_NN,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_JPS_VN2,'Q') Z_JPS_VN2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZSI_INFO','NAZIV','ID',DOG_ISPKV.ID_Z_JPS_NN2,'Q') Z_JPS_NN2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_POC_GL1,'Q')ID_Z_TELE_POC_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_KRAJ_GL1,'Q') ID_Z_TELE_KRAJ_GL1,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_POC_GL2,'Q') ID_Z_TELE_POC_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_KRAJ_GL2,'Q') ID_Z_TELE_KRAJ_GL2,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_POC_GL3,'Q') ID_Z_TELE_POC_GL3,
       ted.TD_NAZIVI.TD_DAJ_SIF('ZAS_DV_TELE_V','INFOR_NAZIV','INFOR_ID',DOG_ISPKV.ID_Z_TELE_KRAJ_GL3,'Q') ID_Z_TELE_KRAJ_GL3,
SNAGA snaga

    FROM DOG_ISPKV
    WHERE DOG_ISPKV.ID_DOG_SMENE = :1
) A

ORDER BY A.VREPOC_SORT_ISPKV
`

	rowsi, err := m.DB.QueryContext(ctx, detailQueryI, id)
	if err != nil {
		return nil, err
	}
	defer rowsi.Close()

	var ispadi []models.Ispad

	for rowsi.Next() {
		var isp models.Ispad

		err := rowsi.Scan(
			&isp.VrepocSortIspkv,
			&isp.ID,
			&isp.RB,
			&isp.Vrepoc,
			&isp.Vrezav,
			&isp.VrstaDog,
			&isp.VrDogSif,
			&isp.Gruzr1,
			&isp.Uzrok1,
			&isp.Objekat,
			&isp.TipOb,
			&isp.Napon,
			&isp.DvTrafoPolje,
			&isp.Fup,
			&isp.Opis,
			&isp.SmPk,
			&isp.RadApu,
			&isp.VremUsl,
			&isp.StatusIspkv1,
			&isp.DopunaIspkv1,
			&isp.IDStavkeI,

			// GL1
			&isp.ZDsdfGl1,
			&isp.ZKvarGl1,
			&isp.ZPrstGl1,
			&isp.ZZMSPGl1,
			&isp.ZUzmsGl1,
			&isp.ZRapuGl1,
			&isp.ZLokkGl1,

			// GL2
			&isp.ZDsdfGl2,
			&isp.ZKvarGl2,
			&isp.ZPrstGl2,
			&isp.ZZMSPGl2,
			&isp.ZUzmsGl2,
			&isp.ZRapuGl2,
			&isp.ZLokkGl2,

			// GL3
			&isp.ZDsdfGl3,
			&isp.ZKvarGl3,
			&isp.ZPrstGl3,
			&isp.ZZMSPGl3,
			&isp.ZUzmsGl3,
			&isp.ZRapuGl3,
			&isp.ZLokkGl3,

			// REZ
			&isp.ZDisRez,
			&isp.ZKvarRez,
			&isp.ZPrstRez,
			&isp.ZZMSPRez,

			&isp.ZDisRez2,
			&isp.ZKvarRez2,
			&isp.ZPrstRez2,
			&isp.ZZMSPRez2,

			// ostalo
			&isp.ZPrekVn,
			&isp.ZPrekNn,
			&isp.ZNel1,
			&isp.ZNel2,
			&isp.ZNel3,
			&isp.ZSabzSab,
			&isp.ZSabzSab2,
			&isp.ZOtprSab,
			&isp.ZOtprSab2,
			&isp.ZJpsVn,
			&isp.ZJpsNn,
			&isp.ZJpsVn2,
			&isp.ZJpsNn2,

			// tele
			&isp.IdZTelePocGl1,
			&isp.IdZTeleKrajGl1,
			&isp.IdZTelePocGl2,
			&isp.IdZTeleKrajGl2,
			&isp.IdZTelePocGl3,
			&isp.IdZTeleKrajGl3,

			&isp.Snaga,
		)
		if err != nil {
			return nil, err
		}

		ispadi = append(ispadi, isp)
	}

	if err := rowsi.Err(); err != nil {
		return nil, err
	}

	if len(ispadi) > 0 {
		d.Ispad = &ispadi
	}

	/**** Kraj dela za ispad ****/

	/*** Deo za manipulacije ***/

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

	if len(manipulacije) > 0 {
		d.Manipulacije = manipulacije
	}

	/**** Kraj dela za manipulacije ****/

	return &d, nil
}

func (m *OracleDBRepo) GetPrekidPById(ctx context.Context, id int) (*models.DogadjajDetaljno, error) {

	//  MASTER QUERY
	masterQuery := `
SELECT 
  d.id,
  d.rb_dog,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_DOG','TIP','ID', d.id_tip_dog,'Q') AS tip_dog,
  d.naslov,
  d.ID_SMENA,
  d.ID_DOG_SMENE AS VEZA_SA,
  d2.rb_dog AS RB_DOG_VEZA_SA,
  s3.DATDNEV datum_veze,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s3.id_tip_smena,'Q') AS tip_smene_veze,
  d.DOPUNA,
  d.ID_SMENA_D,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_GRRAZ','NAZIV','ID', 
    ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','ID_S_GRRAZ','ID', d.ID_S_RAZLOG,'Q'),'Q') AS grazlog,

  ted.TD_NAZIVI.TD_DAJ_SIF('S_RAZLOG','NAZIV','ID', d.ID_S_RAZLOG,'Q') AS razlog,

  d.id_s_razlog,
  d.uzrok_tekst,
  d.man_tekst,
  d.POSL_TEKST,

  s1.DATDNEV AS datum_smene,
  s2.DATDNEV AS datum_dopune,

  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s1.id_tip_smena,'Q') AS tip_smene,
  ted.TD_NAZIVI.TD_DAJ_SIF('TIP_SMENA','SKR_NAZ','ID', s2.id_tip_smena,'Q') AS tip_smene_dopune

FROM dog_smene d
JOIN smena s1 
  ON d.ID_SMENA = s1.id
LEFT JOIN smena s2 
  ON d.ID_SMENA_D = s2.id
LEFT JOIN dog_smene d2 
  ON d.id_dog_smene = d2.id
LEFT JOIN smena s3 
  ON d2.ID_SMENA = s3.id
WHERE d.id = :1
`

	var d models.DogadjajDetaljno

	row := m.DB.QueryRowContext(ctx, masterQuery, id)

	err := row.Scan(
		&d.ID,
		&d.RbDog,
		&d.TipDog,
		&d.Naslov,
		&d.IDSmena,
		&d.VezaSa,
		&d.RbDogVezaSa,
		&d.DatumVezaSa,
		&d.TipSmenaVezaSa,
		&d.Dopuna,
		&d.IDSmenaD,
		&d.Grazlog,
		&d.Razlog,
		&d.IDSrazlog,
		&d.UzrokTekst,
		&d.ManTekst,
		&d.Posledice,
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
SELECT A.VREPOC_SORT_ISPKV5678,A.POLJE,A.FUP_ID,A.FUNKC,A.OB_ID,A.TIP_OB,A.ID_P2_TRAF,
A.VREPOC,A.VREZAV,
  case when A.tip_ob in ('DV','TD','KB','TK') THEN decode(A.ID_P2_TRAF,null,'',TD_NAZIVI.COL_V_S_OB_P2_TRAF(ID_P2_TRAF,'OPIS','Q'))
     else OBJEKAT_NAZIV
     end objekat,
    
     A.GENERATOR,A.VR_PREK,A.UZROK_PREK,A.TIP_PREK,a.SNAGA,A.OPIS,A.DOPUNA_ISPKV5678, STATUS
FROM (SELECT 
       to_char(DOG_PREKID_P.vrepoc,'hh24:mi') vrepoc,to_char(DOG_PREKID_P.vrezav,'hh24:mi') vrezav,
       DOG_PREKID_P.OPIS,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_TIPOB','SIFRA','ID',DOG_PREKID_P.ID_TIPOB,'Q') tip_ob,
       pgi.PI_NAZIVI.NAZ_OBJ(DOG_PREKID_P.ID_TIPOB,DOG_PREKID_P.OB_ID,'OB_SIF','Q') OBJEKAT_SIFRA,
       pgi.PI_NAZIVI.NAZ_OBJ(DOG_PREKID_P.ID_TIPOB,DOG_PREKID_P.OB_ID,'OPIS','Q') OBJEKAT_NAZIV,
       ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','POLJE','P2_TRAF_ID',DOG_PREKID_P.ID_P2_TRAF,'Q') POLJE,
       ted.td_nazivi.td_daj_sif('V_S_POLJE_SVA_AP','IME_PO','P2_TRAF_ID',DOG_PREKID_P.ID_P2_TRAF,'Q') GENERATOR,  
       ted.TD_NAZIVI.TD_DAJ_ID('V_S_POLJE_SVA_AP','FUP_ID','P2_TRAF_ID',DOG_PREKID_P.ID_P2_TRAF,'Q') FUP_ID,
       ted.TD_NAZIVI.TD_DAJ_SIF('V_S_POLJE_SVA_AP','FUNKC','P2_TRAF_ID',DOG_PREKID_P.ID_P2_TRAF,'Q') FUNKC,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_VR_PREK','NAZIV','ID',DOG_PREKID_P.ID_S_VR_PREK,'Q') VR_PREK,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_UZROK_PREK','NAZIV','ID',DOG_PREKID_P.ID_S_UZROK_PREK,'Q') UZROK_PREK,
       ted.TD_NAZIVI.TD_DAJ_SIF('S_TIP_PREK','NAZIV','ID',DOG_PREKID_P.ID_S_TIP_PREK,'Q') TIP_PREK,
       
       DOG_PREKID_P.ob_id,DOG_PREKID_P.ID_P2_TRAF,
DOG_PREKID_P.vrepoc VREPOC_SORT_ISPKV5678,
DOG_PREKID_P.DOPUNA DOPUNA_ISPKV5678,DOG_PREKID_P.SNAGA, STATUS
FROM DOG_PREKID_P
WHERE DOG_PREKID_P.ID_DOG_SMENE=:1
) A
order by A.VREPOC_SORT_ISPKV5678
`

	rows, err := m.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pps []models.PrekidP

	for rows.Next() {
		var pp models.PrekidP

		err := rows.Scan(
			&pp.VrepocSort,
			&pp.Polje,
			&pp.FupID,
			&pp.Funkc,
			&pp.ObID,
			&pp.TipOb,
			&pp.IDP2Traf,

			&pp.Vrepoc,
			&pp.Vrezav,

			&pp.Objekat,

			&pp.Generator,
			&pp.VrPrek,
			&pp.UzrokPrek,
			&pp.TipPrek,

			&pp.Snaga,
			&pp.Opis,
			&pp.Dopuna,
			&pp.Status,
		)
		if err != nil {
			return nil, err
		}

		pps = append(pps, pp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(pps) > 0 {
		d.PrekidP = &pps
	}

	return &d, nil
}

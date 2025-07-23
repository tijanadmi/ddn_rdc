package oraclerepo

import (
	"context"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"
)

func (m *OracleDBRepo) GetPGDRadapuMes(ctx context.Context, arg models.ListPGD) ([]*models.PGDRadapuMes, error) {

	query := `select a.mesec,td_nazivi.td_daj_sif('s_meseci','naziv','broj',a.mesec,'Q') mesec_naz,
         a.godina,a.napon,decode(c.broj_apu,null,0,c.broj_apu) broj_apu,decode(b.broj_apu_n,null,0,b.broj_apu_n) broj_apu_n
from (select distinct EXTRACT(MONTH from datizv) mesec,EXTRACT(year from datizv) godina,s_nap.naziv napon,s_nap.id
      from s_nap,pi_dd
      where naziv > 100
        and EXTRACT(year from datizv) like :1) a,
     (SELECT EXTRACT(MONTH from DATIZV) mesec,EXTRACT(year from DATIZV) godina,v_s_ob.nn_id napon,COUNT (*) broj_apu_n
      FROM PI_dd,v_s_ob,s_tipob
      WHERE pi_dd.id_tipob = V_S_OB.TIPOB
        and pi_dd.id_tipob = s_tipob.id
        and s_tipob.sifra not in ('TD','TS')
        and pi_dd.OB_ID = v_s_ob.ob_id
        and ID_S_TIPD = td_nazivi.td_daj_id('s_tipd','id','sifra',1,'Q')
        AND ID_S_VRPD = (select id
                          from s_vrpd
                          where id_s_tipd = pi_dd.id_s_tipd
                            and sifra = '6')
        and poc_pp is null
        and zav_pp is null
        and id1 is not null
        and EXTRACT(year from DATIZV) like :2
       GROUP BY EXTRACT(MONTH from DATIZV),EXTRACT(year from DATIZV),v_s_ob.nn_id) b,
     (SELECT EXTRACT(MONTH from DATIZV) mesec,EXTRACT(year from DATIZV) godina,v_s_ob.nn_id napon,COUNT (*) broj_apu
      FROM pi_dd,v_s_ob,s_tipob
      WHERE pi_dd.id_tipob = V_S_OB.TIPOB
        and pi_dd.id_tipob = s_tipob.id
        and s_tipob.sifra not in ('TD','TS')
        and pi_dd.OB_ID = v_s_ob.ob_id
        and ID_S_TIPD = td_nazivi.td_daj_id('s_tipd','id','sifra',1,'Q')
        AND ID_S_VRPD = (select id
                          from s_vrpd
                          where id_s_tipd = pi_dd.id_s_tipd
                            and sifra = '1')
        and poc_pp is null
        and zav_pp is null
        and id1 is not null
        and EXTRACT(year from DATIZV) like :3
       GROUP BY EXTRACT(MONTH from DATIZV),EXTRACT(year from DATIZV),v_s_ob.nn_id) c       
where a.mesec = b.mesec(+)
  and a.godina = b.godina(+)
  and a.id = b.napon(+)
  and a.mesec = c.mesec(+)
  and a.godina = c.godina(+)
  and a.id = c.napon(+)
order by godina desc,mesec asc,napon desc`

	// fmt.Println(arg.Ind, arg.Mrc, arg.StartDate, arg.EndDate, arg.Offset,arg.Limit)
	rows, err := m.DB.QueryContext(ctx, query, arg.Godina, arg.Godina, arg.Godina)

	//fmt.Println(arg.StartDate, arg.EndDate, arg.Tipd)
	if err != nil {
		fmt.Println("Pogresan upit ili nema rezultata upita")
		return nil, err
	}
	defer rows.Close()

	var ues []*models.PGDRadapuMes

	for rows.Next() {
		var ue models.PGDRadapuMes
		err := rows.Scan(
			&ue.MesecBr,
			&ue.Mesec,
			&ue.Godina,
			&ue.Napon,
			&ue.BrojApu,
			&ue.BrojApuN,
		)

		if err != nil {
			return nil, err
		}

		ues = append(ues, &ue)

	}

	//fmt.Println(totalCount)
	return ues, nil
}

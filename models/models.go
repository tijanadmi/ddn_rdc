package models

// Zastita is the type for all signals of protective devices
type Signal struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type MalfunctionIn struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Order string `json:"order"`
}

type Apu struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Code      string `json:"code"`
	Status    string `json:"status"`
}

type GroupOfCause struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Status    string `json:"status"`
	DDNCode   string `json:"ddn_code"`
	Sort      int    `json:"sort"`
}

type Cause struct {
	ID           int    `json:"id"`
	GroupCauseId int    `json:"-"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	Status       string `json:"status"`
	Sort         int    `json:"sort"`
	GroupOfCause GroupOfCause
}

type GroupOfReason struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Status    string `json:"status"`
	Sort      int    `json:"sort"`
}

type Reason struct {
	ID            int    `json:"id"`
	GroupReasonId int    `json:"-"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	ShortName     string `json:"short_name"`
	Status        string `json:"status"`
	Sort          int    `json:"sort"`
	GroupOfReason GroupOfReason
}

type WeatherCondition struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CategoryOfEvent struct {
	ID          int    `json:"id"`
	TypeEventId int    `json:"event_type_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Status      string `json:"status"`
	Sort        int    `json:"sort"`
}

type OverheadLine struct {
	ID                      int    `json:"tis_id_ohl"`
	IpsIdDV                 string `json:"ips_id_ohl"`
	SapIdDV                 string `json:"sap_id_ohl"`
	IdDal                   string `json:"code_ohl"`
	IPSName                 string `json:"ips_name_ohl"`
	CategoryNameDV          string `json:"category_name_ohl"`
	IdNN                    int    `json:"id_voltage_level"`
	NNName                  string `json:"name_voltage_level"`
	CategoryId              string `json:"category_id"`
	CategoryName            string `json:"category_name"`
	TisIdFirstFeeder        int    `json:"tis_id_first_feeder"`
	IpsIdFirstFeeder        string `json:"ips_id_first_feeder"`
	SapIdFirstFeeder        string `json:"sap_id_first_feeder"`
	FirstFeederNumber       string `json:"first_feeder_number"`
	FirstFeederName         string `json:"first_feeder_name"`
	FirstFeederCategoryName string `json:"first_feeder_category_name"`
	TisIdLastFeeder         int    `json:"tis_id_last_feeder"`
	IpsIdLastFeeder         string `json:"ips_id_last_feeder"`
	SapIdLastFeeder         string `json:"sap_id_last_feeder"`
	LastFeederNumber        string `json:"last_feeder_number"`
	LastFeederName          string `json:"last_feeder_name"`
	LastFeederCategoryName  string `json:"last_feeder_category_name"`
	ObjTypeId               int    `json:"object_type_id"`
	ObjTypeCode             string `json:"object_type_code"`
	ObjTypeName             string `json:"object_type_name"`
	Mrc1Id                  string `json:"rdc1_id"`
	Mrc1Name                string `json:"rdc1_name"`
	Mrc2Id                  string `json:"rdc2_id"`
	Mrc2Name                string `json:"rdc2_name"`
}

type PowerCable struct {
	ID                      int    `json:"tis_id_cabel"`
	IpsIdKB                 string `json:"ips_id_cabel"`
	SapIdKB                 string `json:"sap_id_cabel"`
	IdKab                   string `json:"code_cabel"`
	IPSName                 string `json:"ips_name_cabel"`
	CategoryNameDV          string `json:"category_name_cabel"`
	CategoryId              string `json:"category_id"`
	CategoryName            string `json:"category_name"`
	TisIdFirstFeeder        int    `json:"tis_id_first_feeder"`
	IpsIdFirstFeeder        string `json:"ips_id_first_feeder"`
	SapIdFirstFeeder        string `json:"sap_id_first_feeder"`
	FirstFeederNumber       string `json:"first_feeder_number"`
	FirstFeederName         string `json:"first_feeder_name"`
	FirstFeederCategoryName string `json:"first_feeder_category_name"`
	TisIdLastFeeder         int    `json:"tis_id_last_feeder"`
	IpsIdLastFeeder         string `json:"ips_id_last_feeder"`
	SapIdLastFeeder         string `json:"sap_id_last_feeder"`
	LastFeederNumber        string `json:"last_feeder_number"`
	LastFeederName          string `json:"last_feeder_name"`
	LastFeederCategoryName  string `json:"last_feeder_category_name"`
	ObjTypeId               int    `json:"object_type_id"`
	ObjTypeCode             string `json:"object_type_code"`
	ObjTypeName             string `json:"object_type_name"`
	Mrc1Id                  string `json:"rdc1_id"`
	Mrc1Name                string `json:"rdc1_name"`
	Mrc2Id                  string `json:"rdc2_id"`
	Mrc2Name                string `json:"rdc2_name"`
}

type Substation struct {
	ID              int    `json:"tis_id"`
	IpsId           string `json:"ips_id"`
	SapId           string `json:"sap_id"`
	TISName         string `json:"name_sub"`
	CategoryNameSub string `json:"category_name_sub"`
	ObjTypeId       int    `json:"object_type_id"`
	ObjTypeCode     string `json:"object_type_code"`
	ObjTypeName     string `json:"object_type_name"`
	Mrc1Id          int    `json:"rdc1_id"`
	Mrc1Name        string `json:"rdc1_name"`
	Mrc2Id          string `json:"rdc2_id"`
	Mrc2Name        string `json:"rdc2_name"`
	OrgId           int    `json:"org_id"`
	OrgName         string `json:"org_name"`
}

type Feeder struct {
	ID                 int    `json:"tis_id_feeder"`
	IpsIdFeeder        string `json:"ips_id_feeder"`
	SapIdFeeder        string `json:"sap_id_feeder"`
	TisIdSub           int    `json:"tis_id_substation"`
	IpsIdSub           string `json:"ips_id_substation"`
	SapIdSub           string `json:"sap_id_substation"`
	NameSub            string `json:"name_substation"`
	IdNN               int    `json:"id_voltage_level"`
	NNName             string `json:"name_voltage_level"`
	FeederNumber       string `json:"feeder_number"`
	FeederName         string `json:"feeder_name"`
	FeederCategoryName string `json:"feeder_category_name"`
	CategoryId         string `json:"category_id"`
	CategoryName       string `json:"category_name"`
	FeederFunId        string `json:"feeder_fun_id"`
	FeederFunName      string `json:"feeder_fun_name"`
	Equipped           string `json:"equipped"`
	Active             string `json:"active"`
	Completely         string `json:"completely"`
	FeederType         string `json:"feeder_type"`
}

type ProtectionDevice struct {
	ID           int    `json:"tis_id_prt"`
	IpsIdDV      string `json:"ips_id_prt"`
	TisIdFeeder  int    `json:"tis_id_feeder"`
	IpsIdFeeder  string `json:"ips_id_feeder"`
	SapIdFFeeder string `json:"sap_id_feeder"`
	Manufacturer string `json:"manufacturer"`
	Type         string `json:"type"`
	TypeOfEqu    string `json:"type_of_equ"`
	Indicator    string `json:"indicator"`
	Technology   string `json:"tehnology"`
	StatusIPS    string `json:"status_ips"`
	Status       string `json:"status"`
	SerialNumber string `json:"serial_number"`
	DeviceTag    string `json:"device_tag"`
	APU          string `json:"apu"`
}

type PowerTransformer struct {
	ID              int    `json:"tis_id_tr"`
	IpsIdDV         string `json:"ips_id_tr"`
	SapIdDV         string `json:"sap_id_tr"`
	TisIdSub        int    `json:"tis_id_substation"`
	IpsIdSub        string `json:"ips_id_substation"`
	SapIdSub        string `json:"sap_id_substation"`
	NameSub         string `json:"name_substation"`
	NameTr          string `json:"name_tr"`
	CategoryNameTr  string `json:"category_name_tr"`
	CategoryId      string `json:"category_id"`
	CategoryName    string `json:"category_name"`
	Un              string `json:"un"`
	TisIdPrimFeeder int    `json:"tis_id_prim_feeder"`
	IpsIdPrimFeeder string `json:"ips_id_prim_feeder"`
	SapIdPrimFeeder string `json:"sap_id_prim_feeder"`
	PrimFeederName  string `json:"prim_feeder_name"`
	TisIdSec1Feeder string `json:"tis_id_sec1_feeder"`
	IpsIdSec1Feeder string `json:"ips_id_sec1_feeder"`
	SapIdSec1Feeder string `json:"sap_id_sec1_feeder"`
	Sec1FeederName  string `json:"sec1_feeder_name"`
	TisIdSec2Feeder string `json:"tis_id_sec2_feeder"`
	IpsIdSec2Feeder string `json:"ips_id_sec2_feeder"`
	SapIdSec2Feeder string `json:"sap_id_sec2_feeder"`
	Sec2FeederName  string `json:"sec2_feeder_name"`
	TisIdTerFeeder  string `json:"tis_id_ter_feeder"`
	IpsIdTerFeeder  string `json:"ips_id_ter_feeder"`
	SapIdTerFeeder  string `json:"sap_id_ter_feeder"`
	TerFeederName   string `json:"ter_feeder_name"`
}

type Disconnector struct {
	ID              int    `json:"tis_id_dis"`
	IpsIdDs         string `json:"ips_id_dis"`
	SapIdDs         string `json:"sap_id_dis"`
	TisIdFeeder     int    `json:"tis_id_feeder"`
	IpsIdFeeder     string `json:"ips_id_feeder"`
	SapIdFeeder     string `json:"sap_id_feeder"`
	FeederName      string `json:"feeder_name"`
	DisCategoryName string `json:"dis_category_name"`
	FunDisId        string `json:"fun_dis_id"`
	FunDis          string `json:"fun_dis"`
	CategoryId      string `json:"category_id"`
	CategoryName    string `json:"category_name"`
}

type WorkPermission struct {
	Code               int    `json:"code"`
	PerNum1            string `json:"per_num1"`
	PerNum2            string `json:"per_num2"`
	ScheduledTimeStart string `json:"scheduled_time_start"`
	ScheduledDateStart string `json:"scheduled_date_start"`
	Note               string `json:"note"`
	Status             string `json:"status"`
	EndTime            string `json:"ens_time"`
	EndDate            string `json:"end_date"`
}

type WorkPermissionAll struct {
	IdZahteva       int    `json:"id_zahteva"`
	Grupa           string `json:"grupa"`
	Ukljucenost     string `json:"ukljucenost"`
	IntPl           string `json:"int_pl"`
	Br1Z1Gr         string `json:"br1_z_1gr"`
	Br2Z1Gr         string `json:"br2_z_1gr"`
	BrZRDC2Gr       string `json:"br_z_rdc_2gr"`
	BrSag           string `json:"br_sag"`
	PlDatumOdZ      string `json:"pl_datum_od_z"`
	PlVremeOdZ      string `json:"pl_vreme_od_z"`
	PlDatumDoZ      string `json:"pl_datum_do_z"`
	PlVremeDoZ      string `json:"pl_vreme_do_z"`
	RukRadova       string `json:"ruk_radova"`
	OpisRadova      string `json:"opis_radova"`
	NapomenaVeza    string `json:"napomena_veza"`
	SagUslovi       string `json:"sag_uslovi"`
	SagNapomenaVeza string `json:"sag_napomena_veza"`
	IdDozvole       string `json:"id_dozvole"`
	Br1Doz1Gr       string `json:"br1_doz_1gr"`
	Br1Doz2Gr       string `json:"br1_doz_2gr"`
	Br2Doz          string `json:"br2_doz"`
	PlVremeOd       string `json:"pl_vreme_od"`
	PlVremeDo       string `json:"pl_vreme_do"`
	ZavVreme        string `json:"zav_vreme"`
	ZavDatum        string `json:"zav_datum"`
	Status          string `json:"status"`
}

type WorkInEENetwork struct {
	MaxNum      string `json:"max_num"`
	Num         string `json:"number"`
	EEElements  string `json:"ee_elements"`
	Workplace   string `json:"workplace"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Time        string `json:"time"`
	Link        string `json:"link"`
}

type UnopenedPermit struct {
	BrojIsk    string `json:"broj_isk"`
	BrojIsk2   string `json:"broj_isk2"`
	BrojIskRDC string `json:"broj_isk_rdc"`
	RbrIsk     string `json:"rbr_isk"`
	RazlogO    string `json:"razlog_o"`
	DatumO     string `json:"datum_o"`
	Rbr        string `json:"rbr"`
	Korisnik   string `json:"korisnik"`
	Opis       string `json:"opis"`
	LiceRuk    string `json:"lice_ruk"`
}

type PiPiDDNIspad struct {
	DatSmene       string `json:"dat_smene"`
	IdSMrc         string `json:"id_s_mrc"`
	IdSTipd        string `json:"id_s_tipd"`
	IdSVrpd        string `json:"id_s_vrpd"`
	IdRadAPU       string `json:"id_s_rpu"`
	IdTipob        string `json:"id_tipob"`
	ObId           string `json:"ob_id"`
	TrafoId        string `json:"trafo_id"`
	Vrepoc         string `json:"vrepoc"`
	Vrezav         string `json:"vrezav"`
	Id1SGruzr      string `json:"id_s_gruzr"`
	Id1SUzrok      string `json:"id_s_uzrok"`
	Snaga          string `json:"snaga"`
	Opis           string `json:"opis"`
	IdSNap         string `json:"id_s_nap"`
	P2TrafId       string `json:"p2_traf_id"`
	KorUneo        string `json:"kor_uneo"`
	IdZDsdfGL1     string `json:"id_z_dsdf_gl1"`
	IdZKvarGL1     string `json:"id_z_kvar_gl1"`
	IdZRapuGL1     string `json:"id_z_rapu_gl1"`
	IdZPrstGL1     string `json:"id_z_prst_gl1"`
	IdZZmspGL1     string `json:"id_z_zmsp_gl1"`
	IdZUzmsGL1     string `json:"id_z_uzms_gl1"`
	ZLokkGL1       string `json:"z_lokk_gl1"`
	IdZDsdfGL2     string `json:"id_z_dsdf_gl2"`
	IdZKvarGL2     string `json:"id_z_kvar_gl2"`
	IdZRapuGL2     string `json:"id_z_rapu_gl2"`
	IdZPrstGL2     string `json:"id_z_prst_gl2"`
	IdZZmspGL2     string `json:"id_z_zmsp_gl2"`
	IdZUzmsGL2     string `json:"id_z_uzms_gl2"`
	ZLokkGL2       string `json:"z_lokk_gl2"`
	IdZPrekVN      string `json:"id_z_prek_vn"`
	IdZDisREZ      string `json:"id_z_dis_rez"`
	IdZKvarREZ     string `json:"id_z_kvar_rez"`
	IdZPrstREZ     string `json:"id_z_prst_rez"`
	IdZZmspREZ     string `json:"id_z_zmsp_rez"`
	IdZNel1        string `json:"id_z_nel1"`
	IdZNel2        string `json:"id_z_nel2"`
	IdZNel3        string `json:"id_z_nel3"`
	IdZPrekNN      string `json:"id_z_prek_nn"`
	IdZSabzSAB     string `json:"id_z_sabz_sab"`
	IdZOtprSAB     string `json:"id_z_otpr_sab"`
	IdSVremUSL     string `json:"id_s_vrem_usl"`
	UzrokTekst     string `json:"uzrok_tekst"`
	IdZJpsVN       string `json:"id_z_jps_vn"`
	IdZJpsNN       string `json:"id_z_jsp_vn"`
	PoslTekst      string `json:"posl_tekst"`
	IdZTelePocGL1  string `json:"id_s_z_tele_poc_gl1"`
	IdZTeleKrajGL1 string `json:"id_z_tele_kraj_gl1"`
	IdZTelePocGL2  string `json:"id_z_tele_poc_gl2"`
	IdZTeleKrajGL2 string `json:"id_z_tele_kraj_gl2"`
	SynsoftId      string `json:"ed_id"`
}

type PiPiDDNIsklj struct {
	DatSmene  string `json:"dat_smene"`
	IdSMrc    string `json:"id_s_mrc"`
	TipMan    string `json:"tip_man"`
	IdTipob   string `json:"id_tipob"`
	ObId      string `json:"ob_id"`
	TrafoId   string `json:"trafo_id"`
	Vrepoc    string `json:"vrepoc"`
	Vrezav    string `json:"vrezav"`
	IdSGrraz  string `json:"id_s_grraz"`
	IdSRazlog string `json:"id_s_razlog"`
	ManTekst  string `json:"man_tekst"`
	IdSNap    string `json:"id_s_nap"`
	P2TrafId  string `json:"p2_traf_id"`
	KorUneo   string `json:"kor_uneo"`
	SynsoftId string `json:"ed_id"`
}

type PiPiDDN struct {
	Datizv         string `json:"datizv"`
	IdSMrc         string `json:"id_s_mrc"`
	IdSTipd        string `json:"id_s_tipd"`
	IdSVrpd        string `json:"id_s_vrpd"`
	IdTipob        string `json:"id_tipob"`
	ObId           string `json:"ob_id"`
	TrafoId        string `json:"trafo_id"`
	Vrepoc         string `json:"vrepoc"`
	PocPP          string `json:"poc_pp"`
	Vrezav         string `json:"vrezav"`
	ZavPP          string `json:"zav_pp"`
	Id1SGruzr      string `json:"id_s_gruzr"`
	Id1SUzrok      string `json:"id_s_uzrok"`
	IdSGrraz       string `json:"id_s_grraz"`
	IdSRazlog      string `json:"id_s_razlog"`
	Snaga          string `json:"snaga"`
	Opis           string `json:"opis"`
	IdSNap         string `json:"id_s_nap"`
	P2TrafId       string `json:"p2_traf_id"`
	KorUneo        string `json:"kor_uneo"`
	Status         string `json:"status"`
	Datpri         string `json:"datpri"`
	IdZDsdfGL1     string `json:"id_z_dsdf_gl1"`
	IdZKvarGL1     string `json:"id_z_kvar_gl1"`
	IdZRapuGL1     string `json:"id_z_rapu_gl1"`
	IdZPrstGL1     string `json:"id_z_prst_gl1"`
	IdZZmspGL1     string `json:"id_z_zmsp_gl1"`
	IdZUzmsGL1     string `json:"id_z_uzms_gl1"`
	ZLokkGL1       string `json:"z_lokk_gl1"`
	IdZDsdfGL2     string `json:"id_z_dsdf_gl2"`
	IdZKvarGL2     string `json:"id_z_kvar_gl2"`
	IdZRapuGL2     string `json:"id_z_rapu_gl2"`
	IdZPrstGL2     string `json:"id_z_prst_gl2"`
	IdZZmspGL2     string `json:"id_z_zmsp_gl2"`
	IdZUzmsGL2     string `json:"id_z_uzms_gl2"`
	ZLokkGL2       string `json:"z_lokk_gl2"`
	IdZPrekVN      string `json:"id_z_prek_vn"`
	IdZDisREZ      string `json:"id_z_dis_rez"`
	IdZKvarREZ     string `json:"id_z_kvar_rez"`
	IdZPrstREZ     string `json:"id_z_prst_rez"`
	IdZZmspREZ     string `json:"id_z_zmsp_rez"`
	IdZNel1        string `json:"id_z_nel1"`
	IdZNel2        string `json:"id_z_nel2"`
	IdZNel3        string `json:"id_z_nel3"`
	IdZPrekNN      string `json:"id_z_prek_nn"`
	IdZSabzSAB     string `json:"id_z_sabz_sab"`
	IdZOtprSAB     string `json:"id_z_otpr_sab"`
	IdSVremUSL     string `json:"id_s_vrem_usl"`
	UzrokTekst     string `json:"uzrok_tekst"`
	IdZJpsVN       string `json:"id_z_jps_vn"`
	IdZJpsNN       string `json:"id_z_jsp_vn"`
	PoslTekst      string `json:"posl_tekst"`
	IdZTelePocGL1  string `json:"id_s_z_tele_poc_gl1"`
	IdZTeleKrajGL1 string `json:"id_z_tele_kraj_gl1"`
	IdZTelePocGL2  string `json:"id_z_tele_poc_gl2"`
	IdZTeleKrajGL2 string `json:"id_z_tele_kraj_gl2"`
	SynsoftId      string `json:"ed_id"`
}

type UnfinishedEvents struct {
	Datizv         string `json:"datizv"`
	IdSMrc         string `json:"-"`
	IdSTipd        string `json:"id_s_tipd"`
	IdSVrpd        string `json:"id_s_vrpd"`
	IdTipob        string `json:"id_tipob"`
	ObId           string `json:"ob_id"`
	TrafoId        string `json:"trafo_id"`
	Vrepoc         string `json:"vrepoc"`
	Vrezav         string `json:"vrezav"`
	Id1SGruzr      string `json:"id_s_gruzr"`
	Id1SUzrok      string `json:"id_s_uzrok"`
	Snaga          string `json:"snaga"`
	Opis           string `json:"opis"`
	IdSNap         string `json:"id_s_nap"`
	P2TrafId       string `json:"p2_traf_id"`
	PgiKor         string `json:"pgi_kor"`
	IdZDsdfGL1     string `json:"id_z_dsdf_gl1"`
	IdZKvarGL1     string `json:"id_z_kvar_gl1"`
	IdZRapuGL1     string `json:"id_z_rapu_gl1"`
	IdZPrstGL1     string `json:"id_z_prst_gl1"`
	IdZZmspGL1     string `json:"id_z_zmsp_gl1"`
	IdZUzmsGL1     string `json:"id_z_uzms_gl1"`
	ZLokkGL1       string `json:"z_lokk_gl1"`
	IdZDsdfGL2     string `json:"id_z_dsdf_gl2"`
	IdZKvarGL2     string `json:"id_z_kvar_gl2"`
	IdZRapuGL2     string `json:"id_z_rapu_gl2"`
	IdZPrstGL2     string `json:"id_z_prst_gl2"`
	IdZZmspGL2     string `json:"id_z_zmsp_gl2"`
	IdZUzmsGL2     string `json:"id_z_uzms_gl2"`
	ZLokkGL2       string `json:"z_lokk_gl2"`
	IdZPrekVN      string `json:"id_z_prek_vn"`
	IdZDisREZ      string `json:"id_z_dis_rez"`
	IdZKvarREZ     string `json:"id_z_kvar_rez"`
	IdZPrstREZ     string `json:"id_z_prst_rez"`
	IdZZmspREZ     string `json:"id_z_zmsp_rez"`
	IdZNel1        string `json:"id_z_nel1"`
	IdZNel2        string `json:"id_z_nel2"`
	IdZNel3        string `json:"id_z_nel3"`
	IdZPrekNN      string `json:"id_z_prek_nn"`
	IdZSabzSAB     string `json:"id_z_sabz_sab"`
	IdZOtprSAB     string `json:"id_z_otpr_sab"`
	IdSVremUSL     string `json:"id_s_vrem_usl"`
	UzrokTekst     string `json:"uzrok_tekst"`
	IdZJpsVN       string `json:"id_z_jps_vn"`
	IdZJpsNN       string `json:"id_z_jsp_vn"`
	PoslTekst      string `json:"posl_tekst"`
	IdZTelePocGL1  string `json:"id_s_z_tele_poc_gl1"`
	IdZTeleKrajGL1 string `json:"id_z_tele_kraj_gl1"`
	IdZTelePocGL2  string `json:"id_z_tele_poc_gl2"`
	IdZTeleKrajGL2 string `json:"id_z_tele_kraj_gl2"`
	SynsoftId      string `json:"ed_id"`
}

type UnfinishedEventsUpdate struct {
	DatSmene       string `json:"dat_smene"`
	TipSmene       string `json:"tip_smene"`
	Vrepoc         string `json:"vrepoc"`
	Vrezav         string `json:"vrezav"`
	Id1SGruzr      string `json:"id_s_gruzr"`
	Id1SUzrok      string `json:"id_s_uzrok"`
	Opis           string `json:"opis"`
	IdZDsdfGL1     string `json:"id_z_dsdf_gl1"`
	IdZKvarGL1     string `json:"id_z_kvar_gl1"`
	IdZRapuGL1     string `json:"id_z_rapu_gl1"`
	IdZPrstGL1     string `json:"id_z_prst_gl1"`
	IdZZmspGL1     string `json:"id_z_zmsp_gl1"`
	IdZUzmsGL1     string `json:"id_z_uzms_gl1"`
	ZLokkGL1       string `json:"z_lokk_gl1"`
	IdZDsdfGL2     string `json:"id_z_dsdf_gl2"`
	IdZKvarGL2     string `json:"id_z_kvar_gl2"`
	IdZRapuGL2     string `json:"id_z_rapu_gl2"`
	IdZPrstGL2     string `json:"id_z_prst_gl2"`
	IdZZmspGL2     string `json:"id_z_zmsp_gl2"`
	IdZUzmsGL2     string `json:"id_z_uzms_gl2"`
	ZLokkGL2       string `json:"z_lokk_gl2"`
	IdZPrekVN      string `json:"id_z_prek_vn"`
	IdZDisREZ      string `json:"id_z_dis_rez"`
	IdZKvarREZ     string `json:"id_z_kvar_rez"`
	IdZPrstREZ     string `json:"id_z_prst_rez"`
	IdZZmspREZ     string `json:"id_z_zmsp_rez"`
	IdZNel1        string `json:"id_z_nel1"`
	IdZNel2        string `json:"id_z_nel2"`
	IdZNel3        string `json:"id_z_nel3"`
	IdZPrekNN      string `json:"id_z_prek_nn"`
	IdZSabzSAB     string `json:"id_z_sabz_sab"`
	IdZOtprSAB     string `json:"id_z_otpr_sab"`
	IdSVremUSL     string `json:"id_s_vrem_usl"`
	IdZJpsVN       string `json:"id_z_jps_vn"`
	IdZJpsNN       string `json:"id_z_jsp_vn"`
	IdZTelePocGL1  string `json:"id_s_z_tele_poc_gl1"`
	IdZTeleKrajGL1 string `json:"id_z_tele_kraj_gl1"`
	IdZTelePocGL2  string `json:"id_z_tele_poc_gl2"`
	IdZTeleKrajGL2 string `json:"id_z_tele_kraj_gl2"`
	SynsoftId      string `json:"ed_id"`
}

type DDNInterruptionOfDelivery struct {
	IdSMrc            string `json:"id_s_mrc"`
	IdSTipd           string `json:"id_s_tipd"`
	IdSVrpd           string `json:"id_s_vrpd"`
	IdTipob           string `json:"id_tipob"`
	ObId              string `json:"ob_id"`
	Vrepoc            string `json:"vrepoc"`
	Vrezav            string `json:"vrezav"`
	TipDogadjaja      string `json:"tip_dogadjaja"`
	Uzrok             string `json:"uzrok"`
	IdSVrPrek         string `json:"id_s_vr_prek"`
	SVrPrek     	  SVrPrek     `json:"vrsta_prek"`
	IdSUzrokPrek      string `json:"id_s_uzrok_prek"`
	SUzrokPrek     	  SUzrokPrek     `json:"uzrok_prek"`
	Snaga             string `json:"snaga"`
	Opis              string `json:"opis"`
	KorUneo           string `json:"kor_uneo"`
	IdSMernaMesta     string `json:"id_s_mrena_mesta"`
	SMernaMesta       SMernaMesta     `json:"merna_mesta"`
	BrojMesta         string `json:"broj_mesta"`
	Ind               string `json:"ind"`
	P2TrafId          string `json:"p2_traf_id"`
	Bi                string `json:"bi"`
	IdSPoduzrokPrek   string `json:"id_s_poduzrok_prek"`
	SPoduzrokPrek     SPoduzrokPrek     `json:"poduzrok_prek"`
	IdDogPrekidP      string `json:"id_dog_prekid_p"`
	IdTipObjektaNdc   string `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc string `json:"id_tip_dogadjaja_ndc"`
	SynsoftId         string `json:"ed_id"`
}

type DDNInterruptionOfDeliveryPayload struct {
	IdSMrc       string `json:"id_s_mrc"`
	IdSTipd      string `json:"id_s_tipd"`
	IdTipob      string `json:"id_tipob"`
	ObId         string `json:"ob_id"`
	Vrepoc       string `json:"vrepoc"`
	Vrezav       string `json:"vrezav"`
	TipDogadjaja string `json:"tip_dogadjaja"`
	Uzrok        string `json:"uzrok"`
	Snaga        string `json:"snaga"`
	Opis         string `json:"opis"`
	KorUneo      string `json:"kor_uneo"`
	P2TrafId     string `json:"p2_traf_id"`
	TipObjekta   string `json:"tip_objekta"`
	SynsoftId    string `json:"ed_id"`
}

type VSOb struct {
	IpsId           string `json:"ips_id"`
	Tipob           string `json:"tipob"`
	ObId         string `json:"ob_id"`
	ObSif string `json:"ob_sif"`
	Nazob       int    `json:"nazob"`
	Opis     string `json:"opis"`
	NnId     string `json:"nn_id"`
	NnSifra          int    `json:"nn_sifra"`
	Skt        string `json:"skr"`
	IdSMrc1          string `json:"id_s_mrc1"`
	Mrc1        string `json:"mrc1"`
	IdSMrc2           int    `json:"id_s_mrc2"`
	Mrc2         string `json:"mrc2"`
	IdSOrg1           int    `json:"id_s_org1"`
	IdSOrg2           int    `json:"id_s_org2"`
	Status         string `json:"status"`
}

type VSPoljeSvaAP struct {
	IdSOrg                int    `json:"id_s_org"`
	P2TrafId        string `json:"p2_traf_id"`
	ObId        string `json:"ob_id"`
	P1TrafId           int    `json:"p1_traf_id"`
	NnId           string `json:"nn_id"`
	Polje           string `json:"polje"`
	FupId            string `json:"fup_id"`
	Funkc               int    `json:"funkc"`
	TstId             string `json:"tst_id"`
	Sabir       string `json:"sabir"`
	Status         string `json:"status"`
	ImePo string `json:"ime_po"`
	IdKat         string `json:"id_kat"`
	Oprem       string `json:"oprem"`
	Aktne        string `json:"aktne"`
	Potpun      string `json:"potpun"`
	NormUkl           string `json:"norm_ukl"`
	Opis             string `json:"opis"`
	SapId         string `json:"sap_id"`
}


type STipPrek struct{
	ID     int    `json:"id"`
	Code string  `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SVrPrek struct{
	ID     int    `json:"id"`
	Code string  `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SUzrokPrek  struct{
	ID     int    `json:"id"`
	Code string  `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SPoduzrokPrek  struct{
	ID     int    `json:"id"`
	Code string  `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type SMernaMesta  struct{
	ID     int    `json:"id"`
	Code string  `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Description string `json:"description"`
}

type Permission struct {
	RbrPk           int64  `json:"rbr_pk"`
	BrZahtevaFK     int64  `json:"br_zahteva_fk"`
	SapSifra        string `json:"sap_sifra"`
	BrojIsk         string `json:"broj_isk"`
	BrDozvole       string `json:"br_dozvole"`
	TipRadova       string `json:"tip_radova"`
	EeObjekat       string `json:"ee_objekat"`
	Tip             string `json:"tip"`
	Oznaka          string `json:"oznaka"`
	StatusEl        string `json:"status_el"`
	Izuzev          string `json:"izuzev"`
	Napomena        string `json:"napomena"`
	DozIzdao        string `json:"doz_izdao"`
	DozPrimio       string `json:"doz_primio"`
	DatPrijemaDoz   string `json:"dat_prijema_doz"`
	VremePrijemaDoz string `json:"vreme_prijema_doz"`
	StatusDoz       string `json:"status_doz"`
	NapomenaZavRad  string `json:"napomena_zav_rad"`
	DozZavIzdao     string `json:"doz_zav_izdao"`
	DozZavPrimio    string `json:"doz_zav_primio"`
	DatZavRadova    string `json:"dat_zav_radova"`
	VremeZavRad     string `json:"vreme_zav_rad"`
}

type Request struct {
	RbrPk      int64  `json:"rbr_pk"`
	Rco        string `json:"rco"`
	BrojIsk    string `json:"broj_isk"`
	BrZahteva  string `json:"br_zahteva"`
	BrojKps    string `json:"broj_kps"`
	Grupa      string `json:"grupa"`
	Int        string `json:"int"`
	SapSifra   string `json:"sap_sifra"`
	EeObjekat  string `json:"ee_objekat"`
	Tip        string `json:"tip"`
	Oznaka     string `json:"oznaka"`
	Opis       string `json:"opis"`
	PlDatumOd  string `json:"pl_datum_od"`
	PlVremeOd  string `json:"pl_vreme_od"`
	PlDatumDo  string `json:"pl_datum_do"`
	PlVremeDo  string `json:"pl_vreme_do"`
	Nodob      string `json:"nodob"`
	TipRadova  string `json:"tip_radova"`
	Uslovi     string `json:"uslovi"`
	NapVeza    string `json:"nap_veza"`
	PodZahteva string `json:"pod_zhteva"`
	IskOdobrio string `json:"isk_odobrio"`
}

type Request1Gr struct {
	IdZahteva    int64  `json:"id_zahteva"`
	Grupa        string `json:"grupa"`
	Ukljucenost  string `json:"ukljucenost"`
	IntPl        string `json:"int_pl"`
	Br1Z1Gr      string `json:"br1_z_1gr"`
	Br2Z1Gr      string `json:"br2_z_1gr"`
	PlDatumOdZ   string `json:"pl_datum_od_z"`
	PlVremeOdZ   string `json:"pl_vreme_od_z"`
	PlDatumDoZ   string `json:"pl_datum_do_z"`
	PlVremeDoZ   string `json:"pl_vreme_do_z"`
	RukRadova    string `json:"ruk_radova"`
	Elementi     string `json:"elementi"`
	OpisRadova   string `json:"opis_radova"`
	NapomenaVeza string `json:"napomena_veza"`
}

type Request2Gr struct {
	IdZahteva       int64  `json:"id_zahteva"`
	Grupa           string `json:"grupa"`
	Ukljucenost     string `json:"ukljucenost"`
	IntPl           string `json:"int_pl"`
	BrZRDC2Gr       string `json:"br_z_rdc_2gr"`
	BrSag           string `json:"br_sag"`
	PlDatumOdZ      string `json:"pl_datum_od_z"`
	PlVremeOdZ      string `json:"pl_vreme_od_z"`
	PlDatumDoZ      string `json:"pl_datum_do_z"`
	PlVremeDoZ      string `json:"pl_vreme_do_z"`
	RukRadova       string `json:"ruk_radova"`
	Elementi        string `json:"elementi"`
	OpisRadova      string `json:"opis_radova"`
	NapomenaVeza    string `json:"napomena_veza"`
	SagUslovi       string `json:"sag_uslovi"`
	SagNapomenaVeza string `json:"sag_napomena_veza"`
}

type Outage struct {
	Datizv        string `json:"datizv"`
	Vrepoc        string `json:"vrepoc"`
	Vrezav        string `json:"vrezav"`
	Traj          string `json:"trajanje"`
	IpsId         string `json:"ips_id"`
	TipOb         string `json:"tipob"`
	Opis          string `json:"opis"`
	ImeDalekovoda string `json:"ime_dalekovoda"`
	PoljeTrafo    string `json:"polje_trafo"`
	Org1          string `json:"org1"`
	Org2          string `json:"org2"`
	Nazvrpd       string `json:"nazvrpd"`
	Uzrok         string `json:"uzrok"`
	VrmUsl        string `json:"vrm_usl"`
	Tekst         string `json:"tekst"`
	IdStavke      string `json:"id_stavke"`
	IdSeq         string `json:"id_seq"`
}

type Exclusion struct {
	Datizv        string `json:"datizv"`
	Vrepoc        string `json:"vrepoc"`
	Vrezav        string `json:"vrezav"`
	Traj          string `json:"trajanje"`
	IpsId         string `json:"ips_id"`
	TipOb         string `json:"tipob"`
	Opis          string `json:"opis"`
	ImeDalekovoda string `json:"ime_dalekovoda"`
	PoljeTrafo    string `json:"polje_trafo"`
	Org1          string `json:"org1"`
	Org2          string `json:"org2"`
	Nazvrpd       string `json:"nazvrpd"`
	Razlog        string `json:"razlog"`
	Tekst         string `json:"tekst"`
}

type Plan struct {
	IdPogOdr  string `json:"pog_odr_id"`
	IdSapFLok string `json:"sap_id"`
	IdIPSFLok string `json:"ips_id"`
	Opis      string `json:"opis"`
	IdPogPlan string `json:"id_pog_plan"`
	PlOdr     string `json:"pl_odr"`
	TksStOd   string `json:"tks_st_od"`
	TipNaloga string `json:"tip_naloga"`
	DatumPoc  string `json:"datum_poc"`
	DatumZav  string `json:"datum_zav"`
	Id        string `json:"id"`
}

// User is the type for users
type User struct {
	ID       int
	Username string
	Password string
	UserRole []UserRole
}

type Role struct {
	ID   int
	Code string
	Name string
}
type UserRole struct {
	ID       int
	IdUser   int
	IdRole   int
	RoleCode string
	RoleName string
}

/*func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}*/

// UnbalancedTraders
type UnbalancedTrader struct {
	Code      string `json:"code"`
	Deviation string `json:"deviation"`
}
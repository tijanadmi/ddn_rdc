package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DDNInterruptionOfDelivery struct {
	Id                string `json:"id"`
	IdSMrc            string `json:"id_s_mrc"`
	Mrc               string `json:"mrc"`
	IdSTipd           string `json:"id_s_tipd"`
	IdSVrpd           string `json:"id_s_vrpd"`
	IdTipob           string `json:"id_tipob"`
	ObId              string `json:"ob_id"`
	ObNaziv           string `json:"ob_naziv"`
	ObOpis            string `json:"ob_opis"`
	Vrepoc            string `json:"vrepoc"`
	Vrezav            string `json:"vrezav"`
	TipDogadjaja      string `json:"tip_dogadjaja"`
	IdSVrPrek         string `json:"id_s_vr_prek"`
	VrstaPrek         string `json:"vrsta_prek"`
	PodvrstaPrek      string `json:"podvrsta_prek"`
	IdSUzrokPrek      string `json:"id_s_uzrok_prek"`
	Uzrok             string `json:"uzrok"`
	Snaga             string `json:"snaga"`
	Opis              string `json:"opis"`
	KorUneo           string `json:"kor_uneo"`
	IdSMernaMesta     string `json:"id_s_merna_mesta"`
	MernaMesta        string `json:"merna_mesta"`
	BrojMesta         string `json:"broj_mesta"`
	Ind               string `json:"ind"`
	P2TrafId          string `json:"p2_traf_id"`
	PoljeNaziv        string `json:"polje_naziv"`
	PoljeOpis         string `json:"polje_opis"`
	Bi                string `json:"bi"`
	IdSPoduzrokPrek   string `json:"id_s_poduzrok_prek"`
	PoduzrokPrek      string `json:"poduzrok_prek"`
	IdDogPrekidP      string `json:"id_dog_prekid_p"`
	IdTipObjektaNdc   string `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc string `json:"id_tip_dogadjaja_ndc"`
	SynsoftId         string `json:"ed_id"`
}

type DDNInterruptionOfDeliveryRez struct {
	Id                int           `json:"id"`
	IdSMrc            int           `json:"id_s_mrc"`
	SMrc              SMrc          `json:"mrc"`
	IdSTipd           int           `json:"id_s_tipd"`
	IdSVrpd           int           `json:"id_s_vrpd"`
	IdTipob           int           `json:"id_tipob"`
	ObId              int           `json:"ob_id"`
	VSOb              VSOb          `json:"objekat"`
	Vrepoc            string        `json:"vrepoc"`
	Vrezav            string        `json:"vrezav"`
	IdSVrPrek         int           `json:"id_s_vr_prek"`
	SVrPrek           *SVrPrek      `json:"vrsta_prek"`
	IdSUzrokPrek      int           `json:"id_s_uzrok_prek"`
	SUzrokPrek        *SUzrokPrek   `json:"uzrok_prek"`
	Snaga             string        `json:"snaga"`
	Opis              string        `json:"opis"`
	KorUneo           string        `json:"kor_uneo"`
	IdDogSmene        int           `json:"id_dog_smene"`
	IdStavke          int           `json:"id_stavke"`
	Mod               string        `json:"mod"`
	IdSMernaMesta     sql.NullInt64 `json:"id_s_mrena_mesta"`
	SMernaMesta       *SMernaMesta  `json:"merna_mesta"`
	BrojMesta         int           `json:"broj_mesta"`
	Ind               string        `json:"ind"`
	P2TrafId          int           `json:"p2_traf_id"`
	VSPoljeSvaAP      *VSPoljeSvaAP `json:"polja"`
	Bi                int           `json:"bi"`
	IdSPoduzrokPrek   int           `json:"id_s_poduzrok_prek"`
	SPoduzrokPrek     SPoduzrokPrek `json:"poduzrok_prek"`
	IdDogPrekidP      int           `json:"id_dog_prekid_p"`
	IdTipObjektaNdc   int           `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc int           `json:"id_tip_dogadjaja_ndc"`
	SynsoftId         int           `json:"synsoft_id"`
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
	IpsId   string `json:"ips_id"`
	Tipob   string `json:"tipob"`
	ObId    string `json:"ob_id"`
	ObSif   string `json:"ob_sif"`
	Nazob   int    `json:"nazob"`
	Opis    string `json:"opis"`
	NnId    string `json:"nn_id"`
	NnSifra int    `json:"nn_sifra"`
	Skt     string `json:"skr"`
	IdSMrc1 string `json:"id_s_mrc1"`
	Mrc1    string `json:"mrc1"`
	IdSMrc2 int    `json:"id_s_mrc2"`
	Mrc2    string `json:"mrc2"`
	IdSOrg1 int    `json:"id_s_org1"`
	IdSOrg2 int    `json:"id_s_org2"`
	Status  string `json:"status"`
}

type VSPoljeSvaAP struct {
	IdSOrg   int    `json:"id_s_org"`
	P2TrafId string `json:"p2_traf_id"`
	ObId     string `json:"ob_id"`
	P1TrafId int    `json:"p1_traf_id"`
	NnId     string `json:"nn_id"`
	Polje    string `json:"polje"`
	FupId    string `json:"fup_id"`
	Funkc    int    `json:"funkc"`
	TstId    string `json:"tst_id"`
	Sabir    string `json:"sabir"`
	Status   string `json:"status"`
	ImePo    string `json:"ime_po"`
	IdKat    string `json:"id_kat"`
	Oprem    string `json:"oprem"`
	Aktne    string `json:"aktne"`
	Potpun   string `json:"potpun"`
	NormUkl  string `json:"norm_ukl"`
	Opis     string `json:"opis"`
	SapId    string `json:"sap_id"`
}

type ListLimitOffsetParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListInterruptionParams struct {
	Mrc       string  `json:"mrc"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Ind       string `json:"ind"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type STipPrek struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SVrPrek struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SUzrokPrek struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SPoduzrokPrek struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type SMernaMesta struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type SMrc struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	NameCir string `json:"name_cir"`
}

type ListObjectLimitOffsetParams struct {
	Mrc    int32 `json:"mrc"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
type ObjLOV struct {
	Opis     string `json:"opis"`
	IdSMrc1  string `json:"id_s_mrc1"`
	ObId     string `json:"ob_id"`
	Tipob    string `json:"tipob"`
	SifTipob string `json:"sif_tipob"`
	ObSif    string `json:"ob_sif"`
}

type ListPoljaLimitOffsetParams struct {
	ObjId  int32 `json:"obj_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
type PoljaLOV struct {
	Id         string `json:"polje_id"`
	Polje      string `json:"polje"`
	PoljeNaziv string `json:"polje_naziv"`
	NNId       string `json:"nn_id"`
	NNNaziv    string `json:"nn_naziv"`
}

// User is the type for users
type User struct {
	ID       int
	Username string
	Password string
	FullName string
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

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type FailuresAll struct {
	Datizv        string `json:"datizv"`
	IdSMrc        string `json:"id_s_mrc"`
	Id1           string `json:"id1"`
	Id2           string `json:"id2"`
	Vrepoc        string `json:"vrepoc"`
	Vrezav        string `json:"vrezav"`
	Traj          string `json:"trajnje"`
	IpsId         string `json:"-"`
	ObId          string `json:"-"`
	TipOb         string `json:"-"`
	Objekat       string `json:"objekat"`
	ImeDalekovoda string `json:"ime_dalekovoda"`
	P2TrafId      string `json:"p2_traf_id"`
	IdSNap        string `json:"id_s_nap"`
	TrafoId       string `json:"trafo_id"`
	PoljeTrafo    string `json:"facility"`
	Org1          string `json:"object_belongs_to1"`
	Org2          string `json:"object_belongs_to2"`
	Vrpd          string `json:"vrpd"`
	Nazvrpd       string `json:"event_type"`
	Uzrok         string `json:"event_cause_type"`
	VrmUsl        string `json:"weather_conditions"`
	Tekst         string `json:"description"`
	ZDsdfGl1      string `json:"z_dsdf_gl1"`
	ZKvarGl1      string `json:"z_kvar_gl1"`
	ZPrstGl1      string `json:"z_prst_gl1"`
	ZZmspGl1      string `json:"z_zmsp_gl1"`
	ZUzmsGl1      string `json:"z_uzms_gl1"`
	ZRapuGl1      string `json:"z_rapu_gl1"`
	ZLokkGl1      string `json:"z_lokk_gl1"`
	ZDsdfGl2      string `json:"z_dsdf_gl2"`
	ZKvarGl2      string `json:"z_kvar_gl2"`
	ZPrstGl2      string `json:"z_prst_gl2"`
	ZZmspGl2      string `json:"z_zmsp_gl2"`
	ZUzmsGl2      string `json:"z_uzms_gl2"`
	ZRapuGl2      string `json:"z_rapu_gl2"`
	ZLokkGl2      string `json:"z_lokk_gl2"`
	ZDisRez       string `json:"z_dis_rez"`
	ZKvarRez      string `json:"z_kvar_rez"`
	ZPrstRez      string `json:"z_prst_rez"`
	ZZmspRez      string `json:"z_zmsp_rez"`
	ZPrekVn       string `json:"z_prek_vn"`
	ZPrekNn       string `json:"z_prek_nn"`
	ZNel1         string `json:"z_nel1"`
	ZNel2         string `json:"z_nel2"`
	ZNel3         string `json:"z_nel3"`
	ZSabzSab      string `json:"z_sabz_sab"`
	ZOtprSab      string `json:"z_otpr_sab"`
	ZJpsVn        string `json:"z_jps_vn"`
	ZJpsNn        string `json:"z_jps_nn"`
	ZTelePocGl1   string `json:"z_poc_kraj_gl1"`
	ZTeleKrajGl1  string `json:"z_tele_kraj_gl1"`
	ZTelePocGl2   string `json:"z_tele_poc_gl2"`
	ZTeleKrajGl2  string `json:"z_tele_kraj_gl2"`
	Snaga         string `json:"snaga"`
}

type PiMM struct {
	Id         int    `json:"id"`
	Stav       int    `json:"stav"`
	Tipd       string `json:"tipd"`
	Datizv     string `json:"datizv"`
	Id1        string `json:"id1"`
	Id2        string `json:"id2"`
	Vrepoc     string `json:"vrepoc"`
	Vrezav     string `json:"vrezav"`
	Traj       string `json:"trajnje"`
	ObId       string `json:"ob_id"`
	TipOb      string `json:"ob_tip"`
	ObSif      string `json:"ob_sif"`
	NazOb      string `json:"naz_ob"`
	PoljeTrafo string `json:"polje_trafo"`
	IdSNap     string `json:"id_s_nap"`
	TrafoId    string `json:"trafo_id"`
	P2TrafId   string `json:"p2_traf_id"`
	Napon      string `json:"napon"`
	Polje      string `json:"polje"`
	ImePo      string `json:"ime_po"`
	Funkc      string `json:"funkc"`
	Vrpd       string `json:"vrpd"`
	Nazvrpd    string `json:"naziv_vrpd"`
	Uzrok      string `json:"uzrok"`
	Razlog     string `json:"razlog"`
	VrmUsl     string `json:"weather_conditions"`
	Tekst      string `json:"description"`
	Snaga      string `json:"snaga"`

	/*NazSop      string `json:"naz_sop"`
	SopNaziv      string `json:"sop_naziv"`
	IdSSop      string `json:"id_s_sop"`
	IdSop      string `json:"id_sop"`*/

	IdZDsdfGl1 string `json:"id_z_dsdf_gl1"`
	ZDsdfGl1   string `json:"z_dsdf_gl1"`
	IdZKvarGl1 string `json:"id_z_kvar_gl1"`
	ZKvarGl1   string `json:"z_kvar_gl1"`
	IdZRapuGl1 string `json:"id_z_rapu_gl1"`
	ZRapuGl1   string `json:"z_rapu_gl1"`
	IdZPrstGl1 string `json:"id_z_prst_gl1"`
	ZPrstGl1   string `json:"z_prst_gl1"`
	IdZZmspGl1 string `json:"id_z_zmsp_gl1"`
	ZZmspGl1   string `json:"z_zmsp_gl1"`
	IdZUzmsGl1 string `json:"id_z_uzms_gl1"`
	ZUzmsGl1   string `json:"z_uzms_gl1"`
	ZLokkGl1   string `json:"z_lokk_gl1"`

	IdZDsdfGl2     string `json:"id_z_dsdf_gl2"`
	ZDsdfGl2       string `json:"z_dsdf_gl2"`
	IdZKvarGl2     string `json:"id_z_kvar_gl2"`
	ZKvarGl2       string `json:"z_kvar_gl2"`
	IdZRapuGl2     string `json:"id_z_rapu_gl2"`
	ZRapuGl2       string `json:"z_rapu_gl2"`
	IdZPrstGl2     string `json:"id_z_prst_gl2"`
	ZPrstGl2       string `json:"z_prst_gl2"`
	IdZZmspGl2     string `json:"id_z_zmsp_gl2"`
	ZZmspGl2       string `json:"z_zmsp_gl2"`
	IdZUzmsGl2     string `json:"id_z_uzms_gl2"`
	ZUzmsGl2       string `json:"z_uzms_gl2"`
	ZLokkGl2       string `json:"zLokkGl2"`
	IdZDisRez      string `json:"idZDisRez"`
	ZDisRez        string `json:"zDisRez"`
	IdZKvarRez     string `json:"idZKvarRez"`
	ZKvarRez       string `json:"zKvarRez"`
	IdZPrstRez     string `json:"idZPrstRez"`
	ZPrstRez       string `json:"zPrstRez"`
	IdZZmspRez     string `json:"idZZmspRez"`
	ZZmspRez       string `json:"zZmspRez"`
	IdZPrekVn      string `json:"idZPrekVn"`
	ZPrekVn        string `json:"zPrekVn"`
	IdZPrekNn      string `json:"idZPrekNn"`
	ZPrekNn        string `json:"zPrekNn"`
	IdZNel1        string `json:"idZNel1"`
	ZNel1          string `json:"zNel1"`
	IdZNel2        string `json:"idZNel2"`
	ZNel2          string `json:"zNel2"`
	IdZNel3        string `json:"idZNel3"`
	ZNel3          string `json:"zNel3"`
	IdZSabzSab     string `json:"idZSabzSab"`
	ZSabzSab       string `json:"zSabzSab"`
	IdZOtprSab     string `json:"idZOtprSab"`
	ZOtprSab       string `json:"zOtprSab"`
	IdZJpsVn       string `json:"idZJpsVn"`
	ZJpsVn         string `json:"zJpsVn"`
	IdZJpsNn       string `json:"idZJpsNn"`
	ZJpsNn         string `json:"zJpsNn"`
	IdZTelePocGl1  string `json:"idZTelePocGl1"`
	ZTelePocGl1    string `json:"zTelePocGl1"`
	IdZTeleKrajGl1 string `json:"idZTeleKrajGl1"`
	ZTeleKrajGl1   string `json:"zTeleKrajGl1"`
	IdZTelePocGl2  string `json:"idZTelePocGl2"`
	ZTelePocGl2    string `json:"zTelePocGl2"`
	IdZTeleKrajGl2 string `json:"idZTeleKrajGl2"`
	ZTeleKrajGl2   string `json:"zTeleKrajGl2"`
	Fup            string `json:"fup"`
}

type PiMMT4 struct {
	Id     int    `json:"id"`
	Stav   string `json:"stav"`
	Datizv string `json:"datizv"`
	Id1    string `json:"id1"`
	Mrc    string `json:"mrc"`
	Tekst  string `json:"tekst"`
	Kom1   string `json:"kom1"`
	Kom2   string `json:"kom2"`
	Kom3   string `json:"kom3"`
	Kom4   string `json:"kom4"`
	Kom5   string `json:"kom5"`
	Kom6   string `json:"kom6"`
	Kom7   string `json:"kom7"`
	Kom8   string `json:"kom8"`
	Opist4 string `json:"opist4"`
}

type ListPiMMParams struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Tipd      int    `json:"tipd"`
}

type ListPiMMParamsByPage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Tipd      int    `json:"tipd"`
	Fup       string `json:"fup"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}
type ListPiMMT4Params struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type ListPiDDT4Params struct {
	Datizv string `json:"datizv"`
	IdSMrc string `json:"id_s_mrc"`
}

type PiMMResponse struct {
	PiMm  []PiMM `json:"pi_mm"`
	Count int    `json:"count"`
}

type PiDD struct {
	Id         int    `json:"id"`
	IdSMrc     string `json:"id_s_mrc"`
	Mrc        string `json:"mrc"`
	Tipd       string `json:"tipd"`
	Datizv     string `json:"datizv"`
	Id1        string `json:"id1"`
	Id2        string `json:"id2"`
	Vrepoc     string `json:"vrepoc"`
	VrepocHHMI string `json:"vrepoc_hhmi"`
	PocPP      string `json:"poc_pp"`
	Vrezav     string `json:"vrezav"`
	VrezavHHMI string `json:"vrezav_hhmi"`
	ZavPP      string `json:"zav_pp"`
	Traj       string `json:"trajnje"`
	ObId       string `json:"ob_id"`
	TipOb      string `json:"ob_tip"`
	ObSif      string `json:"ob_sif"`
	NazOb      string `json:"naz_ob"`
	PoljeTrafo string `json:"polje_trafo"`
	IdSNap     string `json:"id_s_nap"`
	TrafoId    string `json:"trafo_id"`
	P2TrafId   string `json:"p2_traf_id"`
	Napon      string `json:"napon"`
	Polje      string `json:"polje"`
	ImePo      string `json:"ime_po"`
	Funkc      string `json:"funkc"`
	Vrpd       string `json:"vrpd"`
	Nazvrpd    string `json:"naziv_vrpd"`
	GrUzrok    string `json:"gr_uzrok"`
	Uzrok      string `json:"uzrok"`
	GrRazlog   string `json:"gr_razlog"`
	Razlog     string `json:"razlog"`
	VrmUsl     string `json:"weather_conditions"`
	Opis       string `json:"description"`
	Snaga      string `json:"snaga"`

	NazSop   string `json:"naz_sop"`
	SopNaziv string `json:"sop_naziv"`
	IdSSop   string `json:"id_s_sop"`
	IdSop    string `json:"id_sop"`

	IdZDsdfGl1 string `json:"id_z_dsdf_gl1"`
	ZDsdfGl1   string `json:"z_dsdf_gl1"`
	IdZKvarGl1 string `json:"id_z_kvar_gl1"`
	ZKvarGl1   string `json:"z_kvar_gl1"`
	IdZRapuGl1 string `json:"id_z_rapu_gl1"`
	ZRapuGl1   string `json:"z_rapu_gl1"`
	IdZPrstGl1 string `json:"id_z_prst_gl1"`
	ZPrstGl1   string `json:"z_prst_gl1"`
	IdZZmspGl1 string `json:"id_z_zmsp_gl1"`
	ZZmspGl1   string `json:"z_zmsp_gl1"`
	IdZUzmsGl1 string `json:"id_z_uzms_gl1"`
	ZUzmsGl1   string `json:"z_uzms_gl1"`
	ZLokkGl1   string `json:"z_lokk_gl1"`

	IdZDsdfGl2     string `json:"id_z_dsdf_gl2"`
	ZDsdfGl2       string `json:"z_dsdf_gl2"`
	IdZKvarGl2     string `json:"id_z_kvar_gl2"`
	ZKvarGl2       string `json:"z_kvar_gl2"`
	IdZRapuGl2     string `json:"id_z_rapu_gl2"`
	ZRapuGl2       string `json:"z_rapu_gl2"`
	IdZPrstGl2     string `json:"id_z_prst_gl2"`
	ZPrstGl2       string `json:"z_prst_gl2"`
	IdZZmspGl2     string `json:"id_z_zmsp_gl2"`
	ZZmspGl2       string `json:"z_zmsp_gl2"`
	IdZUzmsGl2     string `json:"id_z_uzms_gl2"`
	ZUzmsGl2       string `json:"z_uzms_gl2"`
	ZLokkGl2       string `json:"zLokkGl2"`
	IdZDisRez      string `json:"idZDisRez"`
	ZDisRez        string `json:"zDisRez"`
	IdZKvarRez     string `json:"idZKvarRez"`
	ZKvarRez       string `json:"zKvarRez"`
	IdZPrstRez     string `json:"idZPrstRez"`
	ZPrstRez       string `json:"zPrstRez"`
	IdZZmspRez     string `json:"idZZmspRez"`
	ZZmspRez       string `json:"zZmspRez"`
	IdZPrekVn      string `json:"idZPrekVn"`
	ZPrekVn        string `json:"zPrekVn"`
	IdZPrekNn      string `json:"idZPrekNn"`
	ZPrekNn        string `json:"zPrekNn"`
	IdZNel1        string `json:"idZNel1"`
	ZNel1          string `json:"zNel1"`
	IdZNel2        string `json:"idZNel2"`
	ZNel2          string `json:"zNel2"`
	IdZNel3        string `json:"idZNel3"`
	ZNel3          string `json:"zNel3"`
	IdZSabzSab     string `json:"idZSabzSab"`
	ZSabzSab       string `json:"zSabzSab"`
	IdZOtprSab     string `json:"idZOtprSab"`
	ZOtprSab       string `json:"zOtprSab"`
	IdZJpsVn       string `json:"idZJpsVn"`
	ZJpsVn         string `json:"zJpsVn"`
	IdZJpsNn       string `json:"idZJpsNn"`
	ZJpsNn         string `json:"zJpsNn"`
	IdZTelePocGl1  string `json:"idZTelePocGl1"`
	ZTelePocGl1    string `json:"zTelePocGl1"`
	IdZTeleKrajGl1 string `json:"idZTeleKrajGl1"`
	ZTeleKrajGl1   string `json:"zTeleKrajGl1"`
	IdZTelePocGl2  string `json:"idZTelePocGl2"`
	ZTelePocGl2    string `json:"zTelePocGl2"`
	IdZTeleKrajGl2 string `json:"idZTeleKrajGl2"`
	ZTeleKrajGl2   string `json:"zTeleKrajGl2"`
	Fup            string `json:"fup"`
}

type ListPiDDParams struct {
	Datizv string `json:"datizv"`
	Tipd   int    `json:"tipd"`
	IdSMrc string `json:"id_s_mrc"`
}

type ListPiDDParamsByPage struct {
	Datizv string `json:"datizv"`
	Tipd   int    `json:"tipd"`
	Fup    string `json:"fup"`
	IdSMrc string `json:"id_s_mrc"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type PiDDResponse struct {
	PiDD  []PiDD `json:"pi_dd"`
	Count int    `json:"count"`
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

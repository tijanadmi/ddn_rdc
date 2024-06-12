package models

import (
	"time"

	"github.com/google/uuid"
)

type DDNInterruptionOfDelivery struct {
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
	SVrPrek           SVrPrek       `json:"vrsta_prek"`
	IdSUzrokPrek      int           `json:"id_s_uzrok_prek"`
	SUzrokPrek        SUzrokPrek    `json:"uzrok_prek"`
	Snaga             string        `json:"snaga"`
	Opis              string        `json:"opis"`
	KorUneo           string        `json:"kor_uneo"`
	IdDogSmene        string        `json:"id_dog_smene"`
	IdStavke          string        `json:"id_stavke"`
	Mod               string        `json:"mod"`
	IdSMernaMesta     int           `json:"id_s_mrena_mesta"`
	SMernaMesta       SMernaMesta   `json:"merna_mesta"`
	BrojMesta         int           `json:"broj_mesta"`
	Ind               string        `json:"ind"`
	P2TrafId          int           `json:"p2_traf_id"`
	Bi                int           `json:"bi"`
	IdSPoduzrokPrek   int           `json:"id_s_poduzrok_prek"`
	SPoduzrokPrek     SPoduzrokPrek `json:"poduzrok_prek"`
	IdDogPrekidP      int           `json:"id_dog_prekid_p"`
	IdTipObjektaNdc   int           `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc int           `json:"id_tip_dogadjaja_ndc"`
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

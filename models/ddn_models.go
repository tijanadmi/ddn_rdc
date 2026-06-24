package models

import (
	"time"
)

type ListShiftsWithPaginationParams struct {
	Mrc       string `json:"mrc"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type Smena struct {
	IdSmene    int       `json:"id_smene"`
	DatDnev    time.Time `json:"dat_dnev"`
	DatDnevStr string    `json:"datdnev_str"`

	IdSMRC int    `json:"id_s_mrc"`
	Rdc    string `json:"rdc"`

	DezDisp1    string `json:"dez_disp1"`
	DezDisp1Ime string `json:"dez_disp1_ime"`

	DezDisp2    string `json:"dez_disp2"`
	DezDisp2Ime string `json:"dez_disp2_ime"`

	DezDisp3    string `json:"dez_disp3"`
	DezDisp3Ime string `json:"dez_disp3_ime"`

	DezDisp4    string `json:"dez_disp4"`
	DezDisp4Ime string `json:"dez_disp4_ime"`

	IDTipSmena int    `json:"id_tip_smena"`
	TipSkr     string `json:"tip_skr"`
	TipSmena   string `json:"tip_smena"`
	IntSmena   string `json:"int_smena"`
	MrcSif     string `json:"mrc_sif"`

	PredaoDisp1 string `json:"predao_disp1"`
	PredaoDisp2 string `json:"predao_disp2"`
	PrimDisp1   string `json:"prim_disp1"`
	PrimDisp2   string `json:"prim_disp2"`
	PredaoDisp3 string `json:"predao_disp3"`
	PrimDisp3   string `json:"prim_disp3"`

	KomentZat string `json:"koment_zat"`
	OtvSpec   string `json:"otv_spec"`
	ZatSpec   string `json:"zat_spec"`

	IDKatDok int    `json:"id_kat_dok"`
	IdBroj   string `json:"id_broj"`
	Dan      string `json:"dan"`

	// kasnije dodaješ ovo:
	Dogadjaji []Dogadjaj `json:"dogadjaji,omitempty"`
}

type Dogadjaj struct {
	ID        int    `json:"id"`
	IDSmena   int    `json:"id_smena"`
	RbDog     string `json:"rb_dog"`
	Naslov    string `json:"naslov"`
	IDTipDog  int    `json:"id_tip_dog"`
	TipDog    string `json:"tip_dog"`
	TipDogCir string `json:"tip_dog_cir"`
	Tip       string `json:"tip"`
	TipObav   string `json:"tip_obav"`
	Dopuna    string `json:"dopuna"`
	Status    string `json:"status"`
}

type DogadjajDetaljno struct {
	ID     int    `json:"id"`
	RbDog  string `json:"rb_dog"`
	TipDog string `json:"tip_dog"`

	Naslov    string `json:"naslov"`
	Podnaslov string `json:"podnaslov"`

	IDSmena        int        `json:"id_smena"`
	VezaSa         *int       `json:"veza_sa"`
	RbDogVezaSa    *string    `json:"rb_dog_veza_sa"`
	DatumVezaSa    *time.Time `json:"datum_veza_sa"`
	TipSmenaVezaSa *string    `json:"tip_smene_veza_sa"`

	Dopuna   *string `json:"dopuna"`     // može biti null
	IDSmenaD *int    `json:"id_smena_d"` // može biti null

	IDSrazlog *int `json:"id_s_razlog"`

	Grazlog string `json:"grazlog"`
	Razlog  string `json:"razlog"`

	UzrokTekst *string `json:"uzrok_tekst"`
	ManTekst   *string `json:"man_tekst"`
	Posledice  *string `json:"posledice"`

	DatumSmene  time.Time  `json:"datum_smene"`  // već formatiran
	DatumDopune *time.Time `json:"datum_dopune"` // može biti null

	TipSmene       string  `json:"tip_smene"`
	TipSmeneDopune *string `json:"tip_smene_dopune"`

	Manipulacije        []Manipulacija       `json:"manipulacije"`
	ObavBeleske         *ObavBeleska         `json:"obav_beleske,omitempty"`
	TSU                 *[]TSU               `json:"tsu"`
	TK                  *[]TK                `json:"tk"`
	SOP                 *[]SOP               `json:"sop"`
	Ispad               *[]Ispad             `json:"ispad"`
	PrekidP             *[]PrekidP           `json:"prekid_proizvodnje"`
	AngazovaniRukovalac *AngazovaniRukovalac `json:"angazovani_rukovalac"`
}

/************** Iskljucenja/Ukljucenja ****************/

type Manipulacija struct {
	IDDogSmene int  `json:"id_dog_smene"`
	Rb         *int `json:"rb"`

	Vrepoc string  `json:"vrepoc"` // već formatiran hh:mm
	Vrezav *string `json:"vrezav"` // može biti null

	Manipulacija string  `json:"manipulacija"`
	TekstMan     *string `json:"tekst_man"`

	Ees      *string `json:"ees"`
	TekstEes *string `json:"tekst_ees"`
	Napomena *string `json:"napomena"`

	Objekat      string  `json:"objekat"`
	DvTrafoPolje *string `json:"dv_trafo_polje"`

	StatusMan string  `json:"status_man"`
	DopunaMan *string `json:"dopuna_man"`

	IDStavkeM int `json:"id_stavke_m"`
}

type ObjekatView struct {
	Naziv  string      `json:"naziv"`
	MinRb  int         `json:"-"`
	Stavke []ManipView `json:"stavke"`
}

type ManipView struct {
	DopunaDaNe  string `json:"dopuna_da_ne"`
	Vrepoc      string `json:"vrepoc"`
	Vrezav      string `json:"vrezav"`
	RecenicaMan string `json:"recenica_man"`
	Rb          int    `json:"-"`
}

/************** Kraj Iskljucenja/Ukljucenja ****************/
type ObavBeleska struct {
	Napomena string `json:"napomena"`  // napomena, može biti prazan string
	TekstObv string `json:"tekst_obv"` // ovde ide CLOB sadržaj iz baze
	TipObv   string `json:"tip_obv"`   // tip obaveštenja
	Dopuna   string `json:"dopuna"`    // dopuna, može biti prazan string
}

type TK struct {
	ID            int       `json:"id"`
	Vrepoc        string    `json:"vrepoc"`
	Vrezav        *string   `json:"vrezav"`
	VrstaDog      string    `json:"vrstaDog"`
	Opis          *string   `json:"opis"`
	ObID          int       `json:"obId"`
	ObjekatNaziv  string    `json:"objekatNaziv"`
	ObID2         *int      `json:"obId2"`
	ObjekatNaziv2 *string   `json:"objekatNaziv2"`
	Vropr         string    `json:"vropr"`
	Vrpd          int       `json:"vrpd"`
	Status        *string   `json:"status"`
	Dopuna        *string   `json:"dopuna"`
	VrepocSort    time.Time `json:"vrepocSort"`
}

type DetaljT567 struct {
	DopunaDaNe string `json:"dopuna_da_ne"`
	Recenica1  string `json:"recenica1"`
	Recenica2  string `json:"recenica2"`
	Opis       string `json:"opis"`
}

type TSU struct {
	ID           int       `json:"id"`
	Vrepoc       string    `json:"vrepoc"`
	Vrezav       *string   `json:"vrezav"`
	VrstaDog     string    `json:"vrstaDog"`
	Opis         *string   `json:"opis"`
	ObID         int       `json:"obId"`
	ObjekatNaziv string    `json:"objekatNaziv"`
	Vropr        string    `json:"vropr"`
	Vrpd         int       `json:"vrpd"`
	Status       *string   `json:"status"`
	Dopuna       *string   `json:"dopuna"`
	VrepocSort   time.Time `json:"vrepocSort"`
}

type SOP struct {
	ID           int       `json:"id"`
	Vrepoc       string    `json:"vrepoc"`
	Vrezav       *string   `json:"vrezav"`
	VrstaDog     string    `json:"vrstaDog"`
	Opis         *string   `json:"opis"`
	ObID         int       `json:"obId"`
	ObjekatNaziv string    `json:"objekatNaziv"`
	NazSop       string    `json:"nazSop"`
	RBrSop       string    `json:"rbrSop"`
	IdSSop       string    `json:"idSSop"`
	IdSop        string    `json:"idSop"`
	Vrpd         int       `json:"vrpd"`
	Status       *string   `json:"status"`
	Dopuna       *string   `json:"dopuna"`
	VrepocSort   time.Time `json:"vrepocSort"`
}

type Ispad struct {
	VrepocSortIspkv time.Time `json:"vrepocSortIspkv"`
	ID              int       `json:"id"`
	RB              *int      `json:"rb"`

	Vrepoc string `json:"vrepoc"`
	Vrezav string `json:"vrezav"`

	VrstaDog string `json:"vrstaDog"`
	VrDogSif string `son:"vrDogSif"`
	TipDog   string `json:"tipDog"`
	Gruzr1   string `json:"gruzr1"`
	Uzrok1   string `json:"uzrok1"`

	Objekat      string `json:"objekat"`
	TipOb        string `json:"tipOb"`
	Napon        string `json:"napon"`
	DvTrafoPolje string `json:"dvTrafoPolje"`
	Fup          string `json:"fup"`

	Opis    string `json:"opis"`
	SmPk    string `json:"smPk"`
	RadApu  string `json:"radApu"`
	VremUsl string `json:"vremUsl"`

	StatusIspkv1 string  `json:"statusIspkv1"`
	DopunaIspkv1 *string `json:"dopunaIspkv1"`

	IDStavkeI int `json:"idStavkeI"`

	// GL1
	ZDsdfGl1 *string `json:"zDsdfGl1"`
	ZKvarGl1 *string `json:"zKvarGl1"`
	ZPrstGl1 *string `json:"zPrstGl1"`
	ZZMSPGl1 *string `json:"zZmspGl1"`
	ZUzmsGl1 *string `json:"zUzmsGl1"`
	ZRapuGl1 *string `json:"zRapuGl1"`
	ZLokkGl1 *string `json:"zLokkGl1"`

	// GL2
	ZDsdfGl2 *string `json:"zDsdfGl2"`
	ZKvarGl2 *string `json:"zKvarGl2"`
	ZPrstGl2 *string `json:"zPrstGl2"`
	ZZMSPGl2 *string `json:"zZmspGl2"`
	ZUzmsGl2 *string `json:"zUzmsGl2"`
	ZRapuGl2 *string `json:"zRapuGl2"`
	ZLokkGl2 *string `json:"zLokkGl2"`

	// GL3
	ZDsdfGl3 *string `json:"zDsdfGl3"`
	ZKvarGl3 *string `json:"zKvarGl3"`
	ZPrstGl3 *string `json:"zPrstGl3"`
	ZZMSPGl3 *string `json:"zZmspGl3"`
	ZUzmsGl3 *string `json:"zUzmsGl3"`
	ZRapuGl3 *string `json:"zRapuGl3"`
	ZLokkGl3 *string `json:"zLokkGl3"`

	// REZ
	ZDisRez  *string `json:"zDisRez"`
	ZKvarRez *string `json:"zKvarRez"`
	ZPrstRez *string `json:"zPrstRez"`
	ZZMSPRez *string `json:"zZmspRez"`

	ZDisRez2  *string `json:"zDisRez2"`
	ZKvarRez2 *string `json:"zKvarRez2"`
	ZPrstRez2 *string `json:"zPrstRez2"`
	ZZMSPRez2 *string `json:"zZmspRez2"`

	// ostalo
	ZPrekVn *string `json:"zPrekVn"`
	ZPrekNn *string `json:"zPrekNn"`
	ZNel1   *string `json:"zNel1"`
	ZNel2   *string `json:"zNel2"`
	ZNel3   *string `json:"zNel3"`

	ZSabzSab  *string `json:"zSabzSab"`
	ZSabzSab2 *string `json:"zSabzSab2"`
	ZOtprSab  *string `json:"zOtprSab"`
	ZOtprSab2 *string `json:"zOtprSab2"`

	ZJpsVn  *string `json:"zJpsVn"`
	ZJpsNn  *string `json:"zJpsNn"`
	ZJpsVn2 *string `json:"zJpsVn2"`
	ZJpsNn2 *string `json:"zJpsNn2"`

	// tele
	IdZTelePocGl1  *string `json:"idZTelePocGl1"`
	IdZTeleKrajGl1 *string `json:"idZTeleKrajGl1"`
	IdZTelePocGl2  *string `json:"idZTelePocGl2"`
	IdZTeleKrajGl2 *string `json:"idZTeleKrajGl2"`
	IdZTelePocGl3  *string `json:"idZTelePocGl3"`
	IdZTeleKrajGl3 *string `json:"idZTeleKrajGl3"`

	Snaga *string `json:"snaga"`
}

type PrekidP struct {
	VrepocSort string  `db:"VREPOC_SORT_ISPKV5678" json:"vrepoc_sort"`
	Polje      *string `db:"POLJE" json:"polje"`
	FupID      *int    `db:"FUP_ID" json:"fup_id"`
	Funkc      *string `db:"FUNKC" json:"funkc"`
	ObID       *int    `db:"OB_ID" json:"ob_id"`
	TipOb      *string `db:"TIP_OB" json:"tip_ob"`
	IDP2Traf   *int    `db:"ID_P2_TRAF" json:"id_p2_traf"`

	Vrepoc string  `db:"VREPOC" json:"vrepoc"`
	Vrezav *string `db:"VREZAV" json:"vrezav"`

	Objekat *string `db:"OBJEKAT" json:"objekat"`

	Generator *string `db:"GENERATOR" json:"generator"`
	VrPrek    *string `db:"VR_PREK" json:"vr_prek"`
	UzrokPrek *string `db:"UZROK_PREK" json:"uzrok_prek"`
	TipPrek   *string `db:"TIP_PREK" json:"tip_prek"`

	Snaga  *string `db:"SNAGA" json:"snaga"`
	Opis   *string `db:"OPIS" json:"opis"`
	Status *string `json:"status"`
	Dopuna *string `db:"DOPUNA_ISPKV5678" json:"dopuna"`
}

type ObavSlika struct {
	Base64 string `json:"base64"`
	Format string `json:"format"`
	RB     int    `json:"rb"`
}

type AngazovaniRukovalac struct {
	VremeNaloga  *time.Time `json:"vreme_naloga"`
	ImeNaloga    *string    `json:"ime_naloga"`
	VremeDolaska *time.Time `json:"vreme_dolaska"`
	VremeOdlaska *time.Time `json:"vreme_odlaska"`

	Rukovalac *string `json:"rukovalac"`
	Objekat   *string `json:"objekat"`

	Opis *string `json:"opis"`
}

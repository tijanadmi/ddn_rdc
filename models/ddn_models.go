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
	IdSmene int       `json:"id_smene"`
	DatDnev time.Time `json:"dat_dnev"`

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

	IDKatDok int `json:"id_kat_dok"`

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

	DatumSmene  time.Time  `json:"datum_smene"`  // već formatiran
	DatumDopune *time.Time `json:"datum_dopune"` // može biti null

	TipSmene       string  `json:"tip_smene"`
	TipSmeneDopune *string `json:"tip_smene_dopune"`

	Manipulacije *[]Manipulacija `json:"manipulacije"`
	ObavBeleske  *ObavBeleska    `json:"obav_beleske,omitempty"`
	TSU          *[]TSU          `json:"tsu"`
	TK           *[]TK           `json:"tk"`
	SOP          *[]SOP          `json:"sop"`
}

/************** Iskljucenja/Ukljucenja ****************/

type Manipulacija struct {
	IDDogSmene int `json:"id_dog_smene"`
	Rb         int `json:"rb"`

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
	Stavke []ManipView `json:"stavke"`
}

type ManipView struct {
	DopunaDaNe  string `json:"dopuna_da_ne"`
	Vrepoc      string `json:"vrepoc"`
	Vrezav      string `json:"vrezav"`
	RecenicaMan string `json:"recenica_man"`
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

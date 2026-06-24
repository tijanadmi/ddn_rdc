package models

type ShiftReport struct {
	Smena       Smena
	Dogadjaji   []DogadjajPDF
	Proizvodnja string `json:"proizvodnja"`
}

type DogadjajPDF struct {
	RbDog     string
	Naslov    string
	Podnaslov string
	TipDog    string
	Tip       string

	UzrokTekst string
	ManTekst   string
	Posledice  string

	Detalji      []DetaljT567
	Manipulacije []ObjekatView

	//  ZA ISKLJUČENJE
	Objekti []ObjekatView

	Grazlog string
	Razlog  string

	AngazovaniRukovalac *AngazovaniRukovalac
	ObavBeleske         *ObavBeleska
	ObavSlike           []ObavSlika
}

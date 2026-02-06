package models

type ListPiMMByParam struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Tipd      int    `json:"tipd"`
	Kom       string `json:"komisija"`
}

type Report struct {
	TipdGroups []TipdGroup
	StartDate  string
	EndDate    string
}

type TipdGroup struct {
	Tipd  string
	Naziv string
	Days  []DayGroup
}

type DayGroup struct {
	Date   string
	Events []EventGroup
}

type EventGroup struct {
	ID1   int
	Tekst string
	Rows  []DetailRow
}

type DetailRow struct {
	Vrepoc         string
	Vrezav         string
	Traj           string
	Objekat        string
	Polje          string
	ImePolja       string
	Snaga          string
	VrstaDogadjaja string
	GrupaUzroka    string
	Uzrok          string
	VremUsl        string
	GrRazlog       string
	Razlog         string
}

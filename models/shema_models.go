package models

type Shema struct {
	ID      int     `db:"ID" json:"id"`
	ImeDok  *string `db:"IME_DOK" json:"ime_dok"`
	Putanja *string `db:"PUTANJA" json:"putanja"`
	TipDok  *string `db:"TIP_DOK" json:"tip_dok"`
	IdSOrg  *int    `db:"ID_S_ORG" json:"id_s_org"`
	Datpri  *string `db:"DATPRI" json:"datpri"`
}

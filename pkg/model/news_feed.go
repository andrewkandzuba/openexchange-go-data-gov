package model

type Article struct {
	Type    string `json:"type"`
	Nid     int `json:"nid"`
	Label   string `json:"label"`
	Created int    `json:"created"`
	Update  int    `json:"updated"`
	Href    string `json:"href"`
	Body    string `json:"body"`
	Status  string `json:"release_status"`
	UUID    string `json:"uuid"`
	AdminOfficials []AdminOfficial `json:"admin_officials"`
}

type AdminOfficial struct {
	Id string `json:"id"`
	Label   string `json:"label"`
	Href    string `json:"href"`
}
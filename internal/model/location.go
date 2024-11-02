package model

type Country struct {
	Code      string               `json:"code"`
	Provinces map[string]*Province `json:"provinces"`
}

type Province struct {
	Code   string           `json:"code"`
	Cities map[string]*City `json:"cities"`
}

type City struct {
	Code string `json:"code"`
}

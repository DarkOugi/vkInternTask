package js

// get JSON

type GetUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetAdv {
	Security string `json:"security"`
	Header string `json:"header"`
	About   string `json:"about"`
	Picture string `json:"picture"`
	Price float32 `json:"price"`
}

type GetSecurity struct {
	Security string `json:"security"`
}
type GetSort struct {
	typeSort string `json:"type"`
	columnSort string `json:"column"`
}
type GetFilter struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
type GetAdvs struct {
	Security string `json:"security,omitempty"`
	Page int `json:"page"`
	Sort *GetSort `json:"sort,omitempty"`
	Filter *GetFilter `json:"filter,omitempty"`
}

// parse to JSON

type ToToken struct {
	Token string `json:"token"`
}

type ToAdvertisement struct {
	header  string `json:"header"`
	about   string `json:"about"`
	picture string `json:"picture"`
	price   float32 `json:"price"`
	author  string `json:"author,omitempty"`
}

type ToError struct {
	Errors string `json:"errors"`
}

package entity

type Filter struct {
	min float32
	max float32
}
type Advertisement struct {
	header  string
	about   string
	picture string
	price   float32
	login   string
}

type User struct {
	Login    string
	Password string
}

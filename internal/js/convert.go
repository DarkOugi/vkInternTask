package js

import (
	"encoding/json"
)

// get from JSON
func GetFromJSUser(user []byte) (*GetUser, error) {
	ujs := GetUser{}

	err := json.Unmarshal(user, &ujs)

	return &ujs, err
}

func GetFromJSAdv(adv []byte) (*GetAdv, error) {
	ajs := GetAdv{}

	err := json.Unmarshal(adv, &ajs)

	return &ajs, err
}
func GetFromJSAdvs(advs []byte) (*GetAdvs, error) {
	ajs := GetAdvs{}

	err := json.Unmarshal(advs, &ajs)

	return &ajs, err
}

// parse to JSON
func ToJSError(err string) ([]byte, error) {
	e := ToError{Errors: err}

	be, errM := json.Marshal(e)

	return be, errM
}
func ToJSToken(token string) ([]byte, error) {
	t := ToToken{Token: token}

	bt, err := json.Marshal(t)

	return bt, err
}

func ToJsAdvertisement(adv []*entity.Advertisement) ([]byte, error) {
	advJs := []*ToAdvertisement{}
	for _, el := range adv {
		advJs = append(advJs, &ToAdvertisement{
			header:  el.header,
			about:   el.about,
			picture: el.picture,
			price:   el.price,
			author:  el.author,
		})
	}

	ba, err := json.Marshal(advJs)

	return ba, err
}

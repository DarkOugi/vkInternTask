package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"vk/internal/entity"
	"vk/pkg/auth"
	"vk/pkg/jwt"
)

type Repository interface {
	GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error)
	InitUser(ctx context.Context, login, password string) error
	InitAdv(ctx context.Context, login, header, about, pucture string, price float32) error
	GetAdvs(ctx context.Context, page, pageSize int, colSort, typeSort string, filter *entity.Filter) ([]*entity.Advertisement, error)
}

var typeOrder = map[string]bool{
	"DESC": true,
	"ASC":  true,
}
var columnsSort = map[string]bool{
	"price": true,
	"date":  true,
}

var typeFilter = map[string]bool{
	"min": true,
	"max": true,
}

var columnsFilter = map[string]bool{
	"price": true,
}

const minLenPassword = 4
const pageSize = 50

var (
	

	ErrUnCorrectColSort   = errors.New("uncorrect cols for sort")
	ErrUnCorrectTypeSort  = errors.New("uncorrect type sort")

	ErrUnCorrectLogin = errors.New("bad login")
	ErrUnCorrectPass  = errors.New("bad pass")
	ErrBadPassword    = errors.New("bad password")
	ErrUnCorrectJWT   = errors.New("not correct JWT")
	ErrLoginBusy      = errors.New("login error")
	ErrLongAbout	  = errors.New("about too long")
	ErrSmallAbout		= errors.New("about too small")
	ErrLongHeader = errors.New("header too long")
	ErrSmallHeader= errors.New("header too small")
	ErrLongUrl= errors.New("picUrl too long")
	ErrSmallUrl= errors.New("picUrl too long")
	ErrPictureUrl     = errors.New("uncorrect url picture")
	ErrNegativePrice  = errors.New("price - negative")
)

type Service struct {
	rep Repository
}

func NewService(rep Repository) *Service {
	return &Service{
		rep: rep,
	}
}
func CheckPictureURL(url string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".webp", ".svg"}

	validExt := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			validExt = true
			break
		}
	}
	if !validExt {
		return false
	}

	resp, err := http.Head(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return false
	}

	return true
}
func (sv *Service) baseUserCheck(ctx context.Context, login, password string) (*entity.User, bool, error) {
	if ok := auth.CheckLogin(login); !ok {
		return nil, false, fmt.Errorf("%w: not valid login", ErrUnCorrectLogin)
	}

	if len(password) < minLenPassword {
		return nil, false, fmt.Errorf("%w: very short password", ErrUnCorrectPass)
	}
	password = auth.HashPassword(password)

	user, ok, err := sv.rep.GetUserInfo(ctx, login)
	if err != nil {
		return nil, false, fmt.Errorf("can't get user info: %w", err)
	}
	return user, ok, err
}

func (sv *Service) Reg(ctx context.Context, login, password string) (string, error) {
	// Логика регистрации
	// Если такой логин уже есть - возвращаем ошибку
	// Иначе регистрируем и возвращаем jwt tokken
	_, ok, err := sv.baseUserCheck(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("can't create user: %w", err)
	}
	if !ok {
		cErr := sv.rep.InitUser(ctx, login, password)
		if cErr != nil {
			return "", fmt.Errorf("can't create user: %w", cErr)
		}

		tokenJWT, errJWT := jwt.GenerateTokenAccess(login)
		if errJWT != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}
	return "", fmt.Errorf("can't create user: %w", ErrLoginBusy)
}

func (sv *Service) Auth(ctx context.Context, login, password string) (string, error) {
	// логика авторизации
	// 1 проверить логин и пароль на валидность
	// Проверяем существует ли такой пользователь
	// Если существует то
	// а) получаем его хэшированный пароль и сравниваем с переданным пользователейм в форме
	// (нужно его перед этим захэшировать)
	// если пароли не равны - ошибка, иначе вернем jwt токен
	// б) ошибка - такого пользователя нет

	u, ok, err := sv.baseUserCheck(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("can't create user: %w", err)
	}

	if !ok {
		return "", fmt.Errorf("This User not created")
	}

	if password == u.Password {
		tokenJWT, errJWT := jwt.GenerateTokenAccess(login)
		if errJWT != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}

	return "", fmt.Errorf("uncorrect login/password: %w %s", ErrBadPassword, password)
}

func (sv *Service) CreateAdv(ctx context.Context, sec, header, about, picture string, price float32) error {
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return fmt.Errorf("uncorrect security: %w", ErrUnCorrectJWT)
	}

	if len(header) > 56 {
		return fmt.Errorf("uncorrect header: %w", ErrLongHeader)
	}
	if len(header) < 10 {
		return fmt.Errorf("uncorrect header: %w", ErrSmallHeader)
	}

	if len(about) > 512 {
		return fmt.Errorf("uncorrect about: %w", ErrLongAbout)
	}
	if len(about) < 10 {
		return fmt.Errorf("uncorrect about: %w", ErrSmallAbout)
	}

	if len(picture) > 512 {
		return fmt.Errorf("uncorrect picUrl: %w", ErrLongUrl)
	} else {
		if len(picture) < 10 {
			return fmt.Errorf("uncorrect picUrl: %w", ErrSmallUrl)
		}
		flag := CheckPictureURL(picture)
		if !flag {
			return fmt.Errorf("uncorrect picUrl: %w", ErrPictureUrl)
		}
	}

	if price < 0 {
		return fmt.Errorf("uncorrect price: %w", ErrNegativePrice)
	}

	cErr := sv.rep.InitAdv(ctx, t.User, header, about, picture, price)
	if cErr != nil {
		return fmt.Errorf("can't create adv: %w", cErr)
	}
	return nil
}

func (sv *Service) ShowAdvs(ctx context.Context, sec string, page int,sorter js.GetSort, filter js.GetFilter) ([]*entity.Advertisement, error) {
	if _, ok := columnsSort[sorter.columnSort]; !ok {
		return nil, fmt.Errorf("can't sort: %w", ErrUnCorrectColSort)
	}
	if _, ok := typeOrder[sorter.typeSort]; !ok {
		return nil, fmt.Errorf("can't sort: %w", ErrUnCorrectTypeSort)
	}
	login := ''
	if len(sec) > 0 { 
		t, err := jwt.GetInfoFromToken(sec)
		if err != nil {
			return nil,fmt.Errorf("uncorrect security: %w", ErrUnCorrectJWT)
		}
		login = t.User
	}

	advs, errSql := sv.rep.GetAdvs(ctx, t.User, page, pageSize, colSort, typeSort, valueFilters)
	if errSql != nil {
		return nil, fmt.Errorf("can't get correct info: %w", errSql)
	}
	advsEn := []*entity.Advertisement{}
	for _,adv := range advs {
		if len(login) > 0 {
			if adv.login == login {
				advsEn = append(advsEn, *entity.Advertisement{
					header: adv.header,
					about:  adv.about,
					picture:adv.picture,
					price:adv.price,
					login:"Вы являетесь автором объявления",
				})
			} else {
				advsEn = append(advsEn, *entity.Advertisement{
					header: adv.header,
					about:  adv.about,
					picture:adv.picture,
					price:adv.price,
					login:"Вы не являетесь автором объявления",
				})
			}
		} else {
			advsEn = append(advsEn, *entity.Advertisement{
				header: adv.header,
				about:  adv.about,
				picture:adv.picture,
				price:adv.price,
				login:"",
			})
		}
	}
	return advsEn, nil
}


package server

import (
	"errors"
	"vk/internal/js"
	"vk/internal/service"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Server struct {
	service *service.Service
}

func NewServer(sv *service.Service) *Server {
	return &Server{service: sv}
}

func setError(ctx *fasthttp.RequestCtx, err error, codeErr int, errMsg string) {
	errStr, errConv := js.ToJSError(errMsg)
	if errConv != nil {
		log.Err(errConv).Msg("Error js.ToJSError")
	}
	log.Info().Err(err).Msgf("Error")
	ctx.SetStatusCode(codeErr)
	ctx.SetBody(errStr)
}

func (sv *Server) Auth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	us, err := js.GetFromJSUser(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	tokenJWT, errAuth := sv.service.Auth(ctx, us.Login, us.Password)
	if errAuth != nil {
		switch {
		case errors.Is(errAuth, service.ErrUnCorrectLogin):
			setError(ctx, errAuth, fasthttp.StatusBadRequest, "Невалидный логин")
			return
		case errors.Is(errAuth, service.ErrUnCorrectPass):
			setError(ctx, errAuth, fasthttp.StatusUnauthorized, "Невалидный пароль")
			return
		case errors.Is(errAuth, service.ErrBadPassword):
			setError(ctx, errAuth, fasthttp.StatusUnauthorized, "Неверный логин пароль")
			return
		default:
			setError(ctx, errAuth, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	tokenJS, errJS := js.ToJSToken(tokenJWT)
	if errJS != nil {
		setError(ctx, errJS, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(tokenJS)
}

func (sv *Server) Reg(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	us, err := js.GetFromJSUser(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	tokenJWT, errReg := sv.service.Reg(ctx, us.Login, us.Password)
	if errReg != nil {
		switch {
		case errors.Is(errReg, service.ErrLoginBusy):
			setError(ctx, errReg, fasthttp.StatusBadRequest, "Данный логин уже занят")
			return
		case errors.Is(errReg, service.ErrUnCorrectPass):
			setError(ctx, errReg, fasthttp.StatusUnauthorized, "Пароль несоответствует требованиям")
			return
		case errors.Is(errReg, service.ErrUnCorrectLogin):
			setError(ctx, errReg, fasthttp.StatusUnauthorized, "Логин несоответствует требованиям")
			return
		default:
			setError(ctx, errReg, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	tokenJS, errJS := js.ToJSToken(tokenJWT)
	if errJS != nil {
		setError(ctx, errJS, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(tokenJS)
}

func (sv *Server) CreateAdv(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	adv, err := js.GetFromJSAdv(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	errAdv := sv.service.CreateAdv(ctx, adv.Security, adv.Header, adv.About, adv.Picture, adv.Price)
	if errAdv != nil {
		switch {
		case errors.Is(errAdv, service.ErrUnCorrectJWT):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Невалидный токен")
			return

		case errors.Is(errAdv, service.ErrLongAbout):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком длинное описание объявления")
			return
		case errors.Is(errAdv, service.ErrSmallAbout):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком короткое описание объявления")
			return
		case errors.Is(errAdv, service.ErrLongHeader):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком длинный заголовок")
			return
		case errors.Is(errAdv, service.ErrSmallHeader):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком короткий заголовок")
			return
		case errors.Is(errAdv, service.ErrLongUrl):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком длинная ссылка")
			return
		case errors.Is(errAdv, service.ErrSmallUrl):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Слишком короткая ссылка")
			return
		case errors.Is(errAdv, service.ErrPictureUrl):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Невалидная ссылка")
			return
		case errors.Is(errAdv, service.ErrNegativePrice):
			setError(ctx, errAdv, fasthttp.StatusUnauthorized, "Цена не может быть отрицательной")
			return

		default:
			setError(ctx, errAdv, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte{})
}

func (sv *Server) GetAllAdv(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	advs, err := js.GetFromJSAdvs(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	advs, errAdvs := sv.service.ShowAdvs(ctx, advs.Security, advs.Page, advs.Sort, advs.Filter)
	if errAdvs != nil {
		switch {
		case errors.Is(errAdvs, service.ErrUnCorrectJWT):
			setError(ctx, errAdvs, fasthttp.StatusUnauthorized, "Невалидный токен")
			return
		case errors.Is(errAdvs, service.ErrUnCorrectTypeSort):
			setError(ctx, errAdvs, fasthttp.StatusUnauthorized, "Не известный тип сортировки")
			return
		case errors.Is(errAdvs, service.ErrUnCorrectColSort):
			setError(ctx, errAdvs, fasthttp.StatusUnauthorized, "По данному полю нельзя сортировать")
			return
		default:
			setError(ctx, errAdvs, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	jsonData, errJS := js.ToJsAdvertisement(advs)
	if errJS != nil {
		setError(ctx, errJS, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jsonData)
}

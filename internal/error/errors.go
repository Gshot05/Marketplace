package errors2

import "errors"

var (
	ErrNotCustomer        = errors.New("Только заказчики имеют доступ к этой функции!")
	ErrNotPerformer       = errors.New("Только исполнители имеют доступ к этой функции!")
	ErrEmptyName          = errors.New("Имя не может быть пустым!")
	ErrEmptyRole          = errors.New("Некорректная роль!")
	ErrWrongJson          = errors.New("Неверный формат запроса!")
	ErrWrongUpdateOffer   = errors.New("Оффер не найден или принадлежит не вам!")
	ErrWrongUpdateService = errors.New("Услуга не найдена или принадлежит не вам!")
	ErrNotFindService     = errors.New("Сервис не найден!")
	ErrNoAuth             = errors.New("Нет авторизации!")
	ErrBadToken           = errors.New("Фиговый токен!")
	ErrTokenExpired       = errors.New("Токен истёк!")
	ErrWrongPassOrLog     = errors.New("Неверный логин или пароль!")
	ErrEmptyToken         = errors.New("Пустой токен!")
)

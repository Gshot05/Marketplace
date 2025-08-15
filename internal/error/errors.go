package errors

import "errors"

var (
	ErrNotCustomer        = errors.New("Только заказчики имеют доступ к этой функции")
	ErrNotPerformer       = errors.New("Только исполнители имеют доступ к этой функции")
	ErrEmptyName          = errors.New("Имя не может быть пустым")
	ErrEmptyRole          = errors.New("Некорректная роль")
	ErrWrongJson          = errors.New("Неверный формат запроса")
	ErrWrongUpdateOffer   = errors.New("Оффер не найден или принадлежит не вам!")
	ErrWrongUpdateService = errors.New("Услуга не найден или принадлежит не вам!")
	ErrNotFindService     = errors.New("Сервис не найден!")
)

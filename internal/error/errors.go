package errors2

import "errors"

var (
	//ErrNot group
	ErrNotCustomer    = errors.New("Только заказчики имеют доступ к этой функции!")
	ErrNotPerformer   = errors.New("Только исполнители имеют доступ к этой функции!")
	ErrNotFindService = errors.New("Сервис не найден!")

	//ErrEmpty group
	ErrEmptyName        = errors.New("Имя не может быть пустым!")
	ErrEmptyRole        = errors.New("Некорректная роль!")
	ErrEmptyToken       = errors.New("Пустой токен!")
	ErrEmptyTitle       = errors.New("Заголовок не может быть пустым!")
	ErrEmptyDescription = errors.New("Описание не может быть пустым!")
	ErrEmptyPrice       = errors.New("Цена не может быть пустой!")
	ErrEmptyOffers      = errors.New("Офферов пока нет:(")
	ErrEmptyServices    = errors.New("Услуг пока нет:(")
	ErrEmptyFav         = errors.New("Избранное пока пусто:(")

	//ErrWrong group
	ErrWrongJson          = errors.New("Неверный формат запроса!")
	ErrWrongUpdateOffer   = errors.New("Оффер не найден или принадлежит не вам!")
	ErrWrongUpdateService = errors.New("Услуга не найдена или принадлежит не вам!")
	ErrWrongPassOrLog     = errors.New("Неверный логин или пароль!")

	//ErrAuth group
	ErrNoAuth       = errors.New("Нет авторизации!")
	ErrBadToken     = errors.New("Фиговый токен!")
	ErrTokenExpired = errors.New("Токен истёк!")
	ErrEmailSent    = errors.New("Не удалось отправить письмо!")
)

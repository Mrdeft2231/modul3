package auth

import "errors"

var (
	ErrSamePassword      = errors.New("новый пароль не должен совпадать со старым")
	ErrIncorrectPassword = errors.New("неверный старый пароль")
	ErrPassword          = errors.New("неверный пароль пользователя")
	ErrBlock             = errors.New("вы достигли маскимальное колчество попыток входа, пользователь заблокирован")
	ErrBlockDate         = errors.New("пользователь заблокирован")
	ErrBlockStatus       = errors.New("вас заблокировали, вход невозможен")
	ErrUniqName          = errors.New("такое имя существует")
)

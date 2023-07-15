package models //package models

type User struct { //описание структуры User (пользователь\клиент)
	IdUser    int    //id в БД
	UserLogin string //логин
	UserPswrd string //пароль
	UserMail  string //почта
	UserName  string //имя и фамилия
}

type Hotel struct { //описание структуры Hotel (отель)
	IdHotel    int    //id в БД
	HotelTitle string //название отеля
}

type Number struct { //описание структуры Number (номер)
	IdHotel   int    //id отеля, к которому он принадлежит
	IdNumber  int    //id номера в БД
	HotelName string //название отеля, к которому он принадлежит
	Size      int    //на сколько человек рассчитан номер
	Status    string //статус номера в БД (заброинрован он или свободен)
	Payment   string //статус оплаты
}

type BookingRequest struct { //описание структуры BookingRequest (запроса на бронирование)
	Id        int    //id запроса в БД
	IdFrom    int    //id пользователя, от которого поступил запрос
	IdTo      int    //id менеджера, которому поступил запрос
	IdHotel   int    //id отеля
	IdNumber  int    //id номера в этом отеле, который пользователь собирается забронировать
	HotelName string //название отеля
	Size      int    //на сколько человек рассчитан номер
	Status    string //статус номера (свободен он или нет)
	Price     int    //цена для броинрования номера
	PayMethod string //способ оплаты
	Nights    int    //кол-во ночей
}

type Admin struct { //структура Admin (вложение структуры User)
	UserStruct User 
}

type Manager struct { //структура Manager (вложение структур User и BookingRequest)
	UserStruct User
}

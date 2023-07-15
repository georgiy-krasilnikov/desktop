package db //package db

import ( //импортирование пакетов
	"context"

	"main/models" //импортирование локального package'a "models"
)

func (dbase *DB) AddBookingRequest(ctx context.Context, BR models.BookingRequest) error { //функция добавления запроса в БД
	query := `INSERT INTO booking_requests(id, id_from, id_to, id_hotel, hotel_name, id_room, room_size, br_status, price, pay_method, number_of_nights)
					VALUES 
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	//запрос для postgreSQL
	rows, err := dbase.cnct.Query(ctx, query, 1, BR.IdFrom, BR.IdTo, BR.IdHotel, BR.HotelName, BR.IdNumber, BR.Size, BR.Status, BR.Price, BR.PayMethod, BR.Nights)
	//опр. данных структуры BookingRequest для занесения их в БД
	if err != nil { //если произошла ошибка, то возвращаем её - иначе, возвращаем пустую ошибку
		return err
	}
	defer rows.Close()
	return nil
}

func (dbase *DB) ScanBooking(ctx context.Context) (*models.BookingRequest, error) { //функция считывания данных бронирования из БД
	var BR models.BookingRequest //объявление BookingRequest-структуры BR
	//запрос для postgreSQL
	rows, err := dbase.cnct.Query(ctx, "SELECT id_from, hotel_name, id_room, room_size, br_status, price, pay_method, number_of_nights FROM booking_requests")
	if err != nil {
		return nil, err //если произошла ошибка, то возвращаем её и пустую структуру
	}
	for rows.Next() {
		defer rows.Close()
		//считывание данных из БД в структуру BR
		err := rows.Scan(&BR.IdFrom, &BR.HotelName, &BR.IdNumber, &BR.Size, &BR.Status, &BR.Price, &BR.PayMethod, &BR.Nights)
		if err != nil {
			return nil, err //если произошла ошибка, то возвращаем её и пустую структуру
		}
	}
	return &BR, nil //возвращаем адрес на структуру BR и пустуюу ошибку, в случае успешного выполнения запроса
}

package db //package db

import ( //импортирование пакетов
	"context"
	"strconv"

	"main/models" //импортирование локального пакета "models", где описаны необходимые структуры
)

func (dbase *DB) RoomData(ctx context.Context, hotel models.Hotel, number string) (*models.Number, error) { //функция для получения данных о конкретном номере из БД
	var room models.Number         //объявление Number-структуры room из package'a models
	num, _ := strconv.Atoi(number) //конвертация number из string в int
	rows, err := dbase.cnct.Query(ctx, "SELECT hotel_id, id, room_size, status FROM rooms_of_hotels WHERE hotel_id = $1 AND room_size = $2", hotel.IdHotel, num)
	//запрос для postgreSQL
	if err != nil {
		return nil, err //возвращаем пустую структуру и ошибку, если она есть
	}
	for rows.Next() {
		defer rows.Close() //функция, закрывающая канал после обработки запроса
		//сканирование данных из БД в структуру Number
		if err := rows.Scan(&room.IdHotel, &room.IdNumber, &room.Size, &room.Status); err != nil {
			return nil, err //возвращаем пустую структуру и ошибку, если она есть
		}
	}
	return &room, nil //возвращаем адрес памяти и пустую ошибку, если все прошло успешно
}

func (dbase *DB) UpdateStatus(ctx context.Context, br models.BookingRequest) error { //функция обновления статуса номера
	roomid := br.IdNumber                                                       //получаем id номера, статус которого нужно обновить
	query := `UPDATE rooms_of_hotels SET status = 'Забронирован' WHERE id = $1` //запрос для postgreSQL
	rows, err := dbase.cnct.Query(ctx, query, roomid)                           //выполнение запроса
	if err != nil {                                                             //если есть ошибка, возвращаем её
		return err
	}
	defer rows.Close() //закрываем канал
	return nil         //если запрос выолнился успешно, возваращем пустую ошибку
}

func (dbase *DB) UpdatePayStatus(ctx context.Context, br models.BookingRequest) error { //функция обновления статуса номера
	roomid := br.IdNumber                                                         //получаем id номера, статус которого нужно обновить
	query := `UPDATE rooms_of_hotels SET paymentstatus = 'Оплачен' WHERE id = $1` //запрос для postgreSQL
	rows, err := dbase.cnct.Query(ctx, query, roomid)                             //выполнение запроса
	if err != nil {                                                               //если есть ошибка, возвращаем её
		return err
	}
	defer rows.Close() //закрываем канал
	return nil         //если запрос выолнился успешно, возваращем пустую ошибку
}

func (dbase *DB) CheckRoom(ctx context.Context) (map[int]int, error) { //функция, опр. статуса всех номеров в каждом отеле
	data := make(map[int]int, 4) //создаем хэш-таблицу для каждого отеля
	for i := 1; i < 5; i++ {     //цикл для выполнения каждого запроса
		var count int
		//запрос для postgreSQL
		rows, err := dbase.cnct.Query(context.Background(), "SELECT COUNT(*) AS count FROM rooms_of_hotels WHERE hotel_id = $1 AND status = 'Свободен'", i)
		if err != nil { //если есть ошибка, возвращаем пустые хэш-таблицы и ошибку
			return nil, err
		}
		for rows.Next() {
			defer rows.Close()                        //функция закрытия канала для БД после выполнения запроса
			if err := rows.Scan(&count); err != nil { //сканирование COUNT из запроса в count
				return nil, err //если есть ошибка, возвращаем пустые хэш-таблицы и ошибку
			}
			data[i] = count //заносим в хэш-таблицу i отеля кол-во свободных номеров
		}
	}
	return data, nil //если все прошло успешно, возвращаем хэш-таблицы и пустую ошибку
}

func (dbase *DB) DeleteRooms(ctx context.Context, hotel models.Hotel) error {
	rows, err := dbase.cnct.Query(context.Background(), "DELETE FROM rooms_of_hotels WHERE hotel_id = $1", hotel.IdHotel)
	//запрос для postgreSQL
	if err != nil {
		return err //возвращение ошибки, если она произошла
	}
	defer rows.Close() //закрываем канал подключения к БД
	return nil         //возвращаем пустую ошибку в случае успешного выполнения запроса
}

func (dbase *DB) AddRoomsOfNewHotel(ctx context.Context, hotel models.Hotel, size, id int) error {
	query := `INSERT INTO rooms_of_hotels(hotel_id, id, room_size, status)
			VALUES 
			($1, $2, $3, $4)`
	rows, err := dbase.cnct.Query(ctx, query, hotel.IdHotel, id, size, "Свободен")
	if err != nil {
		return err //если произошла ошибка, возвращаем её
	}
	defer rows.Close()
	return nil
}

package db //package db

import ( //импортирование пакетов
	"context"

	"main/models" //импортирование локального пакета "models", где описаны необходимые структуры
)

func (dbase *DB) HotelsData(ctx context.Context) ([]models.Hotel, error) { //функция получения данных всех отелей из БД
	var number int
	rows1, err := dbase.cnct.Query(context.Background(), "SELECT COUNT(*) AS number FROM list_of_hotels")
	//запрос для postgreSQL для получения кол-ва отелей
	if err != nil {
		return nil, err //если произошла ошибка, возвращаем пустые структуры и ошибку
	}
	for rows1.Next() {
		defer rows1.Close()        //функция закрытия канала для БД после выполнения запроса
		err := rows1.Scan(&number) //сканируем COUNT в number
		if err != nil {
			return nil, err //если произошла ошибка, возвращаем пустые структуры и ошибку
		}
	}
	H := make([]models.Hotel, number) //создаем массив H из структур Hotel
	rows2, err := dbase.cnct.Query(context.Background(), "SELECT * FROM list_of_hotels")
	//запрос для postgreSQL
	if err != nil {
		return nil, err //если произошла ошибка, возвращаем пустые структуры и ошибку
	}
	i := 0
	for rows2.Next() {
		defer rows2.Close()
		err := rows2.Scan(&H[i].IdHotel, &H[i].HotelTitle) //сканируем данные из БД в структуру H под индексом i
		if err != nil {
			return nil, err //если произошла ошибка, возвращаем пустые структуры и ошибку
		}
		i++
	}
	tmp := H[3]
	H[3] = H[1]
	H[1] = tmp
	return H, nil //возвращаем массив структур и пустую ошибку, если все прошло успешно
}

func (dbase *DB) DeleteHotel(ctx context.Context, hotel models.Hotel) error { //функция удаления отеля из списка
	rows, err := dbase.cnct.Query(context.Background(), "DELETE FROM list_of_hotels WHERE id = $1", hotel.IdHotel)
	//запрос для postgreSQL
	if err != nil {
		return err //возвращение ошибки, если она произошла
	}
	defer rows.Close() //закрываем канал подключения к БД
	return nil         //возвращаем пустую ошибку в случае успешного выполнения запроса
}

func (dbase *DB) AddNewHotel(ctx context.Context, hotel models.Hotel) error {
	query := `INSERT INTO list_of_hotels(id, hotel_name)
			VALUES 
			($1, $2)`
	rows, err := dbase.cnct.Query(ctx, query, hotel.IdHotel, hotel.HotelTitle)
	//данные структуры (нового пользователя), которые нам нужно занести в БД
	if err != nil {
		return err //если произошла ошибка, возвращаем её
	}
	defer rows.Close()
	return nil
}

package db //package db

import ( //импортирование пакетов
	"context"

	"main/models" //импортирование локального package'a "models" с описанием структур
)

func (dbase *DB) UsersData(ctx context.Context) ([]models.User, error) { //функция получения данных пользователей из БД
	var number int //запрос для postgreSQL на опр. кол-ва пользователей в БД
	rows1, err := dbase.cnct.Query(context.Background(), "SELECT COUNT(*) AS number FROM clients_database")
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
	users := make([]models.User, number)                                                  //создание массива сusers из структур User
	rows, err := dbase.cnct.Query(context.Background(), "SELECT * FROM clients_database") //запрос для postgreSQL
	if err != nil {
		return nil, err //если произошла ошибка, возвразащаем пустой массив и ошибку
	}
	i := 0
	for rows.Next() {
		defer rows.Close() //функция для закрытия канала подключения к БД после выполнения запроса
		//сканируем данные клиентов из БД в i структуру массива users
		err := rows.Scan(&users[i].IdUser, &users[i].UserLogin, &users[i].UserPswrd, &users[i].UserMail, &users[i].UserName)
		if err != nil {
			return nil, err
		}
		if i+1 < len(users) { //проверка, чтобы не было ошибки длины массива
			i++
		}
	}
	return users, nil //если все прошло успешно, возвращаем массив и пустую ошибку
}

func (dbase *DB) AddClient(ctx context.Context, U models.User) error { //функция добавления нового клиента в БД
	//запрос для postgreSQL
	query := `INSERT INTO clients_database(id, client_login, client_pswrd, client_mail, client_name)
			VALUES 
			($1, $2, $3, $4, $5)`
	rows, err := dbase.cnct.Query(ctx, query, 4, U.UserLogin, U.UserPswrd, U.UserMail, U.UserName)
	//данные структуры (нового пользователя), которые нам нужно занести в БД
	if err != nil {
		return err //если произошла ошибка, возвращаем её
	}
	defer rows.Close()
	return nil
}

func (dbase *DB) ScanClient(ctx context.Context, id int) (*models.User, error) { //функция считывания данных из БД опр. клиента
	var user models.User //объявление User-структуры user
	rows, err := dbase.cnct.Query(context.Background(), "SELECT * FROM clients_database WHERE id = $1", id)
	//запрос для postgreSQL
	if err != nil {
		return nil, err //если произошла ошибка, возвращаем её и пустую структуру
	}
	for rows.Next() {
		defer rows.Close()
		//сканирование данных из БД в User-структуру user
		err := rows.Scan(&user.IdUser, &user.UserLogin, &user.UserPswrd, &user.UserMail, &user.UserName)
		if err != nil {
			return nil, err //если произошла ошибка, возвращаем её и пустую структуру
		}
	}
	return &user, nil //если все прошло успешно, возвращаем адрес на структуру и пустую ошибку
}

package db //package db

import ( //импортирование пакетов
	"context"

	"main/models" //импортирование локального package'a "models"
)

func (dbase *DB) AdminData(ctx context.Context) (*models.Admin, error) { //функция считывания данных администратора из БД
	var admin models.Admin //объявление Admin-структуры admin
	//запрос для postgreSQL
	rows, err := dbase.cnct.Query(context.Background(), "SELECT * FROM administrator")
	if err != nil {
		return nil, err //если произошла ошибка, возвращаем её и пустую структуру
	}
	for rows.Next() {
		defer rows.Close()
		//считывание данных из БД в структуру admin
		err := rows.Scan(&admin.UserStruct.IdUser, &admin.UserStruct.UserLogin, &admin.UserStruct.UserPswrd, &admin.UserStruct.UserMail, &admin.UserStruct.UserName)
		if err != nil {
			return nil, err //если произошла ошибка, возвращаем её и пустую структуру
		}
	}
	return &admin, nil //если не было ошибкиЮ возвращаем адрес на структуру и пустую ошибку
}

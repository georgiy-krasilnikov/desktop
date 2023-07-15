package db //package db

import ( //импортирование пакетов
	"context"

	"main/models" //импортирование локального package'а "models"
)

func (dbase *DB) ManagerData(ctx context.Context) ([]models.Manager, error) { //функция считывания данных менеджера из БД
	number := 0
	rows1, err := dbase.cnct.Query(context.Background(), "SELECT COUNT(*) AS number FROM manager")
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
	mng := make([]models.Manager, number) //объявление manager-структуры mng
	//запрос для postgreSQL
	rows, err := dbase.cnct.Query(ctx, "SELECT * FROM manager")
	if err != nil {
		return nil, err //если произошла ошибка, возвращаем пустую структуру и ошибку
	}
	i := 0
	for rows.Next() {
		defer rows.Close()
		//сканирование данных менеджера из БД в структуру mng
		err := rows.Scan(&mng[i].UserStruct.IdUser, &mng[i].UserStruct.UserLogin, &mng[i].UserStruct.UserPswrd, &mng[i].UserStruct.UserName, &mng[i].UserStruct.UserMail)
		if err != nil {
			return nil, err //если произошла ошибка, возвращаем пустую структуру и ошибку
		}
		if i < number {
			i++
		}
	}
	return mng, nil //если нет ошибки, то возаращем адрес структуры и пустую ошибку
}

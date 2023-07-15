package db //package db

import ( //импортирование пакетов
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5" //импортирование драйвера pgx для работы с posgtreSQL
)

type DB struct { //определение структуры DB, состоящей из подключения к БД (*pgx.Conn)
	cnct *pgx.Conn
}

func (dbase *DB) New(ctx context.Context) error { //создание канала для подключения к БД
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", "postgres", "empty123", "localhost", "5432", "mssd")
	//создание строки для подключения
	x := 0 //кол-во попыток подключения
	var err error
	for x < 5 { //несколько попыток
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //опр. времени на выполнение функции (5 секунд)
		defer cancel()
		dbase.cnct, err = pgx.Connect(ctx, dsn) //попытка выполнения подключения к БД (5 секунд)
		x++                                     //счетчик попыток
	}
	if err != nil { //если так и не смогли подключиться к БД, то возвращаем ошибку подключения
		return err
	}
	return nil //иначе - пустую ошибку
}

package main

import ( //импортирование пакетов
	"context"

	"fyne.io/fyne" //импортировоание пакетов для работы с граф. интерфейсом
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	"main/db" //импортирование локального package'a для работы с базой данных и созданием окон
)

func main() {
	ctx := context.Background()            //базовый контекст (пустой)
	var dbase db.DB                        //структура из package'a db - состоит лишь из подключения к БД (*pgx.Conn)
	if err := dbase.New(ctx); err != nil { //подключаемся к БД, проверка на ошибку подключения
		panic(err)
	}
	app := app.New()                  //создание граф. приложения
	menu := app.NewWindow("Турфирма") //создаем окно
	icon, _ := fyne.LoadResourceFromPath("icons/icon.png")
	menu.SetIcon(icon)                  //устанавливаем иконку
	menu.Resize(fyne.NewSize(400, 250)) //опр. размер окна
	menu.SetContent(widget.NewVBox(     //устанавливаем виджеты для окна
		widget.NewButton("Войти", func() { //кнопка для "входа"
			dbase.Authentication(ctx, app) //функция авторизации из package'a db
			menu.Close()
		}),
		widget.NewButton("Зарегистрироваться", func() { //кнопка для "регистрации"
			dbase.Registration(ctx, app, menu) //функция регистрации package'a db
		}),
		widget.NewButton("Выйти", func() { //кнопка для "выхода"
			app.Quit()
		}),
	))
	menu.ShowAndRun() //работа изначального окна
}

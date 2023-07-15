package db //package db

import ( //импортирование пакетов
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne" //импортирование пакета fyne для работы с граф. интерфейсом
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"

	"main/models" //импортирование локального package'a "models"
)

func (dbase *DB) Authentication(ctx context.Context, app fyne.App) { //окно стандартной авторизации
	users, err := dbase.UsersData(ctx) //создаем массив пользователей
	if err != nil {                    //если есть ошибка, выводим её и выходим из программы
		panic(err)
	}
	count := 0 //счетчик попыток ввода логина/пароля
	check := widget.NewLabel("")
	login := widget.NewEntry()                       //виджет ввода логина
	password := widget.NewPasswordEntry()            //виоджет ввода пароля
	enter := app.NewWindow("Аутентификация клиента") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/reg.png")
	enter.SetIcon(icon)
	enter.Resize(fyne.NewSize(300, 200))
	enter.SetContent(widget.NewVBox( //устанавливаем виджеты для окна enter
		widget.NewLabel("Логин:"),
		login,
		widget.NewLabel("Пароль:"),
		password,
		widget.NewButton("Войти", func() { //оработка кнопки
			for _, u := range users { //цикл проверки введенных данных с данными всех пользоваталей
				if login.Text == u.UserLogin && password.Text == u.UserPswrd { //если есть совпадение
					dbase.HotelChoice(ctx, u, app, enter) //переходим дальше
					break                                 //выход
				} else if count == 0 { //первая попытка
					check.SetText("Неправильный пароль или логин!") //сообщение об ошибке
					count++
					break //выход
				} else { //после первой попытки переходим к регистрации
					dbase.Registration(ctx, app, enter)
					break //выход
				}
			}
		}),
		check,
	))
	enter.Show() //показ окна
}

func (dbase *DB) Registration(ctx context.Context, app fyne.App, menu fyne.Window) { //окно регистрации
	var u models.User         //объявление User-структуры u
	newl := widget.NewEntry() //виджеты ввода для данных нового пользователя
	newpswrd := widget.NewPasswordEntry()
	newmail := widget.NewEntry()
	newname := widget.NewEntry()
	check := widget.NewLabel("")
	registration := app.NewWindow("Регистрация") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/add.png")
	registration.SetIcon(icon)
	registration.Resize(fyne.NewSize(300, 200))
	registration.SetContent(widget.NewVBox(
		widget.NewLabel("Логин:"),
		newl,
		widget.NewLabel("Пароль:"),
		newpswrd,
		widget.NewLabel("Почта:"),
		newmail,
		widget.NewLabel("Фамилия и имя:"),
		newname,
		widget.NewButton("Зарегистрироваться", func() { //
			u = models.User{ //присваивание значений полям структуры u
				UserLogin: newl.Text,
				UserPswrd: newpswrd.Text,
				UserMail:  newmail.Text,
				UserName:  newname.Text,
			}
			if err := dbase.AddClient(ctx, u); err != nil { //добавление нового клиента в БД
				panic(err) //если есть ошибка, выводим её и выходим из программы
			}
			check.SetText("Вы зарегистрировались!")
		}),
		check,
		widget.NewButton("Далее", func() { //обработка кнопки
			dbase.HotelChoice(ctx, u, app, registration) //переходим дальше
		}),
	))
	dialog.ShowInformation("Здравствуйте", "Пожалуйста, зарегистрируйтесь!", registration) //показ диалогового окна
	registration.Show()
	menu.Close()
}

func (dbase *DB) HotelChoice(ctx context.Context, u models.User, app fyne.App, w fyne.Window) { //окно выбора отеля
	hotels, err := dbase.HotelsData(ctx) //получение данных об отелях из БД
	if err != nil {
		panic(err) //проверка на наличие ошибки
	}
	var ch string
	img1 := canvas.NewImageFromFile("radisson/rad.jpg") //получение фотографий отелей
	img2 := canvas.NewImageFromFile("azimut/azimut.jpg")
	img3 := canvas.NewImageFromFile("vremena/vremena.jpg")
	img4 := canvas.NewImageFromFile("standart/stdart.jpg")
	img5 := canvas.NewImageFromFile("vega/vega.jpg")
	hotellist := app.NewWindow("Выбор отеля") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/icon.png")
	hotellist.SetIcon(icon)
	hotellist.Resize(fyne.NewSize(500, 500))
	var hotel *widget.RadioGroup //объявление списка с одним вариантом выбора
	switch len(hotels) {         //switch-case с занесением названий отелей в список, в зависимости от кол-ва отелей в БД
	case 5:
		hotel = widget.NewRadioGroup([]string{hotels[0].HotelTitle, hotels[1].HotelTitle, hotels[2].HotelTitle, hotels[3].HotelTitle, hotels[4].HotelTitle}, func(s string) {})
	case 4:
		hotel = widget.NewRadioGroup([]string{hotels[0].HotelTitle, hotels[1].HotelTitle, hotels[2].HotelTitle, hotels[3].HotelTitle}, func(s string) {})
	case 3:
		hotel = widget.NewRadioGroup([]string{hotels[0].HotelTitle, hotels[1].HotelTitle, hotels[2].HotelTitle}, func(s string) {})
	default:
		hotel = widget.NewRadioGroup([]string{hotels[0].HotelTitle, hotels[1].HotelTitle}, func(s string) {})
	}
	btn := widget.NewButton("Далее", func() { //обработка кнопки
		ch = hotel.Selected                                  //выбранное название отелч
		dbase.RoomChoice(ctx, hotels, u, app, hotellist, ch) //переходим дальше
		hotellist.Close()
	})
	box := container.NewVBox( //контейнер со списком
		hotel,
	)
	switch len(hotels) {
	case 5:
		hotellist.SetContent((container.NewGridWithColumns(2, img1, img2, img3, img4, img5, box, btn))) //установка виджетов для окна list
	case 4:
		hotellist.SetContent((container.NewGridWithColumns(2, img1, img2, img3, img4, box, btn))) //установка виджетов для окна list
	case 3:
		hotellist.SetContent((container.NewGridWithColumns(2, img1, img2, img3, box, btn))) //установка виджетов для окна list
	case 2:
		hotellist.SetContent((container.NewGridWithColumns(2, img1, img2, box, btn))) //установка виджетов для окна list
	}
	hotellist.Show()
	w.Close()
}

func (dbase *DB) RoomChoice(ctx context.Context, hotels []models.Hotel, u models.User, app fyne.App, hotellist fyne.Window, ch string) {
	//функция подбора номера
	var hotel models.Hotel //объявление Hotel-структуры hotel
	for i := 0; i < len(hotels); i++ {
		if ch == hotels[i].HotelTitle { //если название выбранного отеля совпадает с название i отеля
			hotel = hotels[i] //присваиваем структуре hotel значение полей i отеля
		}
	}
	roomchoice := app.NewWindow("Бронирование номера") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/check.png")
	roomchoice.SetIcon(icon)
	roomchoice.Resize(fyne.NewSize(300, 200))
	number := widget.NewEntry()
	data := widget.NewEntry()
	roomchoice.SetContent(widget.NewVBox( //установка виджетов
		widget.NewLabel("Введите количество гостей:"),
		number,
		widget.NewLabel("Введите количество ночей:"),
		data,
		widget.NewButton("Подобрать номер", func() { //обработка кнопки
			r, err := dbase.RoomData(ctx, hotel, number.Text) //получение адреса в памяти структуры Number
			room := *r                                        //обращаемся через указатель
			if err != nil {                                   //если есть ошибка, выводим её и выходим из программы
				panic(err)
			}
			info := dbase.SetContent(ctx, room, hotel, u, data.Text, app) //функция установки виджетов
			info.Show()
			roomchoice.Close()
		}),
	))
	roomchoice.Show()
}

// функция установки виджетов в зависимости от отеля и номера
func (dbase *DB) SetContent(ctx context.Context, room models.Number, hotel models.Hotel, u models.User, data string, app fyne.App) fyne.Window {
	z, _ := strconv.Atoi(data)                      //конвертация из string в int
	roomset := app.NewWindow("Бронирование номера") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/check.png")
	roomset.SetIcon(icon)
	roomset.Resize(fyne.NewSize(600, 600))
	box := widget.NewLabel("Мы подобрали Вам номер:")
	switch hotel.IdHotel { //switch-case по отелям
	case 1:
		switch room.Size { //switch case по номерам
		case 1: //для первого номера 1 отеля
			t := (rand.Intn(1000) + 1500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("radisson/1.jpg")
			img2 := canvas.NewImageFromFile("radisson/1 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		case 2: //для второго номера
			t := (rand.Intn(1500) + 2000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("radisson/2.jpg")
			img2 := canvas.NewImageFromFile("radisson/2 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() {
				dbase.Authentication(ctx, app) //выход в окно авторизации
			}) //обработка кнопки
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		default: //для третьего
			t := (rand.Intn(2000) + 4000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("radisson/3.jpg")
			img2 := canvas.NewImageFromFile("radisson/3 (1).jpg")
			img3 := canvas.NewImageFromFile("radisson/3 (2).jpg")
			img4 := canvas.NewImageFromFile("radisson/3 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		}
	case 2: //второй отель
		switch room.Size { //первый номер
		case 1:
			t := (rand.Intn(1000) + 2000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("azimut/1.jpg")
			img2 := canvas.NewImageFromFile("azimut/1 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		case 2:
			t := (rand.Intn(2000) + 2500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("azimut/2.jpg")
			img2 := canvas.NewImageFromFile("azimut/2 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		default:
			t := (rand.Intn(2500) + 3500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("azimut/3.jpg")
			img2 := canvas.NewImageFromFile("azimut/3 (1).jpg")
			img3 := canvas.NewImageFromFile("azimut/3 (2).jpg")
			img4 := canvas.NewImageFromFile("azimut/3 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		}
	case 3:
		switch room.Size {
		case 1:
			t := (rand.Intn(2000) + 2500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vremena/1.jpg")
			img2 := canvas.NewImageFromFile("vremena/1 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		case 2:
			t := (rand.Intn(2000) + 3000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vremena/2.jpg")
			img2 := canvas.NewImageFromFile("vremena/2 (1).jpg")
			img3 := canvas.NewImageFromFile("vremena/2 (2).jpg")
			img4 := canvas.NewImageFromFile("vremena/2 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		default:
			t := (rand.Intn(3000) + 3000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vremena/3.png")
			img2 := canvas.NewImageFromFile("vremena/3 (1).png")
			img3 := canvas.NewImageFromFile("vremena/3 (2).png")
			img4 := canvas.NewImageFromFile("vremena/3 (3).png")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		}
	case 4:
		switch room.Size {
		case 1:
			t := (rand.Intn(1000) + 1500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("standart/1.jpg")
			img2 := canvas.NewImageFromFile("standart/1 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		case 2:
			t := (rand.Intn(1500) + 2000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("standart/2.jpg")
			img2 := canvas.NewImageFromFile("standart/2 (1).jpg")
			img3 := canvas.NewImageFromFile("standart/2 (2).jpg")
			img4 := canvas.NewImageFromFile("standart/2 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		default:
			t := (rand.Intn(2000) + 4000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("standart/3.jpg")
			img2 := canvas.NewImageFromFile("standart/3 (1).jpg")
			img3 := canvas.NewImageFromFile("standart/3 (2).jpg")
			img4 := canvas.NewImageFromFile("standart/3 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		}
	case 5:
		switch room.Size {
		case 1:
			t := (rand.Intn(2000) + 1500) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vega/1.jpg")
			img2 := canvas.NewImageFromFile("vega/1 (1).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, btn1, btn2))) //установка виджетов
		case 2:
			t := (rand.Intn(2000) + 2000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vega/2.jpg")
			img2 := canvas.NewImageFromFile("vega/2 (1).jpg")
			img3 := canvas.NewImageFromFile("vega/2 (2).jpg")
			img4 := canvas.NewImageFromFile("vega/2 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		case 3:
			t := (rand.Intn(2500) + 3000) * z
			v := strconv.Itoa(t) + " рублей"
			price := widget.NewLabel(v)
			img1 := canvas.NewImageFromFile("vega/3.jpg")
			img2 := canvas.NewImageFromFile("vega/3 (1).jpg")
			img3 := canvas.NewImageFromFile("vega/3 (2).jpg")
			img4 := canvas.NewImageFromFile("vega/3 (3).jpg")
			btn1 := widget.NewButton("Сделать запрос на бронирование", func() { //обработка кнопки
				dbase.BookingData(ctx, hotel, room, u, t, roomset, app, price.Text, data) //выполнение функции BookingData
			})
			btn2 := widget.NewButton("Отказаться", func() { //обработка кнопки
				dbase.Authentication(ctx, app) //выход в окно авторизации
			})
			roomset.SetContent((container.NewGridWithColumns(2, box, price, img1, img2, img3, img4, btn1, btn2))) //установка виджетов
		}
	}
	return roomset
}

// окно бронирования
func (dbase *DB) BookingData(ctx context.Context, hotel models.Hotel, room models.Number, u models.User, t int, ch fyne.Window, app fyne.App, price, data string) {
	mngs, err := dbase.ManagerData(ctx) //получаем массив менеджеров из БД
	if err != nil {                     //если есть ошибка, выводим её и выходим из программы
		panic(err)
	}
	BR := make([]models.BookingRequest, len(mngs))
	bookingcheck := app.NewWindow("Бронирование") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/survey.png")
	bookingcheck.Resize(fyne.NewSize(300, 200))
	bookingcheck.SetIcon(icon)
	paymethod := widget.NewRadioGroup([]string{"Qiwi", "Картой онлайн"}, func(s string) {}) //список с выбором
	bookingcheck.SetContent(widget.NewVBox(                                                 //установка виджетов
		widget.NewLabel("Проверьте данные бронирования:"),
		widget.NewLabel("На кого происходит бронирование:"),
		widget.NewLabel(u.UserName), //имя пользователя
		widget.NewLabel("Цена:"),
		widget.NewLabel(price), //цена
		widget.NewLabel("Количество гостей:"),
		widget.NewLabel(strconv.Itoa(room.Size)), //кол-во человек
		widget.NewLabel("Количество ночей:"),
		widget.NewLabel(data), //кол-во ночей
		widget.NewLabel("Выберите способ оплаты:"),
		paymethod, //выбор способа оплаты
		widget.NewButton("Далее", func() {
			method := paymethod.Selected   //выбраннный способ оплаты
			count, _ := strconv.Atoi(data) //конвертация из string в int
			room.HotelName = hotel.HotelTitle
			for i := 0; i < len(BR); i++ {
				BR[i] = models.BookingRequest{ //авт. опр. BookingRequest, присваиваем значение полям BR
					IdFrom:    u.IdUser,
					IdTo:      mngs[i].UserStruct.IdUser,
					IdHotel:   room.IdHotel,
					IdNumber:  room.IdNumber,
					HotelName: room.HotelName,
					Size:      room.Size,
					Status:    room.Status,
					Price:     t,
					PayMethod: method,
					Nights:    count,
				}
				if err := dbase.AddBookingRequest(ctx, BR[i]); err != nil { //добавление запроса в БД, првоерка на ошибку
					panic(err)
				}
			}
			dialog.ShowInformation("Запрос отправлен!", "ОК", bookingcheck) //диалоговое окно
			time.Sleep(time.Second * 5)
			dbase.AuthenticationManager(ctx, mngs, app, bookingcheck) //авторизация менеджера
		}),
	))
	bookingcheck.Show()
	ch.Close()
}

func (dbase *DB) AuthenticationManager(ctx context.Context, mngs []models.Manager, app fyne.App, bookingcheck fyne.Window) { //окно авторизации менеджера
	br, err := dbase.ScanBooking(ctx) //получение адреса памяти на структуру BookingRequest
	BR := *br                         //обращение с помощью указателя
	if err != nil {                   //проверка на ошибку
		panic(err)
	}
	login := widget.NewEntry()
	password := widget.NewPasswordEntry()
	entermng := app.NewWindow("Аутентификация менеджера") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/mng.png")
	entermng.SetIcon(icon)
	entermng.Resize(fyne.NewSize(300, 200))
	entermng.SetContent(widget.NewVBox(
		widget.NewLabel("Логин:"),
		login,
		widget.NewLabel("Пароль:"),
		password,
		widget.NewButton("Войти", func() { //обработка кнопки
			for i := 0; i < len(mngs); i++ {
				if login.Text == mngs[i].UserStruct.UserLogin && password.Text == mngs[i].UserStruct.UserPswrd { //если логин и пароль совпадают
					mng := mngs[i]
					dbase.BookingDecision(ctx, mng, BR, app) //идем дальше
					entermng.Close()
					bookingcheck.Close()
				}
			}
		}),
	))
	entermng.Show()
	bookingcheck.Close()
}

func (dbase *DB) BookingDecision(ctx context.Context, mng models.Manager, br models.BookingRequest, app fyne.App) { //окно для обработки менеджером запроса на бронирование
	var uname string
	users, err := dbase.UsersData(ctx) //получаем данные структур пользоваталей
	if err != nil {
		panic(err) //проверка на ошибку
	}
	var u models.User
	for _, v := range users {
		if v.IdUser == br.IdFrom { //если id v пользователя совпадает с id пользователя, от которого поступил запрос, то опр. его имя
			uname = v.UserName
			u = v
		}
	}
	name := fmt.Sprintf("Уважаемая %s,", mng.UserStruct.UserName)
	s := fmt.Sprintf("Вам поступил запрос на бронирование %d номера в %s", br.IdNumber, br.HotelName)
	card := widget.NewCard(name, s, widget.NewLabel("Проверьте данные бронирования"))
	p, n, c := strconv.Itoa(br.Price), strconv.Itoa(br.Size), strconv.Itoa(br.Nights) //конвертация из int в стринг
	bookingmng := app.NewWindow("Проверка запроса")                                   //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/survey.png")
	bookingmng.SetIcon(icon)
	bookingmng.Resize(fyne.NewSize(400, 200))
	bookingmng.SetContent(widget.NewVBox( //установка виджетов
		card,
		widget.NewLabel("Бронирование происходит на имя:"),
		widget.NewLabel(uname),
		widget.NewLabel("Цена:"),
		widget.NewLabel(p),
		widget.NewLabel("Количество гостей:"),
		widget.NewLabel(n),
		widget.NewLabel("Количество ночей:"),
		widget.NewLabel(c),
		widget.NewLabel("Статус номера: "),
		widget.NewLabel(br.Status),
		widget.NewLabel("Способ оплаты:"),
		widget.NewLabel(br.PayMethod),
		widget.NewButton("Одобрить бронирование", func() {
			if err := dbase.UpdateStatus(ctx, br); err != nil { //обновление статуса номера в случае одобрения бронирования
				panic(err)
			}
			pay := dbase.Pay(ctx, br, app) //окно оплаты
			pay.Show()
			bookingmng.Close()
		}),
		widget.NewButton("Отказать в бронировании", func() { //отказ (в случае, если номер уже забронирован)
			dbase.Sorry(ctx, app, u, bookingmng)
		}),
	))
	bookingmng.Show()
}

func (dbase *DB) Sorry(ctx context.Context, app fyne.App, u models.User, bookingmng fyne.Window) {
	s := "Уважаемый " + u.UserName
	s1 := "Вам отказано в бронировании"
	card := widget.NewCard(s, s1, widget.NewLabel("извините"))
	sorry := app.NewWindow("Отказано")
	icon, _ := fyne.LoadResourceFromPath("icons/err.png")
	sorry.SetIcon(icon)
	sorry.SetContent(widget.NewVBox(
		card,
	))
	sorry.Show()
	bookingmng.Close()
}

func (dbase *DB) Pay(ctx context.Context, br models.BookingRequest, app fyne.App) fyne.Window { //окно оплаты
	u, err := dbase.ScanClient(ctx, br.IdFrom) //получаем данные клиента
	user := *u                                 //
	if err != nil {                            //проверка на ошибку
		panic(err)
	}
	s := "Цена: " + strconv.Itoa(br.Price)
	name := fmt.Sprintf("Уважаемый %s,", user.UserName)
	s1 := fmt.Sprintf("Вам одобрили бронирвание. Пожалуйста, оплатите номер!")
	card := widget.NewCard(name, s1, widget.NewLabel(s))
	pay := app.NewWindow("Оплата") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/pay.png")
	pay.SetIcon(icon)
	pay.Resize(fyne.NewSize(400, 200))

	set := widget.NewLabel("")
	switch br.PayMethod { //switch-case в зависимости от выбора способа оплаты
	case "Qiwi":
		pass := widget.NewPasswordEntry()
		pay.SetContent(widget.NewVBox( //установка виджетов
			card,
			widget.NewLabel("Введите номер телефона:"),
			widget.NewEntry(),
			widget.NewLabel("Введите пароль:"),
			pass,
			widget.NewButton("Оплатить", func() { //обработка кнопки
				if pass.Text == user.UserPswrd { //если логин и пароль клиента совпадают с данными из БД
					set.SetText("Производится оплата...")
					time.Sleep(time.Second * 2)
					dialog.ShowCustom("Оплата", "ОК", widget.NewLabel("Оплата успешно произведена"), pay) //показываем диалоговое окно
					time.Sleep(time.Second * 5)
					if err := dbase.UpdatePayStatus(ctx, br); err != nil {
						panic(err)
					}
					dbase.AuthenticationAdmin(ctx, app, pay) //переходим к авторизации даминистратора
				}
			}),
			set,
		))
	case "Картой онлайн":
		num := widget.NewEntry()
		name := widget.NewEntry()
		data := widget.NewEntry()
		cvv := widget.NewPasswordEntry()
		pay.SetContent(widget.NewVBox( //установка виджетов
			widget.NewLabel(s),
			widget.NewLabel("Введите данные карты:"),
			container.NewGridWithColumns(2, //установка виджетов по 2 колонкам
				widget.NewLabel("Введите номер карты:"), num,
				widget.NewLabel("Введите имя и фамилию:"), name,
				widget.NewLabel("Введите срок действия карты:"), data,
				widget.NewLabel("Введите cvv код:"), cvv),
			widget.NewButton("Оплатить", func() { //обработка кнопки
				if name.Text == user.UserName { //если имя и фамилия совпадают с данными клиента из БД
					set.SetText("Производится транзакция...")
					time.Sleep(time.Second * 2)
					dialog.ShowCustom("Оплата", "ОК", widget.NewLabel("Оплата успешно произведена"), pay) //показываем диалоговое окно
					time.Sleep(time.Second * 5)
					if err := dbase.UpdatePayStatus(ctx, br); err != nil {
						panic(err)
					}
					dbase.AuthenticationAdmin(ctx, app, pay) //переходим к авторизации даминистратора
				}
			}),
			set,
		))
	}
	return pay
}

func (dbase *DB) AuthenticationAdmin(ctx context.Context, app fyne.App, pay fyne.Window) { //окно авторизации администратора
	a, err := dbase.AdminData(ctx) //получаем данные администратора и заносим их в структуру
	admin := *a
	if err != nil { //проверка на ошибку
		panic(err)
	}
	login := widget.NewEntry()
	password := widget.NewPasswordEntry()
	enteradmin := app.NewWindow("Аутентификация администратора") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/admin.png")
	enteradmin.SetIcon(icon)
	enteradmin.Resize(fyne.NewSize(300, 200))
	enteradmin.SetContent(widget.NewVBox( //установка виджетов
		widget.NewLabel("Логин:"),
		login,
		widget.NewLabel("Пароль:"),
		password,
		widget.NewButton("Войти", func() { //обработка кнопки
			if login.Text == admin.UserStruct.UserLogin && password.Text == admin.UserStruct.UserPswrd { //если логин и пароль совпадают с данными из БД
				dbase.ChoiceWindow(ctx, app, enteradmin)
			}
		}),
	))
	enteradmin.Show()
	pay.Close()
}

func (dbase *DB) ChoiceWindow(ctx context.Context, app fyne.App, enteradmin fyne.Window) {
	makelist := app.NewWindow("Составление списка гостиниц и номеров") //параметры окна
	icon, _ := fyne.LoadResourceFromPath("icons/list.png")
	makelist.SetIcon(icon)
	makelist.Resize(fyne.NewSize(300, 200))
	makelist.SetContent(widget.NewVBox(
		widget.NewLabel("Выберите действие, которые хотите выполнить со списком"),
		widget.NewButton("Удалить отель", func() {
			dbase.DeleteList(ctx, app, makelist)
		}),
		widget.NewButton("Добавить отель", func() {
			dbase.NewHotel(ctx, app, makelist)
		}),
	)) //установка виджетов
	makelist.Show()
	enteradmin.Close()
}

func (dbase *DB) CreateCard(ctx context.Context, hotels []models.Hotel, s []string) *widget.Card { //функция создания виджета card
	var card *widget.Card
	switch len(hotels) {
	case 5:
		card = widget.NewCard(
			"ИНФОРМАЦИЯ", //заголовок
			"Ознакомьтесь с информацией, представленной в списке гостиниц:", //сообщение
			widget.NewAccordion( //формирование выпадающего списка
				widget.NewAccordionItem(hotels[0].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[0], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[0]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[0]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[1].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[1], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[1]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[1]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[2].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[2], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[2]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[2]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[3].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[3], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[3]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[3]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[4].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[4], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[4]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[4]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
			),
		)
	case 4:
		card = widget.NewCard(
			"ИНФОРМАЦИЯ", //заголовок
			"Ознакомьтесь с информацией, представленной в списке гостиниц:", //сообщение
			widget.NewAccordion( //формирование выпадающего списка
				widget.NewAccordionItem(hotels[0].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[0], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[0]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[0]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[1].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[1], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[1]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[1]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[2].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[2], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[2]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[2]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[3].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[3], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[3]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[3]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
			),
		)
	case 3:
		card = widget.NewCard(
			"ИНФОРМАЦИЯ", //заголовок
			"Ознакомьтесь с информацией, представленной в списке гостиниц:", //сообщение
			widget.NewAccordion( //формирование выпадающего списка
				widget.NewAccordionItem(hotels[0].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[0], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[0]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[0]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[1].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[1], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[1]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[1]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[2].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[2], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[2]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[2]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
			),
		)
	case 2:
		card = widget.NewCard(
			"ИНФОРМАЦИЯ", //заголовок
			"Ознакомьтесь с информацией, представленной в списке гостиниц:", //сообщение
			widget.NewAccordion( //формирование выпадающего списка
				widget.NewAccordionItem(hotels[0].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[0], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[0]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[0]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
				widget.NewAccordionItem(hotels[1].HotelTitle, //название отеля
					widget.NewAccordion(
						widget.NewAccordionItem(
							s[1], //данные о кол-ве свободных номеров
							widget.NewButton("Выбрать", func() { //обработка кнопки
								if err := dbase.DeleteHotel(ctx, hotels[1]); err != nil { //если администратор выбрал отель, то удаляем его из БД
									panic(err) //обработка ошибки
								}
								if err := dbase.DeleteRooms(ctx, hotels[1]); err != nil {
									panic(err)
								}
							}),
						),
					),
				),
			),
		)
	}
	return card
}

func (dbase *DB) DeleteList(ctx context.Context, app fyne.App, makelist fyne.Window) { //окно формирования списка гостиниц администратором
	hotels, err := dbase.HotelsData(ctx) //получаем массив структур Hotel
	if err != nil {
		panic(err) //обработка и зверешение программы в случае ошибки
	}
	s := make([]string, len(hotels))       //массив для занесения в него данных о кол-ве свободных номеров
	data := make(map[int]int, len(hotels)) //хэш-таблицы для каждого отеля, где ключ - id отеля, а значение - кол-во свободных номеров
	for i := 0; i < len(hotels); i++ {     //цикл по всем отелям
		x, err := dbase.CheckRoom(ctx) //проверка статуса номеров i отеля
		if err != nil {
			panic(err) //завершение программы и вывод ошибки в случае её возникновения
		}
		data[i] = x[i+1]                                                //заносим в data данные хэш-таблицы x
		s[i] = "Количество свободных номеров: " + strconv.Itoa(data[i]) //конвертация в string и занесения данных в массив s
	}
	card := dbase.CreateCard(ctx, hotels, s) //создание выпадающего списка
	deletelist := app.NewWindow("Составление списка гостиниц ")
	icon, _ := fyne.LoadResourceFromPath("icons/survey.png")
	deletelist.SetIcon(icon)
	deletelist.Resize(fyne.NewSize(300, 400))
	deletelist.SetContent(widget.NewVBox(
		widget.NewLabel("Выберите гостиницу, которую хотите удалить:"),
		card,
	))
	deletelist.Show()
	makelist.Close()
}

func (dbase *DB) NewHotel(ctx context.Context, app fyne.App, makelist fyne.Window) {
	hotels, err := dbase.HotelsData(ctx)
	if err != nil {
		panic(err)
	}
	id := len(hotels) + 1
	var hotel models.Hotel
	addlist := app.NewWindow("Добавление нового отеля")
	icon, _ := fyne.LoadResourceFromPath("icons/add.png")
	addlist.SetIcon(icon)
	addlist.Resize(fyne.NewSize(300, 300))
	title := "Hotel Vega"
	var img11, img12, img21, img22, img23, img24, img31, img32, img33, img34 *canvas.Image
	for i := 1; i < 4; i++ {
		switch i {
		case 1:
			img11 = canvas.NewImageFromFile("vega/1.jpg")
			img12 = canvas.NewImageFromFile("vega/1 (1).jpg")
		case 2:
			img21 = canvas.NewImageFromFile("vega/2.jpg")
			img22 = canvas.NewImageFromFile("vega/2 (1).jpg")
			img23 = canvas.NewImageFromFile("vega/2 (2).jpg")
			img24 = canvas.NewImageFromFile("vega/2 (3).jpg")
		case 3:
			img31 = canvas.NewImageFromFile("vega/3.jpg")
			img32 = canvas.NewImageFromFile("vega/3 (1).jpg")
			img33 = canvas.NewImageFromFile("vega/3 (2).jpg")
			img34 = canvas.NewImageFromFile("vega/3 (3).jpg")
		}
	}
	card := widget.NewCard("Информация об отеле:", title, widget.NewLabel("Ознакомьтесь с номерами отеля"))
	box1 := container.NewGridWithColumns(2, img11, img12)
	box2 := container.NewGridWithColumns(2, img21, img22, img23, img24)
	box3 := container.NewGridWithColumns(2, img31, img32, img33, img34)
	btn1 := CreateButton(box1, app)
	btn2 := CreateButton(box2, app)
	btn3 := CreateButton(box3, app)
	addlist.SetContent(widget.NewVBox(
		card,
		container.NewGridWithColumns(3, widget.NewLabel("Посмотреть одноместный вариант?"),
			widget.NewLabel("Посмотреть двухместный вариант?"), widget.NewLabel("Посмотреть трёхместный вариант?"),
			btn1, btn2, btn3),
		widget.NewButton("Добавить", func() {
			hotel = models.Hotel{
				IdHotel:    id,
				HotelTitle: title,
			}
			if err := dbase.AddNewHotel(ctx, hotel); err != nil {
				panic(err)
			}
			id := 13
			for i := 1; i < 4; i++ {
				if err := dbase.AddRoomsOfNewHotel(ctx, hotel, i, id); err != nil {
					panic(err)
				}
				id++
			}
		}),
	))
	addlist.Show()
	makelist.Close()
}

func CreateButton(box *fyne.Container, app fyne.App) *widget.Button {
	btn := widget.NewButton("Да", func() { ShowRoom(box, app) })
	return btn
}

func ShowRoom(box *fyne.Container, app fyne.App) {
	info := app.NewWindow("Информация")
	icon, _ := fyne.LoadResourceFromPath("icons/show.png")
	info.SetIcon(icon)
	info.Resize(fyne.NewSize(500, 300))
	info.SetContent(box)
	info.Show()
}

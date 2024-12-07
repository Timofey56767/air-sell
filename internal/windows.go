package internal

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"github.com/TOMMy-Net/air-sell/db"
	"github.com/TOMMy-Net/air-sell/models"
	"github.com/TOMMy-Net/air-sell/tools"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/widget"
)

type Settings struct {
	Window  fyne.Window
	Storage *db.Storage
	Account *models.User
	App     fyne.App
}

type RetFunc struct {
	F    func()
	Next *RetFunc
}

// Функция создания структуры настроек
func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) SignInWindow() {
	email := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("Почта", email),
		widget.NewFormItem("Пароль", password),
	)

	form.SubmitText = "Войти"
	form.Orientation = widget.Vertical

	form.OnSubmit = func() {
		m := models.UserEntry{
			Email:    email.Text,
			Password: password.Text,
		}
		if err := tools.Validate(m); err != nil {
			dialog.ShowError(ErrValidRow, s.Window)
		} else {
			res, err := s.Storage.GetUser(models.User{Email: m.Email, Password: m.Password})
			if err != nil {
				dialog.ShowInformation("Ошибка", "Ошибка базы данных", s.Window)
			} else {
				s.Account = &models.User{
					ID:       res.ID,
					Email:    m.Email,
					Password: m.Password,
				}
				s.MainWindow(models.TicketsSearch{})
			}
		}
	}

	label := widget.NewLabel("для того чтобы продолжить войдите в систему")

	fTxt := canvas.NewText("Добро Пожаловать", color.White)
	fTxt.TextSize = 20
	fTxt.TextStyle.Bold = true
	fTxt.Alignment = fyne.TextAlignCenter

	regButton := widget.NewButtonWithIcon("Зарегестрироваться", theme.AccountIcon(), func() {
		s.SignUpWindow()
	})

	/*withOutReg := widget.NewButton("Продолжить без регистрации", func() {
		s.MainWindow()
	})*/
	regButton.Enable()

	s.Window.SetContent(
		container.NewBorder(nil, nil, nil, nil,
			container.NewCenter(container.NewVBox(fTxt, label, canvas.NewLine(color.White), form, widget.NewActivity(), regButton))))

}

func (s *Settings) SignUpWindow() {
	email := widget.NewEntry()
	password := widget.NewPasswordEntry()
	confirmPassword := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("Почта", email),
		widget.NewFormItem("Пароль", password),
		widget.NewFormItem("Повторите пароль", confirmPassword),
	)

	form.SubmitText = "Зарегистрироваться"
	form.Orientation = widget.Vertical

	form.OnSubmit = func() {
		if password.Text == confirmPassword.Text {
			if len([]rune(password.Text)) > 6 {
				m := models.UserEntry{
					Email:    email.Text,
					Password: password.Text,
				}
				if err := tools.Validate(m); err != nil {
					dialog.ShowInformation("Ошибка", "Не все поля заполнены", s.Window)
				} else {
					if res, err := s.Storage.CreateUser(models.User{Email: m.Email, Password: m.Password}); err != nil {
						dialog.ShowInformation("Ошибка", "Такой аккаунт уже существует", s.Window)
					} else {
						s.Account = &models.User{
							ID:       int(res),
							Email:    m.Email,
							Password: m.Password,
						}
						s.MainWindow(models.TicketsSearch{})
					}
				}
			} else {
				dialog.ShowInformation("Ошибка", "Пароль меньше 6 символов", s.Window)
			}
		} else {
			dialog.ShowInformation("Ошибка", "Пароли не совпадают", s.Window)
		}
	}

	fTxt := canvas.NewText("Регистрация", color.White)
	fTxt.TextSize = 20
	fTxt.TextStyle.Bold = true
	fTxt.Alignment = fyne.TextAlignCenter

	label := widget.NewLabel("для того чтобы продолжить зарегестрируйтесь в системе")

	signButton := widget.NewButtonWithIcon("Войти", theme.AccountIcon(), func() {
		s.SignInWindow()
	})
	signButton.Enable()

	s.Window.SetContent(
		container.NewBorder(nil, nil, nil, nil,
			container.NewCenter(container.NewVBox(fTxt, label, canvas.NewLine(color.White), form, widget.NewActivity(), signButton))))
}

func (s *Settings) MainWindow(t models.TicketsSearch) {
	var fromEntry = widget.NewEntry()
	fromEntry.SetPlaceHolder("Откуда")
	if t.From != "" {
		fromEntry.Text = t.From
	}

	var toEntry = widget.NewEntry()
	toEntry.SetPlaceHolder("Куда")
	if t.To != "" {
		toEntry.Text = t.To
	}

	var dateButtonFrom = s.DateButton()
	dateButtonFrom.Text = "Туда (дата)"
	if t.Date_from != "" {
		dateButtonFrom.Text = t.Date_from
	}

	//var dateButtonTo = s.DateButton()
	//dateButtonTo.Text = "Обратно (дата)"

	grid := container.NewGridWithColumns(3, fromEntry, toEntry, dateButtonFrom)

	var buttonFind = widget.NewButtonWithIcon("Поиск", theme.SearchIcon(), func() {})

	var searchMenu = container.NewVBox(grid, widget.NewActivity(), buttonFind)
	buttonFind.OnTapped = func() {
		stack := container.NewVBox()
		ticket := models.NewTicketSearch() // инициализация структуры билета
		ticket.From = fromEntry.Text
		ticket.To = toEntry.Text
		ticket.Date_from = dateButtonFrom.Text
		//ticket.Date_to = dateButtonTo.Text

		tickets, err := s.FindTickets(ticket)
		if err != nil {
			dialog.ShowInformation("Ошибка", err.Error(), s.Window)
		} else {
			fmt.Println(tickets, err)
			if len(tickets) > 0 {
				for i := 0; i < len(tickets); i++ {
					t := tickets[i]
					stack.Add(widget.NewButton(fmt.Sprintf("%s (%s) \u27F6 %s (%s) \n %s \u27F6 %s \n Цена: %.2f \u20BD", t.DepartureFrom.City, t.DepartureFrom.Iata, t.ArrivalAt.City, t.ArrivalAt.Iata, t.DepartureTime, t.ArrivalTime, t.Price), func() {
						s.TicketWindow(&t)
					}))
				}
			} else {
				stack.Add(widget.NewLabel("По вашему запросу билеты не найдены"))
			}
			searchMenu.Refresh()
			s.Window.SetContent(container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, nil, s.ProfileButton()), container.NewCenter(searchMenu)), nil, nil, nil, widget.NewCard("", "", container.NewVScroll(stack))))
		}
	}

	//grid2 := container.NewGridWrap(fyne.NewSize(50, 100), fromEntry, toEntry, dateButtonFrom, dateButtonTo)
	s.Window.SetContent(container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, nil, s.ProfileButton())), nil, nil, nil, container.NewCenter(searchMenu)))
}

// окно билета
func (s *Settings) TicketWindow(t *models.Ticket) {

	var ticketInfo = widget.NewLabel(fmt.Sprintf("Авиакомпания %s", t.Airline))
	var wayInfo = widget.NewLabel(fmt.Sprintf("%s(%s) \u27F6 %s(%s)", t.DepartureFrom.City, t.DepartureFrom.Iata, t.ArrivalAt.City, t.ArrivalAt.Iata))
	var timeInfo = widget.NewLabel(fmt.Sprintf("%s \u27F6 %s", t.DepartureTime, t.ArrivalTime))
	var count = widget.NewLabel(fmt.Sprintf("Колличество: %d", t.Quantity))
	var baggInfo = widget.NewLabel(fmt.Sprintf("Багаж: %s  Ручная кладь: %s", t.Luggage, t.HandBaggage))
	var priceInfo = widget.NewLabel(fmt.Sprintf("Цена: %.2f \u20BD", t.Price))

	var buyButton = widget.NewButton("Купить", func() {
		s.BuyWindow(t)
	})

	var retButton = widget.NewButton("Назад", func() {

		ticket := returnTicketEntry(t)
		s.MainWindow(ticket)
	})

	s.Window.SetContent(container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				retButton,
				s.ProfileButton())),
		nil,
		nil,
		nil,
		container.NewCenter(
			container.NewVBox(
				ticketInfo,
				wayInfo,
				timeInfo,
				count,
				baggInfo,
				priceInfo,
				buyButton,
			),
		)))
}

func (s *Settings) BuyWindow(t *models.Ticket) {
	var ticketCols = 0
	var ticketColsText = ""
	var cols = binding.BindString(&ticketColsText)

	var selectPassports = &widget.CheckGroup{}
	var selectedPassports []string

	var buttonBack = widget.NewButton("Назад", func() {
		s.TicketWindow(t)
	})

	passports, err := s.Storage.GetPassportData(s.Account.ID)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Не возможно получить данные аккаунта", s.Window)
		return
	}

	var buttonBuy = widget.NewButton("Оплатить", func() {
		if ticketCols > 0 && t.Quantity >= ticketCols {
			dialog.ShowConfirm("Подтверждение оплаты", fmt.Sprintf("Колличество билетов: %d\nЦена: %.2f \u20BD", ticketCols, t.Price*float64(ticketCols)), func(b bool) {
				if b {
					// После подтверждения происходит оплата и обновление данных в БД
					// Затем происходит отправка билетов пользователю
					err := s.Storage.MinusTicketCount(t.ID, ticketCols)
					if err != nil {
						dialog.ShowInformation("Ошибка", "Не возможно обновить данные билетов", s.Window)
						return
					}
					for _, v := range selectedPassports {
						// Создание процесса создания билетов
						fmt.Println(v)
					}

					err = s.Storage.SetBuyHistory(models.BuyHistory{
						Ticket:  *t,
						UserId:  s.Account.ID,
						BuyTime: time.Now().Format(time.RFC3339),
						Count:   ticketCols,
					})
					if err != nil {
						dialog.ShowInformation("Ошибка", "Не возможно обновить данные истории покупок", s.Window)
						return
					}
					ticket := returnTicketEntry(t)

					s.MainWindow(ticket)
					dialog.ShowInformation("Успешно", "Билет успешно куплен", s.Window)

				}
			}, s.Window)
		} else {
			dialog.ShowInformation("Ошибка", "Не возможно купить билет, выберите паспорта или проверьте доступность билета", s.Window)
		}
	})

	var addPassport = widget.NewButton("Добавить паспорт", func() {
		pass := &models.Passport{}
		form := s.PassportEntry(pass)
		s.Window.SetContent(
			container.NewBorder(
				container.NewBorder(
					nil,
					nil,
					widget.NewButton("Назад", func() {
						s.BuyWindow(t)
					}),
					nil,
				),
				nil,
				nil,
				nil,
				container.NewCenter(
					form,
				),
			),
		)
	})

	if len(passports) > 0 {
		convPass := convertPassports(passports)
		selectPassports = widget.NewCheckGroup(
			convPass, func(s []string) {
				ticketCols = len(s)
				ticketColsText = fmt.Sprintf("Колличество выбранных пассажиров: %d", ticketCols)
				selectedPassports = s
				cols.Reload()
			},
		)
	}

	s.Window.SetContent(
		container.NewBorder(
			container.NewBorder(
				nil,
				nil,
				buttonBack,
				nil,
			),
			nil,
			nil,
			nil,
			container.NewCenter(
				container.NewVBox(
					widget.NewLabelWithData(cols),
					selectPassports,
					addPassport,
					buttonBuy,
				),
			),
		),
	)
}

// Окно профиля
func (s *Settings) ProfileWindow() {
	var userInfo = widget.NewLabel(fmt.Sprintf("ID: %d\nEmail: %s\nПароль: %s", s.Account.ID, s.Account.Email, s.Account.Password))
	var exitButton = widget.NewButtonWithIcon("Выйти", theme.CancelIcon(), func() {
		s.Account = &models.User{}
		s.SignInWindow()
	})

	var buyHistory = widget.NewButton("История покупок", func() {
		s.HistoryWindow()
	})

	var cardInfo = widget.NewCard("", "Информация аккаунта", userInfo)
	var buttonBack = widget.NewButton("Назад", func() {
		s.MainWindow(models.TicketsSearch{})
	})

	var passportsButton = widget.NewButton("Паспорта", func() {
		s.PassportWindow()
	})

	var addPasportbutton = widget.NewButton("Добавить паспорт", func() {
		passport := &models.Passport{}
		form := s.PassportEntry(passport)
		s.Window.SetContent(
			container.NewBorder(
				container.NewBorder(
					nil,
					nil,
					widget.NewButton("Назад", func() {
						s.ProfileWindow()
					}),
					nil,
				),
				nil,
				nil,
				nil,
				container.NewCenter(
					form,
				),
			),
		)
	})

	//	var bar = widget
	s.Window.SetContent(container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			buttonBack,
			nil),
		nil,
		nil,
		nil,
		container.NewCenter(
			container.NewVBox(
				cardInfo,
				buyHistory,
				passportsButton,
				addPasportbutton,
				exitButton,
			),
		)))
}

func (s *Settings) PassportWindow() {
	var buttonBack = widget.NewButton("Назад", func() {
		s.ProfileWindow()
	})
	var stack = container.NewGridWithColumns(3)
	var passports, err = s.Storage.GetPassportData(s.Account.ID)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Ошибка получения поспартов", s.Window)
		return
	}

	for _, v := range passports {
		passport := v
		stack.Add(widget.NewCard("", fmt.Sprintf("%s %s (%s)", v.Name, v.Surname, v.SeriesAndNumber), container.NewHBox(
			widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
				err := s.Storage.DeletePassport(passport.ID)
				if err != nil {
					dialog.ShowInformation("Ошибка", "Ошибка удаления паспорта", s.Window)
					return
				}
				dialog.ShowInformation("Успех", "Паспорт удален", s.Window)
			}),
		)))
	}
	s.Window.SetContent(container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			buttonBack,
			nil),
		nil,
		nil,
		nil,

		container.NewVScroll(
			container.NewVBox(
				stack,
			),
		),
	))
}

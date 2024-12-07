package internal

import (
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/TOMMy-Net/air-sell/models"
	"github.com/TOMMy-Net/air-sell/tools"
	datepicker "github.com/sdassow/fyne-datepicker"
)

var (
	ErrValidRow = errors.New("Не все поля были заполнены")
)

// Поле ввода даты
func (s *Settings) DateButton() *widget.Button {
	dateInput := widget.NewButton("", func() {})

	dateInput.OnTapped = func() {
		var d *dialog.CustomDialog

		when, err := time.Parse("2006-01-02", dateInput.Text)
		if err != nil {
			when = time.Now()
		}

		datepicker := datepicker.NewDatePicker(when, time.Monday, func(when time.Time, ok bool) {
			if ok {
				dateInput.SetText(when.Format("2006-01-02"))
			}
			d.Hide()
		})

		d = dialog.NewCustomWithoutButtons("Выберите дату", datepicker, s.Window)
		d.Show()
	}
	return dateInput

}

// Поиск билетов
func (s *Settings) FindTickets(ticket *models.TicketsSearch) ([]models.Ticket, error) {

	err := tools.Validate(ticket)
	if err != nil {
		return []models.Ticket{}, ErrValidRow
	}
	fmt.Println(ticket)
	if validateDateFormat(ticket.Date_from) {
		tickets, err := s.Storage.FindTickets(ticket)
		if err != nil {
			return nil, err
		}
		return tickets, nil
	} else {
		return []models.Ticket{}, ErrValidRow
	}
}

// кнопка профиля
func (s *Settings) ProfileButton() *widget.Button {
	button := widget.NewButtonWithIcon("Профиль", theme.AccountIcon(), func() {
		s.ProfileWindow()
	})
	return button
}

func validateDateFormat(d string) bool {
	_, err := time.Parse("2006-01-02", d)
	if err != nil {
		return false
	}
	return true
}

func convertPassports(p []models.Passport) []string {
	res := []string{}
	for _, v := range p {
		res = append(res, fmt.Sprintf("%d: %s %s %s", v.ID, v.Name, v.Surname, v.SeriesAndNumber))
	}
	return res
}

func (s *Settings) PassportEntry(passport *models.Passport) *widget.Form {

	var name = widget.NewEntry()
	var surname = widget.NewEntry()
	var patronymic = widget.NewEntry()
	var seriesAndNumber = widget.NewEntry()
	var gender = widget.NewSelect([]string{"m", "w"}, func(s string) {
		passport.Gender = s
	})
	var validityPeriod = s.DateButton()
	var birthDate = s.DateButton()
	var passportType = widget.NewSelect([]string{"international_passport", "passport"}, func(s string) {
		passport.Type = s
	})
	var citizenship = widget.NewEntry()

	var form = widget.NewForm(
		widget.NewFormItem("Имя", name),
		widget.NewFormItem("Фамилия", surname),
		widget.NewFormItem("Отчество", patronymic),
		widget.NewFormItem("Серия и номер", seriesAndNumber),
		widget.NewFormItem("Пол", gender),
		widget.NewFormItem("Срок действия", validityPeriod),
		widget.NewFormItem("Дата рождения", birthDate),
		widget.NewFormItem("Тип документа", passportType),
		widget.NewFormItem("Гражданство", citizenship),
	)
	form.SubmitText = "Добавить"
	form.OnSubmit = func() {
		VP, err := time.Parse(time.RFC3339, validityPeriod.Text)
		if err != nil {
			VP = time.Time{}
			passport.ValidityPeriod.Valid = false

		}
		BD := validateDateFormat(birthDate.Text)
		if !BD {
			dialog.ShowInformation("Ошибка", "Дата введена неверно", s.Window)
			return
		}

		passport.Name = name.Text
		passport.Surname = surname.Text
		passport.Patronymic = patronymic.Text
		passport.SeriesAndNumber = seriesAndNumber.Text
		passport.ValidityPeriod.Time = VP
		passport.Birthday = birthDate.Text
		passport.Citizenship = citizenship.Text
		passport.UserID = s.Account.ID

		err = tools.Validate(passport)
		if err != nil {
			dialog.ShowInformation("Ошибка", "Ошибка валидации", s.Window)
			return
		}

		_, err = s.Storage.AddPassport(*passport)
		if err != nil {
			dialog.ShowInformation("Ошибка", "Ошибка записи", s.Window)
			return
		}
		dialog.ShowInformation("Успех", "Паспорт успешно добавлен", s.Window)
	}

	return form
}

func returnTicketEntry(t *models.Ticket) models.TicketsSearch {

	ticket := models.TicketsSearch{
		From: t.DepartureFrom.City,
		To:   t.ArrivalAt.City,
	}
	time, err := time.Parse(time.DateTime, t.DepartureTime)
	if err != nil {
		ticket.Date_from = ""
	} else {
		ticket.Date_from = time.Format("2006-01-02")
	}
	return ticket
}


func (s *Settings) HistoryWindow()  {
	retButton := widget.NewButton("Назад", func() {
		s.ProfileWindow()
	})

	history, err := s.Storage.GetBuyHistory(s.Account.ID)
	if err != nil {
		dialog.ShowInformation("Ошибка", "Ошибка получения истории", s.Window)
	}else{
		grid := container.NewGridWithColumns(2)
		for _, v := range history {
			grid.Add(widget.NewCard("", fmt.Sprintf("ID билета: %s", v.Ticket.ID), container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Дата покупки: %s", v.BuyTime)),
				widget.NewLabel(fmt.Sprintf("Цена: %.2f \u20BD", v.Ticket.Price)),
				widget.NewLabel(fmt.Sprintf("Колличество: %d", v.Count)),
			)))
		}
		s.Window.SetContent(
			container.NewBorder(
				container.NewBorder(
					nil, 
					nil,
					retButton,
					nil,
				),
				nil,
				nil, 
				nil,
				container.NewVScroll(
					grid,
				),
			),
		)
	}
}
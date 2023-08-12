package desktop

import (
	"context"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run(ctx context.Context, addr string) error {
	myApp := app.New()

	// controller := new(ControllerMock)
	controller, err := newController_gRPC(addr)
	if err != nil {
		return err
	}

	defer controller.Close()

	coreBox := container.NewVBox()
	authBox := container.NewAppTabs(
		container.NewTabItemWithIcon("Авторизация", theme.LoginIcon(), renderLogin(ctx, controller)),
		container.NewTabItemWithIcon("Регистрация", theme.AccountIcon(), renderRegister(ctx, controller)),
		container.NewTabItemWithIcon("Супер-кнопка", theme.HomeIcon(), renderButtonArea(ctx, controller)),
	)
	coreBox.Add(authBox)

	w := myApp.NewWindow("Desktop client")
	w.Resize(fyne.NewSize(800, 500))
	w.SetContent(coreBox)
	w.Show()

	myApp.Run()

	return nil
}

func renderLogin(ctx context.Context, c Controller) fyne.CanvasObject {
	login := widget.NewEntry()
	password := widget.NewPasswordEntry()

	alertBox, setError := makeAlert()

	loginCallback := func() {
		alertBox.Hide()

		err := c.Login(ctx, login.Text, password.Text)
		if err != nil {
			setError(err.Error())
			alertBox.Show()

			return
		}
	}

	alertBox.Hide()

	form := widget.NewForm(
		widget.NewFormItem("Логин", login),
		widget.NewFormItem("Пароль", password),
		widget.NewFormItem("", widget.NewButton("Войти", loginCallback)),
	)

	return container.NewVBox(alertBox, form)
}

func renderRegister(ctx context.Context, c Controller) fyne.CanvasObject {
	login := widget.NewEntry()
	password1 := widget.NewPasswordEntry()
	password2 := widget.NewPasswordEntry()

	alertBox, setError := makeAlert()

	registerCallback := func() {
		alertBox.Hide()

		if password1.Text != password2.Text {
			setError("Пароль не совпал")
			alertBox.Show()

			return
		}

		err := c.Register(ctx, login.Text, password1.Text)
		if err != nil {
			setError(err.Error())
			alertBox.Show()

			return
		}
	}

	alertBox.Hide()

	form := widget.NewForm(
		widget.NewFormItem("Логин", login),
		widget.NewFormItem("Пароль", password1),
		widget.NewFormItem("Повторите пароль", password2),
		widget.NewFormItem("", widget.NewButton("Зарегистрироваться", registerCallback)),
	)

	return container.NewVBox(alertBox, form)
}

func makeAlert() (fyne.CanvasObject, func(string)) {
	text := binding.NewString()

	box := container.NewHBox(
		widget.NewIcon(theme.WarningIcon()),
		widget.NewLabelWithData(text),
	)

	return container.NewCenter(box), func(s string) {
		_ = text.Set(s)
	}
}

func renderButtonArea(ctx context.Context, c Controller) fyne.CanvasObject {
	durationEntry := widget.NewEntry()
	chanceEntry := widget.NewEntry()

	alertBox, setError := makeAlert()

	registerCallback := func() {
		alertBox.Hide()

		if durationEntry.Text == "" {
			setError("Не указана продолжительность")
			alertBox.Show()

			return
		}

		if chanceEntry.Text == "" {
			setError("Не указана вероятность провала")
			alertBox.Show()

			return
		}

		duration, err := strconv.ParseInt(durationEntry.Text, 10, 64)
		if err != nil {
			setError(err.Error())
			alertBox.Show()

			return
		}

		chance, err := strconv.ParseInt(chanceEntry.Text, 10, 64)
		if err != nil {
			setError(err.Error())
			alertBox.Show()

			return
		}

		err = c.ButtonClick(ctx, duration, chance)
		if err != nil {
			setError(err.Error())
			alertBox.Show()

			return
		}
	}

	alertBox.Hide()

	form := widget.NewForm(
		widget.NewFormItem("Продолжительность", durationEntry),
		widget.NewFormItem("Вероятность провала", chanceEntry),
		widget.NewFormItem("", widget.NewButton("Отправить запрос", registerCallback)),
	)

	return container.NewVBox(alertBox, form)
}

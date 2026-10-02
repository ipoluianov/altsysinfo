package texts

import "fmt"

var en = Strings{
	Error: "Error",

	ToolRefresh:   "Refresh",
	ToolPdfReport: "Save PDF report",
	ToolTray:      "Minimize to tray",
	TrayShow:      "Show",
	TrayQuit:      "Quit",
	Settings:      "Settings",
	Help:          "Help",
	About:         "About",

	Categories: "Categories",
	Category: CategoryStrings{
		Common:          "Common",
		RAM:             "RAM",
		System:          "System",
		CPU:             "CPU",
		Storage:         "Storage",
		Graphics:        "Graphics & Displays",
		Network:         "Network",
		PCI:             "PCI Devices",
		USB:             "USB Devices",
		Audio:           "Audio",
		Sensors:         "Sensors",
		Battery:         "Battery",
		PCIVendorDevice: "PCI Vendor/Device",
	},

	ColName:          "Name",
	ColValue:         "Value",
	ColModel:         "Model",
	ColSize:          "Size",
	ColSpeed:         "Speed",
	ColVendorID:      "Vendor ID",
	ColDeviceID:      "Device ID",
	ColVendorName:    "Vendor Name",
	ColDeviceName:    "Device Name",
	VendorDeviceName: "Vendor/Device Name",
	CPUModel:         "CPU Model",
	CPUCores:         "CPU Cores",
	Drive:            func(n int) string { return fmt.Sprintf("Drive %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Driver",
	Copied:           "Copied",
	Loading:          "Loading…",
	NoData:           "No data available",
	CannotLoad:       func(err string) string { return "Cannot load information: " + err },

	SaveReportTitle:      "Save PDF report",
	PDFFiles:             "PDF files",
	CannotCollect:        func(err string) string { return "Cannot collect system information: " + err },
	CannotGenerate:       func(err string) string { return "Cannot generate report: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Cannot open save dialog: " + err },
	CannotSave:           func(err string) string { return "Cannot save report: " + err },
	SavedTo:              func(path string) string { return "Report saved to " + path },

	AlwaysOnTop:    "Always on top",
	Language:       "Language:",
	LanguageSystem: "As in the system",
	Theme:          "Theme:",
	ThemeDark:      "Dark",
	ThemeLight:     "Light",

	AboutTitle:   func(name string) string { return "About " + name },
	Version:      "Version",
	Author:       "Author:",
	License:      "License:",
	VisitWebsite: "Visit Website",
	Close:        "Close",

	Install: "Install",
	InstallAsk: func(dir string) string {
		return "AltSysInfo will be copied to " + dir + " and added to the Start menu, the desktop and the list of installed apps. If it is already installed there, it is updated. Settings stay as they are.\n\nAltSysInfo will then restart from there."
	},
	InstallFailed: func(err string) string { return "AltSysInfo could not be installed:\n\n" + err },
	Installed:     "AltSysInfo is installed",
	UninstallAsk: func(dir string) string {
		return "Remove AltSysInfo from this computer?\n\nSettings stay in " + dir
	},
	UninstallRunning: "AltSysInfo is running. Close it and try again.",
	Uninstalled:      "AltSysInfo has been removed.",
}

var ru = Strings{
	Error: "Ошибка",

	ToolRefresh:   "Обновить",
	ToolPdfReport: "Сохранить отчёт PDF",
	ToolTray:      "Свернуть в трей",
	TrayShow:      "Показать",
	TrayQuit:      "Выход",
	Settings:      "Настройки",
	Help:          "Справка",
	About:         "О программе",

	Categories: "Категории",
	Category: CategoryStrings{
		Common:          "Общее",
		RAM:             "Память",
		System:          "Система",
		CPU:             "Процессор",
		Storage:         "Накопители",
		Graphics:        "Графика и дисплеи",
		Network:         "Сеть",
		PCI:             "Устройства PCI",
		USB:             "Устройства USB",
		Audio:           "Звук",
		Sensors:         "Датчики",
		Battery:         "Батарея",
		PCIVendorDevice: "Справочник PCI",
	},

	ColName:          "Название",
	ColValue:         "Значение",
	ColModel:         "Модель",
	ColSize:          "Объём",
	ColSpeed:         "Частота",
	ColVendorID:      "ID производителя",
	ColDeviceID:      "ID устройства",
	ColVendorName:    "Производитель",
	ColDeviceName:    "Устройство",
	VendorDeviceName: "Производитель/устройство",
	CPUModel:         "Модель процессора",
	CPUCores:         "Ядра процессора",
	Drive:            func(n int) string { return fmt.Sprintf("Диск %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("Видеокарта %d", n) },
	Driver:           "Драйвер",
	Copied:           "Скопировано",
	Loading:          "Загрузка…",
	NoData:           "Нет данных",
	CannotLoad:       func(err string) string { return "Не удалось получить сведения: " + err },

	SaveReportTitle: "Сохранить отчёт PDF",
	PDFFiles:        "Файлы PDF",
	CannotCollect: func(err string) string {
		return "Не удалось собрать сведения о системе: " + err
	},
	CannotGenerate: func(err string) string { return "Не удалось создать отчёт: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "Не удалось открыть окно сохранения: " + err
	},
	CannotSave: func(err string) string { return "Не удалось сохранить отчёт: " + err },
	SavedTo:    func(path string) string { return "Отчёт сохранён: " + path },

	AlwaysOnTop:    "Поверх всех окон",
	Language:       "Язык:",
	LanguageSystem: "Как в системе",
	Theme:          "Тема:",
	ThemeDark:      "Тёмная",
	ThemeLight:     "Светлая",

	AboutTitle:   func(name string) string { return "О программе " + name },
	Version:      "Версия",
	Author:       "Автор:",
	License:      "Лицензия:",
	VisitWebsite: "Открыть сайт",
	Close:        "Закрыть",

	Install: "Установить",
	InstallAsk: func(dir string) string {
		return "AltSysInfo будет скопирован в " + dir + " и добавлен в меню «Пуск», на рабочий стол и в список установленных приложений. Если он там уже установлен, он будет обновлён. Настройки сохранятся.\n\nЗатем AltSysInfo перезапустится оттуда."
	},
	InstallFailed: func(err string) string { return "Не удалось установить AltSysInfo:\n\n" + err },
	Installed:     "AltSysInfo установлен",
	UninstallAsk: func(dir string) string {
		return "Удалить AltSysInfo с этого компьютера?\n\nНастройки останутся в " + dir
	},
	UninstallRunning: "AltSysInfo запущен. Закройте его и попробуйте снова.",
	Uninstalled:      "AltSysInfo удалён.",
}

package texts

import "fmt"

var pl = Strings{
	Error: "Błąd",

	ToolRefresh:   "Odśwież",
	ToolPdfReport: "Zapisz raport PDF",
	ToolTray:      "Zminimalizuj do zasobnika",
	TrayShow:      "Pokaż",
	TrayQuit:      "Zakończ",
	Settings:      "Ustawienia",
	Help:          "Pomoc",
	About:         "O programie",

	Categories: "Kategorie",
	Category: CategoryStrings{
		Common:          "Ogólne",
		RAM:             "Pamięć RAM",
		System:          "System",
		CPU:             "Procesor",
		Storage:         "Dyski",
		Graphics:        "Grafika i ekrany",
		Network:         "Sieć",
		PCI:             "Urządzenia PCI",
		USB:             "Urządzenia USB",
		Audio:           "Dźwięk",
		Sensors:         "Czujniki",
		Battery:         "Bateria",
		PCIVendorDevice: "Producenci/urządzenia PCI",
	},

	ColName:          "Nazwa",
	ColValue:         "Wartość",
	ColModel:         "Model",
	ColSize:          "Rozmiar",
	ColSpeed:         "Szybkość",
	ColVendorID:      "ID producenta",
	ColDeviceID:      "ID urządzenia",
	ColVendorName:    "Producent",
	ColDeviceName:    "Urządzenie",
	VendorDeviceName: "Producent/urządzenie",
	CPUModel:         "Model procesora",
	CPUCores:         "Rdzenie procesora",
	Drive:            func(n int) string { return fmt.Sprintf("Dysk %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Sterownik",
	Copied:           "Skopiowano",
	Loading:          "Wczytywanie…",
	NoData:           "Brak danych",
	CannotLoad:       func(err string) string { return "Nie można odczytać informacji: " + err },

	SaveReportTitle:      "Zapisz raport PDF",
	PDFFiles:             "Pliki PDF",
	CannotCollect:        func(err string) string { return "Nie można zebrać informacji o systemie: " + err },
	CannotGenerate:       func(err string) string { return "Nie można utworzyć raportu: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Nie można otworzyć okna zapisu: " + err },
	CannotSave:           func(err string) string { return "Nie można zapisać raportu: " + err },
	SavedTo:              func(path string) string { return "Raport zapisano: " + path },

	AlwaysOnTop:    "Zawsze na wierzchu",
	Language:       "Język:",
	LanguageSystem: "Jak w systemie",
	Theme:          "Motyw:",
	ThemeDark:      "Ciemny",
	ThemeLight:     "Jasny",

	AboutTitle:   func(name string) string { return "O programie " + name },
	Version:      "Wersja",
	Author:       "Autor:",
	License:      "Licencja:",
	VisitWebsite: "Odwiedź stronę",
	Close:        "Zamknij",

	Install: "Zainstaluj",
	InstallAsk: func(dir string) string {
		return "AltSysInfo zostanie skopiowany do " + dir + " i dodany do menu Start, na pulpit i do listy zainstalowanych aplikacji. Jeśli jest już tam zainstalowany, zostanie zaktualizowany. Ustawienia pozostaną bez zmian.\n\nNastępnie AltSysInfo uruchomi się ponownie stamtąd."
	},
	InstallFailed: func(err string) string { return "Nie udało się zainstalować AltSysInfo:\n\n" + err },
	Installed:     "AltSysInfo jest zainstalowany",
	UninstallAsk: func(dir string) string {
		return "Usunąć AltSysInfo z tego komputera?\n\nUstawienia pozostaną w " + dir
	},
	UninstallRunning: "AltSysInfo jest uruchomiony. Zamknij go i spróbuj ponownie.",
	Uninstalled:      "AltSysInfo został usunięty.",
}

var sr = Strings{
	Error: "Грешка",

	ToolRefresh:   "Освежи",
	ToolPdfReport: "Сачувај PDF извештај",
	ToolTray:      "Умањи у системску траку",
	TrayShow:      "Прикажи",
	TrayQuit:      "Изађи",
	Settings:      "Подешавања",
	Help:          "Помоћ",
	About:         "О програму",

	Categories: "Категорије",
	Category: CategoryStrings{
		Common:          "Опште",
		RAM:             "RAM меморија",
		System:          "Систем",
		CPU:             "Процесор",
		Storage:         "Складиштење",
		Graphics:        "Графика и екрани",
		Network:         "Мрежа",
		PCI:             "PCI уређаји",
		USB:             "USB уређаји",
		Audio:           "Звук",
		Sensors:         "Сензори",
		Battery:         "Батерија",
		PCIVendorDevice: "PCI произвођачи/уређаји",
	},

	ColName:          "Назив",
	ColValue:         "Вредност",
	ColModel:         "Модел",
	ColSize:          "Величина",
	ColSpeed:         "Брзина",
	ColVendorID:      "ID произвођача",
	ColDeviceID:      "ID уређаја",
	ColVendorName:    "Произвођач",
	ColDeviceName:    "Уређај",
	VendorDeviceName: "Произвођач/уређај",
	CPUModel:         "Модел процесора",
	CPUCores:         "Језгра процесора",
	Drive:            func(n int) string { return fmt.Sprintf("Диск %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Драјвер",
	Copied:           "Копирано",
	Loading:          "Учитавање…",
	NoData:           "Нема података",
	CannotLoad:       func(err string) string { return "Није могуће учитати податке: " + err },

	SaveReportTitle: "Сачувај PDF извештај",
	PDFFiles:        "PDF датотеке",
	CannotCollect: func(err string) string {
		return "Није могуће прикупити податке о систему: " + err
	},
	CannotGenerate: func(err string) string { return "Није могуће направити извештај: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "Није могуће отворити прозор за чување: " + err
	},
	CannotSave: func(err string) string { return "Није могуће сачувати извештај: " + err },
	SavedTo:    func(path string) string { return "Извештај је сачуван: " + path },

	AlwaysOnTop:    "Увек на врху",
	Language:       "Језик:",
	LanguageSystem: "Као у систему",
	Theme:          "Тема:",
	ThemeDark:      "Тамна",
	ThemeLight:     "Светла",

	AboutTitle:   func(name string) string { return "О програму " + name },
	Version:      "Верзија",
	Author:       "Аутор:",
	License:      "Лиценца:",
	VisitWebsite: "Посети сајт",
	Close:        "Затвори",

	Install: "Инсталирај",
	InstallAsk: func(dir string) string {
		return "AltSysInfo ће бити копиран у " + dir + " и додат у мени Старт, на радну површину и у листу инсталираних апликација. Ако је тамо већ инсталиран, биће ажуриран. Подешавања остају непромењена.\n\nAltSysInfo ће се затим поново покренути одатле."
	},
	InstallFailed: func(err string) string { return "Није могуће инсталирати AltSysInfo:\n\n" + err },
	Installed:     "AltSysInfo је инсталиран",
	UninstallAsk: func(dir string) string {
		return "Уклонити AltSysInfo са овог рачунара?\n\nПодешавања остају у " + dir
	},
	UninstallRunning: "AltSysInfo је покренут. Затворите га и покушајте поново.",
	Uninstalled:      "AltSysInfo је уклоњен.",
}

var de = Strings{
	Error: "Fehler",

	ToolRefresh:   "Aktualisieren",
	ToolPdfReport: "PDF-Bericht speichern",
	ToolTray:      "In den Infobereich minimieren",
	TrayShow:      "Anzeigen",
	TrayQuit:      "Beenden",
	Settings:      "Einstellungen",
	Help:          "Hilfe",
	About:         "Über",

	Categories: "Kategorien",
	Category: CategoryStrings{
		Common:          "Allgemein",
		RAM:             "Arbeitsspeicher",
		System:          "System",
		CPU:             "Prozessor",
		Storage:         "Laufwerke",
		Graphics:        "Grafik & Displays",
		Network:         "Netzwerk",
		PCI:             "PCI-Geräte",
		USB:             "USB-Geräte",
		Audio:           "Audio",
		Sensors:         "Sensoren",
		Battery:         "Akku",
		PCIVendorDevice: "PCI-Hersteller/Geräte",
	},

	ColName:          "Name",
	ColValue:         "Wert",
	ColModel:         "Modell",
	ColSize:          "Größe",
	ColSpeed:         "Takt",
	ColVendorID:      "Hersteller-ID",
	ColDeviceID:      "Geräte-ID",
	ColVendorName:    "Hersteller",
	ColDeviceName:    "Gerät",
	VendorDeviceName: "Hersteller/Gerät",
	CPUModel:         "Prozessormodell",
	CPUCores:         "Prozessorkerne",
	Drive:            func(n int) string { return fmt.Sprintf("Laufwerk %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Treiber",
	Copied:           "Kopiert",
	Loading:          "Wird geladen…",
	NoData:           "Keine Daten verfügbar",
	CannotLoad:       func(err string) string { return "Informationen können nicht geladen werden: " + err },

	SaveReportTitle:      "PDF-Bericht speichern",
	PDFFiles:             "PDF-Dateien",
	CannotCollect:        func(err string) string { return "Systeminformationen können nicht erfasst werden: " + err },
	CannotGenerate:       func(err string) string { return "Bericht kann nicht erstellt werden: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Speicherdialog kann nicht geöffnet werden: " + err },
	CannotSave:           func(err string) string { return "Bericht kann nicht gespeichert werden: " + err },
	SavedTo:              func(path string) string { return "Bericht gespeichert unter " + path },

	AlwaysOnTop:    "Immer im Vordergrund",
	Language:       "Sprache:",
	LanguageSystem: "Wie im System",
	Theme:          "Design:",
	ThemeDark:      "Dunkel",
	ThemeLight:     "Hell",

	AboutTitle:   func(name string) string { return "Über " + name },
	Version:      "Version",
	Author:       "Autor:",
	License:      "Lizenz:",
	VisitWebsite: "Website besuchen",
	Close:        "Schließen",

	Install: "Installieren",
	InstallAsk: func(dir string) string {
		return "AltSysInfo wird nach " + dir + " kopiert und zum Startmenü, zum Desktop und zur Liste der installierten Apps hinzugefügt. Ist es dort bereits installiert, wird es aktualisiert. Die Einstellungen bleiben erhalten.\n\nAltSysInfo startet danach von dort neu."
	},
	InstallFailed: func(err string) string { return "AltSysInfo konnte nicht installiert werden:\n\n" + err },
	Installed:     "AltSysInfo ist installiert",
	UninstallAsk: func(dir string) string {
		return "AltSysInfo von diesem Computer entfernen?\n\nDie Einstellungen bleiben in " + dir
	},
	UninstallRunning: "AltSysInfo läuft. Schließen Sie es und versuchen Sie es erneut.",
	Uninstalled:      "AltSysInfo wurde entfernt.",
}

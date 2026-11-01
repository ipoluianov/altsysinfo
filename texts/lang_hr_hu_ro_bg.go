package texts

import "fmt"

var hr = Strings{
	Error: "Pogreška",

	ToolRefresh:   "Osvježi",
	ToolPdfReport: "Spremi PDF izvješće",
	ToolTray:      "Minimiziraj u područje obavijesti",
	TrayShow:      "Prikaži",
	TrayQuit:      "Izlaz",
	Settings:      "Postavke",
	Help:          "Pomoć",
	About:         "O programu",

	Categories: "Kategorije",
	Category: CategoryStrings{
		Common:          "Općenito",
		RAM:             "Memorija",
		System:          "Sustav",
		CPU:             "Procesor",
		Storage:         "Pohrana",
		Graphics:        "Grafika i zasloni",
		Network:         "Mreža",
		PCI:             "PCI uređaji",
		USB:             "USB uređaji",
		Audio:           "Zvuk",
		Sensors:         "Senzori",
		Battery:         "Baterija",
		PCIVendorDevice: "PCI proizvođači/uređaji",
	},

	ColName:          "Naziv",
	ColValue:         "Vrijednost",
	ColModel:         "Model",
	ColSize:          "Veličina",
	ColSpeed:         "Brzina",
	ColVendorID:      "ID proizvođača",
	ColDeviceID:      "ID uređaja",
	ColVendorName:    "Proizvođač",
	ColDeviceName:    "Uređaj",
	VendorDeviceName: "Proizvođač/uređaj",
	CPUModel:         "Model procesora",
	CPUCores:         "Jezgre procesora",
	Drive:            func(n int) string { return fmt.Sprintf("Disk %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("Grafička kartica %d", n) },
	Driver:           "Upravljački program",
	Copied:           "Kopirano",
	Loading:          "Učitavanje…",
	NoData:           "Nema dostupnih podataka",
	CannotLoad:       func(err string) string { return "Nije moguće učitati podatke: " + err },

	SaveReportTitle: "Spremi PDF izvješće",
	PDFFiles:        "PDF datoteke",
	CannotCollect: func(err string) string {
		return "Nije moguće prikupiti podatke o sustavu: " + err
	},
	CannotGenerate: func(err string) string { return "Nije moguće izraditi izvješće: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "Nije moguće otvoriti dijalog za spremanje: " + err
	},
	CannotSave: func(err string) string { return "Nije moguće spremiti izvješće: " + err },
	SavedTo:    func(path string) string { return "Izvješće spremljeno u " + path },

	AlwaysOnTop:    "Uvijek na vrhu",
	Language:       "Jezik:",
	LanguageSystem: "Kao u sustavu",
	Theme:          "Tema:",
	ThemeDark:      "Tamna",
	ThemeLight:     "Svijetla",

	AboutTitle:   func(name string) string { return "O programu " + name },
	Version:      "Verzija",
	Author:       "Autor:",
	License:      "Licenca:",
	VisitWebsite: "Posjeti web-mjesto",
	Close:        "Zatvori",

	Install:   "Instaliraj",
	Update:    "Ažuriraj",
	Uninstall: "Deinstaliraj",
	InstallAsk: func(dir string) string {
		return "AltSysInfo će se kopirati u " + dir + " i dodati u izbornik Start, na radnu površinu i na popis instaliranih aplikacija. Ako je ondje već instaliran, bit će ažuriran. Postavke ostaju nepromijenjene.\n\nAltSysInfo će se zatim ponovno pokrenuti odande."
	},
	InstallFailed: func(err string) string { return "AltSysInfo nije moguće instalirati:\n\n" + err },
	Installed:     "AltSysInfo je instaliran",
	UninstallAsk: func(dir string) string {
		return "Ukloniti AltSysInfo s ovog računala?\n\nPostavke ostaju u " + dir
	},
	UninstallRunning: "AltSysInfo je pokrenut. Zatvorite ga i pokušajte ponovno.",
	Uninstalled:      "AltSysInfo je uklonjen.",
}

var hu = Strings{
	Error: "Hiba",

	ToolRefresh:   "Frissítés",
	ToolPdfReport: "PDF-jelentés mentése",
	ToolTray:      "Kis méret az értesítési területre",
	TrayShow:      "Megjelenítés",
	TrayQuit:      "Kilépés",
	Settings:      "Beállítások",
	Help:          "Súgó",
	About:         "Névjegy",

	Categories: "Kategóriák",
	Category: CategoryStrings{
		Common:          "Általános",
		RAM:             "Memória",
		System:          "Rendszer",
		CPU:             "Processzor",
		Storage:         "Tárolók",
		Graphics:        "Grafika és kijelzők",
		Network:         "Hálózat",
		PCI:             "PCI-eszközök",
		USB:             "USB-eszközök",
		Audio:           "Hang",
		Sensors:         "Érzékelők",
		Battery:         "Akkumulátor",
		PCIVendorDevice: "PCI gyártók/eszközök",
	},

	ColName:          "Név",
	ColValue:         "Érték",
	ColModel:         "Modell",
	ColSize:          "Méret",
	ColSpeed:         "Sebesség",
	ColVendorID:      "Gyártóazonosító",
	ColDeviceID:      "Eszközazonosító",
	ColVendorName:    "Gyártó",
	ColDeviceName:    "Eszköz",
	VendorDeviceName: "Gyártó/eszköz",
	CPUModel:         "Processzormodell",
	CPUCores:         "Processzormagok",
	Drive:            func(n int) string { return fmt.Sprintf("%d. meghajtó", n) },
	GPU:              func(n int) string { return fmt.Sprintf("%d. videokártya", n) },
	Driver:           "Illesztőprogram",
	Copied:           "Másolva",
	Loading:          "Betöltés…",
	NoData:           "Nincs elérhető adat",
	CannotLoad:       func(err string) string { return "Az adatok nem tölthetők be: " + err },

	SaveReportTitle: "PDF-jelentés mentése",
	PDFFiles:        "PDF-fájlok",
	CannotCollect: func(err string) string {
		return "A rendszeradatok nem gyűjthetők össze: " + err
	},
	CannotGenerate: func(err string) string { return "A jelentés nem hozható létre: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "A mentési párbeszédpanel nem nyitható meg: " + err
	},
	CannotSave: func(err string) string { return "A jelentés nem menthető: " + err },
	SavedTo:    func(path string) string { return "Jelentés mentve ide: " + path },

	AlwaysOnTop:    "Mindig felül",
	Language:       "Nyelv:",
	LanguageSystem: "A rendszer szerint",
	Theme:          "Téma:",
	ThemeDark:      "Sötét",
	ThemeLight:     "Világos",

	AboutTitle:   func(name string) string { return name + " névjegye" },
	Version:      "Verzió",
	Author:       "Szerző:",
	License:      "Licenc:",
	VisitWebsite: "Webhely megnyitása",
	Close:        "Bezárás",

	Install:   "Telepítés",
	Update:    "Frissítés",
	Uninstall: "Eltávolítás",
	InstallAsk: func(dir string) string {
		return "Az AltSysInfo a következő helyre lesz másolva: " + dir + ". Bekerül a Start menübe, az asztalra és a telepített alkalmazások listájába. Ha már telepítve van ott, frissül. A beállítások változatlanok maradnak.\n\nEzután az AltSysInfo onnan indul újra."
	},
	InstallFailed: func(err string) string { return "Az AltSysInfo nem telepíthető:\n\n" + err },
	Installed:     "Az AltSysInfo telepítve van",
	UninstallAsk: func(dir string) string {
		return "Eltávolítja az AltSysInfót erről a számítógépről?\n\nA beállítások itt maradnak: " + dir
	},
	UninstallRunning: "Az AltSysInfo fut. Zárja be, és próbálja újra.",
	Uninstalled:      "Az AltSysInfo el lett távolítva.",
}

var ro = Strings{
	Error: "Eroare",

	ToolRefresh:   "Reîmprospătare",
	ToolPdfReport: "Salvare raport PDF",
	ToolTray:      "Minimizare în zona de notificare",
	TrayShow:      "Afișare",
	TrayQuit:      "Ieșire",
	Settings:      "Setări",
	Help:          "Ajutor",
	About:         "Despre",

	Categories: "Categorii",
	Category: CategoryStrings{
		Common:          "General",
		RAM:             "Memorie",
		System:          "Sistem",
		CPU:             "Procesor",
		Storage:         "Stocare",
		Graphics:        "Grafică și afișaje",
		Network:         "Rețea",
		PCI:             "Dispozitive PCI",
		USB:             "Dispozitive USB",
		Audio:           "Audio",
		Sensors:         "Senzori",
		Battery:         "Baterie",
		PCIVendorDevice: "Producători/dispozitive PCI",
	},

	ColName:          "Nume",
	ColValue:         "Valoare",
	ColModel:         "Model",
	ColSize:          "Dimensiune",
	ColSpeed:         "Viteză",
	ColVendorID:      "ID producător",
	ColDeviceID:      "ID dispozitiv",
	ColVendorName:    "Producător",
	ColDeviceName:    "Dispozitiv",
	VendorDeviceName: "Producător/dispozitiv",
	CPUModel:         "Model procesor",
	CPUCores:         "Nuclee procesor",
	Drive:            func(n int) string { return fmt.Sprintf("Unitatea %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("Placa video %d", n) },
	Driver:           "Driver",
	Copied:           "Copiat",
	Loading:          "Se încarcă…",
	NoData:           "Nu sunt date disponibile",
	CannotLoad:       func(err string) string { return "Informațiile nu pot fi încărcate: " + err },

	SaveReportTitle: "Salvare raport PDF",
	PDFFiles:        "Fișiere PDF",
	CannotCollect: func(err string) string {
		return "Informațiile despre sistem nu pot fi colectate: " + err
	},
	CannotGenerate: func(err string) string { return "Raportul nu poate fi generat: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "Fereastra de salvare nu poate fi deschisă: " + err
	},
	CannotSave: func(err string) string { return "Raportul nu poate fi salvat: " + err },
	SavedTo:    func(path string) string { return "Raport salvat în " + path },

	AlwaysOnTop:    "Mereu deasupra",
	Language:       "Limbă:",
	LanguageSystem: "Ca în sistem",
	Theme:          "Temă:",
	ThemeDark:      "Întunecată",
	ThemeLight:     "Luminoasă",

	AboutTitle:   func(name string) string { return "Despre " + name },
	Version:      "Versiune",
	Author:       "Autor:",
	License:      "Licență:",
	VisitWebsite: "Vizitați site-ul web",
	Close:        "Închidere",

	Install:   "Instalare",
	Update:    "Actualizare",
	Uninstall: "Dezinstalare",
	InstallAsk: func(dir string) string {
		return "AltSysInfo va fi copiat în " + dir + " și adăugat în meniul Start, pe desktop și în lista de aplicații instalate. Dacă este deja instalat acolo, va fi actualizat. Setările rămân neschimbate.\n\nApoi AltSysInfo va reporni de acolo."
	},
	InstallFailed: func(err string) string { return "AltSysInfo nu a putut fi instalat:\n\n" + err },
	Installed:     "AltSysInfo este instalat",
	UninstallAsk: func(dir string) string {
		return "Eliminați AltSysInfo de pe acest computer?\n\nSetările rămân în " + dir
	},
	UninstallRunning: "AltSysInfo rulează. Închideți-l și încercați din nou.",
	Uninstalled:      "AltSysInfo a fost eliminat.",
}

var bg = Strings{
	Error: "Грешка",

	ToolRefresh:   "Опресняване",
	ToolPdfReport: "Запазване на PDF отчет",
	ToolTray:      "Минимизиране в областта за уведомяване",
	TrayShow:      "Показване",
	TrayQuit:      "Изход",
	Settings:      "Настройки",
	Help:          "Помощ",
	About:         "Относно",

	Categories: "Категории",
	Category: CategoryStrings{
		Common:          "Общи",
		RAM:             "Памет",
		System:          "Система",
		CPU:             "Процесор",
		Storage:         "Съхранение",
		Graphics:        "Графика и дисплеи",
		Network:         "Мрежа",
		PCI:             "PCI устройства",
		USB:             "USB устройства",
		Audio:           "Звук",
		Sensors:         "Сензори",
		Battery:         "Батерия",
		PCIVendorDevice: "Справочник PCI",
	},

	ColName:          "Име",
	ColValue:         "Стойност",
	ColModel:         "Модел",
	ColSize:          "Размер",
	ColSpeed:         "Скорост",
	ColVendorID:      "ID на производител",
	ColDeviceID:      "ID на устройство",
	ColVendorName:    "Производител",
	ColDeviceName:    "Устройство",
	VendorDeviceName: "Производител/устройство",
	CPUModel:         "Модел на процесора",
	CPUCores:         "Ядра на процесора",
	Drive:            func(n int) string { return fmt.Sprintf("Диск %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("Видеокарта %d", n) },
	Driver:           "Драйвер",
	Copied:           "Копирано",
	Loading:          "Зареждане…",
	NoData:           "Няма налични данни",
	CannotLoad: func(err string) string {
		return "Информацията не може да бъде заредена: " + err
	},

	SaveReportTitle: "Запазване на PDF отчет",
	PDFFiles:        "PDF файлове",
	CannotCollect: func(err string) string {
		return "Информацията за системата не може да бъде събрана: " + err
	},
	CannotGenerate: func(err string) string { return "Отчетът не може да бъде създаден: " + err },
	CannotOpenSaveDialog: func(err string) string {
		return "Прозорецът за запазване не може да бъде отворен: " + err
	},
	CannotSave: func(err string) string { return "Отчетът не може да бъде запазен: " + err },
	SavedTo:    func(path string) string { return "Отчетът е запазен в " + path },

	AlwaysOnTop:    "Винаги отгоре",
	Language:       "Език:",
	LanguageSystem: "Както в системата",
	Theme:          "Тема:",
	ThemeDark:      "Тъмна",
	ThemeLight:     "Светла",

	AboutTitle:   func(name string) string { return "Относно " + name },
	Version:      "Версия",
	Author:       "Автор:",
	License:      "Лиценз:",
	VisitWebsite: "Посетете уебсайта",
	Close:        "Затваряне",

	Install:   "Инсталиране",
	Update:    "Актуализиране",
	Uninstall: "Деинсталиране",
	InstallAsk: func(dir string) string {
		return "AltSysInfo ще бъде копиран в " + dir + " и добавен в менюто „Старт“, на работния плот и в списъка с инсталирани приложения. Ако вече е инсталиран там, ще бъде актуализиран. Настройките остават непроменени.\n\nСлед това AltSysInfo ще се рестартира оттам."
	},
	InstallFailed: func(err string) string {
		return "AltSysInfo не можа да бъде инсталиран:\n\n" + err
	},
	Installed: "AltSysInfo е инсталиран",
	UninstallAsk: func(dir string) string {
		return "Да се премахне ли AltSysInfo от този компютър?\n\nНастройките остават в " + dir
	},
	UninstallRunning: "AltSysInfo работи. Затворете го и опитайте отново.",
	Uninstalled:      "AltSysInfo е премахнат.",
}

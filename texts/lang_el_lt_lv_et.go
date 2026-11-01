package texts

import "fmt"

var el = Strings{
	Error: "Σφάλμα",

	ToolRefresh:   "Ανανέωση",
	ToolPdfReport: "Αποθήκευση αναφοράς PDF",
	ToolTray:      "Ελαχιστοποίηση στην περιοχή ειδοποιήσεων",
	TrayShow:      "Εμφάνιση",
	TrayQuit:      "Έξοδος",
	Settings:      "Ρυθμίσεις",
	Help:          "Βοήθεια",
	About:         "Πληροφορίες",

	Categories: "Κατηγορίες",
	Category: CategoryStrings{
		Common:          "Γενικά",
		RAM:             "Μνήμη RAM",
		System:          "Σύστημα",
		CPU:             "Επεξεργαστής",
		Storage:         "Αποθήκευση",
		Graphics:        "Γραφικά και οθόνες",
		Network:         "Δίκτυο",
		PCI:             "Συσκευές PCI",
		USB:             "Συσκευές USB",
		Audio:           "Ήχος",
		Sensors:         "Αισθητήρες",
		Battery:         "Μπαταρία",
		PCIVendorDevice: "Κατάλογος PCI",
	},

	ColName:          "Όνομα",
	ColValue:         "Τιμή",
	ColModel:         "Μοντέλο",
	ColSize:          "Μέγεθος",
	ColSpeed:         "Ταχύτητα",
	ColVendorID:      "ID κατασκευαστή",
	ColDeviceID:      "ID συσκευής",
	ColVendorName:    "Κατασκευαστής",
	ColDeviceName:    "Συσκευή",
	VendorDeviceName: "Κατασκευαστής/συσκευή",
	CPUModel:         "Μοντέλο επεξεργαστή",
	CPUCores:         "Πυρήνες επεξεργαστή",
	Drive:            func(n int) string { return fmt.Sprintf("Δίσκος %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Πρόγραμμα οδήγησης",
	Copied:           "Αντιγράφηκε",
	Loading:          "Φόρτωση…",
	NoData:           "Δεν υπάρχουν δεδομένα",
	CannotLoad: func(err string) string {
		return "Δεν ήταν δυνατή η φόρτωση των πληροφοριών: " + err
	},

	SaveReportTitle: "Αποθήκευση αναφοράς PDF",
	PDFFiles:        "Αρχεία PDF",
	CannotCollect: func(err string) string {
		return "Δεν ήταν δυνατή η συλλογή πληροφοριών συστήματος: " + err
	},
	CannotGenerate: func(err string) string {
		return "Δεν ήταν δυνατή η δημιουργία της αναφοράς: " + err
	},
	CannotOpenSaveDialog: func(err string) string {
		return "Δεν ήταν δυνατό το άνοιγμα του παραθύρου αποθήκευσης: " + err
	},
	CannotSave: func(err string) string {
		return "Δεν ήταν δυνατή η αποθήκευση της αναφοράς: " + err
	},
	SavedTo: func(path string) string { return "Η αναφορά αποθηκεύτηκε στο " + path },

	AlwaysOnTop:    "Πάντα στο προσκήνιο",
	Language:       "Γλώσσα:",
	LanguageSystem: "Όπως στο σύστημα",
	Theme:          "Θέμα:",
	ThemeDark:      "Σκούρο",
	ThemeLight:     "Ανοιχτόχρωμο",

	AboutTitle:   func(name string) string { return "Πληροφορίες για το " + name },
	Version:      "Έκδοση",
	Author:       "Δημιουργός:",
	License:      "Άδεια χρήσης:",
	VisitWebsite: "Επίσκεψη στην τοποθεσία web",
	Close:        "Κλείσιμο",

	Install:   "Εγκατάσταση",
	Update:    "Ενημέρωση",
	Uninstall: "Κατάργηση εγκατάστασης",
	InstallAsk: func(dir string) string {
		return "Το AltSysInfo θα αντιγραφεί στο " + dir + " και θα προστεθεί στο μενού Έναρξη, στην επιφάνεια εργασίας και στη λίστα εγκατεστημένων εφαρμογών. Αν είναι ήδη εγκατεστημένο εκεί, θα ενημερωθεί. Οι ρυθμίσεις παραμένουν ως έχουν.\n\nΣτη συνέχεια το AltSysInfo θα επανεκκινηθεί από εκεί."
	},
	InstallFailed: func(err string) string {
		return "Δεν ήταν δυνατή η εγκατάσταση του AltSysInfo:\n\n" + err
	},
	Installed: "Το AltSysInfo εγκαταστάθηκε",
	UninstallAsk: func(dir string) string {
		return "Κατάργηση του AltSysInfo από αυτόν τον υπολογιστή;\n\nΟι ρυθμίσεις παραμένουν στο " + dir
	},
	UninstallRunning: "Το AltSysInfo εκτελείται. Κλείστε το και δοκιμάστε ξανά.",
	Uninstalled:      "Το AltSysInfo καταργήθηκε.",
}

var lt = Strings{
	Error: "Klaida",

	ToolRefresh:   "Atnaujinti",
	ToolPdfReport: "Įrašyti PDF ataskaitą",
	ToolTray:      "Sumažinti į pranešimų sritį",
	TrayShow:      "Rodyti",
	TrayQuit:      "Išeiti",
	Settings:      "Parametrai",
	Help:          "Žinynas",
	About:         "Apie",

	Categories: "Kategorijos",
	Category: CategoryStrings{
		Common:          "Bendra",
		RAM:             "Operatyvioji atmintis",
		System:          "Sistema",
		CPU:             "Procesorius",
		Storage:         "Diskai",
		Graphics:        "Grafika ir ekranai",
		Network:         "Tinklas",
		PCI:             "PCI įrenginiai",
		USB:             "USB įrenginiai",
		Audio:           "Garsas",
		Sensors:         "Jutikliai",
		Battery:         "Baterija",
		PCIVendorDevice: "PCI gamintojai/įrenginiai",
	},

	ColName:          "Pavadinimas",
	ColValue:         "Reikšmė",
	ColModel:         "Modelis",
	ColSize:          "Dydis",
	ColSpeed:         "Sparta",
	ColVendorID:      "Gamintojo ID",
	ColDeviceID:      "Įrenginio ID",
	ColVendorName:    "Gamintojas",
	ColDeviceName:    "Įrenginys",
	VendorDeviceName: "Gamintojas/įrenginys",
	CPUModel:         "Procesoriaus modelis",
	CPUCores:         "Procesoriaus branduoliai",
	Drive:            func(n int) string { return fmt.Sprintf("Diskas %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Tvarkyklė",
	Copied:           "Nukopijuota",
	Loading:          "Įkeliama…",
	NoData:           "Duomenų nėra",
	CannotLoad:       func(err string) string { return "Nepavyko gauti informacijos: " + err },

	SaveReportTitle:      "Įrašyti PDF ataskaitą",
	PDFFiles:             "PDF failai",
	CannotCollect:        func(err string) string { return "Nepavyko surinkti sistemos informacijos: " + err },
	CannotGenerate:       func(err string) string { return "Nepavyko sukurti ataskaitos: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Nepavyko atidaryti įrašymo lango: " + err },
	CannotSave:           func(err string) string { return "Nepavyko įrašyti ataskaitos: " + err },
	SavedTo:              func(path string) string { return "Ataskaita įrašyta į " + path },

	AlwaysOnTop:    "Visada viršuje",
	Language:       "Kalba:",
	LanguageSystem: "Kaip sistemoje",
	Theme:          "Tema:",
	ThemeDark:      "Tamsi",
	ThemeLight:     "Šviesi",

	AboutTitle:   func(name string) string { return "Apie " + name },
	Version:      "Versija",
	Author:       "Autorius:",
	License:      "Licencija:",
	VisitWebsite: "Apsilankyti svetainėje",
	Close:        "Uždaryti",

	Install:   "Įdiegti",
	Update:    "Naujinti",
	Uninstall: "Pašalinti",
	InstallAsk: func(dir string) string {
		return "AltSysInfo bus nukopijuotas į " + dir + " ir įtrauktas į meniu Pradžia, darbalaukį ir įdiegtų programų sąrašą. Jei jis ten jau įdiegtas, bus atnaujintas. Parametrai liks nepakeisti.\n\nTada AltSysInfo bus paleistas iš naujo iš ten."
	},
	InstallFailed: func(err string) string { return "Nepavyko įdiegti AltSysInfo:\n\n" + err },
	Installed:     "AltSysInfo įdiegtas",
	UninstallAsk: func(dir string) string {
		return "Pašalinti AltSysInfo iš šio kompiuterio?\n\nParametrai liks aplanke " + dir
	},
	UninstallRunning: "AltSysInfo veikia. Uždarykite jį ir bandykite dar kartą.",
	Uninstalled:      "AltSysInfo pašalintas.",
}

var lv = Strings{
	Error: "Kļūda",

	ToolRefresh:   "Atsvaidzināt",
	ToolPdfReport: "Saglabāt PDF atskaiti",
	ToolTray:      "Minimizēt paziņojumu apgabalā",
	TrayShow:      "Rādīt",
	TrayQuit:      "Iziet",
	Settings:      "Iestatījumi",
	Help:          "Palīdzība",
	About:         "Par",

	Categories: "Kategorijas",
	Category: CategoryStrings{
		Common:          "Vispārīgi",
		RAM:             "Operatīvā atmiņa",
		System:          "Sistēma",
		CPU:             "Procesors",
		Storage:         "Diski",
		Graphics:        "Grafika un displeji",
		Network:         "Tīkls",
		PCI:             "PCI ierīces",
		USB:             "USB ierīces",
		Audio:           "Skaņa",
		Sensors:         "Sensori",
		Battery:         "Akumulators",
		PCIVendorDevice: "PCI ražotāji/ierīces",
	},

	ColName:          "Nosaukums",
	ColValue:         "Vērtība",
	ColModel:         "Modelis",
	ColSize:          "Izmērs",
	ColSpeed:         "Ātrums",
	ColVendorID:      "Ražotāja ID",
	ColDeviceID:      "Ierīces ID",
	ColVendorName:    "Ražotājs",
	ColDeviceName:    "Ierīce",
	VendorDeviceName: "Ražotājs/ierīce",
	CPUModel:         "Procesora modelis",
	CPUCores:         "Procesora kodoli",
	Drive:            func(n int) string { return fmt.Sprintf("Disks %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Draiveris",
	Copied:           "Nokopēts",
	Loading:          "Ielādē…",
	NoData:           "Nav datu",
	CannotLoad:       func(err string) string { return "Neizdevās iegūt informāciju: " + err },

	SaveReportTitle:      "Saglabāt PDF atskaiti",
	PDFFiles:             "PDF faili",
	CannotCollect:        func(err string) string { return "Neizdevās savākt sistēmas informāciju: " + err },
	CannotGenerate:       func(err string) string { return "Neizdevās izveidot atskaiti: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Neizdevās atvērt saglabāšanas logu: " + err },
	CannotSave:           func(err string) string { return "Neizdevās saglabāt atskaiti: " + err },
	SavedTo:              func(path string) string { return "Atskaite saglabāta: " + path },

	AlwaysOnTop:    "Vienmēr virspusē",
	Language:       "Valoda:",
	LanguageSystem: "Kā sistēmā",
	Theme:          "Motīvs:",
	ThemeDark:      "Tumšs",
	ThemeLight:     "Gaišs",

	AboutTitle:   func(name string) string { return "Par " + name },
	Version:      "Versija",
	Author:       "Autors:",
	License:      "Licence:",
	VisitWebsite: "Apmeklēt tīmekļa vietni",
	Close:        "Aizvērt",

	Install:   "Instalēt",
	Update:    "Atjaunināt",
	Uninstall: "Atinstalēt",
	InstallAsk: func(dir string) string {
		return "AltSysInfo tiks iekopēts mapē " + dir + " un pievienots izvēlnei Sākt, darbvirsmai un instalēto programmu sarakstam. Ja tas tur jau ir instalēts, tas tiks atjaunināts. Iestatījumi paliks nemainīti.\n\nPēc tam AltSysInfo tiks restartēts no turienes."
	},
	InstallFailed: func(err string) string { return "AltSysInfo neizdevās instalēt:\n\n" + err },
	Installed:     "AltSysInfo ir instalēts",
	UninstallAsk: func(dir string) string {
		return "Vai noņemt AltSysInfo no šī datora?\n\nIestatījumi paliks mapē " + dir
	},
	UninstallRunning: "AltSysInfo darbojas. Aizveriet to un mēģiniet vēlreiz.",
	Uninstalled:      "AltSysInfo ir noņemts.",
}

var et = Strings{
	Error: "Viga",

	ToolRefresh:   "Värskenda",
	ToolPdfReport: "Salvesta PDF-aruanne",
	ToolTray:      "Minimeeri teavitusalasse",
	TrayShow:      "Kuva",
	TrayQuit:      "Välju",
	Settings:      "Sätted",
	Help:          "Spikker",
	About:         "Teave",

	Categories: "Kategooriad",
	Category: CategoryStrings{
		Common:          "Üldine",
		RAM:             "Mälu",
		System:          "Süsteem",
		CPU:             "Protsessor",
		Storage:         "Kettad",
		Graphics:        "Graafika ja kuvarid",
		Network:         "Võrk",
		PCI:             "PCI-seadmed",
		USB:             "USB-seadmed",
		Audio:           "Heli",
		Sensors:         "Andurid",
		Battery:         "Aku",
		PCIVendorDevice: "PCI tootjad/seadmed",
	},

	ColName:          "Nimi",
	ColValue:         "Väärtus",
	ColModel:         "Mudel",
	ColSize:          "Maht",
	ColSpeed:         "Kiirus",
	ColVendorID:      "Tootja ID",
	ColDeviceID:      "Seadme ID",
	ColVendorName:    "Tootja",
	ColDeviceName:    "Seade",
	VendorDeviceName: "Tootja/seade",
	CPUModel:         "Protsessori mudel",
	CPUCores:         "Protsessori tuumad",
	Drive:            func(n int) string { return fmt.Sprintf("Ketas %d", n) },
	GPU:              func(n int) string { return fmt.Sprintf("GPU %d", n) },
	Driver:           "Draiver",
	Copied:           "Kopeeritud",
	Loading:          "Laadimine…",
	NoData:           "Andmed puuduvad",
	CannotLoad:       func(err string) string { return "Teavet ei saanud laadida: " + err },

	SaveReportTitle:      "Salvesta PDF-aruanne",
	PDFFiles:             "PDF-failid",
	CannotCollect:        func(err string) string { return "Süsteemiteavet ei saanud koguda: " + err },
	CannotGenerate:       func(err string) string { return "Aruannet ei saanud luua: " + err },
	CannotOpenSaveDialog: func(err string) string { return "Salvestusakent ei saanud avada: " + err },
	CannotSave:           func(err string) string { return "Aruannet ei saanud salvestada: " + err },
	SavedTo:              func(path string) string { return "Aruanne salvestatud: " + path },

	AlwaysOnTop:    "Alati pealmine",
	Language:       "Keel:",
	LanguageSystem: "Nagu süsteemis",
	Theme:          "Kujundus:",
	ThemeDark:      "Tume",
	ThemeLight:     "Hele",

	AboutTitle:   func(name string) string { return "Teave: " + name },
	Version:      "Versioon",
	Author:       "Autor:",
	License:      "Litsents:",
	VisitWebsite: "Külasta veebisaiti",
	Close:        "Sule",

	Install:   "Installi",
	Update:    "Värskenda",
	Uninstall: "Desinstalli",
	InstallAsk: func(dir string) string {
		return "AltSysInfo kopeeritakse kausta " + dir + " ning lisatakse menüüsse Start, töölauale ja installitud rakenduste loendisse. Kui see on seal juba installitud, siis seda värskendatakse. Sätted jäävad muutmata.\n\nSeejärel käivitub AltSysInfo sealt uuesti."
	},
	InstallFailed: func(err string) string { return "AltSysInfot ei saanud installida:\n\n" + err },
	Installed:     "AltSysInfo on installitud",
	UninstallAsk: func(dir string) string {
		return "Kas eemaldada AltSysInfo sellest arvutist?\n\nSätted jäävad kausta " + dir
	},
	UninstallRunning: "AltSysInfo töötab. Sulgege see ja proovige uuesti.",
	Uninstalled:      "AltSysInfo on eemaldatud.",
}

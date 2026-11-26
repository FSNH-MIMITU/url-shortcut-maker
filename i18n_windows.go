package main

import (
	"errors"

	"github.com/FSNH-MIMITU/url-shortcut-maker/core"
)

// language 描述一种界面语言及其推荐字体（保证该语言字形可正常显示）。
type language struct {
	code   string
	name   string
	family string
}

var languages = []language{
	{"zh-CN", "简体中文 (Chinese Simplified)", "Microsoft YaHei UI"},
	{"zh-TW", "繁體中文 (Chinese Traditional)", "Microsoft YaHei UI"},
	{"en", "English", "Segoe UI"},
	{"ja", "日本語 (Japanese)", "Yu Gothic UI"},
	{"ko", "한국어 (Korean)", "Malgun Gothic"},
	{"de", "Deutsch", "Segoe UI"},
	{"fr", "Français", "Segoe UI"},
	{"es", "Español", "Segoe UI"},
	{"pt", "Português", "Segoe UI"},
	{"ru", "Русский", "Segoe UI"},
}

type uiStrings struct {
	Title    string
	Subtitle string

	LabelURL       string
	LabelIcon      string
	Browse         string
	LabelOpenWith  string
	DefaultBrowser string
	HintBrowser    string
	LabelName      string
	HintName       string
	LabelPreview   string
	Desc           string
	Generate       string
	StatusInit     string
	LabelLang      string

	DlgPickTitle string
	FilterImages string
	FilterAll    string

	ErrImageProcess  string
	ErrImageInvalid  string
	StatusIconPicked string
	ErrNeedURL       string
	ErrNeedIcon      string
	ErrIconsDir      string
	ErrIconProcess   string

	BrowserFallbackFmt string
	ShortcutDescFmt    string
	SystemDefault      string

	ErrCreateFailed   string
	ErrCreateShortcut string
	StatusCreated     string
	InfoTitle         string
	InfoCreatedFmt    string
	InfoOpenWithFmt   string
	InfoTail          string

	ErrWindowClass  string
	ErrWindowCreate string

	ErrURLEmpty  string
	ErrURLFormat string
	ErrURLHost   string
	ErrURLScheme string
}

var translations = map[string]uiStrings{
	"zh-CN": {
		Title:    "网址快捷方式生成器",
		Subtitle: "填写网址 · 选择图标 · 一键生成桌面快捷方式",

		LabelURL:       "网址",
		LabelIcon:      "图标",
		Browse:         "浏览…",
		LabelOpenWith:  "打开方式",
		DefaultBrowser: "默认浏览器（推荐）",
		HintBrowser:    "未安装所选浏览器时，\n自动使用系统默认浏览器",
		LabelName:      "名称",
		HintName:       "留空自动生成",
		LabelPreview:   "图标预览",
		Desc:           "快捷方式会保存到桌面。\n双击即可用所选浏览器打开网址。\n图标将根据所选图片自动生成。",
		Generate:       "生成快捷方式",
		StatusInit:     "填写网址并选择图标，然后点击「生成快捷方式」。",
		LabelLang:      "语言",

		DlgPickTitle: "选择图标图片",
		FilterImages: "图片文件",
		FilterAll:    "所有文件",

		ErrImageProcess:  "无法处理该图片：",
		ErrImageInvalid:  "该图片无法作为图标使用，请换一张试试。",
		StatusIconPicked: "图标已选择：",
		ErrNeedURL:       "请先填写要打开的网址。",
		ErrNeedIcon:      "请先点击「浏览…」选择图标图片。",
		ErrIconsDir:      "无法创建图标目录：",
		ErrIconProcess:   "处理图标失败：",

		BrowserFallbackFmt: "%s 无法使用，已改用系统默认浏览器。\n\n",
		ShortcutDescFmt:    "打开 %s",
		SystemDefault:      "系统默认浏览器",

		ErrCreateFailed:   "生成失败：",
		ErrCreateShortcut: "生成快捷方式失败：\n",
		StatusCreated:     "已生成：",
		InfoTitle:         "生成成功",
		InfoCreatedFmt:    "已在桌面生成快捷方式：\n%s",
		InfoOpenWithFmt:   "\n\n打开方式：%s",
		InfoTail:          "\n双击快捷方式即可打开该网址。",

		ErrWindowClass:  "窗口类注册失败",
		ErrWindowCreate: "窗口创建失败",

		ErrURLEmpty:  "请先填写要打开的网址。",
		ErrURLFormat: "网址格式不正确。",
		ErrURLHost:   "请填写完整的网址，例如 www.example.com",
		ErrURLScheme: "仅支持 http / https 网址。",
	},
	"zh-TW": {
		Title:    "網址捷徑產生器",
		Subtitle: "填寫網址 · 選擇圖示 · 一鍵產生桌面捷徑",

		LabelURL:       "網址",
		LabelIcon:      "圖示",
		Browse:         "瀏覽…",
		LabelOpenWith:  "開啟方式",
		DefaultBrowser: "預設瀏覽器（建議）",
		HintBrowser:    "未安裝所選瀏覽器時，\n自動使用系統預設瀏覽器",
		LabelName:      "名稱",
		HintName:       "留空自動產生",
		LabelPreview:   "圖示預覽",
		Desc:           "捷徑會儲存到桌面。\n按兩下即可用所選瀏覽器開啟網址。\n圖示將依據所選圖片自動產生。",
		Generate:       "產生捷徑",
		StatusInit:     "填寫網址並選擇圖示，然後點擊「產生捷徑」。",
		LabelLang:      "語言",

		DlgPickTitle: "選擇圖示圖片",
		FilterImages: "圖片檔案",
		FilterAll:    "所有檔案",

		ErrImageProcess:  "無法處理此圖片：",
		ErrImageInvalid:  "此圖片無法作為圖示使用，請換一張試試。",
		StatusIconPicked: "圖示已選擇：",
		ErrNeedURL:       "請先填寫要開啟的網址。",
		ErrNeedIcon:      "請先點擊「瀏覽…」選擇圖示圖片。",
		ErrIconsDir:      "無法建立圖示資料夾：",
		ErrIconProcess:   "處理圖示失敗：",

		BrowserFallbackFmt: "%s 無法使用，已改用系統預設瀏覽器。\n\n",
		ShortcutDescFmt:    "開啟 %s",
		SystemDefault:      "系統預設瀏覽器",

		ErrCreateFailed:   "產生失敗：",
		ErrCreateShortcut: "產生捷徑失敗：\n",
		StatusCreated:     "已產生：",
		InfoTitle:         "產生成功",
		InfoCreatedFmt:    "已在桌面產生捷徑：\n%s",
		InfoOpenWithFmt:   "\n\n開啟方式：%s",
		InfoTail:          "\n按兩下捷徑即可開啟該網址。",

		ErrWindowClass:  "視窗類別註冊失敗",
		ErrWindowCreate: "視窗建立失敗",

		ErrURLEmpty:  "請先填寫要開啟的網址。",
		ErrURLFormat: "網址格式不正確。",
		ErrURLHost:   "請填寫完整的網址，例如 www.example.com",
		ErrURLScheme: "僅支援 http / https 網址。",
	},
	"en": {
		Title:    "URL Shortcut Maker",
		Subtitle: "Enter a URL · Choose an icon · Create a desktop shortcut",

		LabelURL:       "URL",
		LabelIcon:      "Icon",
		Browse:         "Browse…",
		LabelOpenWith:  "Open with",
		DefaultBrowser: "Default browser (recommended)",
		HintBrowser:    "If the selected browser is not\ninstalled, the default one is used",
		LabelName:      "Name",
		HintName:       "Leave empty to auto-generate",
		LabelPreview:   "Icon preview",
		Desc:           "The shortcut is saved to your desktop.\nDouble-click it to open the URL in the chosen browser.\nThe icon is generated from the image you pick.",
		Generate:       "Create shortcut",
		StatusInit:     "Enter a URL and choose an icon, then click “Create shortcut”.",
		LabelLang:      "Language",

		DlgPickTitle: "Choose icon image",
		FilterImages: "Image files",
		FilterAll:    "All files",

		ErrImageProcess:  "Cannot process this image: ",
		ErrImageInvalid:  "This image cannot be used as an icon. Please try another one.",
		StatusIconPicked: "Icon selected: ",
		ErrNeedURL:       "Please enter the URL first.",
		ErrNeedIcon:      "Please click “Browse…” to choose an icon image first.",
		ErrIconsDir:      "Cannot create icon folder: ",
		ErrIconProcess:   "Failed to process the icon: ",

		BrowserFallbackFmt: "%s is not available, using the system default browser.\n\n",
		ShortcutDescFmt:    "Open %s",
		SystemDefault:      "system default browser",

		ErrCreateFailed:   "Failed: ",
		ErrCreateShortcut: "Failed to create the shortcut:\n",
		StatusCreated:     "Created: ",
		InfoTitle:         "Success",
		InfoCreatedFmt:    "Shortcut created on your desktop:\n%s",
		InfoOpenWithFmt:   "\n\nOpens with: %s",
		InfoTail:          "\nDouble-click the shortcut to open the URL.",

		ErrWindowClass:  "Failed to register window class",
		ErrWindowCreate: "Failed to create window",

		ErrURLEmpty:  "Please enter the URL first.",
		ErrURLFormat: "The URL format is invalid.",
		ErrURLHost:   "Please enter a complete URL, e.g. www.example.com",
		ErrURLScheme: "Only http / https URLs are supported.",
	},
	"ja": {
		Title:    "URL ショートカット作成ツール",
		Subtitle: "URL を入力 · アイコンを選択 · デスクトップにショートカットを作成",

		LabelURL:       "URL",
		LabelIcon:      "アイコン",
		Browse:         "参照…",
		LabelOpenWith:  "開く方法",
		DefaultBrowser: "既定のブラウザ（推奨）",
		HintBrowser:    "選択したブラウザが未インストールの場合は\n既定のブラウザを使用します",
		LabelName:      "名前",
		HintName:       "空欄なら自動生成",
		LabelPreview:   "アイコンのプレビュー",
		Desc:           "ショートカットはデスクトップに保存されます。\nダブルクリックで選択したブラウザから URL を開きます。\nアイコンは選択した画像から自動生成されます。",
		Generate:       "ショートカットを作成",
		StatusInit:     "URL を入力してアイコンを選択し、「ショートカットを作成」をクリックしてください。",
		LabelLang:      "言語",

		DlgPickTitle: "アイコン画像を選択",
		FilterImages: "画像ファイル",
		FilterAll:    "すべてのファイル",

		ErrImageProcess:  "この画像を処理できません：",
		ErrImageInvalid:  "この画像はアイコンとして使用できません。別の画像をお試しください。",
		StatusIconPicked: "アイコンを選択しました：",
		ErrNeedURL:       "先に開きたい URL を入力してください。",
		ErrNeedIcon:      "先に「参照…」からアイコン画像を選択してください。",
		ErrIconsDir:      "アイコンフォルダを作成できません：",
		ErrIconProcess:   "アイコンの処理に失敗しました：",

		BrowserFallbackFmt: "%s は使用できないため、既定のブラウザを使用します。\n\n",
		ShortcutDescFmt:    "%s を開く",
		SystemDefault:      "既定のブラウザ",

		ErrCreateFailed:   "作成に失敗しました：",
		ErrCreateShortcut: "ショートカットの作成に失敗しました：\n",
		StatusCreated:     "作成しました：",
		InfoTitle:         "作成完了",
		InfoCreatedFmt:    "デスクトップにショートカットを作成しました：\n%s",
		InfoOpenWithFmt:   "\n\n開く方法：%s",
		InfoTail:          "\nショートカットをダブルクリックすると URL が開きます。",

		ErrWindowClass:  "ウィンドウクラスの登録に失敗しました",
		ErrWindowCreate: "ウィンドウの作成に失敗しました",

		ErrURLEmpty:  "先に開きたい URL を入力してください。",
		ErrURLFormat: "URL の形式が正しくありません。",
		ErrURLHost:   "完全な URL を入力してください（例：www.example.com）",
		ErrURLScheme: "http / https の URL のみ対応しています。",
	},
	"ko": {
		Title:    "URL 바로가기 생성기",
		Subtitle: "URL 입력 · 아이콘 선택 · 바탕 화면 바로가기 생성",

		LabelURL:       "URL",
		LabelIcon:      "아이콘",
		Browse:         "찾아보기…",
		LabelOpenWith:  "열기 방식",
		DefaultBrowser: "기본 브라우저(권장)",
		HintBrowser:    "선택한 브라우저가 설치되어 있지 않으면\n시스템 기본 브라우저를 사용합니다",
		LabelName:      "이름",
		HintName:       "비워 두면 자동 생성",
		LabelPreview:   "아이콘 미리보기",
		Desc:           "바로가기는 바탕 화면에 저장됩니다.\n두 번 클릭하면 선택한 브라우저로 URL이 열립니다.\n아이콘은 선택한 이미지로 자동 생성됩니다.",
		Generate:       "바로가기 만들기",
		StatusInit:     "URL을 입력하고 아이콘을 선택한 뒤 「바로가기 만들기」를 클릭하세요.",
		LabelLang:      "언어",

		DlgPickTitle: "아이콘 이미지 선택",
		FilterImages: "이미지 파일",
		FilterAll:    "모든 파일",

		ErrImageProcess:  "이 이미지를 처리할 수 없습니다: ",
		ErrImageInvalid:  "이 이미지는 아이콘으로 사용할 수 없습니다. 다른 이미지를 선택하세요.",
		StatusIconPicked: "아이콘 선택됨: ",
		ErrNeedURL:       "먼저 열 URL을 입력하세요.",
		ErrNeedIcon:      "먼저 「찾아보기…」에서 아이콘 이미지를 선택하세요.",
		ErrIconsDir:      "아이콘 폴더를 만들 수 없습니다: ",
		ErrIconProcess:   "아이콘 처리 실패: ",

		BrowserFallbackFmt: "%s을(를) 사용할 수 없어 시스템 기본 브라우저로 전환합니다.\n\n",
		ShortcutDescFmt:    "%s 열기",
		SystemDefault:      "시스템 기본 브라우저",

		ErrCreateFailed:   "생성 실패: ",
		ErrCreateShortcut: "바로가기 생성 실패:\n",
		StatusCreated:     "생성됨: ",
		InfoTitle:         "생성 완료",
		InfoCreatedFmt:    "바탕 화면에 바로가기를 만들었습니다:\n%s",
		InfoOpenWithFmt:   "\n\n열기 방식: %s",
		InfoTail:          "\n바로가기를 두 번 클릭하면 URL이 열립니다.",

		ErrWindowClass:  "창 클래스 등록 실패",
		ErrWindowCreate: "창 생성 실패",

		ErrURLEmpty:  "먼저 열 URL을 입력하세요.",
		ErrURLFormat: "URL 형식이 올바르지 않습니다.",
		ErrURLHost:   "전체 URL을 입력하세요. 예: www.example.com",
		ErrURLScheme: "http / https URL만 지원합니다.",
	},
	"de": {
		Title:    "URL-Verknüpfung erstellen",
		Subtitle: "URL eingeben · Symbol wählen · Desktop-Verknüpfung erstellen",

		LabelURL:       "URL",
		LabelIcon:      "Symbol",
		Browse:         "Durchsuchen…",
		LabelOpenWith:  "Öffnen mit",
		DefaultBrowser: "Standardbrowser (empfohlen)",
		HintBrowser:    "Ist der gewählte Browser nicht installiert,\nwird der Standardbrowser verwendet",
		LabelName:      "Name",
		HintName:       "Leer lassen für automatischen Namen",
		LabelPreview:   "Symbolvorschau",
		Desc:           "Die Verknüpfung wird auf dem Desktop gespeichert.\nDoppelklick öffnet die URL im gewählten Browser.\nDas Symbol wird aus dem gewählten Bild erzeugt.",
		Generate:       "Verknüpfung erstellen",
		StatusInit:     "URL eingeben und Symbol wählen, dann auf „Verknüpfung erstellen“ klicken.",
		LabelLang:      "Sprache",

		DlgPickTitle: "Symbolbild auswählen",
		FilterImages: "Bilddateien",
		FilterAll:    "Alle Dateien",

		ErrImageProcess:  "Dieses Bild kann nicht verarbeitet werden: ",
		ErrImageInvalid:  "Dieses Bild kann nicht als Symbol verwendet werden. Bitte ein anderes wählen.",
		StatusIconPicked: "Symbol gewählt: ",
		ErrNeedURL:       "Bitte zuerst die URL eingeben.",
		ErrNeedIcon:      "Bitte zuerst über „Durchsuchen…“ ein Symbolbild wählen.",
		ErrIconsDir:      "Symbolordner kann nicht erstellt werden: ",
		ErrIconProcess:   "Symbol konnte nicht verarbeitet werden: ",

		BrowserFallbackFmt: "%s ist nicht verfügbar, es wird der Standardbrowser verwendet.\n\n",
		ShortcutDescFmt:    "%s öffnen",
		SystemDefault:      "Standardbrowser",

		ErrCreateFailed:   "Fehlgeschlagen: ",
		ErrCreateShortcut: "Verknüpfung konnte nicht erstellt werden:\n",
		StatusCreated:     "Erstellt: ",
		InfoTitle:         "Erfolg",
		InfoCreatedFmt:    "Verknüpfung auf dem Desktop erstellt:\n%s",
		InfoOpenWithFmt:   "\n\nÖffnen mit: %s",
		InfoTail:          "\nDoppelklick auf die Verknüpfung öffnet die URL.",

		ErrWindowClass:  "Fensterklasse konnte nicht registriert werden",
		ErrWindowCreate: "Fenster konnte nicht erstellt werden",

		ErrURLEmpty:  "Bitte zuerst die URL eingeben.",
		ErrURLFormat: "Das URL-Format ist ungültig.",
		ErrURLHost:   "Bitte eine vollständige URL eingeben, z. B. www.example.com",
		ErrURLScheme: "Es werden nur http-/https-URLs unterstützt.",
	},
	"fr": {
		Title:    "Créateur de raccourci URL",
		Subtitle: "Saisir une URL · Choisir une icône · Créer un raccourci sur le bureau",

		LabelURL:       "URL",
		LabelIcon:      "Icône",
		Browse:         "Parcourir…",
		LabelOpenWith:  "Ouvrir avec",
		DefaultBrowser: "Navigateur par défaut (recommandé)",
		HintBrowser:    "Si le navigateur choisi n'est pas installé,\nle navigateur par défaut est utilisé",
		LabelName:      "Nom",
		HintName:       "Laisser vide pour générer automatiquement",
		LabelPreview:   "Aperçu de l'icône",
		Desc:           "Le raccourci sera enregistré sur le bureau.\nDouble-cliquez pour ouvrir l'URL dans le navigateur choisi.\nL'icône est générée à partir de l'image sélectionnée.",
		Generate:       "Créer le raccourci",
		StatusInit:     "Saisissez une URL, choisissez une icône, puis cliquez sur « Créer le raccourci ».",
		LabelLang:      "Langue",

		DlgPickTitle: "Choisir une image d'icône",
		FilterImages: "Fichiers image",
		FilterAll:    "Tous les fichiers",

		ErrImageProcess:  "Impossible de traiter cette image : ",
		ErrImageInvalid:  "Cette image ne peut pas servir d'icône. Essayez-en une autre.",
		StatusIconPicked: "Icône sélectionnée : ",
		ErrNeedURL:       "Veuillez d'abord saisir l'URL.",
		ErrNeedIcon:      "Veuillez d'abord cliquer sur « Parcourir… » pour choisir une icône.",
		ErrIconsDir:      "Impossible de créer le dossier d'icônes : ",
		ErrIconProcess:   "Échec du traitement de l'icône : ",

		BrowserFallbackFmt: "%s n'est pas disponible, utilisation du navigateur par défaut.\n\n",
		ShortcutDescFmt:    "Ouvrir %s",
		SystemDefault:      "navigateur par défaut",

		ErrCreateFailed:   "Échec : ",
		ErrCreateShortcut: "Échec de la création du raccourci :\n",
		StatusCreated:     "Créé : ",
		InfoTitle:         "Succès",
		InfoCreatedFmt:    "Raccourci créé sur le bureau :\n%s",
		InfoOpenWithFmt:   "\n\nOuvre avec : %s",
		InfoTail:          "\nDouble-cliquez sur le raccourci pour ouvrir l'URL.",

		ErrWindowClass:  "Échec de l'enregistrement de la classe de fenêtre",
		ErrWindowCreate: "Échec de la création de la fenêtre",

		ErrURLEmpty:  "Veuillez d'abord saisir l'URL.",
		ErrURLFormat: "Le format de l'URL est incorrect.",
		ErrURLHost:   "Saisissez une URL complète, par ex. www.example.com",
		ErrURLScheme: "Seules les URL http / https sont prises en charge.",
	},
	"es": {
		Title:    "Creador de accesos directos URL",
		Subtitle: "Introduce una URL · Elige un icono · Crea un acceso directo en el escritorio",

		LabelURL:       "URL",
		LabelIcon:      "Icono",
		Browse:         "Examinar…",
		LabelOpenWith:  "Abrir con",
		DefaultBrowser: "Navegador predeterminado (recomendado)",
		HintBrowser:    "Si el navegador elegido no está instalado,\nse usará el predeterminado",
		LabelName:      "Nombre",
		HintName:       "Déjalo vacío para generarlo automáticamente",
		LabelPreview:   "Vista previa del icono",
		Desc:           "El acceso directo se guardará en el escritorio.\nHaz doble clic para abrir la URL con el navegador elegido.\nEl icono se genera a partir de la imagen seleccionada.",
		Generate:       "Crear acceso directo",
		StatusInit:     "Introduce una URL y elige un icono; luego pulsa «Crear acceso directo».",
		LabelLang:      "Idioma",

		DlgPickTitle: "Elegir imagen de icono",
		FilterImages: "Archivos de imagen",
		FilterAll:    "Todos los archivos",

		ErrImageProcess:  "No se puede procesar esta imagen: ",
		ErrImageInvalid:  "Esta imagen no sirve como icono. Prueba con otra.",
		StatusIconPicked: "Icono seleccionado: ",
		ErrNeedURL:       "Introduce primero la URL.",
		ErrNeedIcon:      "Primero pulsa «Examinar…» para elegir una imagen de icono.",
		ErrIconsDir:      "No se puede crear la carpeta de iconos: ",
		ErrIconProcess:   "Error al procesar el icono: ",

		BrowserFallbackFmt: "%s no está disponible; se usará el navegador predeterminado.\n\n",
		ShortcutDescFmt:    "Abrir %s",
		SystemDefault:      "navegador predeterminado",

		ErrCreateFailed:   "Error: ",
		ErrCreateShortcut: "No se pudo crear el acceso directo:\n",
		StatusCreated:     "Creado: ",
		InfoTitle:         "Correcto",
		InfoCreatedFmt:    "Acceso directo creado en el escritorio:\n%s",
		InfoOpenWithFmt:   "\n\nAbrir con: %s",
		InfoTail:          "\nHaz doble clic en el acceso directo para abrir la URL.",

		ErrWindowClass:  "No se pudo registrar la clase de ventana",
		ErrWindowCreate: "No se pudo crear la ventana",

		ErrURLEmpty:  "Introduce primero la URL.",
		ErrURLFormat: "El formato de la URL no es válido.",
		ErrURLHost:   "Introduce una URL completa, p. ej. www.example.com",
		ErrURLScheme: "Solo se admiten URL http / https.",
	},
	"pt": {
		Title:    "Criador de atalho de URL",
		Subtitle: "Digite uma URL · Escolha um ícone · Crie um atalho na área de trabalho",

		LabelURL:       "URL",
		LabelIcon:      "Ícone",
		Browse:         "Procurar…",
		LabelOpenWith:  "Abrir com",
		DefaultBrowser: "Navegador padrão (recomendado)",
		HintBrowser:    "Se o navegador escolhido não estiver instalado,\no navegador padrão será usado",
		LabelName:      "Nome",
		HintName:       "Deixe em branco para gerar automaticamente",
		LabelPreview:   "Prévia do ícone",
		Desc:           "O atalho será salvo na área de trabalho.\nDê um clique duplo para abrir a URL no navegador escolhido.\nO ícone é gerado a partir da imagem selecionada.",
		Generate:       "Criar atalho",
		StatusInit:     "Digite uma URL, escolha um ícone e clique em «Criar atalho».",
		LabelLang:      "Idioma",

		DlgPickTitle: "Escolher imagem do ícone",
		FilterImages: "Arquivos de imagem",
		FilterAll:    "Todos os arquivos",

		ErrImageProcess:  "Não é possível processar esta imagem: ",
		ErrImageInvalid:  "Esta imagem não pode ser usada como ícone. Tente outra.",
		StatusIconPicked: "Ícone selecionado: ",
		ErrNeedURL:       "Digite primeiro a URL.",
		ErrNeedIcon:      "Clique primeiro em «Procurar…» para escolher uma imagem de ícone.",
		ErrIconsDir:      "Não é possível criar a pasta de ícones: ",
		ErrIconProcess:   "Falha ao processar o ícone: ",

		BrowserFallbackFmt: "%s não está disponível; usando o navegador padrão.\n\n",
		ShortcutDescFmt:    "Abrir %s",
		SystemDefault:      "navegador padrão",

		ErrCreateFailed:   "Falha: ",
		ErrCreateShortcut: "Falha ao criar o atalho:\n",
		StatusCreated:     "Criado: ",
		InfoTitle:         "Concluído",
		InfoCreatedFmt:    "Atalho criado na área de trabalho:\n%s",
		InfoOpenWithFmt:   "\n\nAbrir com: %s",
		InfoTail:          "\nDê um clique duplo no atalho para abrir a URL.",

		ErrWindowClass:  "Falha ao registrar a classe da janela",
		ErrWindowCreate: "Falha ao criar a janela",

		ErrURLEmpty:  "Digite primeiro a URL.",
		ErrURLFormat: "O formato da URL é inválido.",
		ErrURLHost:   "Digite uma URL completa, por exemplo www.example.com",
		ErrURLScheme: "Somente URLs http / https são suportadas.",
	},
	"ru": {
		Title:    "Создание ярлыка URL",
		Subtitle: "Введите URL · Выберите значок · Создайте ярлык на рабочем столе",

		LabelURL:       "URL",
		LabelIcon:      "Значок",
		Browse:         "Обзор…",
		LabelOpenWith:  "Открыть с помощью",
		DefaultBrowser: "Браузер по умолчанию (рекомендуется)",
		HintBrowser:    "Если выбранный браузер не установлен,\nбудет использован браузер по умолчанию",
		LabelName:      "Имя",
		HintName:       "Оставьте пустым для автогенерации",
		LabelPreview:   "Предпросмотр значка",
		Desc:           "Ярлык будет сохранён на рабочем столе.\nДважды щёлкните, чтобы открыть URL в выбранном браузере.\nЗначок создаётся из выбранного изображения.",
		Generate:       "Создать ярлык",
		StatusInit:     "Введите URL, выберите значок и нажмите «Создать ярлык».",
		LabelLang:      "Язык",

		DlgPickTitle: "Выберите изображение значка",
		FilterImages: "Файлы изображений",
		FilterAll:    "Все файлы",

		ErrImageProcess:  "Не удалось обработать изображение: ",
		ErrImageInvalid:  "Это изображение нельзя использовать как значок. Попробуйте другое.",
		StatusIconPicked: "Значок выбран: ",
		ErrNeedURL:       "Сначала введите URL.",
		ErrNeedIcon:      "Сначала нажмите «Обзор…» и выберите изображение значка.",
		ErrIconsDir:      "Не удалось создать папку значков: ",
		ErrIconProcess:   "Ошибка обработки значка: ",

		BrowserFallbackFmt: "%s недоступен, используется браузер по умолчанию.\n\n",
		ShortcutDescFmt:    "Открыть %s",
		SystemDefault:      "браузер по умолчанию",

		ErrCreateFailed:   "Ошибка: ",
		ErrCreateShortcut: "Не удалось создать ярлык:\n",
		StatusCreated:     "Создано: ",
		InfoTitle:         "Готово",
		InfoCreatedFmt:    "Ярлык создан на рабочем столе:\n%s",
		InfoOpenWithFmt:   "\n\nОткрыть с помощью: %s",
		InfoTail:          "\nДважды щёлкните ярлык, чтобы открыть URL.",

		ErrWindowClass:  "Не удалось зарегистрировать класс окна",
		ErrWindowCreate: "Не удалось создать окно",

		ErrURLEmpty:  "Сначала введите URL.",
		ErrURLFormat: "Неверный формат URL.",
		ErrURLHost:   "Введите полный URL, например www.example.com",
		ErrURLScheme: "Поддерживаются только URL http / https.",
	},
}

func tr(code string) uiStrings {
	if s, ok := translations[code]; ok {
		return s
	}
	return translations["en"]
}

func fontFor(code string) string {
	for _, l := range languages {
		if l.code == code {
			return l.family
		}
	}
	return "Segoe UI"
}

func languageIndex(code string) int {
	for i, l := range languages {
		if l.code == code {
			return i
		}
	}
	return -1
}

// detectSystemLanguage 根据 Windows 界面语言返回支持的界面语言代码。
func detectSystemLanguage() string {
	id := uint32(syscallN(procGetUserDefaultUILanguage))
	primary := id & 0x3FF
	sub := (id >> 10) & 0x3F
	switch primary {
	case 0x04: // 中文
		if sub == 0x02 || sub == 0x03 {
			return "zh-TW"
		}
		return "zh-CN"
	case 0x09:
		return "en"
	case 0x11:
		return "ja"
	case 0x12:
		return "ko"
	case 0x07:
		return "de"
	case 0x0C:
		return "fr"
	case 0x0A:
		return "es"
	case 0x16:
		return "pt"
	case 0x19:
		return "ru"
	}
	return "en"
}

func isImageError(err error) bool {
	return errors.Is(err, core.ErrImageUnsupported) ||
		errors.Is(err, core.ErrImageEmpty) ||
		errors.Is(err, core.ErrImageSizeInvalid) ||
		errors.Is(err, core.ErrImageTooLarge)
}

// errText 把 core 返回的哨兵错误转换成当前界面语言的文案。
func errText(err error, s uiStrings) string {
	switch {
	case errors.Is(err, core.ErrURLRequired):
		return s.ErrURLEmpty
	case errors.Is(err, core.ErrURLInvalid):
		return s.ErrURLFormat
	case errors.Is(err, core.ErrURLHostRequired):
		return s.ErrURLHost
	case errors.Is(err, core.ErrURLSchemeUnsupported):
		return s.ErrURLScheme
	case isImageError(err):
		return s.ErrImageInvalid
	}
	return err.Error()
}

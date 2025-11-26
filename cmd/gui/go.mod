module gui

go 1.25.4

require (
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.52.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero v1.15.0
	github.com/vegidio/go-sak v0.0.0-20251125122114-8b01823c4a11
	github.com/vegidio/umd v0.0.0-20250918022752-c66c1259c887
	github.com/wailsapp/wails/v2 v2.10.2
)

require (
	github.com/PuerkitoBio/goquery v1.11.0 // indirect
	github.com/Velocidex/json v0.0.0-20220224052537-92f3c0326e5a // indirect
	github.com/Velocidex/ordereddict v0.0.0-20250626035939-2f7f022fc719 // indirect
	github.com/Velocidex/yaml/v2 v2.2.8 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/bep/debounce v1.2.1 // indirect
	github.com/browserutils/kooky v0.2.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/go-sqlite/sqlite3 v0.0.0-20180313105335-53dd8e640ee7 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gonuts/binary v0.2.0 // indirect
	github.com/google/go-github/v74 v74.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect
	github.com/keybase/go-keychain v0.0.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/labstack/echo/v4 v4.13.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.1 // indirect
	github.com/leaanthony/gosod v1.0.4 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/leaanthony/u v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/tkrajina/go-reflector v0.5.8 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/go-webview2 v1.0.21 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	github.com/zalando/go-keyring v0.2.6 // indirect
	github.com/zeebo/assert v1.3.0 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	www.velocidex.com/golang/go-ese v0.2.0 // indirect
)

// replace github.com/wailsapp/wails/v2 v2.9.2 => /Users/vegidio/go/pkg/mod

// Local shared code
replace github.com/vegidio/shared => ../shared

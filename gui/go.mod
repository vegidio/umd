module gui

go 1.23.2

require (
	github.com/google/uuid v1.6.0
	github.com/samber/lo v1.51.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero v1.14.0
	github.com/vegidio/shared v0.0.0-00010101000000-000000000000
	github.com/vegidio/umd-lib v0.0.0-20250814050101-d1f8efa57498
	github.com/wailsapp/wails/v2 v2.10.1
)

require (
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/bep/debounce v1.2.1 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/dromara/dongle v1.0.1 // indirect
	github.com/emmansun/gmsm v0.30.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/go-rod/rod v0.116.2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/labstack/echo/v4 v4.13.3 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.1 // indirect
	github.com/leaanthony/gosod v1.0.4 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/leaanthony/u v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mixpanel/mixpanel-go v1.2.1 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/tkrajina/go-reflector v0.5.8 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/go-webview2 v1.0.21 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	github.com/ysmood/fetchup v0.3.0 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.41.0 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	github.com/ysmood/leakless v0.9.0 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

// replace github.com/wailsapp/wails/v2 v2.9.2 => /Users/vegidio/go/pkg/mod

// Local shared code
replace github.com/vegidio/shared => ../shared

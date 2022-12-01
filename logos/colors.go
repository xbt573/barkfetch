package logos

var Colors = map[string]string{
	"c0":     "\x1b[38;5;0m",
	"c1":     "\x1b[38;5;1m",
	"c2":     "\x1b[38;5;2m",
	"c3":     "\x1b[38;5;3m",
	"c4":     "\x1b[38;5;4m",
	"c5":     "\x1b[38;5;5m",
	"c6":     "\x1b[38;5;6m",
	"c7":     "\x1b[38;5;7m",
	"c8":     "\x1b[38;5;8m",
	"c9":     "\x1b[38;5;9m",
	"c10":    "\x1b[38;5;10m",
	"c11":    "\x1b[38;5;11m",
	"c12":    "\x1b[38;5;12m",
	"c13":    "\x1b[38;5;13m",
	"c14":    "\x1b[38;5;14m",
	"c15":    "\x1b[38;5;15m",
	"creset": "\x1b[0m",
}

func ColorExpand(color string) string {
	return Colors[color]
}

func init() {
	logo, err := GuessLogo()
	if err != nil {
		panic(err)
	}

	Colors["caccent"] = Colors[logo.AccentColor]
}

package routines

import (
	"barkfetch/info"
)

var Routines = map[string]func() any{
	"logo": func() any {
		logo, err := info.GetLogo()
		if err != nil {
			panic(err)
		}

		return logo
	},
	"userline": func() any {
		userline, err := info.GetUserline()
		if err != nil {
			panic(err)
		}

		return userline
	},
	"userunderline": func() any {
		underline, err := info.GetUserUnderline()
		if err != nil {
			panic(err)
		}

		return underline
	},
	"kernel": func() any {
		kernel, err := info.GetKernel()
		if err != nil {
			panic(err)
		}

		return kernel
	},
	"uptime": func() any {
		uptime, err := info.GetUptime()
		if err != nil {
			panic(err)
		}

		return uptime
	},
	"shell": func() any {
		return info.GetShell()
	},
	"memory": func() any {
		mem, err := info.GetMemory()
		if err != nil {
			panic(err)
		}

		return mem
	},
}

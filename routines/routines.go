package routines

import "barkfetch/info"
import "barkfetch/console"

var Routines = map[string]func(*console.Console){
	"logo": func(c *console.Console) {
		lastline := c.GetLastLine()
		c.PushAt(0, c.GetLastLine(), info.GetLogo())
		c.SetOffset(info.GetLogoMaxLength() + 1)

		c.SetLastLine(lastline)
	},
	"userline": func(c *console.Console) {
		userline, err := info.GetUserline()
		if err != nil {
			panic(err)
		}

		c.PushAt(0, c.GetLastLine(), userline)
	},
	"userunderline": func(c *console.Console) {
		underline, err := info.GetUserUnderline()
		if err != nil {
			panic(err)
		}

		c.PushAt(0, c.GetLastLine(), underline)
	},
	"kernel": func(c *console.Console) {
		kernel, err := info.GetKernel()
		if err != nil {
			panic(err)
		}

		c.PushAt(0, c.GetLastLine(), kernel)
	},
	"uptime": func(c *console.Console) {
		uptime, err := info.GetUptime()
		if err != nil {
			panic(err)
		}

		c.PushAt(0, c.GetLastLine(), uptime)
	},
	"shell": func(c *console.Console) {
		c.PushAt(0, c.GetLastLine(), info.GetShell())
	},
	"memory": func(c *console.Console) {
		mem, err := info.GetMemory()
		if err != nil {
			panic(err)
		}

		c.PushAt(0, c.GetLastLine(), mem)
	},
}

package ascii

import (
	"github.com/fatih/color"
)

// Banner returns the ASCII art banner for go-starter
func Banner() string {
	// Create color functions
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	
	banner := `
` + cyan(`  ██████╗  ██████╗ ██╗      █████╗ ███╗   ██╗ ██████╗`) + `
` + cyan(`  ██╔════╝ ██╔═══██╗██║     ██╔══██╗████╗  ██║██╔════╝`) + `
` + blue(`  ██║  ███╗██║   ██║██║     ███████║██╔██╗ ██║██║  ███╗`) + `
` + blue(`  ██║   ██║██║   ██║██║     ██╔══██║██║╚██╗██║██║   ██║`) + `
` + green(`  ╚██████╔╝╚██████╔╝███████╗██║  ██║██║ ╚████║╚██████╔╝`) + `
` + green(`   ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝`) + `
` + cyan(``) + `
` + cyan(`  ███████╗████████╗ █████╗ ██████╗ ████████╗███████╗██████╗`) + `
` + blue(`  ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗`) + `
` + blue(`  ███████╗   ██║   ███████║██████╔╝   ██║   █████╗  ██████╔╝`) + `
` + green(`  ╚════██║   ██║   ██╔══██║██╔══██╗   ██║   ██╔══╝  ██╔══██╗`) + `
` + green(`  ███████║   ██║   ██║  ██║██║  ██║   ██║   ███████╗██║  ██║`) + `
` + green(`  ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝`) + `

`
	return banner
}

// Logo returns a smaller ASCII logo for go-starter
func Logo() string {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	return cyan(`
 _____ _____ 
|   __|     |
|  |  |  |  |
|_____|_____|  STARTER
`)
}

// Gopher returns a small ASCII gopher
func Gopher() string {
	blue := color.New(color.FgBlue).SprintFunc()
	return blue(`
        ,_---~~~~~----._         
 _,,_,*^____      _____*g*"*,   
/ __/ /'     ^.  /      \ ^@q   f 
[  @f | @))    |  | @))   l  0 _/  
 \ /   \~____ / __ \_____/    \   
  |           _l__l_           I   
  }          [______]           I  
  ]            | | |            |  
  ]             ~ ~             |  
  |                            |   
   |                           |
`)
}
/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	. "fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	ptime "github.com/yaa110/go-persian-calendar"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	usage := `jalalicli is a CLI frontend for https://github.com/yaa110/go-persian-calendar
	
	Usage:
	  jalalicli today [--jalali-format=<jalali-format>]
	  jalalicli tojalali [--gregorian-format=<gregorian-format> --jalali-format=<jalali-format>] <date>
	  jalalicli togregorian [--gregorian-format=<gregorian-format>] <date>
	  jalalicli -h | --help
	
	  togregorian's input should be in a "yyyy/MM/dd" format.

	Options:
	  -j --jalali-format=<jalali-format>  Jalali format (see the readme of the backend).
	  -g --gregorian-format=<gregorian-format>  Gregorian format (go style). [Default: 2006/01/02]
	  -h --help  Show this screen.`

	debug := os.Getenv("DEBUGME") != ""
	arguments, _ := docopt.ParseDoc(usage)
	if debug {
		log.Println(os.Args)
		log.Println(arguments)
	}

	todayMode := arguments["today"].(bool)
	var jalaliFormat string = "yyyy/MM/dd"
	if arguments["--jalali-format"] != nil {
		jalaliFormat = arguments["--jalali-format"].(string)
	}
	layout := arguments["--gregorian-format"].(string)

	tojalaliMode := arguments["tojalali"].(bool)
	togregorianMode := arguments["togregorian"].(bool)

	var date string
	if arguments["<date>"] != nil {
		date = arguments["<date>"].(string)
	}

	if todayMode {
		// Get a new instance of ptime.Time representing the current time
		pt := ptime.Now(ptime.Iran())
		Println(pt.Format(jalaliFormat))
	} else if tojalaliMode {
		t, err := time.Parse(layout, date)
		check(err)
		pt := ptime.New(t)
		Println(pt.Format(jalaliFormat))
	} else if togregorianMode {
		dateParts := strings.Split(date, "/")
		year, err := strconv.Atoi(dateParts[0])
		check(err)
		month, err := strconv.Atoi(dateParts[1])
		check(err)
		day, err := strconv.Atoi(dateParts[2])
		check(err)
		var pt ptime.Time = ptime.Date(year, ptime.Month(month), day, 12, 59, 59, 0, ptime.Iran())
		t := pt.Time()
		Println(t.Format(layout))
	}
}

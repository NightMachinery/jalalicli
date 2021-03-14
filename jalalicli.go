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
	  jalalicli today [--jalali-format=<jalali-format> --inc-year=<years> --inc-month=<monthss> --inc-day=<days>]
	  jalalicli tojalali [--gregorian-format=<gregorian-format> --jalali-format=<jalali-format> --inc-year=<years> --inc-month=<months> --inc-day=<days>] <date>
	  jalalicli togregorian [--gregorian-format=<gregorian-format> --inc-year=<years> --inc-month=<months> --inc-day=<days>] [<date>]
	  jalalicli -h | --help
	
	  togregorian's input should be in a "yyyy/MM/dd" format.
      Date increments are always done in Jalali. Negative numbers are supported.

	Options:
	  -j --jalali-format=<jalali-format>  Jalali format (see the readme of the backend).
	  -g --gregorian-format=<gregorian-format>  Gregorian format (go style). [Default: 2006/01/02]
      -y --inc-year=<years>  Increment output's year by specified amount.
      -m --inc-month=<months>  Increment output's month by specified amount.
      -d --inc-day=<days>  Increment output's day by specified amount.
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

	var incYear int = 0
	var err error = nil
	if arguments["--inc-year"] != nil {
		incYear, err = strconv.Atoi(arguments["--inc-year"].(string))
		check(err)
	}
	var incMonth int = 0
	if arguments["--inc-month"] != nil {
		incMonth, err = strconv.Atoi(arguments["--inc-month"].(string))
		check(err)
	}
	var incDay int = 0
	if arguments["--inc-day"] != nil {
		incDay, err = strconv.Atoi(arguments["--inc-day"].(string))
		check(err)
	}

	var date string = ""
	if arguments["<date>"] != nil {
		date = arguments["<date>"].(string)
	}

	if todayMode {
		// Get a new instance of ptime.Time representing the current time
		pt := ptime.Now(ptime.Iran())
		pt = pt.AddDate(incYear, incMonth, incDay)
		if jalaliFormat != "unix" {
			Println(pt.Format(jalaliFormat))
		} else {
			Println(pt.Unix())
		}
	} else if tojalaliMode {
		var t time.Time
		if layout != "unix" {
			_t, err := time.Parse(layout, date)
			check(err)
			t = _t
		} else {
			i, err := strconv.ParseInt(date, 10, 64)
			check(err)
			t = time.Unix(i, 0)
		}
		// t = t.AddDate(incYear, incMonth, incDay)
		pt := ptime.New(t)
		pt = pt.AddDate(incYear, incMonth, incDay)
		/// different from t.AddDate:
		// jalalicli tojalali 2001/09/11 --inc-month 1 => 1380/07/20
		// jalalicli tojalali 2001/09/11 --inc-month 1 => 1380/07/19 # t.AddDate
		///
		if jalaliFormat != "unix" {
			Println(pt.Format(jalaliFormat))
		} else {
			Println(t.Unix())
		}
	} else if togregorianMode {
		var pt ptime.Time
		if date != "" {
			dateParts := strings.Split(date, "/")
			year, err := strconv.Atoi(dateParts[0])
			check(err)
			month, err := strconv.Atoi(dateParts[1])
			check(err)
			day, err := strconv.Atoi(dateParts[2])
			check(err)
			pt = ptime.Date(year, ptime.Month(month), day, 12, 59, 59, 0, ptime.Iran())
		} else {
			pt = ptime.Now(ptime.Iran())
		}
		pt = pt.AddDate(incYear, incMonth, incDay)
		t := pt.Time()
		// t = t.AddDate(incYear, incMonth, incDay)
		if layout != "unix" {
			Println(t.Format(layout))
		} else {
			Println(t.Unix())
		}
	}
}

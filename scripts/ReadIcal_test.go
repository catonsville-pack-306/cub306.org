package main

/** ReadIcal test file, this will test most of the functions in ReadIcal.go
Most important is the full test at the bottom TestWork() which will run a sample
ical file through the work code and generate markdown responses.

Other tests are unit level tests of key functions in the app.
*/

import (
    "fmt"
    "testing"
    "time"
    //"bufio"
    "strings"
    "reflect"
)

const SAMPLE_ICAL_1 = `BEGIN:VCALENDAR
PRODID:-//Google Inc//Google Calendar 70.9054//EN
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:PUBLISH
X-WR-CALNAME:Cub Scout Pack 306
X-WR-TIMEZONE:America/New_York
X-WR-CALDESC:Cub Scout Pack 306 of Catonsville\, MD
BEGIN:VTIMEZONE
TZID:America/New_York
X-LIC-LOCATION:America/New_York
BEGIN:DAYLIGHT
TZOFFSETFROM:-0500
TZOFFSETTO:-0400
TZNAME:EDT
DTSTART:19700308T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
END:DAYLIGHT
BEGIN:STANDARD
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
TZNAME:EST
DTSTART:19701101T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
DTSTART:20210602T230000Z
DTEND:20210603T000000Z
DTSTAMP:20210530T183248Z
UID:24m08jgh1kscp879k328mebro1@google.com
CREATED:20210530T174621Z
DESCRIPTION:There is no meeting on this day
LAST-MODIFIED:20210530T174621Z
LOCATION:
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:No Meeting
TRANSP:OPAQUE
END:VEVENT
END:VCALENDAR`

func Test (t *testing.T) {
	
	json := `[
	{"Keys": ["hike", "walk"],
	 "Icon": "<i class=\"fas fa-leaf\"></i>"},
	{"Keys": ["camp", "camping"],
	 "Icon": "<i class=\"fas fa-campground\"></i>"},
	{"Keys": ["meeting", "mappings"],
	 "Icon": "<i class=\"fas fa-user-friends\"></i>"},
	{"Keys": [], "Icon": "<i class=\"far fa-calendar-alt\"></i>"}
]`
	
	tester := func (given, expected, msg string) {
		key_data := icon_keys_from_string(json)
		actual := findSummaryIcons(given, key_data)
		
		if expected != actual {
			t.Error(msg, actual)
		}
	}
	
	tester("this is my summary about camping",
		"<i class=\"fas fa-campground\"></i>",
		"Camping found")
	
	tester("this is generic",
		"<i class=\"far fa-calendar-alt\"></i>",
		"Nothing to find, default")
}

func TestEventDayAfter (t *testing.T) {
	start := "20210529T000000Z"
	actualEnd := "20210530T000000Z"
	
	tester := func(start, stop string, expected bool, msg string) {
		event := CalEvent{start: start,
    	    end: stop,
        	actualEnd: stop,
	        summary: "summary",
    	    location: "location",
        	description: "description"}
		actual := eventDayAfter(event, 0)

		if actual!=expected {
			t.Error(msg, "test:", stop, "is not a day after", start, "\n\n")
		}
	}
	
	tester(start, actualEnd, true, "Good")
	tester("20220102", "20220102", false, "Same day")
	tester("20220102", "20220103", true, "Next day")
}

func TestDayAfter(t *testing.T) {
	tester := func(before, after int, expected bool, msg string) {
	    expected1 := time.Date(2021, 5, before, 0, 0, 0, 0, time.UTC)
    	expected2 := time.Date(2021, 5, after, 0, 0, 0, 0, time.UTC)
    	
		actual := dayAfter(expected1, expected2)
		if expected != actual {
			report := fmt.Sprintf("%s test: %d is not a day from %d.",
				msg, after, before)
			t.Error(report)
		}
	}

	tester(3, 3, false, "same day")
	tester(3, 4, true, "good")
	tester(4, 3, false, "end before begin")
}


func TestDateStrToObj2 (t *testing.T) {
	tester := func (given string, expected time.Time, msg string) {
		actual := dateStrToObj(given, -4)
		if expected != actual {
			report := fmt.Sprintf("%s test: %s != %s.", msg, expected, actual) 
			t.Error(report)
		}
	}
	
	given1 := "20210530T080000Z"
    expected1 := time.Date(2021, 5, 30, 4, 0, 0, 0, time.UTC)
	
	tester(given1, expected1, "Normal")
}

func TestZeroOutTime (t *testing.T) {
	tester := func (given, expected time.Time, msg string) {
		actual := zeroOutTimeFromDate(given)
		if expected != actual {
			report := fmt.Sprintf("%s test: %s != %s.", msg, expected, actual) 
			t.Error(report)
		}
	}
	
	given1 := time.Date(2021, 5, 30, 14, 32, 48, 0, time.UTC)
    expected1 := time.Date(2021, 5, 30, 0, 0, 0, 0, time.UTC)
	
	tester(given1, expected1, "Normal")
}

func TestDateStrToHuman (t *testing.T) {
	tester := func (given, expected string, offset int, msg string) {
		actual := dateStrToHuman(given, offset)
		if expected != actual {
			report := fmt.Sprintf("%s test: %s != %s.", msg, expected, actual) 
			t.Error(report)
		}
	}
	
	tester("20100408T083246Z", "April 08, 2010: 04:32 AM", -4, "Full Date & Time")
	tester("20100408T040000Z", "April 08, 2010: 12:00 AM", -4, "Just Date")
	tester("", "", -4, "Empty")
}

func TestDateStrToString (t *testing.T) {
	tester := func (given, expected string, offset int, msg string) {
		actual := dateStrToString(given, offset)
		if expected != actual {
			report := fmt.Sprintf("%s test: %s != %s.", msg, expected, actual) 
			t.Error(report)
		}
	}
	
	tester("20100408T083246Z", "2010-04-08 04:32 AM", -4, "Full Date & Time")
	tester("20100408T040000Z", "2010-04-08", -4, "Just Date")
	tester("", "empty", -4, "Empty")
}

func TestBetweenDates(t *testing.T) {
	early := time.Date(2010, 4, 8, 14, 32, 48, 0, time.UTC)
	middle := time.Date(2015, 1, 1, 14, 32, 48, 0, time.UTC)
	late := time.Date(2020, 9, 10, 14, 32, 48, 0, time.UTC)

	tester := func(past, current, future time.Time, expected bool, msg string) {
		actual := betweenDates(past, current, future)
		if expected != actual {
			report := fmt.Sprintf("Failed %s test, %s < %s < %s",
				msg, past, current, future)
			t.Error(report)
		}
	}
	
	tester(early, middle, late, true, "normal")
	tester(early, late, middle, false, "middle,late swapped")
	tester(middle, early, late, false, "normal-early swapped")
	tester(middle, late, early, false, "early last")
	tester(late, early, middle, false, "late before all")
	tester(late, middle, early, false, "early-late swapped")
}

func TestOutputFileName(t *testing.T) {
	tester := func (prefix string, start time.Time, end time.Time,
			expected string, msg string) {
		actual := outputFileName(prefix, start, end)
		if expected!=actual {
			given := fmt.Sprintf("(%s-%s-%s)", prefix, start, end)
			report := fmt.Sprintf("Failed %s test, from %s: \n%s != %s", msg,
				given, expected, actual)
			t.Error(report)
		}
	}
	
	d1 := time.Date(2020, 9, 10, 14, 32, 48, 0, time.UTC)
	d2 := time.Date(2020, 9, 11, 14, 32, 48, 0, time.UTC)
	d3 := time.Date(2020, 9, 10, 0, 0, 0, 0, time.UTC)
	var d4 time.Time
	
	tester("cal", d1, d2, "cal-2020-09-10-to-2020-09-11.md", "Basic two date")
	tester("cal", d1, d1, "cal-2020-09-10-to-2020-09-10.md", "Same, with time")
	tester("cal", d3, d3, "cal-2020-09-10-to-2020-09-10.md", "Same, no time")
	tester("swp", d2, d1, "swp-2020-09-11-to-2020-09-10.md", "Swapped dates")
	tester("blk", d1, d4, "blk-2020-09-10-to-0001-01-01.md", "blank date")
}

func TestStrToNumber(t *testing.T) {
	tester := func (given string, fallback int, expected int, msg string){
		actual := strToNumber(given, fallback)
		if expected!=actual {
			report := fmt.Sprintf("Failed %s test, %s!=%d", msg, given, 
				expected)
			t.Error(report)
		}
	}
	
	tester("-1", -1, -1, "negative one")
	tester("1", 0, 1, "one")
	tester("one", -1, -1, "word one")
	tester("3.1459", -1, -1, "pi")
	tester("256", -1, 256, "some large number")
}

func TestNumberToMonth(t *testing.T) {
	tester := func(given int, expected time.Month, msg string) {
		actual := numberToMonth(given)
		if expected!=actual {
			t.Error(fmt.Sprintf("Number to month failed for %s : %v != %v(%d)",
				msg, expected, actual, given))
		}
	}
	
	tester(0, time.December, "Low Month")
	tester(1, time.January, "January")
	tester(6, time.June, "June")
	tester(12, time.December, "December")
	tester(255, time.December, "High Month")
}

func TestAppData(t *testing.T) {
	//now := time.Now()
	now := time.Date(2020, 9, 10, 14, 32, 48, 0, time.UTC)
	actual := InitApp(now)


	tester := func(actual AppData, name string, exp interface{}){
		a := reflect.ValueOf(&actual).Elem()
		
		switch tv := a.FieldByName(name).Interface().(type) {
			case string:
				act := string(a.FieldByName(name).String())
				if exp != act {
					t.Errorf(fmt.Sprintf("Bad %v %v!=%v", name, exp, act))
				}
			case int64, int:
				act := int(a.FieldByName(name).Int())
				if exp != act {
					t.Errorf(fmt.Sprintf("Bad %v %v=%v", name, exp, act))
				}
			case bool:
				act := bool(a.FieldByName(name).Bool())
				if exp != act {
					t.Errorf(fmt.Sprintf("Bad %v %v=%v", name, exp, act))
				}
			default:
				t.Error(fmt.Sprintf("unknown type tested: %s %t", name, tv))
		}
	}
	
	tester(actual, "FilePrefix", "calendar")
	tester(actual, "IconKeyFile", "icon_keys.json")
	tester(actual, "OutputPath", "out")
	tester(actual, "TimeZone", -4)
	tester(actual, "Date", "2020-9-10")
	tester(actual, "Limit", 0)
	tester(actual, "Markdown", true)
	tester(actual, "OutMonths", -1)
	tester(actual, "OutDays", 0)
	tester(actual, "AfterMonths", 0)
	tester(actual, "AfterDays", 1)
}

func TestDateStrToObj(t *testing.T) {
    expected := time.Date(2021, 5, 30, 14, 32, 48, 0, time.UTC)
    actual := dateStrToObj("20210530T183248Z", -4)
    if expected != actual {
        t.Errorf("dateStrToObj does not match: [%s] vs [%s].", expected, actual)
    }
}

func TestWork(t *testing.T) {
	expected := `
## <i class="fas fa-user-friends"></i> 2021-06-02 07:00 PM to 2021-06-02 08:00 PM - No Meeting

There is no meeting on this day

* When: June 02, 2021: 07:00 PM

`
    reader := strings.NewReader(SAMPLE_ICAL_1)

	date_of_test := time.Date(2021,5,30,14,32,48,0,time.UTC)
    app_data := InitApp(date_of_test)
    app_data.OutputPath = "/tmp"
    app_data.TimeZone = -4
    app_data.Markdown = false

    work(reader, date_of_test, app_data)
    actual := readFile("/tmp/calendar-2021-05-02-to-2021-06-03.md")
    
    if expected != actual {
        t.Errorf("not matching")
    }
}
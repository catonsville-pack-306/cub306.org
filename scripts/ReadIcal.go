package main

//go get github.com/arran4/golang-ical

import ("fmt"
    "bufio"
    "os"
    //"io"
    "time"
    //"io/ioutil"
    //"bytes"
    "github.com/arran4/golang-ical"
    "strings"
    "flag"
    "text/template"
    //"syscall"
    //"unsafe"
    "io"
    )

const FILE_TEMPLATE string = `
## {{ .Icon }} {{ .When }} - {{ .Summary }}
{{ if .Description}}
{{ .Description}}
{{end}}
{{if .WhenHuman}}* When: {{ .WhenHuman }}{{end}}
{{if .Where}}* {{ .Where }}{{end}}
`

const (
    ICON_LEAF = "<i class=\"fas fa-leaf\"></i>"
    ICON_CAMP = "<i class=\"fas fa-campground\"></i>"
    ICON_FRIENDS = "<i class=\"fas fa-user-friends\"></i>"
    ICON_CALENDAR = "<i class=\"far fa-calendar-alt\"></i>"
    ICON_MAP = "<i class=\"far fa-map\"></i>"
    //ICON_SMILE "<i class=\"far fa-smile\"></i>"
    //ICON_FLAG_US "<i class=\"fas fa-flag-usa\"></i>"
    //ICON_PASSPORT "<i class=\"fas fa-passport\"></i>"
)

/** 4 kids of Dates
Start date
Start Date & Time
Stand Date and End Date
Stand Date & Time and End Date & Time
*/
type cubevent struct {
    start string
    end string
    actualEnd string
    summary string
    location string
    description string
}

type template_data struct {
    Icon string
    When string
    Ending string
    Summary string
    Description string
    WhenHuman string
    Where string
}

func summaryIcon(summary string,  flags []string, icon string) (bool, string) {
    for _, value := range flags {
        if strings.Contains(strings.ToLower(summary), strings.ToLower(value)) {
            return true, icon
        }
    }
    return false, ""
}

func findSummaryIcons(summary string) string {
    found, icon := summaryIcon (summary, []string{"hike", "walk"}, ICON_LEAF)
    if found {return icon}
    
    found, icon = summaryIcon (summary, []string{"camping", "camp"}, ICON_CAMP)
    if found {return icon}
    
    found, icon = summaryIcon (summary, []string{"meeting", "meetings"}, ICON_FRIENDS)
    if found {return icon}
    
    return ICON_CALENDAR
}

/** convert a calender event to a cubevent
*/
func eventToCubEvent (event *ics.VEvent) cubevent {
    start := valueIfExists(event, ics.ComponentPropertyDtStart)
    actualEnd := valueIfExists(event, ics.ComponentPropertyDtEnd)
    summary := valueIfExists(event, ics.ComponentPropertySummary)
    location := valueIfExists(event, ics.ComponentPropertyLocation)
    description := valueIfExists(event, ics.ComponentPropertyDescription)
    
    if len(actualEnd)<1 {actualEnd = start}
    
    //if len(actualEnd)>8 {fmt.Println(start, " - ", actualEnd)}

    return cubevent{start: start,
        end: actualEnd,
        actualEnd: actualEnd,
        summary: summary,
        location: location,
        description: description}
}

func eventDayAfter(event cubevent, offset int) bool {
    start := dateStrToObj(event.start, offset)
    end := dateStrToObj(event.end, offset)
    return dayAfter(start, end)
}

func dayAfter(start time.Time, end time.Time) bool {
    return start.Hour()==0 && start.Minute()==0 && start.Second()==0 &&
        end.Hour()==0 && end.Minute()==0 && end.Second()==0 &&
        start.Day()==(end.Day()-1)
}

/**

*/
func handleEvent(event cubevent, offset int) string {
    temp, err := template.New("event").Parse(FILE_TEMPLATE)
    if err!=nil {
        fmt.Println (err)
    } else {
        i := findSummaryIcons(event.summary)
        data := template_data{
            Icon: i,
            Summary: event.summary,
            Description: event.description,
            Where: event.location,
        }
    
        if event.start == event.end && !eventDayAfter(event, offset) {
            data.When = dateStrToString(event.start, offset)[:10]
        } else {
            data.When = fmt.Sprintf("%s to %s", dateStrToString(event.start, offset),
                dateStrToString(event.end, offset))
        }
        data.WhenHuman = dateStrToHuman(event.start, offset)
        
        buf := new(strings.Builder)
        err = temp.Execute(buf, data)
        if err!=nil {
            fmt.Println (err)
        }
        return buf.String()
    }
    return ""
}

func valueIfExists(event *ics.VEvent, propName ics.ComponentProperty) string {
    prop := event.GetProperty(propName)
    if prop!=nil {
        return prop.Value
    }
    return ""
}

func dateStrToObj(raw string, offset int) time.Time {
    ans, _ := time.Parse("20060102T150405Z", raw)
    if ans.IsZero() {
        ans, _ = time.Parse("20060102T150405Z", raw)
    }
    ans = ans.Add(time.Hour * time.Duration(offset))
    if ans.IsZero() {
        ans, _ = time.Parse("20060102", raw)
    }
    return ans
}

func zeroOutTimeFromDate(datetime time.Time) time.Time {
    datetime = datetime.Add(time.Duration(-datetime.Hour()) * time.Hour)
    datetime = datetime.Add(time.Duration(-datetime.Minute()) * time.Minute)
    datetime = datetime.Add(time.Duration(-datetime.Second()) * time.Second)
    datetime = datetime.Add(time.Duration(-datetime.Nanosecond()) * time.Nanosecond)
    return datetime
}

func dateStrToHuman(raw string, offset int) string {
    if len(raw) < 1 {return ""}
    ans := dateStrToObj(raw, offset)
    if ans.IsZero() {return ""}
    return ans.Format("January 02, 2006: 03:04 PM")
}
func dateStrToString(raw string, offset int) string {
    if len(raw) < 1 {return "empty"}
    ans := dateStrToObj(raw, offset)
    if ans.IsZero() {return "zero"}
    formated_date := ans.Format("2006-01-02 03:04 PM")
    if ans.Hour()==0 && ans.Minute()==0 && ans.Second()==0 {
        formated_date = formated_date[:10]
    }
    return formated_date
}

func betweenDates(past time.Time, current time.Time, future time.Time) bool {
    return current.After(past) && current.Before(future)
}

func writeFileByDate(path string, start time.Time, end time.Time, content string) {
    template := "calendar-%s-to-%s.md"
    startString := fmt.Sprintf("%04d-%02d-%02d", start.Year(), start.Month(), start.Day())
    endString := fmt.Sprintf("%04d-%02d-%02d", end.Year(), end.Month(), end.Day())
    name := fmt.Sprintf(template, startString, endString)
    writeFile(path + "/" + name, content)
}

func writeFile(file string, content string) {
    //d1 := []byte(content)
    //err := ioutil.WriteFile(file, d1, 0644)
    fmt.Printf("write " + file + " with\n" + content + "\n********************\n")
}

func work(reader io.Reader, timezone int, outPath string, today time.Time) {
    cal1, err2 := ics.ParseCalendar(reader)  //configure this
    if err2!=nil {
        fmt.Println ("Error:")
        fmt.Println (err2==nil)
    } else {
        for _, comp := range cal1.Events() {
            event := eventToCubEvent(comp)
            dtStart := dateStrToObj(event.start, timezone)
            dtEnd := dateStrToObj(event.end, timezone)
            if dtEnd.IsZero() { dtEnd = dtStart}

            dStart := zeroOutTimeFromDate(dtStart)
            dEnd := zeroOutTimeFromDate(dtEnd)
            monthBefore := dStart.AddDate(0, -1, 0)             //configure this
            dayAfter := dEnd.AddDate(0, 0, 1)                   //configure this

            if !betweenDates(monthBefore, today, dayAfter) {continue}

            eventContent := handleEvent(event, timezone)
            if len(outPath) > 0 {
                writeFileByDate(outPath, monthBefore, dayAfter, eventContent);
            } else {
                fmt.Println (eventContent)
            }
        }
    }
}

func main() {
    //args := os.Args
    
    outPath := flag.String("out", "out", "directory to write output to")
    timezone := flag.Int("timezone", -4, "timezone offset")
    
    flag.Parse()

    /*
    b, err1 := ioutil.ReadFile("../public/events/events.ics")
    if err1 != nil { fmt.Print(err1) }
    content := string(b)
    fmt.Println (content[:500])
    cal1, err2 := ics.ParseCalendar(strings.NewReader(content))
    */

    today := zeroOutTimeFromDate(time.Now())
    reader := bufio.NewReader(os.Stdin)

    work(reader, *timezone, *outPath, today)


}
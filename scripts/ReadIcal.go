package main

/* *****************************************************************************

File: ReadIcal.go
Author: thomas.cherry@gmail.com
Abstract:
Read in an ical file from standard in and output markdown or partial HTML files
for each event in a time block, normally one month in the future, one day in the
past. These output files are then used in a web page to provide timely event
content for a web site.
***************************************************************************** */

// be sure to load these dependencies
// go get github.com/arran4/golang-ical
// go get gitlab.com/golang-commonmark/markdown

import (
    "fmt"
    "bufio"
    "encoding/json"
    "errors"
    "flag"
    "github.com/arran4/golang-ical"
    "gitlab.com/golang-commonmark/markdown"
    "io"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
    "text/template"
    "time"
    )

/* ************************************************************************** */
// MARK: Consts

const APP_BY = "ReadIcal by Thomas.Cherry@gmail.com"
const APP_VERSION = "1.0.0"

/**
Raw JSON containing a a default list of maps containing a list of keywords and
an HTML snip-it representing the topic of a message. Used to apply a custom icon
for a category of topics. By default, a scout/camping theme is used.
*/
const DEFAULT_MAPPINGS string = `[
    {"Keys": ["hike", "walk"],
     "Icon": "<i class=\"fas fa-leaf\"></i>"},
    {"Keys": ["camp", "camping"],
     "Icon": "<i class=\"fas fa-campground\"></i>"},
    {"Keys": ["meeting", "mappings"],
     "Icon": "<i class=\"fas fa-user-friends\"></i>"},
    {"Keys": ["happy", "smile"],
     "Icon": "<i class=\"far fa-smile\"></i>"},
    {"Keys": ["flag", "usa"],
     "Icon": "<i class=\"fas fa-flag-usa\"></i>"},
    {"Keys": ["world", "travel"],
     "Icon": "<i class=\"fas fa-passport\"></i>"},
    {"Keys": [],
     "Icon": "<i class=\"far fa-calendar-alt\"></i>"}
]
`

/**
Default markdown template for each event
*/
var FILE_TEMPLATE string = `
## {{ .Icon }} {{ .WhenHuman }} - {{ .Summary }}
{{ if .Description}}
{{ .Description}}
{{end}}
{{if .When}}* When: {{ .When }}{{end}}
{{if .Where}}* {{ .Where }}{{end}}
`

/* ************************************************************************** */
// MARK: - Data Types

type ExitCode int
const (
    EXIT_NORMAL ExitCode = iota //0
    EXIT_VERSION
    EXIT_SHOW_ICONS
)
func (ec ExitCode) int() int {
    return int(ec)
}


/** App Data holding the current state of the application */
type AppData struct {
    Markdown bool
    FilePrefix string
    IconKeyFile string
    OutputPath string
    TimeZone int
    Date string
    Template string
    IconKeyData []IconKeys
    Limit int
    Verbose bool
    OutMonths, OutDays, AfterMonths, AfterDays int
}

/**
Event data representing data of interest from one parsed ical event.
4 kids of Dates
Start date
Start Date & Time
Stand Date and End Date
Stand Date & Time and End Date & Time
*/
type CalEvent struct {
    start string
    end string
    actualEnd string
    summary string
    location string
    description string
}

/**
Data to be passed to the template engine to be swapped out with template
variables
*/
type template_data struct {
    Icon string
    When string
    Ending string
    Summary string
    Description string
    WhenHuman string
    Where string
}

/** Icon configuration data from a JSON configuration like DEFAULT_MAPPINGS */
type IconKeys struct {
    Keys []string
    Icon string
}

/**
Filter out lines from a stream, used to remove apple specific tags from a
calendar as it confuses the library
*/
func filterStream (input io.Reader, output *io.PipeWriter, matchList []string) {
    scanner := bufio.NewScanner(input)
    defer output.Close()
    for scanner.Scan() {
        line := scanner.Text()
        skip := false
        for _, match := range matchList {
            if len(match)>0 && strings.Contains(line, match) {
                skip = true
                break
            }
        }
        if skip {continue}
        output.Write([]byte(line + "\n"))
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}

/* ************************************************************************** */
// MARK: - dependency functions

/** wrapper function to the markdown package
not tested
@param input markdown
@return html
*/
func MarkdownToHTML(input string) string {
    md := markdown.New(markdown.HTML(true))
    output := md.RenderToString([]byte(input))
    return output
}

/**
Not Tested
Return a value from a ics property
@param event ics event to search for property
@param propName field in event to search for
*/
func valueIfExists(event *ics.VEvent, propName ics.ComponentProperty) string {
    prop := event.GetProperty(propName)
    if prop!=nil {
        return prop.Value
    }
    return ""
}

/**
convert a calender event to a CalEvent
Not Tested
*/
func eventToCalEvent (event *ics.VEvent) CalEvent {
    start := valueIfExists(event, ics.ComponentPropertyDtStart)
    actualEnd := valueIfExists(event, ics.ComponentPropertyDtEnd)
    summary := valueIfExists(event, ics.ComponentPropertySummary)
    location := valueIfExists(event, ics.ComponentPropertyLocation)
    description := valueIfExists(event, ics.ComponentPropertyDescription)

    if len(actualEnd)<1 {actualEnd = start}

    return CalEvent{start: start,
        end: actualEnd,
        actualEnd: actualEnd,
        summary: summary,
        location: location,
        description: description}
}

/* ************************************************************************** */
// MARK: - functions

/**
Check to see what icon can be made from a summary
@param summary text to search
@param flags keyword to search for
@param icon text to use when keyword is found in summary
@return true if found, false otherwise ; icon text - do we need return 2?
*/
//write a test
func summaryIcon(summary string,  flags []string, icon string) (bool) {
    for _, value := range flags {
        if strings.Contains(strings.ToLower(summary), strings.ToLower(value)) {
            return true
        }
    }
    return false
}

/**
tested through other tests
*/
func icon_keys_from_string(data string) []IconKeys {
    var icon_keys []IconKeys
    err := json.Unmarshal([]byte(data), &icon_keys)
    if err!=nil {
        icon_keys = []IconKeys{}
        msg := fmt.Sprintf ("Icon key file: %s.\n", err.Error())
        os.Stderr.WriteString(msg)
    }
    return icon_keys
}

/**
load the icon file json into memory and return the resulting structure
Not Tested
@param icon_file path to a json file holding a list of Keys and Icon
*/
func load_icon_keys(icon_file string) []IconKeys {
    var icon_keys []IconKeys
    data := DEFAULT_MAPPINGS
    if testFile(icon_file) {
        file_data := readFile(icon_file)
        if len(file_data)>0 {
            data = file_data
        }
    }
    icon_keys = icon_keys_from_string(data)
    return icon_keys
}

/**
Look through all the posible icons and try to assign one based on the summary
@param summary text to search through
@param icon_file path to the icon file
*/
func findSummaryIcons(summary string, key_data []IconKeys) string {
    default_icon := ""
    for _, k := range key_data {
        if 0==len(k.Keys) {
            //if no keys are given then treat the icon as a default
            default_icon = k.Icon
        } else {
            if summaryIcon (summary, k.Keys, k.Icon) {
                return k.Icon
            }
        }
    }
    return default_icon
}

func eventDayAfter(event CalEvent, offset int) bool {
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
Process one event from the ical file
Not tested
@param CalEvent calendar event
@param app_data application configuration information
@return formated event
*/
func handleEvent(event CalEvent, app_data AppData) string {
    offset := app_data.TimeZone
    temp, err := template.New("event").Parse(FILE_TEMPLATE)
    if err!=nil {
        os.Stderr.WriteString(err.Error() + "\n")
    } else {
        icon := findSummaryIcons(event.summary, app_data.IconKeyData)
        data := template_data{
            Icon: icon,
            Summary: event.summary,
            Description: event.description,
            Where: event.location,
        }

        //Calculate When
        if event.start == event.end && !eventDayAfter(event, offset) {
            // an all day event
            data.When = dateStrToString(event.start, offset)[:10]
        } else {
            // multi day event
            data.When=fmt.Sprintf("%s to %s",
                dateStrToString(event.start, offset),
                dateStrToString(event.end, offset))
        }
        //create Human readable When
        data.WhenHuman = dateStrToHuman(event.start, offset)

        buf := new(strings.Builder)
        err = temp.Execute(buf, data)
        if err!=nil {
            os.Stderr.WriteString(err.Error() + "\n")
        }
        output := buf.String()

        if app_data.Markdown {
            output = MarkdownToHTML(output)
        }
        return output
    }
    return ""
}

/**
convert a date string and timezone offset to a Time Object
@param raw date in a string in ISO format
@param offset hours offset from UTC
@return time.Time object
*/
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

/**
Remove set time fields hour, minute, second, and nanoseconds to zero
*/
func zeroOutTimeFromDate(datetime time.Time) time.Time {
    datetime = datetime.Add(time.Duration(-datetime.Hour()) * time.Hour)
    datetime = datetime.Add(time.Duration(-datetime.Minute()) * time.Minute)
    datetime = datetime.Add(time.Duration(-datetime.Second()) * time.Second)
    datetime = datetime.Add(time.Duration(-datetime.Nanosecond()) * time.Nanosecond)
    return datetime
}

/**
Convert raw string containing an iso date to a pretty view of the data
@return Month day, year: hour:min AM/PM
*/
func dateStrToHuman(raw string, offset int) string {
    if len(raw) < 1 {return ""}
    ans := dateStrToObj(raw, offset)
    if ans.IsZero() {return ""}
    return ans.Format("January 02, 2006: 03:04 PM")
}

/**
Convert raw string containing an iso date to an iso style output
@return year-month-day hour:minute AM/PM
*/
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

/** Util Function to test if a time is between two other dates */
func betweenDates(past time.Time, current time.Time, future time.Time) bool {
    return current.After(past) && current.Before(future)
}

/**
Delete all files with a prefix from a given directory
Not tested!
*/
func clearOutDataFiles(path string, app_data AppData) {
    open_path, err1 := os.Open(path)
    if err1!=nil {
        os.Stderr.WriteString(err1.Error() + "\n")
        return
    }
    allFiles, err2 := open_path.Readdir(0)
    if err2!=nil {
        os.Stderr.WriteString(err2.Error() + "\n")
        return
    }
    for f := range(allFiles) {
        file := allFiles[f]
        fileName := file.Name()
        if strings.HasPrefix(fileName, app_data.FilePrefix + "-") {
            filePath := path + "/" + fileName
            err3 := os.Remove(filePath)
            if err3 != nil {
                os.Stderr.WriteString(err3.Error() + "\n")
            } else {
                fmt.Println("The file has been deleted: ", filePath)
            }
        }
    }
}

/**
generate the markdown save name
tested
*/
func outputFileName(prefix string, start time.Time, end time.Time) string {
    template := "%s-%s-to-%s.md"
    startString := fmt.Sprintf("%04d-%02d-%02d", start.Year(), start.Month(), start.Day())
    endString := fmt.Sprintf("%04d-%02d-%02d", end.Year(), end.Month(), end.Day())
    name := fmt.Sprintf(template, prefix, startString, endString)
    return name
}

/**
Not tested!
*/
func writeFileByDate(path string, prefix string, start time.Time, end time.Time, content string) {
    name := outputFileName(prefix, start, end)
    writeFile(path + "/" + name, content)
}

/**
Util/ Function to Test if a file exists
NOT tested!
*/
func testFile(file string) bool {
    _, err := os.Stat(file)
    return ! errors.Is(err, os.ErrNotExist)
}

/**
Util Function Write a file
Not tested!
*/
func writeFile(file string, content string) {
    d1 := []byte(content)
    err := ioutil.WriteFile(file, d1, 0644)
    if err!=nil {
        fmt.Printf("write " + file + " with\n" + content + "\n**************\n")
        os.Stderr.WriteString(err.Error() + "\n")
    }
}

/**
Util Function Read a file
Not tested!
@param full path to read
@return empty string on error, file contents otherwise
*/
func readFile(file_path string) string {
    content, err := ioutil.ReadFile(file_path)
    if err != nil {
        os.Stderr.WriteString(err.Error() + "\n")
        return ""
    }
    return string(content)
}

/**
Util Function to convert a string to a numberToMonth
@param raw string posible containing a number
@param fallback value to use if raw fails to contain a number
@return a number
*/
func strToNumber (raw string, fallback int) int{
    converted, err := strconv.Atoi(raw)
    if err != nil {
        converted = fallback
        msg := fmt.Sprintf("%s using %d.\n", err.Error(), fallback)
        os.Stderr.WriteString(msg)
    }
    return converted
}

/** Util Function Convert a 1-12 number to a Golang time.Month object */
func numberToMonth(month int) time.Month {
    var month_obj time.Month
    switch month {
        case 1: month_obj = time.January
        case 2: month_obj = time.February
        case 3: month_obj = time.March
        case 4: month_obj = time.April
        case 5: month_obj = time.May
        case 6: month_obj = time.June
        case 7: month_obj = time.July
        case 8: month_obj = time.August
        case 9: month_obj = time.September
        case 10: month_obj = time.October
        case 11: month_obj = time.November
        default: month_obj = time.December
    }
    return month_obj
}

/* ************************************************************************** */
// MARK: - Application

/** create a default app_data structure with default values */
func InitApp(now time.Time) AppData {
    year := now.Year()
    month := int(now.Month())
    day := now.Day()

    date := fmt.Sprintf("%d-%d-%d", year, month, day)

    app_data := AppData{
        FilePrefix: "calendar",
        IconKeyFile: "icon_keys.json",
        OutputPath: "out",
        TimeZone: -4,
        Date: date,
        Template: FILE_TEMPLATE,
        Limit: 0,
        Markdown: true,
        Verbose: false,

        OutMonths: -1,
        OutDays: 0,
        AfterMonths: 0,
        AfterDays: 1}
    return app_data
}

/**
Tested!
do the work for this command, load ical files through a stream and processes it
for a given day
@param reader - stream of ical file
@param today - date to base results around
@param app_data - application configurations
*/
func work(reader io.Reader, today time.Time, app_data AppData) {
    cal1, err := ics.ParseCalendar(reader)  //configure this
    if err!=nil {
        os.Stderr.WriteString(err.Error() + "\n")
    } else {
        app_data.IconKeyData = load_icon_keys(app_data.IconKeyFile)
        event_count := 0
        for _, comp := range cal1.Events() {
            event := eventToCalEvent(comp)
            dtStart := dateStrToObj(event.start, app_data.TimeZone)
            dtEnd := dateStrToObj(event.end, app_data.TimeZone)
            if dtEnd.IsZero() { dtEnd = dtStart}

            dStart := zeroOutTimeFromDate(dtStart)
            dEnd := zeroOutTimeFromDate(dtEnd)
            monthBefore := dStart.AddDate(0,app_data.OutMonths,app_data.OutDays)
            dayAfter := dEnd.AddDate(0,app_data.AfterMonths,app_data.AfterDays)

            if !betweenDates(monthBefore, today, dayAfter) {continue}
            //at least one event has made it this far
            event_count = event_count + 1
            eventContent := handleEvent(event, app_data)
            if len(app_data.OutputPath) > 0 {
                if event_count==1 {
                    clearOutDataFiles(app_data.OutputPath, app_data)
                }
                writeFileByDate(app_data.OutputPath,
                    app_data.FilePrefix,
                    monthBefore,
                    dayAfter,
                    eventContent);
            } else {
                fmt.Println (eventContent)
            }
            if 0 < app_data.Limit && app_data.Limit <= event_count {
                //Limit of 0 means unlimited events
                break
            }
        }
    }
}

/** command line interface */
func main() {
    //overwrite the usage function
    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), APP_BY + "n\n")
        fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()
    }

    //process command line arguments
    version := flag.Bool("version", false, "Report application version information")
    verbose := flag.Bool("verbose", false, "send more text to err")

    outPath := flag.String("out", "out", "directory to write output to")
    timezone := flag.Int("timezone", -4, "timezone offset")
    date := flag.String("date", "", "date to use")
    template := flag.String("template", "", "Event Template File")
    icon_file := flag.String("icons",
        "icon_keys.json","File defining icon maps, - to output internal file")
    limit := flag.Int("limit", 0, "max number of events to return")
    no_html := flag.Bool("no-html", false, "Do not convert markdown to HTML")

    out_months := flag.Int("out-months", -1, "months to look forward")
    out_days := flag.Int("out-days", 0, "months to look forward")
    after_months := flag.Int("after-months",0,"months to keep event after date")
    after_days := flag.Int("after-days", 1, "days to keep event after date")

    flag.Parse()

    app_data := InitApp(time.Now())

    if *version {
        fmt.Printf("%s\nVesion %s\n", APP_BY, APP_VERSION)
        os.Exit(EXIT_VERSION.int())
    }

    //process verbose
    if *verbose {app_data.Verbose = *verbose}

    //process no_html
    if *no_html {app_data.Markdown = false}

    //process limit
    if 0 <= *limit {app_data.Limit = *limit}

    //process out and after times
    if -12 < *out_months && *out_months < 12 {app_data.OutMonths = *out_months}
    if -365 < *out_days && *out_days < 365 {app_data.OutDays = *out_days}
    if -12 < *after_months && *after_months < 12 {app_data.AfterMonths = *after_months}
    if -365 < *after_days && *after_days < 365 {app_data.AfterDays = *after_days}

    //process icon key file
    if 0<len(*icon_file) {
        if *icon_file == "-" || *icon_file == "internal" {
            fmt.Println (DEFAULT_MAPPINGS)
            os.Exit(EXIT_SHOW_ICONS.int())
        } else if *icon_file == "ignore" {
            app_data.IconKeyFile = ""
            if app_data.Verbose {
                fmt.Fprintln(os.Stderr, "ignoring external file, using internal")
            }
        } else {
            app_data.IconKeyFile = *icon_file
        }
    }

    //process out path
    if 0<len(*outPath) { app_data.OutputPath = *outPath}

    //process timezone
    if -12 < *timezone && *timezone < 12 {app_data.TimeZone = *timezone}

    //process template
    if 1<len(*template) {
        raw_template := readFile(*template)
        if 0<len(raw_template) {
            app_data.Template = raw_template
            FILE_TEMPLATE = raw_template //delete
        }
    }

    //process date
    today := zeroOutTimeFromDate(time.Now())
    if *date != "" {
        app_data.Date = *date
        date_parts := strings.Split(*date, "-")
        if len(date_parts)>2 {
            //today := time.Date(2018, time.August, 30, 0, 0, 0, 0, time.UTC)
            year := strToNumber(date_parts[0], time.Now().Year())
            month_num := strToNumber(date_parts[1], int(time.Now().Month()))
            day := strToNumber(date_parts[2], time.Now().Day())
            month := numberToMonth(month_num)
            today = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
        }
    }
    
    if *verbose {
        fmt.Println (app_data)
        fmt.Println (today)
    }

    // go to work
        
    //golang-ical seams to choke on tags from Apple, so filter them out
    reader, writer := io.Pipe()
    go filterStream(os.Stdin, writer, []string{"X-APPLE-"})
    
    work(reader, today, app_data)
}

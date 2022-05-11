# built in
require 'date'      # iso date method
require 'erb'       # template pages and inject pages can be ERB files
require 'json'

# 

# Types of start and end dates:
        #if key == "DTSTART" or 
            #key == "DTSTART;VALUE=DATE"
            #or key == "DTSTART;TZID=America/New_York"
        #if key=="DTEND" 
            #or key=="DTEND;VALUE=DATE"
            #or key=="DTEND;TZID=America/New_York"
#

# ******************************************************************************

class String
    def on_list? (list)
        list.any? { |word| self.include?(word) }
    end
end

# ******************************************************************************

class ReadIcal
    ICON_LEAF = '<i class="fas fa-leaf"></i>'
    ICON_CAMP = '<i class="fas fa-campground"></i>'
    ICON_FRIENDS = '<i class="fas fa-user-friends"></i>'
    ICON_CALENDAR = '<i class="far fa-calendar-alt"></i>'
    ICON_MAP = '<i class="far fa-map"></i>'
    #ICON_SMILE = '<i class="far fa-smile"></i>'
    #ICON_FLAG_US = '<i class="fas fa-flag-usa"></i>'
    #ICON_PASSPORT = '<i class="fas fa-passport"></i>'

    EVENTS_FILE = 'public/events/events.ics'

    FILE_TEMPLATE = %q{
## <%=icon%> <%=date%> - <%=summary%>

<%=description%>

* When: <%=when_str%>
* <%=where%>
}

    def initialize()
        @data = {}
        @events_file = nil
        if @events_file.nil?
            @events_file = Dir.glob("#{Dir.pwd}/../**/events.ics")
            unless @events_file.nil?
                @events_file = @events_file.first
            end
        end
    end

    def read
        lines = []
        input = File.exists?(@events_file) ? File.open(@events_file) : STDIN
        puts "did not find #{@events_file}" unless File.exists?(@events_file)
        input.read.split("\r\n").each do |text|
            unless text.empty? then
                if text.start_with?(' ') then
                    # this is a continued line
                    before = lines.pop()
                    text[0] = ''
                    text = "#{before}#{text}"
                end
                #un-escape things
                text.gsub! '\\;', ';'
                text.gsub! '\\,', ','
                text.gsub! '\,', ','
                text.gsub! '\\N', '\n'
                text.gsub! '\\n', '\n'
                text.gsub! '\\\\', '\\'
                lines << text
            end
        end
        rescue Exception => error
            print error
        end
        calender_hash lines
    end

    # recursivly convert an array of ical tags to a hash table. Lines are
    # removed from input as processing progresses.
    # @param lines array of ical statments
    # @return hash table
    def calender_hash (lines)
        obj = {}
        until lines.empty?
            line = lines.delete_at(0) # remove and return first item
            parts = line.split(":", 2)
            key = parts[0]
            value = parts[1]
            #value.strip! unless value.nil? # forgot why this was done

            if key == "BEGIN"
                # start a new block
                name = value
                sub_object = calender_hash lines
                if obj[ name ].nil?
                    # first time here
                    obj[ name ] = sub_object
                elsif obj[ name ].kind_of?(Array)
                    # nth time here
                    obj[ name ] << sub_object
                else
                    #second time here, convert to an array
                    original_item = obj[name]
                    obj[ name ] = []
                    obj[ name ] << original_item
                    obj[ name ] << sub_object
                end
            elsif key == "END"
                #end the current block
                return obj
            else
                # add element to block
                obj[ key ] = value
            end
        end
        return obj
    end
    

    def read_from_stdin
        @data = read
    end

    def sorted_events
        events = @data["VCALENDAR"]["VEVENT"].sort_by do |k|
            ret = ''
            k.each do |key, value|
                if key.start_with? 'DTSTART'
                    ret = value
                    break
                end
            end
            ret
        end
        events
    end
    
    #accepts a block
    def events
        #events = sorted_events
        events = @data["VCALENDAR"]["VEVENT"]

        #extract details from events
        events.each do |item|
            details = {}
            item.each do |key, value|
                if key.start_with? "DTSTART"
                    details['start'] = value
                    unless value.nil?
                        if value.end_with?("Z")
                            details['start_date'] = DateTime.parse(value).new_offset('-4')
                        else
                            details['start_date'] = Date.parse value unless value.nil?
                        end
                    end
                end
                if key.start_with? "DTEND"
                    details['stop'] = value
                    unless value.nil?
                        if value.end_with?("Z")
                            details['end_date'] = DateTime.parse(value).new_offset('-4')
                        #else
                            #details['stop_date'] = Date.parse value unless value.nil?
                        end
                    end
                end
                if key == "SUMMARY"
                    details['summary'] = value
                end
                if key == "DESCRIPTION"
                    details['description'] = value
                end
                if key == "LOCATION"
                    details['location'] = value
                end
            end

            start_date = details["start_date"]
            if details["end_date"].nil?
                end_date = start_date
            else
                end_date = details['end_date']
            end
            show_date = start_date - 30
            hide_date = end_date + 1

            details['show_date'] = show_date
            details['hide_date'] = hide_date
            details['end_date'] = end_date

            if show_date <= Date.today and Date.today <= hide_date
                yield details
            end
        end
    end

    def format_date (raw)
        raw.strftime("%m-%d %I:%M%P")
    end
end

# ******************************************************************************

if __FILE__ == $PROGRAM_NAME
    destination = ARGV[0]

    #remove previouse events
    Dir.glob("#{destination}/calendar-*-to-*.md").each{ |f| File.delete(f)}
    
    reader = ReadIcal.new
    reader.read_from_stdin
    reader.events do |details|
        date = ""
        if details['start_date'] == details['end_date']
            if details['start_date'].to_s.include? "T"
                date = "#{details['start_date'].strftime("%Y-%m-%d %I:%M%P")}"
            else
                date = "#{details['start_date']}"
            end
        else
            if details['start_date'].to_s.include? 'T'
                date = "#{details['start_date'].strftime("%Y-%m-%d %I:%M%P")} to #{details['end_date'].strftime("%I:%M%P")}\n"
            else
                date = "#{details['start_date']} to #{details['end_date']}\n"
            end
        end
        
        base_file_name = "#{details['show_date'].to_s[0..9]}-to-#{details['hide_date'].to_s[0..9]}.md"
        file_name = "#{destination}/calendar-#{base_file_name}"
        alt_file_name = "#{destination}/#{base_file_name}"
        if File.file?(alt_file_name)
            puts "skipping #{file_name} because #{alt_file_name} exists"
            next
        end

        File.open(file_name, 'w') do |file|
            puts "Writing #{file_name}"
            
            unless details['summary'].nil?
                summary = details['summary']
                sum_down = summary.downcase
                if sum_down.on_list? ['hike','walk']
                    icon = ReadIcal::ICON_LEAF
                elsif sum_down.on_list? ['camping','camp']
                    icon = ReadIcal::ICON_CAMP
                elsif sum_down.on_list? ['meeting', 'meetings']
                    icon = ReadIcal::ICON_FRIENDS
                else
                    icon = ReadIcal::ICON_CALENDAR
                end
            else
                summary = ''
            end
            description = details['description'] unless details['description'].nil?
            
            description.gsub! '\n', '
'
            
            w = DateTime.parse details['start']
            #W = DateTime.parse details['stop']
            when_str = "#{w.strftime '%B %d'}"
            #when_str "#{w.strftime '%B %d - %I:%M'} to W.strftime '%I:%M'"

            unless details['location'].nil?
                if details['location'].start_with? "http"
                    where = "#{ReadIcal::ICON_MAP} [Where](#{details['location']})"
                else
                    where = details['location']
                end
            end
            
            template = ERB.new ReadIcal::FILE_TEMPLATE
            file.write template.result(binding)
            
            puts "*"*80
        end
    end
end

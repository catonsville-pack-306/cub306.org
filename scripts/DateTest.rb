# date test

require 'date'      # iso date method

raw='20200213T000000Z'

result = DateTime.parse(raw).new_offset('-5')


puts result.strftime("%m-%d %I:%M%P")
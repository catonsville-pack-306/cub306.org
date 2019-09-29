require 'rqrcode'

url = ARGV[0] ||= "https://cub306.org"
puts url
qrcode = RQRCode::QRCode.new(url)

# NOTE: showing with default options specified explicitly
png = qrcode.as_png(
  bit_depth: 1,
  border_modules: 4,
  color_mode: ChunkyPNG::COLOR_GRAYSCALE,
  color: 'black',
  file: nil,
  fill: 'white',
  module_px_size: 6,
  resize_exactly_to: false,
  resize_gte_to: false,
  size: 120
)

IO.write("out.png", png.to_s)

#puts svg
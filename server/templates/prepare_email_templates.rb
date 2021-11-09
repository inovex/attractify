require 'premailer'

def prepare(file, data)
  premailer = Premailer.new(data, warn_level: Premailer::Warnings::SAFE, with_html_string: true, css: ['mails/style.css'], escape_url_attributes: false)

  html = premailer.to_inline_css
  html.gsub!(/%7B/, '{')
  html.gsub!(/%7D/, '}')
  html.gsub!(/%24/, '$')

  File.open(file.gsub(/\.src/, ''), 'w') do |fh|
    fh.write(html)
  end

  premailer.warnings.each do |w|
    puts("#{w[:message]} (#{w[:level]}) may not render properly in #{w[:clients]}")
  end
end

Dir.glob('mails/**/*.src').each do |f|
  data = File.open(f).read
  puts f
  prepare(f, data)
end

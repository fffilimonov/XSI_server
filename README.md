# XSI_server
XSI_server

setsid ./XSI_server config.cfg &>log.log

kill -HUP `ps axf | awk '$0~"XSI_server"&&$0!~"awk"{print $1}'`

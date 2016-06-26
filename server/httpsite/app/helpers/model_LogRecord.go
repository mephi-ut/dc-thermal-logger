package helpers

import "devel.mephi.ru/dc-thermal-logger/server/httpsite/app/models"

var LogRecord models.LogRecord
func init() { LogRecord.Init(&LogRecord) }


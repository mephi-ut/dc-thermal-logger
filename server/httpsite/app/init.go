package app

import (
	  "time"
	  "os"
	  "github.com/revel/revel"
	  "database/sql"
	  "gopkg.in/reform.v1"
	  "gopkg.in/reform.v1/dialects/postgresql"
	  "devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	_ "github.com/lib/pq"
)

var DB *reform.DB

func initDB() {
	simpleDB, err := sql.Open("postgres", revel.Config.StringDefault("app.db_url", "postgres://localhost/sensors"))
	if (err != nil) {
		revel.ERROR.Printf("Cannot connect to DB: %s", err.Error())
		os.Exit(-1)
	}

	DB = reform.NewDB(simpleDB, postgresql.Dialect, reform.NewPrintfLogger(revel.TRACE.Printf))
}

func initRecordsConverted() {

	go func() {
		for ;; {
			revel.TRACE.Printf("Running converter iteration…")
			rawRecords,err := models.RawRecord.Select(DB)
			if (err != nil) {
				revel.ERROR.Printf("Converter error: %v", err.Error())
				continue
			}

			for _,rawRecord := range rawRecords {
				var err error

				historyRecords := rawRecord.ToHistoryRecords()
				for _,historyRecord := range historyRecords {
					historyRecordFilter := historyRecord
					historyRecordFilter.RawValue       = 0
					historyRecordFilter.ConvertedValue = 0
					historyRecordOld,er := models.HistoryRecord.First(DB, historyRecordFilter)
					err = er

					if err != nil {
						if err == reform.ErrNoRows {
							revel.TRACE.Printf("historyRecord.Insert(DB)")
							historyRecord.Counter = 1
							err = historyRecord.Insert(DB)
						}
					} else {
						historyRecordOld.Merge(historyRecord)
						revel.TRACE.Printf("historyRecordOld.Update(DB)")
						err = historyRecordOld.Update(DB)
					}

					if err != nil {
						break
					}
				}
				if err != nil {
					revel.ERROR.Printf("Converter error: %v", err.Error())
					continue
				}
				err = rawRecord.Delete(DB)
				if err != nil {
					revel.ERROR.Printf("Converter error: %v", err.Error())
				}
			}

			time.Sleep(time.Second)
		}
	}()
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(initDB)
	revel.OnAppStart(initRecordsConverted)
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

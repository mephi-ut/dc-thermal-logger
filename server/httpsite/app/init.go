package app

import (
	  "fmt"
	  "time"
	  "reflect"
	  "math"
	  "os"
	  "github.com/revel/revel"
	  "strings"
	  "encoding/base64"
	  "net/http"
	  "net/url"
	  "database/sql"
	  "gopkg.in/reform.v1"
	//"gopkg.in/reform.v1/dialects/postgresql"
	  "gopkg.in/reform.v1/dialects/mysql"
	  "devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app/models"
	//_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	  "gopkg.in/cas.v1"
)

var DB *reform.DB
var CasClient *cas.Client

func initDB() {
	//simpleDB, err := sql.Open("postgres", revel.Config.StringDefault("app.db_url", "postgres://localhost/sensors"))
	simpleDB, err := sql.Open("mysql", revel.Config.StringDefault("app.db_url", "mysql://localhost/sensors"))
	if (err != nil) {
		revel.ERROR.Printf("Cannot connect to DB: %s", err.Error())
		os.Exit(-1)
	}

	//DB = reform.NewDB(simpleDB, postgresql.Dialect, reform.NewPrintfLogger(revel.TRACE.Printf))
	DB = reform.NewDB(simpleDB, mysql.Dialect, reform.NewPrintfLogger(revel.TRACE.Printf))

	models.HistoryRecord.SetDefaultDB(DB)
	models.RawRecord.SetDefaultDB(DB)
}

func initRecordsConverted() {
	if revel.DevMode {	// Don't corrupt data
		return
	}

	go func() {
		for ;; {
			revel.TRACE.Printf("Running converter iteration…")
			rawRecords,err := models.RawRecord.Limit(100000).Select()
			if (err != nil) {
				revel.ERROR.Printf("Converter error: %v", err.Error())
				continue
			}

			sem := make(chan bool, 32)

			for _,rawRecord := range rawRecords {
				sem <- true

				rawRecordE := models.NewRawRecordE(rawRecord);

				go func(rawRecord models.RawRecordE) {
					defer func() { <-sem }()

					var err error

					historyRecords,er := rawRecord.ToHistoryRecords()
					if er == models.ErrInvalidSensorId {
						revel.WARN.Printf("Got error \"%v\" while parsing {%v}", er.Error(), rawRecord)
						return
					}

					for _,historyRecord := range historyRecords {
						historyRecordFilter := historyRecord
						historyRecordFilter.RawValue       = 0
						historyRecordFilter.ConvertedValue = 0
						historyRecordOld,er := models.HistoryRecord.First(historyRecordFilter)
						err = er

						if err != nil {
							if err == reform.ErrNoRows {
								revel.TRACE.Printf("historyRecord.Insert()")
								historyRecord.Counter = 1
								err = historyRecord.Insert()
							}
						} else {
							if (historyRecord.ConvertedValue < models.MIN_CONVERTED_VALUE) {
								revel.INFO.Printf("skipped: historyRecordOld.Update(): SensorId: %v(%v); %v < %v: %v (%v)", historyRecord.SensorId, rawRecord.RawSensorId, historyRecord.ConvertedValue, models.MIN_CONVERTED_VALUE, historyRecord, rawRecord)
								continue
							}
							historyRecordOld.Merge(historyRecord)
							revel.TRACE.Printf("historyRecordOld.Update()")
							err = historyRecordOld.Update()
						}

						if err != nil {
							return
						}
					}
					if err != nil {
						revel.ERROR.Printf("Converter error: %v", err.Error())
						return
					}
					err = rawRecord.Delete()
					if err != nil {
						revel.ERROR.Printf("Converter error: %v", err.Error())
					}
				}(rawRecordE);
			}

			time.Sleep(time.Second)
		}
	}()
}

func initCasClient() {
	url, _ := url.Parse("https://login.mephi.ru")
	CasClient = cas.NewClient(&cas.Options{
		URL: url,
	})
}

func toFloat64(i interface{}) (float64) {
	switch v := i.(type) {
		case int64:
			return float64(v)
		case int:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return v
	}

	return math.NaN()
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
		ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	revel.OnAppStart(initDB)
	revel.OnAppStart(initRecordsConverted)
	revel.OnAppStart(initCasClient)

	revel.TemplateFuncs["dict"] = func(values ...interface{}) (map[string]interface{}, error) {	// This function is copied from http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters/18276968
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i+=2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, fmt.Errorf("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	}

	revel.TemplateFuncs["hasIndex"] = func(a interface{}, idxs ...interface{}) bool {
		v := reflect.ValueOf(a)
		t := v.Type()

		for _,idxI := range idxs {
			switch t.Kind() {
				case reflect.Slice, reflect.Array:
					idx := idxI.(int)
					if idx >= v.Len() {
						return false
					}
				case reflect.Map:
					r := v.MapIndex(reflect.ValueOf(idxI))
					if (!r.IsValid()) {
						return false
					}
					if r.Interface() == reflect.Zero(r.Type()).Interface() {
						return false
					}
			}
		}
		return true
	}

	revel.TemplateFuncs["isNil"] = func(arg interface{}) (bool) {
		switch arg.(type) {
			case nil:
				return true
		}
		v := reflect.ValueOf(arg)
		t := v.Type().Kind()
		switch (t) {
			case reflect.Ptr:
				return v.IsNil()
		}
		return false;
	}

	revel.TemplateFuncs["isEmpty"] = func(arg interface{}) (bool) {
		value := reflect.ValueOf(arg)
		if (!value.IsValid()) {
			return true
		}
		return arg == reflect.Zero(reflect.TypeOf(arg)).Interface()
	}
	revel.TemplateFuncs["flt"] = func(a, b interface{}) (bool) {
		aV := toFloat64(a)
		bV := toFloat64(b)

		return aV < bV
	}
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	/*c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")*/

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 { return false }

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil { return false }

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 { return false }

	if pair[0] == "" || pair[1] == "" {
		return false
	}

	if pair[0] == revel.Config.StringDefault("user0.login", "") && pair[1] == revel.Config.StringDefault("user0.pass", "") {
		return true
	}

	if pair[0] == revel.Config.StringDefault("user1.login", "") && pair[1] == revel.Config.StringDefault("user1.pass", "") {
		return true
	}

	if pair[0] == revel.Config.StringDefault("user2.login", "") && pair[1] == revel.Config.StringDefault("user2.pass", "") {
		return true
	}

	if pair[0] == revel.Config.StringDefault("user3.login", "") && pair[1] == revel.Config.StringDefault("user3.pass", "") {
		return true
	}

	if pair[0] == revel.Config.StringDefault("user4.login", "") && pair[1] == revel.Config.StringDefault("user4.pass", "") {
		return true
	}

	return false
}

var ActionInvoker = func(c *revel.Controller, f []revel.Filter) {
	isAuthed := checkAuth(c.Response.Out, c.Request.Request)

	if (!isAuthed) {
		h := func(w http.ResponseWriter, r *http.Request) {
			if !cas.IsAuthenticated(r) {
				CasClient.RedirectToLogin(w, r)
				return
			}

			isAuthed = true
		}

		CasClient.HandleFunc(h).ServeHTTP(c.Response.Out, c.Request.Request)
	}

	if (!isAuthed) {
		return
	}

	revel.ActionInvoker(c, f);
}


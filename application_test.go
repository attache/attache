package attache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var CURRENT_TEST *testing.T

const testSchema = `CREATE TABLE items (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name VARCHAR(40) NOT NULL,
	value INTEGER NOT NULL
);`

const testSchemaData = `INSERT INTO items (name, value) VALUES ("init", 12), ("init2", 13);`

type TestApplicationModel struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Value int    `db:"value"`
}

func (m *TestApplicationModel) Table() string { return "items" }

func (m *TestApplicationModel) Key() ([]string, []interface{}) {
	return []string{"id"}, []interface{}{m.ID}
}

func (m *TestApplicationModel) Insert() ([]string, []interface{}) {
	return []string{"name", "value"}, []interface{}{m.Name, m.Value}
}

func (m *TestApplicationModel) Update() ([]string, []interface{}) { return m.Insert() }

func (m *TestApplicationModel) AfterInsert(result sql.Result) {
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	m.ID = id
}

var _ Record = (*TestApplicationModel)(nil)

type TestApplication struct {
	BaseContext

	DefaultViews
	DefaultDB

	TestValue string
}

func (c *TestApplication) CONFIG_Views() ViewConfig {
	return ViewConfig{
		Root: "testdata/views",
	}
}

func (c *TestApplication) CONFIG_DB() DBConfig {
	return DBConfig{
		Driver: "sqlite3",
		DSN:    ":memory:",
	}
}

func (c *TestApplication) Init(_ http.ResponseWriter, _ *http.Request) {}

func writeSuccess(c Context) {
	writeString(c, "here")
}

func writeString(c Context, s string) {
	c.ResponseWriter().Write([]byte(s))
}

func (c *TestApplication) ALL_TestAll()       { writeSuccess(c) }
func (c *TestApplication) GET_TestGet()       { writeSuccess(c) }
func (c *TestApplication) PUT_TestPut()       { writeSuccess(c) }
func (c *TestApplication) POST_TestPost()     { writeSuccess(c) }
func (c *TestApplication) TRACE_TestTrace()   { writeSuccess(c) }
func (c *TestApplication) PATCH_TestPatch()   { writeSuccess(c) }
func (c *TestApplication) DELETE_TestDelete() { writeSuccess(c) }

func (c *TestApplication) GUARD_TestGuard() { c.TestValue = "set_in_guard" }
func (c *TestApplication) ALL_TestGuard()   { writeString(c, c.TestValue) }

func (c *TestApplication) GUARD_TestGuardX() { c.TestValue += " second_value" }
func (c *TestApplication) ALL_TestGuardX()   { writeString(c, c.TestValue) }

func (c *TestApplication) PROVIDE_StringVal(_ *http.Request) interface{} { return "from_provider" }
func (c *TestApplication) ALL_TestProvider(provided string)              { writeString(c, provided) }

func (c *TestApplication) PROVIDE_IntVal(_ *http.Request) interface{} { return 12 }
func (c *TestApplication) ALL_TestProvider2(provided int)             { writeString(c, fmt.Sprint(provided)) }

func (c *TestApplication) GET_TestUtils1() { Error(400) }
func (c *TestApplication) GET_TestUtils2() { ErrorMessage(400, "test_message") }
func (c *TestApplication) GET_TestUtils3() { ErrorFatal(errors.New("test_err")) }
func (c *TestApplication) GET_TestUtils4() { ErrorMessageJSON(400, "test") }

func (c *TestApplication) GET_TestViews() { RenderHTML(c, "test1") }

func (c *TestApplication) GET_TestDB() {
	conn := c.DB()

	if err := conn.Raw().Ping(); err != nil {
		CURRENT_TEST.Error("bad database:", err)
		return
	}

	item := TestApplicationModel{Name: "test", Value: 102}

	if err := conn.Insert(&item); err != nil {
		CURRENT_TEST.Error("bad insert:", err)
		return
	}

	if want := int64(3); item.ID != want {
		CURRENT_TEST.Errorf("bad insert: expected new ID %d, got %d", want, item.ID)
		return
	}

	item = TestApplicationModel{} // reset

	if err := conn.Get(&item, 3); err != nil {
		CURRENT_TEST.Error("bad fetch:", err)
		return
	}

	if want := (TestApplicationModel{3, "test", 102}); item != want {
		CURRENT_TEST.Errorf("bad fetch: expected %#+v, got %#+v", want, item)
	}

	item = TestApplicationModel{} // reset

	if err := conn.Get(&item, 2); err != nil {
		CURRENT_TEST.Error("bad fetch:", err)
		return
	}

	if want := (TestApplicationModel{2, "init2", 13}); item != want {
		CURRENT_TEST.Errorf("bad fetch: expected %#+v, got %#+v", want, item)
	}

	item = TestApplicationModel{} // reset

	if err := conn.Get(&item, 1); err != nil {
		CURRENT_TEST.Error("bad fetch:", err)
		return
	}

	if want := (TestApplicationModel{1, "init", 12}); item != want {
		CURRENT_TEST.Errorf("bad fetch: expected %#+v, got %#+v", want, item)
	}
}

func TestMain(m *testing.M) {
	devNull, _ := os.Open(os.DevNull)
	log.SetOutput(devNull)
	tmpCtx := &TestApplication{}
	app, err := Bootstrap(tmpCtx)
	if err != nil {
		log.Fatalln(err)
	}

	app.NoLogging = true

	db, err := DBFor(tmpCtx.CONFIG_DB())
	if err != nil {
		log.Fatalln(err)
	}

	// build schema
	if _, err := db.Raw().Exec(testSchema); err != nil {
		log.Fatalln(err)
	}

	// insert initial data
	if _, err := db.Raw().Exec(testSchemaData); err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{Addr: ":8080"}
	go app.RunWithServer(server)
	code := m.Run()
	server.Shutdown(context.Background())
	os.Exit(code)
}

func TestApplicationEndpoints(t *testing.T) {
	CURRENT_TEST = t

	allMethods := []string{
		"GET",
		"PUT",
		"POST",
		// "HEAD", ignored for testing
		"TRACE",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}

	type S = []string

	cases := []struct {
		url      string
		wantCode int
		wantBody string
		methods  []string
	}{
		{"/test/all", 200, "here", allMethods},
		{"/test/get", 200, "here", S{"GET"}},
		{"/test/put", 200, "here", S{"PUT"}},
		{"/test/post", 200, "here", S{"POST"}},
		{"/test/trace", 200, "here", S{"TRACE"}},
		{"/test/patch", 200, "here", S{"PATCH"}},
		{"/test/delete", 200, "here", S{"DELETE"}},

		{"/test/guard", 200, "set_in_guard", allMethods},
		{"/test/guard/x", 200, "set_in_guard second_value", allMethods},
		{"/test/provider", 200, "from_provider", allMethods},
		{"/test/provider2", 200, "12", allMethods},

		{"/test/utils1", 400, "", S{"GET"}},
		{"/test/utils2", 400, "test_message", S{"GET"}},
		{"/test/utils3", 500, "", S{"GET"}},
		{"/test/utils4", 400, "{\"error\":\"test\"}\n", S{"GET"}},

		{"/test/views", 200, "CONTENT[from_test_1]", S{"GET"}},

		{"/test/db", 200, "", S{"GET"}},
	}

	for i, c := range cases {
		good := c.methods
		bad := []string{}
		for _, m := range allMethods {
			if !contains(good, m) {
				bad = append(bad, m)
			}
		}

		for _, okMethod := range good {
			body, code, err := tryHttp(c.url, okMethod, "")
			if err != nil {
				t.Fatal("case:", i, "method:", okMethod, "error:", err)
				continue
			}

			if code != c.wantCode {
				t.Fatal("case:", i, "method:", okMethod, "expected:", c.wantCode, "got:", code)
				continue
			}

			if body != c.wantBody {
				t.Fatal("case:", i, "method:", okMethod, "expected:", c.wantBody, "got:", body)
				continue
			}
		}

		for _, badMethod := range bad {
			_, code, err := tryHttp(c.url, badMethod, "")
			if err != nil {
				t.Fatal("case:", i, "method:", badMethod, "error:", err)
				continue
			}

			if code != http.StatusMethodNotAllowed {
				t.Fatal("case:", i, "method:", badMethod, "expected:", http.StatusMethodNotAllowed, "got:", code)
				continue
			}
		}
	}
}

func tryHttp(url string, method string, data string) (string, int, error) {
	req, err := http.NewRequest(method, "http://localhost:8080"+url, strings.NewReader(data))
	if err != nil {
		return "", 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	return string(bytes), res.StatusCode, nil
}

func contains(list []string, want string) bool {
	// simple linear search
	for _, v := range list {
		if want == v {
			return true
		}
	}
	return false
}

package command

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/kohkimakimoto/hq/internal/client"
)

func testApp(t *testing.T) *cli.App {
	t.Helper()

	app := cli.NewApp()
	app.Writer = &bytes.Buffer{}
	app.ErrWriter = &bytes.Buffer{}
	app.Commands = Commands

	return app
}

func testRegisterTestClient(t *testing.T, app *cli.App, fn RoundTripFunc) {
	t.Helper()

	setClientFactory(app, func(ctx *cli.Context) *client.Client {
		c := client.New(ctx.String("address"))
		c.HttpClient = testHttpClient(t, fn)
		return c
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testHttpClient(t *testing.T, fn RoundTripFunc) *http.Client {
	t.Helper()

	return &http.Client{
		Transport: fn,
	}
}

func testTempFile(t *testing.T, b []byte) *os.File {
	t.Helper()
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(tempFile.Name(), b, 0644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	})
	return tempFile
}

package router

import (
	"strings"
	"testing"

	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
)

func mk(input string, t *testing.T) (*ctx.Context, *RoutingIterator) {
	t.Helper()
	parsed, err := p.ParseInput(input)
	if err != nil {
		t.Fatalf("ParseInput error: %v", err)
	}
	c := ctx.NewContext(parsed)
	it := NewRoutingIterator(c)
	return c, it
}

func TestRoute_SimpleEndpoint_HandlerCalledAndValuesParsed(t *testing.T) {
	c, it := mk("server --host=localhost --port=8080 --debug", t)

	r := NewRouter()
	called := false

	r.Endpoint("server").
		Description("Start server").
		StringOption("host").
		IntOption("port").
		BoolOption("debug").
		Handler(func(cc ctx.Context) error {
			host, str_err := cc.GetValueAsString("host")
			if str_err != nil {
				t.Fatal(str_err)
			}

			port, int_err := cc.GetValueAsInt("port")
			if int_err != nil {
				t.Fatal(int_err)
			}

			debug := cc.GetValueAsBool("debug")

			if host != "localhost" {
				t.Errorf("host = %q, want localhost", host)
			}
			if port != 8080 {
				t.Errorf("port = %d, want 8080", port)
			}
			if !debug {
				t.Errorf("debug = false, want true")
			}
			called = true
			return nil
		}).
		Register()

	r.Route(*c, it)

	if !called {
		t.Fatalf("handler was not called")
	}
}

func TestRoute_UnknownCommand_CallsCustomErrorHandler(t *testing.T) {
	c, it := mk("unknowncmd", t)

	r := NewRouter()

	var gotErr error
	var gotCmd string

	r.CustomErrorHandler(func(err error, cc ctx.Context) {
		gotErr = err
		gotCmd = cc.GetCommand()
	})

	r.Route(*c, it)

	if gotErr == nil {
		t.Fatalf("expected error for unknown command")
	}
	if gotCmd != "unknowncmd" {
		t.Errorf("ctx.GetCommand() = %q, want %q", gotCmd, "unknowncmd")
	}
	if msg := gotErr.Error(); !strings.Contains(strings.ToLower(msg), "unknown") &&
		!strings.Contains(strings.ToLower(msg), "not found") {
		t.Errorf("unexpected error: %v", gotErr)
	}
}

func TestRoute_RequiredOptions_ValidationError(t *testing.T) {
	c, it := mk("config --count=5 --ratio=3.14 --enable", t)
	r := NewRouter()

	var gotErr error
	r.CustomErrorHandler(func(err error, _ ctx.Context) { gotErr = err })

	r.Endpoint("config").
		RequiredString("config").
		IntOption("count").
		RequiredInt("port").
		FloatOption("ratio").
		RequiredFloat("threshold").
		BoolOption("enable").
		Handler(func(ctx.Context) error { return nil }).
		Register()

	r.Route(*c, it)

	if gotErr == nil {
		t.Fatalf("expected validation error for required options")
	}
	msg := strings.ToLower(gotErr.Error())
	wantAny := []string{"missing", "required", "--config", "--port", "--threshold"}
	ok := false
	for _, w := range wantAny {
		if strings.Contains(msg, strings.ToLower(w)) {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("unexpected error msg: %v", gotErr)
	}
}

func TestRoute_ExclusiveGroups_Conflict(t *testing.T) {
	c, it := mk("backup --local --path=/data --remote --host=example.com --user=admin --s3", t)
	r := NewRouter()

	var gotErr error
	r.CustomErrorHandler(func(err error, _ ctx.Context) { gotErr = err })

	r.Endpoint("backup").
		Description("Create backup").
		ExclusiveGroup("local", "--local").
		RequiredString("path").
		BoolOption("compress").
		StringOption("encryption").
		EndGroup().
		ExclusiveGroup("remote", "--remote").
		RequiredString("host").
		RequiredString("user").
		StringOption("key").
		IntOption("port").
		EndGroup().
		ExclusiveGroup("cloud", "--s3").
		StringOption("bucket").
		StringOption("region").
		StringOption("access-key").
		EndGroup().
		SetGroupsCanBeIgnored(false).
		Handler(func(ctx.Context) error { return nil }).
		Register()

	r.Route(*c, it)

	if gotErr == nil {
		t.Fatalf("expected conflict error when multiple exclusive groups are used")
	}
	msg := strings.ToLower(gotErr.Error())
	if !strings.Contains(msg, "group requires solitude") {
		t.Errorf("unexpected error msg: %v", gotErr)
	}
}

func TestRoute_NestedCommands_HandlerReached(t *testing.T) {
	c, it := mk("docker container run --image=nginx --name=web --port=80 --detach", t)
	r := NewRouter()

	called := false

	r.NewCmd("docker").
		NewSub("container").
		Endpoint("run").
		Description("Run a container").
		RequiredString("image").
		StringOption("name").
		IntOption("port").
		BoolOption("detach").
		Handler(func(cc ctx.Context) error {
			img, _ := cc.GetValueAsString("image")
			if img != "nginx" {
				t.Errorf("image = %q, want nginx", img)
			}
			called = true
			return nil
		}).
		Build().
		Endpoint("stop").
		Description("Stop container").
		RequiredString("container").
		BoolOption("force").
		Handler(func(ctx.Context) error { return nil }).
		Build().
		Build().
		Register()

	r.Route(*c, it)

	if !called {
		t.Fatalf("nested handler was not called")
	}
}

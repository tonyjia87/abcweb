package abcserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nullbio/lolwtf/rendering"
	"github.com/volatiletech/abcmiddleware"
	"github.com/volatiletech/abcweb/abcconfig"
)

func TestNotFound(t *testing.T) {
	t.Parallel()

	// Only run the non-compiled assets hotpath test if it can
	// find the robots.txt asset file to run against
	_, err := os.Stat(filepath.FromSlash("../public/robots.txt"))
	if err != nil {
		t.Skip("cannot find robots.txt asset to run against, so skipping NotFound test")
	}

	state := &State{}
	state.AppConfig = &abcconfig.AppConfig{
		RenderRecompile: true,
		PublicPath:      filepath.FromSlash("../public"),
	}
	state.InitLogger()
	state.Render = rendering.InitRenderer(state.AppConfig, filepath.FromSlash("../templates"))

	// test the non-compiled assets hotpath first
	r := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()

	r = r.WithContext(context.WithValue(r.Context(), abcmiddleware.CtxLoggerKey, state.Log))

	// Call the handler
	state.NotFound(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected http 200, but got http %d", w.Code)
	}
	loc, err := w.Result().Location()
	if err == nil {
		t.Error("did not expect a redirect, but got one to:", loc.String())
	}

	// Only run the compiled assets hotpath test if it can
	// find the assets/main.css asset file to run against
	_, err = os.Stat(filepath.FromSlash("../public/assets/css/main.css"))
	if err != nil {
		t.Skip("cannot find main.css asset to run against, so skipping NotFound test")
	}

	// test the compiled assets hotpath with non-manifest
	r = httptest.NewRequest("GET", "/assets/css/main.css", nil)
	w = httptest.NewRecorder()

	r = r.WithContext(context.WithValue(r.Context(), abcmiddleware.CtxLoggerKey, state.Log))

	// Call the handler
	state.NotFound(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected http 200, but got http %d", w.Code)
	}
	loc, err = w.Result().Location()
	if err == nil {
		t.Error("did not expect a redirect, but got one to:", loc.String())
	}

	// test the compiled assets hotpath with manifest
	r = httptest.NewRequest("GET", "/assets/css/main-manifestmagic.css", nil)
	w = httptest.NewRecorder()

	r = r.WithContext(context.WithValue(r.Context(), abcmiddleware.CtxLoggerKey, state.Log))

	// Set asset manifest to test manifest hotpath
	rendering.AssetsManifest = map[string]string{
		"css/main-manifestmagic.css": "css/main.css",
	}

	// Call the handler
	state.NotFound(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected http 200, but got http %d", w.Code)
	}
	loc, err = w.Result().Location()
	if err == nil {
		t.Error("did not expect a redirect, but got one to:", loc.String())
	}
}
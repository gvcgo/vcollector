package official

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gvcgo/vcollector/internal/req"
)

// gradlePageFixture mimics the markup of https://gradle.org/release-checksums/.
// Version headings used to look like "v8.14.3", but gradle.org dropped the "v"
// prefix when Gradle 9 was released, so both formats must be supported.
const gradlePageFixture = `
<html><body>
<h3 class="u-text-with-icon">
<span class="text-with-icon__icon" aria-hidden="true"></span>
<span><a id="9.0.0" href="#9.0.0">9.0.0</a></span>
</h3>
<ul style="list-style-type: none;">
<li>Binary-only (-bin) ZIP Checksum: <code>bin900</code></li>
<li>Complete (-all) ZIP Checksum: <code>all900</code></li>
<li>Wrapper JAR Checksum: <code>wrapper900</code></li>
</ul>
<h3 class="u-text-with-icon">
<span class="text-with-icon__icon" aria-hidden="true"></span>
<span><a id="v8.14.3"></a>v8.14.3</span>
</h3>
<ul style="list-style-type: none;">
<li>Binary-only (-bin) ZIP Checksum: <code>bin814</code></li>
<li>Complete (-all) ZIP Checksum: <code>all814</code></li>
<li>Wrapper JAR Checksum: <code>wrapper814</code></li>
</ul>
<h3 class="u-text-with-icon">
<span class="text-with-icon__icon" aria-hidden="true"></span>
<span><a id="0.7" href="#0.7">0.7</a></span>
</h3>
<ul style="list-style-type: none;">
<li>Binary-only (-bin) ZIP Checksum: <code>bin07</code></li>
<li>Complete (-all) ZIP Checksum: <code>all07</code></li>
<li>Wrapper JAR Checksum: <code>wrapper07</code></li>
</ul>
</body></html>
`

func crawlGradleFixture(t *testing.T) *Gradle {
	t.Helper()
	// Keep the local test request away from a developer's configured proxy.
	t.Setenv(req.ProxyEnvName, "")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, gradlePageFixture)
	}))
	t.Cleanup(srv.Close)

	g := NewGradle()
	g.DownloadUrl = srv.URL
	g.Start()
	return g
}

func TestGradleParsesVersionWithoutVPrefix(t *testing.T) {
	g := crawlGradleFixture(t)

	v, ok := g.Version["9.0.0"]
	if !ok {
		t.Fatalf("version 9.0.0 was not parsed, got: %v", g.Version)
	}
	if len(v) != 1 {
		t.Fatalf("version 9.0.0 should have exactly 1 item, got %d", len(v))
	}
	wantUrl := "https://services.gradle.org/distributions/gradle-9.0.0-all.zip"
	if v[0].Url != wantUrl {
		t.Errorf("url = %q, want %q", v[0].Url, wantUrl)
	}
	if v[0].Sum != "all900" {
		t.Errorf("sum = %q, want %q (Complete (-all) checksum)", v[0].Sum, "all900")
	}
}

func TestGradleParsesVersionWithVPrefix(t *testing.T) {
	g := crawlGradleFixture(t)

	if _, ok := g.Version["8.14.3"]; !ok {
		t.Fatalf("version 8.14.3 was not parsed, got: %v", g.Version)
	}
}

func TestGradleExcludesVersion07(t *testing.T) {
	g := crawlGradleFixture(t)

	if _, ok := g.Version["0.7"]; ok {
		t.Fatalf("version 0.7 should be excluded, got: %v", g.Version)
	}
}

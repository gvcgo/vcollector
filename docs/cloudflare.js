/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run "npm run dev" in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run "npm run deploy" to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

async function getAllowedSites() {
  let durl = new URL("https://raw.githubusercontent.com/gvcgo/vcollector/main/docs/allowed_sites.conf")
  try {
      let response = await fetch(durl);
      return await response.json();
  } catch (error) {
      console.log('get allowed_sites failed', error);
      return [];
  }
}

function makeRes(body, status = 200, headers = {}) {
  headers['access-control-allow-origin'] = '*'
  return new Response(body, {status, headers})
}
  
addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})
  
async function handleRequest(request) {
  const url = new URL(request.url);
  const actualUrlStr = url.pathname.replace("/proxy/","") + url.search + url.hash

  let allowed_sites = await getAllowedSites()
  if (!("github.com/gvcgo" in allowed_sites)) {
    allowed_sites = [
      "github.com/gvcgo",
      "github.com/moqsien",
      "raw.githubusercontent.com",
      "github.com/coursier",
      "github.com/zigtools",
      "github.com/nvarner",
      "github.com/elixir-lang",
      "dlcdn.apache.org/maven",
      "github.com/msys2",
      "ziglang.org/builds",
      "archive.apache.org/dist",
      "github.com/BurntSushi",
      "repo.anaconda.com/miniconda",
      "ziglang.org/download",
      "github.com/junegunn",
      "dl.k8s.io/release",
      "github.com/Pure-D",
      "cygwin.com/setup-x86_64.exe",
      "github.com/erlang",
      "github.com/sharkdp",
      "github.com/gvcgo",
      "github.com/v-analyzer",
      "github.com/charmbracelet",
      "github.com/exaloop",
      "github.com/VirtusLab",
      "github.com/vlang",
      "github.com/gerardog",
      "storage.googleapis.com/flutter_infra_release",
      "julialang-s3.julialang.org/bin",
      "dl.google.com/android",
      "github.com/protocolbuffers",
      "github.com/odin-lang",
      "go.dev/dl",
      "nodejs.org/download",
      "github.com/gleam-lang",
      "github.com/neovim",
      "vscode.download.prss.microsoft.com/dbazure",
      "github.com/git-for-windows",
      "github.com/asciinema",
      "github.com/JetBrains",
      "download.visualstudio.microsoft.com/download",
      "downloads.dlang.org/releases",
      "github.com/Enter-tainer",
      "github.com/denoland",
      "github.com/jesseduffield",
      "github.com/Kitware",
      "gradle.org/releases",
      "static.rust-lang.org/rustup",
      "github.com/clojure",
      "github.com/upx",
      "github.com/tree-sitter",
      "github.com/oven-sh",
      "windows.php.net/downloads",
      "github.com/bell-sw"
    ]
  }
  
  for (var key in allowed_sites) {
    if (actualUrlStr.includes(allowed_sites[key])) {
      const actualUrl = new URL(actualUrlStr)

      const modifiedRequest = new Request(actualUrl, {
        headers: request.headers,
        method: request.method,
        body: request.body,
        redirect: 'follow'
      });
      const response = await fetch(modifiedRequest);
      const modifiedResponse = new Response(response.body, response);
      // 添加允许跨域访问的响应头
      modifiedResponse.headers.set('Access-Control-Allow-Origin', '*');
      return modifiedResponse;
    }
  }

  var infoStr = allowed_sites.join("\n")
  var resp = makeRes("unsupported url: "+ actualUrlStr + "\n \nallowed urls: \n\n" + infoStr, 502)
  return resp
}

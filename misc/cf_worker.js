/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run "npm run dev" in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run "npm run deploy" to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

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
  
    const allowed_sites = [
      "vo.msecnd.net",
      "www.cygwin.com",
      "github.com/microsoft/vcpkg",
      "github.com/git-for-windows/git",
      "go.dev",
      "gradle.org",
      "github.com/gerardog/gsudo",
      "github.com/moqsien",
      "github.com/JohyC/Hosts",
      "github.com/ineo6/hosts",
      "github.com/sengshinlee/hosts",
      "githubusercontent.com/JohyC/Hosts",
      "githubusercontent.com/ineo6/hosts",
      "githubusercontent.com/sengshinlee/hosts",
      "oracle.com/java",
      "julialang-s3.julialang.org/bin",
      "dlcdn.apache.org/maven",
      "github.com/lyc8503/sing-box-rules",
      "github.com/Loyalsoldier/v2ray-rules-dat",
      "github.com/neovim/neovim",
      "github.com/protocolbuffers/protobuf",
      "github.com/pyenv/pyenv",
      "github.com/pyenv-win/pyenv-win",
      "rust-lang.org",
      "github.com/typst/typst",
      "github.com/vlang/v",
      "github.com/v-analyzer/v-analyzer",
      "github.com/JetBrains/kotlin",
      "github.com/lampepfl/dotty",
      "github.com/msys2/msys2-installer",
      "github.com/zigtools/zls",
      "github.com/neovide/neovide",
      "github.com/AstroNvim/AstroNvim",
      "github.com/tree-sitter/tree-sitter",
      "rustup-init",
      "rustup.rs",
      "rust-lang.org",
      "github.com/tree-sitter/tree-sitter",
      "github.com/junegunn/fzf"
    ]
  
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
  
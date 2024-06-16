/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run "npm run dev" in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run "npm run deploy" to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

async function getInstallScript(sUrl) {
    let durl = new URL(sUrl)
    try {
        let response = await fetch(durl);
        return response;
    } catch (error) {
        console.log('get allowed_sites failed', error);
        let resp  = makeRes("get installation script failed!", 502)
        return ;
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
    let actualUrl = "https://raw.githubusercontent.com/gvcgo/version-manager/main/scripts/install.preview.sh"
    const url = new URL(request.url);
    if (url.pathname.match("/windows")) {
        actualUrl = "https://raw.githubusercontent.com/gvcgo/version-manager/main/scripts/install.preview.ps1"
    }

    let response = await getInstallScript(actualUrl);
    const modifiedResponse = new Response(response.body, response);
    modifiedResponse.headers.set('Access-Control-Allow-Origin', '*');
    return modifiedResponse;

    // var infoStr = allowed_sites.join("\n")
    // var resp = makeRes("unsupported url: "+ actualUrlStr + "\n \nallowed urls: \n\n" + infoStr, 502)
    // return resp
}

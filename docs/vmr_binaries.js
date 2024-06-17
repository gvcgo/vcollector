/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run "npm run dev" in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run "npm run deploy" to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

async function getLatestTag() {
    let durl = "https://github.com/gvcgo/version-manager/releases/latest"
    try {
        let response = await fetch(durl);
        return response;
    } catch (error) {
        console.log('get allowed_sites failed', error);
        let resp  = makeRes("get installation script failed!", 502)
        return ;
    }
}

async function getResponse(dUrl) {
    try {
        let response = await fetch(dUrl);
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

// https://github.com/gvcgo/version-manager/releases/download/v0.6.2/vmr_darwin-amd64.zip

addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})
    
async function handleRequest(request) {
    let rUrl = new URL(request.url)
    let resp = await getLatestTag();
    let tag = resp.url.split("/").pop()
 
    let response = makeRes("download failed", 502)
    const downloadUrl = "https://github.com/gvcgo/version-manager/releases/download/"
    if (tag) {
        if (rUrl.pathname.match("/darwin/amd64")) {
            response = await getResponse(downloadUrl + tag + "/vmr_darwin-amd64.zip")
        }else if (rUrl.pathname.match("/darwin/arm64")) {
            response = await getResponse(downloadUrl + tag + "/vmr_darwin-arm64.zip")
        } else if (rUrl.pathname.match("/linux/amd64"))  {
            response = await getResponse(downloadUrl + tag + "/vmr_linux-amd64.zip")
        } else if (rUrl.pathname.match("/linux/arm64"))  {
            response = await getResponse(downloadUrl + tag + "/vmr_linux-arm64.zip")
        } else if (rUrl.pathname.match("/windows/amd64"))  {
            response = await getResponse(downloadUrl + tag + "/vmr_windows-amd64.zip")
        } else if (rUrl.pathname.match("/windows/arm64"))  {
            response = await getResponse(downloadUrl + tag + "/vmr_windows-arm64.zip")
        }
    }
    
    const modifiedResponse = new Response(response.body, response);
    modifiedResponse.headers.set('Access-Control-Allow-Origin', '*');
    return modifiedResponse;
}

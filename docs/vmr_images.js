/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run "npm run dev" in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run "npm run deploy" to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */


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

// https://raw.githubusercontent.com/gvcgo/vmrdocs/main/src/assets/vmr.gif

addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
})
    
async function handleRequest(request) {
    let rUrl = new URL(request.url)

    let fileName = rUrl.pathname.split("/").pop()
 
    let response = makeRes("download failed", 502)
    const downloadUrl = "https://raw.githubusercontent.com/gvcgo/vmrdocs/main/src/assets/"
    if (fileName) {
        response = await getResponse(downloadUrl + fileName)
    }
    
    const modifiedResponse = new Response(response.body, response);
    modifiedResponse.headers.set('Access-Control-Allow-Origin', '*');
    return modifiedResponse;
}

any("#go").on('click', ev => openWs(ev));
let pwField = any("#password")

let params = new URLSearchParams(location.search);
let requestID = params.get("id")

function openWs(ev) {
    let password = pwField.value
    const socket = new WebSocket('ws://' + location.host + '/challenge');

    socket.onopen = function (event) {
        console.log('Connection is open');
        socket.send(requestID);
    };

    socket.onmessage = async function (event) {
        console.log(event);
        if (event.type == "message") {
            let challengeString = password + event.data
            let challenge = new TextEncoder().encode(challengeString)
            let exportedKeySha = await window.crypto.subtle.digest("SHA-512", challenge);

            socket.send(exportedKeySha);
        }
    };

}

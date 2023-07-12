any("#go").on('mouseup', ev => openWs(ev));
me("#password").on('any', ev => halt(ev))

let params = new URLSearchParams(location.search);
let requestID = params.get("id")

let password;
let iv;

function openWs(ev) {
    let encodedData = me("#password").value
    let passwordStruct = JSON.parse(atob(encodedData))
    console.log("Struct: ", passwordStruct);
    let password = passwordStruct.key
    let iv = passwordStruct.iv

    const socket = new WebSocket('ws://' + location.host + '/challenge');

    socket.onopen = function (event) {
        console.log('Connection is open');
        socket.send(requestID);
    };

    socket.onmessage = async function (event) {
        console.log("EV: ", event);
        console.log("PWS: ", password)
        console.log("PW: ", password.k)
        console.log("IV: ", iv)
        if (typeof event.data === "string") {
            let encodedKey = new TextEncoder().encode(password.k)
            let keySha = await window.crypto.subtle.digest("SHA-512", encodedKey);
            let keyShaDecoded = buf2hex(keySha)

            console.debug("SHA (decoded):", keyShaDecoded)

            let challengeString = new TextEncoder().encode(keyShaDecoded + event.data);
            let challengeSha = await window.crypto.subtle.digest("SHA-512", challengeString);

            console.debug("SHA (decoded):", buf2hex(challengeSha))

            socket.send(challengeSha);
        } else if (event.data instanceof Blob) {
            console.debug("BLOB: ", event.data)

            let fileReader = new FileReader();
            fileReader.onload = function () {
                let arrayBuffer = this.result;
                let dataView = new DataView(arrayBuffer);
                console.debug("BLOB: ", arrayBuffer)
                console.debug("BLOB: ", dataView)

                decryptFile(dataView, passwordStruct)
                    .then(decryptedBlob => {
                        const url = window.URL.createObjectURL(decryptedBlob);
                        const downloadLink = document.createElement('a');

                        downloadLink.href = url;
                        downloadLink.download = 'decrypted_file'; // replace this with your preferred file name
                        document.body.appendChild(downloadLink);

                        downloadLink.click();

                        // Remember to revoke the URL after usage to free memory
                        window.URL.revokeObjectURL(url);
                    })
                    .catch(error => console.error("Error in file decryption: ", error));
            };

            fileReader.readAsArrayBuffer(event.data);
        } else {
            console.log("Data not readable.")
        }
    };
}

async function decryptFile(encryptedFile, passwordStruct) {
    let importedKey = await importAesKey(passwordStruct.key)
    let importedIv = new Uint8Array(Object.values(passwordStruct.iv));

    let fileBLob = await window.crypto.subtle.decrypt(
        { name: "AES-GCM", iv: importedIv },
        importedKey,
        encryptedFile
    ).catch((err) => {
        console.log("shit", err);
    });

    return new Blob([fileBLob])
}

async function importAesKey(jwk) {
    return await window.crypto.subtle.importKey(
        'jwk',
        jwk,
        {
            name: 'AES-GCM'
        },
        true,
        ["encrypt", "decrypt"]
    );
}

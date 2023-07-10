any(".file-input").on('click', ev => me(any("#form_file")).click());

any("#form_file").on("change", ev => fileHandler(ev));

any(".file-input").on("dragover", ev => halt(ev));

any(".file-input").on("drop", ev => { halt(ev); dropHandler(ev); });

function dropHandler(ev) {
  if (ev.dataTransfer.items) {
    if (ev.dataTransfer.items[0].kind === 'file') {
      let dT = new DataTransfer();
      dT.items.add(ev.dataTransfer.items[0].getAsFile())
      me("#form_file").files = dT.files
      uploadFile(dT.files[0].name)
    }
  } else {
    fileInput.files = ev.dataTransfer.files;
    uploadFile(ev.dataTransfer.files[0].name)
  }
}

function fileHandler(ev) {
  let file;
  if (ev.target.files) {
    file = ev.target.files[0];
  } else {
    return;
  }

  if (file) {
    let fileName = file.name;
    if (file.size > 200000000) {
      alert("You can't upload Files larger than 200 MB");
      return;
    }
    if (fileName.length >= 12) {
      let splitName = fileName.split('.');
      fileName = splitName[0].substring(0, 16) + "... ." + splitName[1];
    }
    uploadFile(fileName, file);
  }
}

async function uploadFile(name, f) {
  let fileName = name;
  let encodedData;
  let file = f;
  let data = new FormData(me(any("form")));

  if (me(any("#form_encrypt")).checked) {
    let { encryptedFile, key, sha, iv, base64 } = await encryptFile(file);

    encodedData = base64;

    // debug info
    console.log('Key:', key);
    console.log('IV:', iv);
    console.log('SHA512 (DECODED):', buf2hex(sha));
    console.log('SHA512:', sha);
    console.log('b64:', base64);

    data.append('file', new Blob([encryptedFile]), fileName);
    data.append('sha', new Blob([sha]))
    data.append('isEnc', true)
  } else {
    data.append('file', file, fileName);
  }

  let xhr = new XMLHttpRequest(); //AJAX request
  xhr.open("POST", "/upload"); //sending post request to the specified URL
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      callback(xhr.responseText, fileName, (me(any("#form_encrypt")).checked), encodedData);
    }
  }

  let element = createElement("div");
  me(element).classAdd("uploading");
  me(any("#progress-area")).prepend(element);

  xhr.upload.addEventListener("progress", ({ loaded, total }) => { //file uploading progress event
    let fileLoaded = Math.floor((loaded / total) * 100);  //getting percentage of loaded file size

    element.innerHTML = `
                        <div>${name} / ${fileLoaded}%</li>
                        <div class="progress-bar" style="width: ${fileLoaded}%;"}></li>
                        `

    if (loaded == total) {
      element.remove();
    }
  });

  xhr.send(data);
}

function callback(string, fileName, encrypted, b64) {
  let jsonStr = JSON.parse(string)
  console.log(jsonStr)

  element = createElement("div");
  me(element).classAdd("upload");
  me(any("#progress-area")).prepend(element);

  if (jsonStr.Ok === true) {
    let link = ""
    if (encrypted) {
      link = "http://" + location.host + "/dec.html?id=" + jsonStr.fileID;
    } else {
      link = "http://" + location.host + "/pv/" + jsonStr.fileID;
    }
    var qrcode = new QRious({
      value: link
    })
    element.innerHTML = `
                          <div class="uploaded-file-info">
                            <span>${fileName}</span>
                            <span><a href="${link}">${link}</a></span>
                            <spane${b64}</span>
                          </div>
                          <div class="qrcode">
                            <img src="${qrcode.toDataURL()}"/>
                          </div>`;
  } else {
    element.innerHTML = `
                          <div>
                            <span class="error">Upload Failed: ${jsonStr.reason}</span>
                          </div>`;
  }
}

async function encryptFile(file) {
  let fileBuffer = await file.arrayBuffer();

  // key generation
  let encKey = await window.crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt']);
  let iv = window.crypto.getRandomValues(new Uint8Array(12));
  let exportedKey = await window.crypto.subtle.exportKey("jwk", encKey);

  let encryptedFile = await window.crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, encKey, fileBuffer);

  let keyBuffer = new TextEncoder("utf-8").encode(exportedKey.k)
  let keySha = await window.crypto.subtle.digest("SHA-512", keyBuffer);
  let decodedKeySha = buf2hex(keySha)

  console.debug("SHA (decoded):", decodedKeySha)
  console.debug("KEY:", exportedKey.k)

  //base64 encode "password"
  let encodedData = btoa(JSON.stringify({ key: exportedKey, iv: iv }))

  return {
    encryptedFile: encryptedFile,
    key: keyBuffer,
    sha: keySha,
    iv: Array.from(iv).join(','),
    base64: encodedData
  };
}

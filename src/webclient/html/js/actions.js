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
  let fileSize;
  let file = f;
  let data = new FormData(me(any("form")));

  if (me(any("#form_encrypt")).checked) {
    let { encryptedFile, sha, key, iv } = await encryptFile(file);
    console.log('Key:', key);
    console.log('IV:', iv);
    data.append('file', new Blob([encryptedFile]), fileName);
    data.append('sha', sha)
    data.append('isEnc', true)
  } else {
    data.append('file', file, fileName);
  }

  let xhr = new XMLHttpRequest(); //AJAX request
  xhr.open("POST", "/upload"); //sending post request to the specified URL
  xhr.onreadystatechange = function () {
    if (xhr.readyState === 4) {
      callback(xhr.responseText, fileName, (me(any("#form_encrypt")).checked));
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

function callback(string, fileName, encrypted) {
  let jsonStr = JSON.parse(string)
  console.log(jsonStr)

  element = createElement("div");
  me(element).classAdd("upload");
  me(any("#progress-area")).prepend(element);

  if (jsonStr.Ok === true) {
    let link = ""
    if (encrypted) {
      link = "http://" + location.host + "/enc/" + jsonStr.fileID;
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
  let arrayBuffer = await file.arrayBuffer();
  let key = await window.crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, ['encrypt', 'decrypt']);
  let iv = window.crypto.getRandomValues(new Uint8Array(12));
  let encryptedFile = await window.crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, key, arrayBuffer);

  let exportedKey = await window.crypto.subtle.exportKey("jwk", key);
  let exportedKeySha = await window.crypto.subtle.digest("SHA-512", new TextEncoder().encode(exportedKey.k));
  console.log(exportedKeySha)
  console.log(buf2hex(exportedKeySha))

  return {
    encryptedFile,
    key: exportedKey.k,
    sha: exportedKeySha,
    iv: Array.from(iv).join(',')
  };
}

function buf2hex(buffer) { // buffer is an ArrayBuffer
  return Array.prototype.map.call(new Uint8Array(buffer), x => ('00' + x.toString(16)).slice(-2)).join('');
}
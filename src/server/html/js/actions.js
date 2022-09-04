any(".file-input").on('click', ev => me(any("#file")).click());

any("#file").on("change", ev => fileHandler(ev));

any(".file-input").on("dragover", ev => halt(ev));

any(".file-input").on("drop", ev => { halt(ev); dropHandler(ev); });

function dropHandler(ev) {
  if (ev.dataTransfer.items) {
    if (ev.dataTransfer.items[0].kind === 'file') {
      let dT = new DataTransfer();
      dT.items.add(ev.dataTransfer.items[0].getAsFile())
      fileInput.files = dT.files
      uploadFile(dT.files[0].name)
    }
  } else {
    fileInput.files = ev.dataTransfer.files;
    uploadFile(ev.dataTransfer.files[0].name)
  }
}

function fileHandler(ev) {
  let file;
  if (ev.target.files){
    file = ev.target.files[0];
  } else {
    return;
  }

  if(file){
    let fileName = file.name;
    if(file.size > 200000000){
      alert("You can't upload Files larger than 200 MB");
      return;
    }
    if(fileName.length >= 12){
      let splitName = fileName.split('.');
      fileName = splitName[0].substring(0, 16) + "... ." + splitName[1];
    }
    uploadFile(fileName);
  }
}

function uploadFile(name){
  let fileName = name;
  let fileSize;
  let xhr = new XMLHttpRequest(); //AJAX request
  xhr.open("POST", "/upload"); //sending post request to the specified URL
  xhr.onreadystatechange = function() {
    if (xhr.readyState === 4) {
      callback(xhr.responseText, fileName, fileSize);
    }
  }

  let element = createElement("div");
  me(element).classAdd("uploading");
  me(any("#progress-area")).prepend(element);

  xhr.upload.addEventListener("progress", ({loaded, total}) =>{ //file uploading progress event
    let fileLoaded = Math.floor((loaded / total) * 100);  //getting percentage of loaded file size
    
    element.innerHTML = `
                        <div>${name} / ${fileLoaded}%</li>
                        <div class="progress-bar" style="width: ${fileLoaded}%;"}></li>
                        `

    if(loaded == total){
        element.remove();
    }
  });
  
  let data = new FormData(me(any("form")));
  xhr.send(data);
}

function callback(string, fileName, fileSize) {
    let jsonStr = JSON.parse(string)
    let uploadedHTML;
    console.log(jsonStr)

    element = createElement("div");
    me(element).classAdd("upload");
    me(any("#progress-area")).prepend(element);

    if (jsonStr.Ok === true){
      let link = "http://" + location.host + "/pv/" + jsonStr.fileID;
      var qrcode = new QRious({
        value: link
      })
      element.innerHTML = `
                          <div>
                            <span class="name">${fileName}</span>
                            <span class="name"><a href="${link}">${link}</a></span>
                          </div>
                          <div>
                            <img src="${qrcode.toDataURL()}"/>
                          </div>`;
    } else {
      element.innerHTML = `
                          <div>
                            <span class="error">Upload Failed: ${jsonStr.reason}</span>
                          </div>`;
    }
  }

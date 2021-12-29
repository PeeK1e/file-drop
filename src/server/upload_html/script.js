//Original File/Code By CodingNepal - youtube.com/codingnepal
const form = document.querySelector("form"),
fileInput = document.querySelector(".file-input"),
progressArea = document.querySelector(".progress-area"),
uploadedArea = document.querySelector(".uploaded-area");

// form click event
form.addEventListener("click", () =>{
  fileInput.click();
});

//prevent default drag n drop behavior
function dragOverHandler(ev) {
  ev.preventDefault();
}

function dropHandler(ev) {
  // Prevent default behavior (Prevent file from being opened)
  ev.preventDefault();

  if (ev.dataTransfer.items) {
    // Use DataTransferItemList interface to access the file(s)
    // If dropped items aren't files, reject them
    if (ev.dataTransfer.items[0].kind === 'file') {
      let dT = new DataTransfer();
      dT.items.add(ev.dataTransfer.items[0].getAsFile())
      fileInput.files = dT.files
      uploadFile(dT.files[0].name)
    }
  } else {
    // Use DataTransfer interface to access the file(s)
    fileInput.files = ev.dataTransfer.files;
    uploadFile(ev.dataTransfer.files[0].name)
  }
}


fileInput.onchange = ({target})=>{
  let file = target.files[0]; //getting file [0] this means if user has selected multiple files then get first one only
  if(file){
    let fileName = file.name; //getting file name
    if(file.size > 200000000){
      alert("You can't upload Files larger than 200 MB")
      return
    }
    if(fileName.length >= 12){ //if file name length is greater than 12 then split it and add ...
      let splitName = fileName.split('.');
      fileName = splitName[0].substring(0, 13) + "... ." + splitName[1];
    }
    uploadFile(fileName); //calling uploadFile with passing file name as an argument
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
  xhr.upload.addEventListener("progress", ({loaded, total}) =>{ //file uploading progress event
    let fileLoaded = Math.floor((loaded / total) * 100);  //getting percentage of loaded file size
    let fileTotal = Math.floor(total / 1000); //getting total file size in KB from bytes
    // if file size is less than 1024 then add only KB else convert this KB into MB
    (fileTotal < 1024) ? fileSize = fileTotal + " KB" : fileSize = (loaded / (1024*1024)).toFixed(2) + " MB";
    let progressHTML = `<li class="row">
                          <i class="fas fa-file-alt"></i>
                          <div class="content">
                            <div class="details">
                              <span class="name">${name} • Uploading</span>
                              <span class="percent">${fileLoaded}%</span>
                            </div>
                            <div class="progress-bar">
                              <div class="progress" style="width: ${fileLoaded}%"></div>
                            </div>
                          </div>
                        </li>`;
    uploadedArea.classList.add("onprogress");
    progressArea.innerHTML = progressHTML;
    if(loaded == total){
      progressArea.innerHTML = "";
      //toRemove += uploadedArea.getElementsByClassName("onprogress")[0];
      uploadedArea.classList.remove("onprogress");
    }
  });
  let data = new FormData(form);
  xhr.send(data);
}

function callback(string, fileName, fileSize) {
  let jsonStr = JSON.parse(string)
  let uploadedHTML;
  console.log(jsonStr)
  if (jsonStr.Ok === true){
    let link = "http://" + location.host + "/pv/" + jsonStr.fileID;
    uploadedHTML = `<li class="row">
                            <div class="content upload">
                              <i class="fas fa-file-alt"></i>
                              <div class="details">
                                <span class="name">${fileName} • Uploaded</span>
                                <span class="name">Link • <a href="${link}">${link}</a></span>
                                <span class="size">${fileSize}</span>
                              </div>
                            </div>
                            <i class="fas fa-check"></i>
                          </li>`;
  } else {
    uploadedHTML = `<li class="row">
                            <div class="content upload">
                              <i class="fas fa-file-alt"></i>
                              <div class="details">
                                <span class="error">Upload Failed: ${jsonStr.reason}</span>
                              </div>
                            </div>
                            <i class="fas fa-check"></i>
                          </li>`;
  }
  uploadedArea.insertAdjacentHTML("afterbegin", uploadedHTML);
}

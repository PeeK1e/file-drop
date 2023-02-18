# File Drop
This is a fairly stright forward FOSS file share service. The goal is simple drop a file, upload it and share the link to family, friends and so on.

## The Goal

This Project was created beacause there were a lot of projects that didn't satisfy my needs of
> "I need a file-host that i can deploy myself, that isn't a pain in the ass to deploy as Container and auto deletes files after a given period"

*Yes I am aware of [XKCD 927](https://xkcd.com/927/)*

Hence this project was born.

## Features

**A Web GUI to upload the files which includes**
* encryption
* selectable expiry date
* a QR code to the download link

**An API**
* uploads
* downloads
* burn files after N downloads

**The Cleaner**
* Checks for expired files and deletes them
* Checks for "burned" files


## The Architecture

<img src="./img/architecture.svg" />

The Project consists of

* A Web Client which is a Container that deploys the Frontend in the Browser
* The API which handles uploads and downloads of the files
* A cleaner which runs every few minutes to remove files which expired or are marked as burned

## Deployment
You can either deploy the service with docker-compose with the`docker-compose.yml` or run every service by hand by either compiling the go source or pulling and running the image(s) from dockerhub.

If you deploy by hand you will need to setup Postgresql and modify the `dbSettings.json` in the `server/template/` directory.
___
This Project is subject to the GNU/AGPL License

(c) Ryuko

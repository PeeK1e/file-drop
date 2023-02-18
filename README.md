# File Drop
This is a fairly stright forward FOSS file share service. The goal is simple drop a file, upload it and share the link to family, friends and so on.

## The Goal

This Project was created beacause there were a lot of projects that didn't satisfy my needs of
> "I need a file-host that i can deploy myself, that isn't a pain in the ass to deploy as Container and auto deletes files after a given period"

*Yes I am aware of [XKCD 927](https://xkcd.com/927/)*

Hence this project was born.

## Features

**A Web GUI to upload the files which includes**
* encryption (soon)
* selectable expiry date (soon)
* a QR code to the download link

**An API handling**
* uploads
* downloads
* burn files after N downloads (soon when the frontend is done)

**The Cleaner**
* Checks for expired files and deletes them
* Checks for "burned" files

## Configuration and Deployment
This project was made to run inside containers and as scalable microservice architecture.

**I would not reccommend running it as plain binaries**

If you want to use the `docker-compose.yml` you will need a reverse proxy routing the API endpoints to the respective services. You can use the [./src/nginx.conf](./src/nginx.conf) for an example routing.

Otherwise use the `docker-compose-nginx.yml` file for a simple setup.


Configuring the filedrop server and cleaner

```env
#########################################
#           Server Only Values          #
#########################################

HTTP_LISTEN_ADDRESS="0.0.0.0"   # The Address the API will listen on
HTTP_PORT="8080"                # The port the API will listen on

#########################################
#        Server And Cleaner Values      #
#########################################

DATABASE_HOSTNAME="db"          # Postgres Hostname
DATABASE_PORT="5432"            # Postgres Port
DATABASE_USERNAME="user"        # Postgres Username
DATABASE_PASSWORD="pass"        # Postgres Password
DATABASE_DATABASENAME="uploas"  # Postgres Database Name
DATABASE_SSL="disable"          # Postgres SSL Mode {disable, verify-ca, ...}
```

## The Architecture

<img src="./img/architecture.svg" />

The Project consists of

* A Web Client which is a Container that deploys the Frontend in the Browser
* The API which handles uploads and downloads of the files
* A cleaner which runs every few minutes to remove files which expired or are marked as burned

## Deployment

___
This Project is subject to the GNU/AGPL License

(c) Ryuko

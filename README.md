## File Drop
This is a fairly stright forward FOSS file share service. The goal is simple drop a file, upload it and share the link to family, friends and so on.
___
#### Deployment
You can either deploy the service with docker-compose with the`docker-compose.yml` or run every service by hand by either compiling the go source or pulling and running the image(s) from dockerhub.

If you deploy by hand you will need to setup Postgresql and modify the `dbSettings.json` in the `server/template/` directory.
___
This Project is subject to the GNU/GPLv3 License

(c) Daniel Lehmann

## Steps to setup the project
1. Download the code into directory named fampay.
2. Install Mysql with the configuraions as per the yaml file.
3. Run the file db.sql is for creating DB and tables required.

Use the following code to build the docker image.
`docker build -t fampaay .  `

Use the following code to run the docker image.
`docker run -p 8080:8080/tcp fampaay`
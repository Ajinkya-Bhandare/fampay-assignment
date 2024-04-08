## Steps to setup the project
1. Download the code into directory named fampay.
2. Install Mysql with the configuraions as per the yaml file.
3. Run the file db.sql is for creating DB and tables required.

Use the following code to build the docker image.
`docker build -t fampaay .  `

Use the following code to run the docker image.
`docker run -p 8080:8080/tcp fampaay`

For the getting videos api use the following link
`http://localhost:8080/getVideos?PageSize=8&PageNumber=1`
While retrieving videos, all records are listed in descending order of publishing time. we have to mention page size and page number to get the required set of videos.

For the Search api, use the following link
`http://localhost:8080/getSearch?Query=cat`
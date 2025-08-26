# go-api-integrated-sql
the most of the code is from a tutorial but i tried to integrate the code with mysql inside a docker 

this was the command for running an sql server 
`docker run --name booksdb -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=booksdb -p 3306:3306 -d mysql:8`
and the sql command in teh .sql file

# resty
RestApi framework GO Lang

<h1>Before start</h1>

Download docker environment
<pre>
https://github.com/VladimirYalumov/docker_pgdb_mongodb_redis
cd [path to docker-compose.yml]
docker-compose up -d
</pre>
if need do this? or delete mongodb service (it needs to logging)
<pre>
chmod -R 777 mongodb
</pre>

<h1>Start</h1>
1. If you use my docker add this to configs

config.yml
<pre>
db_host: "127.0.0.1"
db_name: "deewave"
db_user: "user"
db_password: "qwerty"
db_port: "5432"
smtp_server: "{yoursmtp server (make it empty? if you haven't smtp server)}"
email_user: "{your email on smtp server}"
email_password: "{your password for email on smtp server}"
redis_host: "127.0.0.1:6379"
</pre>
tests/functional_tests/config.yml
<pre>
db_host: "127.0.0.1"
db_name: "test"
db_user: "user"
db_password: "qwerty"
db_port: "5433"
smtp_server: ""
email_user: ""
email_password: ""
redis_host: "127.0.0.1:6378"
</pre>

start project
<pre>
go run .
</pre>

Use postman collection go rest api.postman_collection.json to test api requests

<h1>Tests</h1>
Run this command
<pre>
    cd tests/functional_tests
    go test -v
</pre>
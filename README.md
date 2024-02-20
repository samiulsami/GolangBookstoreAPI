## Description

A simple REST API written in Go

## Running the server

`git clone git@github.com:samiulsami/GolangBookstoreAPI.git`

`cd GolangBookstoreAPI`

`go install GoBookstoreAPI`

`GoBookstoreAPI serve --port=3000`

## API Endpoints

|method| url                                          | body                                        | actions                            |
|---|----------------------------------------------|---------------------------------------------|------------------------------------|
|GET| `http://localhost:3000/api/v1/get-token`     |  | Set cookie and receive a jwt token |
|GET| `http://localhost:3000/api/v1/books`         |                                             | returns all books in a JSON array  |
|GET| `http://localhost:3000/api/v1/books/{uuid}`  |                                             | returns book with given uuid       |
|POST| `http://localhost:3000/api/v1/books`         | a json object of appropriate format         | adds a book                        |
|PUT| `http://localhost:3000/api/v1/books/{uuid}`  | a json object of appropriate format                                 | updates book with given uuid       |
|DELETE| `http://localhost:3000/api/v1/books/{id}` |   | deletes book with given uuid       |

---

## Example Book json object
``````
{
    "uuid": null
    "name": "The Lord Of The Rings",
    "authorList": ["J.R.R. Tolkien"],
    "publishDate": "August 14, 2012",
    "ISBN": "978-0544003415"
}
``````

## Some cURL commands
#### Login using the credentials in the example .env file to receive a JWT $TOKEN and set cookie
```
curl --location 'http://localhost:3000/api/v1/get-token' \
--header 'Authorization: Basic QGRtaW5VU0VSbmFtZTpzdHJvbmdwYXNzd29yZA=='
```
#### Show all books
```
curl --location 'http://localhost:3000/api/v1/books' \
--header 'Authorization: Bearer $TOKEN'
```
#### Show book with given {id}
```
curl --location 'http://localhost:3000/api/v1/books/{uuid}' \
--header 'Authorization: Bearer $TOKEN'
```
#### Add book
```
curl --location 'http://localhost:3000/api/v1/books' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer $TOKEN' \
--data '{
    "name": "Lasdfdasford asdfadsasdfffasdfof the rings",
    "authorList": ["jrr tolfdaksien", "asldksalk", "asdfsdakjk"],
    "publishDate": "2023",
    "ISBN": "idk"
}'
```
#### Update book with given {id}
```
curl --location --request PUT 'http://localhost:3000/api/v1/books/{uuid}' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYxNjc1NTIsImlhdCI6MTcwNjE2NzI1Miwic3ViIjoiQGRtaW5VU0VSbmFtZSJ9.sw-ESZpt-Zhldo30xTMAhiKONmYy2W0CRecaCWyltD8' \
--data '{
    "name": "updated books",
    "authorList": ["jrr tolfdaksien", "asldksalk", "asdfsdakjk"],
    "publishDate": "2023",
    "ISBN": "idk"
}'
```
#### Delete book with given {id}
```    
curl --location --request DELETE 'http://localhost:3000/api/v1/books/{1}' \
--header 'Authorization: Bearer $TOKEN'
```
----

## References

- https://github.com/shn27/BookStoreApi-Go
- https://github.com/MobarakHsn/api-server
- https://github.com/golang-jwt/jwt
- https://pkg.go.dev/github.com/golang-jwt/jwt/v5#New

## Build and run Docker Image

Create a .env file with the appropriate username, password, and JWTSecretKey (see the .env.example file)

Run the following commands:

``docker build -t sami7786/gobookstoreapi:latest .``

``docker run -p 3000:3000 sami7786/gobookstoreapi``

## Run with Kind

``kind create cluster``

``kubectl apply -f deploy/raw/deployment.yaml``

~~kubectl get pods~~

~~kubectl logs -f $pod_name~~

~~kubectl port-forward $pod_name 3000:3000~~

``kubectl get nodes -o wide``

The API endpoints are exposed on the internal IP of the node on port 30000

``kubectl delete deployment --all``

## Run using Helm chart

``helm install bookstorehelm deploy/helm``

## Makefile

Compile and run the binary:

``make run``

Run with docker:

``make dockerRun``

Run with helm:

``make helmInstall``

Push docker image to dockerhub:

``make dockerPush DOCKER_USER='username' DOCKER_PASSWORD='password'``

Stop and delete all running binaries/containers/pods:

``make cleanAll``
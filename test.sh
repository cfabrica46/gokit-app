#!/bin/bash

#SignUp
# response=$(curl -X POST -k http://localhost:8080/signup -d '{"username":"cesar","password":"01234","email":"cfabrica46@gmail.com"}')

#Signin
response=$(curl -X POST -Lk http://localhost:8080/signin -d '{"username":"cesar","password":"01234"}')
# curl -X POST -Lk http://localhost:8080/signin -d '{"username":"cesar","password":"01234"}'

# echo "$response"

token=$(echo "$response" | jq -r '.token')

# echo "$token"

#ShowUsers
# curl -X GET -Lk http://localhost:8080/users

#Profile
curl -X POST -k http://localhost:8080/profile -H "Authorization: $token"

#Delete
# curl -X DELETE -Lk http://localhost:8080/profile -H "Authorization: $token"

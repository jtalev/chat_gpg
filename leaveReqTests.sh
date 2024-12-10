#!/bin/bash
echo "POST"
curl -X POST -d "name=John&type="annual"&from=2024-12-01&to=2024-12-05&note=Vacation" http://localhost/post-leave-request
echo "POST"
curl -X POST -d "name=John&type="annual"&from=2024-12-01&to=2024-12-05&note=Vacation" http://localhost/post-leave-request
echo "get all"
curl -X GET http://localhost/get-leave-requests
echo "get by id"
curl -X GET "http://localhost/get-leave-request-by-id?requestId=87654321"
echo "update"
curl -X PUT "http://localhost/put-leave-request?requestId=87654321&name=John&type=updated&from=2024-12-01&to=2024-12-05&note=updated" 
echo "delete"
curl -X DELETE "http://localhost/delete-leave-request?requestId=87654321" 
echo "get all"
curl -X GET http://localhost/get-leave-requests
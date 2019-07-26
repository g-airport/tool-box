## Upload File


```bash

curl -X POST \
  http://127.0.0.1:8080/upload \
  -H 'Accept: */*' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Host: 127.0.0.1:8080' \
  -H 'Postman-Token: eed16e0e-0db0-41e5-ba6a-f81b9cca6535,10bf4433-6d45-4ce3-9a1e-913b12153dc4' \
  -H 'User-Agent: PostmanRuntime/7.15.0' \
  -H 'accept-encoding: gzip, deflate' \
  -H 'cache-control: no-cache' \
  -H 'content-length: 9674' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F type=csv \
  -F upload_file=@/Users/tqll/Downloads/email_sample.csv
```
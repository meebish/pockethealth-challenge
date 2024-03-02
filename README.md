# PocketHealth Backend Programming Challenge

## Overview

This programming challenge aims to deal with DICOM files with the following 3 goals:

* Allowing a user to upload a DICOM file
* Allow users to retrieve any DICOM Header Attribute of a given DICOM file through DICOM tags
* Allow users to retrieve any DICOM file as a PNG to view in-browser

## Getting Started

You'll need to have `Go` installed. I'm running this on version 1.22.0

```sh
git clone git@github.com:meebish/pocket-health.git

cd pocket-health

go run cmd/main.go
```

# API Endpoints
By default, the API currently runs on `localhost:8080`

| HTTP Verbs | Endpoints | Action |
| --- | --- | --- |
| POST | /upload | Uploads a new DICOM File |
| GET | /dicomFile/:fileName? | Retrieves Header Attributes or a PNG of a given dicom file |


## Upload DICOM File
### Request
`POST /upload`

Example
```sh
curl -i -X POST \
-H 'content-type: multipart/form-data' \
-F file=@/Users/mibo/Downloads/SE000001/IM000001 \
http://localhost:8080/upload
```

### Sample Response
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Thu, 29 Feb 2024 22:54:42 GMT
Content-Length: 78

{"message":"File saved as: ec122fa4-9626-401c-a9b1-f48b8818b0fc-IM000001.dcm"}
```

## Get DICOM Header Attributes
### Request
`GET /dicomFile/:filename?tag=(xxxx,yyyy)`

Where xxxx = tag group and yyyy = tag element

Example
```sh
curl -i -X GET \
  'http://localhost:8080/dicomFile/89b3857f-f13e-48a6-92f7-852a27a33420-IM000001.dcm?tag=(0002%2C0002)'
```

### Sample Response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 29 Feb 2024 23:33:17 GMT
Content-Length: 146

{"data":{"header-attribute-name":"MediaStorageSOPClassUID","header-attribute-value":"[1.2.840.10008.5.1.4.1.1.1.1.1]","tag-values":"(0002,0002)"}}
```

## Get DICOM as a png
### Request
`GET /dicomFile/:filename?png`

Example
```sh
curl -v -X GET \
  -o image.png \
  'http://localhost:8080/dicomFile/89b3857f-f13e-48a6-92f7-852a27a33420-IM000001.dcm?png'
```

Or go to the same link in a browser

### Sample Response
```
HTTP/1.1 200 OK
Content-Type: image/png
Date: Fri, 01 Mar 2024 00:12:48 GMT
Transfer-Encoding: chunked

[36651 bytes data]
```

An image similar to 

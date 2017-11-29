## Backend API:

# /register POST

Request JSON body contains:
 - *login : string
 - *password : string
    
Response JSON body contains:
 - token : string
 - expires : datetime string in ISO8601 format

# /login POST

Request JSON body contains:
 - *login : string
 - *password : string
 
Response JSON body contains:
 - token : string
 - expires : datetime string in ISO8601 format


# /search POST
JSON body contains:
 - from : string
 - to : string
 - date : datetime string in ISO8601 format
    
Headers contain:
 - *token : string

Response JSON body contains:
 - tickets : array of objects, each of which contains the following fields
    - from : string
    - to : string
    - departure : datetime string in ISO8601 format
    - carriage : int
    - seat : int


# /directions GET
Headers contain:
 - *token : string
 
Response JSON body contains:
 - locations : array of strings

# /departures GET
Headers contain:
 - *token : string

Response JSON body contains:
 - locations : array of strings


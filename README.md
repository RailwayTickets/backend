##Backend API:

#/register
POST
Request JSON body contains:
    - *login (string)
    - *password (string)
Response JSON body contains:
    - token (string)
    - expires (datetime string in ISO8601 format)

#/login
POST
Request JSON body contains:
    - *login (string)
    - *password (string)
Response JSON body contains:
    - token (string)
    - expires (datetime string in ISO8601 format)


#/search
POST
JSON body contains:
    - from (string)
    - to (string)
    - date (datetime string in ISO8601 format)
Headers contain:
    - *token (string)


#/directions
GET
Headers contain:
    - *token (string)

#/departures
GET
Headers contain:
    - *token (string)

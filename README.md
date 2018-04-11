## Distributed Key-Value Store


The Proxy node balances load using the by created a shasum hash of the requested key modulo with the number of servers available.
To Run:

`./execute`

---

The system accepts simple get and set request for the key-value store with the following API:

`http://localhost:5555/set`

Request body:

`{ "key": "the_key", "value" : "the_value" }`

Request type: `POST/PUT/GET`

Response body:

{ "message": "Success" }

---

`http://localhost:5555/get`

Request type: `POST/PUT/GET`

with body as:

`{ "key": "the_key" }`

Response body:

`{ "value": "the_value" }`

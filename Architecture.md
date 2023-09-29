Architecture
=============

API -> business logic -> redis

API Design
-----------

- The user will POST a json that creates a Resource called "ShortURL".
- Response will contain Location header to point to the Resource created.
- If request url contains path param, then we'll redirect to the resolved url else 404

Code Design (Packages and their responsibilites)
-------------------------------------------------

- the "rest"" package will work as front controller
  - It will allow CRUD operation on the resource called ShortURL, in a RESTful way
  - It also facillates redirection of short url on GET request

- the store package abstracts underlying key-value store

- the svc package stores core business logic that shortens the url

- we are separating these packages so that future change can be easily accommodated
- we can easily test each layer(package) independently
- We can plug additional layers, for example authentication/authorization, with least disturbance to existing code



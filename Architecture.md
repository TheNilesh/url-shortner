Architecture
=============

API -> business logic -> redis

API Design
-----------

- The user will POST a json that creates a Resource called "ShortURL".
- Response will contain Location header to point to the Resource created.
- If request url contains path param, then we'll redirect to the resolved url else 404
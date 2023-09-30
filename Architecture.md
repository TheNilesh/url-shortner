Architecture
=============

API -> business logic -> redis

# API Design


- The user will POST a json that creates a Resource called "ShortURL".
- Response will contain Location header to point to the Resource created.
- If request url contains path param, then we'll redirect to the resolved url else 404

# Code Design (Packages and their responsibilites)

- the "rest"" package will work as front controller
  - It will allow CRUD operation on the resource called ShortURL, in a RESTful way
  - It also facillates redirection of short url on GET request

- the store package abstracts underlying key-value store

- the svc package stores core business logic that shortens the url

- we are separating these packages so that future change can be easily accommodated
- we can easily test each layer(package) independently
- We can plug additional layers, for example authentication/authorization, with least disturbance to existing code

# Considering Requirements

## Map known URL with known short path

If I again ask for the same URL, itshould give me the same URL as it gave before instead
of generating a new one

To implement this we need to allow reverse loopkup. To allow that, we will create an additional k-v store that maps, hash of the URL with
the short_path(id of short url).
To save the additional lookup, we can redesign existing implementation such taht code will calculate hash of the url everytime shortening is requested and returning hash as short url. This way URL will be linked to its short ID by a hash function.
This improvement is not worth because common hash functions typically return 8 bytes. 4 english letters are required represent a byte, hence 32(8*4) letters short url cant be called short.

What should be key in reverse lookup store?

1. Put string urls as key directly
2. Base64 encode the url and then use it as a key
3. Compute MD5 hash of the url and then use it as a key

The 2 and 3 requires additional processing but makes debugging easier. The CPUs are expensive than memory hence lets go with 1.

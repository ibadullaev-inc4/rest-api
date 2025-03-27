GET     /users        200, 404, 500
GET     /users/:id    200, 404, 500
POST    /users/:id    204, 404, 400
PUT     /users/:id    204/200, 404, 400, 500
PATCH   /users/:id    204/200, 404, 400, 500
DELERE  /users/:id    204, 404, 400
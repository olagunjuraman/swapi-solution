# Bush-swars

## How to run the application

A live version is currently running on 

Clone the application with `git clone https://github.com/bushaHQ/swapi-solution.git`

### Configuration
Copy the .env.example to .env file and configure your database, redis and host/port details.

Ensure Redis and Postgres is running

### Docker

```
cd bush-swapi
docker-compose up --build
```

### Locally

To run the application locally, please use the following command -


        go run ./
 
> Note: By default the port number its being run on is **8100**.

## Endpoints Description

### Get All Movies

This makes a request to swapi.dev to fetch all movies from the /films endpoint. 
However, responses are cached on redis for all subsequent requests. This enpoint returns the Total comments for each movie, Name, Opening Crawl and Release data.

```
    URL - /movies
    Method - GET
    Response - (content-type = application/json)
``` 
```JSON
    {
        "data": [
            {
                "name": "A New Hope",
                "opening_crawl": "It is a period of civil war....",
                "total_comment": 6,
                "release_date": "1977-05-25"
            }
        ],
        "message": "Movies",
        "status": 200
    }
```

### Add Comment to a Movie (POST Request)
This endpoint adds a comment to a movie. Parameters needed are Film_id (int) and comment (string).
However, if the movieID does not exists, it would return an error.

```
    URL - /comment
    Method - POST
    Request - (content-type = application/json)
```
```JSON
    {
        "film_id" :1,
        "comment": "Testing comment"
    }
```

### Get All comments for a movie
This endpoint returns all comments posted/stored for a movie. It returns the IP address of the person who made the comment, the comment alongside the date teh comment was created.
```
    URL - /comments/{movieID}
    Method - GET
    Response - (content-type = application/json)
``` 
```JSON
    {
      "data": [
        {
          "createdAt": "2021-06-02T11:48:58.90173+01:00",
          "comment": "It is a period of civil war.\r\nRebel spaceships, striking\r\nf.",
          "ipaddress": "::1"
        }
      ],
      "message": "Comments",
      "status": 200
    }
```

### Get All characters for a movie
This endpoint returns all characters in a movie. 
It has a metadata object that contains the total count of characters based on filers or sorts (if any), total height in both 
cm and ft/inches.

To filter by gender; You can pass a get parameter such as ?fliter_by=female || ?filter_by=male

To sort by name, gender or height; You can pass a get parameter such as ?sort_by=name || ?sort_by=height || ?sort_by=gender

By default, it responds in Ascending order, but to change that. You can pass ?order_by=asc for ascending || ?order_by=desc for descending order.

```
    URL - /characters/{movieID}
    Method - GET
    Response - (content-type = application/json)
``` 
```JSON
    {
  "data": {
    "metadata": {
      "count": 1,
      "total_height_meters": "170",
      "total_height_feet": "5ft and 6.93 inches"
    },
    "characters": [
      {
        "gender": "n/a",
        "name": "R5-D4",
        "height": "170"
      }
    ]
  },
  "message": "Characters in this movie",
  "status": 400
}
```


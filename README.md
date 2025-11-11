# how to run

go build .

go run .

server starts on :8080 port, accepts /execute requests

``POST http://localhost:8080/execute``

with following JSON body (example):
```
{
	"execution_timeout": 800,
	"url_requests": [
		{
			"url": "https://google.com",
			"timeout": 800,
			"headers": {
				"X-Device-IP": "10.20.3.15",
				"other-custom-header":"other-custom-value"
			}
		},
		{
			"url": "https://stackoverflow.com/",
			"timeout": 900,
			"headers": {
				"X-Device-IP": "10.20.3.15"
			}
		}
	]
}
```
execution_timeout - optional, defaults to 800ms
url_request.timeout - optional, defaults to execution_timeout. 
If timeout > execution_timeout, execution_timeout will be used instead.


Response
```
{
	"results": [
		{
			"url": "https://google.com",
			"code": 200,
			"payload": "<!doctype html><html>...</html>"
		},
	    {
			"url": "https://amazon.com",
			"code": 500,
			"error": "Get \"https://amazon.com\": context deadline exceeded"
		},
		{
			"url": "https://stackoverflow.com/",
			"code": 403
		}
	]
}
```

# Technical decisions rationale

Gin router - go-to package for routing, allows for easy context handling and response writing

go-playground/validator - convenient input request validation, with model tags

Project structure:
* controller - api read/write logic
* models - request and response structures
* service - main execution logic
* utils - helper functions

Utils functions are generic and in larger project would be shared between multiple services

Service logic is purely functional, in larger project definitely would want to use interface
and service struct, which would make code much more easily mockable and testable.
However, that would add some overhead to properly initialize service(s), which in project of this scale is unwarranted.

Implementation details.

Running GET requests for given list of URLs concurrently, with semaphore to limit simultaneous concurrent requests.

Using context with timeout to cancel requests based on global or url request timeout.

Added a simple output cleanup function to make requests aborted due to timeout more presentable.
# HTTP Server-Timing for Gin (a HTTP web framework written in Golang)
[![Godoc](https://godoc.org/github.com/p768lwy3/gin-server-timing?status.svg)](https://godoc.org/github.com/p768lwy3/gin-server-timing)

This is a project for myself to know more about API server in Golang and Thank you for 
repository [go-server-timing](https://github.com/mitchellh/go-server-timing/) by 
[mitchellh](https://github.com/mitchellh). It is my first time to play with Golang and 
this is quite cool.

## Description

This is a library to build a middleware of HTTP Server-Timing for Gin, which this header 
allow a server to send timing metrics from the backend to show processing times on the 
browser as below:

![Server Timing Example](https://raw.githubusercontent.com/p768lwy3/gin-server-timing/master/examples/screenshot.png)

## Usage

Example usage is shown below.

```go
func main() {
    // Build router from gin-gonic/gin
	router := gin.Default()

	// Wrap handler with timing middleware
	router.Use(servertiming.Middleware())

	// Build a testing example of routing
	router.GET("/", Handler)

	// Run gin application
	router.Run()
}


func Handler(c *gin.Context) {
	// Create a wait group to wait for the response of HTTP fetches
	var wg sync.WaitGroup

	// Get timing header builder from the context
	timing := servertiming.FromContext(c)

	// Samples for testing
	for i := 0; i < 5; i++ {
		// Increment the WaitGroup counter.
		wg.Add(1)

		// Launch a goroutine to fetch the URL.
		name := fmt.Sprintf("service-%d", i)
		go func(name string) {
			// Imagine handler performing tasks in a goroutine
			defer timing.NewMetric(name).Start().Stop()
			time.Sleep(random(25, 75))

			// Decrement the counter when the goroutine completes.
			wg.Done()
		}(name)
	}

	// Imagine blocking code in handler
	m := timing.NewMetric("sql").WithDesc("SQL query").Start()
	time.Sleep(random(20, 50))
	m.Stop()

	// Wait for all HTTP fetches to complete.
	wg.Wait()

	// Write header to the response adter all fetches done
	servertiming.WriteHeader(c)

	// Create response of gin
	c.String(http.StatusOK, "Done. Check your browser inspector timing details.")
}

func random(min, max int) time.Duration {
	return (time.Duration(rand.Intn(max-min) + min)) * time.Millisecond
}
```

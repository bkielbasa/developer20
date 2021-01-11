---
title: "Refactoring for better testability"
publishdate: 2020-11-17
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
    - system design
    - tests
---

When we talk about software design we very often use very generic and abstract words. But, how about the practice? How does it look in a real-world project? Today, you and I will start refactoring a small to-do app for better testability and maintainability. In this article, we will make the application testable. We'll write black-box tests that will prevent from some bugs and make future refactoring easier and safer.

This article is the first part of mini-series where we do a code review of existing code and try to improve it. In future articles, we'll extract the domain from the project as well as fix some bugs. So, stay tuned!

The full source code is available on [Github](https://github.com/bkielbasa/gotodo). When you take a look at the project you can learn about the project structure. It isn't that bad but there's one big issue that won't let us refactor the logic without any worries - it has no tests. This is our goal for today - make the project testable. To be able to refactor the application we have to have some guardian that will keep an eye on the changes you're making. How to do it?

## Refactoring the main() function

In the main.go file we have two functions. The `main` function holds the whole application initialization as well as starts the HTTP server. The `getDB()` just returns a connection to the database. We can see some handler's definitions. I read the code for you and prepared a few sample requests that we can use to test the application manually.

Create a new project:

{{< highlight bash >}}
curl --request POST \
  --url http://localhost:8090/project/create \
  --header 'content-type: application/json' \
  --data '{
	"name": "Home"
}
'
{{< / highlight >}}

List all available projects:

{{< highlight bash >}}
curl --request GET \
  --url http://localhost:8090/projects
{{< / highlight >}}

Add the todo:

{{< highlight bash >}}
curl --request POST \
  --url http://localhost:8090/todo/create \
  --header 'content-type: application/json' \
  --data '{
	"name": "lala"
}
'
{{< / highlight >}}

Mark the todo as done:

{{< highlight bash >}}
curl --request POST \
  --url http://localhost:8090/todo/TASK_ID/done \
  --header 'content-type: application/json'
{{< / highlight >}}

Mark as undone

{{< highlight bash >}}
curl --request POST \
  --url http://localhost:8090/todo/TASK_ID/undone \
  --header 'content-type: application/json'
{{< / highlight >}}

and so on. Manual testing works but it's not scalable. You want to replace them with automatic tests, right? Let's do it the simple way. Create a new file main_test.go where you'll keep all of the tests of the application. Think about how you want to interact with the program. You don't want to change a lot in the project so you'll ignore the fact that there's a real database inside. You just want to run the application. Here's how it can look like.

{{< highlight golang >}}
func TestRunServer(t *testing.T) {
  ctx := context.Background()
  run, shutdown := todo.App(ctx, port)
  defer shutdown()
  go run()

  // run your tests here
}
{{< / highlight >}}

What you do is creating a new application that uses the specified context and running it on the port. Then, you start the application, and when the tests end - shuts it down. Sounds simple, isn't it?

To achive that, we have to rename (for now) the `main()` function to `App.`

{{< highlight golang >}}
func App(ctx context.Context, port int) (func() error, func() error)
{{< / highlight >}}

The function accepts the context, a port the HTTP server will be running on, and returns two functions: for starting the server and for closing it. This will give us full control over the server that will be useful in a moment. The noticeable change made in the function is creating an HTTP server directly so we can control it. Here's the full function after changes.

{{< highlight golang >}}
func App(ctx context.Context, port int) (func() error, func() error) {
	m := http.NewServeMux()
	s := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: m}
	log.Printf("starting on port %d", port)

	db := getDB()
	repo := repositories.NewPostgres(db)
	todoHandler := handlers.ToDo{Repo: repo}
	projectHandler := handlers.Project{Repo: repo}
	m.HandleFunc("/projects", projectHandler.List)
	m.HandleFunc("/project/create", projectHandler.Create)
	m.HandleFunc("/project/{id:[0-9a-z\\-]+}/archive", projectHandler.Archive)
	m.HandleFunc("/todos", todoHandler.List)
	m.HandleFunc("/todo/create", todoHandler.Create)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}", todoHandler.Get)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}/done", todoHandler.MarkAsDone)
	m.HandleFunc("/todo/{id:[0-9a-z\\-]+}/undone", todoHandler.MarkAsUndone)
	return s.ListenAndServe, func() error {
		return s.Shutdown(ctx)
	}
}
{{< / highlight >}}

At this point, the test almost run. What we have to do is to add the missing main() function we removed. we can find the current code below. The deferred shutdown() function doesn't make sense right now because when it's executed the server is already gone but we'll take care of it next.

{{< highlight golang >}}
func main() {
	ctx := context.Background()
	run, shutdown := App(ctx, 8090)
	defer shutdown()
	err := run()
	if !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(err)
		os.Exit(1)
	}
}
{{< / highlight >}}

At this point, the test passes and the application can still run! You test almost nothing but the test will be improved later. Before that, there's one thing missing in the main() function - there's no other way of closing the application than just killing it. we need a [graceful shutdown](https://developer20.com/golang-tips-and-trics-iii/).

{{< highlight golang >}}
func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx := context.Background()
	run, shutdown := App(ctx, 8090)

	go func() {
		_ = <-gracefulStop
		fmt.Println("shutting down...")
		err := shutdown()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	err := run()
	if !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(err)
		os.Exit(1)
	}
}
{{< / highlight >}}

Before starting the server a new `os.Signal` channel is created where we'll receive a signal that it's the time to stop the process. We use our `shutdown()` function to stop the HTTP server and exit.

## Writing the first test

Our test will contain three steps: creating a new project, getting a list of available projects, and checking if our new brand project is visible on the list. Let's rename the test name to `TestAddingNewProject` and update its code to fit the requirements. Every time a new (unique) project name will be created to make sure if the project was created with the correct name.

{{< highlight golang >}}
  name := uuid.New().String()
	reqBody := fmt.Sprintf(`{"name": "%s"}`, name)
	url := "http://localhost:8090/project/create"
	client := http.Client{
		Timeout: time.Second,
	}

	// create a new project
	resp, err := client.Post(url, "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
    r := struct{ ID string }{}
	json.Unmarshal(respBody, &r)
	require.NotEmpty(t, r.ID)
{{< / highlight >}}

One noticeable thing is that we create a new instance of the standard HTTP client. The `http.DefaultClient` doesn't have any timeout set what, in some cases, may slow the test down. Waiting for timeouts can take some time :) 

We received the feedback that we created the project correctly and got the ID of it into an anonymous struct. Now, it's the time for checking if the new project was persisted in the database and can be read from it.

{{< highlight golang >}}
  // list existing projects
	url = "http://localhost:8090/projects"
	resp, err = client.Get(url)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listResp := projectsListResponse{}
	respBody, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	json.Unmarshal(respBody, &listResp)

	// check if the projects are on the list
	found := false

	for _, proj := range listResp.Projects {
		if proj.ID == r.ID {
			require.Equal(t, name, proj.Name)
			require.False(t, proj.Archived)
			found = true
		}
	}

	require.True(t, found, "cannot found the project on the list")
}
{{< / highlight >}}

It was done by checking the `/projects` endpoint. Firstly, the HTTP status is checked, and when it succeeds we check every project, one by one, and look for the brand new one. When you find it, you make sure that the name is OK and the project isn't archived from the very beginning. If everything's OK, the test passes!

## Summary

Today, we refactored the project a bit what helped us writing very first tests for it. We can find the diff with changes we made today in this PR [https://github.com/bkielbasa/gotodo/pull/1](https://github.com/bkielbasa/gotodo/pull/1).

Your homework is to write tests for other endpoints. One of scenarios we can write is creating a new project, archiving it and checking if it's status is changed to archived. As I said, this is the first part of refactoring mini-series. The project has more issues in both design and Go good practices. Before fixing that we have to have some tests, right? :)

I hope you liked the post and if you have any questions, leave them in the comments section below.
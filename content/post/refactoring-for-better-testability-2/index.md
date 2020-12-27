---
title: "Refactoring for better testability #2"
publishdate: 2020-12-23
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
    - system design
    - tests
    - ddd
---

In the [previous article](https://developer20.com/refactoring-for-better-testability/), we wrote a few tests for a project to make sure that our refactoring won’t break anything. This time, we’ll extract part of the domain and add a test to it. It will help us better understand what the project does and make those tests more reliable.

We have an issue with the end-to-end (e2e) tests: database under the hood. The approach has some problems. Firstly, those tests are slower. We use a real database connection that has an overhead. The database runs on the same machine, so the latency isn’t significant right now, but it can be when the number of tests increases.

The second consequence is that those tests aren't as stable as isolated once. We have to remember about launching the database before running tests, running all migrations, and (sometimes) purging tables. If something can crash - it will eventually happen. If we want to have the CI useful, we have to run those tests there as well. We have to configure the CI the same way we did it on our local machines. The setup is much more complicated than just running `go test ./...`.

From my experience, integration tests **are** helpful. Unit tests should be the core of our tests' sets. This is our motivation to write them.

We have to understand the core domain first. From the order of requests we send, we can assume that creating the project is the starting point. Let's see what the handler contains.

```go
func (p Project) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("cannot read the body: %s", err)
		http.Error(w, "cannot read the body", http.StatusBadRequest)
		return
	}

	req := httpmodels.CreateProjectRequest{}
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("cannot read the body: %s", err)
		http.Error(w, "invalid JSON provided", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Printf("the name is required")
		http.Error(w, "the name is required", http.StatusBadRequest)
		return
	}

	id, err := p.Repo.CreateProject(req.Name)
	if err != nil {
		log.Printf("internal server error: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := httpmodels.CreateProjectResponse{id}
	b, _ = json.Marshal(resp)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
```

First lines aren't very useful for now. It's just standard reading the request and unmarshalling to a struct. There's an `if` statement that looks suspiciously.

```go
	if req.Name == "" {
		log.Printf("the name is required")
		http.Error(w, "the name is required", http.StatusBadRequest)
		return
	}
```
It says we require providing the name. It’s obligatory. When we continue reading, notice that we return the `id` of the project to the API caller. The id is retrieved from the function that creates the project. To create a new project we have to provide its name. Every project has an `ID`. Let’s model this in the code.

```go
type Project struct {
    id string
    name string
}

func (p Project) Name() string {
	return p.name
}

func (p Project) ID() string {
	return p.id
}
```

We created a new struct with private methods. Why? We want to make sure the `Project` is always in a correct state. Private fields help us with this. To let to get those fields we wrote two getter methods. How about creating? Go doesn't have constructors. It doesn't mean we cannot create a custom constructor function.

In our domain, the project has to have a valid (not empty) name and an ID. We can achieve this in at least two ways. The first one is creating the constructor method. The constructor method will do all the checks. The second approach is to create a function `func (p Project) IsValid() bool`. We'll call it everytime we want to check if the project is a valid object.

Personally, I prefer the first option, but the second one is valid as well. It's all about preferences and the specific case. It's the time for the test. Create a new file called `domain/project_test.go` and put the test as shown below. Please notice that we created a new `domain` package.

I prefer the first option, but the second one is valid as well. It’s all about preferences and the specific case. It’s the time for the test. Create a new file called `domain/project_test.go` and put the test as shown below. Please notice that we created a new domain package.

{{< info title="What's in the domain package?" msg="In Domain-Driven Design (DDD), the Domain is the core of our application. It holds all the business logic of the application. It cannot contain any code that interacts with the infrastructure. The Domain should be both platform and framework agnostic." >}}

```go
package domain

import "testing"

func TestProject_Test_Validation(t *testing.T) {
	testCases := map[string]struct {
		id string
		name string
	}{
		"empty ID": {
			name: "jfslfjal",
		},

		"empty name": {
			id: "jfslfjal",
		},
	}

	for _, tc := range testCases {
		_, err := NewProject(tc.id, tc.name)
		if err == nil {
			t.Error("expected that the validation fails but got no error")
		}
	}
}

```

We make sure we check all the requirements. The test is red (doesn't compile). There's no such a `NewProject` yet. It's the time to add it in `domain/project.go` file.

```go
func NewProject(id, name string) (Project, error)  {
	if id == "" {
		return Project{}, errors.New("the ID cannot be empty")
	}

	if name == "" {
		return Project{}, errors.New("the name cannot be empty")
	}

	return Project{id: id, name: name}, nil
}
```

Tests should be green now. We extracted the first part of the domain! The domain cannot talk with other parts of the code directly. We need an additional layer. Let’s create a new package and call it the `app` (for an application layer).
{{< info title="What's in the app package?" msg="The application layer is responsible for orchestrating the communication between the external world (DB, HTTP, etc) and your application. The flow is generally like this: get a domain object from a repository, execute an action, and put it back there." >}}

When we take a look at the HTTP handler for creating a project we’ll notice a simple flow: the user provides the name, we create a new project, and return its ID. Let’s write a test that will model it.

```go
package app

import (
	"context"
	"github.com/bkielbasa/gotodo/domain"
	"github.com/google/uuid"
	"testing"
)

func TestAddNewProject(t *testing.T) {
	name := "my name:" + uuid.New().String()
	ctx := context.Background()

	projectServ := NewProjectService()
	p, err := projectServ.Add(ctx, name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if p.ID() == "" {
		t.Errorf("ID is empty")
	}

	if name != p.Name() {
		t.Errorf("name don't match, expected (%s) but got (%s)", name, p.Name())
	}
}
```

What we do is creating a new application service. The service accepts the name of the project and returns a freshly created project followed by an error (if occurs). Just after that, we make sure the name is as we provided, and the ID isn’t an empty string (this is what we know about the ID right now). The test doesn’t compile. Let’s fix it.

```go
type ProjectService struct {}

func NewProjectService() ProjectService {
	return ProjectService{}
}

func (serv ProjectService) Add(ctx context.Context, name string) (domain.Project, error) {
	return domain.Project{}, nil
}
```

We create the missing constructor function for our new type - the application service. The service has a simple method with initial code - to make the code compile. When we run the test, we’ll notice that it fails. Nothing surprising because we do nothing in the Add function.

```go
func (serv ProjectService) Add(ctx context.Context, name string) (domain.Project, error) {
	id := "gopher"
	return domain.NewProject(id, name)
}
```

From now, the test is green. We can add one more test that will check if we validate the name correctly.

```go
func TestAddNewProjectWithEmptyName(t *testing.T) {
	name := ""

	projectServ := NewProjectService()
	_, err := projectServ.Add(context.Background(), name)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}
```


The test should are green. We do not check too much, so it's time to change it. We'll update the first test with getting a project for the particular ID and check if the `Get` method still returns the same project.

```go
func TestAddNewProject(t *testing.T) {
	name := "my name:" + uuid.New().String()
	ctx := context.Background()

	projectServ := NewProjectService()
	p, err := projectServ.Add(ctx, name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if p.ID() == "" {
		t.Errorf("ID is empty")
	}

	if name != p.Name() {
		t.Errorf("name don't match, expected (%s) but got (%s)", name, p.Name())
	}

	p2, err := projectServ.Get(ctx, p.ID())
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if p.ID() !=  p2.ID() {
		t.Errorf("expected ID %s but %s given", p.ID(),  p2.ID())
	}

	if p.Name() !=  p2.Name() {
		t.Errorf("expected name %s but %s given", p.Name(),  p2.Name())
	}
}
```

Hmm, the code looks a bit unreadable... We can refactor the code by providing a helper function `requireProject`.

```go
func TestAddNewProject(t *testing.T) {
	name := "my name:" + uuid.New().String()
	ctx := context.Background()

	projectServ := NewProjectService(newStoreMock())
	p, err := projectServ.Add(ctx, name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	checkProjectName(t, p, name)

	p2, err := projectServ.Get(ctx, p.ID())
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	checkProjectName(t, p2, p.Name())
	checkProjectID(t, p2, p.ID())
}

func checkProjectID(t *testing.T, p domain.Project, expectedID string) {
	if p.ID() !=  expectedID {
		t.Errorf("expected ID %s but %s given", expectedID,  p.ID())
	}
}

func checkProjectName(t *testing.T, p domain.Project, expectedName string) {
	if p.Name() !=  expectedName {
		t.Errorf("expected name %s but %s given", expectedName,  p.Name())
	}
}
```

Much better, isn't it? :) The code doesn't compile. To fix it we have to add the missing `Get` function.

```go
func (serv ProjectService) Get(ctx context.Context, id string) (domain.Project, error) {
	return domain.NewProject(id, "fjsfsl")
}
```

The test is still red. To make it work, we have to add storage that will keep the list of projects we created with the ability to fetch it back. This is how I designed its interface and update `Add()` and `Get` functions to use.

```go
type Storage interface {
	Store(ctx context.Context, p domain.Project) error
	Get(ctx context.Context, id string) (domain.Project, error)
}

func (serv ProjectService) Add(ctx context.Context, name string) (domain.Project, error) {
	id := "gopher"
	p, err := domain.NewProject(id, name)
	if err != nil {
		return domain.Project{}, err
	}

	err = serv.storage.Store(ctx, p)
	if err != nil {
		return domain.Project{}, err
	}

	return p, err
}

func (serv ProjectService) Get(ctx context.Context, id string) (domain.Project, error) {
	return serv.storage.Get(ctx, id)
}
```

The `ProjectService` doesn't contain the new functionality so let's add it now.

```go
type ProjectService struct {
	storage Storage
}

func NewProjectService(storage Storage) ProjectService {
	return ProjectService{storage: storage}
}
```

Almost there. We have to put the new dependency everywhere we create a new `ProjectService` struct.
We need a new struct that will implement the interface. Let's create a new one with a map that will hold the instances of `domain.Project`.

```go
type storeMock struct {
	data map[string]domain.Project
}

func newStoreMock() *storeMock {
	return &storeMock{
		data: make(map[string]domain.Project),
	}
}
func (s *storeMock) Store(ctx context.Context, p domain.Project) error {
	s.data[p.ID()] = p
	return nil
}

func (s *storeMock) Get(ctx context.Context, id string) (domain.Project, error) {
	return s.data[id], nil
}
```

It's green again! I'd add one more test because we did not cover one important case. What if the project doesn't exist? Shouldn't `Get` function return an error?
The storage knows if the project exists or not so the error should come from it. Let's create a separate error just for this case.

```go
// in app/project.go
var ErrProjectNotFound = errors.New("the project is not found")
```

To make our testing easier, we need to add a new error to the mock `storeMock` and create a new method to set the given error.

```go
type storeMock struct {
	data map[string]domain.Project
	err error // new field
}

func (s *storeMock) Get(ctx context.Context, id string) (domain.Project, error) {
	return s.data[id], s.err // added the error here
}

func (s *storeMock) withError(err error) *storeMock {
	s.err = err
	return s
}
```

When we are guarded with new helper methods, it's time to write the test.

```go
func TestAGetNotExistingProject(t *testing.T) {
	id := "not exists"
	ctx := context.Background()
	storage := newStoreMock().withError(ErrProjectNotFound)

	projectServ := NewProjectService(storage)

	_, err := projectServ.Get(ctx, id)
	if !errors.Is(err, ErrProjectNotFound) {
		t.Errorf("expected error ErrProjectNotFound but got %v", err)
	}
}
```

Almost done! If you're perceptive you noticed that we have a hardcoded ID for every project ID: `gopher`. Let's prepare a test that will force us to fix it.

```go
func TestEveryProjectShouldHaveUniqueID(t *testing.T) {
	name := "a name"

	projectServ := NewProjectService(newStoreMock())
	p1, err := projectServ.Add(context.Background(), name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	p2, err := projectServ.Add(context.Background(), name)
	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}

	if p1.ID() == p2.ID() {
		t.Error("every project should have a unique ID")
	}
}
```

It's red now. There are many ways to generate a unique ID. We'll use one of the simplest - [uuid](https://github.com/google/uuid).

```go
id := uuid.New().String()
```

That's all! Tests pass. We extracted the domain from the current code. Of course, it's not the whole business logic we have to refactor, but it's a good starting point.
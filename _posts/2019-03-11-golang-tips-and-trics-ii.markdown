---
layout: post
title: "Golang Tips & Tricks #2" 
mainPhoto: gotips-logo-02.png 
categories: Golang
tags: [golang, tipstricks]
---

When it comes to interfaces, a good practice is to create an interface where you'll use it. Creating interfaces in advanced is not recommended in Go. There are two exceptions: 

 * you're creating a library which will be used in different projects
 * you'll have more than 1 implementation

In the example below, we have a storage implementation.

```go
type inMemoryStorage struct {
   mutex *sync.Mutex
   storage map[string]*Value
}

func NewStorage() *inMemoryStorage {
   return &inMemoryStorage{
      storage: map[string]*Value{},
      mutex: &sync.Mutex{},
   }
}

func (s inMemoryStorage) Set(ctx context.Context, value *Value) error  {
   s.mutex.Lock()
   s.storage[value.key] = value
   s.mutex.Unlock()
   return nil
}

func (s inMemoryStorage) Get(ctx context.Context, key string)  (*Value, error)  {
   if val, ok := s.storage[key]; ok {
      return val, nil
   }

   return nil, nil
}

func (s inMemoryStorage) GetAll(ctx context.Context)  map[string]*Value  {
   return s.storage
}

func (s inMemoryStorage) Remove(ctx context.Context, key string) error  {
   s.mutex.Lock()
   delete(s.storage, key)
   s.mutex.Unlock()
   return nil
}
```

As you can see, we skipped the interface(s) because they are not needed here.

Why do we have this rule? Imagine the situation when you add the interface next to the implementation. The interface is used for abstracting (dependency injection). The interface can look like the below. 

```go
type Storager interface {
   Set(ctx context.Context, value *Value) error
   Get(ctx context.Context, key string) (*Value, error)
   GetAll(ctx context.Context) map[string]*Value
   Remove(ctx context.Context, key string) error
}
```

The problem with the approach comes, for example, in testing. In our production code, you may use only Set method, but you have to mock all of them. It's better to split the interface to smaller parts and define only those methods which are really needed. 

```go
type Remover interface {
   Remove(ctx context.Context, key string) error
}
```

---
title: "Services in DDD finally explained"
publishDate: 2018-07-15
---
I’ve noticed that there is always a challenge of understanding what services are in a context of domain-driven development and what is the difference between a service in an application, domain, and infrastructure layer.

Domain-driven design made a lot of cleanup in the IT environment and conquered the hearts of programmers. [Eric Evans](https://domainlanguage.com/) is one of the most famous people who promote this not so a new way of developing software. Unfortunately, from time to time naming elements of it may cause some problems with fully understanding the idea behind some parts of the DDD. One of them is a service.

I’ll explain this topic in an example – a library. In the library, you have readers who can borrow a book. The book can be available or not. It has the title, author and so on. We have domain objects like the library, the reader, books, and librarian.

Below you can find some simple example for the core domain objects. The first one is a Book.

```java
    package ddd.explained.library.domain;
    public class Book {
       private int bookId;
       private boolean borrowed;
       private boolean readingRoomOnly;
       Book(int bookId, boolean borrowed, boolean readingRoomOnly) {
           this.bookId = bookId;
           this.borrowed = borrowed;
           this.readingRoomOnly = readingRoomOnly;
       }
       boolean canBeBorrowed() {
           return !borrowed && !readingRoomOnly;
       }
    }
```

The second one implements Reader object.

```java
    package ddd.explained.library.domain;
    public class Reader {
       private int readerId;
       Reader(int readerId)
       {
           this.readerId = readerId;
       }
       public boolean lend(Book book) {
           if (book.canBeBorrowed()) {
               return false;
           }
           // some impl
           return true;
       }
    }
```

I decided to not implement the whole domain logic because it could make it less readable. The domain objects shouldn’t care about low-level details. Or even know about them.

## Domain service

Services in the domain layer are responsible for actions in this area. In our case, the service will be responsible to answer the question if the reader can borrow the book. The copies of the book could end or the reader can already have the maximum number of books or has a disabled account. The book can be only available in the reading room. As you can see, the service holds many domain objects and orchestrate them. In our case, the LibraryService service will call actions on the domain objects.

```java
    package ddd.explained.library.service;
    public class LibraryService {
       public boolean lend(Reader reader, Book book)
       {
           if (reader.canLend(book)) {
              reader.lend(book);
               return true;
           }
           return false;
       }
    }
```

## Application service

In most of the cases, we don’t only want to save changes. We want to have some interaction e. g. with a user or another service. This is the place where the application services enter the stage. Services in the application layer are responsible for fetching input data from outside of the domain, returns information about a result of the action, listens for an answer and decides if the communicational message should be sent.

You may notice the application service’s responsibility is to contact with the world. However, you should remember that in this context it’s not important where the information comes from or what exactly sends the message. It’s important that it happens.

```java
    package ddd.explained.library.application.service;
    import ddd.explained.application.service.NotificationService;
    import ddd.explained.domain.ReaderRepository;
    import ddd.explained.domain.service.BookRepository;
    public class LibraryService {
       private ddd.explained.domain.service.LibraryService libraryService;
       private NotificationService notifications;
       private ReaderRepository readerRepository;
       private BookRepository bookRepository;
       public LibraryService(ddd.explained.domain.service.LibraryService libraryService, NotificationService notifications, ReaderRepository readerRepository, BookRepository bookRepository) {
           this.libraryService = libraryService;
           this.notifications = notifications;
           this.readerRepository = readerRepository;
           this.bookRepository = bookRepository;
       }
       public boolean lend(int readerId, int bookId)
       {
           Reader reader = this.readerRepository->find(readerId);
           Book book = this.bookRepository->find(bookId);
           if (libraryService.borrow(reader, book)) {
               bookRreaderRepositoryepository->update(book);
               bookRepository->update(reader);
               notifications.sendEmailAboutLoan(book, reader);
               return true;
           }
           return false;
       }
    }
```

The ideal application service should be taken to your application, called in the proper way and all the magic should be done. The book will be booked (or not if it’s not available or other requirements are not fulfilled) and sent the e-mail with the loan details. In this context, it’s not important how exactly domain looks like or what library/vendor you use to send the email. It says what exactly have to happen – not how.
## Infrastructure service

This is the place where you’re the nearest to the plate. Services in the infrastructure layer are services which implement an interface from the domain layer. A good example is a notification service. In the domain layer, you define an interface with actions we want to have eg sendEmailAboutLoan and in this layer we implement it. This is the place what exactly mean to send the email and decide which provider we will use or bind infrastructure-specific error. This is a place where you implement repositories or anything that depends on libraries we use in the program.

You use the service in the application layer, however, your program should choose what direct implementation is more proper.
## Some explanations

This way, you can use the different implementation in different situations. Why is it so important? Imagine that our application uses a relational database like PostgreSQL or MySQL. It works like a charm but you decided to add a unit test to the ApplicationService from the application layer. If you have database related operations implemented in the same class you have a problem – you need a separate database for tests because you don’t want to mix them with your local data. In this situation, you have a few ways to resolve it:


### Create the whole database in every test you need it and drop it in the end

Pros:

* you have tests which are very similar to the production env
* you can test the whole flow in the service and be a quite sure everything works

Cons:

* they are extremely slow. Dropping and recreating tables in every test can increase the time execution dramatically
* it may be hard to debug if your class you try to test works OK but something is wrong somewhere else
* everyone who works on the project (or even CI server) has to have a separate database for tests

If your application requires some other stuff like SMTP server or another software it’ll become even more problematic.

### Mock methods or classes with direct access to the database

It’s a more elegant solution. Mocking makes you test faster and independent from the environment. In many cases, I’ve seen that these database-related methods are private or protected and it’s cool because me, as the user, I do not want to know about how it works under the hood. I’m focused on what it does. The same with an object which has these methods grouped in one place. Unfortunately, using mocks have some issues. You can read about them in articles like Why [I Don’t Use Mocking Frameworks Anymore](http://www.tddfellow.com/blog/2016/06/21/why-i-dont-use-mocking-frameworks-anymore/), [When Writing Unit Tests, Don’t Use Mocks](http://www.tddfellow.com/blog/2016/06/21/why-i-dont-use-mocking-frameworks-anymore/) or [To Mock or Not to Mock: Is That Even a Question](https://www.solutionsiq.com/resource/blog-post/to-mock-or-not-to-mock-is-that-even-a-question/)?

### Add another layer to separate our class with the environment

The last option is to separate our class from the environment and it’s what we did earlier. An extra layer help has another advantage – you’re independent on the environment. When you mock MySQL server, you know that you use MySQL in your test. While writing tests, what you use for caching isn’t important. It’s important that it works. It gives you the ability to change it every time without any impact on your domain code. Good separation helps you keep the most important part of your project completely independent of the framework you use. You can upgrade its version without worrying it will break your core code. You can even change the whole framework with simple cut&past and just integrate the whole project in only a few places. It makes the maintenance easier.

---
title: Entity and value object  
publishdate: 2018-08-04
categories: [Programming]

resources:
    - name: header
    - src: featured.png

tags:
  - ddd
  - entity
  - value object
---
Knowing the basics is the key to understanding more complex concepts. After reading this post you will know what are entities and value objects and find out differences between them.

When you pay for something at a shop it’s not important which exactly coin you choose. The most important thing to the shop assistant is their value. It does not matter if you give him coin from the left or right pocket. The key is the value.

On the other hand, when you drive a car and a policeman stops you, it’s important for the police officer to know exactly who you are. Maybe you’re a most wanted criminal? That’s why you show him your ID or driving license. Thanks to that he can identify you without any problem.

Basically, this is the difference between entities and value objects. Value object can be copied or replaced without worrying about consequences. The entity has its identity which is important in the situation.

A great example of a value object is a Money class.

```java
class Money {
    int value;
    String currency;

    public Money(int value, String currency) {
        this.value = value;
        this.currency = currency;
    }

    int getValue() {
        return value;
    }

    String getCurrency() {
        return currency;
    }

    public boolean equals(Money m) {
        return value == m.value and currency == m.currency;
    }
}
```

As long as the values are the same, both instances are identical and whichever you choose it will have exactly the same result.

```java
money1 = new Money(10, "PLN")
money2 = new Money(10, "PLN")

someMethod(money1) == someMethod(money2) // true
Money1 == money2 // true
```

Entities have their own identity. For the policeman, the ID card number directly points to a specific person and he cannot confuse with someone else. The police will not confuse two people even if they have the same name and were born at the same day.

An entity can be a user with his unique email address. This is a very popular scenario in web applications. However, it does not have to be only 1 parameter which undeniably identifies the entity.

Imagine you have many instances of an application run separately. Every instance has his own infrastructure including the database. In such situations, the user’s id will repeat at different markets. To identify only one specific user you have to take both user’s id and market code.

When you buy products at an online shop you receive a unique order’s id which identifies your specific order which can be used while complaint etc.

```java
class Order {
    ProductsList products;
    String number;

    public Order(String number) {
        this.number = number;
        products = new ArrayList<>();
    }

    public addProduct(Product p){
        products.add(p);
    }

    public boolean equals(Order order) {
        return number == order.number;
    }
}
```


Worth noticing fact is how the equals functions look like. In value objects, we compare all the parameters and this is the only way to say if they are equal or not. If they differ with at least one value they are not equal anymore.

```java
order1 = new Order("123")
order2 = new Order("123")
order1.addProduct(new Product)

order1 == order2 // true
```

On the other hand, while working with entities, we need to compare only this values which identify the entity in your system. When you buy a ticket to a movie you can be treated as an entity because for the cinema it’s important who you are. A spot in the cinema hall you bought can be an entity too because when you buy a specific spot it should be reserved for you. To identify the spot you have a few information on the ticket eg cinema hall number, the line number, and the seat number.

## Summary
As you can see both types of object are very simple and similar but have one huge difference – the way we can treat them while comparing and identifying. While it’s not important which exactly value object we choose. On the other hand, for entities uniqueness is very important it is their key feature.

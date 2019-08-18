---
title: Mutational tests
publishdate: 2018-03-24

---
When you have a very simple application it’s not so important to test every edge case but in a project, in the very complex domain, the priority of it will increase. The more high-quality the tests, the more high-quality the code. Mutational tests will help you with making sure you did not miss some a variant of the flow in your code.

## How it works


* Mutational testing is a test which runs other test several times but with a bit changed production code in every iteration. Everything happens on the fly. The idea behind it is to check if your test will fail if you change something in your production code.
* Runs Tests without modifications. The tests should be green here.
* Parses the whole production code to find places where it can make a change.
* For every small change runs tests again and one of your tests should fail.
* If one of the mutations fail, the whole mutational test fails, too.
* Otherwise, the test succeeded.

I’ll show it in an example when it may be useful. Imagine you have a code like the above.

```golang
    function calculateTotalPrice($products = [])
    {
        $totalPrice = 0;
        foreach($products as $product)
        {
            $totalPrice += $product->getPrice();
        }
        return $totalPrice;
    }
```

You wrote a test to make sure it works. The tests are green. Cool, isn’t it? When you change 0 to another number like 1 or you change plus to minus, your tests will fail. You go to drink another cup of coffee. You’ve done your job.

But life goes on… We change. Our code changes, too. Imagine that someone else got a task to add support for discounts. What a great feature! The developer (let’s call them “the bad guy”) changed your code to the above and did not write a test! In CI the tests are still green but if you are observant you’ll see something’s wrong. The bad guy is going for his coffee but the CI never sleeps.

    function calculateTotalPrice($products)
    {
        $totalPrice = 0;
        foreach($products as $product)
        {
            $price = $product->getPrice();
            if ($product->hasDiscount())
            {
                $price = $price + ($product->getDiscount()/100) * $price;
            }
            $totalPrice += $price;
        }
        return $totalPrice;
    }

Let me remind you how mutational tests work — they change your production code a bit and check if your tests fail. In this example, in line 10 the plus sign will change to minus or multiplication will be replaced to a division. A mutational framework will detect this not tested cases and let us know that’s something to cover.

## Advantages of mutational testing

In short, they show you where you should do the second thought. If it found a problem while mutating, maybe you have some unused variable or some part of the code is not covered by tests or it’s not covered the way it should be. It helps you to make sure your code is more resistant to changes.

## Disadvantages of mutational testing

As everything, mutational tests have disadvantages. First of all, they are slow. Very slow. They run all test every time they change something in production code on the fly. In a large project with lots of tests, it may take hours.

Secondly, a framework I tested found some stupid “bugs” which you’ll have to resolve with… more tests. In some situations, it may be cumbersome.

Finally, it’s hard to introduce them in large, already existing, projects. Why? Imagine you have a project with about 2k lines of code and mutational framework found 100 errors. It means you have to write about 100 tests to cover the edge cases.

In almost every popular programming language you can find a mutational framework/tool for mutational testing. Here are some of them:

* python — https://github.com/sixty-north/cosmic-ray
* java — http://mutation-testing.org/
* php — https://infection.github.io/
* .net — https://visualmutator.github.io/web/
* javascript — https://stryker-mutator.github.io/

## Summary

As everything, it may help or maybe a loss of time. If a high-quality code is extremely important to you, mutational testing is what you should test. Have you used tools like that? What is your experience? Tell me in the comments.

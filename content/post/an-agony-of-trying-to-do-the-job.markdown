---
title: An agony of trying to do the job
publishdate: 2017-12-27

---
I’ll tell you a story of Igor. Igor is a web developer. He’s a young man with a girlfriend and some ambitious plan in the future. Igor sit at his desk because he has some work to do. In front of him is a PC. He turns it on and sees some system updates. 30 minutes have passed, and he can see his desktop with a beautiful wallpaper. Funny cats always make his day better. Some people are walking around his working desk, talking and drinking coffee. A typical day.

The IDE is on, but he cannot start working immediately because it’s indexing your project. It takes two extra minutes. Igor is okay with it – more time to finish the coffee is always a good idea.

In the meantime, Igor checks what exactly he has to do. Change a footer of a newly released project? What can be easier? – He thought. Igor opens a file where he has to make the change. Text replaced. File saved.

The boy wants to check how it looks now to push the fix further. You cannot even imagine a surprise on his face when he noticed that the project does not compile.

It worked a week ago! – Shouted. What does the log say? Some NodeJS module cannot be installed because of missing a system library. Someone added a module to compress images into a build step.

Igor asked his friend, Steven, who worked with him on the project and asked why he has to install some system libraries? It how NodeJS works – Steven said. – You know NodeJS – sometimes he has some strange dependencies. Just install the package and go further.

Steven was right. Installing the missing library fixed the problem. Everything looks OK now, and it’s time for `git commit` and `git push`. It failed… `git push` command was rejected because of a problem with ssh. The admin changed some keys, and he had to update a file in his `~/.ssh` folder. One of the admins came to Igor and tried to explain why it happened. Igor did not understand any word he was saying. He was talking in some different languages. Nevermind, it works now.

The change is on a build server, so it’s time for second breakfast. After 15 minutes Igor comes back, and he was almost sure that it was built and ready for production. Nothing more wrong. Build agents are down because of a regular update which disabled all the agents. The deployment was impossible! How long will it take? Some admin said it should be working before lunch.

A small change could not be pushed in few minutes because of many minor problems. It’s a simple story just to show you how many time we spend on solving problems not with the project itself but with tools we use every day and which should make our life easier.

Try to imagine it’s your first day of a new job. Can you just copy the project to your newly installed system and just start working on it? How many hours or even days you need to set up the whole environment?
## Summary

I am curious how much time you devote to fighting the tools and what are your most common problems. Are they operating system problems or maybe hardware issues? Write in comments!

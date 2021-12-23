---
title: "History of WWW"
publishdate: 2018-05-02

resources:
    - name: header
    - src: featured.png
---
The history of the WWW wasn’t simple and still isn’t. In the beginning, it was a complete mess. Every browser tried to meet developers halfway without any standards or cooperation with other web browser’s developers.

It all started with Memex – a theoretical machine created by Vannevar Bush with the ability to create links between documents and books and add comments to them. Until early 90’ Memex was derided and described as unreal because we had no tools to implement this idea. But then modern computers appeared…

## The beginnings

Let’s go back to 1969-70 when IBM took a couple macros and created a language called GML (Generalized Markup Language). It used a colon to call a given macro. It looked like something along these lines:

    :h1.Chapter 1: Introduction
    :p.GML supported hierarchical containers, such as
    :ol
    :li.Ordered lists (like this one),
    :li.Unordered lists, and
    :li.Definition lists
    :eol.
    as well as simple structures.
    :p.Markup minimization (later generalized and formalized in SGML),
    allowed the end-tags to be omitted for the "h1" and "p" elements.

GML was an underline for another language called SGML which was more standardized. The main difference between these two languages is that SGML uses angle brackets instead of colons. Sample SGML documents are very similar to today’s HTML or XML documents:

    <QUOTE TYPE="example">
    typically something like <ITALICS>this</ITALICS>
    </QUOTE>

Both GML and SGML were created to help to display documents in a more user-friendly way e.g. with bolds, underlines, lists and so on. But there was a problem – there was no easy and scalable way to create links between the documents. There were some attempts to achieve it but they did not pass the test of time. One of them was ENQUIRE which was more like wiki pages with bi-directional hyperlinks. If you are interested in this subject, you can read the [ENQUIRE manual with scans of the original manual from 1980](https://www.w3.org/People/Berners-Lee/EnquireManual.htm)!

## The birth of HTML and HTTP

![Tim](/assets/posts/tim.jpg)

### Tim Berners-Lee

At the end of 80’ two scientists, Tim Berners-Lee and Dom Connolly, started working on a new language called HTML (Hypertext Markup Language). HTML allowed to create hyperlinks between documents and add to it some text-formatting functionality. Next, they started working on HTTP (Hypertext Transfer Protocol) which was an extremely simple way to serve HTML documents (prior solutions would display an icon which when clicked, would download and open the graphics file in an external application). Another project he was working on was URI (Uniform Resource Identifier). It was a unique identifier for server-hosted resources. It used the widely adopted approach of grouping documents in folders which worked very similarly to directory structures in Unix systems.

HTML and HTTP were easier and simpler than their predecessors – their negligible learning curve helped them take over the crawling Internet’s developing need for easily accessible content…

### First WWW browsers

In 1993 the HTTP protocol contributed about 0.1 percent of total network traffic in America’s National Science Foundation. In the same year, a new web browser was created. It was called Mosaic which used both HTML and HTTP protocol and was graphical. Mosaic had some unique features like embedding images or sending forms. On the foundation built by Mosaic, the next generation of browsers was born. Two of them are Mosaic Netscape (called later Netscape Navigator) and Spyglass Mosaic, later bought by Microsoft and called Internet Explorer.

![Mocaic Netscape](/assets/posts/mosaic.jpg)

![Internet Explorer 1](/assets/posts/ie1.png)

There were other projects of web browsers completely independent from Mocaic, e.g. Opera which was first publicly released in 1996 but works on the browser started in 1994.

However, they were not the very first projects which helped to browse WWW. The first ever WWW browser was called WorldWideWeb and was created in 1990 by [Tim Berners-Lee](https://en.wikipedia.org/wiki/Tim_Berners-Lee). Other similar projects were [MidasWWW](https://en.wikipedia.org/wiki/MidasWWW) and [ViolaWWW](https://en.wikipedia.org/wiki/ViolaWWW). MidasWWW’s source code is available in [GitHub](https://github.com/dckc/MidasWWW)!

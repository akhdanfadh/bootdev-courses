# Learn the HTTP Protocol in Go

Or, as I like to call it, **HTTP from TCP**.

In this course, we'll **build our own HTTP 1.1 server from scratch** in Go.

<video src="https://storage.googleapis.com/qvault-webapp-dynamic-assets/lesson_videos/http-intro.mp4" controls=""></video>

## Prerequisites

- Understanding of the Go programming language. If you're fuzzy on that, [have we got a course for you](https://www.boot.dev/courses/learn-golang).
- What [binary](https://en.wikipedia.org/wiki/Binary_number) is, which is covered in our [Learn to Code in Python](https://www.boot.dev/courses/learn-code-python) course.
- Understanding of how to _use_ HTTP, because in this course we want to _build_ it. We have an [HTTP Clients](https://www.boot.dev/courses/learn-http-clients-golang) course and an [HTTP Servers](https://www.boot.dev/courses/learn-http-servers-golang) course that can help with that.

## Learning Goals

1. Understand the big ideas of HTTP, and implement the `HTTP/1.1` protocol. (Look, the RFC is long, but the main points are fairly straightforward). You'll understand from a high level what happens when you `fetch("google.com")`.
2. To feel like a wizard, building the protocol of the internet from scratch.
3. To gain more challenging programming practice, and a deeper understanding of how web applications work.
4. To have a good time.

## Story Time

I still remember sitting in class, fall of 2007, as I watched my first data structure coded in Java. The magic I felt seeing a class reference itself...

```java
class Node<T> {
    public Node prev;
    public Node next;
    public T data;
    ...
}
```

In that moment, I knew I was a computer scientist through and through. This feeling reached a pinnacle when I had to create a video transfer protocol for a government robot.

- **Packet ordering**? Nope.
- **Reliable transport**? Ackshually yes. (pikachu shocked face)
- **Corrupt Data**? Regularly.
- **But did _I_ create it**? You're darn tootin'.

I felt empowered. And the best part? _It was all in C_.

Were there security issues? Yep. Memory leaks? Probably. I wouldn't know. Did it successfully transfer video while the other teams struggled to get basic hardware components to even work? Yes. I felt like a wizard. And I want you to feel like a wizard too.

Welcome to TCP to HTTP.

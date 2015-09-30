Interfaces in Go
================

In this example, you end up with a session manager and a database manager
that need to use each other, but Go doesn't allow cyclic package dependencies
- that is, package "a" can't import package "b" if package "b" imports 
package "a". Go does have interfaces, and what's interesting about them is
that they're implicit. That means that you can interface out a third-party
structure without modifying its package. Just define an interface in your
package with the methods you need to access on that structure, and treat
the structure as your interface type.


Example
-------

In this (contrived) example, we have a `datastore` package containing a `SQLDatabase`
structure that needs some info from the `session` package's `RedisSession`
struture. However, the `RedisSession` has a couple cases where it needs
information from the `SQLDatabase`. You get an error about cyclic package
dependencies, and immediately shake your fist at the language. The thing is,
the language is telling you that there's a 
[code smell](https://en.wikipedia.org/wiki/Code_smell). So, let's clean up
by taking advantage of Go's implicit interfaces. You'll notice neither package
imports each other.

prj - Simple Filesystem Project Tool
====================================

What?
-----

This is a tool to mark a directory as containing a project and to generate a
hash of its contents. Hashes can be marked at a point in time and a list of
files changed since a particular mark can be queried.


How?
----

Easy!

    go get -u github.com/shabbyrobe/prj/cmd/prj

Initialise:

    cd /path/to/your/thing/that/is/a/project
    prj init

Modify:

    touch foo
    prj status
    prj mark -m 'Touched foo'

Modify again:

    rm foo
    prj mark -m 'Removed foo'

Show me the log:

    prj log

Show me what I've changed:

    touch pants
    prj status

Show me all the projects in all descendents of the current folder:

    prj find
    prj find "$(pwd)" # long-form equivalent


Why?
----

I've got a lot of software projects on my hard drive. They're easy to identify
because they almost always contain a `.git`, `.hg` or `.svn` folder, or some
other easy-to-identify way that this folder is the root of a "project".

I've also got nearly two decades worth of other, miscellaneous projects on my
hard drive. They're hard to identify, harder to compare, they're strewn across
various backups in various states of disorganisation, and it's very hard to sort
out what may have changed since it was put on that hard drive to go back into
the studio to do more tracking back in 2010 then put back into a folder called
"unsorted", which is the same "unsorted" folder we put it into last time we did
the same dance in 2008. Help!

> Why don't you just be more organised?

The true answer to this question is, of course, n-gate.com, but yes, that has
definitely been tried. Obviously it hasn't worked, because here we are!


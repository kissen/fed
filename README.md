fed
===

Trying to build a server for [ActivityPub](https://www.w3.org/TR/activitypub/)
that can handle all objects defined in the standard. Development happens in
the `experiment/` branches, `master` doesn't even build right now.

What I Would Like
-----------------

Current services like [Mastodon](https://joinmastodon.org/) or
[Pleroma](https://pleroma.social/) do a bit too much for my liking and
are difficult to install.  I would like `fed` to be easy to deploy and
self-host. The core service should be without any user interface,
leaving it up to the user how to interact with it.


I'm not sure if I will ever create something working. I do know that I
want this kind of software; hosting a private ActivityPub instance for
your own needs to way easier than it is right now.

License
-------

The code in this repository contains comments copied from the
[go-fed/activity](https://github.com/go-fed/activity) project which is
licensed under a BSD 3-Clause License, *Copyright (c) 2018, go-fed*.

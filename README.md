# fed

Trying to build a server for [ActivityPub](https://www.w3.org/TR/activitypub/)
that can handle all objects defined in the standard. Development happens in
the `experiment/` branches, `master` doesn't even build right now.

* `fedd` contains code for the backend so to speak, the federation daemon.

* `fedweb` contains code for a basic web interface that speaks to
  `fedd`. Ideally it should be able to speak with any ActivityPub
  service.

* `fedutil` contains shared code between `fedd` and `fedwweb`. In
  particular, it contains convenience functions that work on types
  from the `go-fed/activity` library.

## What I Would Like

Current services like [Mastodon](https://joinmastodon.org/) or
[Pleroma](https://pleroma.social/) do a bit too much for my liking and
are difficult to install.  I would like `fed` to be easy to deploy and
self-host. The core service should be without any user interface,
leaving it up to the user how to interact with it.

I'm not sure if I will ever create something working. I do know that I
want this kind of software; hosting a private ActivityPub instance for
your own needs to way easier than it is right now.

## Credit

Even this small prototype wouldn't work with the help of many open
source projects. Some are directly included in this repository.

### go-fed

The code in this repository contains comments copied from the
[go-fed/activity](https://github.com/go-fed/activity) project
licensed under the following terms.

	BSD 3-Clause License

	Copyright (c) 2018, go-fed
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	* Redistributions of source code must retain the above copyright notice, this
	list of conditions and the following disclaimer.

	* Redistributions in binary form must reproduce the above copyright notice,
	this list of conditions and the following disclaimer in the documentation
	and/or other materials provided with the distribution.

	* Neither the name of the copyright holder nor the names of its
	contributors may be used to endorse or promote products derived from
	this software without specific prior written permission.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

### Feather

This code contains icons the [Feather](https://feathericons.com/) icon set
licensed under the following terms.

	The MIT License (MIT)

	Copyright (c) 2013-2017 Cole Bemis

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

The affected files, some of which were modified, are

	like-active.svg  like.svg  repeat-active.svg  repeat.svg
	reply-active.svg  reply.svg external.svg send.svg inbox.svg
	warning.svg error.svg

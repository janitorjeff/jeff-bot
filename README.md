# JanitorJeff
A general purpose, cross-platform bot.

*Very much under active developement, there will most likely be breaking changes.*

## Structure

There's 3 main components to the bot:

The [frontend](https://git.slowtyper.com/slowtyper/janitorjeff/src/branch/main/frontends)
layer is responsible for providing methods for working with a specific frontend
(receiving messages, sending messages, etc.). An
[interface](https://git.slowtyper.com/slowtyper/janitorjeff/src/branch/main/core/message.go#L20-L60)
exists in the core and for a frontend to be added it must implement that
interface.

The [message](https://git.slowtyper.com/slowtyper/janitorjeff/src/branch/main/core/message.go#L74-L85)
layer is responsible for creating a common struct under which all messages from
all frontends are processed.

The [command](https://git.slowtyper.com/slowtyper/janitorjeff/src/branch/main/commands)
layer where all of the actual commands are implemented. For most use cases it
is frontend agnostic, an exception is discord because of the special
"rendering" the messages get (embeds as opposed to plain text).

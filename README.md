discordgo 2 POC
===============

This is a proof-of-concept of a proposed architecture for a rework of the [discordgo](https://github.com/bwmarrin/discordgo) library, taking into account various complaints, suggestions and concerns that have been raised in the associated discord channels. It also contains my own attempts at making the API more idiomatic.

Please note that this is a rough proof of concept, published with intent to spark discussion, from which something more usable can be created. There will be sharp corners, bugs aplenty and design choices that on second thought could be better.

Architecture
============

Separation of concerns.
-----------------------

Currently, Discordgo has a single Session type that wraps the REST-, WS- and Voice APIs in one. While convenient in some ways, this has more than once resulted in confusion as to which methods actually touch what, and it results in a lot of complexity being concentrated in one place.

This POC splits it into three separate pieces: Client, WSClient, and VoiceConn. Each one depends on the one before it, but never the other way around.

Interfaces and mocking.
-----------------------

Speaking of, Client, WSClient and VoiceConn are all interfaces. This means we can implement mocking for unit tests in a sane manner.

OAuth2 and tokens.
------------------

A Session is created with a call to `New(token)`. Innocious as it may be, I believe there's major room for improvement.

The current `New()` actually takes a varadic `...interface{}`, but will error on anything but a `string` (or, for some reason, `[]string`) is passed in. It also has three separate behaviours depending on if how many parameters were passed:

1. Use what's passed (a token) verbatim for the Authorization header.
2. Use what's passed as an email and password, and make a login call to Discord to retrieve a token.
3. Use args 1 and 2 to login if the token (3) isn't valid.

This actually historically made sense: back before we had proper applications and bots, we ran out bots using regular user accounts, and for those, Discord recommends (recommended?) an equivalent of (3). At the time, they also hadn't quite made up their mind on whether or not tokens should expire, among other things.

That said, convenient as it may be, the single most commonly used function in the library is now one whose arguments can't be reasoned about from its definition, and where a blocking network call can be silently made from a constructor.

To make matters worse, the introduction of the "Bot" prefix for bot users (without which you get read-only access) means that one of the single most common questions that pop up on Discord is "Why do I get a 403 whenever I do anything?". While not the API's fault (Discord's own docs are completely useless in this regards, fueling the confusion), I believe the API could be made clearer.

This POC's `NewClient()` function takes a `*oauth2.Token`, from the official [`x/oauth2`](https://godoc.org/golang.org/x/oauth2) package. It also provides two functions to wrap the two different kinds of tokens: `UserToken(string)` and `BotToken(string)`.

The x/oauth2 library also provides functionality for authenticating with an OAuth2 server, which means we can make this step explicit, as well as support other forms of authentication than email/password (eg. three-legged OAuth2 for webapps).

Future-proofing and method bloat.
---------------------------------

Discordgo also has historically had problems whenever Discord change their API to add more functionality. Over the years, discordgo has accumulated all of these functions for sending messages to a channel:

* `ChannelMessageSend(channel, content string)`
* `ChannelMessageSendTTS(channel, content string)`
* `ChannelFileSend(channel, name string, r io.Reader)`
* `ChannelFileSendWithMessage(channel, content, name string, r io.Reader)`
* `ChannelMessageSendEmbed(channel string, embed *MessageEmbed)`
* `ChannelMessageSendComplex(channel string, data *MessageSend)`

The last one, taking a struct with every single possible field, is the logical response to this growing amount of bloat. A perfectly reasonable course of action for a v2 would be to cut most (all?) of the bloat and keep Complex as the only API, possibly taking a lesson from the past and using this idiom in a more widespread manner.

However, this would turn the most common use case - `ChannelMessageSend(cid, "hi!")` - into something rather bulky (but manageable!): `ChannelMessageSend(cid, &discordgo.MessageSend{Contents: "hi!"})`.

An idiom that has popped up recently to combat this was [first published by Dave Cheney](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis): "functional options", and this paradigm has since been picked up by everything from the [go-micro toolkit](https://godoc.org/github.com/micro/go-micro/server#NewHandler), to Google's official [gRPC implementation](https://godoc.org/google.golang.org/grpc#ClientConn.Invoke). Essentially, this is a way of combatting Go's lack of optional parameters by using varargs and "modifier" functions.

Thus, instead of having to instantiate a `MessageSend` for every message, or having multiple functions like status quo, you'd have something like this:

```go
// Send a message.
c.ChannelMessageSend(cid, "hi!")

// Send a message with an attachment.
f, err := os.Open("myfile.jpg")
c.ChannelMessageSend(cid, "hi!", SendWithFile("myfile.jpg", f))

// Send a message with embeds.
c.ChannelMessageSend(cid, "hi!", SendWithEmbed(...), SendWithEmbed(...))
```

Behind the scenes, these options would be simple functions along these lines:

```go
type SendOpt func(*MessageSend)

func SendWithFile(name string, r io.Reader) SendOpt {
    return SendOpt(func(s *MessageSend) {
        s.Files = append(s.Files, &File{Name: name, Reader: r})
    })
}
```

Message Handlers.
-----------------

One of the major criticisms against the current API is the `AddHandler()` function. It takes an `interface{}` argument and uses a code-generated type switch to correctly slot it in - improved from the original implementation, which used reflection. An incorrect signature simply gets you a log message telling you so, at runtime; there are no compile time checks of any kind.

Using the same approach as above, we can make a type-safe AddHandler API:

```go
c.AddHandler(OnReady(func(r *Ready) {
    // ...
}))
```

The implementation is something along these lines:

```go
type wsHandlers struct {
    Ready []func(*Ready)
}

type WSHandler func(*wsHandlers)

func OnReady(fn func(*Ready)) WSHandler {
    return WSHandler(func(hls *wsHandlers) {
        hls.Ready = append(hls.Ready, fn)
    })
}

type WSClient struct {
    // ...
    handlers wsHandlers
}

func (c *WSClient) AddHandler(hl WSHandler) {
    hl(&c.handlers)
}
```

The actual thing is code-generated and somewhat more complex, involving proper locking and a way to remove handlers once added, but the basic semantics remain the same.

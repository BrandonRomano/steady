# Steady

Server generated static sites.

## What is this?

In short, Steady is a tool that will allow you to build a server side application (written in Rails, Django, Sinatra, Laravel, Revel, etc.) and compile it down to be a static site.

This site generator works entirely over HTTP, so if your server side application can communicate over HTTP (it can) Steady can work with your setup.

## Steady vs Traditional Static Site Generators

The first question you might ask is "Why this over a more traditional static site generator?".  There are many answers to this question, but I will outline a few.

- Traditional static site generators don't scale well with large amounts of content (long compile times), this is particularly painful in development.
- Instant previews of data changes (editors will thank you).
- If business needs change, you can cut out Steady from your build process and switch back to a traditional SSA hosting strategy (without changing a line of app code).
- Generally, you will be able to build a smaller API (or none at all) as a server side app can safely access the database directly.
- There are lots of development advantages to server side apps (I will avoid specifics, as this is a large topic)
- Your team might already be comfortable with a server side app stack.
- Your project might already be a server side app and you want to host it as a static site.

## Installation

TODO

## Setup

There is a small amount of setup involved with getting Steady up and running.

You will have to generate the endpoint `/steady-files`, which will return an object with a key named `files` which is an array that points to all of the pages your site can generate, as well as all of your assets.

For cach file object, you can specify:

- `url`: The absolute path to the resource
- `type` (optional): The expected file extension if no extension is provided in the `url`.

Here's an example of what you should output:

```js
{
  "files":[
    {
      "url":"/posts",
      "type":"html"
    },
    {
      "url":"/posts/hello-world",
      "type":"html"
    },
    {
      "url":"/posts/foo-bar",
      "type":"html"
    },
    {
      "url":"/assets/index.css"
    },
    {
      "url":"/assets/index.js"
    }
  ]
}
```

## Compiling

Before you compile, you'll need to have a running instance of your server application (either locally or hosted).

Assuming that your server running at `http://localhost:5000`, you can run the following command:

```
steady http://localhost:5000
```

After this command has finished executing, your static site will be located in your current directory at `./public/`.

## Caveats

Most server side frameworks allow you to handle form submissions in browser which result in a direct request to the server.

Static sites will not support this type of interaction. You will have to make sure all of your client side interactions that need to speak to a server are performed via AJAX requests to an external API.  Keep in mind that this API can be the same codebase, but it will have to be hosted separately from your static site (and referenced as a remote resource in your requests).

## License

[MIT](LICENSE.md)

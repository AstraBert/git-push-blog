# git-push Blog

**git push** is a blog space for everyone who wants to share their journey, tips and tricks in coding.

It's Open Source, built around this GitHub repository and with a precise mission: create an open and inclusive environment where everyone can help others grow and thrive in the awesome land of programming.

## Usage as a blog template

You can use this repository as a template to build your own blog. In this sense, you will need to have go 1.24.5+ installed on your machine.

Apart from this requirement, using the template is simple, as the repository is structured as follows:

- in `contents` you can place the markdown files that make up your blog. If you want a file to be included as source for the blog, you have to name it `{number}.md` (where `{number}` is a positive integer number). The file will then have to contain the following header:

```md
---
title: your title here
publishing_date: publishing date here
author: author name here
---
```

Your header should always follow this format, otherwise there might be some unexpected behaviors.

- in `templates`, you can find the HTML templates written with [Templ](https://templ.guide). They are based on HTML and [Tailwind CSS](https://tailwind.com), and you can adjust them as you wish. Templ uses a Go-like syntax for displaying some content dynamically, so a little knowledge of Go would be good here, unless you just want to change the style or the structure of the HTML document.
- In `models` there are the data models and the associated helper functions
- In `page_reader` there are the utilities to transpile markdown to HTML code and sprinkle some CSS styling over the converted HTML.

> _If you are simply planning to adopt this template as a blog and simply adjust the style, you do not need to change anythng in the `models` or `page_reader` folder._

`main.go` is the entrypoint for the application, and provides the web server architecture.

In order to run the blog locally, you can then simply execute the main file:

```bash
go run main.go
```

Or use a Docker build:

```bash
docker build . -t your-name/git-push-blog
docker run -p 8000:8000 your-name/git-push-blog
```

As the `docker run` command suggests, you will find the local server running on `localhost:8000`.

If you are interested in deployment solutions, services like [Koyeb](https://koyeb.com) (the one that hosts the git-push blog) generally provide easy and fast solutions, especially centered around Docker.

## Contributing

Contributions (both for the blog and for the source code) are more than welcome! You can find a detail contribution guide [here](./CONTRIBUTING.md).

## License

This project is distributed under [MIT license](./LICENSE).

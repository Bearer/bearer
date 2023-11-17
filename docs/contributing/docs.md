---
title: Contributing to documentation
---

# Contribute to Bearer CLI's documentation

If you're interested in contributing to Bearer CLI's documentation, this guide will help you do it. If you haven't already, review the [contributing overview](/contributing/) for different ways you can contribute.

## Contribute directly on GitHub

If you're making a change to an existing page, or fixing something small like a typo, you can do so directly on GitHub. Select the "Edit this page" link in the right column on any page in the docs to edit the page's source.

## Setting up local development

The documentation is built with [eleventy](https://www.11ty.dev) and [tailwindcss](https://tailwindcss.com/), and written primarily in markdown. To get started, you'll need NodeJS 12 or newer.

Once you've forked and cloned the repo, navigate to the docs directory, and `npm install` to bootstrap the dependencies.

```bash
# github.com/bearer/bearer/docs
npm install
```

With dependencies installed, you can start a local dev server by running the `start` command:

```bash
npm start
```

Your local setup will cache some requests, so if you find for instance that the latest rules aren't pulled in, you can clean the cache and local temp folder with the `clean` script

```bash
npm run clean
```

## Source structure

As mentioned previously, The Bearer CLI docs run on [11ty](https://www.11ty.dev/), a static site generator written in JavaScript. For most doc contributions, you won't need to worry about the inner workings, but if you're adding something more complex this overview may help. Documentation lives in the `/docs` folder of the [bearer/bearer](https://github.com/bearer/bearer) repo. The directory structure is as follows.

### Source files

Source files make up the workings of our documentation site. They include data sources, styles, layouts, and more.

- `.eleventy.js`, `.eleventyignore`, `11tydata.json`: Configuration files for 11ty.
- `_redirects`: Whenever a breaking change to a path is made, make sure to acknowledge it in the redirects file. Netlify uses this to set redirects on the server.
- `_data`: The source for all derived and shared data for documentation page. Some files, like `bearer_scan.yml` come directly from the CLI build, while others like `rules.js` convert source data into JSON that our page templates use.
- `_src`: This directory includes partial files, plugins, javascript, and styles used by the pages. For example, the `_includes` folder has page templates, icons, etc.
- `assets`: Any images, shared files, etc. live in the assets folder. They are then copied to the root for use in the live docs.

### Pages

Pages are written in markdown or nunjucks. See the [creating a new documentation page](#creating-a-new-documentation-page) section below for details on creating a new page.

- `docs.md`: The docs homepage.
- `quickstart.md`: The quickstart guide.
- `/guides`: The guides directory. Guides are tutorial-style pages that walk through how to perform a task or tasks.
- `/explanations`: Explanations are longer, more detailed reasoning for specific concepts.
- `/reference`: Reference pages primarily use the CLI's source as a foundation to display reference data in a nicer way.
- `/contributing`: All docs related to contributing to the project live here.

## Creating a new documentation page

To create a new page, first confirm that an existing page doesn't serve the same purpose. In general, **it's best to enhance an existing page** rather than create a new one. If you do need to create a new page, start by determining the location within the hierarchy. See the _Pages_ section above for details on each section.

For traditional static documentation, use markdown and create a `.md` file in the selected category (guides, explanations, reference, contributing). Include frontmatter data for the `title` at minimum. This is used to create the page title in HTML. Then, starting with a top-level heading, write your doc in markdown. Review similar documentation pages for examples of how to create different elements.

```md
---
title: This is the html and seo title
---

# This is the title seen on the page
```

If your page needs to be more dynamic or pull in data at build time, use the `.njk` nunjucks format. This templating language is used throughout the docs.

Finally, add the new page to the navigation. The nav template is built from the `_data/nav.js` data file. Find the right section, and add your new page. In addition, you may want to add a link to the main `docs.md` page.

### Page tips

Here are some common solutions you may need when working with pages.

#### Table of contents

A table of contents is auto-generated for each page using the headings (h2-4).

#### Using markdown in nunjucks files.

If you want to write markdown in a nunjucks file, the `renderTemplate` shortcode makes this possible. See [11ty's Render docs](https://www.11ty.dev/docs/plugins/render/) for details.

#### Using images in markdown

To add an image to a page, first optimize it and add it to the `docs/assets/img` folder. From markdown, use the final output location to reference the image. Note that in the example below, we reference the root path.

```md
![This is alt text](/assets/img/example.png)
```

#### Building data-driven pages

Some pages in the documentation, like the [rules](/reference/rules/) or [data types](/reference/datatypes/) pages, are data-driven. They are still static, but rely on a data source at build time to create the HTML. This is done by pairing a `data` file (.js) with a `page` file (.njk). The [11ty docs](https://www.11ty.dev/docs/data-js/) has explanations on how JavaScript data files work.

In addition, the [data types](/reference/datatypes/) page should act as a good example. The process looks like this:

1. At build time, `docs/_data/datatypes.js` reads source files from Bearer CLI.
2. `docs/reference/datatypes.njk` looks for the global `datatypes` object created in step 1 and uses it to populate the template.

## Getting help

If you're unsure where to start, have questions, or need help contributing to the documentation, join our [community discord]({{meta.links.discord}}) and we'll be happy to help out.

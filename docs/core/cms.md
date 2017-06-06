# Flamingo CMS Bundle

## CMS Pages

This bundle registers `cms.page.view` for CMS pages.

The default URL is `/page/:name`

## CMS Blocks

CMS Blocks are usable in pug templates via

```pug
h1 #{get("cms.block", {"block": "blockname"})}.title
p !{get("cms.block", {"block": "blockname"})}.content
```

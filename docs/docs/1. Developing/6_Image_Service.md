# Image Service Usage in Templates

Example Usage in Templates:

```
img(src=image("pim","1200x1100","catalog/0/0/0/0/00003b92d2702b3513749e53aacfdd699675cc13_product_image_595fab5992ced.png"),width="500", height="500")
```

## Parameters explained

The `image` Method accepts three parameters: Source, Options and an Image Path.

### Source

Source for Images is set to either :

 - "pim" for Akeneo
 - "mdp" for Master Data Portal
 - "cms" for Magento CMS
 
The Applications you´re requesting Images from may have differing Storage Backends (AWS S3, Filesystem, Database, etc.) which the Image Service hides from you.

### Options

The Imageservice uses the same Parse Options as Willnorris ImageProxy. You can find a documentation on how to chain Parse Options here:
https://godoc.org/willnorris.com/go/imageproxy#ParseOptions

Examples
```
0x0         - no resizing
200x        - 200 pixels wide, proportional height
x0.15       - 15% original height, proportional width
100x150     - 100 by 150 pixels, cropping as needed
100         - 100 pixels square, cropping as needed
150,fit     - scale to fit 150 pixels square, no cropping
100,r90     - 100 pixels square, rotated 90 degrees
100,fv,fh   - 100 pixels square, flipped horizontal and vertical
200x,q60    - 200 pixels wide, proportional height, 60% quality
200x,png    - 200 pixels wide, converted to PNG format
cw100,ch100 - crop image to 100px square, starting at (0,0)
cx10,cy20,cw100,ch200 - crop image starting at (10,20) is 100px wide and 200px tall
```

### Image Path

There will probably rarely a usecase for hardcoded Image Paths, they will most likely be fetched from an API
(Products, Categories, Brand Response from the Master Data Portal, etc.)

The Path is always a relative Path and does not begin with a Slash.

### Image Format

Even after applying Filters according to the Options Parameter, the Image Format stays the same, an image/jpeg still has the
same mime Type.

### Security

The Image Helper Method will sign the request to the image service with a SHA256-HMAC Signature automatically.

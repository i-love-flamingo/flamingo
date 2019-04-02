# Package: Pugtemplate


## About
Main PUG Template Enging.

Provides the following template Funnctions:

* get: Give access to DataControllers in Flamingo. Example get("retailer", "retailerCode")
* priceFormat: Formats price. Example: priceFormat(4.666, "EUR")
* debug


## Configuations:

### Configure the priceFormat
```
accounting:
  decimal:    ','
  thousand:    ','
  formatZero:    '%s -,-'
  format:    '%s %v'
``

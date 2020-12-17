# Locale module (Localization)

This package provides localization features:

 * Translations of Labels: with the template func `__()`  (which uses [github.com/nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n))
 * Local display of Dates (In local timezone or given timezone): with the template func `dateTimeFormat` or `dateTimeFormatFromIso`
 * Local display of prices and numbers: with the template func `priceFormat` and `numberFormat` (Using [github.com/leekchan/accounting](https://github.com/leekchan/accounting))

## Configuration

```yaml
locale:
  locale: en-gb                                    # the locale used for labels
  fallbackLocales:                                 # a list of (optional) locales that should be used for fallback
  - en-gb
  translationFile: translations/en_GB.all.json     # the label file location
  translationFiles:                                # or a list of label file locations
  - translations/merged/en-gb.all.yaml
  - translations/translated/en-gb.adjusted.yaml
  accounting:
    default:                                    # configure display of prices
        thousand: ','
        decimal: '.'
        formatZero: '%s -.-'
        format: "%v %s"
    GBP:                                        # configure display of prices in currency GBP
        thousand: ','
        decimal: '.'
        formatZero: '%s -.-'
        format: "%v %s"
  numbers:                                         # configure display of numbers
    thousand: ','
    decimal: '.'
    precision: 2
  date:
    dateFormat:  02 Jan 2006
    timeFormat: 15:04:05
    dateTimeFormat: 02 Jan 2006 15:04:05
    location: LOCATIONCODE                          # required for formatLocaleTime
```

By providing different configurations for the different configuration areas (see prefixrouter module) you can easily build multilanguage applications.

## Usage in Templates:

### Localisation of Labels:

```pug
  # Display the label for the given key:
  __("key")
  
  # Display the label and pass a default that is used if the key is not existent in a labelfile
	__("key").setDefaultLabel("default")
	
	# Display the label with dynamic values replaces:
	__("key").setDefaultLabel("Hello Mr {{.UserName}}").setTranslationArguments({UserName: "Max"})
	
	# Support the usage of plural labels idepending on a given count:
	__("unread_mails").setCount(5)
	
	
	# Force the usage of another local by passing languange code as 5th paramater:
	__("switch_to_german").setLocale("de-DE")
	
```

#### Proposed organisation of label files:

We propose to put the translations in a folder *translations* like this:

```
translations
└───src (put the original label files here. This is where developers should work.)
└───merged (contains the generated files)
└───translated (optional - can contain the files returned from a translation tool or agency)
```

The label files in `translations/src` can either be json or yaml.
Example:
 
```json
[
  {
    "id": "attribute.clothingSize",
    "translation": "Size"
  },
  {
    "id": "error404.headline",
    "translation": "Page not found!"
  }
]
```

You can then run this command to generate the merged label file for the contained language codes:

```bash
goi18n merge -sourceLanguage en-us -format yaml -outdir translations/merged/ translations/src/*.json
```

Read more about label files and translation workflows here: [github.com/nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n)

### Formatting of dates:

Two template functions are provided:

 * `dateTimeFormatFromIso` - can get an ISO date format and returns the formatter object
 * `dateTimeFormat` - need to get a go `time.Time` object as input and returns the formatter

The formatter can format a date in the configured format - either in the timezone passed - or converted to the local timezone. 

E.g.:
```pug
dateTimeFormatFromIso("2006-01-02T15:04:05Z").formatDate()
dateTimeFormat(timeObject).formatToLocalDate()
```
Other functions are `formalToLocalDate()` or `formatTime()` etc..

**Note:** For displaying locale formats set the correct date `locationcode` - see [golang.org/pkg/time/#LoadLocation](https://golang.org/pkg/time/#LoadLocation)

### Formatting of prices:

```pug
priceFormat(90.25,"£")
// £ 90.25
priceFormatLong(42049.99,"$","USD")
// $ 42,049.99 USD
```

### Formatting of numbers:

Formatting of numbers can be configured like described above. The delimiter for thousand and
decimal can be configured. The precision for the decimal places can be configured with a default
value, but can also be overwritten.

```pug
// with defaul precision
numberFormat(12300)
// with overwritten precision
numberFormat(12300, 2)
```

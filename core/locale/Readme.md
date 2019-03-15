# Locale module (Localization)

This package provides localization features:

 * Translations of Labels: with the template func `__()`  (which uses https://github.com/nicksnyder/go-i18n)
 * Local display of Dates (In local timezone or given timezone): with the template func `dateTimeFormat` or `dateTimeFormatFromIso`
 * Local display of prices and numbers: with the template func `priceFormat` and `numberFormat` (Using https://github.com/leekchan/accounting)

## Configuration

```
locale:
  locale: en-gb                                    # the locale used for labels
  translationFile: translations/en_GB.all.json     # the label file location
  translationFiles:                                # or a list of label file locations
  - translations/merged/en-gb.all.yaml
  - translations/translated/en-gb.adjusted.yaml
  accounting:                                      # configure display of prices
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

By providing different configurations for the different configuration areas (see prefixrouter) you can easily build multilanguage applications.

## Usage in Templates:

### Localisation of Labels:

```pug
  # Display the label for the given key:
  __("key")
  
  # Display the label and pass a default that is used if the key is not existend in a labelfile
	__("key").SetDefaultLabel("default")
	
	# Display the label with dynamic values replaces:
	__("key").SetDefaultLabel("Hello Mr {{.UserName}}").SetTranslationArguments({UserName: "Max"})
	
	# Support the usage of plural labels idepending on a given count:
	__("unread_mails").SetCount(5)
	
	
	# Force the usage of another local by passing languange code as 5th paramater:
	__("switch_to_german").SetLocale("de-DE")
	
```

#### Proposed organisation of label files:

We propose to put the translations in a folder *translations* like this:

* translations
    * src (put the original label files here. This is where developers should work.)
    * merged (contains the generated files)
    * translated (optional - can contain the files returned from a translation tool or agency)

The label files in *translations/src* can either be json or yaml.
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
```sh
goi18n merge -sourceLanguage en-us -format yaml -outdir translations/merged/ translations/src/*.json
```


Read more about label files and translation workflows here:  https://github.com/nicksnyder/go-i18n

### Formatting of dates:

Two tenplatefunctions are provided:

 * dateTimeFormatFromIso - can get an ISO date format and returns the formatter object
 * dateTimeFormat - need to get a go time.Time object as input and returns the formatter

The formatter can format a date in the configured format - either in the timezone passed - or converted to the local timezone. 

E.g.:
```
dateTimeFormatFromIso("2006-01-02T15:04:05Z").formatDate()
dateTimeFormat(timeObject).formatToLocalDate()
```
Other functions are formalToLocalDate() or formatTime() etc..

Note: For displaying locale formats set the correct date locationcode - see https://golang.org/pkg/time/#LoadLocation

### Formatting of prices:

```
priceFormat(90,"GBP")
priceFormatLong(90,"GBP","british pound")
```

### Formatting of numbers:

Formatting of numbers can be configured like described above. The delimiter for thousand and
decimal can be configured. The precision for the decimal places can be configured with a default
value, but can also be overwritten.

```
// with defaul precision
numberFormat(12300)
// with overwritten precision
numberFormat(12300, 2)
```

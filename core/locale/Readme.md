## Locale Package

This package provides localization features:
 * the template func `__()`  which uses https://github.com/nicksnyder/go-i18n
 * templatefuncs to format DateTime ISO Strings 
 * templatefunc to format prices

### Configuration

```
locale:
  locale: en_GB
  translationFile: translations/en_GB.all.json
  translationFiles:
  - translations/en_GB.1.json
  - translations/en_GB.2.json
  accounting:
    thousand: ','
    decimal: '.'
    formatZero: '%s -.-'
    format: "%v %s"
  numbers:
    thousand: ','
    decimal: '.'
    precision: 2
  date:
    dateFormat:  02 Jan 2006
    timeFormat: 15:04:05
    dateTimeFormat: 02 Jan 2006 15:04:05
    location: LOCATIONCODE (for formatLocaleTime)
```

Planned for later:
 * configure and load multiple translationFile

### Usage in Templates:

#### Localisation of Labels:

```
  __("key")
	__("key","default")
	
	
	__("key","Hello Mr {{.UserName}}",{UserName: "Max"})
	
	//Use mehrzahl:
	__("unread_mails","",{_count: 5})
	
	
	// Force other than configured locale: 
	__("switch_to_german","",{},"de-DE")
	
```
#### Formatting of dates:

Two tenplatefunctions are provided:
 * dateTimeFormatFromIso - can get an ISO date format and returns the formatter object
 * dateTimeFormat - need to get a go time.Time object as input and returns the formatter

The formatter can format a date in the configured format - either in the format passed - or converted to the local timezone 

E.g.:
```
dateTimeFormatFromIso(flight.scheduledDateTimeStringInIsoFormat).formatDate()
dateTimeFormat(flight.scheduledDateTime).formalLocalDate()
```
Other functions are formalLocalDate() or formatTime() etc..

#### Formatting of prices:

```
priceFormat(90,"GBP")
```

#### Formatting of numbers:

Formatting of numbers can be configured like described above. The delimiter for thousand and
decimal can be configured. The precision for the decimal places can be configured with a default
value, but can also be overwritten.

```
// with defaul precision
numberFormat(12300)
// with overwritten precision
numberFormat(12300, 2)
```

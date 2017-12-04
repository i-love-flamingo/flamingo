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
  accounting:
    thousand: ','
    decimal: '.'
    formatZero: '%s -.-'
    format: "%v %s"
  dateTime:
    dateFormat:  02 Jan 2006
    timeFormat: 15:04:05
    dateTimeFormat: 02 Jan 2006 15:04:05
```

Planned for later:
 * configure and load multiple translationFile

### Usage in Templates:

```
  __("key")
	__("key","default")
	
	
	__("key","Hello Mr {{.UserName}}",{UserName: "Max"})
	
	//Use mehrzahl:
	__("unread_mails","",{_count: 5})
	
	
	// Force other than configured locale: 
	__("switch_to_german","",{},"de-DE")
	
```

## Locale Package

This package provides the template func `__()`  

Uses: https://github.com/nicksnyder/go-i18n

### Configuration

```
locale:
  locale: en_GB
  translationFile: translations/en_GB.all.json
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

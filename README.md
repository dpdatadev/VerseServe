# Go Bible

Gobible is a library for interacting with the Bible in Go.

## Bible Formats

Initially, Gobible was created as an effort to support an easy to use JSON format for the Bible for use in other projects,
but through development has grown to instead be a library for working with that format, as well and importing other formats into the Gobible format.
### Supported Formats
- Gobible JSON [Example](https://raw.githubusercontent.com/gobible/gobible/master/data/KJV.json)
- OSIS XML [Example](https://raw.githubusercontent.com/gobible/gobible/master/data/WEB.xml)


Please note that since Gobible works around a data structure defined by the Gobible JSON format, some features of other formats may not be supported.

## Usage

### Initialization

There are a few ways to initialize a bible:

```go
// Initialize an empty GoBible object
b := NewGoBible()

// Load a GoBible formatted JSON file
b.Load("data/KJV.json")

// Load an alternate support format
b.LoadFormat("data/WEB.xml", "osis")
```

One you have a Bible object with some translations loaded into it, you can use this object in a few ways.

### Interacting with Individual Translations

You can load an individual translation, and use it similarly to this:

```go
// Load the KJV translation
kjv := bible.GetTranslation("KJV")

// get a verse
verse := kjv.GetBook("Genesis").GetChapter(1).GetVerse(1)

fmt.Println("Genesis 1:1 - " + verse.Text)
// Genesis 1:1 - In the beginning God created the heaven and the earth.
```

You can also use this to generate GoBible format JSON files from other formats:

```go
b := NewGoBible()
// Load an OSIS XML file
b.LoadFormat("WEB.xml", "osis")

// Save the translation as GoBible JSON
goBibleJson := b.GetBibleJSON("WEB")

// Save the translation as GoBible JSON to a file
os.WriteFile("WEB.json", []byte(goBibleJson), 0644)
```

### Interacting by Reference

Once you have 1 or multiple translations loaded, you can interact with them by reference notation.

(If you have multiple translations loaded, the first translation loaded will be used as the default)


```go

// Examples of some valid references:
//   John 3:16
//   John 3:16-18
//   1 Cor 1:2-5
//   1 Cor 1:2-2:5 
// You can also tag the translation at the end of the reference
//   John 3:16 WEB 
//   John 3:16 KJV

// Lookup by reference
refs, _ := b.ParseReference("Gen 1:31-2:1")
for _, r := range refs {
    fmt.Printf("%s %d:%d %s", r.Book, r.Chapter, r.Verse, r.VerseRef.Text)
}
// Genesis 1:31 God saw everything that he had made, and, behold, it was very good. There was evening and there was morning, the sixth day.
// Genesis 2:1 The heavens and the earth were finished, and all the host of them.
```

### Create HTTP Handlers to Embed Bible functionality into your web services or apps
```go
//server.go
func VerseServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Initialize an empty GoBible object
	b := NewGoBible()

	// Load a GoBible formatted JSON file
	b.Load("data/KJV.json")

	ref := GetReference(r)
	if ref == "" {
		http.Error(w, "missing reference argument", http.StatusBadRequest)
		return
	}

	log.Println("parsing BIBLE REF..", ref)
	verses, err := b.ParseReference(ref)
	if err != nil {
		log.Println("failed to parse reference:", err)
		http.Error(w, "invalid reference format", http.StatusBadRequest)
		return
	}
	verseList := make([]ReferenceDTO, len(verses))
	for i, v := range verses {
		verseList[i] = ReferenceDTO{
			Book:    v.Book,
			Chapter: v.Chapter,
			Verse:   v.Verse,
			Text:    v.VerseRef.Text,
		}
	}

	WriteHttpJson(verseList, w)
}
```
### Now Query:
```sh
curl http://127.0.0.1:7777/verse?ref=Genesis%201:1-3
```
### Or using a Ruby client:
```ruby
 reference_response = HTTParty.get("http://127.0.1:7777/verse?ref=#{readings.join(',')}")
 if reference_response.success?
	Scrapers::ServiceUtils.debug_log("Successfully retrieved verse text from VerseServe API for readings: #{readings.join(', ')}")
	verse_texts = reference_response.parsed_response['Verses']
	return_readings = readings.zip(verse_texts).map do |ref, text|
		Bible::Reading.new(ref, text)
 	end
 else
    Scrapers::ServiceUtils.debug_log("Failed to retrieve verse text from VerseServe API for readings: #{readings.join(', ')}. HTTP Status: #{reference_response.code}. Falling back to manual parsing of readings without verse text.")
    return_readings = manually_strip_readings(readings)
 end
```

## Looking for Bibles?

We have published a few Public Domain bibles, as well as a Go program to generate additional ones in the [gobible-gen](https://github.com/solafide-dev/gobible-gen) repo.

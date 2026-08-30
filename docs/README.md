# Document Fingerprint & Similarity

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/fingerprint)
![Build](https://github.com/lordofscripts/fingerprint/actions/workflows/build.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/fingerprint.svg)](https://pkg.go.dev/github.com/lordofscripts/fingerprint)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/fingerprint)](https://github.com/lordofscripts/fingerprint/releases/latest)
[![Created](https://badges.pufler.dev/created/lordofscripts/fingerprint)](https://badges.pufler.dev)
[![Updated](https://badges.pufler.dev/updated/lordofscripts/fingerprint)](https://badges.pufler.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fingerprinting and similarity module for Go.

This is an implementation of the algorithm described in
the [Winnowing: Local Algorithms for Document Fingerprinting](http://igm.univ-mlv.fr/~mac/ENS/DOC/sigmod03-1.pdf) paper.
It can be used to fingerprint document text and calculate a similarity score using
the [Jaccard index (similarity coefficient)](https://en.wikipedia.org/wiki/Jaccard_index).

### The Source & Credits

This Fingerprint module is a **fork** of the [original](https://github.com/grahambrooks/fingerprint) 
by Graham Brooks. He wrote it in 2021 but it hasn't had any functional update since. After that, his
repository receives only CI updates. 

Based on certain needs and curiousitites,  I recently made a few suggestions (2026) regarding its 
documentation, examples, wrong GO version reporting and lack of file support. Since my 
suggestions appeared to be ignored, I decided to pitch in. 

> This is not just another fork without updates, it contains useful enhancements and improved documentation.

### A Repository of Interest

The original implementation got under my radar because several years ago I was taking an online
Creative Writing Specialization (not software writing). Since there we had to review the works of
others, soon it became obvious that some of the *wannabe-authors* enrolled in the course to **plagiarize**.
Yes! several times I had to review assignments that were a blatant copy of someone elses! That
I know to this date the course platform does not check for that. Therefore, I got interested
in the possibility of detecting such things. But I needed it to find the similarity in files,
not just simple strings.

## (New) Key Features

- **Winnowing Algorithm**: Robust local fingerprinting that is insensitive to small changes in the document.
- **Jaccard Similarity**: Easy comparison of fingerprints to determine document similarity. *Jaccard's* score
  measures a percentage of shared elements. It is best for finding structural overlap.
- **Text Cleaning**: Built-in text normalization (removes non-letters and converts to lowercase) for more accurate
  matching.
- **Size-adjusted, Multi-set and Coverage** scores are now added (not present in the original) 
- **File comparison** which was neither supported nor included in the original module.
- **Improved documentation** to help end-user decide HOW to use it and which metric to use.

<table align="center">
<tr>
<th align="center" colspan="2">
Show your support for continued development of these useful software applications
</th>
</tr>
<tr>
<td>
<img src="./assets/allmylinks.png?raw=true" alt="AllMyLinks logo" />
</td>
<td>
Visit <br> Lord of Scripts&trade; on<br><a href="https://allmylinks.com/lordofscripts">AllMyLinks.com</a>
</td>
</tr>
<tr>
<td>
<img src="./assets/buymecoffee-dark.png?raw=true" alt="Buy LordOfScripts Coffee" />
</td>
<td>
Buy Lord of Scripts&trade; <br> a Capuccino on <br><a href="https://www.buymeacoffee.com/lostinwriting">BuyMeACoffee.com</a>
</td>
</tr>
</table>

![Divider](./assets/lordofscripts-divider-golden.png)

## Installation

For development purposes:

```bash
go get github.com/lordofscripts/fingerprint
```

For using the demo application to compare documents:

```bash
go install github.com/lordofscripts/fingerprint
```

### Usage

The library provides a high-level `similarity` package for quick comparisons and a `fingerprinter` package for more
granular control. It now also provides file-based functionality.

#### As an end-user

Here is how to use this module's application command as an end-user.

General Usage:

> `similarity [OPTIONS] [ARGUMENTS]`

Fingerprinting a single item gives you a hash of the item. It is not
a hash of its contents but a hash taken over the structure, repetitions,
and other attributes that compose the document.

To *fingerprint* text from the command-line:

> `similarity [--no-normalize|--letters-only] -str1 "Text string to be analyzed"

To *fingerprint* a TEXT file:

> `similarity [--no-normalize|--letters-only] --use-files FILENAME

But, if on the other hand you need to *compare the similarity* of two files (or two strings) and
take a decision based on the presented results, you have two manners.

> `similarity [--no-normalize|--letters-only] -str1 "Text string to be analyzed" -str2 "Other text to compare to"

To evaluate the *similarity between two files*:

> `similarity [--no-normalize|--letters-only] --use-files FILENAME_1  FILENAME_2

And you can also do:

- `similarity -help` to obtain help about how to invoke the application
- `similarity -version` to show the application/module version.


#### As a Developer

Here is how to integrate and use this module's functionality into your
applications.


##### Simple Similarity String Comparison

```go
package main

import (
  "fmt"

  "github.com/lordofscripts/fingerprint/fingerprinter"
  "github.com/lordofscripts/fingerprint/similarity"
)

func main() {
  text1 := "The quick brown fox jumped over the lazy dog"
  text2 := "The quick brown fox jumps over the lazy dog"

  // Options define the sensitivity of the fingerprinting
  options := fingerprinter.Options{
    GuaranteeThreshold: 4,
    NoiseThreshold:     4,
  }

  score := similarity.StringSimilarity(text1, text2, options)

  fmt.Printf("Similarity score: %f\n", score)
}
```

##### Simple Similarity File Comparison

For a better real-life example see the `cmd/main.go` which also
showcases the extra functionality such as `FileSimilarity` that I added to 
Graham Brooks' original version (for which I am grateful):

```go
package main

import (
	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/lordofscripts/fingerprint/similarity"
)

func main() {
  text1 := "The quick brown fox jumped over the lazy dog"
  text2 := "The quick brown fox jumps over the lazy dog"

  options := fingerprinter.Options{
		GuaranteeThreshold: 4,
		NoiseThreshold:     4,
	}

  if !useFiles {
		score := similarity.StringSimilarity(text1, text2, options)
		fmt.Println("Similarity: ", score)
	} else if flag.NArg() == 2 {
    _, scores, err = similarity.FileSimilarity(
			flag.Arg(0),
			flag.Arg(1),
			options)
		if err != nil {
			DieWithError(EXIT_CODE_SIMILARITY, err)
		}
    fmt.Println("File-1: ", flag.Arg(0))
		fmt.Println("File-2: ", flag.Arg(1))
		fmt.Println("\t*** Similarity scores ***")
		fmt.Println(scores)
  }
}

```

And the output might look like this:

```bash
lordofscripts@munich$ bin/similarity -use-files tests/assets/document_A_20.txt  tests/assets/document_B_20.txt 
File-1:  tests/assets/document_A_20.txt
File-2:  tests/assets/document_B_20.txt
	*** Similarity scores ***
	Raw          : 95.9%
	Size-adjusted: 95.5%
	Multi-set    : 96.0%
	Coverage     : 98.1%
	Fingerprint-1: a6c2fa9c3f58672d
	Fingerprint-2: 49d3e0cd0236562d
```


### Interpreting Results

Here is how to use or interpret the results given by the *Similarity Scores*:

- Use `Coverage` when you want to detect whether one document is contained in another. The
  coverage function is asymmetric.
- Use `Size-adjusted` when you want an overall similarity score.
- Use `MultiSet` if you want repetitions of KGrams in the document to count
- Add Document Length for *overall similarity*
- The `Fingerprint-*` result is important because it is *not* a hash/digest of the 
  data. Such thing changes completely even if a single letter changes. The outputted
  `fingerprint` value is taken over the K-Grams into which the input file was decomposed,
  therefore, it takes into consideration similarity, repetitions, overlaps, etc. It is
  more meaningful than a mere hash that tells you if two things are **identical**.

### Understanding Options

- `NoiseThreshold` (k): Any match shorter than `k` is considered noise and ignored.
- `GuaranteeThreshold` (t): Any match of length at least `t` is guaranteed to be detected.
- `Normalize` causes all text characters to be normalized to *lowercase*

A common configuration is to set `NoiseThreshold` and `GuaranteeThreshold` to the same value (e.g., 4 or 5).

![Divider](./assets/lordofscripts-divider-golden.png)

## Performance Note

The library is designed for relatively small documents and is not currently optimized for use with large streams or
channels.

Do you want to know more? Read the [In Depth Information](IN_DEPTH.md)


### What Changed from the Original

The upstream baseline from which this was taken was: `grahambrooks/fingerprint` tag `2026-08-27-517b8d3`.

- I am using standard software versioning for tags rather than dates.
- The `go.mod` properly reports the minimum required GO version (not the last version)
- Added a `Taskfile.yml` for building with [TaskFile](https://taskfile.dev/)
- Added the `cmd` package with therein an interactive demo application
- Added multiple score alternatives consolidated into `similarity.SimilarityScores`
- Added human-friendly Document fingerprints to `similarity.SimilarityScores`
- Added `tests/assets` directory with test data therein. Tests moved there too.
- Added `similarity.FileSimilarity()` function
- Added the `fingerprinter.StreamFingerprinter` object plus extra test
- Updated `fingerprinter.TextFingerprint()`
- Added `fingerprinter.FileFingerprint()`
- Added the `Normalize` option to `fingerprinter.Options`

#### History

- 2026-08-28 Added File fingerprinting and 3 supplementary *Similarity Scores* (Go v1.22)
- 2026-08-28 `origin` based on Graham's `2026-08-27-517b8d3` update. (Go v1.22)
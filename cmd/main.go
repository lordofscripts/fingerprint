/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           GoFingerprint
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Document Fingerprint and Various document similarity scores
 * between two text files.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lordofscripts/fingerprint"
	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/lordofscripts/fingerprint/similarity"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func HelpMe() {
	fmt.Println("Similarity Usage:")

	fmt.Println("1. Compare similarity of two strings:")
	fmt.Println("\tsimilarity -str1 'sample text' -str2 'another sample text'")

	fmt.Println("2. Compare similarity of two files:")
	fmt.Println("\tsimilarity --use-files FILENAME_1 FILENAME_2")

	flag.PrintDefaults()
}

func Die(exitCode int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(exitCode)
}

func DieWithError(exitCode int, err error) {
	fmt.Fprintf(os.Stderr, "ERROR:\n\t%v\n", err)
	os.Exit(exitCode)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

func main() {
	const EXIT_CODE_OK int = 0
	const EXIT_CODE_ARGUMENTS int = 1
	const EXIT_CODE_SIMILARITY int = 2
	const EXIT_CODE_FINGERPRINT int = 3
	var text1, text2 string
	var useFiles, help, fingerprintOnly, optUnNormalize, optLettersOnly, version bool
	flag.StringVar(&text1, "str1", "", "Plain string 1")
	flag.StringVar(&text2, "str2", "", "Plain string 2")
	flag.BoolVar(&useFiles, "use-files", false, "Use the two filenames given as arguments")
	flag.BoolVar(&optUnNormalize, "no-normalize", false, "Do NOT Normalize text to lowercase")
	flag.BoolVar(&optLettersOnly, "only-letters", false, "Strip non-letters from analysis")
	flag.BoolVar(&fingerprintOnly, "fingerprint", false, "Only do fingerprinting, no similarity (only ONE item)")
	flag.BoolVar(&help, "help", false, "Help me")
	flag.BoolVar(&version, "version", false, "Version")
	flag.Usage = HelpMe
	flag.Parse()

	if help {
		HelpMe()
		os.Exit(EXIT_CODE_OK)
	}
	if version {
		fmt.Printf("Similarity %s\n", fingerprint.ModuleVersion.Short())
		os.Exit(EXIT_CODE_OK)
	}

	options := fingerprinter.Options{
		GuaranteeThreshold: 4,
		NoiseThreshold:     4,
		Normalize:          !optUnNormalize,
		LettersOnly:        optLettersOnly,
	}

	// I: Fingerpint rather than Similarity
	if fingerprintOnly {
		if useFiles {
			// similarity -fingerprint --use-files FILENAME
			if flag.NArg() != 1 {
				HelpMe()
				Die(EXIT_CODE_ARGUMENTS, "-fingerprint -use-files requires ONE filename as argument! have: %d", flag.NArg())
			}

			if fp, err := fingerprinter.FileFingerprint(flag.Arg(0), options); err != nil {
				DieWithError(EXIT_CODE_FINGERPRINT, err)
			} else {
				fmt.Println("File: ", flag.Arg(0))
				fmt.Printf("Fingerprint: %16s\n", fp.HashStr())
			}

			os.Exit(EXIT_CODE_OK)
		} else if len(text1) == 0 {
			HelpMe()
			Die(EXIT_CODE_ARGUMENTS, "-fingerprint -str1 STRING requires ONE string!")
		} else {
			// similarity -fingerprint -str1 "TEXT TO BE ANALYZED"
			fp := fingerprinter.TextFingerprint(text1, options)
			fmt.Println("String: ", text1)
			fmt.Printf("Fingerprint: %16s\n", fp.HashStr())
		}

		os.Exit(EXIT_CODE_OK)
	}

	// II: Similarity between two strings or files, includes fingerprinting
	var scores *similarity.SimilarityScores
	if !useFiles {
		_, scores = similarity.StringSimilarity(text1, text2, options)
		fmt.Println("String-1: ", text1)
		fmt.Println("String-2: ", text2)
		fmt.Println("\t*** Similarity scores ***")
		fmt.Println(scores)
	} else if flag.NArg() == 2 {
		var err error

		_, scores, err = similarity.FileSimilarity(
			flag.Arg(0),
			flag.Arg(1),
			options,
		)
		if err != nil {
			DieWithError(EXIT_CODE_SIMILARITY, err)
		}

		fmt.Println("File-1: ", flag.Arg(0))
		fmt.Println("File-2: ", flag.Arg(1))
		fmt.Println("\t*** Similarity scores ***")
		fmt.Println(scores)
	} else {
		HelpMe()
		Die(EXIT_CODE_ARGUMENTS, "--use-files requires two filenames as arguments! have: %d", flag.NArg())
	}
}

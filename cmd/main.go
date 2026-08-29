package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lordofscripts/fingerprint/fingerprinter"
	"github.com/lordofscripts/fingerprint/similarity"
)

const VERSION string = "1.0.0"

func HelpMe() {
	fmt.Println("Similarity Usage:")

	fmt.Println("1. Compare similarity of two strings:")
	fmt.Println("\tsimilarity -str1 'sample text' -str2 'another sample text'")

	fmt.Println("2. Compare similarity of two strings:")
	fmt.Println("\tsimilarity -str1 'sample text' -str2 'another sample text'")
}

func Die(exitCode int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(exitCode)
}

func DieWithError(exitCode int, err error) {
	fmt.Fprintf(os.Stderr, "ERROR:\n\t%v\n", err)
	os.Exit(exitCode)
}

func main() {
	const EXIT_CODE_OK int = 0
	const EXIT_CODE_ARGUMENTS int = 1
	const EXIT_CODE_SIMILARITY int = 2
	var text1, text2 string
	var useFiles, help, version bool
	flag.StringVar(&text1, "str1", "", "Plain string 1")
	flag.StringVar(&text2, "str2", "", "Plain string 2")
	flag.BoolVar(&useFiles, "use-files", false, "Use the two filenames given as arguments")
	flag.BoolVar(&help, "help", false, "Help me")
	flag.BoolVar(&version, "version", false, "Version")
	flag.Usage = HelpMe
	flag.Parse()

	if help {
		HelpMe()
		os.Exit(EXIT_CODE_OK)
	}
	if version {
		fmt.Printf("Similarity v%s\n", VERSION)
		os.Exit(EXIT_CODE_OK)
	}

	options := fingerprinter.Options{
		GuaranteeThreshold: 4,
		NoiseThreshold:     4,
	}

	if !useFiles {
		score := similarity.StringSimilarity(text1, text2, options)
		fmt.Println("Similarity: ", score)
	} else if flag.NArg() == 2 {
		var err error
		var scores *similarity.SimilarityScores
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
		Die(EXIT_CODE_ARGUMENTS, "--use-files requires to filenames as arguments! have: %d", flag.NArg())
	}
}

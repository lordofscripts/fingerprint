
![Divider](./assets/lordofscripts-divider-golden.png)
## More In-depth Information

The test directory contains these:

- `document_A_20.txt` 20 identical paragraphs of "Lorem ipsum..."
- `document_A_15.txt` 15 identical paragraphs of "Lorem ipsum..."
- `document_B_20.txt` 20 identical paragraphs of "Lorem epsimia..." (just "impsum" replaced by "epsimia")
- `document_C_20.txt` 20 paragraphs of "Lorem ipsum..." where 3 words have been replaced.
- `random_A_20.txt` 20 random paragraphs.
- `random_B_20.txt` 20 random paragraphs in a different order than the previous.

The so-called **raw** similarity score is the plain score as calculated by the original
[fingerprint module](https://github.com/grahambrooks/fingerprint/) by Graham Brooks.
This simplified table of various manual runs, the **raw** similarity scores will be:

| FILENAME 1          | FILENAME 2          | SCORE |
| ------------------- | ------------------- | ----- |
| `document_A_20.txt` | `document_A_20.txt` | 1.000 | Same files with repeated paragraphs
| `document_A_20.txt` | `document_A_15.txt` | 1.000 | Same content but 5 paragraphs removed
| `document_A_20.txt` | `document_B_20.txt` | 0.959 | One word replaced in the 2nd document
| `document_A_20.txt` | `document_C_20.txt` | 0.878 | Three words replaced in the 2nd document
| `document_A_20.txt` | `random_A_20.txt`   | 0.016 | "Lorem ipsum..." vs random paragraphs.
| `document_C_20.txt` | `random_A_20.txt`   | 0.016 | Variation of the previous
| `random_A_20.txt`   | `random_B_20.txt`   | 1.000 | Random in different order

- "Lorem ipsum" paragraphs have very low information density and are highly repetitive.
  Once converted to fingerprints, they have essentially the same fingerprint vocabulary.
- A `1.000` similarity score basically means that *every fingerprint is found on both documents*.
  It does **not** necessarily mean they have the same length, number of repetitions or the
  same document structure.

So, what's the usefulness of this you may ask?

- Detect whether one docuent contains the same distinctive passages as another.
- Finding copied or lightly modified text.
- Searching a large corpus of documents sharing substantial content.
- Detecting near-duplicate documents (like web pages) with changed formatting or boilerplate.

Do you want to compare equality? Then simply use a Hash/Digest like MD5, SHA256, etc.

Explanation why #2 & #7 report a score of `1.000` is because the similarity score
calculated by this module is treating fingerprints as a Set, not as a Sequence or
Multiset. The repeated paragraphs generate the same Winnowing marks. Once the duplicate
marks are discarded, both files could look like:

> `file A: {mark1, mark2, mark3, mark4}`
> `file B: {mark1, mark2, mark3, mark4}`

The **Jaccard Similarity** is defined as:

> `J(A,B)= ∣A∪B∣ / ∣A∩B∣`

And since both sets reduced to the same. The above becomes `4/4 = 1.000`. The algorithm
detects *which content is present*, but not *how many times it occurs*.​

Given that the **raw** score could appear meaningless of ambiguous, I added a few
more scoring features to aid you in better fingerprinting. These are the scores
returned in the new `similarity.SimilarityScores` structure:


| FILENAME 1          | FILENAME 2          |  RAW  | SIZED | MULTI | COVERAGE |
| ------------------- | ------------------- | ----- | ----- | ----- | -------- |
| `document_A_20.txt` | `document_A_20.txt` | 1.000 | 1.000 | 1.000 | 1.000 |
| `document_A_20.txt` | `document_A_15.txt` | 1.000 | 0.750 | 0.750 | 1.000 | 
| `document_A_20.txt` | `document_B_20.txt` | 0.959 | 0.955 | 0.960 | 0.981 |
| `document_A_20.txt` | `document_C_20.txt` | 0.878 | 0.874 | 0.880 | 0.937 |
| `document_A_20.txt` | `random_A_20.txt`   | 0.016 | 0.013 | 0.009 | 0.018 |
| `document_C_20.txt` | `random_A_20.txt`   | 0.016 | 0.013 | 0.009 | 0.018 |
| `random_A_20.txt`   | `random_B_20.txt`   | 1.000 | 1.000 | 1.000 | 1.000 |

### What Changed from the Original

The upstream baseline from which this was taken was: `2026-08-27-517b8d3`.

- I am using standard software versioning for tags rather than dates.
- The `go.mod` properly reports the minimum required GO version (not the last version)
- Added a `Taskfile.yml` for building with [TaskFile](https://taskfile.dev/)
- Added the `cmd` package with therein an interactive demo application
- Added multiple score alternatives consolidated into `similarity.SimilarityScores`
- Added `tests/assets` directory with test data therein. Tests moved there too.
- Added `similarity.FileSimilarity()` function
- Added the `fingerprinter.StreamFingerprinter` object plus extra test
- Updated `fingerprinter.TextFingerprint()`
- Added `fingerprinter.FileFingerprint()`
- Added the `Normalize` option to `fingerprinter.Options`

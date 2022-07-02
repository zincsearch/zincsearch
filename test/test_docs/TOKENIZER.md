## Test Tokenizer
 tokenizer will take a stream of continuous text and break it up into tokens. Below is the functionality of the tokenizer that are tested using TestAnalyer with endpoint `/api/_analyze`:
   |No. | Tokenizer Name                 | Tokenizer Functionality |
  | -- | ------------------------------- | ---------------- | 
  | 1 | Standard                         | Divide text into terms and remove punctuation symbols|
  | 2 | Letter                           | Divide text into terms whenever encouters a character which is not a letter|
  | 3 | Lowercase                        | Divide text into terms whenver encouters a character which is not a letter and also lowercases the letter| 
   | 4 | Whitespace                      | Divide text into terms whenever encouters a whitespace| 
   | 5 | N-Gram                          | Break text into words and return a n-gram of each word e-g quick --> [qu,ui,ic,ck]| 
   | 6 | Edge N-Gram                     | Break text into words and return n-gram of each word which are anchored to the start of word e-g quick --> [q,qu,qui,quic,quick] | 
   | 7 | Keword                          | Accept the text and output the exact same text as a single term | 
   | 8 | Regexp                          | Split text into terms whenever it encouters a non-word character|
   | 9 | Character Group                 | Split the text into terms through sets of character specified , less expensive than regexp|
   | 10 | Path Hierarchy                 | Takes hierarchial value like filesystem and splits on path separator and emits a term for each character in tree. e-g /one/two/three --> [/one,/one/two,/one/two/three] |
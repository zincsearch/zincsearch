## Test Token Filter
Token filter accepts a stream of tokens from the tokenizer and can modify tokens. Below is the functionality of the token filter that are tested using TestAnalyzer with endpoint `/api/_analyze`:
  |No. | Token Filter Name                    | Token Filter Functionality |
  | -- | -------------------------------------| ---------------- | 
  | 1 | Apostrophe                            |Strip all charachters after the apostrophe including the apostrophe as well|
  | 2 | CJK Bigram                            | Form Bigrams out of Chinese,Japanese,Korean (CJK) tokens|
  | 3 | CJK Width                             | Normalizes width differences in CJK characters| 
   | 4 | Dictionary                           | Use specified list of words and brute force approach to find subwords in compound words, if found included in the output token| 
   | 5 |Edge N-Gram                           | Forms an n-gram of a specified length from the beginning of a token| 
   | 6 | N-Gram                               | Forms n-grams of specified lengths from a token| 
   | 7 | Elision                              | Removes specified elisions from the beginning of tokens | 
   | 8 | Stemmer                              | Provides algorithmic stemming for several languages,algorithmic stemmer applies a series of rules to each word to reduce it to its root form|
   | 9 | Keyword                              | Marks specified tokens as keywords, which are not stemmed|
   | 10 | Length                              | Removes token shorter or longer than specified character lengths |
   | 11 | Lower case                            | Changes token text to lower case|
   | 12 | Replace                              | Replace the specified match string to the specified replace string |
   | 13 | Reverse                              | Reverse the order of the token's characters|
   | 14 | Shingle                       | Add shingles, or word n-grams, to a token stream by concatenating adjacent tokens|
   | 15 | Stop                          | Removes stop character from the tokens|
   | 16 | Trim                        | Remove leading and trailing whitespaces from each token | 
   | 17 | Truncate                     | Truncate tokens that exceed a specified character limit| 
   | 18 | Unique                       | Remove duplicate token from a stream | 
   | 19 | Upper case                   | Changes token text to upper case| 
   | 20 | ASCII folding character      | Converts alphabetic, numeric, and symbolic characters that are not in the Basic Latin Unicode block (first 127 ASCII characters) to their ASCII equivalent, if one exists|
   | 21 | HTML strip                   |Strips HTML elements from a text and replaces HTML entities with their decoded value|
   | 22 | Mapping character           |The mapping character filter accepts a map of keys and values. Whenever it encounters a string of characters that is the same as a key, it replaces them with the value associated with that key|
   | 23 | Pattern Replace             | Uses a regular expression to match and replace token substrings| 
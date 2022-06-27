## Test Analyzers 
Analyzers are combination of tokenizers and token filters. Tokenizer split the text into terms and further filters are applied on the terms using specified token filters.
Below is the functionality of the Analyzer that are tested using TestAnalyzer with endpoint `/api/_analyze`:
  |No. | Analyzer Name                                  | Analyzer Functionality |
  | -- | ---------------------------------------------- | ---------------- | 
  | 1 | Standard                                       | Lower case alphabets|
  | 2 | Standard Analyzer with stop words              | Lower case alphabets and drop stop words|
  | 3 | Standard Analyzer with stop words and filter   | Lower case alphabets,drop stop words and apply specified filter| 
   | 4 | Simple                                        | Lower case alphabets and drop numbers| 
   | 5 | Keyword                                       | Accept the text and output the exact same text as a single term | 
   | 6 | Regexp                                        | Lower case alphabets and drop characters (i.e ",#@ etc)  | 
   | 7 | Regexp with pattern                    | Lower case alphabets and look for pattern to drop and tokenize accordingly| 
   | 8 | Stop                                          | Lower case alphabets and drop stop words|
   | 9 | Stop Analyzer with stop words                 | Lower case alphabets and drop specified words|
   | 10 | Whitespace                                   | Divide text to terms based on whitespace |
   | 11 | Web                                          | Lower case alphabets, drop stop words and identify URLs |

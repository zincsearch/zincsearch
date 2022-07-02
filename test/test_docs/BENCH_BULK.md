## Test Bulk Worker Function
1. This benchmark test the BulkWorker function that takes target as an `Index` name and f as a `file name`. 
2. The file should be in ndjson form. 
3. Using os.Open to read the content of ndjson file.
4. A for loop whill run b.N times and examine the performance of the Bulk Worker Function code.
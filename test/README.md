# Unit Test Workflows
## Setting up the test environment (for etcd metadata storage)
1. Need to install **etcd** , add `root` user and enable authentication. You can use this [link](https://etcd.io/docs/v3.4/op-guide/authentication/) to setup etcd. 
2. Update the root password in your .env file.
3. You are good to go.

## For boltdb, badger metadata storage 

You just need to specify ZINC_METADATA_STORAGE environment variable to "bolt" or "badger" . Default value of ZINC_METADATA_STORAGE is "bolt"

## Steps to run tests
1. A bash script at project root named `test.sh` can be used to run these unit tests.
```sh
cd zinc
./test.sh
``` 
2. If you pass an argument `bench` then it will run the benchmark test.
```sh
./test.sh bench
```

## API Tests
 
### Testing Authentication (Test File : auth_test.go)
---
1. #### Testing  Auth API
   1. [Testing Authentication with Authorization](./test_docs/AUTH_API.md)
   2. [Testing Authorization with error Password](.test_docs/AUTH_API.md)
   3. [Testing Authentication without Authorization](.test_docs/AUTH_API.md) 

2. #### Testing  User API
   1. [Testing Login with username and password](./test_docs/USER_AUTH_API.md)
   2. [Testing Login with bad username or password](./test_docs/USER_AUTH_API.md)
   3. [Testing Create User API](./test_docs/USER_AUTH_API.md)
   4. [Test Update User](./test_docs/USER_AUTH_API.md)
   5. [Test Create User with error Input](./test_docs/USER_AUTH_API.md)
   6. [Test Delete User with existing UserID](./test_docs/USER_AUTH_API.md)
   7. [Test Delete User with no existing UserID](./test_docs/USER_AUTH_API.md)
   8. [Test GET User](./test_docs/USER_AUTH_API.md) 
---
### Test Base API (Test File : base_test.go)
---
#### Base API Test
   1. [Test Base API](./test_docs/BASE_API.md) 
   2. [Test Version API](./test_docs/BASE_API.md)
   3. [Test Health API](./test_docs/BASE_API.md)
   4. [Test User Interface API](./test_docs/BASE_API.md)
---

### Test Document Bulk (Test File : document_bulk_test.go)
---
1. #### Test Bulk Document with Endpoint /api/_bulk
   1. [Test Bulk Document Index](./test_docs/BULK_API.md) 
   2. [Test Bulk Document Delete](./test_docs/BULK_API.md)
   3. [Test Bulk Document with error Input](./test_docs/BULK_API.md)

2. #### Test Bulk Document with Enpoint /api/:target/_bulk
   1. [Test Create Document with not existing Index](./test_docs/BULK_API_TARGET.md)
   2. [Test Create Document with existing Index](./test_docs/BULK_API_TARGET.md)
   3. [Test Bulk Document with error Input](./test_docs/BULK_API_TARGET.md)
---

### Test Document (Test File: document_test.go)
---
1. #### Test Create Document with Endpoint /api/:target/_doc
   1. [Test Create Document with not existing Index](./test_docs/DOC_CREATE.md)
   2. [Test Create Document with existing Index](./test_docs/DOC_CREATE.md)
   3. [Test Create Document with existing Index and not existing Document ID](./test_docs/DOC_CREATE.md)
   4. [Test Update Document with existing Index and existing Document ID](./test_docs/DOC_CREATE.md)
   5. [Test Create Document with error Input](./test_docs/DOC_CREATE.md)

2. #### Test Update Document with Endpoint /api/:target/_doc/:id
   1. [Test Update Document with no existing Index](./test_docs/DOC_UPDATE.md)
   2. [Test Update Document with existing Index](./test_docs/DOC_UPDATE.md)
   3. [Test Update Document with existing Index and not existing ID](./test_docs/DOC_UPDATE.md)
   4. [Test Update Document with existing Index and existing ID](./test_docs/DOC_UPDATE.md)
   5. [Test Update Document with error Input](./test_docs/DOC_UPDATE.md)

3. #### Test Delete Document with Endpoint /api/:target/_doc/:id
   1. [Test Delete Document with not existing Index](./test_docs/DOC_DELETE.md)
   2. [Test Delete Document with existing Index and not existing ID](./test_docs/DOC_DELETE.md)
   3. [Test Delete Document with existing Index and existing ID](./test_docs/DOC_DELETE.md)
---

### Test Index (Test File : index_test.go)
---
1. #### Test Index with Endpoint /api/index
   1. [Test Create Index with Payload](./test_docs/INDEX.md)
   2. [Test Create Index with error Input](./test_docs/INDEX.md)
   3. [Test GET Index](./test_docs/INDEX.md)
   4. [Test Delete with existing Index](./test_docs/INDEX.md)
   5. [Test Delete with not existing Index](./test_docs/INDEX.md)

2. #### Test Index Mapping with Endpoint /api/:target/_mapping
   1. [Test Update mapping for Index](./test_docs/INDEX_MAP.md)
   2. [Test GET mappings from Index](./test_docs/INDEX_MAP.md)

---

### Test ElasticSeacrch API (Test File : es_test.go)
---
1. #### Test ES API with Endpoint /es/_bulk
   1. [Test Bulk Document Index](./test_docs/ES_API_1.md)
   2. [Test Bulk Document Delete](./test_docs/ES_API_1.md)
   3. [Test Bulk Document with error Input](./test_docs/ES_API_1.md)

2. #### Test ES API with Endpoint /es/:target/_bulk
   1. [Test Create Document with no existing Index](./test_docs/ES_API_2.md)
   2. [Test Create Document with existing Index](./test_docs/ES_API_2.md)
   3. [Test Bulk with error Input](./test_docs/ES_API_2.md)

3. #### Test ES API with Endpoint /es/:target/_doc 
   1. [Test Create Document with not existing Index](./test_docs/ES_API_3.md)
   2. [Test Create Document with existing Index](./test_docs/ES_API_3.md)
   3. [Test Create Document with existing Index and not existing Document ID](./test_docs/ES_API_3.md)
   4. [Test Update Document with existing Index and existing Document ID](./test_docs/ES_API_3.md)
   5. [Test Create Document with error Input](./test_docs/ES_API_3.md)

4. #### Test ES API with Endpoint /es/:target/_doc/:id
   1. [Test Update Document with no existing Index](./test_docs/ES_API_4.md)
   2. [Test Update Document with existing Index](./test_docs/ES_API_4.md)
   3. [Test Create Document with existing Index and not existing ID](./test_docs/ES_API_4.md)
   4. [Test Update Document with existing Index and existing ID](./test_docs/ES_API_4.md)
   5. [Test Update Document with error Input](./test_docs/ES_API_4.md)

5. #### Test ES API with Endpoint /es/:target/_doc/:id
   1. [Test Delete Document with not existing Index](./test_docs/ES_API_5.md)
   2. [Test Delete Document with existing Index and not existing ID](./test_docs/ES_API_5.md)
   3. [Test Delete Document with existing Index and existing ID](./test_docs/ES_API_5.md)

6. #### Test ES API with Endpoint /es/:target/_create/:id (PUT)
   1. [Test Update Document with no existing Index](./test_docs/ES_API_6.md)
   2. [Test Update Document with existing Index](./test_docs/ES_API_6.md)
   3. [Test Update Document with existing Index and not existing ID](./test_docs/ES_API_6.md)
   4. [Test Update Document with existing Index and existing ID](./test_docs/ES_API_6.md)
   5. [Test Update Document with error Input](./test_docs/ES_API_6.md)

7. #### Test ES API with Endpoint /es/:target/_create/:id (POST)
   1. [Test Update Document with no existing Index](./test_docs/ES_API_7.md)
   2. [Test Update Document with existing Index](./test_docs/ES_API_7.md)
   3. [Test Update Document with existing Index and not existing ID](./test_docs/ES_API_7.md)
   4. [Test Update Document with existing Index and existing ID](./test_docs/ES_API_7.md)
   5. [Test Update Document with error Input](./test_docs/ES_API_7.md)

8. #### Test ES API with Endpoint /es/:target/_update/:id 
   1. [Test Update Document with no existing Index](./test_docs/ES_API_8.md)
   2. [Test Update Document with existing Index](./test_docs/ES_API_8.md)
   3. [Test Update Document with existing Index and not existing ID](./test_docs/ES_API_8.md)
   4. [Test Update Document with existing Index and existing ID](./test_docs/ES_API_8.md)
   5. [Test Update Document with error Input](./test_docs/ES_API_8.md)
---
### Test Document Searchv1 (Test File : search_v1_test.go)
---
1. #### Test Search with Endpoint /api/:target/_search
   1. [Test Search Document with no existing Index](./test_docs/SEARCH_V1_1.md)
   2. [Test Search Document with existing Index](./test_docs/SEARCH_V1_1.md)
   3. [Test Search Document with no existing Term](./test_docs/SEARCH_V1_1.md)
   4. [Test Search Document with existing Term](./test_docs/SEARCH_V1_1.md)
   5. [Test Search Document with search type all documents](./test_docs/SEARCH_V1_1.md)
   6. [Test Search Document with search type wildcard](./test_docs/SEARCH_V1_1.md)
   7. [Test Search Document with search type fuzzy](./test_docs/SEARCH_V1_1.md)
   8. [Test Search Document with search type term](./test_docs/SEARCH_V1_1.md)
   9. [Test Search Document with search type data range](./test_docs/SEARCH_V1_1.md)
   10. [Test Search Document with search type match all](./test_docs/SEARCH_V1_1.md)
   11. [Test Search Document with search type match](./test_docs/SEARCH_V1_1.md)
   12. [Test Search Document with search type match phrase](./test_docs/SEARCH_V1_1.md)
   13. [Test Search Document with search type multiphrase](./test_docs/SEARCH_V1_1.md)
   14. [Test Search Document with search type prefix](./test_docs/SEARCH_V1_1.md)
   15. [Test Search Document with search type query string](./test_docs/SEARCH_V1_1.md)

2. #### Test Search Document with Aggregation
   1. [Test Search Document with Term Aggregation](./test_docs/SEARCH_V1_2.md)
   2. [Test Search Document with Metric Aggregation](./test_docs/SEARCH_V1_2.md)

---
### Test Document Searchv2 (Test File : search_v2_test.go)
---
1. #### Test Search Document with Endpoint /es/:target/_search
   1. [Test Search Document with no existing Index](./test_docs/SEARCH_V2_1.md)
   2. [Test Search Document with existing Index](./test_docs/SEARCH_V2_1.md)
   3. [Test Search Document with no existing Term](./test_docs/SEARCH_V2_1.md)
   4. [Test Search Document with existing Term](./test_docs/SEARCH_V2_1.md)
   5. [Test Search Document with search type match all](./test_docs/SEARCH_V2_1.md)
   6. [Test Search Document with search type wildcard](./test_docs/SEARCH_V2_1.md)
   7. [Test Search Document with search type fuzzy](./test_docs/SEARCH_V2_1.md)
   8. [Test Search Document with search type term](./test_docs/SEARCH_V2_1.md)
   9. [Test Search Document with search type data range](./test_docs/SEARCH_V2_1.md)
   10. [Test Search Document with search type match](./test_docs/SEARCH_V2_1.md)
   11. [Test Search Document with search type match phrase](./test_docs/SEARCH_V2_1.md)
   12. [Test Search Document with search type prefix](./test_docs/SEARCH_V2_1.md)
   13. [Test Search Document with search type query string](./test_docs/SEARCH_V2_1.md)

2. #### Test Search Document with Aggregation
   1. [Test Search Document with Term Aggregation](./test_docs/SEARCH_V2_2.md)
   2. [Test Search Document with Metric Aggregation](./test_docs/SEARCH_V2_2.md)   
--- 

### Test Analyzers (Test File : analyze_test.go)
---
1. #### [Test Analyzers](./test_docs/ANALYZER.md)
2. #### [Test Tokenizers](./test_docs/TOKENIZER.md)
3. #### [Test Token Filter](./test_docs/TOEKN_FILTER.md)

--- 

## Benchmark Test

### Test Benchmark (Test File : bulk_test.go)
---
   1. #### [Test Bulk Worker Function](./test_docs/BENCH_BULK.md)
---


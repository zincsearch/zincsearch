# curl -v -u admin:Complexpass#123 -POST -H "Content-Type: application/json" -d @load_test_data_2.json http://localhost:4080/es/books/_doc

# hey -n 20000 -c 50 -m POST -H "Content-Type: application/json" -D load_test_data_1.json http://localhost:4080/books/_doc
# -a "admin:Complexpass#123" ==> -H "Authorization: Basic YWRtaW46Q29tcGxleHBhc3MjMTIz"
hey -n 100 -c 10 -H "Authorization: Basic YWRtaW46Q29tcGxleHBhc3MjMTIz" -m POST -H "Content-Type: application/json" -D load_test_data_2.json http://localhost:4080/es/books/_doc

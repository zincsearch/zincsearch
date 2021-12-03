
# hey -n 20000 -c 50 -m POST -H "Content-Type: application/json" -D load_test_data_1.json http://localhost:4080/books/_doc
hey -n 20000 -c 50 -m POST -H "Content-Type: application/json" -D load_test_data_2.json http://localhost:4080/books/_doc


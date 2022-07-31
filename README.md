# Warehouse manager

A simple warehouse manager which provide two main endpoints, GetProducts and SellOneProduct.
</br>
It also provides two other endpoints to populate database. </br>

### Requirements:
1. install [docker](https://docs.docker.com/desktop/install/linux-install) and [docker-compose](https://docs.docker.com/compose/install/).
2. install [liquibase](https://www.liquibase.org/download).


### How to run:
1. run ```make run-server```.
2. run ```make migrate```

### How to run tests:
1. run ```make test``` at project directory.

### For further commands and help:
1. run ```make help```

### Further info:
This Project is a simple warehouse manager, its dependency is **PostgreSQL** database, so if you are not using docker-compose or instructions above for running the project
you need to run a postgresql somewhere that be accessible by service. </br>
By default if you use docker to run the project, docker-compose command will bring up a **PostgreSQL** database first and then connects the service to it.


## Endpoints
1. ```POST /products``` used for populating products table.
2. ```GET /products``` used for getting all products and quantity of availability.
3. ```POST /products/sell``` used for selling a product.
4. ```POST /articles``` used for populating articles table.

### TODO (for future development): 
1. Change **CreateOrUpdateArticles** and **CreateOrUpdateProducts** Endpoints so that they can handle large json files
2. It will be nice to add pagination to **GetAllProducts** Endpoint
3. Optimize Database queries
4. Nice to have Integration test
5. Think and discuss how to scale Database when number of products or articles increase to millions or more
6. Nice to add metrics so that we can have monitoring.

# UATI - BACKEND
## A project for leads analysis and management, part of Codenation Acelera Dev program.


This project was built for a group of 4 people participating in a full-stack (React - Golang) program.
We used the following technologies:

 - Golang (API + some services)
 - Python (Web scrappy, data analysis and services)
 - RabbitMQ (messages)
 - MongoDB
 - Docker
![enter image description here](https://i.ibb.co/ScyZ7yp/tec.png)
This solutions consists in services and an REST API used to get and manipulate data from São Paulo public agents and see if they matches any client in a bank database. It also analyze those data using descriptive statistics to find useful information about public agents, potential leads and bank clients.

![architecture](https://i.ibb.co/Y3ZN0ts/arct.png)

You will have to create/configure a ".env" file.
 
You can run services locally using docker-compose (services/docker-compose.yml).

The main API can be started with "main.go".

# Example: Compose sample application

The example is based on the "react-express-mysql" example in Docker's [Awesome Compose](https://github.com/docker/awesome-compose) repository.
It consists of a simple full-stack application with a React application, a NodeJS backend, and a MariaDB database. 
I've also added a phpMyAdmin database management UI to make it a little bit more interesting. 

## The configuration

To clearly show what is part of the traditional Docker Compose setup and what is added for Laebel,
I've separated the two parts.
- [`compose.yaml`](./compose.yaml): The typical Compose file for the example. It defines the services and their configuration.
- [`compose.override.yaml`](./compose.override.yaml): The Laebel configuration file. It extends the original Compose file with Laebel-specific configuration.

Note that it is perfectly fine to have the Laebel configuration in the same file as the Compose configuration.

## The output

[View the full documentation site](https://rawcdn.githack.com/henrikje/laebel/58aba06ade6bb94263164ffedde2560c74956f2b/examples/react-express-mysql/laebel-output.html) generated by Laebel for this example.

<a href="https://rawcdn.githack.com/henrikje/laebel/58aba06ade6bb94263164ffedde2560c74956f2b/examples/react-express-mysql/laebel-output.html" title="Click to see the full generated documentation site."><img src="./laebel-example-screenshot.png" alt="Laebel output screenshot"></a>

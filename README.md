# docodash

Displays a dashboard with information about all services running in a Docker Compose cluster. 

## TODO

- [ ] Make the display nicer visually
- [ ] Load the current container id on startup (it won't change)
- [ ] Add special labels that services can use to display clickable links
  - For example, a basic web server can link to its homepage
  - A service can link both to its primary endpoint and to its API documentation
- [ ] Display a nice error message if the user forgets to mount the Docker socket
- [ ] Display a nice error message if the user does not mount the Docker socket readonly
- [ ] Display a nice error message if the container is not run as part of a Docker Compose cluster
- [ ] Display the project name in the title
- [ ] Introduce a struct to hold the information passed to the template
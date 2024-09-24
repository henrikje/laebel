# docodash

Displays a dashboard with information about all services running in a Docker Compose cluster. 

## TODO

- [ ] Make the display nicer visually, perhaps something like GitHub READMEs
- [ ] Load the current container id on startup (it won't change)
- [x] Add special labels that services can use to display clickable links
  - For example, a basic web server can link to its homepage
  - A service can link both to its primary endpoint and to its API documentation
- [ ] Display a nice error message if the user forgets to mount the Docker socket
- [ ] Display a nice error message if the user does not mount the Docker socket readonly
- [ ] Display a nice error message if the container is not run as part of a Docker Compose cluster
- [ ] Display the project name in the title
- [ ] Introduce a proper struct hierarchy to hold the information passed to the template.
- [ ] Display container names as well as their id
- [ ] Add a custom "description" label that can be used to display a short description of the service.
- [ ] Add support for Markdown.
- [ ] Live-reload the page when the Docker Compose cluster changes
- [x] Support all relevant [OpenContainer labels](https://github.com/opencontainers/image-spec/blob/main/annotations.md)
- [ ] Support `net.henko.docodash` alternatives to the opencontainers labels, for those who only would use this tool
- [ ] Change name? Perhaps "docodogen"? "docodoc"? "composedoc"?
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Hide the Docodash service itself from the list of services

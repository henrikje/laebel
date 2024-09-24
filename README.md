# docodash

Displays a dashboard with information about all services running in a Docker Compose cluster. 

## TODO

- [x] Add special labels that services can use to display clickable links
  - For example, a basic web server can link to its homepage
  - A service can link both to its primary endpoint and to its API documentation
- [x] Make the display nicer visually, perhaps something like GitHub READMEs
- [x] Support all relevant [OpenContainer labels](https://github.com/opencontainers/image-spec/blob/main/annotations.md)
- [x] Display a nice error message if the user forgets to mount the Docker socket
- [x] Display a nice error message if the container is not run as part of a Docker Compose cluster
- [x] Display the project name in the title
- [x] Introduce a proper struct hierarchy to hold the information passed to the template.
- [x] Display container names as well as their id
- [x] Hide the Docodash service itself from the list of services
- [ ] Display a nice error message if the user does not mount the Docker socket readonly
- [ ] Add support for Markdown.
- [ ] Live-reload the page when the Docker Compose cluster changes
- [ ] Support `net.henko.docodash` alternatives to the opencontainers labels, for those who only would use this tool
- [ ] Change name? Perhaps "docodogen"? "docodoc"? "composedoc"?
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Find a way to set an icon/logo for the project (and services?).
- [ ] Make service title a clickable link to the "website" if present?
- [ ] Create a /health endpoint and configure HEALTHCHECK in the Dockerfile
- [ ] Embed Mermaid.js script to avoid external dependency
- [ ] Add status icons for each service, both in the graph and in the list:
   🟢=healthy, ▶️=running, ⏸️=paused, ⏹️=stopped/exited, 🚫=unhealthy, 🔄=restarting, *️⃣=mixed
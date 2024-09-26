# TODO

## v0.3.0

- [ ] Add status icons/colors for each service, both in the graph and in the list:
  🟢=healthy, ▶️=running, ⏸️=paused, ⏹️=stopped/exited, 🚫=unhealthy, 🔄=restarting, *️⃣=mixed

## Future

- [ ] Add support for Markdown.
- [ ] Live-reload the page when the Docker Compose cluster changes; perhaps use HTMX?
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.

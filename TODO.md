# TODO

## v0.3.0

- [x] Add status icons/colors for each service, both in the graph and in the list.
- [x] Add a health column to the list of containers. (Any other interesting things now that we inspect each container?)
- [ ] Update the example to include status icons.

## Future

- [ ] Add support for Markdown.
- [ ] Live-reload the page when the Docker Compose cluster changes; perhaps use HTMX?
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Consider adding ports, volumes, and networks to the service details.
  - They can be left out if they are not used, or if all services have the same value (e.g., using the default network).
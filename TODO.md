# TODO

## v0.4.0

- [ ] Automatically refresh the page on changes.
  - Use the Docker Event API to listen for changes.
  - Implement SSE to tell the page to refresh.
  - Use HTMX to refresh the page.

## Future

- [ ] Add support for Markdown.
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Consider adding ports, volumes, and networks to the service details.
  - They can be left out if they are not used, or if all services have the same value (e.g., using the default network).

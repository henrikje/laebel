# TODO

## v0.4.0

- [ ] Automatically refresh the page on changes.
  - Use the Docker Event API to listen for changes.
  - Implement SSE to tell the page to refresh.
  - Use HTMX to refresh the page.
  - Is it possible to update only the status fields, ideally even without redrawing the Mermaid graph?
  - Display a banner when the SSE connection is lost, including a link to reconnect.

## Future

- [ ] Add support for Markdown.
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Consider adding ports, volumes, and networks to the service details.
  - They can be left out if they are not used, or if all services have the same value (e.g., using the default network).
- [ ] Is it possible to render the Mermaid graph on the server and send it as an image to the client?
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
- [ ] Perhaps add a "hidden details" key/value section, which can be used for any additional information that can be useful, but does not deserve permanent space in the main view. Things like maintainer, and version+revision.
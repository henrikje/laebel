# TODO

## v0.4.0

- [ ] Automatically refresh the page on changes.
  - Use the Docker Event API to listen for changes.
  - Regenerate changed parts of the page.
  - Implement SSE to push changes.
  - Use HTMX to refresh the parts that have changed.
  - Is it possible to update only the status fields, ideally even without redrawing the Mermaid graph?
- [ ] Display a banner when the SSE connection is lost, including a link to reconnect.

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
- [ ] Add a "last updated" timestamp to the page.
- [ ] Consider changing the default port to 8000 as 8080 is often used by other services.
- [ ] Update model so a service can have multiple values for the same label/property.
      For example, image and group name.
      There is no guarantee that all containers based on the same service have the same image.
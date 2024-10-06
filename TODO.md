# TODO

## v0.4.0

- [x] Automatically refresh the page on changes.
  - Use the Docker Event API to listen for changes.
  - Regenerate changed parts of the page.
  - Implement SSE to push changes.
  - Use HTMX to refresh the parts that have changed.
  - Is it possible to update only the status fields, ideally even without redrawing the Mermaid graph?
- [x] Consider changing the default port to 8000 as 8080 is often used by other services.
- [x] Remove verbose logging for each response or event.
- [x] Hide the Laebel service from the list. It is probably not interesting to most users, and it adds noise to the graph.

## Future

- [ ] Add support for Markdown.
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Consider adding ports, volumes, and networks to the service details.
  - They can be left out if they are not used, or if all services have the same value (e.g., using the default network).
- [ ] Is it possible to render the Mermaid graph on the server and send it as an image to the client?
      This would make it easier to smoothly update the page without flicker.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
- [ ] Perhaps add a "hidden details" key/value section, which can be used for any additional information that can be useful, but does not deserve permanent space in the main view. Things like maintainer, and version+revision.
- [ ] Add a "last updated" timestamp to the page.
- [ ] Update model so a service can have multiple values for the same label/property.
      For example, image and group name.
      There is no guarantee that all containers based on the same service have the same image.
- [ ] Add some kind of disconnected state to the client, so it can show a message when the connection is lost.
- [ ] Add a nicer presentation of a service which is depended on by another service, but not running.
- [ ] Update the service graph when the state of a service changes. Can HTMX in combination with Mermaid.js do that? Ideally, I could update just the labels. Otherwise, I may need to refresh the whole graph. Or switch to a server-side rendered graph.
      I think it should be possible with a HTMX hx-get on the graph parent, triggered by any a "status" event, which loads a page with all statuses for all nodes. Then each node has a selector to update its status. The question is how to get the first htmx attributes into the generated graph. But that should be possible too, I think.
- [ ] Consider if we can download the external JS files during (Docker) build. That way we don't need to have Mermaid.js and HTMX manually copied into the repository.
- [ ] Replace the PNG logo with an SVG logo.
- [ ] Lighthouse: [Enable text compression](https://developer.chrome.com/docs/lighthouse/performance/uses-text-compression/)
- [ ] Lighthouse: [Eliminate render-blocking resources](https://developer.chrome.com/docs/lighthouse/performance/render-blocking-resources/)
- [ ] Lighthouse: Image elements do not have explicit width and height
- [ ] Lighthouse: [Serve static assets with an efficient cache policy](https://developer.chrome.com/docs/lighthouse/performance/uses-long-cache-ttl/)
- [ ] Lighthouse: [Document does not have a meta description](https://developer.chrome.com/docs/lighthouse/seo/meta-description/)
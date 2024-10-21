# TODO

## v0.8.0

- [x] Update graph so icons display status, and color displays health.
- [x] Add a legend to the service graph, explaining what the different colors and icons mean.

## Future

- [ ] Add support for Markdown.
- [ ] Add a sort order label to the services, so they can be displayed in a specific order.
      Otherwise, sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Is it possible to render the Mermaid graph on the server and send it as an image to the client?
      This would make it easier to smoothly update the page without flicker.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
      Perhaps a topographical order of services within each group too?
- [ ] Perhaps add a "hidden details" key/value section, which can be used for any additional useful information, 
      that does not deserve permanent space in the main view. Things like maintainer, and version+revision.
- [ ] Add a "last updated" timestamp to the page.
- [ ] Add a nicer presentation of a service which is depended on by another service, but not running.
- [ ] Consider if we can download the external JS files during (Docker) build.
      That way we don't need to have Mermaid.js and HTMX manually copied into the repository.
- [ ] Lighthouse: [Enable text compression](https://developer.chrome.com/docs/lighthouse/performance/uses-text-compression/)
- [ ] Lighthouse: [Eliminate render-blocking resources](https://developer.chrome.com/docs/lighthouse/performance/render-blocking-resources/)
- [ ] Lighthouse: Image elements do not have explicit width and height
- [ ] Lighthouse: [Serve static assets with an efficient cache policy](https://developer.chrome.com/docs/lighthouse/performance/uses-long-cache-ttl/)
- [ ] Lighthouse: [Document does not have a meta description](https://developer.chrome.com/docs/lighthouse/seo/meta-description/)
- [ ] Add "icon description" title and help cursor in the service graph, just like the status summary icon in the service section has.
- [ ] Why does HTMX request the `hx-get` for the services multiple times? It should only be once per service status event.
- [ ] Add a "writer mode" where you can edit all descriptions, and the system produces the labels for you to paste into the compose file.
- [ ] Try generating the images server-side.
  - https://github.com/mermaid-js/mermaid-cli
- [ ] Look into sending out updated Mermaid and use the Mermaid JS API to update the graph.
  - Use shelved solution to update the graph. 
  - [ ] Also ensure the service list is properly updated.
- [ ] Update model so a service can have multiple values for the same label/property.
  For example, image and group name.
  There is no guarantee that all containers based on the same service have the same image.
  Alternatively, we can log a warning if a service has multiple values for the same label and only use the first one. 

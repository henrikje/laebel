# TODO

## v0.7.0

- [x] Support `net.henko.laebel.hidden` flag for volumes and networks
- [x] Look up port description based on the _container_ port number, not the _host_ port number.
- [ ] Display a banner when the SSE connection is lost (e.g. Laebel is shut down) to indicate that the status is no longer live-updated.
  Perhaps replace all status icons with a "live update lost" (e.g. ❓) icon. (Or is it helpful to know the last known status?)
- [ ] Update model so a service can have multiple values for the same label/property.
  For example, image and group name.
  There is no guarantee that all containers based on the same service have the same image.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Fix reader mode: It should display complete content in reader mode.
- [ ] Add dark mode: Switch to dark mode automatically when requested.
- [ ] Display a message when SSE connection is broken.

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
- [-] Add `org.opencontainers.image.authors` support, perhaps displayed as "Image: <image> by <author>". _No, not really necessary. It is not something a human would have added to README.md._
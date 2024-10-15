# TODO

## v0.6.0

- [x] Add support for volumes.
- [x] Add support for networks.
- [ ] Ensure that a service, a volume, and a networks with the same name are not confused.
- [x] Document supported volume and network labels.
- [ ] Consider adding volumes and/or networks to service groups (using `net.henko.laebel.group`).
      If so, they should perhaps be renamed "groups" rather than "service groups".
- [x] Visualize volumes and networks in the service graph.
      - It could feel natural to represent networks as subgraphs, but since an element may not be part of multiple subgraphs, this will not work.
      - Each network could be represented as its own node, but then there will be a lot o extra edges, making the graph harder to read.
      - We could add “badges” to the serices which note what networks they use. (This could work for volumes too.) The badges could either have a text prefix (“network: “) or a emoji prefix (e.g. “🌐 “ or “🛜 “ for network, or "🫙 “, “📦 “, or “🛄 ” for volume).
- [ ] Consider adding status for volumes and networks. What states can they be in? Created/missing?
- [x] Consider suppressing the "default" network if all services are using it.
- [ ] Avoid showing networks that are not explicitly used by any service.
- [ ] Display a list of services which use each volume and network.
- [ ] Add labels for volumes and networks in the react-express-mysql example. Update example.

## Future

- [ ] Add support for Markdown.
- [ ] Add a sort order label to the services, so they can be displayed in a specific order.
      Otherwise, sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Write unit tests for the project.
- [ ] Optimize Dockerfile for caching.
- [ ] Look for more characters that need to be included in the `escape` template function.
- [ ] Consider adding ports, volumes, and networks to the service details.
  - They can be left out if they are not used, or if all services have the same value (e.g., using the default network).
- [ ] Is it possible to render the Mermaid graph on the server and send it as an image to the client?
      This would make it easier to smoothly update the page without flicker.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
      Perhaps a topographical order of services within each group too?
- [ ] Perhaps add a "hidden details" key/value section, which can be used for any additional useful information, 
      that does not deserve permanent space in the main view. Things like maintainer, and version+revision.
- [ ] Add a "last updated" timestamp to the page.
- [ ] Update model so a service can have multiple values for the same label/property.
      For example, image and group name.
      There is no guarantee that all containers based on the same service have the same image.
- [ ] Add a nicer presentation of a service which is depended on by another service, but not running.
- [ ] Update the service graph when the state of a service changes. Can HTMX in combination with Mermaid.js do that? Ideally, I could update just the labels. Otherwise, I may need to refresh the whole graph. Or switch to a server-side rendered graph.
      I think it should be possible with a HTMX hx-get on the graph parent, triggered by any a "status" event, which loads a page with all statuses for all nodes. Then each node has a selector to update its status.
      The question is how to get the first htmx attributes into the generated graph. But that should be possible too, I think.
- [ ] Consider if we can download the external JS files during (Docker) build.
      That way we don't need to have Mermaid.js and HTMX manually copied into the repository.
- [ ] Lighthouse: [Enable text compression](https://developer.chrome.com/docs/lighthouse/performance/uses-text-compression/)
- [ ] Lighthouse: [Eliminate render-blocking resources](https://developer.chrome.com/docs/lighthouse/performance/render-blocking-resources/)
- [ ] Lighthouse: Image elements do not have explicit width and height
- [ ] Lighthouse: [Serve static assets with an efficient cache policy](https://developer.chrome.com/docs/lighthouse/performance/uses-long-cache-ttl/)
- [ ] Lighthouse: [Document does not have a meta description](https://developer.chrome.com/docs/lighthouse/seo/meta-description/)
- [ ] Add "icon description" title and help cursor, just like the status summary icon in the service section has.
- [ ] Why does HTMX request the `hx-get` for the services multiple times? It should only be once per service status event.
- [ ] Display a banner when the SSE connection is lost (e.g. Laebel is shut down) to indicate that the status is no longer live-updated.
      Perhaps replace all status icons with a "live update lost" (e.g. ❓) icon. (Or is it helpful to know the last known status?)
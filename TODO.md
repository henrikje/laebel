# TODO

## v0.9.0 New design

- [ ] New cleaner and more readable design.
  - Inspired by Simple.css.
  - The design still has a README feel to it.
- [ ] Basic semantic markup as a good foundation to build on.
- [x] Rename project "icon" to "logo" to avoid confusion with status icons.
- [ ] Ensure the design looks good even when there are no metadata labels. (Or with no services at all.)
- [x] Make the instance table more narrow so it fits in the design.
  - Remove the "Port Bindings" column from the container table to make it fit better.
  - Merge the "Status" and "Health" columns into one. Like "running (healthy)".
- [ ] Lighthouse: [Enable text compression](https://developer.chrome.com/docs/lighthouse/performance/uses-text-compression/)
- [ ] Lighthouse: [Eliminate render-blocking resources](https://developer.chrome.com/docs/lighthouse/performance/render-blocking-resources/)
- [ ] Lighthouse: Image elements do not have explicit width and height
- [ ] Lighthouse: [Serve static assets with an efficient cache policy](https://developer.chrome.com/docs/lighthouse/performance/uses-long-cache-ttl/)
- [ ] Lighthouse: [Document does not have a meta description](https://developer.chrome.com/docs/lighthouse/seo/meta-description/)

## v0.10.0 Optimized updating

- [ ] Update only status when possible.
- [ ] Update everything when necessary (no "reload" notice).
- [ ] Intercept the new Mermaid source and render it to SVG before updating the page. 
- [ ] Move as much view logic as possible to the server. Templates should be used for rendering only.

## v0.11.0 Improved build

- [ ] Optimize Dockerfile for caching.
- [ ] Consider if we can download the external JS files during (Docker) build.
  That way we don't need to have Mermaid.js and HTMX manually copied into the repository.

## v1.0.0-rc1 Stability

Focus on making it as stable as possible, fixing bugs, and handling edge cases.

## Future

- [ ] Add support for Markdown.
- [ ] Write unit tests for the project.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
      Perhaps a topographical order of services within each group too?
- [ ] Add a "writer mode" where you can edit all descriptions, and the system produces the labels for you to paste into the compose file.
- [ ] Update model so a service can have multiple values for the same label/property.
  For example, image and group name.
  There is no guarantee that all containers based on the same service have the same image.
  Alternatively, we can log a warning if a service has multiple values for the same label and only use the first one.

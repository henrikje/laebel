# TODO

## v0.10.0 Optimized updating

- [ ] Update only status when possible.
- [ ] Update everything when necessary (no "reload" notice).
- [ ] Intercept the new Mermaid source and render it to SVG before updating the page. 
- [ ] Move as much view logic as possible to the server. Templates should be used for rendering only.

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

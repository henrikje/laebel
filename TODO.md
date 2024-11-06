# TODO

This document acts as a roadmap and primitive issue tracker for the project.

## v0.10.0 Optimized updating

- [x] Update everything on any event, to ensure the presented view stays consistent.
- [-] Use idiomorph to ensure minimal changes to the DOM. _Not really needed._
- [x] Intercept the new Mermaid source and render it to SVG before updating the page.
- [x] Avoid flickering when updating the graph.
- [x] Order graph detail entries to avoid random changes in the graph.

## v1.0.0-rc1 Stability

Focus on making it as stable as possible, fixing bugs, and handling edge cases.

- [ ] Use Docker API to determine actual port binding for Laebel, so it can be displayed on startup.

## Future ideas

- [ ] Add support for Markdown.
- [ ] Write unit tests for the project.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
      Perhaps a topographical order of services within each group too?
- [ ] Add a "writer mode" where you can edit all descriptions, and the system produces the labels for you to paste into the compose file.
- [ ] Update model so a service can have multiple values for the same label/property.
  For example, image and group name.
  There is no guarantee that all containers based on the same service have the same image.
  Alternatively, we can log a warning if a service has multiple values for the same label and only use the first one.

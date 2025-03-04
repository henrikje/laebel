# TODO

This document acts as a roadmap and primitive issue tracker for the project.

## 0.11.1 Fixed network and volume presentation

- [ ] Join multiple networks or volumes into a single comma-separated row in the graph. (Like ports.)
- [ ] Fix mixed-up icons between networks and volumes in the graph.
- [ ] Volumes (and networks?) have lost their icon in the resource list.
- [ ] Drop the "Resources:" prefix for project-level resources.

## 0.12.0 Service logs

- [ ] Add service log view: a separate page which displays logs for a selected service.
- [ ] Make service log view live, so it updates as new logs are written.

## 1.0.0-rc1 Stability

Focus on making it as stable as possible, fixing bugs, and handling edge cases.

## Future ideas

- [ ] Add support for Markdown or HTML.
- [ ] Write unit tests for the project.
- [ ] Sort service groups by topological order; the main service should be at the top with its dependencies below.
      Perhaps a topographical order of services within each group too?
- [ ] Add a "writer mode" where you can edit all descriptions, and the system produces the labels for you to paste into the compose file.
- [ ] Update model so a service can have multiple values for the same label/property.
  - For example, image and group name.
  - There is no guarantee that all containers based on the same service have the same image.
  - Alternatively, we can log a warning if a service has multiple values for the same label and only use the first one.
- [ ] Use [template fragments](https://gist.github.com/benpate/f92b77ea9b3a8503541eb4b9eb515d8a) to simplify templates.
- [ ] Consider [embedding static files and using built-in routing](https://jvns.ca/blog/2024/09/27/some-go-web-dev-notes/).
- [ ] Improve dimensions of example image to show some of the service list as well. 
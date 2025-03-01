# TODO

This document acts as a roadmap and primitive issue tracker for the project.

## 0.10.1 Improved graph rendering

- [x] Use Docker API to determine actual port binding for Laebel, so it can be displayed on startup.
- [ ] Layout of the service graph changes after initial render: parts of service nodes become clipped.
- [ ] Initial render of service graph is slow. What is taking so much time? Can I preload the necessary JS?
- [ ] Change container table column "status" to "state".

## 0.11.0 Better links

- [ ] Update the version scheme from "v0.11.0" to "0.11.0".
- [ ] Split "Links" into "Resources" and "Links".
  - "Resources" displays "Website", "Documentation", and "Source code" as a comma-separated list. (Do the same for project links.)
  - "Links" displays all custom links in an unordered list.
- [ ] Add a description label for links, so they can be displayed `[Title](URL): Description`.

## 0.12.0 Service logs

- [ ] Add service log view: a separate page which displays logs for a selected service.
- [ ] Make service log view live, so it updates as new logs are written.

## 1.0.0-rc1 Stability

Focus on making it as stable as possible, fixing bugs, and handling edge cases.

## Future ideas

- [ ] Add support for Markdown.
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
- [ ] Add "featured links"?
  - Any service could contribute links for more important services.
  - For example, Jaeger could add Jaeger UI, Traefik could add its Dashboard, and client-web-app could link to the web page.

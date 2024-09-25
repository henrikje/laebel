# TODO

## Version 0.1.0

- [x] Select a license for the project.
- [ ] Add a GitHub Actions workflow to build and deploy the project to Docker Hub or GitHub Packages.
- [ ] Make the GitHub project public.

## Future

- [ ] Add support for Markdown.
- [ ] Live-reload the page when the Docker Compose cluster changes
- [ ] Add a sort order label to the services, so they can be displayed in a specific order. Otherwise sort alphanumerically. Perhaps also a sort order for groups?
- [ ] Find a way to set an icon/logo for the project (and services?).
- [ ] Add status icons for each service, both in the graph and in the list:
  🟢=healthy, ▶️=running, ⏸️=paused, ⏹️=stopped/exited, 🚫=unhealthy, 🔄=restarting, *️⃣=mixed
- [ ] Make http://laebel.henko.net/ the official website for the project.
- [ ] Set up an example Docker Compose project to demonstrate the tool. This can be used in examples and for manual testing.
- [ ] Ensure the Laebel service itself is displayed last. (Or just go by natural sort?)
- [ ] Add a favicon to the project.
- [ ] Add a `net.henko.laebel.hidden` label to hide certain services.
- [ ] Write unit tests for the project.
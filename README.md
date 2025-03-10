# Laebel

_Automatic documentation site for your Docker Compose project._

Laebel is a small server that runs in your Docker Compose project, serving a website that documents your project.

<figure>
<a href="./examples/react-express-mysql/README.md"><img src="./examples/react-express-mysql/laebel-example-screenshot.png" alt="Laebel output screenshot"></a>
<figcaption><em>An example of a documentation site generated by Laebel. <a href="./examples/react-express-mysql/README.md">View the full example</a>.</em></figcaption>
</figure>
<p></p>

Laebel uses the Docker API to get information about the services in your project,
and then displays that information in an easy-to-read format.
It even includes a graph of how the services are connected.
Additional information about the services can be added using labels in the Docker Compose project.

The goal is to make it feel like someone wrote a nice README for your project, but without the manual work. 😀

## Features

- **Displays a service graph**: An easy-to-understand, visual representation of how the services in your project are connected. 
  It is based on [`depends_on`](https://docs.docker.com/reference/compose-file/services/#depends_on) relations in your Docker Compose project. 
  It also uses icons and color to visualize the status and health of the services.
- **Live-updated status**: The status of the services is updated in real-time, so you can see when a service is starting, running, stopped, healthy, or unhealthy.
- **Describes each service in your project**: Lists all services and the important information about them.
  The information includes service name, image name, how many containers are running, and the status of the containers.
- **Describes each volume and network**: Also displays information about each volume and network in your project.
- **Additional user-configurable metadata**: Allows you to describe your services, volumes, and networks in more detail, 
  using [`labels`](https://docs.docker.com/reference/compose-file/services/#labels) in your Docker Compose project.
  It allows you to specify title, description, group, external links, and more.
  The metadata is based on the [OpenContainers Annotations Spec](https://specs.opencontainers.org/image-spec/annotations/) 
  but is extended through a set of Laebel-specific labels.
- **Light and dark mode**: Supports both light and dark themes, and automatically switches based on the user's preference.

## Usage

### As a service in a Docker Compose project

To get started, add the following service to your Docker Compose project:

```yaml
laebel:
    image: ghcr.io/henrikje/laebel:latest
    # Expose port 8000 to access the Laebel website
    ports:
    - "8000:8000"
    # Mount the Docker socket in read-only mode
    # This allows Laebel to detect the services in your project
    volumes:
    - "/var/run/docker.sock:/var/run/docker.sock:ro"
```

Then run `docker compose up` and open your browser at http://localhost:8000/ (or the host and port you use).

For a full example, see the [React-Express-MySQL example](./examples/react-express-mysql/README.md).

### As a stand-alone container

If you do not want to add Laebel to your Compose Project, you can run Laebel as a stand-alone container.
Then you need to tell Laebel which Compose Project to document using the environment variable `COMPOSE_PROJECT_NAME`.

```bash
docker run \
  -e COMPOSE_PROJECT_NAME=<your-project> \
  -v "/var/run/docker.sock:/var/run/docker.sock:ro" \
  -p 8000:8000 \
  ghcr.io/henrikje/laebel:latest
```

Note that in this case, you cannot put project metadata environment variables in the Compose Project.
You will have to manage those environment variables yourself.

## Configuration

Laebel will work out of the box and provide useful information about your project.
However, it really shines when you add additional metadata about your Docker Compose project.

### Documenting a project using labels

Laebel reads metadata about each service, volume, and network from labels in the Docker Compose project.
This is used to describe the resources in more detail than what can be derived automatically from Docker.

A short example of how to add documentation labels to a service:

```yaml
services:
  my-rest-api:
    image: my-rest-api:latest
    ports:
      - 80:80
    labels:
      org.opencontainers.image.title: My REST API
      org.opencontainers.image.description: My amazing REST API service that does cool things.
      net.henko.laebel.group:  Public Services
      net.henko.laebel.link.health.url: http://localhost/health
      net.henko.laebel.link.health.title: Health Check
      net.henko.laebel.port.80.description: The main HTTP port
```

The most straight-forward way to add documentation is to add labels in your primary `compose.yaml` file.
If you want to keep the labels separate from the main configuration, you can use a `compose.override.yaml` file.
Docker Compose will automatically [merge](https://docs.docker.com/compose/how-tos/multiple-compose-files/merge/) the two files
when you run `docker compose`.

Alternatively, you add labels to a service's `Dockerfile` so they will be built into the image:

```Dockerfile
LABEL org.opencontainers.image.title="My REST API" \
      org.opencontainers.image.description="My amazing REST API service that does cool things." \
      net.henko.laebel.group="Public Services" \
      net.henko.laebel.link.health.url="http://localhost/health" \
      net.henko.laebel.link.health.title="Health Check" \
      net.henko.laebel.port.80.description="The main HTTP port"
```

### Available labels for services

The following [OpenContainers Annotations Spec](https://specs.opencontainers.org/image-spec/annotations/) labels are supported:

- `org.opencontainers.image.title`: A human-readable title of the service.
- `org.opencontainers.image.description`: A longer description of the service.
- `org.opencontainers.image.url`: A URL to the service's homepage.
- `org.opencontainers.image.documentation`: A URL to the documentation of the service.
- `org.opencontainers.image.source`: A URL to the source code of the service.

In addition, Laebel supports the following custom labels:

- `net.henko.laebel.group`: A group name to categorize services.
  Services with the same group name will be displayed together, both in the service graph and in the list.
- `net.henko.laebel.hidden`: If set to `true`, the service will not be displayed in the graph or the list.

You can also add any number of external links.
These are great for linking to documentation, administration interfaces, or other related services.
Each link is specified with two labels, where `<key>` can be any string:

- `net.henko.laebel.link.<key>.url`: The URL of the link.
- `net.henko.laebel.link.<key>.title`: The title of the link.
- `net.henko.laebel.link.<key>.description`: A description of the link.

Finally, you can document the ports that are bound to the service.
This is especially helpful when it is not a well-known port.

- `net.henko.laebel.port.<port>.description`: A description of the purpose of the port.
  The `<port>` should be the port number that the container exposes (which may not be the same as the host port).

See the [full example](./examples/react-express-mysql/README.md) for examples on how to use these labels.

### Available labels for volumes and networks

Laebel also supports adding metadata to volumes and networks.

The following labels are supported:

- `net.henko.laebel.title`: A human-readable title of the volume/network.
- `net.henko.laebel.description`: A longer description of the volume/network.
- `net.henko.laebel.hidden`: If set to `true`, the label/network will not be included in the generated documentation.

### Project metadata through environment variables

Laebel also supports setting metadata for the project as a whole.
Since these values are not associated with any particular service, 
they are specified by adding _environment variables_ to the `laebel` service in the Docker Compose project.

- `LAEBEL_PROJECT_TITLE`: A human-readable title of the project.
- `LAEBEL_PROJECT_DESCRIPTION`: A description of the project.
- `LAEBEL_PROJECT_URL`: A URL to the project's homepage.
- `LAEBEL_PROJECT_DOCUMENTATION`: A URL to the documentation of the project.
- `LAEBEL_PROJECT_SOURCE`: A URL to the source code of the project.
- `LAEBEL_PROJECT_LOGO`: A URL to an image file to use as the logo of the project.
  Can be a `data:` URL to avoid external dependencies.

## Feedback

If you have any thoughts or questions, please [reach out to me](https://henko.net/contact/).

_Sidenote_: Why the name _Laebel_? 
It is a reference to the idea that to label something is to explain what it is,
combined with the fact that Laebel relies on Docker _labels_ to get information about the project and its services.
However, I also wanted to be a cool kid and use a [digraph](https://en.wikipedia.org/wiki/Digraph_(orthography)) like Traefik and Jaeger. 😉

## License

This project is licensed under the terms of the [MIT license](LICENSE.md).

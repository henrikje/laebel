# Astra

The Astra project is a decision support system for performing medical triages.

The foundation of the project is the `engine-library` which can be embedded in various applications to provide decision support.
It can be used for custom implementations by partners. Predicare provides `client-webapp` and `engine-restapi` as official implementations.

---

## traefik <small>[otel]</small>
<!-- net.henko.docodash.subtitle -->
<h3 style="color: gray; margin-top: -0.5em;">Load balancer, reverse proxy, and TLS termination.</h3>

<!-- net.henko.docodash.description -->
Traefik acts as a reverse proxy and load balancer for the services in the cluster.
It terminates TLS connections and routes incoming requests to the appropriate service.
It automatically renews certificates with Let's Encrypt.

It also provides a dashboard that displays the current configuration and the status of the services.

- **Image** `traefik:v3.1.4`  
- **Status**: 1 container running (healthy)
- **Ports**: 443, 8080
- **Links**: [Traefik dashboard](http://localhost:8080/dashboard/) <!-- net.henko.docodash.link.<key>.[url|label] -->  

<details>
<summary>Expand for details on 1 container</summary>
<pre>
3d1fb0c71bbf   traefik:v3.1.4                                             "/entrypoint.sh --ap…"   3 hours ago      Up 16 seconds (healthy)
</pre>
<em><strong>Note</strong>: Should be a more hierarchical view of the containers, with all information given by the Docker API.</em>
</details>

---

## client-webapp <small>[traefik]</small> <small>[otel]</small>
<h3 style="color: gray; margin-top: -0.5em;">Official web application implementation.</h3>

A React-based web application that uses the engine-library to let users perform triages. 

- **Image**: `astradockerregistry.azurecr.io/client-web-app:latest`    
- **Status**: 2 containers running (healthy)  
- **Links**: [Webapp](http://localhost/webapp/)

<details>
<summary>Expand for details on 2 containers</summary>
<pre>
e4018e21d8c4   astradockerregistry.azurecr.io/client-web-app:latest       "/docker-entrypoint.…"   8 hours ago      Up 16 seconds             80/tcp                                                                                                                                                    astra-client-web-app-1
d509883f0e93   astradockerregistry.azurecr.io/client-web-app:latest       "/docker-entrypoint.…"   8 hours ago      Up 16 seconds             80/tcp                                                                                                                                                    astra-client-web-app-2
</pre>
</details>

---

## engine-restapi <small>[traefik]</small> <small>[otel]</small>
<h3 style="color: gray; margin-top: -0.5em;">Official REST API implementation</h3>

An Express app that provides a FHIR RESTful compatible API for the engine-library.

- **Image**: `astradockerregistry.azurecr.io/engine-restapi:latest` 
- **Status**: 1 container running (healthy), <span style="color: darkred">1 container failed</span>  
- **Depends on**: `repository-library` <!-- DO WE WANT THIS? -->
- **Links**: [REST API endpoint](http://localhost/), [Swagger UI](http://localhost/swagger/), [Official documentation]()

<details>
<summary>Expand for details on 2 containers</summary>
<pre>
46aeafed2d52   astradockerregistry.azurecr.io/engine-restapi:latest       "docker-entrypoint.s…"   3 hours ago      Up 16 seconds (healthy)
0f30d19c3b6b   astradockerregistry.azurecr.io/engine-restapi:latest       "docker-entrypoint.s…"   3 hours ago      Up 16 seconds (healthy)
</pre>
</details>

---

## repository-restapi <small>[traefik]</small> <small>[otel]</small>
<h3 style="color: gray; margin-top: -0.5em;">Backend REST API for package and license management.</h3>

An Express app that provides a backend REST API that serves packages, license information, and other data. 

- **Image**: `astradockerregistry.azurecr.io/repository-restapi:latest`  
- **Status**: 2 containers running (1 healthy, <span style="color: darkorange">1 unhealthy</span>)  
- **Links**: [REST API endpoint](http://repository.localhost/), [Swagger UI](http://repository.localhost/swagger/)  

<details>
<summary>Expand for details on 2 containers</summary>
<pre>
4da39bba9898   astradockerregistry.azurecr.io/repository-restapi:latest   "docker-entrypoint.s…"   3 hours ago      Up 16 seconds (healthy)
b5115ed88c87   astradockerregistry.azurecr.io/repository-restapi:latest   "docker-entrypoint.s…"   3 hours ago      Up 16 seconds (healthy)
</pre>
</details>

---

## jaeger
<h3 style="color: gray; margin-top: -0.5em;">Distributed OpenTelemetry-based tracing system.</h3>

All-in-one package that collects, stores, and visualizes traces.
All other services in the cluster are configured to send traces to Jaeger.

Note that storage is currently in-memory and not persisted between redeploys.

- **Image**: `jaegertracing/all-in-one:1.61.0`  
- **Status**: 1 container running (healthy)  
- **Ports**: 4317-4318, 14269, 16686
- **Links**: [Jaeger UI](http://localhost:16686/)  

<details>
<summary>Expand for details on 1 container</summary>
<pre>
1e82d853f208   jaegertracing/all-in-one:1.61.0                            "/go/bin/all-in-one-…"   2 days ago       Up 16 seconds (healthy)
</pre>
</details>

---

## test-e2e
<h3 style="color: gray; margin-top: -0.5em;">End-to-end tests for public services in the cluster.</h3>

End-to-end tests for the services in the cluster.
Exists in two flavors, REST-based tests and UI-based tests powered by Playwright.  

- **Image**: `astradockerregistry.azurecr.io/e2e-tests:latest`
- **Profile**: `test`
- **Status**: <span style="color: gray">1 container stopped</span>  
- **Links**: [Test report](http://localhost/test-results/)

<details>
<summary>Expand for details on 1 container</summary>
<pre>
1e82d853f208   astradockerregistry.azurecr.io/e2e-tests:latest            "/go/bin/all-in-one-…"   2 days ago       Up 16 seconds (healthy)
</pre>
</details>

---

## docodash [this]
<h3 style="color: gray; margin-top: -0.5em;">Displays information about the services running in this Docker Compose cluster.</h3>

Displays a dashboard with information about all services running in a Docker Compose cluster.
Intended to be used as a quick reference for developers.

- **Image**: `henko/docodash:latest`
- **Status**: 1 container running (healthy)  
- **Links**: [Docodash](http://localhost:8080/)

<details>
<summary>Expand for details on 1 container</summary>
<pre>
c71bbf3d1fb0   docodash:latest                                             "/entrypoint.sh --ap…"   3 hours ago      Up 16 seconds (healthy)
</pre>
</details>

---

_This page is generated by [Docodash](https://docodash.henko.net/), a tool to provide a quick overview and documentation of services running in a Docker Compose cluster._
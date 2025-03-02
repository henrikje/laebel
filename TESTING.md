# Testing

## Testing Laebel with watch

To make the feedback loop quicker, 
you can use the Docker Compose `watch` feature to automatically rebuild the project when you make changes.

```yaml
services:
  laebel:
    # ... other configuration
    build:
      context: /path/to/laebel/
      dockerfile: Dockerfile
    develop:
      watch:
        - action: rebuild
          path: /path/to/laebel/
````
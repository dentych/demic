# Pyramid

## Diagrams

**Sequence diagram for frontend connecting to backend over WebSocket**:

```mermaid
sequenceDiagram
  participant F as Frontend
  participant B as Backend

  F->>B: Connect WebSocket
  activate B
  F->>B: msg: Hello backend
  B-->>F: msg: Hello frontend
  deactivate B
```

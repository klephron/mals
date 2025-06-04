# mals-engine

## Configuration format

```js
{
  "models": [
    // global, connection is reused for selected workspace
    {
      "id": "llama-3.2:1B:Q6_K",
      "base_url": "http://localhost:8080/v1", // connection type is selected based on this
      "spec": "OpenAPI" // specification based on engine should communicate with model
      "settings": {
          // need to this about (e.g. context window size, max recommendations)
      }
    }
  ],
  "workspaces": {
    // this is applied for each workspace, for each workspace new language server is started
    "default": {
      "lsp_servers": [
        {
          "name": "bash-language-server",
          "filetypes": [
            // on which filetypes server will start
            "bash",
            "sh"
          ],
          "cmd": ["/home/user/bin/bash-language-server", "start"],
          "settings": {
            // passed directly for language server
          }
        }
      ],
      "model": "llama-3.2:1B:Q6_K" // connection is reused for workspace
    }
  }
}
```

## TODO

- Assign `"model": { object }`, same as in `"models"` to start model for each workspace (may be necessary if it accumulates context changes)

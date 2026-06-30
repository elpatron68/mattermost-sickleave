# Mattermost Sick Leave Plugin

Mattermost plugin for structured sick leave reporting (`Krankmeldung`) via slash commands and interactive dialogs.

**Plugin ID:** `com.elpatron68.mattermost-sickleave`

## Features

- `/sick-leave start` — initial report (first sick day)
- `/sick-leave update` — expected return date and AU certificate status
- `/sick-leave extend` — extend expected return date (with optional AU update)
- `/sick-leave status` — show active sick leave
- State machine enforced via KV store (A → B → C)
- HR channel posts via bot `@sickleave`; updates and extensions as thread replies
- German and English UI strings (server + webapp translations)
- System Console settings: HR channel, default locale, max backdate days

## Requirements

| | |
|---|---|
| Mattermost | 6.2.1+ |
| HR channel | Private channel ID configured in plugin settings |

## Slash commands

| Command | When available | Purpose |
|---------|----------------|---------|
| `/sick-leave start` | No active case | Report first sick day |
| `/sick-leave update` | After start (status: reported) | Expected return + AU certificate |
| `/sick-leave extend` | After update (status: updated or extended) | New expected return date |
| `/sick-leave status` | Any time | Show active case |
| `/sick-leave help` | Any time | Show help |

## Configuration

1. Enable plugin uploads in `config.json` (`PluginSettings.EnableUploads: true`)
2. Upload and enable the plugin
3. **System Console → Plugins → Sick Leave → Configure:**
   - **HR Channel ID** — private HR channel (required)
   - **Default Locale** — `en` or `de`
   - **Maximum Backdate Days** — default `3`

## Development

From the plugin root:

```bash
# Regenerate manifest after plugin.json changes
make apply

# Run tests
make test

# Production bundle (linux-amd64 + linux-arm64 + webapp)
make dist
# Output: dist/com.elpatron68.mattermost-sickleave-<version>.tar.gz

# Deploy to local Mattermost
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<token>
make deploy
```

`make dist` uses `build/package_bundle.py` to set executable bits on Linux plugin binaries in the archive (avoids `permission denied` on install).

For local development with a single-arch build:

```bash
export MM_SERVICESETTINGS_ENABLEDEVELOPER=true
make server
```

macOS server binaries for local Mattermost-on-Mac development:

```bash
make server-darwin
```

See [docs/docs/HANDOFF.md](docs/docs/HANDOFF.md) for session handoff notes and open product decisions.

## License

See [LICENSE](LICENSE).

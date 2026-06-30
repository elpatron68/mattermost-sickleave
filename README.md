# Mattermost Sick Leave Plugin

Mattermost plugin for structured sick leave reporting (`Krankmeldung`) via slash commands, a channel header menu, and date-picker dialogs.

**Plugin ID:** `de.medisoftware.mattermost-sickleave`
## Features

- Slash commands: `start`, `update`, `extend`, `end`, `status`, `help` (default trigger: `/sick-leave`)
- Channel header button with action menu (start / update / extend / end depending on case state)
- Custom webapp modals with HTML date pickers (Mattermost 10.5+)
- State machine per user via KV store: initial report (A) → update (B) → extension (C) → close
- HR channel posts via bot `@sickleave`; updates, extensions, and case closure as thread replies
- Configurable report hashtag (default `#krankmeldung`) and slash command trigger (e.g. `krankmeldung`)
- English and German UI (informal *Du* in DE)
- System Console settings: HR channel, locale, backdate limit, hashtag, command trigger

## Workflow variants

| Variant | Purpose | Fields |
|---------|---------|--------|
| **A** | Initial report | First sick day |
| **B** | Update after A | Expected return date, AU certificate (yes/no) |
| **C** | Extension after B | New expected return date, optional AU update |
| **End** | Close case | HR thread reply, active record cleared |

## Requirements

| | |
|---|---|
| Mattermost | 6.2.1+ (tested with 10.5 for webapp date pickers) |
| HR channel | Private channel ID configured in plugin settings |

## Slash commands

Default trigger is `sick-leave`; configure another name in plugin settings (e.g. `krankmeldung`).

| Command | When available | Purpose |
|---------|----------------|---------|
| `/… start` | No active case | Report first sick day |
| `/… update` | After start (status: reported) | Expected return + AU certificate |
| `/… extend` | After update (status: updated or extended) | New expected return date |
| `/… end` | Active case | Close case and notify HR |
| `/… status` | Any time | Show active case |
| `/… help` | Any time | Show help |

The channel header button opens the same actions in a menu.

## Configuration

1. Enable plugin uploads in `config.json` (`PluginSettings.EnableUploads: true`)
2. Upload and enable the plugin
3. **System Console → Plugins → Sick Leave → Configure:**
   - **HR Channel ID** — private HR channel (required)
   - **Default Locale** — `en` or `de`
   - **Maximum Backdate Days** — default `3`
   - **Report Hashtag** — default `#krankmeldung`
   - **Slash Command Trigger** — default `sick-leave`

After changing the plugin ID or uploading an update: disable → enable the plugin and hard-refresh the browser (Ctrl+Shift+R).

## Development

From the plugin root (use Node from `.nvmrc` for webapp builds):

```bash
# Regenerate manifest after plugin.json changes
make apply

# Run tests
make test

# Production bundle (linux-amd64 + linux-arm64 + webapp)
make dist
# Output: dist/de.medisoftware.mattermost-sickleave-<version>.tar.gz

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

### Releases

Bump version, update the changelog, commit, tag, and push:

```bash
./scripts/release.sh -v 0.1.2
```

Options: `--dry-run`, `--no-push`, `-y`, `--unsigned`. The Makefile also provides `make patch`, `make minor`, and `make major` for signed git tags without updating `plugin.json`.

## Project layout

```
mattermost-sickleave/
├── plugin.json
├── server/          # plugin, API, commands, dialogs, i18n, sickleave domain
└── webapp/          # channel header, menus, date-picker modals, i18n
```

## License

See [LICENSE](LICENSE).

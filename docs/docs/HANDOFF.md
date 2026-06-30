# Sick Leave Plugin — Session Handoff

Use this document to continue development in **WSL** (new Cursor chat or same project opened from Linux).

**Last updated:** 2026-06-30  
**Repo (fork):** https://github.com/elpatron68/mattermost-sickleave  
**Plugin ID:** `com.elpatron68.mattermost-sickleave`

---

## Project paths

| Environment | Path |
|-------------|------|
| Windows (current) | `C:\Users\markus.MEDISOFT\source\repos\mattermost-krankmeldung-slashbefehl\mattermost-sickleave` |
| WSL (same files via `/mnt/c`) | `/mnt/c/Users/markus.MEDISOFT/source/repos/mattermost-krankmeldung-slashbefehl/mattermost-sickleave` |

**Recommendation:** Open the folder in Cursor **from WSL** (`File → Open Folder` with `\\wsl$\<distro>\...` or the `/mnt/c/...` path). For faster builds, optionally clone/copy the repo into `~/src/mattermost-sickleave` inside WSL.

---

## What was planned

Mattermost plugin for structured sick leave (`Krankmeldung`) via slash command and action buttons:

| Variant | Purpose | Fields |
|---------|---------|--------|
| **A** | Initial report | First sick day (`date`) |
| **B** | Update after A | Expected end date, AU certificate yes/no |
| **C** | Extension after B | New expected end date, optional AU update |

- **Languages:** EN (default), DE  
- **Base:** [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template)

---

## What is already implemented (Phase 1 + 2)

- Bootstrap rename (`plugin.json`, `go.mod`, imports)
- `/sick-leave start` → interactive dialog (variant A)
- `/sick-leave update` → interactive dialog (variant B)
- `/sick-leave extend` → interactive dialog (variant C)
- `/sick-leave status`, `/sick-leave help`, `/sick-leave end`
- Date validation (no future start date, configurable backdate limit, expected end rules)
- KV store for active sick leave per user with state machine (A → B → C)
- HR channel post via bot `@sickleave`; B/C/D posted as thread replies
- Server i18n: `server/i18n/en.json`, `de.json`
- Webapp: custom date-picker modal, entry menu with action buttons, channel header button
- Webapp: `registerTranslations` (EN/DE)
- System Console settings: `HRChannelID`, `DefaultLocale`, `MaxBackdateDays`
- Unit tests for validation, state transitions, and command handlers
- Makefile aligned with `mattermost-transcribe` (bundle versioning, `package_bundle.py`, linux release builds)

### Phase 3 (implemented)

- Entry menu with action buttons (start / update / extend / end) via channel header button
- Channel header button (`registerChannelHeaderButtonAction`)
- Close case workflow (`/sick-leave end`) with HR thread reply and KV `ClearActive`

### Not yet implemented

- Configurable slash command name (e.g. `/krankmeldung`)
- Configurable hashtag per report

---

## Code layout

```
mattermost-sickleave/
├── plugin.json
├── docs/HANDOFF.md          ← this file
├── server/
│   ├── plugin.go            # OnActivate, bot, wiring
│   ├── configuration.go     # HR channel, locale, backdate
│   ├── api.go               # GET /api/v1/context, POST dialog/submit, POST /end
│   ├── command/             # Slash command + dialog submit + end
│   ├── dialog/              # OpenDialog builders
│   ├── i18n/                # Embedded EN/DE strings (server)
│   └── sickleave/           # Model, KV store, validation, HR post format
└── webapp/
    ├── i18n/en.json, de.json
    └── src/index.tsx        # channel header button, slash hook, translations
```

---

## WSL setup (required before `make`)

Checked on 2026-06-30: WSL has `make` and `npm`, but **`go` was not installed** in WSL.

```bash
# 1. Enter project (adjust distro name if needed)
cd /mnt/c/Users/markus.MEDISOFT/source/repos/mattermost-krankmeldung-slashbefehl/mattermost-sickleave

# 2. Install Go (example: Ubuntu)
sudo apt update
sudo apt install -y golang-go   # or install latest from https://go.dev/dl/

# 3. Node (template uses .nvmrc)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
source ~/.bashrc
nvm install
nvm use

# 4. Verify
go version
node -v
npm -v
make --version
```

---

## Build & deploy (Linux amd64 for Mattermost)

All commands from plugin root:

```bash
cd /mnt/c/Users/markus.MEDISOFT/source/repos/mattermost-krankmeldung-slashbefehl/mattermost-sickleave

# Regenerate manifest after plugin.json changes
go run ./build/manifest/main.go apply

# Run tests
make test

# Production bundle (linux-amd64 + webapp)
make dist
# Output: dist/com.elpatron68.mattermost-sickleave.tar.gz

# Deploy to local Mattermost (needs MM env vars / local mode)
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=<token>
make deploy
```

Upload `dist/*.tar.gz` via **System Console → Plugins** if not using `make deploy`.

---

## Mattermost configuration

1. Enable plugin uploads in `config.json` (`PluginSettings.EnableUploads: true`)
2. Upload and enable the plugin
3. **Plugins → Sick Leave → Configure:**
   - **HR Channel ID** — private HR channel (required)
   - **Default Locale** — `en` or `de`
   - **Maximum Backdate Days** — default `3`

### Manual test (Phase 1 + 2)

```
/sick-leave help
/sick-leave start    → dialog → submit → HR post + ephemeral confirmation
/sick-leave update   → dialog → submit → HR thread reply
/sick-leave extend   → dialog → submit → HR thread reply
/sick-leave end      → confirm → close case → HR thread reply
/sick-leave status
```

Channel header button opens the action menu (start / update / extend / end depending on state).

---

## Open product decisions (before Phase 2)

1. Who may report? Self only, or managers for others?
2. Can variant B be submitted multiple times (change expected end without C)?
3. How is a case closed? `/sick-leave end`, auto on date, HR-only?
4. AU field: boolean only, or enum (yes / no / pending)?
5. Notifications: HR channel only, or also DM to line manager?
6. Minimum Mattermost server version in production?

---

## Suggested next prompt (new Cursor chat in WSL)

Copy into a new chat after opening this folder from WSL:

```
Continue the Mattermost Sick Leave plugin in this repo.
Read docs/HANDOFF.md for context.

Phase 2:
- Implement /sick-leave update (variant B): expected end date + AU boolean
- Implement /sick-leave extend (variant C): new expected end after B
- Enforce state machine (A → B → C) via KV store
- Post HR updates as thread replies on the initial HR post
- Add EN/DE strings for all new UI text
- Add unit tests for B/C validation and state transitions

Build target: linux-amd64 via `make dist` in WSL.
```

---

## Chat / planning reference

Original planning covered: plugin vs alternatives (webhook, external forms, Playbooks), interactive dialogs, i18n strategy, GDPR notes for health data, and phased rollout. Phase 1 code matches that plan; Phase 2 items are listed above.

---

## Git remote

```
origin  https://github.com/elpatron68/mattermost-sickleave
```

Commit and push from WSL when ready (no commits were made automatically during Phase 1 setup).

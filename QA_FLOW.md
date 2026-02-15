# QA Flow — Web UI & API

Purpose: Provide a manual QA checklist and quick scripts to verify the web UI forms and backend flows (create lesson, create version, publish, get version, diff).

Prerequisites
- Services running: `docker compose up --build` and migrations applied.
- Gateway reachable at `http://127.0.0.1:8080` (adjust `QA_BASE_URL` in scripts if different).

Manual browser QA (forms in the UI)
- Open: `http://127.0.0.1:8080`
- Create Lesson
  - Enter content in the "Create Lesson" textarea (e.g. `parrot is awesome!`).
  - Click "Create Lesson".
  - Expected: `200 OK` and a JSON object with `root_id`, `version_number`, `content_id`. Displayed in the result box.
- Create Version
  - Paste the `root_id` from previous step into the "Root ID" field, add new content, click "Create Version".
  - Expected: JSON response with new `version_number` and `content_id`.
- Publish Version
  - Enter `root_id` and `version_number` to publish; click "Publish".
  - Expected: JSON status `{ "status": "published" }` (or success HTTP code).
- Get Version
  - Enter `root_id` and `version_number`, click "Get Version".
  - Expected: JSON describing the version and content id.
- Diff
  - Enter `from` and `to` content IDs, click "Diff".
  - Expected: JSON diff response from backend.

Quick API smoke tests (local scripts)
- There are two helper scripts: `qa/run_qa.sh` (bash) and `qa/run_qa.ps1` (PowerShell). They exercise the same flow used by the web UI and print HTTP status + response. Use whichever matches your environment.

Usage (bash)
```bash
cd c:\dev\lessons
bash qa/run_qa.sh
```

Usage (PowerShell)
```powershell
cd C:\dev\lessons
.\qa\run_qa.ps1
```

If any step fails, the scripts will print the non-2xx HTTP status and the response body. Also check the browser console and gateway/metadata logs.

Notes
- The scripts assume the gateway listens on `http://127.0.0.1:8080`.
- The PowerShell script uses `Invoke-RestMethod` (parsed JSON) and will show errors in a human-friendly way.

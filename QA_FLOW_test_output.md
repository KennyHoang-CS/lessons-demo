# QA Flow — Test Output
````markdown
# QA Flow — Test Output

Run: `.\qa\run_qa.ps1`
Running QA against http://127.0.0.1:8080

Summary: QA flow completed successfully. Below are the step outputs captured during the run.

1) Create Lesson

Response:

```json
{
    "content_id":  "278cfc7b-e808-4d18-ae20-088a977ad941",
    "root_id":  "0ccf7803-7a07-41f3-801d-94cac44d49e4",
    "version_number":  1
}
```

root_id=0ccf7803-7a07-41f3-801d-94cac44d49e4 content_id=278cfc7b-e808-4d18-ae20-088a977ad941

2) Create Version

Response:

```json
{
    "content_id":  "fab31e32-9b50-4747-a39d-e3141cac2821",
    "root_id":  "0ccf7803-7a07-41f3-801d-94cac44d49e4",
    "version_number":  2
}
```

3) Publish Version

Response:

```json
{
    "status":  "published"
}
```

4) Get Version

Response:

```json
{
    "root_id":  "0ccf7803-7a07-41f3-801d-94cac44d49e4",
    "version_number":  2,
    "content_id":  "fab31e32-9b50-4747-a39d-e3141cac2821"
}
```

5) Diff

Response:

```json
{
    "diff":  "- parrot is awesome!\n+ parrot v2\n"
}
```

6) Clone Lesson (from version 2)

Response:

```json
{
    "root_id":  "e8f31deb-c55c-42db-870f-c03a4725f555",
    "version_number":  1,
    "content_id":  "6f42d6c9-190c-43e6-b744-10f99a29204d"
}
```

7) Verify clone independence: make change on original and diff against clone

Create new version on original (version 3) Response:

```json
{
    "content_id":  "d0fa7f04-fa67-4bee-b991-5484a1eedc71",
    "root_id":  "0ccf7803-7a07-41f3-801d-94cac44d49e4",
    "version_number":  3
}
```

Diff against the clone:

```json
{
    "diff":  "- original changed\n+ parrot v2\n"
}
```

QA flow completed successfully (See attachments for file contents.)
````
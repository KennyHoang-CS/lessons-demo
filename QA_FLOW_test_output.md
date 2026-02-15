# QA Flow — Test Output

![QA Status](https://img.shields.io/badge/QA-Passed-4caf50?style=for-the-badge)

Run: `.\qa\run_qa.ps1`  
Target: `http://127.0.0.1:8080`  
Summary: **QA flow completed successfully**

---

## 1) Create Lesson

**Input:**
```http
POST /lessons
Body:
{
  "content": "parrotapp is awesome!"
}
```

**Output:**
```json
{
  "content_id": "278cfc7b-e808-4d18-ae20-088a977ad941",
  "root_id": "0ccf7803-7a07-41f3-801d-94cac44d49e4",
  "version_number": 1
}
```

---

## 2) Create Version

**Input:**
```http
POST /lessons/0ccf7803-7a07-41f3-801d-94cac44d49e4/versions
Body:
{
  "content": "parrot v2"
}
```

**Output:**
```json
{
  "content_id": "fab31e32-9b50-4747-a39d-e3141cac2821",
  "root_id": "0ccf7803-7a07-41f3-801d-94cac44d49e4",
  "version_number": 2
}
```

---

## 3) Publish Version

**Input:**
```http
POST /lessons/0ccf7803-7a07-41f3-801d-94cac44d49e4/versions/2/publish
```

**Output:**
```json
{
  "status": "published"
}
```

---

## 4) Get Version

**Input:**
```http
GET /lessons/0ccf7803-7a07-41f3-801d-94cac44d49e4/versions/2
```

**Output:**
```json
{
  "root_id": "0ccf7803-7a07-41f3-801d-94cac44d49e4",
  "version_number": 2,
  "content_id": "fab31e32-9b50-4747-a39d-e3141cac2821"
}
```

---

## 5) Diff (v1 → v2)

**Input:**
```http
GET /diff?from=1&to=2&root_id=0ccf7803-7a07-41f3-801d-94cac44d49e4
```

**Output:**
```json
{
  "diff": "- parrot is awesome!\n+ parrot v2\n"
}
```

---

## 6) Clone Lesson (from version 2)

**Input:**
```http
POST /lessons/clone
Body:
{
  "root_id": "0ccf7803-7a07-41f3-801d-94cac44d49e4",
  "version_number": 2
}
```

**Output:**
```json
{
  "root_id": "e8f31deb-c55c-42db-870f-c03a4725f555",
  "version_number": 1,
  "content_id": "6f42d6c9-190c-43e6-b744-10f99a29204d"
}
```

---

## 7) Verify Clone Independence

### 7a) Create new version on original (version 3)

**Input:**
```http
POST /lessons/0ccf7803-7a07-41f3-801d-94cac44d49e4/versions
Body:
{
  "content": "original changed"
}
```

**Output:**
```json
{
  "content_id": "d0fa7f04-fa67-4bee-b991-5484a1eedc71",
  "root_id": "0ccf7803-7a07-41f3-801d-94cac44d49e4",
  "version_number": 3
}
```

### 7b) Diff original v3 against clone v1

**Input:**
```http
GET /diff?from=3&to=1&root_id=e8f31deb-c55c-42db-870f-c03a4725f555
```

**Output:**
```json
{
  "diff": "- original changed\n+ parrot v2\n"
}
```

---

# 🎉 QA flow completed successfully
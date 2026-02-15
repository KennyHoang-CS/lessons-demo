async function createLesson() {
  const body = document.getElementById("lessonInput").value;
  try {
    const res = await fetch("/lessons", { method: "POST", body });
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("createLessonResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("createLessonResult").textContent = `ERROR: ${err.message}`;
  }
}

async function createVersion() {
  const rootId = document.getElementById("rootId").value;
  const body = document.getElementById("versionInput").value;
  try {
    const res = await fetch(`/lessons/${rootId}/versions`, { method: "POST", body });
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("createVersionResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("createVersionResult").textContent = `ERROR: ${err.message}`;
  }
}

async function publishVersion() {
  const rootId = document.getElementById("pubRootId").value;
  const version = document.getElementById("pubVersion").value;
  try {
    const res = await fetch(`/lessons/${rootId}/versions/${version}/publish`, { method: "POST" });
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("publishResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("publishResult").textContent = `ERROR: ${err.message}`;
  }
}

async function getVersion() {
  const rootId = document.getElementById("getRootId").value;
  const version = document.getElementById("getVersion").value;
  try {
    const res = await fetch(`/lessons/${rootId}/versions/${version}`);
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("getVersionResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("getVersionResult").textContent = `ERROR: ${err.message}`;
  }
}

async function getDiff() {
  const from = document.getElementById("diffFrom").value;
  const to = document.getElementById("diffTo").value;
  try {
    const res = await fetch(`/diff?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`);
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("diffResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("diffResult").textContent = `ERROR: ${err.message}`;
  }
}

async function parseJSONorText(res) {
  try {
    return await res.json();
  } catch (e) {
    try {
      return await res.text();
    } catch (e2) {
      return null;
    }
  }
}

async function cloneLesson() {
  const rootId = document.getElementById("cloneRootId").value;
  const version = parseInt(document.getElementById("cloneVersion").value || "1", 10);
  try {
    const res = await fetch(`/lessons/clone`, { method: "POST", headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ from_root_id: rootId, from_version: version }) });
    const out = res.ok ? await parseJSONorText(res) : await parseJSONorText(res);
    document.getElementById("cloneResult").textContent = `${res.status} ${res.statusText}\n${JSON.stringify(out, null, 2)}`;
  } catch (err) {
    document.getElementById("cloneResult").textContent = `ERROR: ${err.message}`;
  }
}
Param(
  [string]$BaseUrl = 'http://127.0.0.1:8080'
)

Write-Host "Running QA against $BaseUrl"

function Check-Response($result, $step) {
  if ($null -eq $result) { Write-Host "${step} failed: no response"; exit 1 }
  if ($result.StatusCode -lt 200 -or $result.StatusCode -ge 300) {
    Write-Host "${step} failed: $($result.StatusCode)"; Write-Host $result.Content; exit 1
  }
}

Write-Host "1) Create Lesson"
try {
  $res = Invoke-RestMethod -Uri "$BaseUrl/lessons" -Method Post -Body 'parrot is awesome!'
} catch {
  Write-Host "Create lesson request failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $res -Depth 4)

$root_id = $res.root_id
$content_id = $res.content_id
if (-not $root_id) { Write-Host "missing root_id"; exit 1 }
Write-Host "root_id=$root_id content_id=$content_id"

Write-Host "2) Create Version"
try {
  $res2 = Invoke-RestMethod -Uri "$BaseUrl/lessons/$root_id/versions" -Method Post -Body 'parrot v2'
} catch {
  Write-Host "Create version request failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $res2 -Depth 4)
$version_number = $res2.version_number
$version_content_id = $res2.content_id

Write-Host "3) Publish Version"
try {
  $res3 = Invoke-RestMethod -Uri "$BaseUrl/lessons/$root_id/versions/$version_number/publish" -Method Post
} catch {
  Write-Host "Publish failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $res3 -Depth 4)

Write-Host "4) Get Version"
try {
  $res4 = Invoke-RestMethod -Uri "$BaseUrl/lessons/$root_id/versions/$version_number" -Method Get
} catch {
  Write-Host "Get version failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $res4 -Depth 4)

Write-Host "5) Diff"
try {
  $res5 = Invoke-RestMethod -Uri "$BaseUrl/diff?from=$content_id&to=$version_content_id" -Method Get
} catch {
  Write-Host "Diff failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $res5 -Depth 4)

Write-Host "6) Clone Lesson (from version $version_number)"
try {
  $clone = Invoke-RestMethod -Uri "$BaseUrl/lessons/clone" -Method Post -Body (ConvertTo-Json @{ from_root_id = $root_id; from_version = $version_number }) -ContentType 'application/json'
} catch {
  Write-Host "Clone failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $clone -Depth 4)

$clone_root = $clone.root_id
$clone_content = $clone.content_id

Write-Host "7) Verify clone independence: make change on original and diff against clone"
try {
  $newVer = Invoke-RestMethod -Uri "$BaseUrl/lessons/$root_id/versions" -Method Post -Body 'original changed'
} catch {
  Write-Host "Create version on original failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $newVer -Depth 4)

$new_original_content = $newVer.content_id
try {
  $diffAfter = Invoke-RestMethod -Uri "$BaseUrl/diff?from=$new_original_content&to=$clone_content" -Method Get
} catch {
  Write-Host "Diff after original change failed: $_"; exit 1
}
Write-Host (ConvertTo-Json $diffAfter -Depth 4)

Write-Host "QA flow completed successfully"
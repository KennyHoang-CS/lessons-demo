CREATE TYPE lesson_status AS ENUM ('DRAFT', 'PUBLISHED');

CREATE TABLE lesson_roots (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE lesson_versions (
    id UUID PRIMARY KEY,
    root_id UUID NOT NULL REFERENCES lesson_roots(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    status lesson_status NOT NULL,
    content_id TEXT NOT NULL,
    parent_version_id UUID REFERENCES lesson_versions(id),
    derived_from_version_id UUID REFERENCES lesson_versions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    UNIQUE (root_id, version_number)
);

CREATE INDEX idx_versions_root_id ON lesson_versions(root_id);
CREATE INDEX idx_versions_root_id_version_number ON lesson_versions(root_id, version_number);
CREATE INDEX idx_versions_published ON lesson_versions(root_id, status, version_number DESC);
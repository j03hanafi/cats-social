CREATE TABLE IF NOT EXISTS matches
(
    id           bytea        NOT NULL PRIMARY KEY,
    match_cat_id bytea        NOT NULL,
    user_cat_id  bytea        NOT NULL,
    message      VARCHAR(120) NOT NULL,
    created_at   TIMESTAMP    NOT NULL,
    updated_at   TIMESTAMP    NOT NULL,
    deleted_at   TIMESTAMP
);

CREATE INDEX idx_matches_match_user_cat_id ON matches (match_cat_id, user_cat_id);
CREATE INDEX idx_matches_created_at_desc_deleted_at_null ON matches (created_at DESC) WHERE deleted_at IS NULL;
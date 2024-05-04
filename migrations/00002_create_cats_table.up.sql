CREATE TYPE cat_race AS ENUM
    (
        'Persian',
        'Maine Coon',
        'Siamese',
        'Ragdoll',
        'Bengal',
        'Sphynx',
        'British Shorthair',
        'Abyssinian',
        'Scottish Fold',
        'Birman'
        );

CREATE TYPE cat_sex AS ENUM (
    'male',
    'female'
    );

CREATE TABLE IF NOT EXISTS cats
(
    id           bytea        NOT NULL PRIMARY KEY,
    name         VARCHAR(30)  NOT NULL,
    race         cat_race     NOT NULL,
    sex          cat_sex      NOT NULL,
    age_in_month INTEGER      NOT NULL CHECK ( age_in_month >= 1 AND age_in_month <= 120082),
    description  VARCHAR(200) NOT NULL,
    user_id      bytea        NOT NULL,
    has_matched  BOOLEAN      NOT NULL,
    created_at   TIMESTAMP    NOT NULL,
    updated_at   TIMESTAMP    NOT NULL,
    deleted_at   TIMESTAMP
);

CREATE INDEX idx_cats_age_in_month ON cats (age_in_month);
CREATE INDEX idx_cats_user_id ON cats (user_id);
CREATE INDEX idx_cats_created_at_desc_deleted_at_null ON cats (created_at DESC) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS cat_images
(
    id         bytea     NOT NULL PRIMARY KEY,
    image_url  TEXT      NOT NULL,
    cat_id     bytea     NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_cat_id ON cat_images (cat_id);
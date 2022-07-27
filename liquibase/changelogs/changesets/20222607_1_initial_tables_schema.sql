create extension if not exists "uuid-ossp";
--rollback drop extension if exists uuid-ossp;

CREATE OR REPLACE FUNCTION sync_updated_at()
    RETURNS trigger AS $$
BEGIN
    NEW.updated_at := now();
RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TABLE "product" (
    product_id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    product_name varchar(30) not null,
    articles_id varchar(30) ARRAY,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

CREATE TABLE "article" (
    article_id integer PRIMARY KEY,
    stock integer DEFAULT 0 not null,
    article_name varchar(30) not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    CONSTRAINT stock_nonnegative CHECK (stock >= 0)
);

CREATE TABLE "product_article" (
    product_id uuid not null,
    article_id integer not null,
    article_count integer not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    CONSTRAINT unique_product_article UNIQUE (product_id,article_id),
    PRIMARY KEY (product_id,article_id)
);
CREATE INDEX "product_article_product_id" ON "product_article" (product_id);
CREATE INDEX "product_article_article_id" ON "product_article" (article_id);

CREATE TRIGGER
    article_updated_at
    BEFORE UPDATE ON
    article
    FOR EACH ROW EXECUTE PROCEDURE
    sync_updated_at();

CREATE TRIGGER
    product_updated_at
    BEFORE UPDATE ON
    product
    FOR EACH ROW EXECUTE PROCEDURE
    sync_updated_at();


CREATE TABLE IF NOT EXISTS sectionItems(
    sectionID BIGSERIAL PRIMARY KEY ,
    isTest BOOLEAN ,
    title TEXT ,
    itemID BIGINT
);

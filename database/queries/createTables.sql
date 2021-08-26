create table if not exists axie_morphed
(
    hash        varchar(256) not null,
    blockNumber int unsigned not null,
    timestamp   datetime     not null,
    tokenId     int unsigned not null,
    primary key (hash),
    unique index tx_hash_uindex (hash)
);

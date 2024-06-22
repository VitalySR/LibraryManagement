CREATE TABLE IF NOT EXISTS author
(
    Id          serial          primary key,
    FirstName   varchar(255)    not null,
    LastName    varchar(255)    not null,
    Biography   text,
    BirthDate   date
);

CREATE TABLE IF NOT EXISTS book
(
    Id          serial          primary key,
	Title       varchar(255)    not null,
	Author_Id   int             REFERENCES author (Id) ON DELETE SET NULL,
    Year        int,
	ISBN        varchar(50)
);
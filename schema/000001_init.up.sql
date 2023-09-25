create schema if not exists praktikum;

create table praktikum.metrics (
    id varchar(50) primary key not null unique,
    m_type varchar(10) not null,
    m_delta bigint,
    m_value double precision
);
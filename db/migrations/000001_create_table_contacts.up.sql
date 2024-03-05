create type contact_status as enum ('new', 'read', 'responded');
create table if not exists contact_requests
(
    id          uuid           default gen_random_uuid(),
    first_name  varchar(50)  not null,
    last_name   varchar(50)  not null,
    email       varchar(100) not null,
    description text         not null,
    status      contact_status default 'new',
    created_at  timestamp      default current_timestamp,
    primary key (id)
);

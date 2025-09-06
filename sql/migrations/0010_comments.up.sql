create table Comments(
    id int8 generated always as identity primary key,
    user_id int8 not null,
    product_id int8 not null,
    comment_text text not null,
    answer_on int8 default null,
    rating int2 check (rating >= 0 and rating <= 100),
    created_at timestamptz not null default now(),
    updated_at timestamptz default null,
    media_ids int8[] not null default '{}',
    is_deleted boolean not null default false
);

create type ReactionType as enum ('like', 'dislike');

create table CommentReactions(
    user_id int8 not null,
    comment_id int8 not null,
    ty ReactionType not null,
    added_at timestamptz not null default now(),
    primary key (user_id, comment_id)
);

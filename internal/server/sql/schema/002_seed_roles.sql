-- +goose Up
insert into roles (name, description)
values
    ('admin', 'Full access to the system'),
    ('tutors', 'Create teaching materials'),
    ('students', 'learners who are taught by various tutors')
on conflict (name) do nothing;

-- +goose Down
delete from roles;

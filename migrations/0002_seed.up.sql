-- users
/* for all users password is password */
INSERT INTO users (id, email, name, password, created_at, updated_at) VALUES
(1, 'user1@example.com', 'User 1', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(2, 'user2@example.com', 'User 2', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(3, 'user3@example.com', 'User 3', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(4, 'user4@example.com', 'User 4', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(5, 'user5@example.com', 'User 5', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(6, 'user6@example.com', 'User 6', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(7, 'user7@example.com', 'User 7', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(8, 'user8@example.com', 'User 8', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(9, 'user9@example.com', 'User 9', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00'),
(10, 'user10@example.com', 'User 10', '$2a$10$51yy0bz8GLRk36YwE3ojeucnYYRptY71AvsMtKmwsEPGFectinKoa', '2026-07-09 09:00:00', '2026-07-09 09:00:00');

-- teams
INSERT INTO teams (id, name, owner_id, created_at) VALUES
(1, 'Team 1', 1, '2026-07-09 09:10:00'),
(2, 'Team 2', 2, '2026-07-09 09:10:00'),
(3, 'Team 3', 3, '2026-07-09 09:10:00'),
(4, 'Team 4', 4, '2026-07-09 09:10:00'),
(5, 'Team 5', 5, '2026-07-09 09:10:00'),
(6, 'Team 6', 6, '2026-07-09 09:10:00'),
(7, 'Team 7', 7, '2026-07-09 09:10:00'),
(8, 'Team 8', 8, '2026-07-09 09:10:00'),
(9, 'Team 9', 9, '2026-07-09 09:10:00'),
(10, 'Team 10', 10, '2026-07-09 09:10:00');

-- team_members
INSERT INTO team_members (team_id, user_id, role, joined_at) VALUES
(1, 1, 'OWNER', '2026-07-09 09:20:00'),
(1, 2, 'MEMBER', '2026-07-09 09:20:00'),
(1, 5, 'ADMIN', '2026-07-09 09:20:00'),
(1, 10, 'MEMBER', '2026-07-09 09:20:00'),
(2, 2, 'OWNER', '2026-07-09 09:20:00'),
(2, 3, 'MEMBER', '2026-07-09 09:20:00'),
(2, 4, 'ADMIN', '2026-07-09 09:20:00'),
(2, 8, 'MEMBER', '2026-07-09 09:20:00'),
(3, 3, 'OWNER', '2026-07-09 09:20:00'),
(3, 1, 'ADMIN', '2026-07-09 09:20:00'),
(3, 6, 'MEMBER', '2026-07-09 09:20:00'),
(3, 9, 'MEMBER', '2026-07-09 09:20:00'),
(4, 4, 'OWNER', '2026-07-09 09:20:00'),
(4, 2, 'MEMBER', '2026-07-09 09:20:00'),
(4, 7, 'ADMIN', '2026-07-09 09:20:00'),
(4, 10, 'MEMBER', '2026-07-09 09:20:00'),
(5, 5, 'OWNER', '2026-07-09 09:20:00'),
(5, 1, 'MEMBER', '2026-07-09 09:20:00'),
(5, 8, 'ADMIN', '2026-07-09 09:20:00'),
(5, 9, 'MEMBER', '2026-07-09 09:20:00');

-- tasks
INSERT INTO tasks (id, team_id, created_by, assignee_id, title, description, status, priority, due_date, created_at, updated_at) VALUES
(1, 1, 1, 2, 'Task 1', 'Description for task 1', 'OPEN', 'MEDIUM', '2026-07-16', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(2, 1, 1, 5, 'Task 2', 'Description for task 2', 'IN_PROGRESS', 'HIGH', '2026-07-17', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(3, 2, 2, 3, 'Task 3', 'Description for task 3', 'DONE', 'LOW', '2026-07-18', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(4, 2, 2, 4, 'Task 4', 'Description for task 4', 'OPEN', 'URGENT', '2026-07-19', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(5, 3, 3, 6, 'Task 5', 'Description for task 5', 'CANCELLED', 'MEDIUM', '2026-07-20', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(6, 3, 3, 1, 'Task 6', 'Description for task 6', 'DONE', 'HIGH', '2026-07-21', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(7, 4, 4, 7, 'Task 7', 'Description for task 7', 'IN_PROGRESS', 'LOW', '2026-07-22', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(8, 4, 4, 10, 'Task 8', 'Description for task 8', 'OPEN', 'MEDIUM', '2026-07-23', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(9, 5, 5, 8, 'Task 9', 'Description for task 9', 'DONE', 'URGENT', '2026-07-24', '2026-07-09 09:30:00', '2026-07-09 09:30:00'),
(10, 5, 5, 9, 'Task 10', 'Description for task 10', 'OPEN', 'MEDIUM', '2026-07-25', '2026-07-09 09:30:00', '2026-07-09 09:30:00');

-- task_history
INSERT INTO task_history (task_id, user_id, action, created_at) VALUES
(1, 1,  'created task', '2026-07-01 08:00:00'),
(1, 2,  'updated status to in_progress', '2026-07-01 09:15:00'),
(1, 3,  'assigned task to user 3', '2026-07-01 10:30:00'),

(2, 2,  'created task', '2026-07-01 08:10:00'),
(2, 4,  'updated title', '2026-07-01 11:05:00'),
(2, 5,  'updated priority to high', '2026-07-01 12:20:00'),

(3, 3,  'created task', '2026-07-01 08:20:00'),
(3, 5,  'updated status to done', '2026-07-01 09:40:00'),
(3, 6,  'added due date', '2026-07-01 13:00:00'),

(4, 4,  'created task', '2026-07-01 08:30:00'),
(4, 6,  'updated description', '2026-07-01 10:10:00'),
(4, 7,  'changed assignee to user 7', '2026-07-01 14:15:00'),

(5, 5,  'created task', '2026-07-01 08:40:00'),
(5, 7,  'updated status to in_review', '2026-07-01 11:25:00'),
(5, 8,  'updated due date', '2026-07-01 15:30:00'),

(6, 6,  'created task', '2026-07-01 08:50:00'),
(6, 8,  'updated title', '2026-07-01 12:35:00'),
(6, 9,  'closed task', '2026-07-01 16:45:00'),

(7, 7,  'created task', '2026-07-01 09:00:00'),
(7, 9,  'updated priority to low', '2026-07-01 13:45:00'),
(7, 10, 'reopened task', '2026-07-01 17:10:00'),

(8, 8,  'created task', '2026-07-01 09:10:00'),
(8, 10, 'updated description', '2026-07-01 14:55:00'),
(8, 1,  'assigned task to user 1', '2026-07-01 18:05:00'),

(9, 9,  'created task', '2026-07-01 09:20:00'),
(9, 1,  'updated status to blocked', '2026-07-01 15:20:00'),
(9, 2,  'changed title', '2026-07-01 19:00:00'),

(10, 10, 'created task', '2026-07-01 09:30:00'),
(10, 2, 'updated status to done', '2026-07-01 16:10:00'),
(10, 3, 'archived task', '2026-07-01 20:15:00');


-- task_comments
INSERT INTO task_comments (task_id, user_id, body, created_at) VALUES
(1, 1,  'Initial comment on task 1', '2026-07-01 09:10:00'),
(1, 2,  'Added some clarification for task 1', '2026-07-01 10:15:00'),
(1, 3,  'Task 1 is ready for review', '2026-07-01 11:20:00'),

(2, 2,  'Started working on task 2', '2026-07-01 09:30:00'),
(2, 4,  'Blocked by missing requirements', '2026-07-01 12:05:00'),

(3, 3,  'Task 3 depends on task 2', '2026-07-01 08:45:00'),
(3, 5,  'Waiting for feedback', '2026-07-01 13:10:00'),

(4, 4,  'Investigating task 4 issue', '2026-07-01 09:55:00'),
(4, 6,  'Found the root cause', '2026-07-01 14:25:00'),

(5, 5,  'Task 5 assigned and in progress', '2026-07-01 10:05:00'),
(5, 7,  'Need updated design assets', '2026-07-01 15:00:00'),

(6, 6,  'Task 6 has been reviewed', '2026-07-01 11:40:00'),
(6, 8,  'Minor fixes required', '2026-07-01 16:20:00'),

(7, 7,  'Working on task 7 implementation', '2026-07-01 12:00:00'),
(7, 9,  'Please check edge cases', '2026-07-01 17:15:00'),

(8, 8,  'Task 8 moved to testing', '2026-07-01 13:35:00'),
(8, 10, 'Testing results look good', '2026-07-01 18:10:00'),

(9, 9,  'Task 9 is awaiting approval', '2026-07-01 14:50:00'),
(9, 1,  'Approved from my side', '2026-07-01 19:05:00'),

(10, 10, 'Final comment for task 10', '2026-07-01 15:30:00'),
(10, 2,  'Ready to close task 10', '2026-07-01 20:00:00');
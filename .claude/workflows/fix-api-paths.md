# Fix API Paths and Supplement Missing API Docs

Parallel task to fix mismatched API paths in frontend code and supplement missing API documentation.

## Phase 1: Fix Frontend API Paths

### Task 1: Fix user.ts - `/system/` → `/admin/`
- `/system/users` → `/admin/users`
- `/system/users/${id}` → `/admin/users/${id}`
- `/system/users/${id}/status` → `/admin/users/${id}/status`
- `/system/users/${id}/reset-password` → `/admin/users/${id}/reset-password`
- `/system/roles` → `/admin/roles`
- `/system/roles/${id}` → `/admin/roles/${id}`
- `/system/roles/${id}/permissions` → `/admin/roles/${id}/permissions`
- `/system/permissions/tree` → `/admin/permissions/tree`

### Task 2: Fix system.ts - Log paths
- `/operation-logs` → `/admin/operation-logs`
- `/login-logs` → `/admin/login-logs`
- `/system-logs` → `/admin/system-logs`

## Phase 2: Supplement Missing API Docs

### Module 1: User Management (用户管理)
Add to `01-用户认证与授权.md` or create `03-用户管理.md`:
- GET `/admin/users` - User list (paginated)
- GET `/admin/users/{id}` - User detail
- PUT `/admin/users/{id}` - Update user
- DELETE `/admin/users/{id}` - Delete user
- PUT `/admin/users/{id}/status` - Toggle user status
- POST `/admin/users/{id}/reset-password` - Reset password
- GET `/admin/roles` - Role list
- GET `/admin/roles/{id}` - Role detail
- POST `/admin/roles` - Create role
- PUT `/admin/roles/{id}` - Update role
- DELETE `/admin/roles/{id}` - Delete role
- PUT `/admin/roles/{id}/permissions` - Update role permissions
- GET `/admin/permissions/tree` - Permission tree

### Module 2: Question Management (题库管理)
Add to `02-题库管理.md`:
- PUT `/admin/questions/{id}/status` - Update question status
- POST `/admin/questions/batch-import` - Batch import
- POST `/admin/questions/batch-export` - Batch export
- GET `/admin/question-bank/options` - Question bank options

### Module 3: Exam Management (考试管理)
Add to `04-考试管理.md`:
- GET `/admin/exam-papers/{id}/questions` - Get paper questions
- POST `/admin/exam-papers/{id}/questions` - Add question to paper
- DELETE `/admin/exam-papers/{id}/questions/{qid}` - Remove question
- POST `/admin/exams/{id}/publish` - Publish exam
- GET `/admin/exams/{id}/monitor` - Exam monitor data
- GET `/admin/exam-records` - Exam records list
- GET `/admin/exam-records/{id}` - Record detail
- POST `/admin/exam-records/{id}/review` - Review record

### Module 4: Class Management (班级管理)
Add to `05-用户与班级.md`:
- GET `/admin/classes/{id}/students` - Class students
- POST `/admin/classes/{id}/students` - Add student
- DELETE `/admin/classes/{id}/students/{uid}` - Remove student
- GET `/admin/class-creator-applications` - Creator applications
- PUT `/admin/class-creator-applications/{id}` - Review application

### Module 5: System Management (系统管理)
Create `06-系统管理.md`:
- GET/POST/PUT/DELETE `/admin/dict-types` - Dictionary types
- GET/POST/PUT/DELETE `/admin/dict-data` - Dictionary data
- GET/POST/PUT/DELETE `/admin/menus` - Menus
- GET/POST/PUT/DELETE `/admin/notices` - Notices
- GET/POST/PUT/DELETE `/admin/configs` - System configs
- GET/POST/PUT/DELETE `/admin/tags` - Tags
- GET/PUT `/admin/real-name-audit` - Real name audit
- GET/PUT `/admin/role-applications` - Role applications
- POST `/admin/cache/clear` - Clear cache
- GET/PUT `/admin/feature-flags` - Feature flags
- GET/PUT `/admin/storage-config` - Storage config

### Module 6: Statistics (数据统计)
Add to `10-数据统计与分析.md`:
- GET `/admin/statistics/exam-stats` - Exam statistics
- GET `/admin/statistics/question-stats` - Question statistics
- GET `/admin/statistics/realtime` - Realtime data

### Module 7: Content Management (内容管理)
Create `07-内容管理.md`:
- GET/DELETE `/admin/comments` - Comments
- GET `/admin/reports` - Reports
- GET/POST/DELETE `/admin/comment-blacklist` - Blacklist
- GET/POST/PUT/DELETE `/admin/moderation-rules` - Rules
- POST `/admin/comments/batch-delete` - Batch delete
- PUT `/admin/comments/{id}/audit` - Audit comment
- POST `/admin/comment-blacklist/batch-add` - Batch add blacklist

### Module 8: File Management (文件管理)
Add to `11-文件管理.md`:
- POST `/admin/files/batch-delete` - Batch delete files

### Module 9: Audit Management (审计管理)
Create `08-审计管理.md`:
- GET `/admin/audit/reports` - Audit reports
- GET `/admin/audit/comments` - Audit comments
- GET/POST/DELETE `/admin/audit/comment-blacklist` - Audit blacklist
- GET/POST/PUT/DELETE `/admin/audit/rules` - Audit rules

## Execution

Run these tasks in parallel:
1. Fix user.ts paths
2. Fix system.ts paths
3. Create/supplement API docs (can be parallelized by module)

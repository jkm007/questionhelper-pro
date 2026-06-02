export const meta = {
  name: 'fix-api-paths',
  description: 'Fix mismatched API paths in frontend code and supplement missing API documentation',
  phases: [
    { title: 'Fix Frontend Paths', detail: 'Fix user.ts and system.ts API paths' },
    { title: 'Supplement API Docs', detail: 'Add missing API documentation by module' },
  ],
}

phase('Fix Frontend Paths')

// Fix user.ts - /system/ → /admin/
await agent('Fix API paths in admin/src/api/user.ts: change all /system/ prefixes to /admin/. Examples: /system/users → /admin/users, /system/roles → /admin/roles, /system/permissions/tree → /admin/permissions/tree', {
  label: 'fix:user.ts',
  phase: 'Fix Frontend Paths',
})

// Fix system.ts - log paths
await agent('Fix API paths in admin/src/api/system.ts: change /operation-logs → /admin/operation-logs, /login-logs → /admin/login-logs, /system-logs → /admin/system-logs', {
  label: 'fix:system.ts',
  phase: 'Fix Frontend Paths',
})

phase('Supplement API Docs')

// Create comprehensive API docs in parallel
await parallel([
  () => agent('Create or update docs/API接口文档/03-用户管理.md with complete API documentation for user management endpoints: GET/POST/PUT/DELETE /admin/users, user status, reset password, roles CRUD, role permissions, permissions tree. Follow the same format as existing API docs.', {
    label: 'doc:user-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Update docs/API接口文档/02-题库管理.md to add missing endpoints: PUT /admin/questions/{id}/status, POST /admin/questions/batch-import, POST /admin/questions/batch-export, GET /admin/question-bank/options', {
    label: 'doc:question-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Update docs/API接口文档/04-考试管理.md to add missing endpoints: paper questions CRUD, exam publish, exam monitor, exam records, record review', {
    label: 'doc:exam-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Update docs/API接口文档/05-用户与班级.md to add missing endpoints: class students CRUD, class creator applications', {
    label: 'doc:class-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Create docs/API接口文档/06-系统管理.md with complete API documentation for: dict-types, dict-data, menus, notices, configs, tags, real-name-audit, role-applications, cache, feature-flags, storage-config', {
    label: 'doc:system-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Update docs/API接口文档/10-数据统计与分析.md to add missing endpoints: exam-stats, question-stats, realtime', {
    label: 'doc:statistics',
    phase: 'Supplement API Docs',
  }),
  () => agent('Create docs/API接口文档/07-内容管理.md with complete API documentation for: comments, reports, comment-blacklist, moderation-rules, batch operations, audit', {
    label: 'doc:content-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Update docs/API接口文档/11-文件管理.md to add missing endpoint: POST /admin/files/batch-delete', {
    label: 'doc:file-management',
    phase: 'Supplement API Docs',
  }),
  () => agent('Create docs/API接口文档/08-审计管理.md with complete API documentation for: audit reports, audit comments, audit blacklist, audit rules', {
    label: 'doc:audit-management',
    phase: 'Supplement API Docs',
  }),
])

log('All API path fixes and documentation supplements completed!')

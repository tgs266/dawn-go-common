package audit

type AuditAction string

const CREATE AuditAction = "create"
const UPDATE AuditAction = "update"
const DELETE AuditAction = "delete"

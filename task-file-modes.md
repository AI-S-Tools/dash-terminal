# Task File Operation Modes

Every dppm task MUST declare its file operation mode to prevent AI from destroying existing work:

## 🆕 MODE: CREATE NEW
```yaml
**✋ FILE OPERATION MODE: CREATE NEW**
- This task CREATES NEW files only
- If files exist, STOP - task is likely done
- Use 'Write' tool for new files
- Never use 'Edit' tool
```

## ✏️ MODE: REFACTOR EXISTING
```yaml
**✋ FILE OPERATION MODE: REFACTOR EXISTING**
- This task MODIFIES existing files
- Files MUST exist or task cannot proceed
- Use 'Read' first, then 'Edit' tool
- Never use 'Write' tool (would overwrite)
```

## 🔄 MODE: EXTEND EXISTING
```yaml
**✋ FILE OPERATION MODE: EXTEND EXISTING**
- This task ADDS to existing files
- Files MUST exist first
- Use 'Read' to check current state
- Use 'Edit' to add new sections
- Preserve ALL existing code
```

## 🗑️ MODE: REPLACE EXISTING
```yaml
**✋ FILE OPERATION MODE: REPLACE EXISTING**
- This task REPLACES entire files
- Backup originals first if needed
- Use 'Read' to verify target
- Use 'Write' to replace completely
```

## Example Task Header:

```yaml
description: |
  ## T2.1: WebSocket Handler Implementation

  **📂 FILE STATUS CHECK - RUN THIS FIRST:**
  ```bash
  ls -la internal/websocket/handler.go || { echo "ERROR: File missing - run T1.1 first"; exit 1; }
  grep "not implemented" internal/websocket/handler.go || { echo "ERROR: Already implemented"; exit 1; }
  ```

  **✏️ FILE OPERATION MODE: REFACTOR EXISTING**
  - This task MODIFIES existing stub files from T1.1
  - Files MUST exist (created by project-setup task)
  - Use 'Read' first to check current stub implementation
  - Use 'Edit' to replace "not implemented" with real code
  - PRESERVE test files and interfaces

  **🔒 PROTECTED FILES - DO NOT MODIFY:**
  - main.go (managed by different task)
  - go.mod (only go get for new deps)
  - Any *_test.go files (separate test task)
```

## Critical Rules:

1. **ALWAYS run FILE STATUS CHECK first**
2. **STOP if mode doesn't match file state**
3. **NEVER Write to existing files (destroys work)**
4. **NEVER Edit non-existent files (causes errors)**
5. **PROTECT files owned by other tasks**
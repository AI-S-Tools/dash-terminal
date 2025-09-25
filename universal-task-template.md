# Universal Task Template - Anti Scope Creep

EVERY SINGLE TASK must have ALL these sections:

## Task Header Format:
```yaml
description: |
  ## TX.Y: Task Name - SPECIFIC SCOPE ONLY

  **📋 SUMMARY FROM PREVIOUS TASKS (What AI before you built):**
  - ✅ **T1.1**: [What was built]
  - ✅ **T1.2**: [What was built]
  - ✅ **TX.Y-1**: [Previous task results]
  - 📍 **You are here**: TX.Y - [Your specific scope]

  **📂 FILE STATUS CHECK - RUN THIS FIRST:**
  ```bash
  # CHECK PREREQUISITES FROM PREVIOUS TASKS
  [Commands to verify dependencies]

  # CHECK IF THIS TASK IS ALREADY DONE
  [Commands to detect if work already exists]
  ```

  **FILE OPERATION MODE:** [CREATE NEW | REFACTOR EXISTING | EXTEND EXISTING | REPLACE EXISTING]
  - [Specific instructions on Read/Edit/Write tools]

  **🔒 PROTECTED FILES - DO NOT MODIFY:**
  - [List of files this task must NOT touch]

  **🤖 AI EXACT IMPLEMENTATION:**
  [Step by step copy-paste ready commands]

  **🚫 FORBIDDEN - DO NOT ADD:**
  - ❌ [15+ specific things NOT to implement]

  **✅ VERIFICATION CHECKLIST:**
  ```bash
  # [5+ commands to verify success]
  ```

  **🎯 SUCCESS CRITERIA:**
  - [Specific measurable goals]

  **📝 HANDOFF SUMMARY - UPDATE WHEN COMPLETE:**
  ```
  ## 🤖 TX.Y COMPLETION SUMMARY
  **COMPLETED BY:** [AI Model Name]
  **DATE:** [Date]
  **WHAT WAS BUILT:** [Specific deliverables]
  **SCOPE DISCIPLINE MAINTAINED:** [What was NOT built]
  **BUILD STATUS:** [Verification results]
  **NEXT AI TASK:** [What T(X.Y+1) can now do]
  ```
```

## What this prevents:

1. **📋 SUMMARY** - AI knows what others built before
2. **📂 FILE STATUS** - Prevents overwriting existing work
3. **FILE OPERATION MODE** - Clear Read/Edit/Write instructions
4. **🔒 PROTECTED FILES** - Can't break other tasks' work
5. **🤖 EXACT IMPLEMENTATION** - No guessing, copy-paste ready
6. **🚫 FORBIDDEN** - Explicit scope boundaries
7. **✅ VERIFICATION** - Must prove it works
8. **📝 HANDOFF** - Documents what was actually built for next AI

Every task becomes:
- ✅ **Self-contained** - AI reads only this one file
- ✅ **Safe** - Cannot destroy existing work
- ✅ **Scoped** - Cannot add features
- ✅ **Verifiable** - Must pass tests
- ✅ **Traceable** - Documents what was built
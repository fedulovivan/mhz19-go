**Prioritized TODO List:**

### **P0 (Critical/Blockers)**
- **Fix bug major**: No cascade deletion of messages after deletion of FK linked device.  
  *(Prevents data loss; urgent for data integrity.)*
- **Fix bug major**: GIT_REV not saved to image when built in Compose.  
  *(Critical for version tracking and deployments.)*
- **Fix bug major**: Reply to `/door` command sent to the wrong bot.  
  *(User-facing issue impacting functionality.)*
- **Fix bug minor**: Device-pinger service stopped working with MQTT "not connected" error.  
  *(Service downtime impacts core functionality.)*

---

### **P1 (High Priority)**
- **Architecture decision**: Preserve custom DB changes after `migrate-reset`.  
  *(Foundational for data persistence during updates.)*
- **Architecture decision**: Execute migrations within running containers or via Postgres.  
  *(Essential for deployment workflows.)*
- **Improve**: Deploy Grafana with dashboards.  
  *(Critical for monitoring and observability.)*
- **Feature**: Emit "device is gone/offline/back" events to a separate topic.  
  *(Improves event handling for buried devices.)*
- **Improve**: Auto DB backup before migrations.  
  *(Prevents data loss during migrations.)*
- **Fix bug minor**: Rebuilding frontend unnecessarily with `mhz19-go`.  
  *(Improves development efficiency.)*

---

### **P2 (Medium Priority)**
- **Improve API**: Toggle rule on/off, update rule strategy, and URL consistency.  
  *(Enhances API usability.)*
- **Improve DB**: Add `created_at`/`updated_at`, `rooms` table, and `comments` column.  
  *(Schema improvements for future scalability.)*
- **Fix bug minor**: Skip duplicates in `make seed`.  
  *(Data consistency fix.)*
- **Improve**: Avoid hardcoded IPs in config (e.g., MQTT_HOST).  
  *(Eases environment migration.)*
- **Investigate**: Slow `getAll` queries and Prometheus benchmark discrepancies.  
  *(Performance optimization.)*
- **Improve UTs**: Create unit tests for `buried_devices/provider.go`.  
  *(Improves code reliability.)*

---

### **P3 (Lower Priority)**
- **Improve**: Rename `DEVICE_CLASS_SONOFF_DIY_PLUG` to a clearer name.  
  *(Minor naming fix.)*
- **Investigate**: Code audit for useless abstractions, "Functional options" pattern.  
  *(Code quality improvements.)*
- **Improve**: Use human-readable timestamps (e.g., "since 8:52" instead of "10m").  
  *(User experience tweak.)*
- **Improve**: Split REST API and engine into separate services.  
  *(Architectural refactor; requires planning.)*
- **Investigate**: Makefile `PHONY` targets and range over channel behavior.  
  *(Minor technical debt.)*

---

**Rationale**:  
- **P0** addresses critical bugs causing data loss, service downtime, or user-facing errors.  
- **P1** focuses on foundational architecture, monitoring, and high-impact features.  
- **P2** includes usability improvements, performance investigations, and technical debt.  
- **P3** covers minor refactors, code quality, and non-urgent enhancements.  

Start with **P0**, then proceed to **P1** to stabilize the system. **P2/P3** can be tackled incrementally.
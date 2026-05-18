# Disaster Recovery

Anleitung zur Wiederherstellung der Datenbank nach einem Ausfall.

## Backup-Strategie

### Automatische Backups (postgres-backup-s3 + RustFS)

Der Backup-Stack besteht aus:
- **postgres-backup-s3**: Automatische PostgreSQL-Backups mit Cron
- **RustFS**: S3-kompatibler Object Storage (self-hosted, Rust-basiert)
- **rustfs-init**: Erstellt automatisch den Backup-Bucket

**Keine manuelle Konfiguration nötig!**

```bash
# Backup-Stack starten
just backup-up

# Backup sofort erstellen
just backup-now

# Backups anzeigen
just backup-list

# RustFS Console öffnen (optional)
# http://localhost:9001/rustfs/console/
# Login: rustfsadmin / rustfsadmin
```

### Konfiguration (Environment Variables)

| Variable | Default | Beschreibung |
|----------|---------|--------------|
| `BACKUP_SCHEDULE` | @daily | Cron-Schedule für Backups |
| `BACKUP_KEEP_DAYS` | 7 | Aufbewahrungsdauer in Tagen |
| `S3_ACCESS_KEY` | rustfsadmin | RustFS Access Key |
| `S3_SECRET_KEY` | rustfsadmin | RustFS Secret Key |

### Manuelle Backups

```bash
# Sofort-Backup erstellen
just backup-now

# Oder direkt mit pg_dump
docker compose -f deploy/compose/docker-compose.dev.yml exec db pg_dump -U postgres -d nextgopg > backup.sql
```

## Disaster Recovery Prozess

### Szenario 1: Datenbank-Daten verloren (Volume gelöscht)

```bash
# 1. Neue Datenbank starten
just db-up

# 2. Migrations ausführen (Schema erstellen)
just migrate-up

# 3. Daten aus Backup wiederherstellen
just backup-restore
```

### Szenario 2: Kompletter Server-Ausfall

```bash
# 1. Repository klonen
git clone https://github.com/atilladeniz/next-go-pg.git
cd next-go-pg

# 2. Dependencies installieren
just install

# 3. Datenbank starten
just db-up

# 4. Migrations ausführen
just migrate-up

# 5. Backup-Stack starten
just backup-up

# 6. Daten wiederherstellen (vom letzten Backup)
just backup-restore

# 7. Anwendung starten
just dev
```

### Szenario 3: Production (Docker)

```bash
# 1. Docker Container starten (führt Migrations automatisch aus)
docker compose up -d

# 2. Backup wiederherstellen
docker compose exec db psql -U postgres -d nextgopg < backup.sql
```

## Backup-Speicherorte

### Development (Lokal)

Backups werden in RustFS gespeichert:
- S3 Bucket: `backups/postgres/`
- Dateiformat: `nextgopg_YYYY-MM-DDTHH:MM:SS.dump`

```bash
# Backups anzeigen
just backup-list
```

### Production (S3/RustFS)

Für Production kann RustFS auf einem separaten Server laufen oder durch AWS S3/andere S3-kompatible Services ersetzt werden.

Ändere die Environment Variables:
```bash
S3_ENDPOINT=https://s3.example.com
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
```

## Backup-Retention

| Umgebung | Häufigkeit | Aufbewahrung |
|----------|------------|--------------|
| Dev | Täglich (@daily) | 7 Tage |
| Staging | Täglich | 14 Tage |
| Production | Stündlich (@hourly) | 30 Tage |

## Wichtige Dateien

| Datei | Beschreibung |
|-------|--------------|
| `deploy/compose/docker-compose.backup.yml` | postgres-backup-s3 + RustFS Service Definition |
| `backend/migrations/*.sql` | Schema-Definitionen |
| `.env` | Datenbank-Credentials |

## Makefile Commands

| Command | Beschreibung |
|---------|--------------|
| `just backup-up` | Startet automatisches Backup-System |
| `just backup-down` | Stoppt Backup-Stack |
| `just backup-now` | Erstellt sofort ein Backup |
| `just backup-list` | Zeigt alle Backups in S3 |
| `just backup-restore` | Stellt vom letzten Backup wieder her |

## Testen der Recovery

**Regelmäßig testen!**

```bash
# 1. Backup erstellen
just backup-now

# 2. Backup prüfen
just backup-list

# 3. Datenbank löschen
just db-reset

# 4. Recovery durchführen
just db-up
just migrate-up
just backup-up
just backup-restore

# 5. Anwendung testen
just dev
```

## Monitoring

- **just backup-list**: Zeigt alle vorhandenen Backups
- **RustFS Console**: http://localhost:9001/rustfs/console/
- Container-Logs: `docker logs next-go-pg-pg-backup-1`

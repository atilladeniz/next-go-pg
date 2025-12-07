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
make backup-up

# Backup sofort erstellen
make backup-now

# Backups anzeigen
make backup-list

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
make backup-now

# Oder direkt mit pg_dump
docker compose -f docker-compose.dev.yml exec db pg_dump -U postgres -d nextgopg > backup.sql
```

## Disaster Recovery Prozess

### Szenario 1: Datenbank-Daten verloren (Volume gelöscht)

```bash
# 1. Neue Datenbank starten
make db-up

# 2. Migrations ausführen (Schema erstellen)
make migrate-up

# 3. Daten aus Backup wiederherstellen
make backup-restore
```

### Szenario 2: Kompletter Server-Ausfall

```bash
# 1. Repository klonen
git clone https://github.com/atilladeniz/next-go-pg.git
cd next-go-pg

# 2. Dependencies installieren
make install

# 3. Datenbank starten
make db-up

# 4. Migrations ausführen
make migrate-up

# 5. Backup-Stack starten
make backup-up

# 6. Daten wiederherstellen (vom letzten Backup)
make backup-restore

# 7. Anwendung starten
make dev
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
make backup-list
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
| `docker-compose.backup.yml` | postgres-backup-s3 + RustFS Service Definition |
| `backend/migrations/*.sql` | Schema-Definitionen |
| `.env` | Datenbank-Credentials |

## Makefile Commands

| Command | Beschreibung |
|---------|--------------|
| `make backup-up` | Startet automatisches Backup-System |
| `make backup-down` | Stoppt Backup-Stack |
| `make backup-now` | Erstellt sofort ein Backup |
| `make backup-list` | Zeigt alle Backups in S3 |
| `make backup-restore` | Stellt vom letzten Backup wieder her |

## Testen der Recovery

**Regelmäßig testen!**

```bash
# 1. Backup erstellen
make backup-now

# 2. Backup prüfen
make backup-list

# 3. Datenbank löschen
make db-reset

# 4. Recovery durchführen
make db-up
make migrate-up
make backup-up
make backup-restore

# 5. Anwendung testen
make dev
```

## Monitoring

- **make backup-list**: Zeigt alle vorhandenen Backups
- **RustFS Console**: http://localhost:9001/rustfs/console/
- Container-Logs: `docker logs next-go-pg-pg-backup-1`

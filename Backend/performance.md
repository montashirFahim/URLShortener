# Performance Monitoring - URL Shortener Backend (Normalized)

Monitoring a normalized system tracks how multiple storage layers interact during a single redirect flow.

## 1. Key Metrics to Monitor

### 1.1 Redirection Latency
- **Target:** < 5ms for cached URLs.
- **Monitoring:** Track `GetLongURL` in the service layer.
- **Cache-Aside:** Redirections use Redis to map `code` -> `long_url` directly, avoiding DB JOINs for public traffic.

### 1.2 Cache Performance
- **Hit Ratio:** Monitor Redis `GET` vs `SET` operations. High hit ratios ensure low redirection latency.
- **Lean Storage:** Cache stores simple strings instead of full objects to minimize memory footprint.

### 1.3 Asynchronous Analytics Pool
- **Dual Write Path:** The worker pool now performs two operations per click:
    1. `INSERT INTO clicks` for granular historical data.
    2. `UPSERT` on `url_stats` for aggregate counts.
- **Queue Monitoring:** Monitor the `analyticsChan` buffer depth in `URLService`.

## 2. Database Performance (PostgreSQL)
- **JOIN Performance:** Monitor the `GetUserURLs` operation which joins `short_urls` and `long_urls`.
- **Query Optimization:** Ensure indexes on `long_urls.url` and `short_urls.code`.
- **Connection Pool:** Monitor `Postgres-MaxOpenConn` as async writes compete with sync reads.

## 3. Practical Commands

**Monitor Click Logs in Real-Time:**
```bash
watch 'docker exec -i pg-container psql -U postgres -d url_shortener -c "SELECT count(*) FROM clicks;"'
```

**Monitor Redis traffic:**
```bash
docker exec -it redis-stack redis-cli monitor
```

**Monitor Table Sizes:**
```bash
docker exec -it pg-container psql -U postgres -d url_shortener -c "SELECT relname, n_live_tup FROM pg_stat_user_tables;"
```

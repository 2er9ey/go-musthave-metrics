-- Таблица metric для хранения значений метрик
CREATE TABLE metrics (
  id           SERIAL PRIMARY KEY,
  metric_id    VARCHAR(255) NOT NULL,
  metric_type  VARCHAR(255) NOT NULL,
  metric_delta bigint,
  metric_value double precision,
  metric_hash  VARCHAR(255)
);

-- Базовый индекс для поиска по названию
CREATE UNIQUE INDEX idx_metrics_type_id ON metrics(metric_type, metric_id);
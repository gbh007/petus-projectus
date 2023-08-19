# Petus-projectus

История о том как сделать нажатие кнопки через микросервисы XD

## TODO

1. Рефакторинг работы с ошибками
2. Рефакторинг логов пользователя (после реализации внутреннего функционала):
   - Необходимо добавить новые поля в сигнатуру через кафку
   - Необходимо добавить новые поля в БД

## Примеры запросов для prometheus

Среднее время бизнес обработки задачи

```plain
  sum by (instance) (rate(petus_projectus_worker_business_handle_time_sum[1m]))
/
  sum by (instance) (rate(petus_projectus_worker_business_handle_time_count[1m]))
```

Размер очереди задач

```plain
sum(petus_projectus_handler_handle_time_count) - sum(petus_projectus_worker_handle_time_count)
```

Среднее время выполнения запроса

```plain
  sum by (instance) (rate(petus_projectus_pb_request_duration_seconds_sum[1m]))
/
  sum by (instance) (rate(petus_projectus_pb_request_duration_seconds_count[1m]))
```

Количество запросов в минуту

```plain
sum by (instance) (rate(petus_projectus_pb_request_duration_seconds_count[1m]))
```

# HTTP Helpers

Модуль реализующий небольшой протокол поверх HTTP. 

Содержит подмодуль HTTP Server с ендпоинтом для метрик, health check, rate limit, timeout limit, request ID, логирование. Все можно настраивать через конифиги. 

Также содержит подмодуль HTTP Client реализующий таймауты, retry, max concurrent request count, circuit breaker паттерны.

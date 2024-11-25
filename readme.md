

## Запуск проекта!
Для запуска проекта достаточно перейти в папку infrastructure cd infrastructure и выполнить команду docker compose up --build по желанию можно добавить и флаг -d.
Сваггер запускается на порту 4444 там можно проверить все ендпоинты (при запуске проекта миграцией добавляется не большое количество данных).

## О ТЗ
Все пункты тз выполнены. Проект написан в формате дева (допущены погрешности такие как - .env-файлы не добавлены в гитигнор и т.д.)

### Моменты для улучшения:
- улучшить валидацию
- унифицировать репозиторий немного отрефакторив
- улучшить миграция.
- закрыть не нужные внешние порты докер композа (они открыты для простоты тестирования)
- имеются некоторые и незначительные дублирования кода 

Весь код в данном проекте было написано мной, не скопировав откуда-то и не посмотрев куда-то (в нем есть интересные утилитарные функции, которых, чуть отрефакторив,  можно использовать почти везде.)
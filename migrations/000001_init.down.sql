DROP TABLE IF EXISTS exchanges;
DROP TABLE IF EXISTS user_actions;

DELETE FROM exchange_rates_history
USING exchange_rates
WHERE exchange_rates_history.exchange_rate_id = exchange_rates.id
  AND exchange_rates.base_currency_code IN (
      'USD', 'RUB', 'EUR', 'GBP', 'CHF', 'PLN', 'CNY', 'JPY',
      'KZT', 'TRY', 'AED', 'GEL', 'UAH', 'CZK', 'CAD', 'AUD'
  )
  AND exchange_rates.quote_currency_code = 'BYN';

DELETE FROM exchange_rates
WHERE base_currency_code IN (
    'USD', 'RUB', 'EUR', 'GBP', 'CHF', 'PLN', 'CNY', 'JPY',
    'KZT', 'TRY', 'AED', 'GEL', 'UAH', 'CZK', 'CAD', 'AUD'
)
  AND quote_currency_code = 'BYN';

DELETE FROM currencies
WHERE code IN (
    'USD', 'BYN', 'RUB', 'EUR', 'GBP', 'CHF', 'PLN', 'CNY',
    'JPY', 'KZT', 'TRY', 'AED', 'GEL', 'UAH', 'CZK', 'CAD', 'AUD'
);

DELETE FROM users
WHERE name IN ('admin', 'operator');

DROP TABLE IF EXISTS exchange_rates_history;
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS currencies;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_role;

DROP EXTENSION IF EXISTS pgcrypto;

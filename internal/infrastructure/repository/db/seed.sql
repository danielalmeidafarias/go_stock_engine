INSERT INTO product_stock_models (id, name, category, current_stock, minimum_stock, average_daily_sales, lead_time_days, unit_cost, criticality_level) VALUES
  (gen_random_uuid(), 'Oil Filter Premium',     'oil',    50,  30, 10, 5,  18.50, 3),
  (gen_random_uuid(), 'Engine Block V8',         'engine', 5,   10, 1,  30, 4500.00, 5),
  (gen_random_uuid(), 'Synthetic Motor Oil 5W30','oil',    200, 80, 25, 3,  32.90, 2),
  (gen_random_uuid(), 'Piston Ring Set',         'engine', 15,  20, 3,  14, 89.00, 4),
  (gen_random_uuid(), 'Oil Drain Plug',          'oil',    300, 50, 12, 2,  4.50, 1),
  (gen_random_uuid(), 'Crankshaft Bearing',      'engine', 0,   8,  2,  21, 120.00, 5),
  (gen_random_uuid(), 'Oil Pressure Sensor',     'oil',    25,  15, 5,  7,  27.00, 3),
  (gen_random_uuid(), 'Timing Chain Kit',        'engine', 3,   5,  1,  28, 310.00, 4),
  (gen_random_uuid(), 'Oil Pan Gasket',          'oil',    80,  40, 8,  4,  15.00, 2),
  (gen_random_uuid(), 'Camshaft Assembly',       'engine', -2,  6,  1,  45, 780.00, 5)
ON CONFLICT DO NOTHING;

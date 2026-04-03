-- +goose Up
-- +goose StatementBegin
-- Inserir Venues (Locais de eventos)
INSERT INTO venues (id, uuid, created_at, updated_at, location, seat_map) VALUES
(1, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', NOW(), NOW(), 'Madison Square Garden, New York', '{"A1": true, "A2": true, "A3": true, "B1": true, "B2": true, "B3": true, "C1": true, "C2": true, "C3": true}'::json),
(2, 'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e', NOW(), NOW(), 'Staples Center, Los Angeles', '{"VIP1": true, "VIP2": true, "VIP3": true, "A1": true, "A2": true, "A3": true, "B1": true, "B2": true}'::json),
(3, 'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5f', NOW(), NOW(), 'Royal Albert Hall, London', '{"BOX1": true, "BOX2": true, "STALL1": true, "STALL2": true, "STALL3": true, "CIRCLE1": true}'::json),
(4, 'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5a', NOW(), NOW(), 'Estádio Maracanã, Rio de Janeiro', '{"A1": true, "A2": true, "A3": true, "A4": true, "B1": true, "B2": true, "B3": true, "B4": true}'::json);

-- Inserir Performers (Artistas/Performers)
INSERT INTO performers (id, uuid, created_at, updated_at, name, age, description) VALUES
(1, 'e5f6a7b8-c9d0-4e5f-8a9b-0c1d2e3f4a5b', NOW(), NOW(), 'Taylor Swift', 34, 'Multi-award winning pop and country music artist known for storytelling in her songs'),
(2, 'f6a7b8c9-d0e1-4f5a-8b9c-0d1e2f3a4b5c', NOW(), NOW(), 'Coldplay', 25, 'British rock band formed in London, known for atmospheric rock and emotional ballads'),
(3, 'a7b8c9d0-e1f2-4a5b-8c9d-0e1f2a3b4c5d', NOW(), NOW(), 'Ed Sheeran', 33, 'English singer-songwriter known for his acoustic pop and heartfelt lyrics'),
(4, 'b8c9d0e1-f2a3-4b5c-8d9e-0f1a2b3c4d5e', NOW(), NOW(), 'Beyoncé', 42, 'Iconic R&B and pop artist, actress, and cultural phenomenon'),
(5, 'c9d0e1f2-a3b4-4c5d-8e9f-0a1b2c3d4e5f', NOW(), NOW(), 'The Weeknd', 34, 'Canadian singer known for his distinctive voice and R&B/pop fusion');

-- Inserir Events (Eventos)
INSERT INTO events (id, uuid, created_at, updated_at, venue_id, date, name, description) VALUES
(1, 'd0e1f2a3-b4c5-4d5e-8f9a-0b1c2d3e4f5a', NOW(), NOW(), 1, '2026-07-15 20:00:00', 'The Eras Tour - New York', 'Taylor Swift brings her iconic Eras Tour to Madison Square Garden for an unforgettable night'),
(2, 'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b', NOW(), NOW(), 2, '2026-08-22 19:30:00', 'Music of the Spheres World Tour', 'Coldplay performs their greatest hits with stunning visual effects'),
(3, 'f2a3b4c5-d6e7-4f5a-8b9c-0d1e2f3a4b5c', NOW(), NOW(), 3, '2026-06-10 21:00:00', 'Ed Sheeran - Mathematics Tour', 'An intimate evening with Ed Sheeran at the Royal Albert Hall'),
(4, 'a3b4c5d6-e7f8-4a5b-8c9d-0e1f2a3b4c5d', NOW(), NOW(), 4, '2026-09-05 18:00:00', 'Rock in Rio 2026', 'The biggest music festival in Brazil featuring multiple top artists'),
(5, 'b4c5d6e7-f8a9-4b5c-8d9e-0f1a2b3c4d5e', NOW(), NOW(), 2, '2026-05-20 20:30:00', 'After Hours til Dawn Tour', 'The Weeknd performs tracks from his latest albums');

-- Inserir relação Many-to-Many entre Events e Performers
INSERT INTO event_performers (event_id, performer_id) VALUES
(1, 1),  -- Taylor Swift no The Eras Tour
(2, 2),  -- Coldplay no Music of the Spheres
(3, 3),  -- Ed Sheeran no Mathematics Tour
(4, 2),  -- Coldplay no Rock in Rio
(4, 4),  -- Beyoncé no Rock in Rio
(4, 5),  -- The Weeknd no Rock in Rio
(5, 5);  -- The Weeknd no After Hours til Dawn

-- Inserir Tickets (Ingressos)
INSERT INTO tickets (id, uuid, created_at, updated_at, event_id, price, seat, status) VALUES
-- Tickets para The Eras Tour
(1, 'c5d6e7f8-a9b0-4c5d-8e9f-0a1b2c3d4e5f', NOW(), NOW(), 1, 35000, 'A1', 'available'),
(2, 'd6e7f8a9-b0c1-4d5e-8f9a-0b1c2d3e4f5a', NOW(), NOW(), 1, 35000, 'A2', 'booked'),
(3, 'e7f8a9b0-c1d2-4e5f-8a9b-0c1d2e3f4a5b', NOW(), NOW(), 1, 35000, 'A3', 'available'),
(4, 'f8a9b0c1-d2e3-4f5a-8b9c-0d1e2f3a4b5c', NOW(), NOW(), 1, 30000, 'B1', 'available'),
(5, 'a9b0c1d2-e3f4-4a5b-8c9d-0e1f2a3b4c5d', NOW(), NOW(), 1, 30000, 'B2', 'booked'),
-- Tickets para Coldplay - Music of the Spheres
(6, 'b0c1d2e3-f4a5-4b5c-8d9e-0f1a2b3c4d5e', NOW(), NOW(), 2, 45000, 'VIP1', 'booked'),
(7, 'c1d2e3f4-a5b6-4c5d-8e9f-0a1b2c3d4e5f', NOW(), NOW(), 2, 45000, 'VIP2', 'available'),
(8, 'd2e3f4a5-b6c7-4d5e-8f9a-0b1c2d3e4f5a', NOW(), NOW(), 2, 28000, 'A1', 'available'),
(9, 'e3f4a5b6-c7d8-4e5f-8a9b-0c1d2e3f4a5b', NOW(), NOW(), 2, 28000, 'A2', 'available'),
-- Tickets para Ed Sheeran
(10, 'f4a5b6c7-d8e9-4f5a-8b9c-0d1e2f3a4b5c', NOW(), NOW(), 3, 25000, 'BOX1', 'available'),
(11, 'a5b6c7d8-e9f0-4a5b-8c9d-0e1f2a3b4c5d', NOW(), NOW(), 3, 25000, 'BOX2', 'booked'),
(12, 'b6c7d8e9-f0a1-4b5c-8d9e-0f1a2b3c4d5e', NOW(), NOW(), 3, 18000, 'STALL1', 'available'),
(13, 'c7d8e9f0-a1b2-4c5d-8e9f-0a1b2c3d4e5f', NOW(), NOW(), 3, 18000, 'STALL2', 'available'),
-- Tickets para Rock in Rio
(14, 'd8e9f0a1-b2c3-4d5e-8f9a-0b1c2d3e4f5a', NOW(), NOW(), 4, 15000, 'A1', 'available'),
(15, 'e9f0a1b2-c3d4-4e5f-8a9b-0c1d2e3f4a5b', NOW(), NOW(), 4, 15000, 'A2', 'booked'),
(16, 'f0a1b2c3-d4e5-4f5a-8b9c-0d1e2f3a4b5c', NOW(), NOW(), 4, 15000, 'A3', 'available'),
(17, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5e', NOW(), NOW(), 4, 12000, 'B1', 'available'),
-- Tickets para The Weeknd
(18, 'b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5f', NOW(), NOW(), 5, 32000, 'VIP1', 'booked'),
(19, 'c3d4e5f6-a7b8-4c5d-8e9f-0a1b2c3d4e5a', NOW(), NOW(), 5, 32000, 'VIP2', 'available'),
(20, 'd4e5f6a7-b8c9-4d5e-8f9a-0b1c2d3e4f5b', NOW(), NOW(), 5, 26000, 'A1', 'available');

-- Atualizar sequências para evitar conflitos de ID no futuro
SELECT setval('venues_id_seq', (SELECT MAX(id) FROM venues));
SELECT setval('performers_id_seq', (SELECT MAX(id) FROM performers));
SELECT setval('events_id_seq', (SELECT MAX(id) FROM events));
SELECT setval('tickets_id_seq', (SELECT MAX(id) FROM tickets));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tickets WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20);
DELETE FROM event_performers WHERE event_id IN (1, 2, 3, 4, 5);
DELETE FROM events WHERE id IN (1, 2, 3, 4, 5);
DELETE FROM performers WHERE id IN (1, 2, 3, 4, 5);
DELETE FROM venues WHERE id IN (1, 2, 3, 4);
-- +goose StatementEnd

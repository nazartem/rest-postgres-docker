--DROP TABLE IF EXISTS product CASCADE;
--DROP TABLE IF EXISTS buyer CASCADE;
--DROP TABLE IF EXISTS product_list CASCADE;
--DROP TABLE IF EXISTS note CASCADE;

CREATE TABLE IF NOT EXISTS public.product
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(100) NOT NULL,
    price DECIMAL DEFAULT 0.00,
    amount INT,

    UNIQUE (name)
);

CREATE TABLE public.buyer
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL
);

CREATE TABLE public.note
(
    number SERIAL PRIMARY KEY,
    date TIMESTAMP,
    buyer_id INT,

    CONSTRAINT buyer_fk FOREIGN KEY (buyer_id) REFERENCES public.buyer (id)
);

CREATE TABLE public.product_list
(
    id   SERIAL PRIMARY KEY,
    note_id INT,
    product_id INT,
    amount INT,

    CONSTRAINT note_id_fk FOREIGN KEY (note_id) REFERENCES public.note (number),
    CONSTRAINT product_id_fk FOREIGN KEY (product_id) REFERENCES public.product (id)
);

-- product
INSERT INTO product (name, description, price, amount)
VALUES ('Колбаса', 'some description', 254.9, 50);
INSERT INTO product (name, description, price, amount)
VALUES ('Сыр', 'some description', 213.9, 21);
INSERT INTO product (name, description, price, amount)
VALUES ('Молоко', 'some description', 61.3, 30);

-- buyer
INSERT INTO buyer (name, surname)
VALUES ('Билли', 'Харингтонов');
INSERT INTO buyer (name, surname)
VALUES ('Гарри', 'Поттер');
INSERT INTO buyer (name, surname)
VALUES ('Рон', 'Уизли');

-- note
INSERT INTO note (date, buyer_id)
VALUES ('2022-03-25T11:11:00Z', 1);
INSERT INTO note (date, buyer_id)
VALUES ('2022-03-27T13:15:00Z', 2);
INSERT INTO note (date, buyer_id)
VALUES ('2022-03-27T16:13:00Z', 3);


-- product_list
INSERT INTO product_list (note_id, product_id, amount)
VALUES (1, 2, 10);
INSERT INTO product_list (note_id, product_id, amount)
VALUES (1, 1, 15);
INSERT INTO product_list (note_id, product_id, amount)
VALUES (2, 3, 50);
INSERT INTO product_list (note_id, product_id, amount)
VALUES (3, 3, 150);

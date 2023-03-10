CREATE TABLE users
(
    id bigserial not null unique,
    uname varchar(250),
    balance bigint,
    reserved bigint
);

CREATE TABLE services
(
    id int not null unique,
    price bigint
);

CREATE TABLE orders (
    user_id bigint references users(id) unique,
    services_id int[],
    amount bigint not null
);

CREATE TABLE accounting
(
    usersid int references users(id) on delete cascade,
    servicesid int[],
    ac_date TIMESTAMP DEFAULT now() 
);

CREATE PROCEDURE transaction_p2p(srsid INT, dstid INT, amount INT)
LANGUAGE plpgsql
AS $$
BEGIN
  IF (SELECT balance FROM users WHERE id = srsid) >= amount THEN
    UPDATE users SET balance = balance - amount, reserved = reserved + amount WHERE id = srsid;
    UPDATE users SET balance = balance + amount WHERE id = dstid;
    UPDATE users SET reserved = reserved - amount WHERE id = srsid;
    COMMIT;
  ELSE
    ROLLBACK;
  END IF;
END;
$$;

CREATE PROCEDURE make_order(s_id INT, u_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
  IF (SELECT balance FROM users WHERE id = u_id) >= (SELECT price FROM services WHERE id = s_id) THEN
    IF EXISTS (SELECT * FROM orders WHERE user_id = u_id) THEN
      UPDATE orders SET services_id = array_append(services_id, s_id), amount = amount + (SELECT price FROM services WHERE id = s_id) WHERE user_id = u_id;
    ELSE
      INSERT INTO orders (services_id, user_id, amount) VALUES (ARRAY[s_id], u_id, (SELECT price FROM services WHERE id = s_id));
    END IF;
    UPDATE users SET balance = balance - (SELECT price FROM services WHERE id = s_id), reserved = reserved + (SELECT price FROM services WHERE id = s_id) WHERE id = u_id;

    COMMIT;
  ELSE
    ROLLBACK;
  END IF;
END;
$$;

CREATE PROCEDURE do_order(u_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO accounting (usersid,servicesid) VALUES (u_id, (SELECT services_id FROM orders WHERE user_id = u_id));
  UPDATE users SET reserved = reserved - (SELECT amount FROM orders WHERE user_id = u_id) WHERE id = u_id;
  DELETE FROM orders WHERE user_id = u_id;
COMMIT;
END;
$$;




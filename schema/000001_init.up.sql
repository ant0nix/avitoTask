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
    id bigserial unique,
    services_id int references services(id) on delete cascade,
    user_id int references users(id) on delete cascade,
    amount bigint not null
);                        
  
CREATE TABLE company (
    order_id bigint references orders(id) on delete cascade,
    balance bigint
); 

CREATE TABLE accounting
(
    usersid int references users(id) on delete cascade,
    servicesid int references services(id) on delete cascade
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

CREATE PROCEDURE make_order (s_id INT, u_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
  IF (SELECT balance FROM users WHERE id = u_id) >= (SELECT price FROM services WHERE id = s_id) THEN
    UPDATE users SET balance = balance - (SELECT price FROM services WHERE id = s_id), reserved = reserved + (SELECT price FROM services WHERE id = s_id) WHERE id = u_id;
    INSERT INTO orders (services_id,user_id,amount) VALUES (s_id,u_id,(SELECT price FROM services WHERE id = s_id));
    COMMIT;
  ELSE
    ROLLBACK;
  END IF;
END;
$$;
create table items(
	item_id serial primary key not null,
	item_code char(3) not null,
	description varchar(255) not null,
	quantity int not null,
	order_id int not null,
	CONSTRAINT fk_items
		FOREIGN KEY(order_id)
			REFERENCES orders(order_id)
);

create table orders(
	order_id serial primary key not null,
	customer_name varchar(255) not null,
	ordered_at timestamp default now() not null
);

select * from orders;
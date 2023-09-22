CREATE TABLE people (
	uid varchar(36) PRIMARY KEY,
	name varchar(50) NOT NULL,
	surname varchar(50) NOT NULL,
	patronymic varchar(50),
	age integer,
	gender varchar(6),
	nation varchar(3),
	error TEXT
);
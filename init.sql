CREATE DATABASE grocery;

\c grocery;

CREATE TABLE grocery_items (
    name varchar(50),
    unit varchar(20),
    cost Numeric(6, 2),
    picture varchar(255),
    PRIMARY KEY (name, unit)
);

CREATE TABLE grocery_list (
    name varchar(50),
    unit varchar(20),
    cost Numeric(6, 2),
    number Decimal,
    PRIMARY KEY (name, unit),
    FOREIGN KEY (name, unit) REFERENCES grocery_items(name, unit)
);

CREATE TABLE recipes (
    name varchar(50) PRIMARY KEY,
    directions text,
    picture varchar(255)
);

CREATE TABLE recipe_ingredients (
    grocery_item varchar(50),
    recipe varchar(50),
    number Decimal,
    unit varchar(20),
    PRIMARY KEY (grocery_item, unit, recipe),
    FOREIGN KEY (grocery_item, unit) REFERENCES grocery_items(name, unit),
    FOREIGN KEY (recipe) REFERENCES recipes(name)
);
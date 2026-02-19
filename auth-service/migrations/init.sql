CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    emp_type VARCHAR(50) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_employees_username ON employees(username);